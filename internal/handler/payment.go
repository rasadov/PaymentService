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

// CreateCheckoutSession creates a new checkout session for a customer to purchase a product.
// It accepts a POST request with customer email, name, and product ID.
//
// Request body:
//
//	{
//	  "customer": {
//		"email": "customer@example.com",
//		"name": "John Doe"
//	  },
//	  "product_cart": [
//		{
//			"quantity": 1,
//			"product_id": "prod_123456"
//		}
//	  ],
//	  "metadata": {
//		"key": "value"
//	  }
//	}
//
// Response:
//
//	{
//	  "url": "https://checkout.dodopayments.com/session_abc123"
//	}
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
	response := map[string]string{"url": url}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetSubscriptionManagementLink retrieves a link for customers to manage their subscription.
// It accepts a POST request with the customer ID.
//
// Request body:
//
//	{
//	  "customer_id": "cus_123456789"
//	}
//
// Response:
//
//	{
//	  "url": "https://billing.dodopayments.com/manage/cus_123456789"
//	}
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

	response := map[string]string{"url": url}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// HandleWebhook processes webhook events from Dodo Payments.
// It verifies the webhook signature and forwards the data to the configured service.
//
// Required headers:
//   - webhook-signature: Webhook signature for verification
//   - webhook-id: Unique webhook identifier
//
// Request body: DodoWebhookPayload with event type and subscription data
//
// Response:
//
//	{
//	  "message": "Webhook processed"
//	}
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
