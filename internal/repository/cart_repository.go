package repository

import (
	"clothing-shop-api/internal/domain/models"
)

type CartRepository interface {
	AddItem(cart *models.Cart) error
	UpdateItem(cart *models.Cart) error
	RemoveItem(id uint) error
	GetUserCart(userID uint) ([]models.Cart, error)
	GetCartItem(id uint) (*models.Cart, error)
	ClearCart(userID uint) error
}
