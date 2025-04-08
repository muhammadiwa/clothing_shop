package repository

import (
	"context"

	"fashion-shop/internal/domain/entity"
)

// ProductRepository defines the interface for product data access
type ProductRepository interface {
	Create(ctx context.Context, product *entity.Product) error
	GetByID(ctx context.Context, id uint) (*entity.Product, error)
	GetBySlug(ctx context.Context, slug string) (*entity.Product, error)
	Update(ctx context.Context, product *entity.Product) error
	Delete(ctx context.Context, id uint) error
	List(ctx context.Context, filter map[string]interface{}, sort string, offset, limit int) ([]*entity.Product, int64, error)
	Search(ctx context.Context, keyword string, filter map[string]interface{}, sort string, offset, limit int) ([]*entity.Product, int64, error)
	UpdateStock(ctx context.Context, variantID uint, quantity int) error
	GetBestSellers(ctx context.Context, limit int) ([]*entity.Product, error)
	GetNewArrivals(ctx context.Context, limit int) ([]*entity.Product, error)
	GetTopRated(ctx context.Context, limit int) ([]*entity.Product, error)
}

// CategoryRepository defines the interface for category data access
type CategoryRepository interface {
	Create(ctx context.Context, category *entity.Category) error
	GetByID(ctx context.Context, id uint) (*entity.Category, error)
	GetBySlug(ctx context.Context, slug string) (*entity.Category, error)
	Update(ctx context.Context, category *entity.Category) error
	Delete(ctx context.Context, id uint) error
	List(ctx context.Context, parentID *uint) ([]*entity.Category, error)
}

// ProductImageRepository defines the interface for product image data access
type ProductImageRepository interface {
	Create(ctx context.Context, image *entity.ProductImage) error
	GetByID(ctx context.Context, id uint) (*entity.ProductImage, error)
	GetByProductID(ctx context.Context, productID uint) ([]*entity.ProductImage, error)
	Update(ctx context.Context, image *entity.ProductImage) error
	Delete(ctx context.Context, id uint) error
	SetPrimary(ctx context.Context, id uint, productID uint) error
}

// ProductVariantRepository defines the interface for product variant data access
type ProductVariantRepository interface {
	Create(ctx context.Context, variant *entity.ProductVariant) error
	GetByID(ctx context.Context, id uint) (*entity.ProductVariant, error)
	GetByProductID(ctx context.Context, productID uint) ([]*entity.ProductVariant, error)
	GetBySKU(ctx context.Context, sku string) (*entity.ProductVariant, error)
	Update(ctx context.Context, variant *entity.ProductVariant) error
	Delete(ctx context.Context, id uint) error
	UpdateStock(ctx context.Context, id uint, quantity int) error
}

// ReviewRepository defines the interface for review data access
type ReviewRepository interface {
	Create(ctx context.Context, review *entity.Review) error
	GetByID(ctx context.Context, id uint) (*entity.Review, error)
	GetByProductID(ctx context.Context, productID uint, offset, limit int) ([]*entity.Review, int64, error)
	GetByUserID(ctx context.Context, userID uint, offset, limit int) ([]*entity.Review, int64, error)
	Update(ctx context.Context, review *entity.Review) error
	Delete(ctx context.Context,  error)
	Update(ctx context.Context, review *entity.Review) error
	Delete(ctx context.Context, id uint) error
	GetAverageRatingByProductID(ctx context.Context, productID uint) (float64, error)
}

// TagRepository defines the interface for tag data access
type TagRepository interface {
	Create(ctx context.Context, tag *entity.Tag) error
	GetByID(ctx context.Context, id uint) (*entity.Tag, error)
	GetByName(ctx context.Context, name string) (*entity.Tag, error)
	List(ctx context.Context) ([]*entity.Tag, error)
	Update(ctx context.Context, tag *entity.Tag) error
	Delete(ctx context.Context, id uint) error
}
