package repository

import (
	"clothing-shop-api/internal/domain/models"
)

type WishlistRepository interface {
	Add(wishlist *models.Wishlist) error
	Remove(id uint) error
	GetUserWishlist(userID uint) ([]models.Wishlist, error)
	FindByUserAndProduct(userID, productID uint) (*models.Wishlist, error)
}
