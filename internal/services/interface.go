package services

import "github.com/rasadov/PaymentService/internal/dto"

type PaymentService interface {
	CreateCheckoutSession(email, name, productID string) string
	GetSubscriptionManagementLink(customerId string) (string, error)
	SendWebhookDataToService(payload dto.DodoWebhookPayload) error
}
