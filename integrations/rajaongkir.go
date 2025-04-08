package integrations

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/fashion-shop/config"
)

// RajaOngkirService handles RajaOngkir API integration
type RajaOngkirService struct {
	config *config.Config
}

// NewRajaOngkirService creates a new RajaOngkir service
func NewRajaOngkirService(config *config.Config) *RajaOngkirService {
	return &RajaOngkirService{
		config: config,
	}
}

// Province represents a province from RajaOngkir API
type Province struct {
	ProvinceID string `json:"province_id"`
	Province   string `json:"province"`
}

// City represents a city from RajaOngkir API
type City struct {
	CityID     string `json:"city_id"`
	ProvinceID string `json:"province_id"`
	Province   string `json:"province"`
	Type       string `json:"type"`
	CityName   string `json:"city_name"`
	PostalCode string `json:"postal_code"`
}

// ShippingCost represents shipping cost from RajaOngkir API
type ShippingCost struct {
	Service     string `json:"service"`
	Description string `json:"description"`
	Cost        []struct {
		Value int    `json:"value"`
		Etd   string `json:"etd"`
		Note  string `json:"note"`
	} `json:"cost"`
}

// GetProvinces gets all provinces from RajaOngkir API
func (s *RajaOngkirService) GetProvinces() ([]Province, error) {
	url := fmt.Sprintf("%s/province", s.config.RajaOngkirBaseURL)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("key", s.config.RajaOngkirAPIKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status code %d", resp.StatusCode)
	}

	var response struct {
		Rajaongkir struct {
			Status struct {
				Code        int    `json:"code"`
				Description string `json:"description"`
			} `json:"status"`
			Results []Province `json:"results"`
		} `json:"rajaongkir"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	if response.Rajaongkir.Status.Code != http.StatusOK {
		return nil, fmt.Errorf("API error: %s", response.Rajaongkir.Status.Description)
	}

	return response.Rajaongkir.Results, nil
}

// GetCities gets all cities from RajaOngkir API
func (s *RajaOngkirService) GetCities(provinceID string) ([]City, error) {
	url := fmt.Sprintf("%s/city", s.config.RajaOngkirBaseURL)
	if provinceID != "" {
		url = fmt.Sprintf("%s?province=%s", url, provinceID)
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("key", s.config.RajaOngkirAPIKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status code %d", resp.StatusCode)
	}

	var response struct {
		Rajaongkir struct {
			Status struct {
				Code        int    `json:"code"`
				Description string `json:"description"`
			} `json:"status"`
			Results []City `json:"results"`
		} `json:"rajaongkir"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	if response.Rajaongkir.Status.Code != http.StatusOK {
		return nil, fmt.Errorf("API error: %s", response.Rajaongkir.Status.Description)
	}

	return response.Rajaongkir.Results, nil
}

// CalculateShippingCost calculates shipping cost from RajaOngkir API
func (s *RajaOngkirService) CalculateShippingCost(origin, destination, weight int, courier string) ([]ShippingCost, error) {
	apiURL := fmt.Sprintf("%s/cost", s.config.RajaOngkirBaseURL)

	data := url.Values{}
	data.Set("origin", fmt.Sprintf("%d", origin))
	data.Set("destination", fmt.Sprintf("%d", destination))
	data.Set("weight", fmt.Sprintf("%d", weight))
	data.Set("courier", courier)

	req, err := http.NewRequest("POST", apiURL, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}

	req.Header.Set("key", s.config.RajaOngkirAPIKey)
	req.Header.Set("content-type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status code %d", resp.StatusCode)
	}

	var response struct {
		Rajaongkir struct {
			Status struct {
				Code        int    `json:"code"`
				Description string `json:"description"`
			} `json:"status"`
			Results []struct {
				Code  string         `json:"code"`
				Name  string         `json:"name"`
				Costs []ShippingCost `json:"costs"`
			} `json:"results"`
		} `json:"rajaongkir"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	if response.Rajaongkir.Status.Code != http.StatusOK {
		return nil, fmt.Errorf("API error: %s", response.Rajaongkir.Status.Description)
	}

	if len(response.Rajaongkir.Results) == 0 {
		return nil, fmt.Errorf("no shipping costs available")
	}

	return response.Rajaongkir.Results[0].Costs, nil
}

// TrackShipment tracks a shipment from RajaOngkir API
func (s *RajaOngkirService) TrackShipment(waybill, courier string) (map[string]interface{}, error) {
	apiURL := fmt.Sprintf("%s/waybill", s.config.RajaOngkirBaseURL)

	data := url.Values{}
	data.Set("waybill", waybill)
	data.Set("courier", courier)

	req, err := http.NewRequest("POST", apiURL, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}

	req.Header.Set("key", s.config.RajaOngkirAPIKey)
	req.Header.Set("content-type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status code %d", resp.StatusCode)
	}

	var response struct {
		Rajaongkir struct {
			Status struct {
				Code        int    `json:"code"`
				Description string `json:"description"`
			} `json:"status"`
			Result map[string]interface{} `json:"result"`
		} `json:"rajaongkir"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	if response.Rajaongkir.Status.Code != http.StatusOK {
		return nil, fmt.Errorf("API error: %s", response.Rajaongkir.Status.Description)
	}

	return response.Rajaongkir.Result, nil
}
