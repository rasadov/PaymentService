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
	mux := http.NewServeMux()

	setupHealthChecks(mux)

	// Try to load config - if it fails, only serve health endpoints
	err := config.LoadConfig()
	if err != nil {
		setupErrorHandler(mux, "Configuration error: "+err.Error())
		return &App{Handler: mux}, nil
	}

	storage, err := db.GetConnection()
	if err != nil {
		setupErrorHandler(mux, "Database connection error: "+err.Error())
		return &App{Handler: mux}, nil
	}

	paymentClient := payments.NewDodoClient(
		config.GetConfig().DodoAPIKey,
		config.GetConfig().Environment == "development",
	)

	paymentService := services.NewPaymentService(storage, paymentClient)
	paymentHandler := handler.NewPaymentHandler(paymentService)

	handler.SetupRoutes(mux, paymentHandler)

	return &App{Handler: mux}, nil
}

func setupErrorHandler(mux *http.ServeMux, errorMsg string) {
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/health" || r.URL.Path == "/echo" {
			return // Let health checks handle these
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusServiceUnavailable)
		w.Write([]byte(`{"error": "` + errorMsg + `"}`))
	})
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
