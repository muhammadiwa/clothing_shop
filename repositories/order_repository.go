package repositories

import (
	"github.com/fashion-shop/database"
	"github.com/fashion-shop/models"
	"gorm.io/gorm"
)

// OrderRepository handles database operations for orders
type OrderRepository struct {
	db *gorm.DB
}

// NewOrderRepository creates a new order repository
func NewOrderRepository() *OrderRepository {
	return &OrderRepository{
		db: database.DB,
	}
}

// Create creates a new order
func (r *OrderRepository) Create(order *models.Order) error {
	return r.db.Create(order).Error
}

// FindByID finds an order by ID
func (r *OrderRepository) FindByID(id string) (*models.Order, error) {
	var order models.Order
	err := r.db.Preload("User").Preload("Items.Product").Preload("Items.Variant").
		Preload("ShippingAddress").Preload("Payment").
		First(&order, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

// FindByOrderNumber finds an order by order number
func (r *OrderRepository) FindByOrderNumber(orderNumber string) (*models.Order, error) {
	var order models.Order
	err := r.db.Preload("User").Preload("Items.Product").Preload("Items.Variant").
		Preload("ShippingAddress").Preload("Payment").
		First(&order, "order_number = ?", orderNumber).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

// Update updates an order
func (r *OrderRepository) Update(order *models.Order) error {
	return r.db.Save(order).Error
}

// UpdateStatus updates an order's status
func (r *OrderRepository) UpdateStatus(id string, status string) error {
	return r.db.Model(&models.Order{}).Where("id = ?", id).Update("status", status).Error
}

// FindByUserID finds orders by user ID with pagination
func (r *OrderRepository) FindByUserID(userID string, page, limit int) ([]models.Order, int64, error) {
	var orders []models.Order
	var count int64

	offset := (page - 1) * limit

	// Get count
	err := r.db.Model(&models.Order{}).Where("user_id = ?", userID).Count(&count).Error
	if err != nil {
		return nil, 0, err
	}

	// Get orders
	err = r.db.Where("user_id = ?", userID).
		Preload("Items.Product").Preload("Items.Variant").Preload("Payment").
		Order("created_at DESC").
		Offset(offset).Limit(limit).Find(&orders).Error
	if err != nil {
		return nil, 0, err
	}

	return orders, count, nil
}

// FindAll finds all orders with pagination and filtering
func (r *OrderRepository) FindAll(page, limit int, filters map[string]interface{}) ([]models.Order, int64, error) {
	var orders []models.Order
	var count int64

	offset := (page - 1) * limit
	query := r.db.Model(&models.Order{})

	// Apply filters
	if status, ok := filters["status"].(string); ok && status != "" {
		query = query.Where("status = ?", status)
	}

	if userID, ok := filters["user_id"].(string); ok && userID != "" {
		query = query.Where("user_id = ?", userID)
	}

	// Get count
	err := query.Count(&count).Error
	if err != nil {
		return nil, 0, err
	}

	// Get orders
	err = query.Preload("User").Preload("Items.Product").Preload("Items.Variant").
		Preload("ShippingAddress").Preload("Payment").
		Order("created_at DESC").
		Offset(offset).Limit(limit).Find(&orders).Error
	if err != nil {
		return nil, 0, err
	}

	return orders, count, nil
}

// AddOrderItem adds an item to an order
func (r *OrderRepository) AddOrderItem(item *models.OrderItem) error {
	return r.db.Create(item).Error
}
