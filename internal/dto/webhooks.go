package dto

type Customer struct {
	CustomerID string `json:"customer_id"`
	Email      string `json:"email"`
	Name       string `json:"name"`
}

// DodoWebhookPayload is the payload received from Dodo payments api
// https://docs.dodopayments.com/developer-resources/webhooks/intents/subscription
type DodoWebhookPayload struct {
	EventType string `json:"type"`
	Data      struct {
		SubscriptionID string `json:"subscription_id"`
		Status         string `json:"status"`
		ProductID      string `json:"product_id"`
		Customer       Customer
	} `json:"data"`
}

// PaymentProcessorResponse is the payload this service will
// send to your api after processing the webhook
type PaymentProcessorResponse struct {
	SubscriptionID string `json:"subscription_id"`
	Status         string `json:"status"`
	ProductID      string `json:"product_id"`
	Customer       Customer
}
