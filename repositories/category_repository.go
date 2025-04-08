package repositories

import (
	"github.com/fashion-shop/database"
	"github.com/fashion-shop/models"
	"gorm.io/gorm"
)

// CategoryRepository handles database operations for categories
type CategoryRepository struct {
	db *gorm.DB
}

// NewCategoryRepository creates a new category repository
func NewCategoryRepository() *CategoryRepository {
	return &CategoryRepository{
		db: database.DB,
	}
}

// Create creates a new category
func (r *CategoryRepository) Create(category *models.Category) error {
	return r.db.Create(category).Error
}

// FindByID finds a category by ID
func (r *CategoryRepository) FindByID(id string) (*models.Category, error) {
	var category models.Category
	err := r.db.Preload("Parent").First(&category, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &category, nil
}

// Update updates a category
func (r *CategoryRepository) Update(category *models.Category) error {
	return r.db.Save(category).Error
}

// Delete deletes a category
func (r *CategoryRepository) Delete(id string) error {
	return r.db.Delete(&models.Category{}, "id = ?", id).Error
}

// FindAll finds all categories
func (r *CategoryRepository) FindAll() ([]models.Category, error) {
	var categories []models.Category
	err := r.db.Preload("Parent").Find(&categories).Error
	if err != nil {
		return nil, err
	}
	return categories, nil
}

// FindByParentID finds categories by parent ID
func (r *CategoryRepository) FindByParentID(parentID *string) ([]models.Category, error) {
	var categories []models.Category
	query := r.db.Preload("Parent")

	if parentID == nil {
		query = query.Where("parent_id IS NULL")
	} else {
		query = query.Where("parent_id = ?", *parentID)
	}

	err := query.Find(&categories).Error
	if err != nil {
		return nil, err
	}
	return categories, nil
}

// GetCategoryProducts gets products by category ID with pagination
func (r *CategoryRepository) GetCategoryProducts(categoryID string, page, limit int) ([]models.Product, int64, error) {
	var products []models.Product
	var count int64

	offset := (page - 1) * limit

	// Get count
	err := r.db.Model(&models.Product{}).Where("category_id = ?", categoryID).Count(&count).Error
	if err != nil {
		return nil, 0, err
	}

	// Get products
	err = r.db.Where("category_id = ?", categoryID).
		Preload("Category").Preload("Images").Preload("Variants").Preload("Tags").
		Offset(offset).Limit(limit).Find(&products).Error
	if err != nil {
		return nil, 0, err
	}

	return products, count, nil
}
