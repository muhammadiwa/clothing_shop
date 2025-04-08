package migrations

import (
	"github.com/fashion-shop/database"
	"github.com/fashion-shop/models"
)

// Migrate runs all migrations
func Migrate() error {
	// Auto migrate the schema
	return database.DB.AutoMigrate(
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

// Seed seeds the database with initial data
func Seed() error {
	// Create admin user if not exists
	var adminCount int64
	database.DB.Model(&models.User{}).Where("role = ?", "admin").Count(&adminCount)

	if adminCount == 0 {
		adminUser := models.User{
			Name:     "Admin",
			Email:    "admin@fashionshop.com",
			Password: "admin123", // This will be hashed by the BeforeCreate hook
			Phone:    "1234567890",
			Role:     "admin",
			IsActive: true,
		}

		if err := database.DB.Create(&adminUser).Error; err != nil {
			return err
		}
	}

	// Create default categories if not exists
	var categoryCount int64
	database.DB.Model(&models.Category{}).Count(&categoryCount)

	if categoryCount == 0 {
		categories := []models.Category{
			{Name: "Men", Slug: "men", Description: "Men's clothing", IsActive: true},
			{Name: "Women", Slug: "women", Description: "Women's clothing", IsActive: true},
			{Name: "Kids", Slug: "kids", Description: "Kids' clothing", IsActive: true},
			{Name: "Accessories", Slug: "accessories", Description: "Fashion accessories", IsActive: true},
		}

		for _, category := range categories {
			if err := database.DB.Create(&category).Error; err != nil {
				return err
			}
		}
	}

	// Create default tags if not exists
	var tagCount int64
	database.DB.Model(&models.Tag{}).Count(&tagCount)

	if tagCount == 0 {
		tags := []models.Tag{
			{Name: "New Arrival", Slug: "new-arrival"},
			{Name: "Best Seller", Slug: "best-seller"},
			{Name: "Sale", Slug: "sale"},
			{Name: "Featured", Slug: "featured"},
			{Name: "Limited Edition", Slug: "limited-edition"},
		}

		for _, tag := range tags {
			if err := database.DB.Create(&tag).Error; err != nil {
				return err
			}
		}
	}

	return nil
}

// Rollback rolls back all migrations
func Rollback() error {
	return database.DB.Migrator().DropTable(
		&models.Notification{},
		&models.Payment{},
		&models.OrderItem{},
		&models.Order{},
		&models.WishlistItem{},
		&models.Wishlist{},
		&models.CartItem{},
		&models.Cart{},
		&models.Review{},
		&models.ProductVariant{},
		&models.ProductImage{},
		&models.Product{},
		&models.Tag{},
		&models.Category{},
		&models.UserActivity{},
		&models.UserAddress{},
		&models.User{},
	)
}

// Reset resets the database by rolling back and then migrating
func Reset() error {
	if err := Rollback(); err != nil {
		return err
	}

	if err := Migrate(); err != nil {
		return err
	}

	return Seed()
}
