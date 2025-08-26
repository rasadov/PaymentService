package handler

import "github.com/gin-gonic/gin"

type PaymentHandler interface {
	CreateCheckoutSession(c *gin.Context)
	GetSubscriptionManagementLink(c *gin.Context)
	GetSubscriptionStatus(c *gin.Context)
	HandleWebhook(c *gin.Context)
}
