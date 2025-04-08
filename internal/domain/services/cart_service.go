// package services

// import (
// 	"clothing-shop-api/internal/domain/interfaces"
// 	"clothing-shop-api/internal/domain/models"
// 	"errors"
// )

// var (
// 	ErrCartItemNotFound  = errors.New("cart item not found")
// 	ErrInsufficientStock = errors.New("insufficient stock")
// )

// type CartService struct {
// 	cartRepo    interfaces.CartRepository
// 	variantRepo interfaces.ProductVariantRepository
// }

// func NewCartService(
// 	cartRepo interfaces.CartRepository,
// 	variantRepo interfaces.ProductVariantRepository,
// ) *CartService {
// 	return &CartService{
// 		cartRepo:    cartRepo,
// 		variantRepo: variantRepo,
// 	}
// }

// func (s *CartService) AddToCart(userID, productVariantID uint, quantity int) error {
// 	// Check if product variant exists and has sufficient stock
// 	variant, err := s.variantRepo.FindByID(productVariantID)
// 	if err != nil {
// 		return err
// 	}
// 	if variant == nil {
// 		return errors.New("product variant not found")
// 	}

// 	if variant.Stock < quantity {
// 		return ErrInsufficientStock
// 	}

// 	// Add to cart
// 	cartItem := &models.Cart{
// 		UserID:           userID,
// 		ProductVariantID: productVariantID,
// 		Quantity:         quantity,
// 	}

// 	return s.cartRepo.AddItem(cartItem)
// }

// func (s *CartService) UpdateCartItem(id uint, quantity int) error {
// 	// Get cart item
// 	cartItem, err := s.cartRepo.GetCartItem(id)
// 	if err != nil {
// 		return err
// 	}
// 	if cartItem == nil {
// 		return ErrCartItemNotFound
// 	}

// 	// Check if product variant has sufficient stock
// 	variant, err := s.variantRepo.FindByID(cartItem.ProductVariantID)
// 	if err != nil {
// 		return err
// 	}
// 	if variant == nil {
// 		return errors.New("product variant not found")
// 	}

// 	if variant.Stock < quantity {
// 		return ErrInsufficientStock
// 	}

// 	// Update cart item
// 	cartItem.Quantity = quantity
// 	return s.cartRepo.UpdateItem(cartItem)
// }

// func (s *CartService) RemoveFromCart(id uint) error {
// 	return s.cartRepo.RemoveItem(id)
// }

// func (s *CartService) GetUserCart(userID uint) ([]models.Cart, error) {
// 	return s.cartRepo.GetUserCart(userID)
// }

// func (s *CartService) ClearCart(userID uint) error {
// 	return s.cartRepo.ClearCart(userID)
// }
