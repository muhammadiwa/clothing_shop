package repository

import (
	"database/sql"

	"clothing-shop-api/internal/domain/models"
)

type wishlistRepositoryImpl struct {
	db *sql.DB
}

func NewWishlistRepository(db *sql.DB) WishlistRepository {
	return &wishlistRepositoryImpl{db: db}
}

func (r *wishlistRepositoryImpl) Add(wishlist *models.Wishlist) error {
	query := `
        INSERT INTO wishlists (user_id, product_id)
        VALUES (?, ?)
        ON DUPLICATE KEY UPDATE created_at = CURRENT_TIMESTAMP
    `
	result, err := r.db.Exec(query, wishlist.UserID, wishlist.ProductID)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	wishlist.ID = uint(id)
	return nil
}

func (r *wishlistRepositoryImpl) Remove(id uint) error {
	query := `DELETE FROM wishlists WHERE id = ?`
	_, err := r.db.Exec(query, id)
	return err
}

func (r *wishlistRepositoryImpl) GetUserWishlist(userID uint) ([]models.Wishlist, error) {
	query := `
        SELECT 
            w.id, w.user_id, w.product_id, w.created_at,
            p.id, p.name, p.description, p.category_id, p.base_price, p.discount, 
            p.weight, p.rating, p.review_count
        FROM wishlists w
        JOIN products p ON w.product_id = p.id
        WHERE w.user_id = ? AND p.deleted_at IS NULL
    `

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var wishlist []models.Wishlist
	for rows.Next() {
		var item models.Wishlist
		var product models.Product

		item.Product = product

		err := rows.Scan(
			&item.ID, &item.UserID, &item.ProductID, &item.CreatedAt,
			&item.Product.ID, &item.Product.Name, &item.Product.Description,
			&item.Product.CategoryID, &item.Product.BasePrice, &item.Product.Discount,
			&item.Product.Weight, &item.Product.Rating, &item.Product.ReviewCount,
		)
		if err != nil {
			return nil, err
		}

		wishlist = append(wishlist, item)
	}

	return wishlist, nil
}

func (r *wishlistRepositoryImpl) FindByUserAndProduct(userID, productID uint) (*models.Wishlist, error) {
	query := `
        SELECT id, user_id, product_id, created_at
        FROM wishlists
        WHERE user_id = ? AND product_id = ?
    `

	row := r.db.QueryRow(query, userID, productID)

	var wishlist models.Wishlist
	err := row.Scan(&wishlist.ID, &wishlist.UserID, &wishlist.ProductID, &wishlist.CreatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &wishlist, nil
}
