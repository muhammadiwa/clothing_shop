package repositories

import (
	"github.com/fashion-shop/database"
	"github.com/fashion-shop/models"
	"gorm.io/gorm"
)

// UserRepository handles database operations for users
type UserRepository struct {
	db *gorm.DB
}

// NewUserRepository creates a new user repository
func NewUserRepository() *UserRepository {
	return &UserRepository{
		db: database.DB,
	}
}

// Create creates a new user
func (r *UserRepository) Create(user *models.User) error {
	return r.db.Create(user).Error
}

// FindByID finds a user by ID
func (r *UserRepository) FindByID(id string) (*models.User, error) {
	var user models.User
	err := r.db.First(&user, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// FindByEmail finds a user by email
func (r *UserRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.db.First(&user, "email = ?", email).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Update updates a user
func (r *UserRepository) Update(user *models.User) error {
	return r.db.Save(user).Error
}

// Delete soft deletes a user
func (r *UserRepository) Delete(id string) error {
	return r.db.Delete(&models.User{}, "id = ?", id).Error
}

// FindAll finds all users with pagination
func (r *UserRepository) FindAll(page, limit int) ([]models.User, int64, error) {
	var users []models.User
	var count int64

	offset := (page - 1) * limit

	err := r.db.Model(&models.User{}).Count(&count).Error
	if err != nil {
		return nil, 0, err
	}

	err = r.db.Offset(offset).Limit(limit).Find(&users).Error
	if err != nil {
		return nil, 0, err
	}

	return users, count, nil
}

// UpdateRefreshToken updates a user's refresh token
func (r *UserRepository) UpdateRefreshToken(userID, token string, expiry interface{}) error {
	return r.db.Model(&models.User{}).Where("id = ?", userID).Updates(map[string]interface{}{
		"refresh_token": token,
		"token_expiry":  expiry,
	}).Error
}

// FindByRefreshToken finds a user by refresh token
func (r *UserRepository) FindByRefreshToken(token string) (*models.User, error) {
	var user models.User
	err := r.db.First(&user, "refresh_token = ?", token).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// LogActivity logs a user activity
func (r *UserRepository) LogActivity(activity *models.UserActivity) error {
	return r.db.Create(activity).Error
}

// GetUserActivities gets a user's activities with pagination
func (r *UserRepository) GetUserActivities(userID string, page, limit int) ([]models.UserActivity, int64, error) {
	var activities []models.UserActivity
	var count int64

	offset := (page - 1) * limit

	err := r.db.Model(&models.UserActivity{}).Where("user_id = ?", userID).Count(&count).Error
	if err != nil {
		return nil, 0, err
	}

	err = r.db.Where("user_id = ?", userID).Offset(offset).Limit(limit).Order("created_at DESC").Find(&activities).Error
	if err != nil {
		return nil, 0, err
	}

	return activities, count, nil
}
