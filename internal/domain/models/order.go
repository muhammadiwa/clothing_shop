package models

import "time"

type Order struct {
    ID          uint      `json:"id" gorm:"primaryKey"`
    UserID      uint      `json:"user_id"`
    ProductID   uint      `json:"product_id"`
    Quantity    int       `json:"quantity"`
    TotalPrice  float64   `json:"total_price"`
    OrderDate   time.Time `json:"order_date"`
    Status      string    `json:"status"`
}