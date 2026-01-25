package observability

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// TracingProvider оборачивает OpenTelemetry tracer для удобного использования
type TracingProvider struct {
	tracer trace.Tracer
}

// NewTracingProvider создаёт новый провайдер трассирования
func NewTracingProvider(name string) *TracingProvider {
	return &TracingProvider{
		tracer: otel.Tracer(name),
	}
}

// StartSpan начинает новый span с заданным именем и атрибутами
func (tp *TracingProvider) StartSpan(ctx context.Context, spanName string, attrs ...attribute.KeyValue) (context.Context, trace.Span) {
	return tp.tracer.Start(ctx, spanName, trace.WithAttributes(attrs...))
}

// StartSpanWithKind начинает новый span с определённым типом
func (tp *TracingProvider) StartSpanWithKind(ctx context.Context, spanName string, kind trace.SpanKind, attrs ...attribute.KeyValue) (context.Context, trace.Span) {
	return tp.tracer.Start(ctx, spanName, trace.WithSpanKind(kind), trace.WithAttributes(attrs...))
}

// AddEvent добавляет событие к span'у
func (tp *TracingProvider) AddEvent(span trace.Span, eventName string, attrs ...attribute.KeyValue) {
	span.AddEvent(eventName, trace.WithAttributes(attrs...))
}

// SetAttribute устанавливает атрибут span'у
func (tp *TracingProvider) SetAttribute(span trace.Span, key attribute.Key, value interface{}) {
	span.SetAttributes(attribute.KeyValue{Key: key, Value: attribute.Value(attribute.StringValue(fmt.Sprintf("%v", value)))})
}

// EndSpan завершает span с опциональной ошибкой
func (tp *TracingProvider) EndSpan(span trace.Span, err error) {
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
	} else {
		span.SetStatus(codes.Ok, "")
	}
	span.End()
}

// OperationSpan вспомогательная структура для работы со span'ами в операциях
type OperationSpan struct {
	span trace.Span
	ctx  context.Context
	tp   *TracingProvider
}

// NewOperationSpan создаёт новый OperationSpan для упрощения работы
func (tp *TracingProvider) NewOperationSpan(ctx context.Context, operationType, documentID string, attrs ...attribute.KeyValue) *OperationSpan {
	// Добавляем стандартные атрибуты
	defaultAttrs := []attribute.KeyValue{
		attribute.String("operation_type", operationType),
		attribute.String("document_id", documentID),
	}
	allAttrs := append(defaultAttrs, attrs...)

	newCtx, span := tp.StartSpan(ctx, fmt.Sprintf("document.%s", operationType), allAttrs...)

	return &OperationSpan{
		span: span,
		ctx:  newCtx,
		tp:   tp,
	}
}

// Context возвращает контекст с span'ом
func (os *OperationSpan) Context() context.Context {
	return os.ctx
}

// End завершает operation span с опциональной ошибкой
func (os *OperationSpan) End(err error) {
	os.tp.EndSpan(os.span, err)
}

// AddEvent добавляет событие к operation span'у
func (os *OperationSpan) AddEvent(eventName string, attrs ...attribute.KeyValue) {
	os.tp.AddEvent(os.span, eventName, attrs...)
}

// SetAttribute устанавливает атрибут operation span'у
func (os *OperationSpan) SetAttribute(key attribute.Key, value interface{}) {
	os.span.SetAttributes(attribute.KeyValue{Key: key, Value: attribute.Value(attribute.StringValue(fmt.Sprintf("%v", value)))})
}

// RecordError записывает ошибку в span
func (os *OperationSpan) RecordError(err error, attrs ...attribute.KeyValue) {
	os.span.RecordError(err, trace.WithAttributes(attrs...))
	os.span.SetStatus(codes.Error, err.Error())
}

// SetStatus устанавливает статус span'а
func (os *OperationSpan) SetStatus(statusCode codes.Code, description string) {
	os.span.SetStatus(statusCode, description)
}

// Global tracer provider для использования во всём приложении
var globalTracingProvider *TracingProvider

// InitTracingProvider инициализирует глобальный провайдер трассирования
func InitTracingProvider(serviceName string) *TracingProvider {
	globalTracingProvider = NewTracingProvider(serviceName)
	return globalTracingProvider
}

// GetTracingProvider возвращает глобальный провайдер трассирования
func GetTracingProvider() *TracingProvider {
	if globalTracingProvider == nil {
		globalTracingProvider = NewTracingProvider("document-server")
	}
	return globalTracingProvider
}

// Common attributes for document operations
const (
	DocumentIDKey     = attribute.Key("document_id")
	OrganizationIDKey = attribute.Key("organization_id")
	OperationTypeKey  = attribute.Key("operation_type")
	StatusKey         = attribute.Key("status")
	UserIDKey         = attribute.Key("user_id")
	ErrorReasonKey    = attribute.Key("error_reason")
	VersionKey        = attribute.Key("version")
	PageKey           = attribute.Key("page")
	PerPageKey        = attribute.Key("per_page")
	ResultCountKey    = attribute.Key("result_count")
	DurationMsKey     = attribute.Key("duration_ms")
)

// Span names for document operations
const (
	SpanCreateDocument = "document.create"
	SpanGetDocument    = "document.get"
	SpanUpdateDocument = "document.update"
	SpanListDocuments  = "document.list"
	SpanSendDocument   = "document.send"
	SpanApproveDoc     = "document.approve"
	SpanRejectDoc      = "document.reject"
	SpanArchiveDoc     = "document.archive"
	SpanValidation     = "validation"
	SpanDatabaseOp     = "database"
	SpanPublishEvent   = "publish_event"
)

// Event names for document operations
const (
	EventDocumentCreated  = "document_created"
	EventDocumentUpdated  = "document_updated"
	EventDocumentSent     = "document_sent"
	EventDocumentApproved = "document_approved"
	EventDocumentRejected = "document_rejected"
	EventDocumentArchived = "document_archived"
	EventValidationPassed = "validation_passed"
	EventValidationFailed = "validation_failed"
	EventDatabaseRead     = "database_read"
	EventDatabaseWrite    = "database_write"
	EventErrorOccurred    = "error_occurred"
	EventCacheMiss        = "cache_miss"
	EventCacheHit         = "cache_hit"
	EventVersionConflict  = "version_conflict"
	EventConcurrentAccess = "concurrent_access"
)
