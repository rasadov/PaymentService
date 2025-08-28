package payments

import (
	"context"

	"github.com/rasadov/PaymentService/internal/models"
)

type PaymentClient interface {
	CreateCheckoutSession(ctx context.Context, customer models.CustomerInput, products []models.Product, metadata map[string]any) (string, error)
	GetCustomerPortalSession(ctx context.Context, customerId string) (string, error)
}
