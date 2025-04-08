package models

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

// User represents a user in the system
type User struct {
	ID           string     `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Email        string     `json:"email" gorm:"uniqueIndex;not null"`
	Password     string     `json:"-" gorm:"not null"` // Password is not exposed in JSON
	Name         string     `json:"name" gorm:"not null"`
	Phone        string     `json:"phone"`
	Role         string     `json:"role" gorm:"default:user"`
	IsActive     bool       `json:"is_active" gorm:"default:true"`
	LastLogin    time.Time  `json:"last_login"`
	RefreshToken string     `json:"-"`
	TokenExpiry  time.Time  `json:"-"`
	CreatedAt    time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt    time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt    *time.Time `json:"-" gorm:"index"`
}

// UserAddress represents a user's address
type UserAddress struct {
	ID         string    `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	UserID     string    `json:"user_id" gorm:"type:uuid;not null"`
	Name       string    `json:"name" gorm:"not null"`
	Phone      string    `json:"phone" gorm:"not null"`
	Address    string    `json:"address" gorm:"not null"`
	City       string    `json:"city" gorm:"not null"`
	Province   string    `json:"province" gorm:"not null"`
	PostalCode string    `json:"postal_code" gorm:"not null"`
	IsDefault  bool      `json:"is_default" gorm:"default:false"`
	CreatedAt  time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt  time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

// UserActivity represents a log of user activity
type UserActivity struct {
	ID        string    `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	UserID    string    `json:"user_id" gorm:"type:uuid;not null"`
	Activity  string    `json:"activity" gorm:"not null"`
	IPAddress string    `json:"ip_address"`
	UserAgent string    `json:"user_agent"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
}

// SetPassword hashes a password and sets it on the user
func (u *User) SetPassword(password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

// CheckPassword verifies a password against the user's hashed password
func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

// TableName specifies the table name for User
func (User) TableName() string {
	return "users"
}

// TableName specifies the table name for UserAddress
func (UserAddress) TableName() string {
	return "user_addresses"
}

// TableName specifies the table name for UserActivity
func (UserActivity) TableName() string {
	return "user_activities"
}
