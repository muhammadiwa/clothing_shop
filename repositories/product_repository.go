package repositories

import (
	"github.com/fashion-shop/database"
	"github.com/fashion-shop/models"
	"gorm.io/gorm"
)

// ProductRepository handles database operations for products
type ProductRepository struct {
	db *gorm.DB
}

// NewProductRepository creates a new product repository
func NewProductRepository() *ProductRepository {
	return &ProductRepository{
		db: database.DB,
	}
}

// Create creates a new product
func (r *ProductRepository) Create(product *models.Product) error {
	return r.db.Create(product).Error
}

// FindByID finds a product by ID
func (r *ProductRepository) FindByID(id string) (*models.Product, error) {
	var product models.Product
	err := r.db.Preload("Category").Preload("Images").Preload("Variants").Preload("Tags").First(&product, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

// Update updates a product
func (r *ProductRepository) Update(product *models.Product) error {
	return r.db.Save(product).Error
}

// Delete soft deletes a product
func (r *ProductRepository) Delete(id string) error {
	return r.db.Delete(&models.Product{}, "id = ?", id).Error
}

// FindAll finds all products with pagination and filtering
func (r *ProductRepository) FindAll(page, limit int, filters map[string]interface{}) ([]models.Product, int64, error) {
	var products []models.Product
	var count int64

	offset := (page - 1) * limit
	query := r.db.Model(&models.Product{})

	// Apply filters
	if categoryID, ok := filters["category_id"].(string); ok && categoryID != "" {
		query = query.Where("category_id = ?", categoryID)
	}

	if minPrice, ok := filters["min_price"].(float64); ok && minPrice > 0 {
		query = query.Where("price >= ?", minPrice)
	}

	if maxPrice, ok := filters["max_price"].(float64); ok && maxPrice > 0 {
		query = query.Where("price <= ?", maxPrice)
	}

	if size, ok := filters["size"].(string); ok && size != "" {
		query = query.Joins("JOIN product_variants ON products.id = product_variants.product_id").
			Where("product_variants.size = ?", size)
	}

	if color, ok := filters["color"].(string); ok && color != "" {
		query = query.Joins("JOIN product_variants ON products.id = product_variants.product_id").
			Where("product_variants.color = ?", color)
	}

	if rating, ok := filters["min_rating"].(int); ok && rating > 0 {
		query = query.Joins("LEFT JOIN reviews ON products.id = reviews.product_id").
			Group("products.id").
			Having("AVG(reviews.rating) >= ?", rating)
	}

	// Apply sorting
	if sortBy, ok := filters["sort_by"].(string); ok {
		switch sortBy {
		case "newest":
			query = query.Order("created_at DESC")
		case "price_asc":
			query = query.Order("price ASC")
		case "price_desc":
			query = query.Order("price DESC")
		case "rating":
			query = query.Joins("LEFT JOIN reviews ON products.id = reviews.product_id").
				Group("products.id").
				Order("AVG(reviews.rating) DESC")
		case "popularity":
			query = query.Joins("LEFT JOIN order_items ON products.id = order_items.product_id").
				Group("products.id").
				Order("COUNT(order_items.id) DESC")
		default:
			query = query.Order("created_at DESC")
		}
	} else {
		query = query.Order("created_at DESC")
	}

	// Get count
	err := query.Count(&count).Error
	if err != nil {
		return nil, 0, err
	}

	// Get products
	err = query.Preload("Category").Preload("Images").Preload("Variants").Preload("Tags").
		Offset(offset).Limit(limit).Find(&products).Error
	if err != nil {
		return nil, 0, err
	}

	return products, count, nil
}

// Search searches for products by keyword
func (r *ProductRepository) Search(keyword string, page, limit int) ([]models.Product, int64, error) {
	var products []models.Product
	var count int64

	offset := (page - 1) * limit
	query := r.db.Model(&models.Product{}).
		Where("name ILIKE ? OR description ILIKE ?", "%"+keyword+"%", "%"+keyword+"%")

	// Get count
	err := query.Count(&count).Error
	if err != nil {
		return nil, 0, err
	}

	// Get products
	err = query.Preload("Category").Preload("Images").Preload("Variants").Preload("Tags").
		Offset(offset).Limit(limit).Find(&products).Error
	if err != nil {
		return nil, 0, err
	}

	return products, count, nil
}

// UpdateStock updates a product's stock
func (r *ProductRepository) UpdateStock(productID string, quantity int) error {
	return r.db.Model(&models.Product{}).Where("id = ?", productID).
		Update("stock", gorm.Expr("stock + ?", quantity)).Error
}

// AddProductImage adds an image to a product
func (r *ProductRepository) AddProductImage(image *models.ProductImage) error {
	return r.db.Create(image).Error
}

// DeleteProductImage deletes a product image
func (r *ProductRepository) DeleteProductImage(id string) error {
	return r.db.Delete(&models.ProductImage{}, "id = ?", id).Error
}

// AddProductVariant adds a variant to a product
func (r *ProductRepository) AddProductVariant(variant *models.ProductVariant) error {
	return r.db.Create(variant).Error
}

// UpdateProductVariant updates a product variant
func (r *ProductRepository) UpdateProductVariant(variant *models.ProductVariant) error {
	return r.db.Save(variant).Error
}

// DeleteProductVariant deletes a product variant
func (r *ProductRepository) DeleteProductVariant(id string) error {
	return r.db.Delete(&models.ProductVariant{}, "id = ?", id).Error
}
