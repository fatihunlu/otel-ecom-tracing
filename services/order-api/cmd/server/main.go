package main

import (
	"log"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	port := os.Getenv("HTTP_PORT")
	if port == "" {
		port = "5002"
	}

	e := echo.New()

	// Logging & Recovery middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Health endpoint
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"status":  "ok",
			"service": "order-api",
		})
	})

	log.Printf("order-api listening on :%s", port)
	if err := e.Start(":" + port); err != nil {
		log.Fatal(err)
	}
}
