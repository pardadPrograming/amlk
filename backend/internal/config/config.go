package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	AppEnv              string
	HTTPAddr            string
	SearchAddr          string
	SearchServiceAddr   string
	SearchServiceToken  string
	MessagingAddr       string
	MessagingServiceURL string
	FilingAddr          string
	FilingServiceURL    string
	FileServiceAddr     string
	FileServiceURL      string
	MediaOptimizerAddr  string
	MediaOptimizerURL   string
	JWTSecret           string
	MongoURI            string
	MongoDatabase       string
	RedisAddr           string
	RedisPassword       string
	RedisDB             int
	RabbitMQURL         string
	EventExchange       string
	ObjectStorageDir    string
	SuperAdminPhones    string
	IdentityCacheTTL    time.Duration
	SearchIndexTTL      time.Duration
	AccessTokenTTL      time.Duration
	RefreshTokenTTL     time.Duration
	OTPTTL              time.Duration
	OTPResendAfter      time.Duration
	OTPMaxAttempts      int
}

func Load() Config {
	return Config{
		AppEnv:              env("APP_ENV", "development"),
		HTTPAddr:            env("HTTP_ADDR", ":8080"),
		SearchAddr:          env("SEARCH_ADDR", ":8091"),
		SearchServiceAddr:   env("SEARCH_SERVICE_ADDR", "127.0.0.1:8091"),
		SearchServiceToken:  env("SEARCH_SERVICE_TOKEN", "dev-search-token"),
		MessagingAddr:       env("MESSAGING_ADDR", ":8092"),
		MessagingServiceURL: env("MESSAGING_SERVICE_URL", ""),
		FilingAddr:          env("FILING_ADDR", ":8093"),
		FilingServiceURL:    env("FILING_SERVICE_URL", ""),
		FileServiceAddr:     env("FILE_SERVICE_ADDR", ":8095"),
		FileServiceURL:      env("FILE_SERVICE_URL", ""),
		MediaOptimizerAddr:  env("MEDIA_OPTIMIZER_ADDR", ":8094"),
		MediaOptimizerURL:   env("MEDIA_OPTIMIZER_URL", ""),
		JWTSecret:           env("JWT_SECRET", "dev-secret-change-me"),
		MongoURI:            env("MONGO_URI", ""),
		MongoDatabase:       env("MONGO_DATABASE", "amlak"),
		RedisAddr:           env("REDIS_ADDR", ""),
		RedisPassword:       env("REDIS_PASSWORD", ""),
		RedisDB:             envInt("REDIS_DB", 0),
		RabbitMQURL:         env("RABBITMQ_URL", ""),
		EventExchange:       env("EVENT_EXCHANGE", "amlak.events"),
		ObjectStorageDir:    env("OBJECT_STORAGE_DIR", "storage/objects"),
		SuperAdminPhones:    env("SUPER_ADMIN_PHONES", ""),
		IdentityCacheTTL:    time.Duration(envInt("IDENTITY_CACHE_TTL_SECONDS", 60)) * time.Second,
		SearchIndexTTL:      time.Duration(envInt("SEARCH_INDEX_TTL_SECONDS", 30)) * time.Second,
		AccessTokenTTL:      time.Duration(envInt("ACCESS_TOKEN_TTL_MINUTES", 30)) * time.Minute,
		RefreshTokenTTL:     time.Duration(envInt("REFRESH_TOKEN_TTL_HOURS", 720)) * time.Hour,
		OTPTTL:              time.Duration(envInt("OTP_TTL_MINUTES", 2)) * time.Minute,
		OTPResendAfter:      time.Duration(envInt("OTP_RESEND_SECONDS", 60)) * time.Second,
		OTPMaxAttempts:      envInt("OTP_MAX_ATTEMPTS", 5),
	}
}

func env(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func envInt(key string, fallback int) int {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	parsed, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}
	return parsed
}
