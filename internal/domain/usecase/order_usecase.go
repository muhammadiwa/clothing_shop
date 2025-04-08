package usecase

import (
	"context"
	"time"

	"fashion-shop/internal/domain/entity"
)

// OrderUseCase defines the interface for order business logic
type OrderUseCase interface {
	CreateOrder(ctx context.Context, userID uint, addressID uint, paymentMethod entity.PaymentMethod, shippingMethod string, notes string) (*entity.Order, error)
	GetOrderByID(ctx context.Context, id uint, userID uint) (*entity.Order, error)
	GetOrderByNumber(ctx context.Context, orderNumber string, userID uint) (*entity.Order, error)
	GetUserOrders(ctx context.Context, userID uint, page, limit int) ([]*entity.Order, int64, error)
	CancelOrder(ctx context.Context, id uint, userID uint) error

	// Admin functions
	GetAllOrders(ctx context.Context, filter map[string]interface{}, page, limit int) ([]*entity.Order, int64, error)
	UpdateOrderStatus(ctx context.Context, id uint, status entity.OrderStatus) error
	UpdateShippingInfo(ctx context.Context, id uint, trackingNumber string) error
	GetSalesReport(ctx context.Context, startDate, endDate time.Time) ([]*entity.Order, float64, error)
}

// PaymentUseCase defines the interface for payment business logic
type PaymentUseCase interface {
	ProcessPayment(ctx context.Context, orderID uint, paymentMethod entity.PaymentMethod) (*entity.Payment, error)
	GetPaymentByID(ctx context.Context, id uint) (*entity.Payment, error)
	GetPaymentByOrderID(ctx context.Context, orderID uint) (*entity.Payment, error)
	HandlePaymentCallback(ctx context.Context, transactionID string, status string) error
	RefundPayment(ctx context.Context, paymentID uint, amount float64, reason string) error
}

// CartUseCase defines the interface for cart business logic
type CartUseCase interface {
	GetCart(ctx context.Context, userID uint) (*entity.Cart, error)
	AddToCart(ctx context.Context, userID, productID, variantID uint, quantity int) error
	UpdateCartItem(ctx context.Context, userID, itemID uint, quantity int) error
	RemoveFromCart(ctx context.Context, userID, itemID uint) error
	ClearCart(ctx context.Context, userID uint) error
	GetCartTotals(ctx context.Context, userID uint) (int, float64, error) // returns total items, total amount, error
}

// WishlistUseCase defines the interface for wishlist business logic
type WishlistUseCase interface {
	GetWishlist(ctx context.Context, userID uint, page, limit int) (*entity.Wishlist, []*entity.WishlistItem, int64, error)
	AddToWishlist(ctx context.Context, userID, productID uint) error
	RemoveFromWishlist(ctx context.Context, userID, itemID uint) error
	IsInWishlist(ctx context.Context, userID, productID uint) (bool, error)
}

// ShippingUseCase defines the interface for shipping business logic
type ShippingUseCase interface {
	GetProvinces(ctx context.Context) ([]map[string]interface{}, error)
	GetCities(ctx context.Context, provinceID string) ([]map[string]interface{}, error)
	CalculateShipping(ctx context.Context, origin, destination, weight int, courier string) ([]map[string]interface{}, error)
	TrackShipment(ctx context.Context, waybill, courier string) (map[string]interface{}, error)
}
