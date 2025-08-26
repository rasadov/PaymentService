package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/rasadov/PaymentService/internal/services"
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
}
