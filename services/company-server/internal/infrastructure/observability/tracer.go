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

// InitializeTracer инициализирует OpenTelemetry трассировку с NoOp экспортером
// На production используйте реальный exporter (Jaeger, Zipkin и т.д.)
func InitializeTracer(ctx context.Context, serviceName, jaegerEndpoint string, logger *zap.Logger) (trace.TracerProvider, func(), error) {
	// Для разработки используем SDK trace provider без экспортера
	// Это позволит избежать проблем с зависимостями
	// На production нужно заменить на реальный exporter

	// Создание ресурса с информацией о сервисе
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
	// TODO: На production добавить реальный exporter (Jaeger, Zipkin и т.д.)
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
		zap.String("jaeger_endpoint", jaegerEndpoint),
	)

	return tp, shutdown, nil
}

// GetTracer возвращает tracer для пакета
func GetTracer(packageName string) trace.Tracer {
	return otel.Tracer(packageName)
}
