// Файл api-gateway/internal/infrastructure/observability/tracer.go содержит реализацию пакета observability.
package observability

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

// InitTracer инициализирует OpenTelemetry трассировку
func InitTracer(ctx context.Context, serviceName, jaegerEndpoint string, logger *zap.Logger) (*sdktrace.TracerProvider, func(), error) {
	// Создаем exporter для Jaeger
	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(jaegerEndpoint)))
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create jaeger exporter: %w", err)
	}

	// Создаем TracerProvider
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exp),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(serviceName),
		)),
	)

	// Устанавливаем глобальный TracerProvider
	otel.SetTracerProvider(tp)

	logger.Info("OpenTelemetry tracer initialized",
		zap.String("service", serviceName),
		zap.String("jaeger_endpoint", jaegerEndpoint),
	)

	// Функция для graceful shutdown
	shutdown := func() {
		if err := tp.Shutdown(ctx); err != nil {
			logger.Error("Failed to shutdown tracer provider", zap.Error(err))
		}
	}

	return tp, shutdown, nil
}

// Tracer обертка для OpenTelemetry tracer
type Tracer struct {
	tracer trace.Tracer
	logger *zap.Logger
}

// NewTracer создает новый Tracer
func NewTracer(serviceName string, logger *zap.Logger) *Tracer {
	return &Tracer{
		tracer: otel.Tracer(serviceName),
		logger: logger,
	}
}

// StartSpan начинает новый span
func (t *Tracer) StartSpan(ctx context.Context, spanName string) (context.Context, func()) {
	ctx, span := t.tracer.Start(ctx, spanName)

	return ctx, func() {
		span.End()
	}
}

// AddEvent добавляет событие в текущий span
func (t *Tracer) AddEvent(ctx context.Context, name string, attributes map[string]interface{}) {
	span := trace.SpanFromContext(ctx)
	if !span.IsRecording() {
		return
	}

	attrs := make([]attribute.KeyValue, 0, len(attributes))
	for k, v := range attributes {
		switch val := v.(type) {
		case string:
			attrs = append(attrs, attribute.String(k, val))
		case int:
			attrs = append(attrs, attribute.Int(k, val))
		case int64:
			attrs = append(attrs, attribute.Int64(k, val))
		case float64:
			attrs = append(attrs, attribute.Float64(k, val))
		case bool:
			attrs = append(attrs, attribute.Bool(k, val))
		}
	}

	span.AddEvent(name, trace.WithAttributes(attrs...))
}

// RecordError записывает ошибку в текущий span
func (t *Tracer) RecordError(ctx context.Context, err error) {
	span := trace.SpanFromContext(ctx)
	if !span.IsRecording() {
		return
	}

	span.RecordError(err)
	span.SetAttributes(attribute.Bool("error", true))
}

// SetAttributes устанавливает атрибуты для текущего span
func (t *Tracer) SetAttributes(ctx context.Context, attributes map[string]interface{}) {
	span := trace.SpanFromContext(ctx)
	if !span.IsRecording() {
		return
	}

	attrs := make([]attribute.KeyValue, 0, len(attributes))
	for k, v := range attributes {
		switch val := v.(type) {
		case string:
			attrs = append(attrs, attribute.String(k, val))
		case int:
			attrs = append(attrs, attribute.Int(k, val))
		case int64:
			attrs = append(attrs, attribute.Int64(k, val))
		case float64:
			attrs = append(attrs, attribute.Float64(k, val))
		case bool:
			attrs = append(attrs, attribute.Bool(k, val))
		}
	}

	span.SetAttributes(attrs...)
}
