package payments

import (
	"context"

	"github.com/dodopayments/dodopayments-go"
	"github.com/dodopayments/dodopayments-go/option"
	"github.com/rasadov/PaymentService/internal/config"
)

type PaymentClient interface {
	GetCustomerPortalSession(ctx context.Context, customerId string) (string, error)
}

type DodoClient struct {
	client        *dodopayments.Client
	webhookSecret string
}

func NewDodoClient(apiKey string, webhookSecret string) PaymentClient {
	if config.GetConfig().Environment == "development" {
		return &DodoClient{
			client: dodopayments.NewClient(
				option.WithBearerToken(apiKey),
				option.WithEnvironmentTestMode(),
			),
			webhookSecret: webhookSecret,
		}
	}

	return &DodoClient{
		client:        dodopayments.NewClient(option.WithBearerToken(apiKey)),
		webhookSecret: webhookSecret,
	}
}

func (dc *DodoClient) GetCustomerPortalSession(ctx context.Context, customerId string) (string, error) {
	customer, err := dc.client.Customers.CustomerPortal.New(ctx, customerId, dodopayments.CustomerCustomerPortalNewParams{})
	if err != nil {
		return "", err
	}
	return customer.Link, nil
}
