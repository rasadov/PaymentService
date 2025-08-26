package dto

type WebhookPayload struct {
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
