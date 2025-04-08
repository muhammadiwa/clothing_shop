package services

import (
	"clothing-shop-api/internal/domain/models"
	"clothing-shop-api/internal/interfaces"
	"errors"
)

var (
	ErrProductNotFound = errors.New("product not found")
	ErrInvalidProduct  = errors.New("invalid product data")
)

type ProductService struct {
	repo         interfaces.ProductRepository
	categoryRepo interfaces.CategoryRepository
	variantRepo  interfaces.ProductVariantRepository
	imageRepo    interfaces.ProductImageRepository
}

func NewProductService(
	repo interfaces.ProductRepository,
	categoryRepo interfaces.CategoryRepository,
	variantRepo interfaces.ProductVariantRepository,
	imageRepo interfaces.ProductImageRepository,
) *ProductService {
	return &ProductService{
		repo:         repo,
		categoryRepo: categoryRepo,
		variantRepo:  variantRepo,
		imageRepo:    imageRepo,
	}
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
	if err := s.repo.Create(product); err != nil {
		return err
	}

	// Create variants if any
	for _, variant := range product.Variants {
		variant.ProductID = product.ID
		if err := s.variantRepo.Create(&variant); err != nil {
			return err
		}
	}

	// Create images if any
	for _, image := range product.Images {
		image.ProductID = product.ID
		if err := s.imageRepo.Create(&image); err != nil {
			return err
		}
	}

	return nil
}

// Ubah return type ke interfaces.ProductResponse
func (s *ProductService) GetProduct(id uint) (*interfaces.ProductResponse, error) {
	// Get product with basic info
	product, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if product == nil {
		return nil, ErrProductNotFound
	}

	// Get variants
	variants, err := s.variantRepo.FindByProductID(id)
	if err != nil {
		return nil, err
	}

	// Get primary image
	primaryImage, err := s.imageRepo.FindPrimaryByProductID(id)
	if err != nil {
		return nil, err
	}

	// Create response with additional info
	response := &interfaces.ProductResponse{
		Product:      product,
		PrimaryImage: primaryImage,
	}

	// Calculate price range and collect available sizes/colors
	sizeMap := make(map[string]bool)
	colorMap := make(map[string]bool)

	if len(variants) > 0 {
		response.MinPrice = variants[0].Price
		response.MaxPrice = variants[0].Price

		for _, variant := range variants {
			if variant.Price < response.MinPrice {
				response.MinPrice = variant.Price
			}
			if variant.Price > response.MaxPrice {
				response.MaxPrice = variant.Price
			}

			if variant.Size != "" {
				sizeMap[variant.Size] = true
			}
			if variant.Color != "" {
				colorMap[variant.Color] = true
			}
		}
	}

	// Convert maps to slices
	for size := range sizeMap {
		response.AvailableSizes = append(response.AvailableSizes, size)
	}
	for color := range colorMap {
		response.AvailableColors = append(response.AvailableColors, color)
	}

	return response, nil
}

func (s *ProductService) ListProducts(filter interfaces.ProductFilter) ([]*interfaces.ProductResponse, int, error) {
	products, total, err := s.repo.FindAll(filter)
	if err != nil {
		return nil, 0, err
	}

	responses := make([]*interfaces.ProductResponse, len(products))
	for i, product := range products {
		response, err := s.GetProduct(product.ID)
		if err != nil {
			return nil, 0, err
		}
		responses[i] = response
	}

	return responses, total, nil
}

func (s *ProductService) UpdateProduct(product *models.Product) error {
	if err := s.validateProduct(product); err != nil {
		return err
	}

	existing, err := s.repo.FindByID(product.ID)
	if err != nil {
		return err
	}
	if existing == nil {
		return ErrProductNotFound
	}

	category, err := s.categoryRepo.FindByID(product.CategoryID)
	if err != nil {
		return err
	}
	if category == nil {
		return errors.New("category not found")
	}

	if err := s.repo.Update(product); err != nil {
		return err
	}

	// Update variants
	for _, variant := range product.Variants {
		variant.ProductID = product.ID
		if variant.ID == 0 {
			if err := s.variantRepo.Create(&variant); err != nil {
				return err
			}
		} else {
			if err := s.variantRepo.Update(&variant); err != nil {
				return err
			}
		}
	}

	// Update images
	for _, image := range product.Images {
		image.ProductID = product.ID
		if image.ID == 0 {
			if err := s.imageRepo.Create(&image); err != nil {
				return err
			}
		}

	}

	return nil
}

func (s *ProductService) DeleteProduct(id uint) error {
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
