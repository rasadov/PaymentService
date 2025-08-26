package app

import (
	"io"
	"net/http"

	"github.com/rasadov/PaymentService/internal/config"
	"github.com/rasadov/PaymentService/internal/db"
	"github.com/rasadov/PaymentService/internal/handler"
	"github.com/rasadov/PaymentService/internal/payments"
	"github.com/rasadov/PaymentService/internal/services"
)

type App struct {
	Handler http.Handler
}

func New() (*App, error) {
	err := config.LoadConfig()
	if err != nil {
		return nil, err
	}

	storage, err := db.GetConnection()
	if err != nil {
		return nil, err
	}

	paymentClient := payments.NewDodoClient(
		config.GetConfig().DodoAPIKey,
		config.GetConfig().Environment == "development",
	)

	paymentService := services.NewPaymentService(storage, paymentClient)

	paymentHandler := handler.NewPaymentHandler(paymentService)

	mux := http.NewServeMux()
	handler.SetupRoutes(mux, paymentHandler)

	setupHealthChecks(mux)

	return &App{
		Handler: mux,
	}, nil
}

func setupHealthChecks(mux *http.ServeMux) {
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	mux.HandleFunc("/echo", func(w http.ResponseWriter, r *http.Request) {
		b, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Failed to read body", http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "text/plain")
		w.Write(b)
	})
}
