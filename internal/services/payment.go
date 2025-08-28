package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/rasadov/PaymentService/internal/config"
	"github.com/rasadov/PaymentService/internal/db"
	"github.com/rasadov/PaymentService/internal/dto"
	"github.com/rasadov/PaymentService/internal/models"
	"github.com/rasadov/PaymentService/internal/payments"
)

type paymentService struct {
	storage       db.Storage
	paymentClient payments.PaymentClient
}

func NewPaymentService(storage db.Storage, paymentClient payments.PaymentClient) PaymentService {
	return &paymentService{storage: storage, paymentClient: paymentClient}
}

func (s *paymentService) CreateCheckoutSession(customer models.CustomerInput, products []models.Product, metadata map[string]any) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return s.paymentClient.CreateCheckoutSession(ctx, customer, products, metadata)
}

func (s *paymentService) GetSubscriptionManagementLink(customerId string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return s.paymentClient.GetCustomerPortalSession(ctx, customerId)
}

func (p *paymentService) SendWebhookDataToService(webhookId string, payload dto.DodoWebhookPayload) error {
	// Check if already processed
	if status, err := p.storage.Get(webhookId); err == nil && status == "processed" {
		return nil
	}

	// Mark as processed immediately to prevent duplicates
	if err := p.storage.PutWithExpiration(webhookId, "processed", 24*time.Hour); err != nil {
		return fmt.Errorf("failed to mark webhook as processed: %w", err)
	}

	response := dto.PaymentProcessorResponse{
		SubscriptionID: payload.Data.SubscriptionID,
		Status:         payload.Data.Status,
		ProductID:      payload.Data.ProductID,
		Customer: models.Customer{
			CustomerId: payload.Data.Customer.CustomerId,
			Email:      payload.Data.Customer.Email,
			Name:       payload.Data.Customer.Name,
		},
		Metadata: payload.Data.Metadata,
	}

	// Marshal the response to JSON
	jsonData, err := json.Marshal(response)
	if err != nil {
		return fmt.Errorf("failed to marshal response: %w", err)
	}

	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	// Retry logic
	maxRetries := 3
	baseWaitTime := 1 * time.Second
	maxWaitTime := 10 * time.Second

	for attempt := 0; attempt <= maxRetries; attempt++ {
		// Create request
		req, err := http.NewRequest("POST", config.GetConfig().PaymentCallbackUrl, bytes.NewBuffer(jsonData))
		if err != nil {
			return fmt.Errorf("failed to create request: %w", err)
		}

		// Set headers
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("User-Agent", "PaymentService/1.0")

		// Make request
		resp, err := client.Do(req)
		if err != nil {
			// Retry on network errors
			if attempt < maxRetries {
				waitTime := time.Duration(attempt+1) * baseWaitTime
				if waitTime > maxWaitTime {
					waitTime = maxWaitTime
				}
				time.Sleep(waitTime)
				continue
			}
			return fmt.Errorf("failed to send request after %d retries: %w", maxRetries, err)
		}

		// Read response body
		body, readErr := io.ReadAll(resp.Body)
		resp.Body.Close()

		// Check status code
		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			return nil // Success
		}

		// Retry on 5xx status codes
		if resp.StatusCode >= 500 && attempt < maxRetries {
			waitTime := time.Duration(attempt+1) * baseWaitTime
			if waitTime > maxWaitTime {
				waitTime = maxWaitTime
			}
			time.Sleep(waitTime)
			continue
		}

		// Non-retryable error or max retries reached
		bodyStr := "unable to read response"
		if readErr == nil {
			bodyStr = string(body)
		}
		return fmt.Errorf("received non-2xx status code: %d, body: %s", resp.StatusCode, bodyStr)
	}

	return fmt.Errorf("max retries exceeded")
}
