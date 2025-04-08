package repositories

import (
	"github.com/fashion-shop/database"
	"github.com/fashion-shop/models"
	"gorm.io/gorm"
)

// WishlistRepository handles database operations for wishlists
type WishlistRepository struct {
	db *gorm.DB
}

// NewWishlistRepository creates a new wishlist repository
func NewWishlistRepository() *WishlistRepository {
	return &WishlistRepository{
		db: database.DB,
	}
}

// FindByUserID finds a wishlist by user ID
func (r *WishlistRepository) FindByUserID(userID string) (*models.Wishlist, error) {
	var wishlist models.Wishlist
	err := r.db.Preload("Items.Product").Preload("Items.Variant").
		First(&wishlist, "user_id = ?", userID).Error
	if err != nil {
		return nil, err
	}
	return &wishlist, nil
}

// CreateWishlist creates a new wishlist
func (r *WishlistRepository) CreateWishlist(wishlist *models.Wishlist) error {
	return r.db.Create(wishlist).Error
}

// AddItem adds an item to a wishlist
func (r *WishlistRepository) AddItem(item *models.WishlistItem) error {
	return r.db.Create(item).Error
}

// DeleteItem deletes a wishlist item
func (r *WishlistRepository) DeleteItem(id string) error {
	return r.db.Delete(&models.WishlistItem{}, "id = ?", id).Error
}

// FindItemByProductAndVariant finds a wishlist item by product ID and variant ID
func (r *WishlistRepository) FindItemByProductAndVariant(wishlistID, productID string, variantID *string) (*models.WishlistItem, error) {
	var item models.WishlistItem
	query := r.db.Where("wishlist_id = ? AND product_id = ?", wishlistID, productID)

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
