package handler

import "net/http"

type PaymentHandler interface {
	CreateCheckoutSession(w http.ResponseWriter, r *http.Request)
	GetSubscriptionManagementLink(w http.ResponseWriter, r *http.Request)
	HandleWebhook(w http.ResponseWriter, r *http.Request)
}
