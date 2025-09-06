package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	otelinit "otel-ecom-tracing/order-api/internal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.opentelemetry.io/contrib/instrumentation/github.com/labstack/echo/otelecho"
)

func main() {
	port := os.Getenv("HTTP_PORT")
	if port == "" {
		port = "5002"
	}

	// ---- OpenTelemetry bootstrapping ----
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	shutdownOTel, err := otelinit.Init(ctx)
	if err != nil {
		log.Fatalf("otel init: %v", err)
	}
	defer func() {
		c, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = shutdownOTel(c)
	}()

	// ---- Echo ----
	e := echo.New()
	e.HideBanner = true
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(otelecho.Middleware(os.Getenv("SERVICE_NAME")))

	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"status":  "ok",
			"service": "order-api",
		})
	})

	// Start + graceful shutdown
	go func() {
		log.Printf("order-api listening on :%s", port)
		if err := e.Start(":" + port); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal(err)
		}
	}()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	ctxShutdown, cancel2 := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel2()
	if err := e.Shutdown(ctxShutdown); err != nil {
		e.Logger.Fatal(err)
	}
}
