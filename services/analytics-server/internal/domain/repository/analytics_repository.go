// Файл analytics-server/internal/domain/repository/analytics_repository.go содержит реализацию пакета repository.
package repository

import (
	"context"
	"time"
)

// AnalyticsStats агрегированная статистика за период
type AnalyticsStats struct {
	TotalRevenue      float64
	TotalInvoices     int32
	AverageAmount     float64
	ActiveContractors int32
	PendingInvoices   int32
	ApprovedInvoices  int32
	RejectedInvoices  int32
}

// SalesDataPoint точка данных для графика продаж
type SalesDataPoint struct {
	Date   time.Time
	Amount float64
}

// StatusDistribution распределение по статусам
type StatusDistribution struct {
	Status string
	Count  int32
	Amount float64
}

// OperationTypeDistribution распределение по типам операций
type OperationTypeDistribution struct {
	OperationType string
	Count         int32
	Amount        float64
}

// TopContractorData данные топ контрагента
type TopContractorData struct {
	ContractorID   string
	ContractorName string
	TotalAmount    float64
	InvoiceCount   int32
}

// MonthlyRevenueData месячная выручка
type MonthlyRevenueData struct {
	Month  time.Time
	Amount float64
}

// AnalyticsRepository интерфейс для работы с аналитическими данными
type AnalyticsRepository interface {
	// GetStats возвращает агрегированную статистику за период
	GetStats(ctx context.Context, startDate, endDate time.Time) (*AnalyticsStats, error)

	// GetSalesData возвращает данные продаж по дням/неделям/месяцам
	GetSalesData(ctx context.Context, startDate, endDate time.Time, granularity string) ([]SalesDataPoint, error)

	// GetStatusDistribution возвращает распределение счетов по статусам
	GetStatusDistribution(ctx context.Context, startDate, endDate time.Time) ([]StatusDistribution, error)

	// GetOperationTypeDistribution возвращает распределение по типам операций
	GetOperationTypeDistribution(ctx context.Context, startDate, endDate time.Time) ([]OperationTypeDistribution, error)

	// GetTopContractors возвращает топ N контрагентов по выручке
	GetTopContractors(ctx context.Context, startDate, endDate time.Time, limit int32) ([]TopContractorData, error)

	// GetMonthlyRevenue возвращает месячную выручку
	GetMonthlyRevenue(ctx context.Context, startDate, endDate time.Time) ([]MonthlyRevenueData, error)
}
