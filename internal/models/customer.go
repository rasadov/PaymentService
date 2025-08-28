package models

type Customer struct {
	CustomerId string `json:"customer_id"`
	Email      string `json:"email"`
	Name       string `json:"name"`
}

type CustomerInput struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}
