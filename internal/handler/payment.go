package handler

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
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

func (h *paymentHandler) CreateCheckoutSession(c *gin.Context) {
}

func (h *paymentHandler) GetSubscriptionManagementLink(c *gin.Context) {
}

func (h *paymentHandler) GetSubscriptionStatus(c *gin.Context) {
}

func (h *paymentHandler) HandleWebhook(c *gin.Context) {
	body, err := c.GetRawData()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read request body"})
		return
	}

	var payload dto.DodoWebhookPayload
	if err = json.Unmarshal(body, &payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid webhook payload"})
		return
	}

	signature := c.GetHeader("webhook-signature")
	if pkg.VerifyWebhookSignature(signature, body) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid webhook signature"})
		return
	}

	err = h.service.SendWebhookDataToService(payload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process webhook"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Webhook processed"})
}
