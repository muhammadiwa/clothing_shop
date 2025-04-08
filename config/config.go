package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

// Config holds all configuration for the application
type Config struct {
	// Server
	ServerPort  string
	Environment string
	BaseURL     string

	// Database
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string

	// Redis
	RedisHost     string
	RedisPort     string
	RedisPassword string
	RedisDB       int

	// JWT
	JWTSecret             string
	JWTExpirationHours    int
	RefreshExpirationDays int

	// Email
	SMTPHost     string
	SMTPPort     string
	SMTPUser     string
	SMTPPassword string
	SMTPFrom     string

	// Third-party APIs
	RajaOngkirAPIKey  string
	RajaOngkirBaseURL string
	MidtransServerKey string
	MidtransClientKey string
	MidtransBaseURL   string

	// Rate Limiting
	RateLimitRequests int
	RateLimitDuration time.Duration
}

// LoadConfig loads configuration from environment variables
func LoadConfig() (*Config, error) {
	redisDB, _ := strconv.Atoi(getEnv("REDIS_DB", "0"))
	jwtExpHours, _ := strconv.Atoi(getEnv("JWT_EXPIRATION_HOURS", "24"))
	refreshExpDays, _ := strconv.Atoi(getEnv("REFRESH_EXPIRATION_DAYS", "7"))
	rateLimitReq, _ := strconv.Atoi(getEnv("RATE_LIMIT_REQUESTS", "100"))
	rateLimitDur, _ := time.ParseDuration(getEnv("RATE_LIMIT_DURATION", "1h"))

	config := &Config{
		// Server
		ServerPort:  getEnv("SERVER_PORT", "8080"),
		Environment: getEnv("ENVIRONMENT", "development"),
		BaseURL:     getEnv("BASE_URL", "http://localhost:8080"),

		// Database
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", "postgres"),
		DBName:     getEnv("DB_NAME", "fashion_shop"),
		DBSSLMode:  getEnv("DB_SSL_MODE", "disable"),

		// Redis
		RedisHost:     getEnv("REDIS_HOST", "localhost"),
		RedisPort:     getEnv("REDIS_PORT", "6379"),
		RedisPassword: getEnv("REDIS_PASSWORD", ""),
		RedisDB:       redisDB,

		// JWT
		JWTSecret:             getEnv("JWT_SECRET", "your-secret-key"),
		JWTExpirationHours:    jwtExpHours,
		RefreshExpirationDays: refreshExpDays,

		// Email
		SMTPHost:     getEnv("SMTP_HOST", ""),
		SMTPPort:     getEnv("SMTP_PORT", "587"),
		SMTPUser:     getEnv("SMTP_USER", ""),
		SMTPPassword: getEnv("SMTP_PASSWORD", ""),
		SMTPFrom:     getEnv("SMTP_FROM", ""),

		// Third-party APIs
		RajaOngkirAPIKey:  getEnv("RAJAONGKIR_API_KEY", ""),
		RajaOngkirBaseURL: getEnv("RAJAONGKIR_BASE_URL", "https://api.rajaongkir.com/starter"),
		MidtransServerKey: getEnv("MIDTRANS_SERVER_KEY", ""),
		MidtransClientKey: getEnv("MIDTRANS_CLIENT_KEY", ""),
		MidtransBaseURL:   getEnv("MIDTRANS_BASE_URL", "https://api.sandbox.midtrans.com"),

		// Rate Limiting
		RateLimitRequests: rateLimitReq,
		RateLimitDuration: rateLimitDur,
	}

	return config, nil
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// GetDSN returns the database connection string
func (c *Config) GetDSN() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.DBHost, c.DBPort, c.DBUser, c.DBPassword, c.DBName, c.DBSSLMode)
}

// GetRedisAddr returns the Redis connection string
func (c *Config) GetRedisAddr() string {
	return fmt.Sprintf("%s:%s", c.RedisHost, c.RedisPort)
}
