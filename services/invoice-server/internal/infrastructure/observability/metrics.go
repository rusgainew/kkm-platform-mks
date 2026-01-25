package observability

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// Metrics содержит метрики Prometheus для invoice-server
type Metrics struct {
	InvoicesCreated      prometheus.Counter
	InvoicesSigned       prometheus.Counter
	InvoicesAccepted     prometheus.Counter
	InvoicesRejected     prometheus.Counter
	InvoicesRevoked      prometheus.Counter
	InvoicesTotal        prometheus.Gauge
	OperationDuration    *prometheus.HistogramVec
	OperationErrors      *prometheus.CounterVec
	InvoicesByStatus     *prometheus.GaugeVec
	DetailsPerInvoice    prometheus.Histogram
	FinancialDataRecords prometheus.Counter
}

// NewMetrics создает новый экземпляр метрик с глобальным реестром Prometheus
func NewMetrics() *Metrics {
	return NewMetricsWithRegisterer(prometheus.DefaultRegisterer)
}

// NewMetricsWithRegisterer создает новый экземпляр метрик с указанным реестром Prometheus (удобно для unit-тестов)
func NewMetricsWithRegisterer(registerer prometheus.Registerer) *Metrics {
	if registerer == nil {
		registerer = prometheus.DefaultRegisterer
	}

	factory := promauto.With(registerer)

	return &Metrics{
		InvoicesCreated: factory.NewCounter(prometheus.CounterOpts{
			Name: "invoice_server_invoices_created_total",
			Help: "Total number of invoices created",
		}),
		InvoicesSigned: factory.NewCounter(prometheus.CounterOpts{
			Name: "invoice_server_invoices_signed_total",
			Help: "Total number of invoices signed",
		}),
		InvoicesAccepted: factory.NewCounter(prometheus.CounterOpts{
			Name: "invoice_server_invoices_accepted_total",
			Help: "Total number of invoices accepted",
		}),
		InvoicesRejected: factory.NewCounter(prometheus.CounterOpts{
			Name: "invoice_server_invoices_rejected_total",
			Help: "Total number of invoices rejected",
		}),
		InvoicesRevoked: factory.NewCounter(prometheus.CounterOpts{
			Name: "invoice_server_invoices_revoked_total",
			Help: "Total number of invoices revoked",
		}),
		InvoicesTotal: factory.NewGauge(prometheus.GaugeOpts{
			Name: "invoice_server_invoices_total",
			Help: "Total number of invoices in the system",
		}),
		OperationDuration: factory.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "invoice_server_operation_duration_seconds",
				Help:    "Duration of invoice operations",
				Buckets: prometheus.DefBuckets,
			},
			[]string{"operation", "status"},
		),
		OperationErrors: factory.NewCounterVec(
			prometheus.CounterOpts{
				Name: "invoice_server_operation_errors_total",
				Help: "Total number of operation errors",
			},
			[]string{"operation", "error_type"},
		),
		InvoicesByStatus: factory.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "invoice_server_invoices_by_status",
				Help: "Number of invoices by status",
			},
			[]string{"status"},
		),
		DetailsPerInvoice: factory.NewHistogram(prometheus.HistogramOpts{
			Name:    "invoice_server_details_per_invoice",
			Help:    "Number of detail items per invoice",
			Buckets: []float64{1, 5, 10, 20, 50, 100},
		}),
		FinancialDataRecords: factory.NewCounter(prometheus.CounterOpts{
			Name: "invoice_server_financial_data_records_total",
			Help: "Total number of financial data records created",
		}),
	}
}
