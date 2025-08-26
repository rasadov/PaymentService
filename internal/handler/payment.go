package handler

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/rasadov/PaymentService/internal/dto"
	"github.com/rasadov/PaymentService/internal/services"
	"github.com/rasadov/PaymentService/pkg"
)

type paymentHandler struct {
	service services.PaymentService
}

func NewPaymentHandler(service services.PaymentService) PaymentHandler {
	return &paymentHandler{service: service}
}

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

	url := h.service.CreateCheckoutSession(request.Email, request.Name, request.ProductID)
	response := map[string]string{"url": url}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

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
