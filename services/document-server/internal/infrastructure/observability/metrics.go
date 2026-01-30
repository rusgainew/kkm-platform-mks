// Файл document-server/internal/infrastructure/observability/metrics.go содержит реализацию пакета observability.
package observability

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// Metrics хранит все Prometheus метрики
type Metrics struct {
	// Document operations counters
	DocumentsCreated  prometheus.Counter
	DocumentsUpdated  prometheus.Counter
	DocumentsApproved prometheus.Counter
	DocumentsRejected prometheus.Counter
	DocumentsArchived prometheus.Counter
	DocumentsSent     prometheus.Counter

	// Request metrics
	RequestsTotal   prometheus.Counter
	RequestDuration prometheus.Histogram
	RequestsActive  prometheus.Gauge

	// Error metrics
	RepositoryErrors      prometheus.Counter
	ValidationErrors      prometheus.Counter
	RabbitMQPublishErrors prometheus.Counter

	// Business metrics
	DocumentsPerStatus prometheus.GaugeVec
	PaginationRequests prometheus.Histogram
}

// InitMetrics инициализирует все метрики
func InitMetrics() *Metrics {
	return &Metrics{
		DocumentsCreated: promauto.NewCounter(prometheus.CounterOpts{
			Name: "document_server_documents_created_total",
			Help: "Total number of documents created",
		}),
		DocumentsUpdated: promauto.NewCounter(prometheus.CounterOpts{
			Name: "document_server_documents_updated_total",
			Help: "Total number of documents updated",
		}),
		DocumentsApproved: promauto.NewCounter(prometheus.CounterOpts{
			Name: "document_server_documents_approved_total",
			Help: "Total number of documents approved",
		}),
		DocumentsRejected: promauto.NewCounter(prometheus.CounterOpts{
			Name: "document_server_documents_rejected_total",
			Help: "Total number of documents rejected",
		}),
		DocumentsArchived: promauto.NewCounter(prometheus.CounterOpts{
			Name: "document_server_documents_archived_total",
			Help: "Total number of documents archived",
		}),
		DocumentsSent: promauto.NewCounter(prometheus.CounterOpts{
			Name: "document_server_documents_sent_total",
			Help: "Total number of documents sent for approval",
		}),
		RequestsTotal: promauto.NewCounter(prometheus.CounterOpts{
			Name: "document_server_requests_total",
			Help: "Total number of requests",
		}),
		RequestDuration: promauto.NewHistogram(prometheus.HistogramOpts{
			Name:    "document_server_request_duration_seconds",
			Help:    "Request duration in seconds",
			Buckets: []float64{.001, .005, .01, .025, .05, .1, .25, .5, 1, 2, 5},
		}),
		RequestsActive: promauto.NewGauge(prometheus.GaugeOpts{
			Name: "document_server_requests_active",
			Help: "Number of active requests",
		}),
		RepositoryErrors: promauto.NewCounter(prometheus.CounterOpts{
			Name: "document_server_repository_errors_total",
			Help: "Total number of repository errors",
		}),
		ValidationErrors: promauto.NewCounter(prometheus.CounterOpts{
			Name: "document_server_validation_errors_total",
			Help: "Total number of validation errors",
		}),
		RabbitMQPublishErrors: promauto.NewCounter(prometheus.CounterOpts{
			Name: "document_server_rabbitmq_publish_errors_total",
			Help: "Total number of RabbitMQ publish errors",
		}),
		PaginationRequests: promauto.NewHistogram(prometheus.HistogramOpts{
			Name:    "document_server_pagination_page_size",
			Help:    "Page size in list requests",
			Buckets: []float64{5, 10, 20, 50, 100},
		}),
	}
}

// RecordDocumentCreated записывает метрику создания документа
func (m *Metrics) RecordDocumentCreated() {
	m.DocumentsCreated.Inc()
	m.RequestsTotal.Inc()
}

// RecordDocumentUpdated записывает метрику обновления документа
func (m *Metrics) RecordDocumentUpdated() {
	m.DocumentsUpdated.Inc()
	m.RequestsTotal.Inc()
}

// RecordDocumentApproved записывает метрику одобрения документа
func (m *Metrics) RecordDocumentApproved() {
	m.DocumentsApproved.Inc()
	m.RequestsTotal.Inc()
}

// RecordDocumentRejected записывает метрику отклонения документа
func (m *Metrics) RecordDocumentRejected() {
	m.DocumentsRejected.Inc()
	m.RequestsTotal.Inc()
}

// RecordDocumentArchived записывает метрику архивирования документа
func (m *Metrics) RecordDocumentArchived() {
	m.DocumentsArchived.Inc()
	m.RequestsTotal.Inc()
}

// RecordDocumentSent записывает метрику отправки документа
func (m *Metrics) RecordDocumentSent() {
	m.DocumentsSent.Inc()
	m.RequestsTotal.Inc()
}

// RecordRepositoryError записывает ошибку репозитория
func (m *Metrics) RecordRepositoryError() {
	m.RepositoryErrors.Inc()
}

// RecordValidationError записывает ошибку валидации
func (m *Metrics) RecordValidationError() {
	m.ValidationErrors.Inc()
}

// RecordRabbitMQPublishError записывает ошибку публикации в RabbitMQ
func (m *Metrics) RecordRabbitMQPublishError() {
	m.RabbitMQPublishErrors.Inc()
}

// RecordRequestDuration записывает длительность запроса
func (m *Metrics) RecordRequestDuration(duration time.Duration) {
	m.RequestDuration.Observe(duration.Seconds())
}

// RecordPaginationRequest записывает размер страницы в запросе
func (m *Metrics) RecordPaginationRequest(pageSize int) {
	m.PaginationRequests.Observe(float64(pageSize))
}

// RecordPaginationPageSize записывает размер страницы
func (m *Metrics) RecordPaginationPageSize(pageSize float64) {
	m.PaginationRequests.Observe(pageSize)
}

// IncActiveRequests увеличивает счетчик активных запросов
func (m *Metrics) IncActiveRequests() {
	m.RequestsActive.Inc()
}

// DecActiveRequests уменьшает счетчик активных запросов
func (m *Metrics) DecActiveRequests() {
	m.RequestsActive.Dec()
}
