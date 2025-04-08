package usecase

import (
	"context"
	"mime/multipart"

	"fashion-shop/internal/domain/entity"
)

// ProductUseCase defines the interface for product business logic
type ProductUseCase interface {
	CreateProduct(ctx context.Context, product *entity.Product) (*entity.Product, error)
	GetProductByID(ctx context.Context, id uint) (*entity.Product, error)
	GetProductBySlug(ctx context.Context, slug string) (*entity.Product, error)
	UpdateProduct(ctx context.Context, id uint, product *entity.Product) (*entity.Product, error)
	DeleteProduct(ctx context.Context, id uint) error
	ListProducts(ctx context.Context, filter map[string]interface{}, sort string, page, limit int) ([]*entity.Product, int64, error)
	SearchProducts(ctx context.Context, keyword string, filter map[string]interface{}, sort string, page, limit int) ([]*entity.Product, int64, error)
	UploadProductImage(ctx context.Context, productID uint, file *multipart.FileHeader, isPrimary bool) (*entity.ProductImage, error)
	DeleteProductImage(ctx context.Context, id uint) error
	SetPrimaryImage(ctx context.Context, id uint, productID uint) error

	// Variant management
	AddVariant(ctx context.Context, productID uint, variant *entity.ProductVariant) (*entity.ProductVariant, error)
	UpdateVariant(ctx context.Context, id uint, variant *entity.ProductVariant) (*entity.ProductVariant, error)
	DeleteVariant(ctx context.Context, id uint) error
	UpdateStock(ctx context.Context, variantID uint, quantity int) error

	// Featured products
	GetBestSellers(ctx context.Context, limit int) ([]*entity.Product, error)
	GetNewArrivals(ctx context.Context, limit int) ([]*entity.Product, error)
	GetTopRated(ctx context.Context, limit int) ([]*entity.Product, error)

	// Bulk operations
	BulkUploadProducts(ctx context.Context, file *multipart.FileHeader) (int, error) // returns number of products created
}

// CategoryUseCase defines the interface for category business logic
type CategoryUseCase interface {
	CreateCategory(ctx context.Context, category *entity.Category) (*entity.Category, error)
	GetCategoryByID(ctx context.Context, id uint) (*entity.Category, error)
	GetCategoryBySlug(ctx context.Context, slug string) (*entity.Category, error)
	UpdateCategory(ctx context.Context, id uint, category *entity.Category) (*entity.Category, error)
	DeleteCategory(ctx context.Context, id uint) error
	ListCategories(ctx context.Context, parentID *uint) ([]*entity.Category, error)
	UploadCategoryImage(ctx context.Context, id uint, file *multipart.FileHeader) (*entity.Category, error)
}

// ReviewUseCase defines the interface for review business logic
type ReviewUseCase interface {
	CreateReview(ctx context.Context, userID, productID, orderID uint, rating int, comment string, images []*multipart.FileHeader) (*entity.Review, error)
	GetReviewByID(ctx context.Context, id uint) (*entity.Review, error)
	GetProductReviews(ctx context.Context, productID uint, page, limit int) ([]*entity.Review, int64, error)
	GetUserReviews(ctx context.Context, userID uint, page, limit int) ([]*entity.Review, int64, error)
	UpdateReview(ctx context.Context, id, userID uint, rating int, comment string) (*entity.Review, error)
	DeleteReview(ctx context.Context, id, userID uint) error
	GetAverageRating(ctx context.Context, productID uint) (float64, error)
}
