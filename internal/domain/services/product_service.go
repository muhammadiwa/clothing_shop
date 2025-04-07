package services

import (
	"errors"
	"clothing-shop-api/internal/domain/models"
	"clothing-shop-api/internal/repository"
)

type ProductService struct {
	repo repository.ProductRepository
}

func NewProductService(repo repository.ProductRepository) *ProductService {
	return &ProductService{repo: repo}
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

func (s *ProductService) ListProducts() ([]models.Product, error) {
	return s.repo.FindAll()
}