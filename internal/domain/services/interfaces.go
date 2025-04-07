package services

import "clothing-shop-api/internal/domain/models"

type UserRepository interface {
	Create(user *models.User) error
	FindByID(id uint) (*models.User, error)
	FindByEmail(email string) (*models.User, error)
	Update(user *models.User) error
	Delete(id uint) error
	FindByVerificationToken(token string) (*models.User, error)
	FindByResetToken(token string) (*models.User, error)
}

type OrderRepository interface {
	Create(order *models.Order) error
	FindByID(id uint) (*models.Order, error)
	FindByUserID(userID uint) ([]models.Order, error)
	FindAll() ([]models.Order, error)
	Update(order *models.Order) error
	Delete(id uint) error
}

type PaymentRepository interface {
	Create(payment *models.Payment) error
	FindByID(id uint) (*models.Payment, error)
	FindByOrderID(orderID uint) (*models.Payment, error)
	FindByTransactionID(transactionID string) (*models.Payment, error)
	Update(payment *models.Payment) error
	Delete(id uint) error
}

type CategoryRepository interface {
	Create(category *models.Category) error
	FindByID(id uint) (*models.Category, error)
	FindAll() ([]*models.Category, error)
	Update(category *models.Category) error
	Delete(id uint) error
}
