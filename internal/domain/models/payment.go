package models

import "time"

type PaymentStatus string

const (
	PaymentStatusPending  PaymentStatus = "pending"
	PaymentStatusPaid     PaymentStatus = "paid"
	PaymentStatusFailed   PaymentStatus = "failed"
	PaymentStatusExpired  PaymentStatus = "expired"
	PaymentStatusRefunded PaymentStatus = "refunded"
)

type Payment struct {
	ID             uint          `json:"id"`
	OrderID        uint          `json:"order_id"`
	PaymentMethod  string        `json:"payment_method"`
	PaymentChannel string        `json:"payment_channel"`
	Amount         float64       `json:"amount"`
	Status         PaymentStatus `json:"status"`
	TransactionID  string        `json:"transaction_id"`
	PaymentToken   string        `json:"payment_token"`
	VANumber       string        `json:"va_number"`
	PaymentURL     string        `json:"payment_url"`
	ExpiryTime     *time.Time    `json:"expiry_time"`
	PaidAt         *time.Time    `json:"paid_at"`
	CreatedAt      time.Time     `json:"created_at"`
	UpdatedAt      time.Time     `json:"updated_at"`
}
