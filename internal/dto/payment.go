package dto

import "github.com/rasadov/PaymentService/internal/models"

type GetCheckoutUrlRequest struct {
	Customer    models.CustomerInput `json:"customer"`
	ProductCart []models.Product     `json:"product_cart"`
	Metadata    map[string]any       `json:"metadata,omitempty"`
}

type CheckoutResponse struct {
	CheckoutURL string `json:"checkout_url"`
}

type GetSubscriptionManagementLinkRequest struct {
	CustomerId string `json:"customer_id"`
}

type UrlResponse struct {
	Url string `json:"url"`
}
