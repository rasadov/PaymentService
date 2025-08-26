package payments

import "context"

type PaymentClient interface {
	GetCustomerPortalSession(ctx context.Context, customerId string) (string, error)
}
