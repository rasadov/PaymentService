package handler

import "net/http"

// PaymentHandler defines the interface for handling payment-related HTTP requests.
// It provides methods for checkout sessions, subscription management, and webhook processing.
type PaymentHandler interface {
	// CreateCheckoutSession creates a new checkout session for a customer.
	// Accepts POST requests with customer details and product ID.
	CreateCheckoutSession(w http.ResponseWriter, r *http.Request)

	// GetSubscriptionManagementLink generates a link for subscription management.
	// Accepts POST requests with customer ID.
	GetSubscriptionManagementLink(w http.ResponseWriter, r *http.Request)

	// HandleWebhook processes incoming webhook events from payment providers.
	// Requires webhook signature verification headers.
	HandleWebhook(w http.ResponseWriter, r *http.Request)
}
