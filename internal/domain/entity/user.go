package entity

import (
	"time"

	"gorm.io/gorm"
)

// Role represents user roles in the system
type Role string

const (
	RoleAdmin  Role = "admin"
	RoleUser   Role = "user"
	RoleSeller Role = "seller"
)

// User represents a user in the system
type User struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Email     string         `gorm:"uniqueIndex;not null" json:"email"`
	Password  string         `gorm:"not null" json:"-"`
	Name      string         `gorm:"not null" json:"name"`
	Phone     string         `json:"phone"`
	Role      Role           `gorm:"type:varchar(20);default:user" json:"role"`
	IsActive  bool           `gorm:"default:true" json:"is_active"`
	LastLogin *time.Time     `json:"last_login,omitempty"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// Address represents a user's address
type Address struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	UserID      uint           `gorm:"index;not null" json:"user_id"`
	User        User           `gorm:"foreignKey:UserID" json:"-"`
	Label       string         `gorm:"not null" json:"label"` // e.g., "Home", "Office"
	Recipient   string         `gorm:"not null" json:"recipient"`
	Phone       string         `gorm:"not null" json:"phone"`
	Province    string         `gorm:"not null" json:"province"`
	City        string         `gorm:"not null" json:"city"`
	District    string         `gorm:"not null" json:"district"`
	PostalCode  string         `gorm:"not null" json:"postal_code"`
	FullAddress string         `gorm:"not null" json:"full_address"`
	IsDefault   bool           `gorm:"default:false" json:"is_default"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}
