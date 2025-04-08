package repositories

import (
	"github.com/fashion-shop/database"
	"github.com/fashion-shop/models"
	"gorm.io/gorm"
)

// ReviewRepository handles database operations for reviews
type ReviewRepository struct {
	db *gorm.DB
}

// NewReviewRepository creates a new review repository
func NewReviewRepository() *ReviewRepository {
	return &ReviewRepository{
		db: database.DB,
	}
}

// Create creates a new review
func (r *ReviewRepository) Create(review *models.Review) error {
	return r.db.Create(review).Error
}

// FindByID finds a review by ID
func (r *ReviewRepository) FindByID(id string) (*models.Review, error) {
	var review models.Review
	err := r.db.Preload("User").Preload("Product").First(&review, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &review, nil
}

// Update updates a review
func (r *ReviewRepository) Update(review *models.Review) error {
	return r.db.Save(review).Error
}

// Delete deletes a review
func (r *ReviewRepository) Delete(id string) error {
	return r.db.Delete(&models.Review{}, "id = ?", id).Error
}

// FindByProductID finds reviews by product ID with pagination
func (r *ReviewRepository) FindByProductID(productID string, page, limit int) ([]models.Review, int64, error) {
	var reviews []models.Review
	var count int64

	offset := (page - 1) * limit

	// Get count
	err := r.db.Model(&models.Review{}).Where("product_id = ?", productID).Count(&count).Error
	if err != nil {
		return nil, 0, err
	}

	// Get reviews
	err = r.db.Where("product_id = ?", productID).
		Preload("User").
		Order("created_at DESC").
		Offset(offset).Limit(limit).Find(&reviews).Error
	if err != nil {
		return nil, 0, err
	}

	return reviews, count, nil
}

// FindByUserID finds reviews by user ID with pagination
func (r *ReviewRepository) FindByUserID(userID string, page, limit int) ([]models.Review, int64, error) {
	var reviews []models.Review
	var count int64

	offset := (page - 1) * limit

	// Get count
	err := r.db.Model(&models.Review{}).Where("user_id = ?", userID).Count(&count).Error
	if err != nil {
		return nil, 0, err
	}

	// Get reviews
	err = r.db.Where("user_id = ?", userID).
		Preload("Product").
		Order("created_at DESC").
		Offset(offset).Limit(limit).Find(&reviews).Error
	if err != nil {
		return nil, 0, err
	}

	return reviews, count, nil
}

// GetAverageRating gets the average rating for a product
func (r *ReviewRepository) GetAverageRating(productID string) (float64, error) {
	var result struct {
		AvgRating float64
	}

	err := r.db.Model(&models.Review{}).
		Select("COALESCE(AVG(rating), 0) as avg_rating").
		Where("product_id = ?", productID).
		Scan(&result).Error

	if err != nil {
		return 0, err
	}

	return result.AvgRating, nil
}
