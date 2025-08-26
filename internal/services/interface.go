package services

import "github.com/rasadov/PaymentService/internal/dto"

type PaymentService interface {
	CreateCheckoutSession() (string, error)
	GetSubscriptionManagementLink() (string, error)
	GetSubscriptionStatus() (string, error)
	SendWebhookDataToService(payload dto.DodoWebhookPayload) error
}
