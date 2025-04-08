package models

import (
	"time"
)

// Cart represents a user's shopping cart
type Cart struct {
	ID        string     `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	UserID    string     `json:"user_id" gorm:"type:uuid;not null;uniqueIndex"`
	User      User       `json:"user" gorm:"foreignKey:UserID"`
	Items     []CartItem `json:"items" gorm:"foreignKey:CartID"`
	CreatedAt time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
}

// CartItem represents an item in a user's shopping cart
type CartItem struct {
	ID        string          `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	CartID    string          `json:"cart_id" gorm:"type:uuid;not null"`
	ProductID string          `json:"product_id" gorm:"type:uuid;not null"`
	Product   Product         `json:"product" gorm:"foreignKey:ProductID"`
	VariantID *string         `json:"variant_id" gorm:"type:uuid"`
	Variant   *ProductVariant `json:"variant" gorm:"foreignKey:VariantID"`
	Quantity  int             `json:"quantity" gorm:"not null"`
	CreatedAt time.Time       `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time       `json:"updated_at" gorm:"autoUpdateTime"`
}

// Wishlist represents a user's wishlist
type Wishlist struct {
	ID        string         `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	UserID    string         `json:"user_id" gorm:"type:uuid;not null;uniqueIndex"`
	User      User           `json:"user" gorm:"foreignKey:UserID"`
	Items     []WishlistItem `json:"items" gorm:"foreignKey:WishlistID"`
	CreatedAt time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
}

// WishlistItem represents an item in a user's wishlist
type WishlistItem struct {
	ID         string          `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	WishlistID string          `json:"wishlist_id" gorm:"type:uuid;not null"`
	ProductID  string          `json:"product_id" gorm:"type:uuid;not null"`
	Product    Product         `json:"product" gorm:"foreignKey:ProductID"`
	VariantID  *string         `json:"variant_id" gorm:"type:uuid"`
	Variant    *ProductVariant `json:"variant" gorm:"foreignKey:VariantID"`
	CreatedAt  time.Time       `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt  time.Time       `json:"updated_at" gorm:"autoUpdateTime"`
}

// TableName specifies the table name for Cart
func (Cart) TableName() string {
	return "carts"
}

// TableName specifies the table name for CartItem
func (CartItem) TableName() string {
	return "cart_items"
}

// TableName specifies the table name for Wishlist
func (Wishlist) TableName() string {
	return "wishlists"
}

// TableName specifies the table name for WishlistItem
func (WishlistItem) TableName() string {
	return "wishlist_items"
}
