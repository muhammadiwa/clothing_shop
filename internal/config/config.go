package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	// Server
	ServerPort string
	AppURL     string

	// Database
	DatabaseURL  string
	DatabaseName string

	// JWT
	JWTSecret      string
	JWTExpireHours int

	// Email
	SMTPHost      string
	SMTPPort      int
	SMTPUsername  string
	SMTPPassword  string
	SMTPFromEmail string
	SMTPFromName  string

	// RajaOngkir
	RajaOngkirAPIKey  string
	RajaOngkirBaseURL string

	// Midtrans
	MidtransServerKey   string
	MidtransClientKey   string
	MidtransEnvironment string

	// File Storage
	StorageDriver string
	StoragePath   string

	// Debug and Environment
	Debug      bool
	Production bool

	// Admin
	AdminEmail    string
	AdminPassword string
}

var config *Config

func LoadConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found, using environment variables")
	}

	debug, _ := strconv.ParseBool(getEnv("DEBUG", "false"))
	production, _ := strconv.ParseBool(getEnv("PRODUCTION", "false"))
	smtpPort, _ := strconv.Atoi(getEnv("SMTP_PORT", "587"))
	jwtExpireHours, _ := strconv.Atoi(getEnv("JWT_EXPIRE_HOURS", "24"))

	config = &Config{
		// Server
		ServerPort: getEnv("PORT", "8080"),
		AppURL:     getEnv("APP_URL", "http://localhost:8080"),

		// Database
		DatabaseURL:  getEnv("DATABASE_URL", "mysql://username:password@tcp(localhost:3306)/clothing_shop"),
		DatabaseName: getEnv("DATABASE_NAME", "clothing_shop"),

		// JWT
		JWTSecret:      getEnv("JWT_SECRET", "your_jwt_secret"),
		JWTExpireHours: jwtExpireHours,

		// Email
		SMTPHost:      getEnv("SMTP_HOST", "smtp.gmail.com"),
		SMTPPort:      smtpPort,
		SMTPUsername:  getEnv("SMTP_USERNAME", ""),
		SMTPPassword:  getEnv("SMTP_PASSWORD", ""),
		SMTPFromEmail: getEnv("SMTP_FROM_EMAIL", "noreply@example.com"),
		SMTPFromName:  getEnv("SMTP_FROM_NAME", "Clothing Shop"),

		// RajaOngkir
		RajaOngkirAPIKey:  getEnv("RAJAONGKIR_API_KEY", ""),
		RajaOngkirBaseURL: getEnv("RAJAONGKIR_BASE_URL", "https://api.rajaongkir.com/starter"),

		// Midtrans
		MidtransServerKey:   getEnv("MIDTRANS_SERVER_KEY", ""),
		MidtransClientKey:   getEnv("MIDTRANS_CLIENT_KEY", ""),
		MidtransEnvironment: getEnv("MIDTRANS_ENVIRONMENT", "sandbox"),

		// File Storage
		StorageDriver: getEnv("STORAGE_DRIVER", "local"),
		StoragePath:   getEnv("STORAGE_PATH", "./uploads"),

		// Debug and Environment
		Debug:      debug,
		Production: production,

		// Admin
		AdminEmail:    getEnv("ADMIN_EMAIL", "admin@example.com"),
		AdminPassword: getEnv("ADMIN_PASSWORD", "admin_password"),
	}

	return config, nil
}

func GetConfig() *Config {
	if config == nil {
		config, _ = LoadConfig()
	}
	return config
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
