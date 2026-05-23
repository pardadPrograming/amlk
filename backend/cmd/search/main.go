package main

import (
	"context"
	"encoding/json"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"amlakcrm/backend/internal/config"
	"amlakcrm/backend/internal/events"
	"amlakcrm/backend/internal/repository"
	"amlakcrm/backend/internal/searchengine"
	"amlakcrm/backend/internal/searchrpc"

	amqp "github.com/rabbitmq/amqp091-go"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
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
	matcher := searchengine.NewOptimizedMatchingService(store, cfg.SearchIndexTTL)
	invalidationConn := startInvalidationConsumer(cfg, logger, matcher)
	if invalidationConn != nil {
		defer invalidationConn.Close()
	}

	listener, err := net.Listen("tcp", cfg.SearchAddr)
	if err != nil {
		logger.Error("search grpc listen failed", "error", err)
		os.Exit(1)
	}
	server := grpc.NewServer()
	searchrpc.RegisterSearchServiceServer(server, &searchServer{
		matcher: matcher,
		token:   cfg.SearchServiceToken,
	})

	go func() {
		logger.Info("search grpc service listening", "addr", cfg.SearchAddr)
		if err := server.Serve(listener); err != nil {
			logger.Error("search grpc service failed", "error", err)
			os.Exit(1)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	server.GracefulStop()
}

type searchServer struct {
	matcher *searchengine.OptimizedMatchingService
	token   string
}

func (s *searchServer) RequestMatches(ctx context.Context, req *searchrpc.RequestMatchesRequest) (*searchrpc.RequestMatchesResponse, error) {
	if !authorized(ctx, s.token) {
		return nil, status.Error(codes.Unauthenticated, "invalid search service token")
	}
	page, err := s.matcher.RequestMatches(ctx, req.UserID, req.BusinessID, req.ContactID, req.RequestID, int(req.Limit), int(req.Offset))
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	return &searchrpc.RequestMatchesResponse{
		Matches: page.Items,
		Total:   int32(page.Total),
		Limit:   int32(page.Limit),
		Offset:  int32(page.Offset),
	}, nil
}

func newPrimaryStore(cfg config.Config, logger *slog.Logger) (repository.Store, *mongo.Client) {
	memoryStore := repository.NewMemoryStore()
	if cfg.MongoURI == "" {
		logger.Info("search store using in-memory adapter; MONGO_URI is empty")
		return memoryStore, nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.MongoURI))
	if err != nil {
		logger.Warn("mongo connect failed; search store falling back to in-memory adapter", "error", err)
		return memoryStore, nil
	}
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		_ = client.Disconnect(ctx)
		logger.Warn("mongo ping failed; search store falling back to in-memory adapter", "error", err)
		return memoryStore, nil
	}
	store := repository.NewMongoPlatformStore(memoryStore, client.Database(cfg.MongoDatabase))
	if err := store.EnsurePlatformIndexes(ctx); err != nil {
		_ = client.Disconnect(ctx)
		logger.Warn("mongo index setup failed; search store falling back to in-memory adapter", "error", err)
		return memoryStore, nil
	}
	logger.Info("search store using mongo", "database", cfg.MongoDatabase)
	return store, client
}

func authorized(ctx context.Context, token string) bool {
	token = strings.TrimSpace(token)
	if token == "" {
		return true
	}
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return false
	}
	for _, value := range md.Get("authorization") {
		if strings.TrimSpace(value) == "Bearer "+token {
			return true
		}
	}
	return false
}

func startInvalidationConsumer(cfg config.Config, logger *slog.Logger, matcher *searchengine.OptimizedMatchingService) *amqp.Connection {
	if cfg.RabbitMQURL == "" {
		logger.Info("search cache invalidation disabled; RABBITMQ_URL is empty")
		return nil
	}
	conn, err := amqp.Dial(cfg.RabbitMQURL)
	if err != nil {
		logger.Warn("rabbitmq unavailable; search cache invalidation disabled", "error", err)
		return nil
	}
	ch, err := conn.Channel()
	if err != nil {
		_ = conn.Close()
		logger.Warn("rabbitmq channel unavailable; search cache invalidation disabled", "error", err)
		return nil
	}
	if err := ch.ExchangeDeclare(cfg.EventExchange, "topic", true, false, false, false, nil); err != nil {
		_ = ch.Close()
		_ = conn.Close()
		logger.Warn("rabbitmq exchange unavailable; search cache invalidation disabled", "error", err)
		return nil
	}
	host, _ := os.Hostname()
	queueName := "amlak.search.invalidate." + host + "." + strconv.Itoa(os.Getpid())
	queue, err := ch.QueueDeclare(queueName, false, true, true, false, nil)
	if err != nil {
		_ = ch.Close()
		_ = conn.Close()
		logger.Warn("rabbitmq queue unavailable; search cache invalidation disabled", "error", err)
		return nil
	}
	if err := ch.QueueBind(queue.Name, events.SearchInvalidateEvent, cfg.EventExchange, false, nil); err != nil {
		_ = ch.Close()
		_ = conn.Close()
		logger.Warn("rabbitmq queue bind failed; search cache invalidation disabled", "error", err)
		return nil
	}
	deliveries, err := ch.Consume(queue.Name, "", false, false, false, false, nil)
	if err != nil {
		_ = ch.Close()
		_ = conn.Close()
		logger.Warn("rabbitmq consume failed; search cache invalidation disabled", "error", err)
		return nil
	}
	go func() {
		defer ch.Close()
		for delivery := range deliveries {
			var event struct {
				Type    string                         `json:"type"`
				Payload events.SearchInvalidatePayload `json:"payload"`
			}
			if err := json.Unmarshal(delivery.Body, &event); err != nil {
				_ = delivery.Nack(false, false)
				continue
			}
			if event.Type == events.SearchInvalidateEvent && event.Payload.BusinessID != "" && event.Payload.UserID != "" {
				matcher.Invalidate(event.Payload.BusinessID, event.Payload.UserID)
			}
			_ = delivery.Ack(false)
		}
	}()
	logger.Info("search cache invalidation enabled", "exchange", cfg.EventExchange, "queue", queue.Name)
	return conn
}
