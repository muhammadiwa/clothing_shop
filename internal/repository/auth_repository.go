package repository

import (
	"clothing-shop-api/internal/domain/models"
	"time"

	"gorm.io/gorm"
)

type AuthRepository interface {
	CreateUser(user *models.User) error
	FindByEmail(email string) (*models.User, error)
	UpdatePassword(userID uint, hashedPassword string) error
	UpdateLastLogin(userID uint) error
	VerifyEmail(userID uint) error
}

type authRepository struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) AuthRepository {
	return &authRepository{db: db}
}

func (r *authRepository) CreateUser(user *models.User) error {
	return r.db.Create(user).Error
}

func (r *authRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *authRepository) UpdatePassword(userID uint, hashedPassword string) error {
	return r.db.Model(&models.User{}).Where("id = ?", userID).
		Update("password", hashedPassword).Error
}

func (r *authRepository) UpdateLastLogin(userID uint) error {
	return r.db.Model(&models.User{}).Where("id = ?", userID).
		Update("last_login", time.Now()).Error
}

func (r *authRepository) VerifyEmail(userID uint) error {
	return r.db.Model(&models.User{}).Where("id = ?", userID).
		Update("email_verified", true).Error
}
