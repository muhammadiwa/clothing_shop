package models

import (
	"time"
)

// Product represents a product in the system
type Product struct {
	ID          string           `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Name        string           `json:"name" gorm:"not null"`
	Description string           `json:"description" gorm:"type:text"`
	Price       float64          `json:"price" gorm:"not null"`
	Stock       int              `json:"stock" gorm:"not null"`
	SKU         string           `json:"sku" gorm:"uniqueIndex"`
	Weight      float64          `json:"weight" gorm:"not null"` // Weight in grams
	CategoryID  string           `json:"category_id" gorm:"type:uuid;not null"`
	Category    Category         `json:"category" gorm:"foreignKey:CategoryID"`
	Images      []ProductImage   `json:"images" gorm:"foreignKey:ProductID"`
	Variants    []ProductVariant `json:"variants" gorm:"foreignKey:ProductID"`
	Tags        []Tag            `json:"tags" gorm:"many2many:product_tags;"`
	IsActive    bool             `json:"is_active" gorm:"default:true"`
	CreatedAt   time.Time        `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time        `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt   *time.Time       `json:"-" gorm:"index"`
}

// ProductImage represents an image for a product
type ProductImage struct {
	ID        string    `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	ProductID string    `json:"product_id" gorm:"type:uuid;not null"`
	URL       string    `json:"url" gorm:"not null"`
	IsPrimary bool      `json:"is_primary" gorm:"default:false"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

// ProductVariant represents a variant of a product (e.g., size, color)
type ProductVariant struct {
	ID        string    `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	ProductID string    `json:"product_id" gorm:"type:uuid;not null"`
	Size      string    `json:"size"`
	Color     string    `json:"color"`
	Stock     int       `json:"stock" gorm:"not null"`
	Price     float64   `json:"price"` // Optional override of base product price
	SKU       string    `json:"sku" gorm:"uniqueIndex"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

// Category represents a product category
type Category struct {
	ID          string    `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Name        string    `json:"name" gorm:"not null"`
	Description string    `json:"description"`
	ParentID    *string   `json:"parent_id" gorm:"type:uuid"`
	Parent      *Category `json:"parent" gorm:"foreignKey:ParentID"`
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

// Tag represents a product tag
type Tag struct {
	ID        string    `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Name      string    `json:"name" gorm:"uniqueIndex;not null"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

// Review represents a product review
type Review struct {
	ID        string    `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	ProductID string    `json:"product_id" gorm:"type:uuid;not null"`
	Product   Product   `json:"product" gorm:"foreignKey:ProductID"`
	UserID    string    `json:"user_id" gorm:"type:uuid;not null"`
	User      User      `json:"user" gorm:"foreignKey:UserID"`
	Rating    int       `json:"rating" gorm:"not null"`
	Comment   string    `json:"comment"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

// TableName specifies the table name for Product
func (Product) TableName() string {
	return "products"
}

// TableName specifies the table name for ProductImage
func (ProductImage) TableName() string {
	return "product_images"
}

// TableName specifies the table name for ProductVariant
func (ProductVariant) TableName() string {
	return "product_variants"
}

// TableName specifies the table name for Category
func (Category) TableName() string {
	return "categories"
}

// TableName specifies the table name for Tag
func (Tag) TableName() string {
	return "tags"
}

// TableName specifies the table name for Review
func (Review) TableName() string {
	return "reviews"
}
