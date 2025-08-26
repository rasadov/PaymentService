package app

import (
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rasadov/PaymentService/internal/config"
	"github.com/rasadov/PaymentService/internal/db"
	"github.com/rasadov/PaymentService/internal/handler"
	"github.com/rasadov/PaymentService/internal/services"
)

type App struct {
	Engine *gin.Engine
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

	paymentService := services.NewPaymentService(storage)

	paymentHandler := handler.NewPaymentHandler(paymentService)

	engine := gin.New()
	handler.SetupRoutes(engine, paymentHandler)

	setupHealthChecks(engine)

	return &App{
		Engine: engine,
	}, nil
}

func setupHealthChecks(engine *gin.Engine) {
	engine.GET("/health", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})

	engine.GET("/echo", func(c *gin.Context) {
		b, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read body"})
			return
		}
		c.Data(http.StatusOK, "text/plain", b)
	})
}
