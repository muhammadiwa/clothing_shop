package repository

import (
	"database/sql"

	"clothing-shop-api/internal/domain/models"
)

type cartRepositoryImpl struct {
	db *sql.DB
}

func NewCartRepository(db *sql.DB) CartRepository {
	return &cartRepositoryImpl{db: db}
}

func (r *cartRepositoryImpl) AddItem(cart *models.Cart) error {
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
        INSERT INTO carts (user_id, product_variant_id, quantity)
        VALUES (?, ?, ?)
    `
	result, err := r.db.Exec(query, cart.UserID, cart.ProductVariantID, cart.Quantity)
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

func (r *cartRepositoryImpl) UpdateItem(cart *models.Cart) error {
	query := `
        UPDATE carts
        SET quantity = ?
        WHERE id = ?
    `
	_, err := r.db.Exec(query, cart.Quantity, cart.ID)
	return err
}

func (r *cartRepositoryImpl) RemoveItem(id uint) error {
	query := `DELETE FROM carts WHERE id = ?`
	_, err := r.db.Exec(query, id)
	return err
}

func (r *cartRepositoryImpl) GetUserCart(userID uint) ([]models.Cart, error) {
	query := `
        SELECT 
            c.id, c.user_id, c.product_variant_id, c.quantity, c.created_at, c.updated_at,
            pv.id, pv.product_id, pv.size, pv.color, pv.sku, pv.price, pv.stock,
            p.id, p.name, p.description, p.base_price, p.discount
        FROM carts c
        JOIN product_variants pv ON c.product_variant_id = pv.id
        JOIN products p ON pv.product_id = p.id
        WHERE c.user_id = ? AND p.deleted_at IS NULL AND pv.deleted_at IS NULL
    `

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cartItems []models.Cart
	for rows.Next() {
		var cart models.Cart
		var productVariant models.ProductVariant
		var product models.Product

		productVariant.Product = &product
		cart.ProductVariant = productVariant

		err := rows.Scan(
			&cart.ID, &cart.UserID, &cart.ProductVariantID, &cart.Quantity,
			&cart.CreatedAt, &cart.UpdatedAt,
			&cart.ProductVariant.ID, &cart.ProductVariant.ProductID,
			&cart.ProductVariant.Size, &cart.ProductVariant.Color,
			&cart.ProductVariant.SKU, &cart.ProductVariant.Price, &cart.ProductVariant.Stock,
			&cart.ProductVariant.Product.ID, &cart.ProductVariant.Product.Name,
			&cart.ProductVariant.Product.Description, &cart.ProductVariant.Product.BasePrice,
			&cart.ProductVariant.Product.Discount,
		)
		if err != nil {
			return nil, err
		}

		cartItems = append(cartItems, cart)
	}

	return cartItems, nil
}

func (r *cartRepositoryImpl) GetCartItem(id uint) (*models.Cart, error) {
	query := `
        SELECT 
            c.id, c.user_id, c.product_variant_id, c.quantity, c.created_at, c.updated_at,
            pv.id, pv.product_id, pv.size, pv.color, pv.sku, pv.price, pv.stock,
            p.id, p.name, p.description, p.base_price, p.discount
        FROM carts c
        JOIN product_variants pv ON c.product_variant_id = pv.id
        JOIN products p ON pv.product_id = p.id
        WHERE c.id = ? AND p.deleted_at IS NULL AND pv.deleted_at IS NULL
    `

	row := r.db.QueryRow(query, id)

	var cart models.Cart
	var productVariant models.ProductVariant
	var product models.Product

	productVariant.Product = &product
	cart.ProductVariant = productVariant

	err := row.Scan(
		&cart.ID, &cart.UserID, &cart.ProductVariantID, &cart.Quantity,
		&cart.CreatedAt, &cart.UpdatedAt,
		&cart.ProductVariant.ID, &cart.ProductVariant.ProductID,
		&cart.ProductVariant.Size, &cart.ProductVariant.Color,
		&cart.ProductVariant.SKU, &cart.ProductVariant.Price, &cart.ProductVariant.Stock,
		&cart.ProductVariant.Product.ID, &cart.ProductVariant.Product.Name,
		&cart.ProductVariant.Product.Description, &cart.ProductVariant.Product.BasePrice,
		&cart.ProductVariant.Product.Discount,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &cart, nil
}

func (r *cartRepositoryImpl) ClearCart(userID uint) error {
	query := `DELETE FROM carts WHERE user_id = ?`
	_, err := r.db.Exec(query, userID)
	return err
}

func (r *cartRepositoryImpl) getCartItemByProductVariant(userID, productVariantID uint) (*models.Cart, error) {
	query := `
        SELECT id, user_id, product_variant_id, quantity, created_at, updated_at
        FROM carts
        WHERE user_id = ? AND product_variant_id = ?
    `
	row := r.db.QueryRow(query, userID, productVariantID)

	var cart models.Cart
	err := row.Scan(
		&cart.ID, &cart.UserID, &cart.ProductVariantID, &cart.Quantity,
		&cart.CreatedAt, &cart.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &cart, nil
}
