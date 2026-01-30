// Файл bank-account-query-server/internal/infrastructure/observability/tracing.go содержит реализацию пакета observability.
package observability

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
)

// InitTracer инициализирует OpenTelemetry tracer provider для Jaeger
func InitTracer(serviceName, jaegerEndpoint string) (trace.TracerProvider, error) {
	// Create Jaeger exporter
	exporter, err := jaeger.New(
		jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(jaegerEndpoint)),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create Jaeger exporter: %w", err)
	}

	// Create resource with service name
	res, err := resource.New(
		context.Background(),
		resource.WithAttributes(
			semconv.ServiceNameKey.String(serviceName),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}

	// Create trace provider
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
	)

	// Set global trace provider
	otel.SetTracerProvider(tp)

	return tp, nil
}

// ShutdownTracer закрывает tracer provider gracefully
func ShutdownTracer(ctx context.Context, tp trace.TracerProvider) error {
	if sdkTP, ok := tp.(*sdktrace.TracerProvider); ok {
		return sdkTP.Shutdown(ctx)
	}
	return nil
}
