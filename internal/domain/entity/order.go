package entity

import (
	"time"

	"gorm.io/gorm"
)

// OrderStatus represents the status of an order
type OrderStatus string

const (
	OrderStatusPending    OrderStatus = "pending"
	OrderStatusProcessing OrderStatus = "processing"
	OrderStatusShipped    OrderStatus = "shipped"
	OrderStatusDelivered  OrderStatus = "delivered"
	OrderStatusCancelled  OrderStatus = "cancelled"
	OrderStatusRefunded   OrderStatus = "refunded"
)

// PaymentStatus represents the status of a payment
type PaymentStatus string

const (
	PaymentStatusPending  PaymentStatus = "pending"
	PaymentStatusPaid     PaymentStatus = "paid"
	PaymentStatusFailed   PaymentStatus = "failed"
	PaymentStatusRefunded PaymentStatus = "refunded"
	PaymentStatusExpired  PaymentStatus = "expired"
)

// PaymentMethod represents the payment method
type PaymentMethod string

const (
	PaymentMethodCreditCard   PaymentMethod = "credit_card"
	PaymentMethodBankTransfer PaymentMethod = "bank_transfer"
	PaymentMethodEWallet      PaymentMethod = "e_wallet"
	PaymentMethodCOD          PaymentMethod = "cod"
)

// Order represents an order in the system
type Order struct {
	ID                     uint           `gorm:"primaryKey" json:"id"`
	UserID                 uint           `gorm:"index;not null" json:"user_id"`
	User                   User           `gorm:"foreignKey:UserID" json:"-"`
	OrderNumber            string         `gorm:"uniqueIndex;not null" json:"order_number"`
	Status                 OrderStatus    `gorm:"type:varchar(20);default:pending" json:"status"`
	TotalAmount            float64        `gorm:"not null" json:"total_amount"`
	ShippingCost           float64        `gorm:"not null" json:"shipping_cost"`
	DiscountAmount         float64        `gorm:"default:0" json:"discount_amount"`
	FinalAmount            float64        `gorm:"not null" json:"final_amount"`
	ShippingAddress        Address        `gorm:"embedded" json:"shipping_address"`
	ShippingMethod         string         `json:"shipping_method"`
	ShippingTrackingNumber string         `json:"shipping_tracking_number,omitempty"`
	PaymentID              *uint          `json:"payment_id,omitempty"`
	Payment                *Payment       `gorm:"foreignKey:PaymentID" json:"payment,omitempty"`
	OrderItems             []OrderItem    `gorm:"foreignKey:OrderID" json:"order_items,omitempty"`
	Notes                  string         `json:"notes,omitempty"`
	CreatedAt              time.Time      `json:"created_at"`
	UpdatedAt              time.Time      `json:"updated_at"`
	DeletedAt              gorm.DeletedAt `gorm:"index" json:"-"`
}

// OrderItem represents an item in an order
type OrderItem struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	OrderID       uint      `gorm:"index;not null" json:"order_id"`
	ProductID     uint      `gorm:"index;not null" json:"product_id"`
	ProductName   string    `gorm:"not null" json:"product_name"`
	VariantID     uint      `json:"variant_id"`
	VariantInfo   string    `json:"variant_info"` // JSON string containing size, color, etc.
	Quantity      int       `gorm:"not null" json:"quantity"`
	Price         float64   `gorm:"not null" json:"price"`
	DiscountPrice float64   `gorm:"default:0" json:"discount_price"`
	FinalPrice    float64   `gorm:"not null" json:"final_price"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// Payment represents a payment for an order
type Payment struct {
	ID            uint          `gorm:"primaryKey" json:"id"`
	OrderID       uint          `gorm:"uniqueIndex;not null" json:"order_id"`
	PaymentMethod PaymentMethod `gorm:"type:varchar(20);not null" json:"payment_method"`
	Amount        float64       `gorm:"not null" json:"amount"`
	Status        PaymentStatus `gorm:"type:varchar(20);default:pending" json:"status"`
	TransactionID string        `json:"transaction_id,omitempty"`
	PaymentURL    string        `json:"payment_url,omitempty"`
	PaidAt        *time.Time    `json:"paid_at,omitempty"`
	ExpiredAt     *time.Time    `json:"expired_at,omitempty"`
	CreatedAt     time.Time     `json:"created_at"`
	UpdatedAt     time.Time     `json:"updated_at"`
}

// Cart represents a user's shopping cart
type Cart struct {
	ID        uint       `gorm:"primaryKey" json:"id"`
	UserID    uint       `gorm:"uniqueIndex;not null" json:"user_id"`
	User      User       `gorm:"foreignKey:UserID" json:"-"`
	Items     []CartItem `gorm:"foreignKey:CartID" json:"items,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

// CartItem represents an item in a user's shopping cart
type CartItem struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CartID    uint      `gorm:"index;not null" json:"cart_id"`
	ProductID uint      `gorm:"index;not null" json:"product_id"`
	Product   Product   `gorm:"foreignKey:ProductID" json:"-"`
	VariantID uint      `json:"variant_id"`
	Quantity  int       `gorm:"not null" json:"quantity"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Wishlist represents a user's wishlist
type Wishlist struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	UserID    uint           `gorm:"uniqueIndex;not null" json:"user_id"`
	User      User           `gorm:"foreignKey:UserID" json:"-"`
	Items     []WishlistItem `gorm:"foreignKey:WishlistID" json:"items,omitempty"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}

// WishlistItem represents an item in a user's wishlist
type WishlistItem struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	WishlistID uint      `gorm:"index;not null" json:"wishlist_id"`
	ProductID  uint      `gorm:"index;not null" json:"product_id"`
	Product    Product   `gorm:"foreignKey:ProductID" json:"-"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
