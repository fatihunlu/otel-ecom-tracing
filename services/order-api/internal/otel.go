package otelinit

import (
	"context"
	"fmt"
	"os"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	"google.golang.org/grpc"
)

type Shutdown func(context.Context) error

func Init(ctx context.Context) (Shutdown, error) {
	endpoint := getenv("OTEL_EXPORTER_OTLP_ENDPOINT", "http://otel-collector:4317")

	// grpc target host:port olmalı; http:// varsa kırp
	grpcTarget := endpoint
	if len(endpoint) > 7 && endpoint[:7] == "http://" {
		grpcTarget = endpoint[7:]
	} else if len(endpoint) > 8 && endpoint[:8] == "https://" {
		grpcTarget = endpoint[8:]
	}

	svcName := getenv("SERVICE_NAME", "ecom.order.api")
	svcVersion := getenv("SERVICE_VERSION", "0.1.0")
	env := getenv("ENVIRONMENT", "local")

	exp, err := otlptracegrpc.New(ctx,
		otlptracegrpc.WithEndpoint(grpcTarget),
		otlptracegrpc.WithInsecure(),
		otlptracegrpc.WithDialOption(grpc.WithBlock()),
	)
	if err != nil {
		return nil, fmt.Errorf("otlp exporter: %w", err)
	}

	res, err := resource.New(ctx,
		resource.WithFromEnv(),
		resource.WithAttributes(
			semconv.ServiceName(svcName),
			semconv.ServiceVersion(svcVersion),
			semconv.DeploymentEnvironment(env),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("resource: %w", err)
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exp),
		sdktrace.WithResource(res),
	)

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.TraceContext{})

	return func(c context.Context) error {
		_ = tp.ForceFlush(c)
		return tp.Shutdown(c)
	}, nil
}

func getenv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}
