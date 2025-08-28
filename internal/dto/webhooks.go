package dto

import "github.com/rasadov/PaymentService/internal/models"

// DodoWebhookPayload is the payload received from Dodo payments api
// https://docs.dodopayments.com/developer-resources/webhooks#request-body
type DodoWebhookPayload struct {
	EventType string `json:"type"`
	Data      struct {
		SubscriptionID string          `json:"subscription_id"`
		Status         string          `json:"status"`
		ProductID      string          `json:"product_id"`
		Customer       models.Customer `json:"customer"`
		Metadata       map[string]any  `json:"metadata,omitempty"`
	} `json:"data"`
}

// PaymentProcessorResponse is the payload this service will
// send to your api after processing the webhook
type PaymentProcessorResponse struct {
	SubscriptionID string          `json:"subscription_id"`
	Status         string          `json:"status"`
	ProductID      string          `json:"product_id"`
	Customer       models.Customer `json:"customer"`
	Metadata       map[string]any  `json:"metadata,omitempty"`
}
