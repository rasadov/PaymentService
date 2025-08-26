package payments

import (
	"context"

	"github.com/dodopayments/dodopayments-go"
	"github.com/dodopayments/dodopayments-go/option"
	"github.com/rasadov/PaymentService/internal/config"
)

type PaymentClient interface {
	GetCustomerPortalSession(ctx context.Context, customerId string) (string, error)
	CreateCustomer(ctx context.Context, email, name string) (string, error)
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

func (dc *DodoClient) CreateCustomer(ctx context.Context, email, name string) (string, error) {
	customer, err := dc.client.Customers.New(ctx, dodopayments.CustomerNewParams{
		Email: dodopayments.F(email),
		Name:  dodopayments.F(name),
	})
	if err != nil {
		return "", err
	}

	return customer.CustomerID, nil
}
