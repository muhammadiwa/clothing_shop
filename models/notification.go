package models

import (
	"time"
)

// Notification represents a notification for a user
type Notification struct {
	ID          string    `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	UserID      string    `json:"user_id" gorm:"type:uuid;not null"`
	User        User      `json:"user" gorm:"foreignKey:UserID"`
	Title       string    `json:"title" gorm:"not null"`
	Message     string    `json:"message" gorm:"not null"`
	Type        string    `json:"type" gorm:"not null"` // order_status, payment, system, etc.
	IsRead      bool      `json:"is_read" gorm:"default:false"`
	RelatedID   *string   `json:"related_id" gorm:"type:uuid"` // ID of related entity (order, payment, etc.)
	RelatedType *string   `json:"related_type"`                // Type of related entity (order, payment, etc.)
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

// TableName specifies the table name for Notification
func (Notification) TableName() string {
	return "notifications"
}
