package repository

import (
    "database/sql"
    "errors"
    "clothing-shop-api/internal/domain/models"
)

type OrderRepository interface {
    CreateOrder(order *models.Order) error
    GetOrderByID(id int) (*models.Order, error)
    GetAllOrders() ([]models.Order, error)
}

type orderRepository struct {
    db *sql.DB
}

func NewOrderRepository(db *sql.DB) OrderRepository {
    return &orderRepository{db: db}
}

func (r *orderRepository) CreateOrder(order *models.Order) error {
    query := "INSERT INTO orders (user_id, product_id, quantity, status) VALUES (?, ?, ?, ?)"
    _, err := r.db.Exec(query, order.UserID, order.ProductID, order.Quantity, order.Status)
    return err
}

func (r *orderRepository) GetOrderByID(id int) (*models.Order, error) {
    query := "SELECT id, user_id, product_id, quantity, status FROM orders WHERE id = ?"
    row := r.db.QueryRow(query, id)

    var order models.Order
    err := row.Scan(&order.ID, &order.UserID, &order.ProductID, &order.Quantity, &order.Status)
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, errors.New("order not found")
        }
        return nil, err
    }
    return &order, nil
}

func (r *orderRepository) GetAllOrders() ([]models.Order, error) {
    query := "SELECT id, user_id, product_id, quantity, status FROM orders"
    rows, err := r.db.Query(query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var orders []models.Order
    for rows.Next() {
        var order models.Order
        if err := rows.Scan(&order.ID, &order.UserID, &order.ProductID, &order.Quantity, &order.Status); err != nil {
            return nil, err
        }
        orders = append(orders, order)
    }
    return orders, nil
}