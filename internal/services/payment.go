package services

import (
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
	return nil
}
