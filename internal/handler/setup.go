package handler

import "github.com/gin-gonic/gin"

func SetupRoutes(
	engine *gin.Engine,
	payment PaymentHandler) {

	router := engine.Group("/payment")
	router.GET("/checkout", payment.CreateCheckoutSession)
	router.GET("/subscriptions", payment.GetSubscriptionManagementLink)
	router.GET("/status/:user_id", payment.GetSubscriptionStatus)
	router.POST("/webhook", payment.HandleWebhook)
}
