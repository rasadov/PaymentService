package payments

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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
	if resp.StatusCode != http.StatusOK {
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
