package handler

import "github.com/gin-gonic/gin"

func SetupRoutes(
	engine *gin.Engine,
	payment PaymentHandler) {

	router := engine.Group("/payment")
	router.GET("/create-checkout-session", payment.CreateCheckoutSession)
	router.GET("/subscription-management-link", payment.GetSubscriptionManagementLink)
	router.GET("/subscription-status", payment.GetSubscriptionStatus)
	router.POST("/webhook", payment.HandleWebhook)
}
