// Package main provides the entry point for the Payment Service API
//
//	@title			Payment Service API
//	@version		1.0
//	@description	A payment processing service that handles checkout sessions, subscription management, and webhooks
//	@termsOfService	http://swagger.io/terms/
//
//	@license.name	MIT
//	@license.url	https://opensource.org/licenses/MIT
//
//	@host		localhost:8080
//	@BasePath	/
//
//	@schemes	https http
package main

import (
	"log"

	"github.com/rasadov/PaymentService/internal/app"
	"github.com/syumai/workers"
)

func main() {
	application, err := app.New()
	if err != nil {
		log.Fatal("Failed to initialize application:", err)
	}

	workers.Serve(application.Handler)
}
