package config

import (
	"os"
	"strconv"
	"time"
)

// Config holds all configuration for the application
type Config struct {
	Server struct {
		Port int
		Mode string
	}
	Database struct {
		Host     string
		Port     int
		User     string
		Password string
		Name     string
	}
	Redis struct {
		Host     string
		Port     int
		Password string
		DB       int
	}
	JWT struct {
		Secret        string
		AccessExpiry  time.Duration
		RefreshExpiry time.Duration
	}
	RateLimit struct {
		Requests int
		Duration time.Duration
	}
	RajaOngkir struct {
		APIKey string
		URL    string
	}
	Midtrans struct {
		ServerKey      string
		ClientKey      string
		Environment    string
		PaymentWebhook string
	}
	SMTP struct {
		Host     string
		Port     int
		Username string
		Password string
		From     string
	}
	Storage struct {
		Type      string // local, s3, etc.
		LocalPath string
		S3Bucket  string
		S3Region  string
	}
}

// NewConfig creates a new Config instance
func NewConfig() *Config {
	cfg := &Config{}

	// Server configuration
	cfg.Server.Port = getEnvAsInt("SERVER_PORT", 8080)
	cfg.Server.Mode = getEnvAsString("SERVER_MODE", "development")

	// Database configuration
	cfg.Database.Host = getEnvAsString("DB_HOST", "localhost")
	cfg.Database.Port = getEnvAsInt("DB_PORT", 5432)
	cfg.Database.User = getEnvAsString("DB_USER", "postgres")
	cfg.Database.Password = getEnvAsString("DB_PASSWORD", "postgres")
	cfg.Database.Name = getEnvAsString("DB_NAME", "fashion_shop")

	// Redis configuration
	cfg.Redis.Host = getEnvAsString("REDIS_HOST", "localhost")
	cfg.Redis.Port = getEnvAsInt("REDIS_PORT", 6379)
	cfg.Redis.Password = getEnvAsString("REDIS_PASSWORD", "")
	cfg.Redis.DB = getEnvAsInt("REDIS_DB", 0)

	// JWT configuration
	cfg.JWT.Secret = getEnvAsString("JWT_SECRET", "your-secret-key")
	cfg.JWT.AccessExpiry = getEnvAsDuration("JWT_ACCESS_EXPIRY", 15*time.Minute)
	cfg.JWT.RefreshExpiry = getEnvAsDuration("JWT_REFRESH_EXPIRY", 7*24*time.Hour)

	// Rate limiting configuration
	cfg.RateLimit.Requests = getEnvAsInt("RATE_LIMIT_REQUESTS", 100)
	cfg.RateLimit.Duration = getEnvAsDuration("RATE_LIMIT_DURATION", time.Minute)

	// RajaOngkir configuration
	cfg.RajaOngkir.APIKey = getEnvAsString("RAJAONGKIR_API_KEY", "")
	cfg.RajaOngkir.URL = getEnvAsString("RAJAONGKIR_URL", "https://api.rajaongkir.com/starter")

	// Midtrans configuration
	cfg.Midtrans.ServerKey = getEnvAsString("MIDTRANS_SERVER_KEY", "")
	cfg.Midtrans.ClientKey = getEnvAsString("MIDTRANS_CLIENT_KEY", "")
	cfg.Midtrans.Environment = getEnvAsString("MIDTRANS_ENVIRONMENT", "sandbox")
	cfg.Midtrans.PaymentWebhook = getEnvAsString("MIDTRANS_PAYMENT_WEBHOOK", "/api/v1/payments/webhook")

	// SMTP configuration
	cfg.SMTP.Host = getEnvAsString("SMTP_HOST", "")
	cfg.SMTP.Port = getEnvAsInt("SMTP_PORT", 587)
	cfg.SMTP.Username = getEnvAsString("SMTP_USERNAME", "")
	cfg.SMTP.Password = getEnvAsString("SMTP_PASSWORD", "")
	cfg.SMTP.From = getEnvAsString("SMTP_FROM", "noreply@fashionshop.com")

	// Storage configuration
	cfg.Storage.Type = getEnvAsString("STORAGE_TYPE", "local")
	cfg.Storage.LocalPath = getEnvAsString("STORAGE_LOCAL_PATH", "./uploads")
	cfg.Storage.S3Bucket = getEnvAsString("STORAGE_S3_BUCKET", "")
	cfg.Storage.S3Region = getEnvAsString("STORAGE_S3_REGION", "")

	return cfg
}

// Helper functions to get environment variables
func getEnvAsString(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvAsDuration(key string, defaultValue time.Duration) time.Duration {
	if value, exists := os.LookupEnv(key); exists {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}
