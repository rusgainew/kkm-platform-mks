package observability

import (
	"github.com/prometheus/client_golang/prometheus"
)

// MetricsCollector holds all Prometheus metrics for invoice-query-server
type MetricsCollector struct {
	RequestsTotal        prometheus.CounterVec
	RequestDuration      prometheus.HistogramVec
	ErrorsTotal          prometheus.CounterVec
	QueryExecutionTime   prometheus.HistogramVec
	DBConnectionPoolSize prometheus.GaugeVec
}

// NewMetricsCollector creates and registers all metrics
func NewMetricsCollector(namespace, subsystem string) *MetricsCollector {
	mc := &MetricsCollector{
		// Total requests counter
		RequestsTotal: *prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Subsystem: subsystem,
				Name:      "requests_total",
				Help:      "Total number of gRPC requests received",
			},
			[]string{"method", "status"},
		),

		// Request duration histogram (in milliseconds)
		RequestDuration: *prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Namespace: namespace,
				Subsystem: subsystem,
				Name:      "request_duration_ms",
				Help:      "gRPC request duration in milliseconds",
				Buckets:   prometheus.ExponentialBuckets(1, 2, 10), // 1, 2, 4, 8, 16, 32, 64, 128, 256, 512 ms
			},
			[]string{"method"},
		),

		// Total errors counter
		ErrorsTotal: *prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Subsystem: subsystem,
				Name:      "errors_total",
				Help:      "Total number of errors",
			},
			[]string{"method", "error_type"},
		),

		// Database query execution time (in milliseconds)
		QueryExecutionTime: *prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Namespace: namespace,
				Subsystem: subsystem,
				Name:      "query_execution_time_ms",
				Help:      "Database query execution time in milliseconds",
				Buckets:   prometheus.ExponentialBuckets(0.5, 2, 10), // 0.5, 1, 2, 4, 8, 16, 32, 64, 128, 256 ms
			},
			[]string{"query_type"},
		),

		// Database connection pool size
		DBConnectionPoolSize: *prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: subsystem,
				Name:      "db_connection_pool_size",
				Help:      "Current number of connections in the database connection pool",
			},
			[]string{"state"},
		),
	}

	// Register all metrics
	prometheus.MustRegister(&mc.RequestsTotal)
	prometheus.MustRegister(&mc.RequestDuration)
	prometheus.MustRegister(&mc.ErrorsTotal)
	prometheus.MustRegister(&mc.QueryExecutionTime)
	prometheus.MustRegister(&mc.DBConnectionPoolSize)

	return mc
}
