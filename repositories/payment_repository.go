package repositories

import (
	"github.com/fashion-shop/database"
	"github.com/fashion-shop/models"
	"gorm.io/gorm"
)

// PaymentRepository handles database operations for payments
type PaymentRepository struct {
	db *gorm.DB
}

// NewPaymentRepository creates a new payment repository
func NewPaymentRepository() *PaymentRepository {
	return &PaymentRepository{
		db: database.DB,
	}
}

// Create creates a new payment
func (r *PaymentRepository) Create(payment *models.Payment) error {
	return r.db.Create(payment).Error
}

// FindByID finds a payment by ID
func (r *PaymentRepository) FindByID(id string) (*models.Payment, error) {
	var payment models.Payment
	err := r.db.Preload("Order").First(&payment, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &payment, nil
}

// FindByOrderID finds a payment by order ID
func (r *PaymentRepository) FindByOrderID(orderID string) (*models.Payment, error) {
	var payment models.Payment
	err := r.db.Preload("Order").First(&payment, "order_id = ?", orderID).Error
	if err != nil {
		return nil, err
	}
	return &payment, nil
}

// FindByTransactionID finds a payment by transaction ID
func (r *PaymentRepository) FindByTransactionID(transactionID string) (*models.Payment, error) {
	var payment models.Payment
	err := r.db.Preload("Order").First(&payment, "transaction_id = ?", transactionID).Error
	if err != nil {
		return nil, err
	}
	return &payment, nil
}

// Update updates a payment
func (r *PaymentRepository) Update(payment *models.Payment) error {
	return r.db.Save(payment).Error
}

// UpdateStatus updates a payment's status
func (r *PaymentRepository) UpdateStatus(id string, status string) error {
	return r.db.Model(&models.Payment{}).Where("id = ?", id).Update("status", status).Error
}
