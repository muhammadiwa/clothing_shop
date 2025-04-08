package integrations

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/fashion-shop/config"
)

// MidtransService handles Midtrans payment gateway integration
type MidtransService struct {
	config *config.Config
}

// NewMidtransService creates a new Midtrans service
func NewMidtransService(config *config.Config) *MidtransService {
	return &MidtransService{
		config: config,
	}
}

// ChargeRequest represents a charge request to Midtrans
type ChargeRequest struct {
	PaymentType        string                 `json:"payment_type"`
	TransactionDetails TransactionDetails     `json:"transaction_details"`
	ItemDetails        []ItemDetail           `json:"item_details"`
	CustomerDetails    CustomerDetails        `json:"customer_details"`
	BankTransfer       *BankTransfer          `json:"bank_transfer,omitempty"`
	CreditCard         *CreditCard            `json:"credit_card,omitempty"`
	Gopay              *Gopay                 `json:"gopay,omitempty"`
	Callbacks          *Callbacks             `json:"callbacks,omitempty"`
	CustomFields       map[string]interface{} `json:"custom_fields,omitempty"`
}

// TransactionDetails represents transaction details
type TransactionDetails struct {
	OrderID     string `json:"order_id"`
	GrossAmount int64  `json:"gross_amount"`
}

// ItemDetail represents an item detail
type ItemDetail struct {
	ID       string `json:"id"`
	Price    int64  `json:"price"`
	Quantity int    `json:"quantity"`
	Name     string `json:"name"`
}

// CustomerDetails represents customer details
type CustomerDetails struct {
	FirstName       string   `json:"first_name"`
	LastName        string   `json:"last_name"`
	Email           string   `json:"email"`
	Phone           string   `json:"phone"`
	BillingAddress  *Address `json:"billing_address,omitempty"`
	ShippingAddress *Address `json:"shipping_address,omitempty"`
}

// Address represents an address
type Address struct {
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Phone       string `json:"phone"`
	Address     string `json:"address"`
	City        string `json:"city"`
	PostalCode  string `json:"postal_code"`
	CountryCode string `json:"country_code"`
}

// BankTransfer represents bank transfer details
type BankTransfer struct {
	Bank string `json:"bank"`
}

// CreditCard represents credit card details
type CreditCard struct {
	Secure        bool         `json:"secure"`
	Channel       string       `json:"channel,omitempty"`
	Bank          string       `json:"bank,omitempty"`
	Installment   *Installment `json:"installment,omitempty"`
	WhitelistBins []string     `json:"whitelist_bins,omitempty"`
	SaveCard      bool         `json:"save_card,omitempty"`
}

// Installment represents installment details
type Installment struct {
	Required bool             `json:"required"`
	Terms    map[string][]int `json:"terms,omitempty"`
}

// Gopay represents Gopay details
type Gopay struct {
	EnableCallback bool   `json:"enable_callback"`
	CallbackUrl    string `json:"callback_url"`
}

// Callbacks represents callback URLs
type Callbacks struct {
	Finish string `json:"finish"`
}

// ChargeResponse represents a charge response from Midtrans
type ChargeResponse struct {
	StatusCode        string              `json:"status_code"`
	StatusMessage     string              `json:"status_message"`
	TransactionID     string              `json:"transaction_id"`
	OrderID           string              `json:"order_id"`
	GrossAmount       string              `json:"gross_amount"`
	PaymentType       string              `json:"payment_type"`
	TransactionTime   string              `json:"transaction_time"`
	TransactionStatus string              `json:"transaction_status"`
	FraudStatus       string              `json:"fraud_status"`
	Actions           []map[string]string `json:"actions,omitempty"`
	VaNumbers         []map[string]string `json:"va_numbers,omitempty"`
	PermataVaNumber   string              `json:"permata_va_number,omitempty"`
	BillerCode        string              `json:"biller_code,omitempty"`
	BillKey           string              `json:"bill_key,omitempty"`
	QrCodeUrl         string              `json:"qr_code_url,omitempty"`
	DeepLinkUrl       string              `json:"deeplink_url,omitempty"`
	RedirectUrl       string              `json:"redirect_url,omitempty"`
}

// CreateTransaction creates a new transaction in Midtrans
func (s *MidtransService) CreateTransaction(req *ChargeRequest) (*ChargeResponse, error) {
	url := fmt.Sprintf("%s/v2/charge", s.config.MidtransBaseURL)

	// Convert request to JSON
	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	// Create HTTP request
	httpReq, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	// Set headers
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "application/json")
	httpReq.SetBasicAuth(s.config.MidtransServerKey, "")

	// Send request
	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		var errorResp map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&errorResp); err != nil {
			return nil, fmt.Errorf("failed to decode error response: %v", err)
		}
		return nil, fmt.Errorf("midtrans error: %v", errorResp)
	}

	// Decode response
	var chargeResp ChargeResponse
	if err := json.NewDecoder(resp.Body).Decode(&chargeResp); err != nil {
		return nil, err
	}

	return &chargeResp, nil
}

// GetTransactionStatus gets the status of a transaction
func (s *MidtransService) GetTransactionStatus(transactionID string) (*ChargeResponse, error) {
	url := fmt.Sprintf("%s/v2/%s/status", s.config.MidtransBaseURL, transactionID)

	// Create HTTP request
	httpReq, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	// Set headers
	httpReq.Header.Set("Accept", "application/json")
	httpReq.SetBasicAuth(s.config.MidtransServerKey, "")

	// Send request
	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		var errorResp map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&errorResp); err != nil {
			return nil, fmt.Errorf("failed to decode error response: %v", err)
		}
		return nil, fmt.Errorf("midtrans error: %v", errorResp)
	}

	// Decode response
	var statusResp ChargeResponse
	if err := json.NewDecoder(resp.Body).Decode(&statusResp); err != nil {
		return nil, err
	}

	return &statusResp, nil
}

// RefundTransaction refunds a transaction
func (s *MidtransService) RefundTransaction(transactionID string, amount float64, reason string) error {
	url := fmt.Sprintf("%s/v2/%s/refund", s.config.MidtransBaseURL, transactionID)

	// Create request body
	reqBody := map[string]interface{}{
		"amount": amount,
		"reason": reason,
	}

	// Convert request to JSON
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return err
	}

	// Create HTTP request
	httpReq, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	// Set headers
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "application/json")
	httpReq.SetBasicAuth(s.config.MidtransServerKey, "")

	// Send request
	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		var errorResp map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&errorResp); err != nil {
			return fmt.Errorf("failed to decode error response: %v", err)
		}
		return fmt.Errorf("midtrans error: %v", errorResp)
	}

	return nil
}
