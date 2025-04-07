package repository

import (
	"database/sql"
	"errors"

	"clothing-shop-api/internal/domain/models"
)

type ProductVariantRepository interface {
	Create(variant *models.ProductVariant) error
	FindByID(id uint) (*models.ProductVariant, error)
	FindByProductID(productID uint) ([]models.ProductVariant, error)
	Update(variant *models.ProductVariant) error
	Delete(id uint) error
	UpdateStock(id uint, quantity int) error
}

type productVariantRepositoryImpl struct {
	db *sql.DB
}

func NewProductVariantRepository(db *sql.DB) ProductVariantRepository {
	return &productVariantRepositoryImpl{db: db}
}

func (r *productVariantRepositoryImpl) Create(variant *models.ProductVariant) error {
	query := `INSERT INTO product_variants (product_id, size, color, sku, price, stock) 
              VALUES (?, ?, ?, ?, ?, ?)`

	result, err := r.db.Exec(query,
		variant.ProductID, variant.Size, variant.Color,
		variant.SKU, variant.Price, variant.Stock)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	variant.ID = uint(id)
	return nil
}

func (r *productVariantRepositoryImpl) FindByID(id uint) (*models.ProductVariant, error) {
	query := `SELECT id, product_id, size, color, sku, price, stock, created_at, updated_at
              FROM product_variants
              WHERE id = ? AND deleted_at IS NULL`

	row := r.db.QueryRow(query, id)

	variant := &models.ProductVariant{}
	err := row.Scan(
		&variant.ID, &variant.ProductID, &variant.Size, &variant.Color,
		&variant.SKU, &variant.Price, &variant.Stock, &variant.CreatedAt, &variant.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return variant, nil
}

func (r *productVariantRepositoryImpl) FindByProductID(productID uint) ([]models.ProductVariant, error) {
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

func (r *productVariantRepositoryImpl) Update(variant *models.ProductVariant) error {
	query := `UPDATE product_variants 
              SET size = ?, color = ?, sku = ?, price = ?, stock = ? 
              WHERE id = ? AND deleted_at IS NULL`

	_, err := r.db.Exec(query,
		variant.Size, variant.Color, variant.SKU,
		variant.Price, variant.Stock, variant.ID)

	return err
}

func (r *productVariantRepositoryImpl) Delete(id uint) error {
	query := `UPDATE product_variants SET deleted_at = CURRENT_TIMESTAMP WHERE id = ? AND deleted_at IS NULL`
	_, err := r.db.Exec(query, id)
	return err
}

func (r *productVariantRepositoryImpl) UpdateStock(id uint, quantity int) error {
	query := `UPDATE product_variants SET stock = stock - ? WHERE id = ? AND deleted_at IS NULL AND stock >= ?`
	result, err := r.db.Exec(query, quantity, id, quantity)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("insufficient stock or variant not found")
	}

	return nil
}
