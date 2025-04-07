package interfaces

import "clothing-shop-api/internal/domain/models"

// ProductFilter defines filter parameters for product queries
type ProductFilter struct {
	CategoryID *uint    `form:"category_id"`
	MinPrice   *float64 `form:"min_price"`
	MaxPrice   *float64 `form:"max_price"`
	Search     string   `form:"search"`
	SortBy     string   `form:"sort_by"`
	Page       int      `form:"page"`
	PageSize   int      `form:"page_size"`
}

// ProductResponse defines the API response structure for products
type ProductResponse struct {
	*models.Product
	AvailableSizes  []string             `json:"available_sizes"`
	AvailableColors []string             `json:"available_colors"`
	MinPrice        float64              `json:"min_price"`
	MaxPrice        float64              `json:"max_price"`
	PrimaryImage    *models.ProductImage `json:"primary_image"`
}

// Repository interfaces
type ProductRepository interface {
	Create(product *models.Product) error
	FindByID(id uint) (*models.Product, error)
	FindAll(filter ProductFilter) ([]*models.Product, int, error)
	Update(product *models.Product) error
	Delete(id uint) error
}

type CategoryRepository interface {
	Create(category *models.Category) error
	FindByID(id uint) (*models.Category, error)
	FindAll() ([]*models.Category, error)
	Update(category *models.Category) error
	Delete(id uint) error
}

type ProductVariantRepository interface {
	Create(variant *models.ProductVariant) error
	FindByID(id uint) (*models.ProductVariant, error)
	FindByProductID(productID uint) ([]models.ProductVariant, error)
	Update(variant *models.ProductVariant) error
	Delete(id uint) error
	UpdateStock(id uint, quantity int) error
}

type ProductImageRepository interface {
	Create(image *models.ProductImage) error
	FindByID(id uint) (*models.ProductImage, error)
	FindByProductID(productID uint) ([]models.ProductImage, error)
	FindPrimaryByProductID(productID uint) (*models.ProductImage, error)
	Delete(id uint) error
	SetPrimary(id uint) error
}
