// Файл document-server/internal/infrastructure/observability/tracer.go содержит реализацию пакета observability.
package observability

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

// InitializeTracer инициализирует OpenTelemetry трассировку
func InitializeTracer(ctx context.Context, serviceName, jaegerEndpoint string, logger *zap.Logger) (trace.TracerProvider, func(), error) {
	// Для разработки используем SDK trace provider без экспортера
	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceNameKey.String(serviceName),
			semconv.ServiceVersionKey.String("1.0.0"),
		),
	)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create resource: %w", err)
	}

	// Создание trace provider без exporter (для разработки)
	tp := tracesdk.NewTracerProvider(
		tracesdk.WithResource(res),
	)

	// Установка глобального trace provider
	otel.SetTracerProvider(tp)

	// Функция для graceful shutdown
	shutdown := func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			logger.Error("Failed to shutdown tracer provider", zap.Error(err))
		}
	}

	logger.Info("OpenTelemetry tracer initialized",
		zap.String("service", serviceName),
		zap.String("endpoint", jaegerEndpoint))

	return tp, shutdown, nil
}
