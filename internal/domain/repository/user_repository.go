package repository

import (
	"context"

	"fashion-shop/internal/domain/entity"
)

// UserRepository defines the interface for user data access
type UserRepository interface {
	Create(ctx context.Context, user *entity.User) error
	GetByID(ctx context.Context, id uint) (*entity.User, error)
	GetByEmail(ctx context.Context, email string) (*entity.User, error)
	Update(ctx context.Context, user *entity.User) error
	Delete(ctx context.Context, id uint) error
	List(ctx context.Context, offset, limit int) ([]*entity.User, int64, error)
	UpdateLastLogin(ctx context.Context, id uint) error
	ChangePassword(ctx context.Context, id uint, hashedPassword string) error
	ToggleActive(ctx context.Context, id uint, isActive bool) error
}

// AddressRepository defines the interface for address data access
type AddressRepository interface {
	Create(ctx context.Context, address *entity.Address) error
	GetByID(ctx context.Context, id uint) (*entity.Address, error)
	GetByUserID(ctx context.Context, userID uint) ([]*entity.Address, error)
	GetDefaultByUserID(ctx context.Context, userID uint) (*entity.Address, error)
	Update(ctx context.Context, address *entity.Address) error
	Delete(ctx context.Context, id uint) error
	SetDefault(ctx context.Context, id uint, userID uint) error
}
