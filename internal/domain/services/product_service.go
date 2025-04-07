package services

import (
	"errors"

	"clothing-shop-api/internal/domain/models"
	"clothing-shop-api/internal/repository"
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

func (s *ProductService) AddProduct(product *models.Product) error {
	if product.Name == "" || product.Price <= 0 {
		return errors.New("invalid product data")
	}
	return s.repo.Create(product)
}

func (s *ProductService) UpdateProduct(product *models.Product) error {
	if product.ID == 0 || product.Name == "" || product.Price <= 0 {
		return errors.New("invalid product data")
	}
	return s.repo.Update(product)
}

func (s *ProductService) DeleteProduct(id uint) error {
	if id == 0 {
		return errors.New("invalid product ID")
	}
	return s.repo.Delete(id)
}

func (s *ProductService) GetProductByID(id uint) (*models.Product, error) {
	if id == 0 {
		return nil, errors.New("invalid product ID")
	}
	return s.repo.FindByID(id)
}

func (s *ProductService) ListProducts(filter ProductFilter) ([]ProductResponse, int, error) {
	if filter.Page <= 0 {
		filter.Page = 1
	}
	if filter.PageSize <= 0 {
		filter.PageSize = 10
	}

	products, total, err := s.repo.FindAll(filter)
	if err != nil {
		return nil, 0, err
	}

	var response []ProductResponse
	for _, product := range products {
		// Get product variants
		variants, err := s.variantRepo.FindByProductID(product.ID)
		if err != nil {
			continue
		}

		// Extract available sizes and colors
		sizeMap := make(map[string]struct{})
		colorMap := make(map[string]struct{})
		for _, variant := range variants {
			sizeMap[variant.Size] = struct{}{}
			colorMap[variant.Color] = struct{}{}
		}

		var availableSizes []string
		for size := range sizeMap {
			availableSizes = append(availableSizes, size)
		}

		var availableColors []string
		for color := range colorMap {
			availableColors = append(availableColors, color)
		}

		// Get primary image
		primaryImage, err := s.imageRepo.FindPrimaryByProductID(product.ID)
		if err != nil {
			continue
		}

		response = append(response, ProductResponse{
			Product:         product,
			AvailableSizes:  availableSizes,
			AvailableColors: availableColors,
			MinPrice:        product.MinPrice,
			MaxPrice:        product.MaxPrice,
			PrimaryImage:    primaryImage,
		})
	}

	return response, total, nil
}
