package repository

import (
	"context"
	"time"

	"github.com/rusgainew/kkm-project-mks/analytics-server/internal/domain/repository"
	"github.com/rusgainew/kkm-project-mks/analytics-server/internal/infrastructure/observability"
	"go.uber.org/zap"
)

// AnalyticsRepositoryWithMetrics оборачивает AnalyticsRepository для добавления метрик
type AnalyticsRepositoryWithMetrics struct {
	repo    repository.AnalyticsRepository
	metrics *observability.Metrics
	logger  *zap.Logger
}

// NewAnalyticsRepositoryWithMetrics создает новый wrapper с метриками
func NewAnalyticsRepositoryWithMetrics(
	repo repository.AnalyticsRepository,
	metrics *observability.Metrics,
	logger *zap.Logger,
) repository.AnalyticsRepository {
	return &AnalyticsRepositoryWithMetrics{
		repo:    repo,
		metrics: metrics,
		logger:  logger,
	}
}

// GetStats возвращает статистику с метриками
func (r *AnalyticsRepositoryWithMetrics) GetStats(ctx context.Context, startDate, endDate time.Time) (*repository.AnalyticsStats, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start).Seconds()
		r.metrics.RecordQueryLatency("get_stats", duration)
	}()

	stats, err := r.repo.GetStats(ctx, startDate, endDate)
	if err != nil {
		r.metrics.RecordAnalyticsError("get_stats", "database_error")
		return nil, err
	}

	return stats, nil
}

// GetSalesData возвращает данные продаж с метриками
func (r *AnalyticsRepositoryWithMetrics) GetSalesData(ctx context.Context, startDate, endDate time.Time, granularity string) ([]repository.SalesDataPoint, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start).Seconds()
		r.metrics.RecordQueryLatency("get_sales_data", duration)
	}()

	data, err := r.repo.GetSalesData(ctx, startDate, endDate, granularity)
	if err != nil {
		r.metrics.RecordAnalyticsError("get_sales_data", "database_error")
		return nil, err
	}

	return data, nil
}

// GetStatusDistribution возвращает распределение по статусам с метриками
func (r *AnalyticsRepositoryWithMetrics) GetStatusDistribution(ctx context.Context, startDate, endDate time.Time) ([]repository.StatusDistribution, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start).Seconds()
		r.metrics.RecordQueryLatency("get_status_distribution", duration)
	}()

	data, err := r.repo.GetStatusDistribution(ctx, startDate, endDate)
	if err != nil {
		r.metrics.RecordAnalyticsError("get_status_distribution", "database_error")
		return nil, err
	}

	return data, nil
}

// GetOperationTypeDistribution возвращает распределение по типам операций с метриками
func (r *AnalyticsRepositoryWithMetrics) GetOperationTypeDistribution(ctx context.Context, startDate, endDate time.Time) ([]repository.OperationTypeDistribution, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start).Seconds()
		r.metrics.RecordQueryLatency("get_operation_type_distribution", duration)
	}()

	data, err := r.repo.GetOperationTypeDistribution(ctx, startDate, endDate)
	if err != nil {
		r.metrics.RecordAnalyticsError("get_operation_type_distribution", "database_error")
		return nil, err
	}

	return data, nil
}

// GetTopContractors возвращает топ контрагентов с метриками
func (r *AnalyticsRepositoryWithMetrics) GetTopContractors(ctx context.Context, startDate, endDate time.Time, limit int32) ([]repository.TopContractorData, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start).Seconds()
		r.metrics.RecordQueryLatency("get_top_contractors", duration)
	}()

	data, err := r.repo.GetTopContractors(ctx, startDate, endDate, limit)
	if err != nil {
		r.metrics.RecordAnalyticsError("get_top_contractors", "database_error")
		return nil, err
	}

	return data, nil
}

// GetMonthlyRevenue возвращает месячную выручку с метриками
func (r *AnalyticsRepositoryWithMetrics) GetMonthlyRevenue(ctx context.Context, startDate, endDate time.Time) ([]repository.MonthlyRevenueData, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start).Seconds()
		r.metrics.RecordQueryLatency("get_monthly_revenue", duration)
	}()

	data, err := r.repo.GetMonthlyRevenue(ctx, startDate, endDate)
	if err != nil {
		r.metrics.RecordAnalyticsError("get_monthly_revenue", "database_error")
		return nil, err
	}

	return data, nil
}
