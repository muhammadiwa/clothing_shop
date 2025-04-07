package services

import (
	"errors"
	"time"

	"clothing-shop-api/internal/domain/models"
	"clothing-shop-api/internal/services/payment"
)

type PaymentService struct {
	paymentRepo    PaymentRepository
	orderRepo      OrderRepository
	midtransClient *payment.MidtransClient
}

func NewPaymentService(
	paymentRepo PaymentRepository,
	orderRepo OrderRepository,
	midtransClient *payment.MidtransClient,
) *PaymentService {
	return &PaymentService{
		paymentRepo:    paymentRepo,
		orderRepo:      orderRepo,
		midtransClient: midtransClient,
	}
}

func (s *PaymentService) CreatePayment(orderID uint, method string) (*models.Payment, error) {
	// Validate order exists
	order, err := s.orderRepo.FindByID(orderID)
	if err != nil {
		return nil, err
	}
	if order == nil {
		return nil, errors.New("order not found")
	}

	// Check if payment already exists for this order
	existingPayment, err := s.paymentRepo.FindByOrderID(orderID)
	if err != nil {
		return nil, err
	}
	if existingPayment != nil {
		return nil, errors.New("payment already exists for this order")
	}

	// Validate payment method
	if !isValidPaymentMethod(method) {
		return nil, errors.New("invalid payment method")
	}

	// Create payment record
	payment := &models.Payment{
		OrderID:       orderID,
		PaymentMethod: method,
		Amount:        order.TotalAmount,
		Status:        models.PaymentStatusPending,
		ExpiryTime:    time.Now().Add(24 * time.Hour),
	}

	// Create payment in Midtrans
	midtransResp, err := s.midtransClient.CreateTransaction(order, order.User, method)
	if err != nil {
		return nil, err
	}

	// Update payment with Midtrans response
	payment.TransactionID = midtransResp.TransactionID
	payment.PaymentToken = midtransResp.Token
	payment.PaymentURL = midtransResp.RedirectURL

	if len(midtransResp.VaNumbers) > 0 {
		payment.PaymentChannel = midtransResp.VaNumbers[0].Bank
		payment.VaNumber = midtransResp.VaNumbers[0].VANumber
	}

	// Save payment to database
	if err := s.paymentRepo.Create(payment); err != nil {
		return nil, err
	}

	return payment, nil
}

func (s *PaymentService) GetPaymentStatus(id uint) (*models.Payment, error) {
	payment, err := s.paymentRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if payment == nil {
		return nil, errors.New("payment not found")
	}
	return payment, nil
}

func (s *PaymentService) UpdatePaymentStatus(id uint, status models.PaymentStatus) error {
	payment, err := s.paymentRepo.FindByID(id)
	if err != nil {
		return err
	}
	if payment == nil {
		return errors.New("payment not found")
	}

	payment.Status = status
	if status == models.PaymentStatusPaid {
		now := time.Now()
		payment.PaidAt = &now
	}

	return s.paymentRepo.Update(payment)
}

func (s *PaymentService) HandlePaymentNotification(notification *payment.PaymentNotification) error {
	payment, err := s.paymentRepo.FindByTransactionID(notification.TransactionID)
	if err != nil {
		return err
	}
	if payment == nil {
		return errors.New("payment not found")
	}

	var status models.PaymentStatus
	switch notification.TransactionStatus {
	case "capture", "settlement":
		status = models.PaymentStatusPaid
	case "deny", "cancel", "failure":
		status = models.PaymentStatusFailed
	case "expire":
		status = models.PaymentStatusExpired
	case "pending":
		status = models.PaymentStatusPending
	case "refund":
		status = models.PaymentStatusRefunded
	default:
		return errors.New("unknown transaction status")
	}

	return s.UpdatePaymentStatus(payment.ID, status)
}

func (s *PaymentService) RefundPayment(id uint, amount float64, reason string) error {
	payment, err := s.paymentRepo.FindByID(id)
	if err != nil {
		return err
	}
	if payment == nil {
		return errors.New("payment not found")
	}

	if payment.Status != models.PaymentStatusPaid {
		return errors.New("payment must be paid to be refunded")
	}

	if amount > payment.Amount {
		return errors.New("refund amount cannot exceed payment amount")
	}

	// Process refund through Midtrans
	if err := s.midtransClient.RefundTransaction(payment.TransactionID, amount, reason); err != nil {
		return err
	}

	// Update payment status
	return s.UpdatePaymentStatus(payment.ID, models.PaymentStatusRefunded)
}

func (s *PaymentService) GetPaymentsByOrderID(orderID uint) ([]*models.Payment, error) {
	return s.paymentRepo.FindByOrderID(orderID)
}

func (s *PaymentService) CancelPayment(id uint) error {
	payment, err := s.paymentRepo.FindByID(id)
	if err != nil {
		return err
	}
	if payment == nil {
		return errors.New("payment not found")
	}

	if payment.Status != models.PaymentStatusPending {
		return errors.New("only pending payments can be cancelled")
	}

	// Cancel payment in Midtrans
	if err := s.midtransClient.CancelTransaction(payment.TransactionID); err != nil {
		return err
	}

	return s.UpdatePaymentStatus(payment.ID, models.PaymentStatusFailed)
}

// Helper functions
func isValidPaymentMethod(method string) bool {
	validMethods := map[string]bool{
		"credit_card":   true,
		"bank_transfer": true,
		"gopay":         true,
		"shopeepay":     true,
		"qris":          true,
	}
	return validMethods[method]
}
