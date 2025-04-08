package models

import "time"

type Product struct {
	ID          string    `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	CategoryID  string    `json:"category_id"`
	Category    Category  `json:"category"`
	Stock       int       `json:"stock"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type ProductVariant struct {
	ID        uint       `json:"id" gorm:"primaryKey"`
	ProductID uint       `json:"product_id"`
	Product   *Product   `json:"product,omitempty"` // Add this field
	Size      string     `json:"size"`
	Color     string     `json:"color"`
	SKU       string     `json:"sku" gorm:"unique"`
	Price     float64    `json:"price" gorm:"type:decimal(12,2)"`
	Stock     int        `json:"stock"`
	CreatedAt time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt *time.Time `json:"-" gorm:"index"`
}

type ProductImage struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	ProductID uint      `json:"product_id"`
	URL       string    `json:"url" gorm:"not null"`
	IsPrimary bool      `json:"is_primary" gorm:"default:false"`
	SortOrder int       `json:"sort_order" gorm:"default:0"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

type ProductReview struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	ProductID uint      `json:"product_id"`
	UserID    uint      `json:"user_id"`
	User      User      `json:"user" gorm:"foreignKey:UserID"`
	Rating    int       `json:"rating" gorm:"not null;check:rating >= 1 AND rating <= 5"`
	Comment   string    `json:"comment"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}
