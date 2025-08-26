package services

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/rasadov/PaymentService/internal/config"
	"github.com/rasadov/PaymentService/internal/db"
	"github.com/rasadov/PaymentService/internal/dto"
	"github.com/rasadov/PaymentService/internal/payments"
)

type paymentService struct {
	storage       db.Storage
	paymentClient payments.PaymentClient
}

func NewPaymentService(storage db.Storage, paymentClient payments.PaymentClient) PaymentService {
	return &paymentService{storage: storage, paymentClient: paymentClient}
}

func (s *paymentService) CreateCheckoutSession(email, name, productID string) string {
	checkoutURL := fmt.Sprintf(config.GetConfig().DodoCheckoutURL,
		productID,
		url.QueryEscape(email),
		url.QueryEscape(name),
		url.QueryEscape(config.GetConfig().DodoCheckoutRedirectUrl),
	)

	return checkoutURL
}

func (s *paymentService) GetSubscriptionManagementLink(customerId string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return s.paymentClient.GetCustomerPortalSession(ctx, customerId)
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
