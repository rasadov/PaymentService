package main

import (
	"bytes"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/syumai/workers"

	"github.com/rasadov/PaymentService/internal/config"
	"github.com/rasadov/PaymentService/internal/db"
	"github.com/rasadov/PaymentService/internal/handler"
	"github.com/rasadov/PaymentService/internal/services"
)

func init() {
	err := config.LoadConfig()
	if err != nil {
		panic(err)
	}
}

func main() {
	storage, err := db.GetConnection()
	if err != nil {
		panic(err)
	}

	paymentService := services.NewPaymentService(storage)

	paymentHandler := handler.NewPaymentHandler(paymentService)

	engine := gin.New()

	handler.SetupRoutes(engine, paymentHandler)

	engine.GET("/hello", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello!")
	})
	engine.GET("/echo", func(c *gin.Context) {
		b, err := io.ReadAll(c.Request.Body)
		if err != nil {
			panic(err)
		}
		io.Copy(c.Writer, bytes.NewReader(b))
	})
	workers.Serve(engine)
}
