package dto

type GetCheckoutUrlRequest struct {
	Email     string `json:"email"`
	Name      string `json:"name"`
	ProductID string `json:"product_id"`
}

type GetSubscriptionManagementLinkRequest struct {
	CustomerId string `json:"customer_id"`
}

type UrlResponse struct {
	Url string `json:"url"`
}
