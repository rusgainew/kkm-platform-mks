// Файл analytics-server/internal/infrastructure/repository/postgres_analytics_repository.go содержит реализацию пакета repository.
package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/rusgainew/kkm-project-mks/analytics-server/internal/domain/repository"
	"go.uber.org/zap"
)

// PostgresAnalyticsRepository реализация AnalyticsRepository для PostgreSQL
type PostgresAnalyticsRepository struct {
	db     *sql.DB
	logger *zap.Logger
}

// NewPostgresAnalyticsRepository создает новый PostgresAnalyticsRepository
func NewPostgresAnalyticsRepository(db *sql.DB, logger *zap.Logger) repository.AnalyticsRepository {
	return &PostgresAnalyticsRepository{
		db:     db,
		logger: logger,
	}
}

// GetStats возвращает агрегированную статистику за период
func (r *PostgresAnalyticsRepository) GetStats(ctx context.Context, startDate, endDate time.Time) (*repository.AnalyticsStats, error) {
	query := `
		SELECT 
			COALESCE(SUM(total_amount), 0) as total_revenue,
			COUNT(*) as total_invoices,
			COALESCE(AVG(total_amount), 0) as average_amount,
			COUNT(DISTINCT contractor_id) as active_contractors,
			COUNT(CASE WHEN status = 'draft' OR status = 'sent' THEN 1 END) as pending_invoices,
			COUNT(CASE WHEN status = 'signed' OR status = 'accepted' THEN 1 END) as approved_invoices,
			COUNT(CASE WHEN status = 'rejected' OR status = 'revoked' THEN 1 END) as rejected_invoices
		FROM invoices
		WHERE created_date >= $1 AND created_date <= $2
	`

	stats := &repository.AnalyticsStats{}
	err := r.db.QueryRowContext(ctx, query, startDate, endDate).Scan(
		&stats.TotalRevenue,
		&stats.TotalInvoices,
		&stats.AverageAmount,
		&stats.ActiveContractors,
		&stats.PendingInvoices,
		&stats.ApprovedInvoices,
		&stats.RejectedInvoices,
	)
	if err != nil {
		r.logger.Error("Failed to get analytics stats", zap.Error(err))
		return nil, fmt.Errorf("get analytics stats: %w", err)
	}

	r.logger.Info("Retrieved analytics stats",
		zap.Time("startDate", startDate),
		zap.Time("endDate", endDate),
		zap.Float64("totalRevenue", stats.TotalRevenue),
		zap.Int32("totalInvoices", stats.TotalInvoices),
	)

	return stats, nil
}

// GetSalesData возвращает данные продаж по дням/неделям/месяцам
func (r *PostgresAnalyticsRepository) GetSalesData(ctx context.Context, startDate, endDate time.Time, granularity string) ([]repository.SalesDataPoint, error) {
	var query string

	switch granularity {
	case "day":
		query = `
			SELECT 
				DATE_TRUNC('day', created_date) as date,
				COALESCE(SUM(total_amount), 0) as amount
			FROM invoices
			WHERE created_date >= $1 AND created_date <= $2
			GROUP BY DATE_TRUNC('day', created_date)
			ORDER BY date ASC
		`
	case "week":
		query = `
			SELECT 
				DATE_TRUNC('week', created_date) as date,
				COALESCE(SUM(total_amount), 0) as amount
			FROM invoices
			WHERE created_date >= $1 AND created_date <= $2
			GROUP BY DATE_TRUNC('week', created_date)
			ORDER BY date ASC
		`
	case "month":
		query = `
			SELECT 
				DATE_TRUNC('month', created_date) as date,
				COALESCE(SUM(total_amount), 0) as amount
			FROM invoices
			WHERE created_date >= $1 AND created_date <= $2
			GROUP BY DATE_TRUNC('month', created_date)
			ORDER BY date ASC
		`
	default:
		return nil, fmt.Errorf("invalid granularity: %s", granularity)
	}

	rows, err := r.db.QueryContext(ctx, query, startDate, endDate)
	if err != nil {
		r.logger.Error("Failed to get sales data", zap.Error(err))
		return nil, fmt.Errorf("get sales data: %w", err)
	}
	defer rows.Close()

	var salesData []repository.SalesDataPoint
	for rows.Next() {
		var dataPoint repository.SalesDataPoint
		err := rows.Scan(&dataPoint.Date, &dataPoint.Amount)
		if err != nil {
			r.logger.Error("Failed to scan sales data point", zap.Error(err))
			return nil, fmt.Errorf("scan sales data: %w", err)
		}
		salesData = append(salesData, dataPoint)
	}

	if err := rows.Err(); err != nil {
		r.logger.Error("Error iterating sales data rows", zap.Error(err))
		return nil, fmt.Errorf("iterate sales data: %w", err)
	}

	r.logger.Info("Retrieved sales data",
		zap.Int("count", len(salesData)),
		zap.String("granularity", granularity),
	)

	return salesData, nil
}

// GetStatusDistribution возвращает распределение счетов по статусам
func (r *PostgresAnalyticsRepository) GetStatusDistribution(ctx context.Context, startDate, endDate time.Time) ([]repository.StatusDistribution, error) {
	query := `
		SELECT 
			status,
			COUNT(*) as count,
			COALESCE(SUM(total_amount), 0) as amount
		FROM invoices
		WHERE created_date >= $1 AND created_date <= $2
		GROUP BY status
		ORDER BY count DESC
	`

	rows, err := r.db.QueryContext(ctx, query, startDate, endDate)
	if err != nil {
		r.logger.Error("Failed to get status distribution", zap.Error(err))
		return nil, fmt.Errorf("get status distribution: %w", err)
	}
	defer rows.Close()

	var distribution []repository.StatusDistribution
	for rows.Next() {
		var dist repository.StatusDistribution
		err := rows.Scan(&dist.Status, &dist.Count, &dist.Amount)
		if err != nil {
			r.logger.Error("Failed to scan status distribution", zap.Error(err))
			return nil, fmt.Errorf("scan status distribution: %w", err)
		}
		distribution = append(distribution, dist)
	}

	if err := rows.Err(); err != nil {
		r.logger.Error("Error iterating status distribution rows", zap.Error(err))
		return nil, fmt.Errorf("iterate status distribution: %w", err)
	}

	return distribution, nil
}

// GetOperationTypeDistribution возвращает распределение по типам операций
func (r *PostgresAnalyticsRepository) GetOperationTypeDistribution(ctx context.Context, startDate, endDate time.Time) ([]repository.OperationTypeDistribution, error) {
	// Примечание: В текущей схеме нет поля operation_type
	// Используем is_resident для определения типа операции (импорт/экспорт)
	query := `
		SELECT 
			CASE 
				WHEN is_resident THEN 'local'
				ELSE 'import'
			END as operation_type,
			COUNT(*) as count,
			COALESCE(SUM(total_amount), 0) as amount
		FROM invoices
		WHERE created_date >= $1 AND created_date <= $2
		GROUP BY is_resident
		ORDER BY count DESC
	`

	rows, err := r.db.QueryContext(ctx, query, startDate, endDate)
	if err != nil {
		r.logger.Error("Failed to get operation type distribution", zap.Error(err))
		return nil, fmt.Errorf("get operation type distribution: %w", err)
	}
	defer rows.Close()

	var distribution []repository.OperationTypeDistribution
	for rows.Next() {
		var dist repository.OperationTypeDistribution
		err := rows.Scan(&dist.OperationType, &dist.Count, &dist.Amount)
		if err != nil {
			r.logger.Error("Failed to scan operation type distribution", zap.Error(err))
			return nil, fmt.Errorf("scan operation type distribution: %w", err)
		}
		distribution = append(distribution, dist)
	}

	if err := rows.Err(); err != nil {
		r.logger.Error("Error iterating operation type distribution rows", zap.Error(err))
		return nil, fmt.Errorf("iterate operation type distribution: %w", err)
	}

	return distribution, nil
}

// GetTopContractors возвращает топ N контрагентов по выручке
func (r *PostgresAnalyticsRepository) GetTopContractors(ctx context.Context, startDate, endDate time.Time, limit int32) ([]repository.TopContractorData, error) {
	query := `
		SELECT 
			contractor_id,
			COALESCE(contractor_id, 'Unknown') as contractor_name,
			COALESCE(SUM(total_amount), 0) as total_amount,
			COUNT(*) as invoice_count
		FROM invoices
		WHERE created_date >= $1 AND created_date <= $2
		GROUP BY contractor_id
		ORDER BY total_amount DESC
		LIMIT $3
	`

	rows, err := r.db.QueryContext(ctx, query, startDate, endDate, limit)
	if err != nil {
		r.logger.Error("Failed to get top contractors", zap.Error(err))
		return nil, fmt.Errorf("get top contractors: %w", err)
	}
	defer rows.Close()

	var contractors []repository.TopContractorData
	for rows.Next() {
		var contractor repository.TopContractorData
		err := rows.Scan(
			&contractor.ContractorID,
			&contractor.ContractorName,
			&contractor.TotalAmount,
			&contractor.InvoiceCount,
		)
		if err != nil {
			r.logger.Error("Failed to scan contractor data", zap.Error(err))
			return nil, fmt.Errorf("scan contractor data: %w", err)
		}
		contractors = append(contractors, contractor)
	}

	if err := rows.Err(); err != nil {
		r.logger.Error("Error iterating contractor rows", zap.Error(err))
		return nil, fmt.Errorf("iterate contractors: %w", err)
	}

	return contractors, nil
}

// GetMonthlyRevenue возвращает месячную выручку
func (r *PostgresAnalyticsRepository) GetMonthlyRevenue(ctx context.Context, startDate, endDate time.Time) ([]repository.MonthlyRevenueData, error) {
	query := `
		SELECT 
			DATE_TRUNC('month', created_date) as month,
			COALESCE(SUM(total_amount), 0) as amount
		FROM invoices
		WHERE created_date >= $1 AND created_date <= $2
		GROUP BY DATE_TRUNC('month', created_date)
		ORDER BY month ASC
	`

	rows, err := r.db.QueryContext(ctx, query, startDate, endDate)
	if err != nil {
		r.logger.Error("Failed to get monthly revenue", zap.Error(err))
		return nil, fmt.Errorf("get monthly revenue: %w", err)
	}
	defer rows.Close()

	var monthlyData []repository.MonthlyRevenueData
	for rows.Next() {
		var data repository.MonthlyRevenueData
		err := rows.Scan(&data.Month, &data.Amount)
		if err != nil {
			r.logger.Error("Failed to scan monthly revenue data", zap.Error(err))
			return nil, fmt.Errorf("scan monthly revenue: %w", err)
		}
		monthlyData = append(monthlyData, data)
	}

	if err := rows.Err(); err != nil {
		r.logger.Error("Error iterating monthly revenue rows", zap.Error(err))
		return nil, fmt.Errorf("iterate monthly revenue: %w", err)
	}

	return monthlyData, nil
}
