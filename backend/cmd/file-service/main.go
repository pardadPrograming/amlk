package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"amlakcrm/backend/internal/config"
	"amlakcrm/backend/internal/repository"
	"amlakcrm/backend/internal/service"
	httptransport "amlakcrm/backend/internal/transport/http"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func main() {
	cfg := config.Load()
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))

	store, mongoClient := newPrimaryStore(cfg, logger)
	if mongoClient != nil {
		defer func() {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			if err := mongoClient.Disconnect(ctx); err != nil {
				logger.Warn("mongo disconnect failed", "error", err)
			}
		}()
	}

	tokens := service.NewTokenService(cfg.JWTSecret, cfg.AccessTokenTTL, cfg.RefreshTokenTTL)
	uploadSvc := service.NewUploadService(store, cfg.ObjectStorageDir)
	uploadSvc.StartExpiredFileCleanup(context.Background(), time.Hour)
	router := httptransport.NewUploadRouter(httptransport.RouterDeps{
		Config:  cfg,
		Logger:  logger,
		Store:   store,
		Tokens:  tokens,
		Uploads: uploadSvc,
	})

	server := &http.Server{
		Addr:              cfg.FileServiceAddr,
		Handler:           router,
		ReadHeaderTimeout: 10 * time.Second,
	}

	go func() {
		logger.Info("file service listening", "addr", cfg.FileServiceAddr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("file service failed", "error", err)
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
		logger.Info("file service store using in-memory adapter; MONGO_URI is empty")
		return memoryStore, nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.MongoURI))
	if err != nil {
		logger.Warn("mongo connect failed; file service falling back to in-memory adapter", "error", err)
		return memoryStore, nil
	}
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		_ = client.Disconnect(ctx)
		logger.Warn("mongo ping failed; file service falling back to in-memory adapter", "error", err)
		return memoryStore, nil
	}
	store := repository.NewMongoPlatformStore(memoryStore, client.Database(cfg.MongoDatabase))
	if err := store.EnsurePlatformIndexes(ctx); err != nil {
		_ = client.Disconnect(ctx)
		logger.Warn("mongo index setup failed; file service falling back to in-memory adapter", "error", err)
		return memoryStore, nil
	}
	logger.Info("file service store using mongo", "database", cfg.MongoDatabase)
	return store, client
}
