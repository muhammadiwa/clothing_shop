package interfaces

import "clothing-shop-api/internal/domain/models"

// ProductFilter defines filter parameters for product queries
type ProductFilter struct {
	CategoryID *uint    `form:"category_id"`
	MinPrice   *float64 `form:"min_price"`
	MaxPrice   *float64 `form:"max_price"`
	Search     string   `form:"search"`
	SortBy     string   `form:"sort_by"`
	Page       int      `form:"page"`
	PageSize   int      `form:"page_size"`
}

type UserRepository interface {
	Create(user *models.User) error
	FindByID(id uint) (*models.User, error)
	FindByEmail(email string) (*models.User, error)
	Update(user *models.User) error
	Delete(id uint) error
	FindByVerificationToken(token string) (*models.User, error)
	FindByResetToken(token string) (*models.User, error)
}

type OrderRepository interface {
	Create(order *models.Order) error
	FindByID(id uint) (*models.Order, error)
	FindByUserID(userID uint) ([]models.Order, error)
	FindAll() ([]models.Order, error)
	Update(order *models.Order) error
	Delete(id uint) error
}

type PaymentRepository interface {
	Create(payment *models.Payment) error
	FindByID(id uint) (*models.Payment, error)
	FindByOrderID(orderID uint) (*models.Payment, error)
	FindByTransactionID(transactionID string) (*models.Payment, error)
	Update(payment *models.Payment) error
	Delete(id uint) error
}

type CategoryRepository interface {
	Create(category *models.Category) error
	FindByID(id uint) (*models.Category, error)
	FindAll() ([]*models.Category, error)
	Update(category *models.Category) error
	Delete(id uint) error
}

type ProductRepository interface {
	Create(product *models.Product) error
	FindByID(id uint) (*models.Product, error)
	FindAll(filter ProductFilter) ([]*models.Product, int, error)
	Update(product *models.Product) error
	Delete(id uint) error
}

// CartRepository defines cart operations
type CartRepository interface {
	AddItem(cart *models.Cart) error
	UpdateItem(cart *models.Cart) error
	RemoveItem(id uint) error
	GetUserCart(userID uint) ([]models.Cart, error)
	GetCartItem(id uint) (*models.Cart, error)
	ClearCart(userID uint) error
}

// ProductVariantRepository defines product variant operations
type ProductVariantRepository interface {
	FindByID(id uint) (*models.ProductVariant, error)
	Create(variant *models.ProductVariant) error
	Update(variant *models.ProductVariant) error
	Delete(id uint) error
	FindByProductID(productID uint) ([]models.ProductVariant, error)
	UpdateStock(id uint, quantity int) error
}
