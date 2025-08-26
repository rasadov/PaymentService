package services

import (
	"fmt"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/rasadov/PaymentService/internal/config"
	"github.com/rasadov/PaymentService/internal/db"
	"github.com/rasadov/PaymentService/internal/dto"
)

type paymentService struct {
	storage db.Storage
}

func NewPaymentService(storage db.Storage) PaymentService {
	return &paymentService{storage: storage}
}

func (p *paymentService) CreateCheckoutSession() (string, error) {
	return "", nil
}

func (p *paymentService) GetSubscriptionManagementLink() (string, error) {
	return "", nil
}

func (p *paymentService) GetSubscriptionStatus() (string, error) {
	return "", nil
}

func (p *paymentService) SendWebhookDataToService(payload dto.DodoWebhookPayload) error {
	response := dto.PaymentProcessorResponse{
		SubscriptionID: payload.Data.SubscriptionID,
		Status:         payload.Data.Status,
		ProductID:      payload.Data.ProductID,
		Customer: dto.Customer{
			CustomerID: payload.Data.Customer.CustomerID,
			Email:      payload.Data.Customer.Email,
			Name:       payload.Data.Customer.Name,
		},
	}

	client := resty.New().
		SetTimeout(30 * time.Second).
		SetRetryCount(3).
		SetRetryWaitTime(1 * time.Second).
		SetRetryMaxWaitTime(10 * time.Second).
		AddRetryCondition(func(r *resty.Response, err error) bool {
			// Retry on network errors or 5xx status codes
			return err != nil || r.StatusCode() >= 500
		})

	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("User-Agent", "PaymentService/1.0").
		SetBody(response).
		Post(config.GetConfig().PaymentCallbackUrl)

	if err != nil {
		return fmt.Errorf("failed to send request after retries: %w", err)
	}

	if resp.StatusCode() < 200 || resp.StatusCode() >= 300 {
		return fmt.Errorf("received non-2xx status code: %d, body: %s", resp.StatusCode(), resp.String())
	}

	return nil
}
