package repository

import (
	"context"

	"fashion-shop/internal/domain/entity"
)

// NotificationRepository defines the interface for notification data access
type NotificationRepository interface {
	Create(ctx context.Context, notification *entity.Notification) error
	GetByID(ctx context.Context, id uint) (*entity.Notification, error)
	GetByUserID(ctx context.Context, userID uint, offset, limit int) ([]*entity.Notification, int64, error)
	GetUnreadByUserID(ctx context.Context, userID uint) ([]*entity.Notification, error)
	MarkAsRead(ctx context.Context, id uint) error
	MarkAllAsRead(ctx context.Context, userID uint) error
	Delete(ctx context.Context, id uint) error
}
