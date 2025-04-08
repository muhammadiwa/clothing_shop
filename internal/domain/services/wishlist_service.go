// package services

// import (
// 	"errors"

// 	"clothing-shop-api/internal/domain/models"
// 	"clothing-shop-api/internal/repository"
// )

// type WishlistService struct {
// 	wishlistRepo repository.WishlistRepository
// 	productRepo  repository.ProductRepository
// }

// func NewWishlistService(
// 	wishlistRepo repository.WishlistRepository,
// 	productRepo repository.ProductRepository,
// ) *WishlistService {
// 	return &WishlistService{
// 		wishlistRepo: wishlistRepo,
// 		productRepo:  productRepo,
// 	}
// }

// func (s *WishlistService) AddToWishlist(userID, productID uint) error {
// 	// Check if product exists
// 	product, err := s.productRepo.FindByID(productID)
// 	if err != nil {
// 		return err
// 	}
// 	if product == nil {
// 		return errors.New("product not found")
// 	}

// 	// Add to wishlist
// 	wishlistItem := &models.Wishlist{
// 		UserID:    userID,
// 		ProductID: productID,
// 	}

// 	return s.wishlistRepo.Add(wishlistItem)
// }

// func (s *WishlistService) RemoveFromWishlist(id uint) error {
// 	return s.wishlistRepo.Remove(id)
// }

// func (s *WishlistService) GetUserWishlist(userID uint) ([]models.Wishlist, error) {
// 	return s.wishlistRepo.GetUserWishlist(userID)
// }

// func (s *WishlistService) IsProductInWishlist(userID, productID uint) (bool, error) {
// 	item, err := s.wishlistRepo.FindByUserAndProduct(userID, productID)
// 	if err != nil {
// 		return false, err
// 	}
// 	return item != nil, nil
// }
