package usecase

import (
	"context"

	"fashion-shop/internal/domain/entity"
)

// NotificationUseCase defines the interface for notification business logic
type NotificationUseCase interface {
	CreateNotification(ctx context.Context, userID uint, notificationType entity.NotificationType, title, message string, data map[string]interface{}) error
	GetNotifications(ctx context.Context, userID uint, page, limit int) ([]*entity.Notification, int64, error)
	GetUnreadNotifications(ctx context.Context, userID uint) ([]*entity.Notification, error)
	MarkAsRead(ctx context.Context, id, userID uint) error
	MarkAllAsRead(ctx context.Context, userID uint) error
	DeleteNotification(ctx context.Context, id, userID uint) error
}
