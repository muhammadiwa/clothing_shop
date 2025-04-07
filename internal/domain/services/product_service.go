package services

import (
	"clothing-shop-api/internal/domain/models"
	"clothing-shop-api/internal/repository"
	"errors"
)

var (
	ErrProductNotFound = errors.New("product not found")
	ErrInvalidProduct  = errors.New("invalid product data")
)

type ProductService struct {
	repo         repository.ProductRepository
	categoryRepo repository.CategoryRepository
	variantRepo  repository.ProductVariantRepository
	imageRepo    repository.ProductImageRepository
}

func NewProductService(
	repo repository.ProductRepository,
	categoryRepo repository.CategoryRepository,
	variantRepo repository.ProductVariantRepository,
	imageRepo repository.ProductImageRepository,
) *ProductService {
	return &ProductService{
		repo:         repo,
		categoryRepo: categoryRepo,
		variantRepo:  variantRepo,
		imageRepo:    imageRepo,
	}
}

type ProductFilter struct {
	CategoryID *uint    `form:"category_id"`
	MinPrice   *float64 `form:"min_price"`
	MaxPrice   *float64 `form:"max_price"`
	Search     string   `form:"search"`
	SortBy     string   `form:"sort_by"` // name_asc, name_desc, price_asc, price_desc, newest
	Page       int      `form:"page"`
	PageSize   int      `form:"page_size"`
}

type ProductResponse struct {
	*models.Product
	AvailableSizes  []string             `json:"available_sizes"`
	AvailableColors []string             `json:"available_colors"`
	MinPrice        float64              `json:"min_price"`
	MaxPrice        float64              `json:"max_price"`
	PrimaryImage    *models.ProductImage `json:"primary_image"`
}

func (s *ProductService) CreateProduct(product *models.Product) error {
	// Validate product data
	if err := s.validateProduct(product); err != nil {
		return err
	}

	// Check if category exists
	category, err := s.categoryRepo.FindByID(product.CategoryID)
	if err != nil {
		return err
	}
	if category == nil {
		return errors.New("category not found")
	}

	// Create product
	return s.repo.Create(product)
}

func (s *ProductService) GetProduct(id uint) (*models.Product, error) {
	product, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if product == nil {
		return nil, ErrProductNotFound
	}
	return product, nil
}

func (s *ProductService) ListProducts(filter ProductFilter) ([]*models.Product, int, error) {
	return s.repo.FindAll(filter)
}

func (s *ProductService) UpdateProduct(product *models.Product) error {
	// Validate product data
	if err := s.validateProduct(product); err != nil {
		return err
	}

	// Check if product exists
	existing, err := s.repo.FindByID(product.ID)
	if err != nil {
		return err
	}
	if existing == nil {
		return ErrProductNotFound
	}

	// Check if category exists
	category, err := s.categoryRepo.FindByID(product.CategoryID)
	if err != nil {
		return err
	}
	if category == nil {
		return errors.New("category not found")
	}

	return s.repo.Update(product)
}

func (s *ProductService) DeleteProduct(id uint) error {
	// Check if product exists
	product, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}
	if product == nil {
		return ErrProductNotFound
	}

	return s.repo.Delete(id)
}

// Helper methods
func (s *ProductService) validateProduct(product *models.Product) error {
	if product.Name == "" {
		return errors.New("product name is required")
	}
	if product.BasePrice <= 0 {
		return errors.New("product price must be greater than zero")
	}
	if product.CategoryID == 0 {
		return errors.New("product category is required")
	}
	if product.Weight <= 0 {
		return errors.New("product weight must be greater than zero")
	}
	return nil
}
