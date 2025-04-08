package third_party

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

// RajaOngkirService defines the interface for RajaOngkir operations
type RajaOngkirService interface {
	GetProvinces() ([]map[string]interface{}, error)
	GetCities(provinceID string) ([]map[string]interface{}, error)
	CalculateShipping(origin, destination, weight int, courier string) ([]map[string]interface{}, error)
	TrackShipment(waybill, courier string) (map[string]interface{}, error)
}

type rajaOngkirService struct {
	apiKey string
	url    string
}

// NewRajaOngkirService creates a new RajaOngkirService instance
func NewRajaOngkirService(apiKey, url string) RajaOngkirService {
	return &rajaOngkirService{
		apiKey: apiKey,
		url:    url,
	}
}

// GetProvinces gets all provinces
func (s *rajaOngkirService) GetProvinces() ([]map[string]interface{}, error) {
	url := fmt.Sprintf("%s/province", s.url)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("key", s.apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("failed to get provinces")
	}

	var result struct {
		Rajaongkir struct {
			Status struct {
				Code        int    `json:"code"`
				Description string `json:"description"`
			} `json:"status"`
			Results []map[string]interface{} `json:"results"`
		} `json:"rajaongkir"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	if result.Rajaongkir.Status.Code != 200 {
		return nil, errors.New(result.Rajaongkir.Status.Description)
	}

	return result.Rajaongkir.Results, nil
}

// GetCities gets cities in a province
func (s *rajaOngkirService) GetCities(provinceID string) ([]map[string]interface{}, error) {
	url := fmt.Sprintf("%s/city", s.url)
	if provinceID != "" {
		url = fmt.Sprintf("%s?province=%s", url, provinceID)
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("key", s.apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("failed to get cities")
	}

	var result struct {
		Rajaongkir struct {
			Status struct {
				Code        int    `json:"code"`
				Description string `json:"description"`
			} `json:"status"`
			Results []map[string]interface{} `json:"results"`
		} `json:"rajaongkir"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	if result.Rajaongkir.Status.Code != 200 {
		return nil, errors.New(result.Rajaongkir.Status.Description)
	}

	return result.Rajaongkir.Results, nil
}

// CalculateShipping calculates shipping cost
func (s *rajaOngkirService) CalculateShipping(origin, destination, weight int, courier string) ([]map[string]interface{}, error) {
	url := fmt.Sprintf("%s/cost", s.url)

	data := map[string]interface{}{
		"origin":      origin,
		"destination": destination,
		"weight":      weight,
		"courier":     courier,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("key", s.apiKey)
	req.Header.Set("content-type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("failed to calculate shipping")
	}

	var result struct {
		Rajaongkir struct {
			Status struct {
				Code        int    `json:"code"`
				Description string `json:"description"`
			} `json:"status"`
			Results []map[string]interface{} `json:"results"`
		} `json:"rajaongkir"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	if result.Rajaongkir.Status.Code != 200 {
		return nil, errors.New(result.Rajaongkir.Status.Description)
	}

	return result.Rajaongkir.Results, nil
}

// TrackShipment tracks a shipment
func (s *rajaOngkirService) TrackShipment(waybill, courier string) (map[string]interface{}, error) {
	url := fmt.Sprintf("%s/waybill", s.url)

	data := map[string]interface{}{
		"waybill": waybill,
		"courier": courier,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("key", s.apiKey)
	req.Header.Set("content-type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("failed to track shipment")
	}

	var result struct {
		Rajaongkir struct {
			Status struct {
				Code        int    `json:"code"`
				Description string `json:"description"`
			} `json:"status"`
			Result map[string]interface{} `json:"result"`
		} `json:"rajaongkir"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	if result.Rajaongkir.Status.Code != 200 {
		return nil, errors.New(result.Rajaongkir.Status.Description)
	}

	return result.Rajaongkir.Result, nil
}
