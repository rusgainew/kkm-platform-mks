// Файл foreign-company-server/internal/infrastructure/observability/metrics.go содержит реализацию пакета observability.
package observability

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// ForeignCompanyOperationsTotal общее количество операций с иностранными компаниями
	ForeignCompanyOperationsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "foreign_company_operations_total",
			Help: "Total number of foreign company operations",
		},
		[]string{"operation", "status"},
	)

	// ForeignCompanyOperationDuration время выполнения операций
	ForeignCompanyOperationDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "foreign_company_operation_duration_seconds",
			Help:    "Duration of foreign company operations in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"operation"},
	)

	// ForeignCompanyErrors количество ошибок
	ForeignCompanyErrors = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "foreign_company_errors_total",
			Help: "Total number of foreign company errors",
		},
		[]string{"operation", "error_type"},
	)

	// ForeignCompaniesActive количество активных иностранных компаний
	ForeignCompaniesActive = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "foreign_companies_active",
			Help: "Number of active foreign companies",
		},
	)

	// ForeignCompaniesTotal общее количество иностранных компаний
	ForeignCompaniesTotal = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "foreign_companies_total",
			Help: "Total number of foreign companies in database",
		},
	)
)

// RecordOperation записывает метрику операции
func RecordOperation(operation, status string) {
	ForeignCompanyOperationsTotal.WithLabelValues(operation, status).Inc()
}

// RecordError записывает метрику ошибки
func RecordError(operation, errorType string) {
	ForeignCompanyErrors.WithLabelValues(operation, errorType).Inc()
}

// ObserveDuration записывает время выполнения операции
func ObserveDuration(operation string, duration float64) {
	ForeignCompanyOperationDuration.WithLabelValues(operation).Observe(duration)
}

// SetActiveCount устанавливает количество активных компаний
func SetActiveCount(count float64) {
	ForeignCompaniesActive.Set(count)
}

// SetTotalCount устанавливает общее количество компаний
func SetTotalCount(count float64) {
	ForeignCompaniesTotal.Set(count)
}
