package models

import (
	"time"
)

type Role string

const (
	RoleAdmin    Role = "admin"
	RoleCustomer Role = "customer"
)

type User struct {
	ID            uint      `json:"id" gorm:"primaryKey"`
	Email         string    `json:"email" gorm:"uniqueIndex;not null"`
	Password      string    `json:"-" gorm:"not null"` // "-" excludes from JSON
	Name          string    `json:"name" gorm:"not null"`
	Phone         string    `json:"phone"`
	IsActive      bool      `json:"is_active" gorm:"default:true"`
	Role          string    `json:"role" gorm:"default:'user'"`
	EmailVerified bool      `json:"email_verified" gorm:"default:false"`
	LastLogin     time.Time `json:"last_login"`
	CreatedAt     time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt     time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}
