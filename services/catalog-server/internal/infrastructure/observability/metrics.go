package observability

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// CatalogItemsTotal количество созданных элементов каталога
	CatalogItemsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "catalog_items_total",
			Help: "Total number of catalog items created",
		},
		[]string{"organization_id"},
	)

	// CatalogItemsActive количество активных элементов
	CatalogItemsActive = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "catalog_items_active",
			Help: "Number of active catalog items",
		},
		[]string{"organization_id"},
	)

	// CatalogOperationDuration время выполнения операций
	CatalogOperationDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "catalog_operation_duration_seconds",
			Help:    "Duration of catalog operations",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"operation"},
	)

	// CatalogOperationErrors количество ошибок операций
	CatalogOperationErrors = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "catalog_operation_errors_total",
			Help: "Total number of catalog operation errors",
		},
		[]string{"operation", "error_type"},
	)

	// CatalogSearchDuration время выполнения поиска
	CatalogSearchDuration = promauto.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "catalog_search_duration_seconds",
			Help:    "Duration of catalog search operations",
			Buckets: prometheus.DefBuckets,
		},
	)
)
