package repositories

import (
	"github.com/fashion-shop/database"
	"github.com/fashion-shop/models"
	"gorm.io/gorm"
)

// CartRepository handles database operations for carts
type CartRepository struct {
	db *gorm.DB
}

// NewCartRepository creates a new cart repository
func NewCartRepository() *CartRepository {
	return &CartRepository{
		db: database.DB,
	}
}

// FindByUserID finds a cart by user ID
func (r *CartRepository) FindByUserID(userID string) (*models.Cart, error) {
	var cart models.Cart
	err := r.db.Preload("Items.Product").Preload("Items.Variant").
		First(&cart, "user_id = ?", userID).Error
	if err != nil {
		return nil, err
	}
	return &cart, nil
}

// CreateCart creates a new cart
func (r *CartRepository) CreateCart(cart *models.Cart) error {
	return r.db.Create(cart).Error
}

// AddItem adds an item to a cart
func (r *CartRepository) AddItem(item *models.CartItem) error {
	return r.db.Create(item).Error
}

// UpdateItem updates a cart item
func (r *CartRepository) UpdateItem(item *models.CartItem) error {
	return r.db.Save(item).Error
}

// DeleteItem deletes a cart item
func (r *CartRepository) DeleteItem(id string) error {
	return r.db.Delete(&models.CartItem{}, "id = ?", id).Error
}

// ClearCart clears all items from a cart
func (r *CartRepository) ClearCart(cartID string) error {
	return r.db.Delete(&models.CartItem{}, "cart_id = ?", cartID).Error
}

// FindItemByProductAndVariant finds a cart item by product ID and variant ID
func (r *CartRepository) FindItemByProductAndVariant(cartID, productID string, variantID *string) (*models.CartItem, error) {
	var item models.CartItem
	query := r.db.Where("cart_id = ? AND product_id = ?", cartID, productID)

	if variantID == nil {
		query = query.Where("variant_id IS NULL")
	} else {
		query = query.Where("variant_id = ?", *variantID)
	}

	err := query.First(&item).Error
	if err != nil {
		return nil, err
	}
	return &item, nil
}
