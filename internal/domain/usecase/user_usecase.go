package usecase

import (
	"context"

	"fashion-shop/internal/domain/entity"
)

// UserUseCase defines the interface for user business logic
type UserUseCase interface {
	Register(ctx context.Context, email, password, name, phone string) (*entity.User, error)
	Login(ctx context.Context, email, password string) (string, string, error) // returns access token, refresh token, error
	RefreshToken(ctx context.Context, refreshToken string) (string, error)     // returns new access token, error
	GetProfile(ctx context.Context, userID uint) (*entity.User, error)
	UpdateProfile(ctx context.Context, userID uint, name, phone string) (*entity.User, error)
	ChangePassword(ctx context.Context, userID uint, oldPassword, newPassword string) error
	RequestPasswordReset(ctx context.Context, email string) error
	ResetPassword(ctx context.Context, token, newPassword string) error
	Logout(ctx context.Context, userID uint) error

	// Admin functions
	GetUsers(ctx context.Context, page, limit int) ([]*entity.User, int64, error)
	GetUserByID(ctx context.Context, id uint) (*entity.User, error)
	ToggleUserActive(ctx context.Context, id uint, isActive bool) error
	ResetUserPassword(ctx context.Context, id uint, newPassword string) error
}

// AddressUseCase defines the interface for address business logic
type AddressUseCase interface {
	CreateAddress(ctx context.Context, userID uint, address *entity.Address) (*entity.Address, error)
	GetAddresses(ctx context.Context, userID uint) ([]*entity.Address, error)
	GetAddressByID(ctx context.Context, id uint, userID uint) (*entity.Address, error)
	UpdateAddress(ctx context.Context, id uint, userID uint, address *entity.Address) (*entity.Address, error)
	DeleteAddress(ctx context.Context, id uint, userID uint) error
	SetDefaultAddress(ctx context.Context, id uint, userID uint) error
	GetDefaultAddress(ctx context.Context, userID uint) (*entity.Address, error)
}
