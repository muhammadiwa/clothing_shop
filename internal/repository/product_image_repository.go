package repository

import (
	"database/sql"

	"clothing-shop-api/internal/domain/models"
)

type ProductImageRepository interface {
	Create(image *models.ProductImage) error
	FindByID(id uint) (*models.ProductImage, error)
	FindByProductID(productID uint) ([]models.ProductImage, error)
	FindPrimaryByProductID(productID uint) (*models.ProductImage, error)
	Delete(id uint) error
	SetPrimary(id uint) error
}

type productImageRepositoryImpl struct {
	db *sql.DB
}

func NewProductImageRepository(db *sql.DB) ProductImageRepository {
	return &productImageRepositoryImpl{db: db}
}

func (r *productImageRepositoryImpl) Create(image *models.ProductImage) error {
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

	// If this image is primary, unset other primary images
	if image.IsPrimary {
		_, err = tx.Exec(
			"UPDATE product_images SET is_primary = FALSE WHERE product_id = ?",
			image.ProductID,
		)
		if err != nil {
			return err
		}
	}

	// Insert the image
	query := `INSERT INTO product_images (product_id, url, is_primary, sort_order)
              VALUES (?, ?, ?, ?)`

	result, err := tx.Exec(query,
		image.ProductID, image.URL, image.IsPrimary, image.SortOrder)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	image.ID = uint(id)

	return tx.Commit()
}

func (r *productImageRepositoryImpl) FindByID(id uint) (*models.ProductImage, error) {
	query := `SELECT id, product_id, url, is_primary, sort_order, created_at, updated_at
              FROM product_images
              WHERE id = ?`

	row := r.db.QueryRow(query, id)

	image := &models.ProductImage{}
	err := row.Scan(
		&image.ID, &image.ProductID, &image.URL, &image.IsPrimary,
		&image.SortOrder, &image.CreatedAt, &image.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return image, nil
}

func (r *productImageRepositoryImpl) FindByProductID(productID uint) ([]models.ProductImage, error) {
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

func (r *productImageRepositoryImpl) FindPrimaryByProductID(productID uint) (*models.ProductImage, error) {
	query := `SELECT id, product_id, url, is_primary, sort_order, created_at, updated_at
              FROM product_images
              WHERE product_id = ? AND is_primary = TRUE
              LIMIT 1`

	row := r.db.QueryRow(query, productID)

	image := &models.ProductImage{}
	err := row.Scan(
		&image.ID, &image.ProductID, &image.URL, &image.IsPrimary,
		&image.SortOrder, &image.CreatedAt, &image.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			// If no primary image, try to get first image
			images, err := r.FindByProductID(productID)
			if err != nil || len(images) == 0 {
				return nil, err
			}
			return &images[0], nil
		}
		return nil, err
	}

	return image, nil
}

func (r *productImageRepositoryImpl) Delete(id uint) error {
	query := `DELETE FROM product_images WHERE id = ?`
	_, err := r.db.Exec(query, id)
	return err
}

func (r *productImageRepositoryImpl) SetPrimary(id uint) error {
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

	// Get product ID of the image
	var productID uint
	err = tx.QueryRow("SELECT product_id FROM product_images WHERE id = ?", id).Scan(&productID)
	if err != nil {
		return err
	}

	// Unset all primary images for this product
	_, err = tx.Exec(
		"UPDATE product_images SET is_primary = FALSE WHERE product_id = ?",
		productID,
	)
	if err != nil {
		return err
	}

	// Set the specified image as primary
	_, err = tx.Exec("UPDATE product_images SET is_primary = TRUE WHERE id = ?", id)
	if err != nil {
		return err
	}

	return tx.Commit()
}
