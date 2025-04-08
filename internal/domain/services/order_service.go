// package services

// import (
// 	"errors"
// 	"time"

// 	"clothing-shop-api/internal/domain/models"
// 	"clothing-shop-api/internal/repository"
// )

// type OrderService struct {
// 	orderRepo repository.OrderRepository
// }

// func NewOrderService(orderRepo repository.OrderRepository) *OrderService {
// 	return &OrderService{orderRepo: orderRepo}
// }

// func (s *OrderService) CreateOrder(order *models.Order) error {
// 	if order == nil {
// 		return errors.New("order cannot be nil")
// 	}
// 	order.CreatedAt = time.Now()
// 	return s.orderRepo.Create(order)
// }

// func (s *OrderService) GetOrderByID(id string) (*models.Order, error) {
// 	return s.orderRepo.FindByID(id)
// }

// func (s *OrderService) GetAllOrders() ([]models.Order, error) {
// 	return s.orderRepo.FindAll()
// }

// func (s *OrderService) UpdateOrder(order *models.Order) error {
// 	if order == nil {
// 		return errors.New("order cannot be nil")
// 	}
// 	return s.orderRepo.Update(order)
// }

// func (s *OrderService) DeleteOrder(id string) error {
// 	return s.orderRepo.Delete(id)
// }