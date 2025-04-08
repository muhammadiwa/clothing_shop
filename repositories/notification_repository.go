package repositories

import (
	"github.com/fashion-shop/database"
	"github.com/fashion-shop/models"
	"gorm.io/gorm"
)

// NotificationRepository handles database operations for notifications
type NotificationRepository struct {
	db *gorm.DB
}

// NewNotificationRepository creates a new notification repository
func NewNotificationRepository() *NotificationRepository {
	return &NotificationRepository{
		db: database.DB,
	}
}

// Create creates a new notification
func (r *NotificationRepository) Create(notification *models.Notification) error {
	return r.db.Create(notification).Error
}

// FindByID finds a notification by ID
func (r *NotificationRepository) FindByID(id string) (*models.Notification, error) {
	var notification models.Notification
	err := r.db.First(&notification, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &notification, nil
}

// FindByUserID finds notifications by user ID with pagination
func (r *NotificationRepository) FindByUserID(userID string, page, limit int) ([]models.Notification, int64, error) {
	var notifications []models.Notification
	var count int64

	offset := (page - 1) * limit

	// Get count
	err := r.db.Model(&models.Notification{}).Where("user_id = ?", userID).Count(&count).Error
	if err != nil {
		return nil, 0, err
	}

	// Get notifications
	err = r.db.Where("user_id = ?", userID).
		Order("created_at DESC").
		Offset(offset).Limit(limit).Find(&notifications).Error
	if err != nil {
		return nil, 0, err
	}

	return notifications, count, nil
}

// MarkAsRead marks a notification as read
func (r *NotificationRepository) MarkAsRead(id string) error {
	return r.db.Model(&models.Notification{}).Where("id = ?", id).Update("is_read", true).Error
}

// MarkAllAsRead marks all notifications for a user as read
func (r *NotificationRepository) MarkAllAsRead(userID string) error {
	return r.db.Model(&models.Notification{}).Where("user_id = ?", userID).Update("is_read", true).Error
}

// CountUnread counts unread notifications for a user
func (r *NotificationRepository) CountUnread(userID string) (int64, error) {
	var count int64
	err := r.db.Model(&models.Notification{}).Where("user_id = ? AND is_read = ?", userID, false).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}
