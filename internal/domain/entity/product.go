package entity

import (
	"time"

	"gorm.io/gorm"
)

// Category represents a product category
type Category struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"not null" json:"name"`
	Slug        string         `gorm:"uniqueIndex;not null" json:"slug"`
	Description string         `json:"description"`
	ParentID    *uint          `json:"parent_id,omitempty"`
	Parent      *Category      `gorm:"foreignKey:ParentID" json:"-"`
	Image       string         `json:"image,omitempty"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// Product represents a product in the system
type Product struct {
	ID            uint             `gorm:"primaryKey" json:"id"`
	Name          string           `gorm:"not null" json:"name"`
	Slug          string           `gorm:"uniqueIndex;not null" json:"slug"`
	Description   string           `json:"description"`
	Price         float64          `gorm:"not null" json:"price"`
	DiscountPrice *float64         `json:"discount_price,omitempty"`
	CategoryID    uint             `gorm:"index;not null" json:"category_id"`
	Category      Category         `gorm:"foreignKey:CategoryID" json:"-"`
	Images        []ProductImage   `gorm:"foreignKey:ProductID" json:"images,omitempty"`
	Variants      []ProductVariant `gorm:"foreignKey:ProductID" json:"variants,omitempty"`
	Tags          []Tag            `gorm:"many2many:product_tags;" json:"tags,omitempty"`
	IsActive      bool             `gorm:"default:true" json:"is_active"`
	CreatedAt     time.Time        `json:"created_at"`
	UpdatedAt     time.Time        `json:"updated_at"`
	DeletedAt     gorm.DeletedAt   `gorm:"index" json:"-"`
}

// ProductImage represents a product image
type ProductImage struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	ProductID uint      `gorm:"index;not null" json:"product_id"`
	URL       string    `gorm:"not null" json:"url"`
	IsPrimary bool      `gorm:"default:false" json:"is_primary"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ProductVariant represents a product variant (size, color, etc.)
type ProductVariant struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	ProductID uint      `gorm:"index;not null" json:"product_id"`
	Size      string    `json:"size"`
	Color     string    `json:"color"`
	SKU       string    `gorm:"uniqueIndex;not null" json:"sku"`
	Stock     int       `gorm:"not null" json:"stock"`
	Weight    float64   `gorm:"not null" json:"weight"` // in grams
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Tag represents a product tag
type Tag struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"uniqueIndex;not null" json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Review represents a product review
type Review struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	ProductID uint      `gorm:"index;not null" json:"product_id"`
	Product   Product   `gorm:"foreignKey:ProductID" json:"-"`
	UserID    uint      `gorm:"index;not null" json:"user_id"`
	User      User      `gorm:"foreignKey:UserID" json:"-"`
	OrderID   uint      `gorm:"index;not null" json:"order_id"`
	Rating    int       `gorm:"not null" json:"rating"` // 1-5
	Comment   string    `json:"comment"`
	Images    []string  `gorm:"-" json:"images,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
