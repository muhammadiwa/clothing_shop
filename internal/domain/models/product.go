package models

import "time"

type Product struct {
	ID          uint       `json:"id" gorm:"primaryKey"`
	Name        string     `json:"name" gorm:"not null"`
	Description string     `json:"description"`
	CategoryID  uint       `json:"category_id"`
	Category    Category   `json:"category" gorm:"foreignKey:CategoryID"`
	BasePrice   float64    `json:"base_price" gorm:"not null;type:decimal(12,2)"`
	Discount    float64    `json:"discount" gorm:"default:0;type:decimal(5,2)"`
	Weight      float64    `json:"weight" gorm:"not null;comment:'Weight in grams'"`
	Rating      float32    `json:"rating" gorm:"default:0;type:decimal(3,2)"`
	ReviewCount int        `json:"review_count" gorm:"default:0"`
	CreatedAt   time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt   *time.Time `json:"-" gorm:"index"`

	Variants []ProductVariant `json:"variants"`
	Images   []ProductImage   `json:"images"`
	Reviews  []ProductReview  `json:"reviews,omitempty" gorm:"foreignKey:ProductID"`
}

type ProductVariant struct {
	ID        uint       `json:"id" gorm:"primaryKey"`
	ProductID uint       `json:"product_id"`
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

type Category struct {
	ID          uint       `json:"id" gorm:"primaryKey"`
	Name        string     `json:"name" gorm:"not null"`
	Slug        string     `json:"slug" gorm:"unique;not null"`
	Description string     `json:"description"`
	ParentID    *uint      `json:"parent_id"`
	Parent      *Category  `json:"parent,omitempty" gorm:"foreignKey:ParentID"`
	Children    []Category `json:"children,omitempty" gorm:"foreignKey:ParentID"`
	CreatedAt   time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt   *time.Time `json:"-" gorm:"index"`
}
