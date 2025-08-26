package services

type PaymentService interface {
	CreateCheckoutSession() (string, error)
	GetSubscriptionManagementLink() (string, error)
	GetSubscriptionStatus() (string, error)
	HandleWebhook()
}
