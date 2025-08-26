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
