// Файл company-server/internal/infrastructure/observability/metrics.go содержит реализацию пакета observability.
package observability

import (
	"context"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"go.uber.org/zap"
)

var (
	// gRPC метрики
	GrpcRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "grpc_requests_total",
			Help: "Total number of gRPC requests",
		},
		[]string{"method", "status"},
	)

	GrpcRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "grpc_request_duration_seconds",
			Help:    "Duration of gRPC requests in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method"},
	)

	// Метрики организаций
	OrganizationsTotal = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "organizations_total",
			Help: "Total number of organizations",
		},
	)

	OrganizationOperations = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "organization_operations_total",
			Help: "Total number of organization operations",
		},
		[]string{"operation", "status"},
	)

	// Метрики участников
	EmployeesTotal = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "employees_total",
			Help: "Total number of employees",
		},
	)

	EmployeeOperations = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "employee_operations_total",
			Help: "Total number of employee operations",
		},
		[]string{"operation", "status"},
	)

	// Метрики событий
	EventsPublished = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "events_published_total",
			Help: "Total number of published events",
		},
		[]string{"event_type", "status"},
	)

	// Метрики базы данных
	DatabaseOperations = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "database_operations_total",
			Help: "Total number of database operations",
		},
		[]string{"operation", "table", "status"},
	)

	DatabaseOperationDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "database_operation_duration_seconds",
			Help:    "Duration of database operations in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"operation", "table"},
	)
)

// RecordOperationDuration записывает длительность операции
func RecordOperationDuration(operation string, duration time.Duration, status string) {
	seconds := duration.Seconds()
	OrganizationOperations.WithLabelValues(operation, status).Inc()

	// Также записываем в общий counter для gRPC методов
	GrpcRequestsTotal.WithLabelValues(operation, status).Inc()
	GrpcRequestDuration.WithLabelValues(operation).Observe(seconds)
}

// RecordOrganizationOperation записывает операцию над организацией
func RecordOrganizationOperation(operation, status string) {
	OrganizationOperations.WithLabelValues(operation, status).Inc()
	GrpcRequestsTotal.WithLabelValues(operation, status).Inc()
}

// RecordEmployeeOperation записывает операцию над сотрудником
func RecordEmployeeOperation(operation, status string) {
	EmployeeOperations.WithLabelValues(operation, status).Inc()
	GrpcRequestsTotal.WithLabelValues(operation, status).Inc()
}

// RecordEventPublished записывает опубликованное событие
func RecordEventPublished(eventType, status string) {
	EventsPublished.WithLabelValues(eventType, status).Inc()
}

// RecordDatabaseOperation записывает операцию БД
func RecordDatabaseOperation(operation, table, status string, duration time.Duration) {
	DatabaseOperations.WithLabelValues(operation, table, status).Inc()
	DatabaseOperationDuration.WithLabelValues(operation, table).Observe(duration.Seconds())
}

// UpdateOrganizationCount обновляет счетчик организаций
func UpdateOrganizationCount(ctx context.Context, count int64, logger *zap.Logger) {
	OrganizationsTotal.Set(float64(count))
	logger.Debug("Organization count updated", zap.Int64("count", count))
}

// UpdateEmployeeCount обновляет счетчик сотрудников
func UpdateEmployeeCount(ctx context.Context, count int64, logger *zap.Logger) {
	EmployeesTotal.Set(float64(count))
	logger.Debug("Employee count updated", zap.Int64("count", count))
}
