package tests

import (
	"io"
	"net/http"
	"testing"

	"github.com/rasadov/PaymentService/internal/config"
	"github.com/rasadov/PaymentService/internal/db"
	"github.com/rasadov/PaymentService/internal/handler"
	"github.com/rasadov/PaymentService/internal/payments"
	"github.com/rasadov/PaymentService/internal/services"
)

func init() {
	config.LoadConfig()
}

func TestCreateCheckoutSession(t *testing.T) {
	storage, err := db.GetConnection()
	if err != nil {
		t.Fatal("Failed to get database connection:", err)
	}

	paymentClient := payments.NewDodoClient(
		config.GetConfig().DodoAPIKey,
		config.GetConfig().Environment == "development",
	)

	paymentService := services.NewPaymentService(storage, paymentClient)
	paymentHandler := handler.NewPaymentHandler(paymentService)

	mux := http.NewServeMux()
	handler.SetupRoutes(mux, paymentHandler)

	req, err := http.NewRequest("POST", "/checkout", nil)
	if err != nil {
		t.Fatal("Failed to create request:", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal("Failed to send request:", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal("Failed to read response body:", err)
	}

	if string(body) != "OK" {
		t.Errorf("Expected response body 'OK', got '%s'", string(body))
	}
}

func TestGetSubscriptionManagementLink(t *testing.T) {

}

func TestSendWebhookDataToService(t *testing.T) {
}
