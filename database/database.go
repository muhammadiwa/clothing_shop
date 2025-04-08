package database

import (
	"fmt"
	"log"
	"time"

	"github.com/fashion-shop/config"
	"github.com/fashion-shop/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DB is the database connection
var DB *gorm.DB

// Initialize initializes the database connection
func Initialize(cfg *config.Config) error {
	var err error

	// Configure GORM logger
	newLogger := logger.New(
		log.New(log.Writer(), "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logger.Info,
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		},
	)

	// Connect to database
	DB, err = gorm.Open(postgres.Open(cfg.GetDSN()), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	// Get underlying SQL DB
	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get database: %w", err)
	}

	// Set connection pool settings
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// Auto migrate the schema
	err = migrateSchema()
	if err != nil {
		return fmt.Errorf("failed to migrate schema: %w", err)
	}

	log.Println("Database connected and migrated successfully")
	return nil
}

// migrateSchema migrates the database schema
func migrateSchema() error {
	return DB.AutoMigrate(
		&models.User{},
		&models.UserAddress{},
		&models.UserActivity{},
		&models.Category{},
		&models.Tag{},
		&models.Product{},
		&models.ProductImage{},
		&models.ProductVariant{},
		&models.Review{},
		&models.Cart{},
		&models.CartItem{},
		&models.Wishlist{},
		&models.WishlistItem{},
		&models.Order{},
		&models.OrderItem{},
		&models.Payment{},
		&models.Notification{},
	)
}
