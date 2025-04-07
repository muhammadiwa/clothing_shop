package models

import "time"

type Role string

const (
	RoleAdmin    Role = "admin"
	RoleCustomer Role = "customer"
)

type User struct {
	ID                uint       `json:"id" gorm:"primaryKey"`
	Username          string     `json:"username" gorm:"unique;not null"`
	Email             string     `json:"email" gorm:"unique;not null"`
	Password          string     `json:"-" gorm:"not null"`
	Role              Role       `json:"role" gorm:"type:enum('admin','customer');default:'customer'"`
	IsVerified        bool       `json:"is_verified" gorm:"default:false"`
	VerificationToken string     `json:"-" gorm:"nullable"`
	ResetToken        string     `json:"-" gorm:"nullable"`
	ResetTokenExpiry  *time.Time `json:"-" gorm:"nullable"`
	PhoneNumber       string     `json:"phone_number" gorm:"nullable"`
	CreatedAt         time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt         time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt         *time.Time `json:"-" gorm:"index"`
}
