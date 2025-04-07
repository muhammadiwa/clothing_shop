package repository

import (
	"clothing-shop-api/internal/domain/models"
	"database/sql"
	"fmt"
)

type CartRepository interface {
	AddItem(cart *models.Cart) error
	UpdateItem(cart *models.Cart) error
	RemoveItem(id uint) error
	GetUserCart(userID uint) ([]models.Cart, error)
	GetCartItem(id uint) (*models.Cart, error)
	ClearCart(userID uint) error
}

type cartRepository struct {
	db *sql.DB
}

func NewCartRepository(db *sql.DB) CartRepository {
	return &cartRepository{db: db}
}

func (r *cartRepository) AddItem(cart *models.Cart) error {
	// Check if the item already exists in the cart
	existingItem, err := r.getCartItemByProductVariant(cart.UserID, cart.ProductVariantID)
	if err != nil && err != sql.ErrNoRows {
		return err
	}

	// If item exists, update quantity
	if existingItem != nil {
		existingItem.Quantity += cart.Quantity
		return r.UpdateItem(existingItem)
	}

	// Otherwise add new item
	query := `
        INSERT INTO carts (
            user_id, product_variant_id, quantity, 
            created_at, updated_at
        ) VALUES (?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
    `
	result, err := r.db.Exec(query,
		cart.UserID,
		cart.ProductVariantID,
		cart.Quantity,
	)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	cart.ID = uint(id)
	return nil
}

func (r *cartRepository) UpdateItem(cart *models.Cart) error {
	query := `
        UPDATE carts 
        SET 
            quantity = ?,
            updated_at = CURRENT_TIMESTAMP
        WHERE id = ?
    `
	result, err := r.db.Exec(query, cart.Quantity, cart.ID)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (r *cartRepository) RemoveItem(id uint) error {
	query := `DELETE FROM carts WHERE id = ?`
	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (r *cartRepository) GetUserCart(userID uint) ([]models.Cart, error) {
	query := `
        SELECT 
            c.id, c.user_id, c.product_variant_id, c.quantity, 
            c.created_at, c.updated_at,
            pv.id, pv.product_id, pv.size, pv.color, pv.sku, 
            pv.price, pv.stock,
            p.id, p.name, p.description, p.category_id,
            p.base_price, p.discount, p.weight, p.rating,
            p.review_count
        FROM carts c
        JOIN product_variants pv ON c.product_variant_id = pv.id
        JOIN products p ON pv.product_id = p.id
        WHERE c.user_id = ? 
        AND p.deleted_at IS NULL 
        AND pv.deleted_at IS NULL
    `

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cartItems []models.Cart
	for rows.Next() {
		var cart models.Cart
		cart.ProductVariant = models.ProductVariant{
			Product: &models.Product{},
		}

		err := rows.Scan(
			&cart.ID, &cart.UserID, &cart.ProductVariantID, &cart.Quantity,
			&cart.CreatedAt, &cart.UpdatedAt,
			&cart.ProductVariant.ID, &cart.ProductVariant.ProductID,
			&cart.ProductVariant.Size, &cart.ProductVariant.Color,
			&cart.ProductVariant.SKU, &cart.ProductVariant.Price,
			&cart.ProductVariant.Stock,
			&cart.ProductVariant.Product.ID, &cart.ProductVariant.Product.Name,
			&cart.ProductVariant.Product.Description,
			&cart.ProductVariant.Product.CategoryID,
			&cart.ProductVariant.Product.BasePrice,
			&cart.ProductVariant.Product.Discount,
			&cart.ProductVariant.Product.Weight,
			&cart.ProductVariant.Product.Rating,
			&cart.ProductVariant.Product.ReviewCount,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning cart item: %v", err)
		}

		cartItems = append(cartItems, cart)
	}

	return cartItems, nil
}

func (r *cartRepository) GetCartItem(id uint) (*models.Cart, error) {
	query := `
        SELECT 
            c.id, c.user_id, c.product_variant_id, c.quantity, 
            c.created_at, c.updated_at,
            pv.id, pv.product_id, pv.size, pv.color, pv.sku, 
            pv.price, pv.stock,
            p.id, p.name, p.description, p.category_id,
            p.base_price, p.discount, p.weight, p.rating,
            p.review_count
        FROM carts c
        JOIN product_variants pv ON c.product_variant_id = pv.id
        JOIN products p ON pv.product_id = p.id
        WHERE c.id = ? 
        AND p.deleted_at IS NULL 
        AND pv.deleted_at IS NULL
    `

	cart := &models.Cart{
		ProductVariant: models.ProductVariant{
			Product: &models.Product{},
		},
	}

	err := r.db.QueryRow(query, id).Scan(
		&cart.ID, &cart.UserID, &cart.ProductVariantID, &cart.Quantity,
		&cart.CreatedAt, &cart.UpdatedAt,
		&cart.ProductVariant.ID, &cart.ProductVariant.ProductID,
		&cart.ProductVariant.Size, &cart.ProductVariant.Color,
		&cart.ProductVariant.SKU, &cart.ProductVariant.Price,
		&cart.ProductVariant.Stock,
		&cart.ProductVariant.Product.ID, &cart.ProductVariant.Product.Name,
		&cart.ProductVariant.Product.Description,
		&cart.ProductVariant.Product.CategoryID,
		&cart.ProductVariant.Product.BasePrice,
		&cart.ProductVariant.Product.Discount,
		&cart.ProductVariant.Product.Weight,
		&cart.ProductVariant.Product.Rating,
		&cart.ProductVariant.Product.ReviewCount,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("error scanning cart item: %v", err)
	}

	return cart, nil
}

func (r *cartRepository) ClearCart(userID uint) error {
	query := `DELETE FROM carts WHERE user_id = ?`
	_, err := r.db.Exec(query, userID)
	return err
}

// Helper methods
func (r *cartRepository) getCartItemByProductVariant(userID, productVariantID uint) (*models.Cart, error) {
	query := `
        SELECT 
            id, user_id, product_variant_id, quantity,
            created_at, updated_at
        FROM carts
        WHERE user_id = ? AND product_variant_id = ?
    `
	row := r.db.QueryRow(query, userID, productVariantID)

	var cart models.Cart
	err := row.Scan(
		&cart.ID, &cart.UserID, &cart.ProductVariantID, &cart.Quantity,
		&cart.CreatedAt, &cart.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &cart, nil
}
