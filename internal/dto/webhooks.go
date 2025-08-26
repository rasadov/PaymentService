package dto

type DodoWebhookPayload struct {
	EventType string `json:"type"`
	Data      struct {
		SubscriptionID string `json:"subscription_id"`
		Status         string `json:"status"`
		ProductID      string `json:"product_id"`
		Customer       struct {
			CustomerID string `json:"customer_id"`
			Email      string `json:"email"`
			Name       string `json:"name"`
		}
	} `json:"data"`
}

// PaymentProcessorResponse is the response this service will
// send to your api after processing the webhook
type PaymentProcessorResponse struct {
	SubscriptionID string `json:"subscription_id"`
	Status         string `json:"status"`
	ProductID      string `json:"product_id"`
	Customer       struct {
		CustomerID string `json:"customer_id"`
		Email      string `json:"email"`
		Name       string `json:"name"`
	}
}
