package handler

import "net/http"

func SetupRoutes(
	mux *http.ServeMux,
	payment PaymentHandler) {

	mux.HandleFunc("/checkout", payment.CreateCheckoutSession)
	mux.HandleFunc("/subscriptions", payment.GetSubscriptionManagementLink)
	mux.HandleFunc("/webhook", payment.HandleWebhook)
}
