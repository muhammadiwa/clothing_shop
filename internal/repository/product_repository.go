package repository

import (
    "database/sql"
    "errors"
    "clothing-shop-api/internal/domain/models"
)

type ProductRepository interface {
    Create(product *models.Product) error
    GetByID(id int) (*models.Product, error)
    GetAll() ([]models.Product, error)
    Update(product *models.Product) error
    Delete(id int) error
}

type productRepository struct {
    db *sql.DB
}

func NewProductRepository(db *sql.DB) ProductRepository {
    return &productRepository{db: db}
}

func (r *productRepository) Create(product *models.Product) error {
    query := "INSERT INTO products (name, description, price, stock) VALUES (?, ?, ?, ?)"
    _, err := r.db.Exec(query, product.Name, product.Description, product.Price, product.Stock)
    return err
}

func (r *productRepository) GetByID(id int) (*models.Product, error) {
    query := "SELECT id, name, description, price, stock FROM products WHERE id = ?"
    row := r.db.QueryRow(query, id)

    var product models.Product
    err := row.Scan(&product.ID, &product.Name, &product.Description, &product.Price, &product.Stock)
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, errors.New("product not found")
        }
        return nil, err
    }
    return &product, nil
}

func (r *productRepository) GetAll() ([]models.Product, error) {
    query := "SELECT id, name, description, price, stock FROM products"
    rows, err := r.db.Query(query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var products []models.Product
    for rows.Next() {
        var product models.Product
        if err := rows.Scan(&product.ID, &product.Name, &product.Description, &product.Price, &product.Stock); err != nil {
            return nil, err
        }
        products = append(products, product)
    }
    return products, nil
}

func (r *productRepository) Update(product *models.Product) error {
    query := "UPDATE products SET name = ?, description = ?, price = ?, stock = ? WHERE id = ?"
    _, err := r.db.Exec(query, product.Name, product.Description, product.Price, product.Stock, product.ID)
    return err
}

func (r *productRepository) Delete(id int) error {
    query := "DELETE FROM products WHERE id = ?"
    _, err := r.db.Exec(query, id)
    return err
}