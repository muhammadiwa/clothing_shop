package rajaongkir

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"clothing-shop-api/internal/config"
)

type RajaOngkirClient struct {
	BaseURL    string
	APIKey     string
	HTTPClient *http.Client
}

type ShippingCostRequest struct {
	Origin      string `json:"origin"`      // City ID
	Destination string `json:"destination"` // City ID
	Weight      int    `json:"weight"`      // Weight in grams
	Courier     string `json:"courier"`     // Courier code: jne, tiki, pos, etc.
}

type ShippingCostResponse struct {
	Code    int    `json:"code"`
	Status  string `json:"status"`
	Results []struct {
		Code  string `json:"code"`
		Name  string `json:"name"`
		Costs []struct {
			Service     string `json:"service"`
			Description string `json:"description"`
			Cost        []struct {
				Value int    `json:"value"`
				Etd   string `json:"etd"`
				Note  string `json:"note"`
			} `json:"cost"`
		} `json:"costs"`
	} `json:"results"`
}

func NewRajaOngkirClient(cfg *config.Config) *RajaOngkirClient {
	return &RajaOngkirClient{
		BaseURL: cfg.RajaOngkirBaseURL,
		APIKey:  cfg.RajaOngkirAPIKey,
		HTTPClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *RajaOngkirClient) GetShippingCost(req ShippingCostRequest) (*ShippingCostResponse, error) {
	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequest("POST", c.BaseURL+"/cost", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	httpReq.Header.Set("key", c.APIKey)
	httpReq.Header.Set("content-type", "application/json")

	resp, err := c.HTTPClient.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("rajaongkir API returned non-200 status: %d", resp.StatusCode)
	}

	var apiResp struct {
		Rajaongkir ShippingCostResponse `json:"rajaongkir"`
	}
	if err = json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, err
	}

	return &apiResp.Rajaongkir, nil
}
