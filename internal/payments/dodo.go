package payments

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/rasadov/PaymentService/internal/dto"
	"github.com/rasadov/PaymentService/internal/models"
)

type DodoClient struct {
	apiKey  string
	baseUrl string
	client  *http.Client
}

type CustomerPortalResponse struct {
	Link string `json:"link"`
}

func NewDodoClient(apiKey string, test bool) PaymentClient {
	if test {
		return &DodoClient{
			apiKey:  apiKey,
			baseUrl: "https://test.dodopayments.com",
			client:  &http.Client{},
		}
	}
	return &DodoClient{
		apiKey:  apiKey,
		baseUrl: "https://live.dodopayments.com",
		client:  &http.Client{},
	}
}

func (dc *DodoClient) CreateCheckoutSession(ctx context.Context, customer models.CustomerInput, products []models.Product, metadata map[string]any) (string, error) {
	url := fmt.Sprintf("%s/checkouts", dc.baseUrl)

	body := map[string]any{
		"customer": map[string]any{
			"email": customer.Email,
			"name":  customer.Name,
		},
		"products": products,
		"metadata": metadata,
	}

	bodyJson, err := json.Marshal(body)
	if err != nil {
		return "", fmt.Errorf("failed to marshal body: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(bodyJson))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+dc.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := dc.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var checkoutResponse dto.CheckoutResponse
	if err := json.NewDecoder(resp.Body).Decode(&checkoutResponse); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	return checkoutResponse.CheckoutURL, nil
}

func (dc *DodoClient) GetCustomerPortalSession(ctx context.Context, customerId string) (string, error) {
	url := fmt.Sprintf("%s/customers/%s/customer-portal/session", dc.baseUrl, customerId)

	req, err := http.NewRequestWithContext(ctx, "POST", url, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Authorization", "Bearer "+dc.apiKey)
	req.Header.Set("Content-Type", "application/json")

	// Make the request
	resp, err := dc.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	// Parse response
	var portalResponse CustomerPortalResponse
	if err := json.NewDecoder(resp.Body).Decode(&portalResponse); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	return portalResponse.Link, nil
}
