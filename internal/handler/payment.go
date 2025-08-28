package handler

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/rasadov/PaymentService/internal/dto"
	"github.com/rasadov/PaymentService/internal/services"
	"github.com/rasadov/PaymentService/pkg"
)

// paymentHandler implements the PaymentHandler interface and handles payment-related HTTP requests.
type paymentHandler struct {
	service services.PaymentService
}

// NewPaymentHandler creates a new payment handler with the provided payment service.
func NewPaymentHandler(service services.PaymentService) PaymentHandler {
	return &paymentHandler{service: service}
}

// CreateCheckoutSession creates a new checkout session for a customer to purchase products
//
//	@Summary		Create checkout session
//	@Description	Creates a new checkout session for a customer to purchase products
//	@Tags			payments
//	@Accept			json
//	@Produce		json
//	@Param			request	body	dto.GetCheckoutUrlRequest	true	"Checkout request"
//	@Success		200	{object}	dto.UrlResponse	"Checkout URL"
//	@Failure		400	{string}	string	"Bad request"
//	@Failure		405	{string}	string	"Method not allowed"
//	@Failure		500	{string}	string	"Internal server error"
//	@Router			/checkout [post]
func (h *paymentHandler) CreateCheckoutSession(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var request dto.GetCheckoutUrlRequest
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}

	if err := json.Unmarshal(body, &request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	url, err := h.service.CreateCheckoutSession(request.Customer, request.ProductCart, request.Metadata)
	if err != nil {
		http.Error(w, "Failed to create checkout session", http.StatusInternalServerError)
		return
	}
	response := dto.UrlResponse{Url: url}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetSubscriptionManagementLink retrieves a link for customers to manage their subscription
//
//	@Summary		Get subscription management link
//	@Description	Retrieves a link for customers to manage their subscription
//	@Tags			subscriptions
//	@Accept			json
//	@Produce		json
//	@Param			request	body	dto.GetSubscriptionManagementLinkRequest	true	"Customer ID request"
//	@Success		200	{object}	dto.UrlResponse	"Management URL"
//	@Failure		400	{string}	string	"Bad request"
//	@Failure		405	{string}	string	"Method not allowed"
//	@Failure		500	{string}	string	"Internal server error"
//	@Router			/subscriptions [post]
func (h *paymentHandler) GetSubscriptionManagementLink(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var request dto.GetSubscriptionManagementLinkRequest
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}

	if err := json.Unmarshal(body, &request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	url, err := h.service.GetSubscriptionManagementLink(request.CustomerId)
	if err != nil {
		http.Error(w, "Failed to get subscription management link", http.StatusInternalServerError)
		return
	}

	response := dto.UrlResponse{Url: url}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// HandleWebhook processes webhook events from Dodo Payments
//
//	@Summary		Handle webhook
//	@Description	Processes webhook events from Dodo Payments and forwards data to configured service
//	@Tags			webhooks
//	@Accept			json
//	@Produce		json
//	@Param			webhook-signature	header	string	true	"Webhook signature for verification"
//	@Param			webhook-id	header	string	true	"Unique webhook identifier"
//	@Param			payload	body	dto.DodoWebhookPayload	true	"Webhook payload"
//	@Success		200	{object}	map[string]string	"Success message"
//	@Failure		400	{string}	string	"Bad request"
//	@Failure		401	{string}	string	"Invalid signature"
//	@Failure		405	{string}	string	"Method not allowed"
//	@Failure		500	{string}	string	"Internal server error"
//	@Router			/webhook [post]
func (h *paymentHandler) HandleWebhook(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}

	var payload dto.DodoWebhookPayload
	if err = json.Unmarshal(body, &payload); err != nil {
		http.Error(w, "Invalid webhook payload", http.StatusBadRequest)
		return
	}

	signature := r.Header.Get("webhook-signature")
	if pkg.VerifyWebhookSignature(signature, body) {
		http.Error(w, "Invalid webhook signature", http.StatusUnauthorized)
		return
	}

	webhookId := r.Header.Get("webhook-id")
	if webhookId == "" {
		http.Error(w, "Missing webhook-id header", http.StatusBadRequest)
		return
	}

	err = h.service.SendWebhookDataToService(webhookId, payload)
	if err != nil {
		http.Error(w, "Failed to process webhook", http.StatusInternalServerError)
		return
	}

	response := map[string]string{"message": "Webhook processed"}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
