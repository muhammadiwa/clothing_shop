package repository

import (
	"database/sql"
	"fmt"

	"clothing-shop-api/internal/domain/models"
	"clothing-shop-api/internal/domain/services"
)

type productRepositoryImpl struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) ProductRepository {
	return &productRepositoryImpl{db: db}
}

func (r *productRepositoryImpl) Create(product *models.Product) error {
	// Begin transaction
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// Insert product
	query := `INSERT INTO products (name, description, category_id, base_price, discount, weight) 
              VALUES (?, ?, ?, ?, ?, ?)`
	result, err := tx.Exec(query, product.Name, product.Description, product.CategoryID,
		product.BasePrice, product.Discount, product.Weight)
	if err != nil {
		return err
	}

	// Get product ID
	productID, err := result.LastInsertId()
	if err != nil {
		return err
	}
	product.ID = uint(productID)

	// Commit transaction
	return tx.Commit()
}

func (r *productRepositoryImpl) FindByID(id uint) (*models.Product, error) {
	query := `SELECT p.id, p.name, p.description, p.category_id, p.base_price, p.discount, p.weight,
              p.rating, p.review_count, p.created_at, p.updated_at,
              c.id, c.name, c.slug, c.description, c.parent_id
              FROM products p
              JOIN categories c ON p.category_id = c.id
              WHERE p.id = ? AND p.deleted_at IS NULL`

	row := r.db.QueryRow(query, id)

	product := &models.Product{}
	var categoryParentID *uint

	err := row.Scan(
		&product.ID, &product.Name, &product.Description, &product.CategoryID,
		&product.BasePrice, &product.Discount, &product.Weight,
		&product.Rating, &product.ReviewCount, &product.CreatedAt, &product.UpdatedAt,
		&product.Category.ID, &product.Category.Name, &product.Category.Slug,
		&product.Category.Description, &categoryParentID,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, services.ErrProductNotFound
		}
		return nil, err
	}

	product.Category.ParentID = categoryParentID

	// Get variants
	variants, err := r.getProductVariants(id)
	if err != nil {
		return nil, err
	}
	product.Variants = variants

	// Get images
	images, err := r.getProductImages(id)
	if err != nil {
		return nil, err
	}
	product.Images = images

	return product, nil
}

func (r *productRepositoryImpl) FindAll(filter services.ProductFilter) ([]*models.Product, int, error) {
	// Build WHERE clause based on filters
	whereClause := "p.deleted_at IS NULL"
	args := []interface{}{}

	if filter.CategoryID != nil {
		whereClause += " AND p.category_id = ?"
		args = append(args, *filter.CategoryID)
	}

	if filter.MinPrice != nil {
		whereClause += " AND p.base_price * (1 - p.discount/100) >= ?"
		args = append(args, *filter.MinPrice)
	}

	if filter.MaxPrice != nil {
		whereClause += " AND p.base_price * (1 - p.discount/100) <= ?"
		args = append(args, *filter.MaxPrice)
	}

	if filter.Search != "" {
		whereClause += " AND (p.name LIKE ? OR p.description LIKE ?)"
		searchArg := "%" + filter.Search + "%"
		args = append(args, searchArg, searchArg)
	}

	// Count total matching products
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM products p WHERE %s", whereClause)
	var total int
	err := r.db.QueryRow(countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Build ORDER BY clause
	orderBy := "p.id DESC"
	switch filter.SortBy {
	case "name_asc":
		orderBy = "p.name ASC"
	case "name_desc":
		orderBy = "p.name DESC"
	case "price_asc":
		orderBy = "p.base_price * (1 - p.discount/100) ASC"
	case "price_desc":
		orderBy = "p.base_price * (1 - p.discount/100) DESC"
	case "newest":
		orderBy = "p.created_at DESC"
	}

	// Add pagination
	limit := filter.PageSize
	offset := (filter.Page - 1) * filter.PageSize

	// Select products with pagination
	query := fmt.Sprintf(`
        SELECT p.id, p.name, p.description, p.category_id, p.base_price, p.discount, p.weight,
        p.rating, p.review_count, p.created_at, p.updated_at,
        c.id, c.name, c.slug, c.description, c.parent_id,
        MIN(pv.price) as min_price, MAX(pv.price) as max_price
        FROM products p
        JOIN categories c ON p.category_id = c.id
        LEFT JOIN product_variants pv ON p.id = pv.product_id AND pv.deleted_at IS NULL
        WHERE %s
        GROUP BY p.id
        ORDER BY %s
        LIMIT ? OFFSET ?
    `, whereClause, orderBy)

	args = append(args, limit, offset)
	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	products := []*models.Product{}
	for rows.Next() {
		product := &models.Product{}
		category := &models.Category{}
		product.Category = *category

		var categoryParentID *uint
		var minPrice, maxPrice float64

		err := rows.Scan(
			&product.ID, &product.Name, &product.Description, &product.CategoryID,
			&product.BasePrice, &product.Discount, &product.Weight,
			&product.Rating, &product.ReviewCount, &product.CreatedAt, &product.UpdatedAt,
			&product.Category.ID, &product.Category.Name, &product.Category.Slug,
			&product.Category.Description, &categoryParentID,
			&minPrice, &maxPrice,
		)
		if err != nil {
			return nil, 0, err
		}

		product.Category.ParentID = categoryParentID
		product.MinPrice = minPrice
		product.MaxPrice = maxPrice

		products = append(products, product)
	}

	return products, total, nil
}

func (r *productRepositoryImpl) Update(product *models.Product) error {
	query := `UPDATE products 
              SET name = ?, description = ?, category_id = ?, base_price = ?, 
              discount = ?, weight = ? 
              WHERE id = ? AND deleted_at IS NULL`

	_, err := r.db.Exec(query,
		product.Name, product.Description, product.CategoryID,
		product.BasePrice, product.Discount, product.Weight, product.ID)

	return err
}

func (r *productRepositoryImpl) Delete(id uint) error {
	// Using soft delete
	query := `UPDATE products SET deleted_at = CURRENT_TIMESTAMP WHERE id = ? AND deleted_at IS NULL`
	_, err := r.db.Exec(query, id)
	return err
}

func (r *productRepositoryImpl) getProductVariants(productID uint) ([]models.ProductVariant, error) {
	query := `SELECT id, product_id, size, color, sku, price, stock, created_at, updated_at
              FROM product_variants
              WHERE product_id = ? AND deleted_at IS NULL`

	rows, err := r.db.Query(query, productID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	variants := []models.ProductVariant{}
	for rows.Next() {
		var variant models.ProductVariant
		err := rows.Scan(
			&variant.ID, &variant.ProductID, &variant.Size, &variant.Color,
			&variant.SKU, &variant.Price, &variant.Stock, &variant.CreatedAt, &variant.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		variants = append(variants, variant)
	}

	return variants, nil
}

func (r *productRepositoryImpl) getProductImages(productID uint) ([]models.ProductImage, error) {
	query := `SELECT id, product_id, url, is_primary, sort_order, created_at, updated_at
              FROM product_images
              WHERE product_id = ?
              ORDER BY is_primary DESC, sort_order ASC`

	rows, err := r.db.Query(query, productID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	images := []models.ProductImage{}
	for rows.Next() {
		var image models.ProductImage
		err := rows.Scan(
			&image.ID, &image.ProductID, &image.URL, &image.IsPrimary,
			&image.SortOrder, &image.CreatedAt, &image.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		images = append(images, image)
	}

	return images, nil
}
