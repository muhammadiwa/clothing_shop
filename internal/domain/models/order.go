package models

import "time"

type OrderStatus string

const (
	OrderStatusPending    OrderStatus = "pending"
	OrderStatusProcessing OrderStatus = "processing"
	OrderStatusShipped    OrderStatus = "shipped"
	OrderStatusDelivered  OrderStatus = "delivered"
	OrderStatusCancelled  OrderStatus = "cancelled"
	OrderStatusRefunded   OrderStatus = "refunded"
)

type PaymentStatus string

const (
	PaymentStatusPending  PaymentStatus = "pending"
	PaymentStatusPaid     PaymentStatus = "paid"
	PaymentStatusFailed   PaymentStatus = "failed"
	PaymentStatusExpired  PaymentStatus = "expired"
	PaymentStatusRefunded PaymentStatus = "refunded"
)

type Order struct {
	ID          uint        `json:"id" gorm:"primaryKey"`
	UserID      uint        `json:"user_id"`
	User        User        `json:"user" gorm:"foreignKey:UserID"`
	OrderNumber string      `json:"order_number" gorm:"unique;not null"`
	Status      OrderStatus `json:"status" gorm:"type:enum('pending','processing','shipped','delivered','cancelled','refunded');default:'pending'"`
	TotalAmount float64     `json:"total_amount" gorm:"type:decimal(12,2);not null"`
	ShippingFee float64     `json:"shipping_fee" gorm:"type:decimal(10,2);not null"`
	GrandTotal  float64     `json:"grand_total" gorm:"type:decimal(12,2);not null"`

	AddressID uint    `json:"address_id"`
	Address   Address `json:"address" gorm:"foreignKey:AddressID"`

	CourierID        uint       `json:"courier_id"`
	Courier          Courier    `json:"courier" gorm:"foreignKey:CourierID"`
	CourierService   string     `json:"courier_service"`
	TrackingNumber   string     `json:"tracking_number"`
	EstimatedArrival *time.Time `json:"estimated_arrival"`

	PaymentID uint    `json:"payment_id"`
	Payment   Payment `json:"payment" gorm:"foreignKey:PaymentID"`

	Notes     string     `json:"notes"`
	CreatedAt time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt *time.Time `json:"-" gorm:"index"`

	Items []OrderItem `json:"items" gorm:"foreignKey:OrderID"`
}

type OrderItem struct {
	ID               uint           `json:"id" gorm:"primaryKey"`
	OrderID          uint           `json:"order_id"`
	ProductVariantID uint           `json:"product_variant_id"`
	ProductVariant   ProductVariant `json:"product_variant" gorm:"foreignKey:ProductVariantID"`
	Quantity         int            `json:"quantity" gorm:"not null"`
	UnitPrice        float64        `json:"unit_price" gorm:"type:decimal(12,2);not null"`
	TotalPrice       float64        `json:"total_price" gorm:"type:decimal(12,2);not null"`
	CreatedAt        time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt        time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
}

type Payment struct {
	ID             uint          `json:"id" gorm:"primaryKey"`
	OrderID        uint          `json:"order_id"`
	PaymentMethod  string        `json:"payment_method" gorm:"not null"`
	PaymentChannel string        `json:"payment_channel" gorm:"not null"`
	Amount         float64       `json:"amount" gorm:"type:decimal(12,2);not null"`
	Status         PaymentStatus `json:"status" gorm:"type:enum('pending','paid','failed','expired','refunded');default:'pending'"`
	TransactionID  string        `json:"transaction_id"`
	PaymentToken   string        `json:"payment_token"`
	VaNumber       string        `json:"va_number"`
	PaymentURL     string        `json:"payment_url"`
	ExpiryTime     *time.Time    `json:"expiry_time"`
	PaidAt         *time.Time    `json:"paid_at"`
	CreatedAt      time.Time     `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt      time.Time     `json:"updated_at" gorm:"autoUpdateTime"`
}

type Cart struct {
	ID               uint           `json:"id" gorm:"primaryKey"`
	UserID           uint           `json:"user_id"`
	User             User           `json:"-" gorm:"foreignKey:UserID"`
	ProductVariantID uint           `json:"product_variant_id"`
	ProductVariant   ProductVariant `json:"product_variant" gorm:"foreignKey:ProductVariantID"`
	Quantity         int            `json:"quantity" gorm:"not null"`
	CreatedAt        time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt        time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
}

type Wishlist struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	UserID    uint      `json:"user_id"`
	User      User      `json:"-" gorm:"foreignKey:UserID"`
	ProductID uint      `json:"product_id"`
	Product   Product   `json:"product" gorm:"foreignKey:ProductID"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
}
