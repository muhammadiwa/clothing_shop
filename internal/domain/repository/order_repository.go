package repository

import (
	"context"
	"time"

	"fashion-shop/internal/domain/entity"
)

// OrderRepository defines the interface for order data access
type OrderRepository interface {
	Create(ctx context.Context, order *entity.Order) error
	GetByID(ctx context.Context, id uint) (*entity.Order, error)
	GetByOrderNumber(ctx context.Context, orderNumber string) (*entity.Order, error)
	GetByUserID(ctx context.Context, userID uint, offset, limit int) ([]*entity.Order, int64, error)
	Update(ctx context.Context, order *entity.Order) error
	UpdateStatus(ctx context.Context, id uint, status entity.OrderStatus) error
	Delete(ctx context.Context, id uint) error
	List(ctx context.Context, filter map[string]interface{}, offset, limit int) ([]*entity.Order, int64, error)
	GetSalesReport(ctx context.Context, startDate, endDate time.Time) ([]*entity.Order, float64, error)
}

// OrderItemRepository defines the interface for order item data access
type OrderItemRepository interface {
	Create(ctx context.Context, item *entity.OrderItem) error
	GetByID(ctx context.Context, id uint) (*entity.OrderItem, error)
	GetByOrderID(ctx context.Context, orderID uint) ([]*entity.OrderItem, error)
	Update(ctx context.Context, item *entity.OrderItem) error
	Delete(ctx context.Context, id uint) error
}

// PaymentRepository defines the interface for payment data access
type PaymentRepository interface {
	Create(ctx context.Context, payment *entity.Payment) error
	GetByID(ctx context.Context, id uint) (*entity.Payment, error)
	GetByOrderID(ctx context.Context, orderID uint) (*entity.Payment, error)
	GetByTransactionID(ctx context.Context, transactionID string) (*entity.Payment, error)
	Update(ctx context.Context, payment *entity.Payment) error
	UpdateStatus(ctx context.Context, id uint, status entity.PaymentStatus) error
	List(ctx context.Context, filter map[string]interface{}, offset, limit int) ([]*entity.Payment, int64, error)
}

// CartRepository defines the interface for cart data access
type CartRepository interface {
	GetOrCreate(ctx context.Context, userID uint) (*entity.Cart, error)
	GetByID(ctx context.Context, id uint) (*entity.Cart, error)
	GetByUserID(ctx context.Context, userID uint) (*entity.Cart, error)
	AddItem(ctx context.Context, cartID uint, item *entity.CartItem) error
	UpdateItem(ctx context.Context, item *entity.CartItem) error
	RemoveItem(ctx context.Context, itemID uint) error
	ClearCart(ctx context.Context, cartID uint) error
	GetTotalItems(ctx context.Context, cartID uint) (int, error)
	GetTotalAmount(ctx context.Context, cartID uint) (float64, error)
}

// WishlistRepository defines the interface for wishlist data access
type WishlistRepository interface {
	GetOrCreate(ctx context.Context, userID uint) (*entity.Wishlist, error)
	GetByID(ctx context.Context, id uint) (*entity.Wishlist, error)
	GetByUserID(ctx context.Context, userID uint) (*entity.Wishlist, error)
	AddItem(ctx context.Context, wishlistID uint, productID uint) error
	RemoveItem(ctx context.Context, itemID uint) error
	IsProductInWishlist(ctx context.Context, wishlistID uint, productID uint) (bool, error)
	GetItems(ctx context.Context, wishlistID uint, offset, limit int) ([]*entity.WishlistItem, int64, error)
}
