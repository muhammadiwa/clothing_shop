package repository

import "clothing-shop-api/internal/domain/models"

// ProductImageRepository defines product image operations
type ProductImageRepository interface {
	Create(image *models.ProductImage) error
	FindByID(id uint) (*models.ProductImage, error)
	FindByProductID(productID uint) ([]models.ProductImage, error)
	FindPrimaryByProductID(productID uint) (*models.ProductImage, error)
	Delete(id uint) error
	SetPrimary(id uint) error
}
