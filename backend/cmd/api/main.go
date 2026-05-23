package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"amlakcrm/backend/internal/cache"
	"amlakcrm/backend/internal/config"
	"amlakcrm/backend/internal/domain"
	"amlakcrm/backend/internal/events"
	"amlakcrm/backend/internal/repository"
	"amlakcrm/backend/internal/service"
	httptransport "amlakcrm/backend/internal/transport/http"

	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func main() {
	cfg := config.Load()
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))

	primaryStore, mongoClient := newPrimaryStore(cfg, logger)
	if mongoClient != nil {
		defer func() {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			if err := mongoClient.Disconnect(ctx); err != nil {
				logger.Warn("mongo disconnect failed", "error", err)
			}
		}()
	}
	publisher := newEventPublisher(cfg, logger)
	defer func() {
		if err := publisher.Close(); err != nil {
			logger.Warn("event publisher close failed", "error", err)
		}
	}()
	store := repository.NewCachedStore(primaryStore, newCachedStoreOptions(cfg, logger, publisher))
	tokens := service.NewTokenService(cfg.JWTSecret, cfg.AccessTokenTTL, cfg.RefreshTokenTTL)
	authSvc := service.NewAuthService(store, tokens, cfg)
	businessSvc := service.NewBusinessService(store)
	invitationSvc := service.NewInvitationService(store)
	fileSvc := service.NewFileService(store)
	uploadSvc := service.NewUploadService(store, cfg.ObjectStorageDir)
	uploadSvc.StartExpiredFileCleanup(context.Background(), time.Hour)
	locationSvc := service.NewLocationService(store)
	propertySvc := service.NewPropertyService(store, cfg.ObjectStorageDir)
	contactSvc := service.NewContactService(store)
	channelSvc := service.NewChannelService(store, cfg.ObjectStorageDir, cfg.MediaOptimizerURL)
	channelSvc.StartMediaRetentionCleanup(context.Background(), 24*time.Hour)
	matchingSvc := service.NewSearchMatchingClient(cfg.SearchServiceAddr, cfg.SearchServiceToken)
	platformSvc := service.NewPlatformService(store, cfg)

	router := httptransport.NewRouter(httptransport.RouterDeps{
		Config:     cfg,
		Logger:     logger,
		Store:      store,
		Tokens:     tokens,
		Auth:       authSvc,
		Businesses: businessSvc,
		Invites:    invitationSvc,
		Files:      fileSvc,
		Uploads:    uploadSvc,
		Locations:  locationSvc,
		Properties: propertySvc,
		Contacts:   contactSvc,
		Channels:   channelSvc,
		Matching:   matchingSvc,
		Platform:   platformSvc,
	})

	server := &http.Server{
		Addr:              cfg.HTTPAddr,
		Handler:           router,
		ReadHeaderTimeout: 10 * time.Second,
	}

	go func() {
		logger.Info("api listening", "addr", cfg.HTTPAddr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("server failed", "error", err)
			os.Exit(1)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_ = server.Shutdown(ctx)
}

func newPrimaryStore(cfg config.Config, logger *slog.Logger) (repository.Store, *mongo.Client) {
	memoryStore := repository.NewMemoryStore()
	if cfg.MongoURI == "" {
		logger.Info("platform catalog store using in-memory adapter; MONGO_URI is empty")
		return memoryStore, nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.MongoURI))
	if err != nil {
		logger.Warn("mongo connect failed; platform catalog falling back to in-memory adapter", "error", err)
		return memoryStore, nil
	}
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		_ = client.Disconnect(ctx)
		logger.Warn("mongo ping failed; platform catalog falling back to in-memory adapter", "error", err)
		return memoryStore, nil
	}
	store := repository.NewMongoPlatformStore(memoryStore, client.Database(cfg.MongoDatabase))
	if err := store.EnsurePlatformIndexes(ctx); err != nil {
		_ = client.Disconnect(ctx)
		logger.Warn("mongo index setup failed; platform catalog falling back to in-memory adapter", "error", err)
		return memoryStore, nil
	}
	logger.Info("platform catalog store using mongo", "database", cfg.MongoDatabase)
	return store, client
}

func newCachedStoreOptions(cfg config.Config, logger *slog.Logger, publisher events.Publisher) repository.CachedStoreOptions {
	opts := repository.CachedStoreOptions{
		TTL:       cfg.IdentityCacheTTL,
		Logger:    logger,
		Publisher: publisher,
	}
	if cfg.RedisAddr == "" {
		logger.Info("identity cache using in-memory ttl cache; REDIS_ADDR is empty", "ttl", cfg.IdentityCacheTTL)
		return opts
	}

	client := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	})
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := client.Ping(ctx).Err(); err != nil {
		_ = client.Close()
		logger.Warn("redis unavailable; identity cache falling back to in-memory ttl cache", "addr", cfg.RedisAddr, "error", err)
		return opts
	}

	logger.Info("identity cache using redis", "addr", cfg.RedisAddr, "ttl", cfg.IdentityCacheTTL)
	opts.UsersByIDCache = cache.NewRedisCache[domain.User](client, "amlak:identity")
	opts.UserIDByPhoneCache = cache.NewRedisCache[string](client, "amlak:identity")
	opts.SessionsByIDCache = cache.NewRedisCache[domain.Session](client, "amlak:identity")
	opts.SessionByTokenCache = cache.NewRedisCache[domain.Session](client, "amlak:identity")
	opts.RefreshByIDCache = cache.NewRedisCache[string](client, "amlak:identity")
	return opts
}

func newEventPublisher(cfg config.Config, logger *slog.Logger) events.Publisher {
	if cfg.RabbitMQURL == "" {
		logger.Info("event publisher disabled; RABBITMQ_URL is empty")
		return events.NoopPublisher{}
	}
	publisher, err := events.NewRabbitMQPublisher(cfg.RabbitMQURL, cfg.EventExchange, logger)
	if err != nil {
		logger.Warn("rabbitmq unavailable; falling back to noop publisher", "error", err)
		return events.NoopPublisher{}
	}
	logger.Info("rabbitmq event publisher enabled", "exchange", cfg.EventExchange)
	return publisher
}
