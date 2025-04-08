package models

import (
	"time"
)

// Order represents a customer order
type Order struct {
	ID                string      `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	UserID            string      `json:"user_id" gorm:"type:uuid;not null"`
	User              User        `json:"user" gorm:"foreignKey:UserID"`
	OrderNumber       string      `json:"order_number" gorm:"uniqueIndex;not null"`
	Status            string      `json:"status" gorm:"not null;default:'pending'"` // pending, processing, shipped, delivered, cancelled
	Items             []OrderItem `json:"items" gorm:"foreignKey:OrderID"`
	ShippingAddressID string      `json:"shipping_address_id" gorm:"type:uuid;not null"`
	ShippingAddress   UserAddress `json:"shipping_address" gorm:"foreignKey:ShippingAddressID"`
	ShippingMethod    string      `json:"shipping_method" gorm:"not null"`
	ShippingCost      float64     `json:"shipping_cost" gorm:"not null"`
	TrackingNumber    string      `json:"tracking_number"`
	SubTotal          float64     `json:"sub_total" gorm:"not null"`
	Tax               float64     `json:"tax" gorm:"not null"`
	Discount          float64     `json:"discount" gorm:"default:0"`
	GrandTotal        float64     `json:"grand_total" gorm:"not null"`
	Notes             string      `json:"notes"`
	PaymentID         *string     `json:"payment_id" gorm:"type:uuid"`
	Payment           *Payment    `json:"payment" gorm:"foreignKey:PaymentID"`
	CreatedAt         time.Time   `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt         time.Time   `json:"updated_at" gorm:"autoUpdateTime"`
}

// OrderItem represents an item in an order
type OrderItem struct {
	ID        string          `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	OrderID   string          `json:"order_id" gorm:"type:uuid;not null"`
	ProductID string          `json:"product_id" gorm:"type:uuid;not null"`
	Product   Product         `json:"product" gorm:"foreignKey:ProductID"`
	VariantID *string         `json:"variant_id" gorm:"type:uuid"`
	Variant   *ProductVariant `json:"variant" gorm:"foreignKey:VariantID"`
	Quantity  int             `json:"quantity" gorm:"not null"`
	Price     float64         `json:"price" gorm:"not null"` // Price at time of purchase
	Total     float64         `json:"total" gorm:"not null"`
	CreatedAt time.Time       `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time       `json:"updated_at" gorm:"autoUpdateTime"`
}

// Payment represents a payment for an order
type Payment struct {
	ID             string     `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	OrderID        string     `json:"order_id" gorm:"type:uuid;not null"`
	Order          Order      `json:"order" gorm:"foreignKey:OrderID"`
	Amount         float64    `json:"amount" gorm:"not null"`
	Method         string     `json:"method" gorm:"not null"`                   // credit_card, bank_transfer, etc.
	Status         string     `json:"status" gorm:"not null;default:'pending'"` // pending, completed, failed, refunded
	TransactionID  string     `json:"transaction_id"`
	PaymentGateway string     `json:"payment_gateway" gorm:"not null"` // midtrans, etc.
	PaymentDetails string     `json:"payment_details" gorm:"type:jsonb"`
	RefundAmount   float64    `json:"refund_amount" gorm:"default:0"`
	RefundReason   string     `json:"refund_reason"`
	RefundDate     *time.Time `json:"refund_date"`
	CreatedAt      time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt      time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
}

// TableName specifies the table name for Order
func (Order) TableName() string {
	return "orders"
}

// TableName specifies the table name for OrderItem
func (OrderItem) TableName() string {
	return "order_items"
}

// TableName specifies the table name for Payment
func (Payment) TableName() string {
	return "payments"
}
