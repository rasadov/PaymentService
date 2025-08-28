package services

import (
	"github.com/rasadov/PaymentService/internal/dto"
	"github.com/rasadov/PaymentService/internal/models"
)

type PaymentService interface {
	CreateCheckoutSession(customer models.CustomerInput, products []models.Product, metadata map[string]any) (string, error)
	GetSubscriptionManagementLink(customerId string) (string, error)
	SendWebhookDataToService(webhookId string, payload dto.DodoWebhookPayload) error
}
