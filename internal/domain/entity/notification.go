package entity

import (
	"time"
)

// NotificationType represents the type of notification
type NotificationType string

const (
	NotificationTypeOrder     NotificationType = "order"
	NotificationTypePayment   NotificationType = "payment"
	NotificationTypePromotion NotificationType = "promotion"
	NotificationTypeSystem    NotificationType = "system"
)

// Notification represents a notification for a user
type Notification struct {
	ID        uint             `gorm:"primaryKey" json:"id"`
	UserID    uint             `gorm:"index;not null" json:"user_id"`
	User      User             `gorm:"foreignKey:UserID" json:"-"`
	Type      NotificationType `gorm:"type:varchar(20);not null" json:"type"`
	Title     string           `gorm:"not null" json:"title"`
	Message   string           `gorm:"not null" json:"message"`
	Data      string           `json:"data,omitempty"` // JSON string containing additional data
	IsRead    bool             `gorm:"default:false" json:"is_read"`
	CreatedAt time.Time        `json:"created_at"`
	UpdatedAt time.Time        `json:"updated_at"`
}
