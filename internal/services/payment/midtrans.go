package payment

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"clothing-shop-api/internal/config"
	"clothing-shop-api/internal/domain/models"
)

type MidtransClient struct {
	BaseURL     string
	ServerKey   string
	ClientKey   string
	Environment string
	HTTPClient  *http.Client
}

type PaymentRequest struct {
	TransactionDetails struct {
		OrderID     string  `json:"order_id"`
		GrossAmount float64 `json:"gross_amount"`
	} `json:"transaction_details"`
	CustomerDetails struct {
		FirstName string `json:"first_name"`
		Email     string `json:"email"`
		Phone     string `json:"phone"`
	} `json:"customer_details"`
	ItemDetails []struct {
		ID       string  `json:"id"`
		Price    float64 `json:"price"`
		Quantity int     `json:"quantity"`
		Name     string  `json:"name"`
	} `json:"item_details"`
	PaymentType  string `json:"payment_type"`
	BankTransfer struct {
		Bank string `json:"bank"`
	} `json:"bank_transfer,omitempty"`
	CreditCard struct {
		Secure bool `json:"secure"`
	} `json:"credit_card,omitempty"`
	Expiry struct {
		Duration int    `json:"duration"`
		Unit     string `json:"unit"`
	} `json:"expiry,omitempty"`
}

type PaymentResponse struct {
	StatusCode        string    `json:"status_code"`
	StatusMessage     string    `json:"status_message"`
	TransactionID     string    `json:"transaction_id"`
	OrderID           string    `json:"order_id"`
	GrossAmount       string    `json:"gross_amount"`
	PaymentType       string    `json:"payment_type"`
	TransactionTime   time.Time `json:"transaction_time"`
	TransactionStatus string    `json:"transaction_status"`
	VaNumbers         []struct {
		Bank     string `json:"bank"`
		VANumber string `json:"va_number"`
	} `json:"va_numbers,omitempty"`
	PermataVaNumber string `json:"permata_va_number,omitempty"`
	RedirectURL     string `json:"redirect_url,omitempty"`
}

func NewMidtransClient(cfg *config.Config) *MidtransClient {
	baseURL := "https://api.sandbox.midtrans.com"
	if cfg.MidtransEnvironment == "production" {
		baseURL = "https://api.midtrans.com"
	}

	return &MidtransClient{
		BaseURL:     baseURL,
		ServerKey:   cfg.MidtransServerKey,
		ClientKey:   cfg.MidtransClientKey,
		Environment: cfg.MidtransEnvironment,
		HTTPClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *MidtransClient) CreateTransaction(order *models.Order, user *models.User, paymentMethod string) (*PaymentResponse, error) {
	req := PaymentRequest{}
	req.TransactionDetails.OrderID = order.OrderNumber
	req.TransactionDetails.GrossAmount = order.GrandTotal

	req.CustomerDetails.FirstName = user.Username
	req.CustomerDetails.Email = user.Email
	req.CustomerDetails.Phone = user.PhoneNumber

	// Add item details
	for _, item := range order.Items {
		req.ItemDetails = append(req.ItemDetails, struct {
			ID       string  `json:"id"`
			Price    float64 `json:"price"`
			Quantity int     `json:"quantity"`
			Name     string  `json:"name"`
		}{
			ID:       fmt.Sprintf("%d", item.ProductVariant.ID),
			Price:    item.UnitPrice,
			Quantity: item.Quantity,
			Name:     fmt.Sprintf("%s - %s/%s", item.ProductVariant.Product.Name, item.ProductVariant.Size, item.ProductVariant.Color),
		})
	}

	// Add shipping fee as an item
	req.ItemDetails = append(req.ItemDetails, struct {
		ID       string  `json:"id"`
		Price    float64 `json:"price"`
		Quantity int     `json:"quantity"`
		Name     string  `json:"name"`
	}{
		ID:       "SHIPPING",
		Price:    order.ShippingFee,
		Quantity: 1,
		Name:     fmt.Sprintf("Shipping Fee (%s)", order.Courier.Name),
	})

	// Set payment method specific details
	switch paymentMethod {
	case "bank_transfer":
		req.PaymentType = "bank_transfer"
		req.BankTransfer.Bank = "bca" // Can be parameterized
		req.Expiry.Duration = 24
		req.Expiry.Unit = "hour"
	case "credit_card":
		req.PaymentType = "credit_card"
		req.CreditCard.Secure = true
	// Add more payment methods as needed
	default:
		return nil, errors.New("unsupported payment method")
	}

	// Encode the request to JSON
	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	// Create the HTTP request
	httpReq, err := http.NewRequest("POST", c.BaseURL+"/v2/charge", strings.NewReader(string(jsonData)))
	if err != nil {
		return nil, err
	}

	// Set headers
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "application/json")
	httpReq.SetBasicAuth(c.ServerKey, "")

	// Send the request
	resp, err := c.HTTPClient.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Check response status code
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("midtrans API returned non-200/201 status: %d", resp.StatusCode)
	}

	// Parse the response
	var paymentResp PaymentResponse
	if err = json.NewDecoder(resp.Body).Decode(&paymentResp); err != nil {
		return nil, err
	}

	return &paymentResp, nil
}
