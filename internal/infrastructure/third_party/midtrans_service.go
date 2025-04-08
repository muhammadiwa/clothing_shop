package third_party

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/coreapi"
	"github.com/midtrans/midtrans-go/snap"
)

// MidtransService defines the interface for Midtrans operations
type MidtransService interface {
	CreateTransaction(orderID string, amount float64, customerDetails map[string]interface{}, itemDetails []map[string]interface{}, paymentType string) (string, error)
	GetTransactionStatus(transactionID string) (map[string]interface{}, error)
	RefundTransaction(transactionID string, amount float64, reason string) error
}

type midtransService struct {
	serverKey   string
	clientKey   string
	environment string
	snapClient  snap.Client
	coreClient  coreapi.Client
}

// NewMidtransService creates a new MidtransService instance
func NewMidtransService(serverKey, clientKey, environment string) MidtransService {
	var env midtrans.EnvironmentType
	if environment == "production" {
		env = midtrans.Production
	} else {
		env = midtrans.Sandbox
	}

	snapClient := snap.Client{}
	snapClient.New(serverKey, env)

	coreClient := coreapi.Client{}
	coreClient.New(serverKey, env)

	return &midtransService{
		serverKey:   serverKey,
		clientKey:   clientKey,
		environment: environment,
		snapClient:  snapClient,
		coreClient:  coreClient,
	}
}

// CreateTransaction creates a new transaction
func (s *midtransService) CreateTransaction(orderID string, amount float64, customerDetails map[string]interface{}, itemDetails []map[string]interface{}, paymentType string) (string, error) {
	req := &snap.Request{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  orderID,
			GrossAmt: int64(amount),
		},
		CustomerDetail: &midtrans.CustomerDetails{
			FName: customerDetails["first_name"].(string),
			LName: customerDetails["last_name"].(string),
			Email: customerDetails["email"].(string),
			Phone: customerDetails["phone"].(string),
		},
	}

	// Add item details
	var items []midtrans.ItemDetails
	for _, item := range itemDetails {
		price, _ := strconv.ParseInt(fmt.Sprintf("%.0f", item["price"].(float64)), 10, 64)
		qty, _ := strconv.ParseInt(fmt.Sprintf("%.0f", item["quantity"].(float64)), 10, 64)

		items = append(items, midtrans.ItemDetails{
			ID:    item["id"].(string),
			Name:  item["name"].(string),
			Price: price,
			Qty:   int32(qty),
		})
	}
	req.Items = &items

	// Set payment type
	if paymentType != "" {
		req.EnabledPayments = []snap.SnapPaymentType{snap.SnapPaymentType(paymentType)}
	}

	resp, err := s.snapClient.CreateTransaction(req)
	if err != nil {
		return "", err
	}

	return resp.RedirectURL, nil
}

// GetTransactionStatus gets the status of a transaction
func (s *midtransService) GetTransactionStatus(transactionID string) (map[string]interface{}, error) {
	resp, err := s.coreClient.CheckTransaction(transactionID)
	if err != nil {
		return nil, err
	}

	result := map[string]interface{}{
		"transaction_id":     resp.TransactionID,
		"order_id":           resp.OrderID,
		"status_code":        resp.StatusCode,
		"status_message":     resp.StatusMessage,
		"transaction_status": resp.TransactionStatus,
		"fraud_status":       resp.FraudStatus,
		"payment_type":       resp.PaymentType,
		"gross_amount":       resp.GrossAmount,
	}

	return result, nil
}

// RefundTransaction refunds a transaction
func (s *midtransService) RefundTransaction(transactionID string, amount float64, reason string) error {
	resp, err := s.coreClient.RefundTransaction(transactionID, &coreapi.RefundReq{
		RefundKey: transactionID,
		Amount:    int64(amount),
		Reason:    reason,
	})

	if err != nil {
		return err
	}

	if resp.StatusCode != "200" {
		return errors.New(resp.StatusMessage)
	}

	return nil
}
