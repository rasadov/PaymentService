package main

import (
	"bytes"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/syumai/workers"

	"github.com/rasadov/PaymentService/internal/config"
)

func main() {
	err := config.LoadConfig()
	if err != nil {
		panic(err)
	}
	engine := gin.New()
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
