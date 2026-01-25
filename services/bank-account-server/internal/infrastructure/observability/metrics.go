package observability

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// BankAccountsTotal количество созданных банковских счетов
	BankAccountsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "bank_accounts_total",
			Help: "Total number of bank accounts created",
		},
		[]string{"organization_id", "currency"},
	)

	// BankAccountsActive количество активных счетов
	BankAccountsActive = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "bank_accounts_active",
			Help: "Number of active bank accounts",
		},
		[]string{"organization_id", "currency"},
	)

	// BankAccountOperationDuration время выполнения операций
	BankAccountOperationDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "bank_account_operation_duration_seconds",
			Help:    "Duration of bank account operations",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"operation"},
	)

	// BankAccountOperationErrors количество ошибок операций
	BankAccountOperationErrors = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "bank_account_operation_errors_total",
			Help: "Total number of bank account operation errors",
		},
		[]string{"operation", "error_type"},
	)

	// DefaultAccountChanges количество изменений счета по умолчанию
	DefaultAccountChanges = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "bank_account_default_changes_total",
			Help: "Total number of default account changes",
		},
		[]string{"organization_id"},
	)
)
