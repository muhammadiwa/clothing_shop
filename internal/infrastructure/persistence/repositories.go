package persistence

import (
	"fashion-shop/internal/domain/repository"

	"gorm.io/gorm"
)

// Repositories holds all repository implementations
type Repositories struct {
	User           repository.UserRepository
	Address        repository.AddressRepository
	Product        repository.ProductRepository
	Category       repository.CategoryRepository
	ProductImage   repository.ProductImageRepository
	ProductVariant repository.ProductVariantRepository
	Review         repository.ReviewRepository
	Tag            repository.TagRepository
	Order          repository.OrderRepository
	OrderItem      repository.OrderItemRepository
	Payment        repository.PaymentRepository
	Cart           repository.CartRepository
	Wishlist       repository.WishlistRepository
	Notification   repository.NotificationRepository
	db             *gorm.DB
}

// NewRepositories creates a new Repositories instance
func NewRepositories(db *gorm.DB) *Repositories {
	return &Repositories{
		User:           NewUserRepository(db),
		Address:        NewAddressRepository(db),
		Product:        NewProductRepository(db),
		Category:       NewCategoryRepository(db),
		ProductImage:   NewProductImageRepository(db),
		ProductVariant: NewProductVariantRepository(db),
		Review:         NewReviewRepository(db),
		Tag:            NewTagRepository(db),
		Order:          NewOrderRepository(db),
		OrderItem:      NewOrderItemRepository(db),
		Payment:        NewPaymentRepository(db),
		Cart:           NewCartRepository(db),
		Wishlist:       NewWishlistRepository(db),
		Notification:   NewNotificationRepository(db),
		db:             db,
	}
}

// GetDB returns the database instance
func (r *Repositories) GetDB() *gorm.DB {
	return r.db
}
