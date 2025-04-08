package persistence

import (
	"context"
	"errors"
	"time"

	"fashion-shop/internal/domain/entity"
	"fashion-shop/internal/domain/repository"

	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

// NewUserRepository creates a new UserRepository instance
func NewUserRepository(db *gorm.DB) repository.UserRepository {
	return &userRepository{
		db: db,
	}
}

// Create creates a new user
func (r *userRepository) Create(ctx context.Context, user *entity.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

// GetByID gets a user by ID
func (r *userRepository) GetByID(ctx context.Context, id uint) (*entity.User, error) {
	var user entity.User
	if err := r.db.WithContext(ctx).First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}

// GetByEmail gets a user by email
func (r *userRepository) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	var user entity.User
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}

// Update updates a user
func (r *userRepository) Update(ctx context.Context, user *entity.User) error {
	return r.db.WithContext(ctx).Save(user).Error
}

// Delete deletes a user
func (r *userRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&entity.User{}, id).Error
}

// List lists users with pagination
func (r *userRepository) List(ctx context.Context, offset, limit int) ([]*entity.User, int64, error) {
	var users []*entity.User
	var count int64

	if err := r.db.WithContext(ctx).Model(&entity.User{}).Count(&count).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.WithContext(ctx).Offset(offset).Limit(limit).Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, count, nil
}

// UpdateLastLogin updates a user's last login time
func (r *userRepository) UpdateLastLogin(ctx context.Context, id uint) error {
	now := time.Now()
	return r.db.WithContext(ctx).Model(&entity.User{}).Where("id = ?", id).Update("last_login", &now).Error
}

// ChangePassword changes a user's password
func (r *userRepository) ChangePassword(ctx context.Context, id uint, hashedPassword string) error {
	return r.db.WithContext(ctx).Model(&entity.User{}).Where("id = ?", id).Update("password", hashedPassword).Error
}

// ToggleActive toggles a user's active status
func (r *userRepository) ToggleActive(ctx context.Context, id uint, isActive bool) error {
	return r.db.WithContext(ctx).Model(&entity.User{}).Where("id = ?", id).Update("is_active", isActive).Error
}
