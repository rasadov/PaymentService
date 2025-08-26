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
	var request dto.GetCheckoutUrlRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	url := h.service.CreateCheckoutSession(request.Email, request.Name, request.ProductID)
	c.JSON(http.StatusOK, gin.H{"url": url})
}

func (h *paymentHandler) GetSubscriptionManagementLink(c *gin.Context) {
	var request dto.GetSubscriptionManagementLinkRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	url, err := h.service.GetSubscriptionManagementLink(request.CustomerId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get subscription management link"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"url": url})
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

	webhookId := c.GetHeader("webhook-id")
	if webhookId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing webhook-id header"})
		return
	}

	err = h.service.SendWebhookDataToService(webhookId, payload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process webhook"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Webhook processed"})
}
