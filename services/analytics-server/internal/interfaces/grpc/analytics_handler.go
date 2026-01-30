// Файл analytics-server/internal/interfaces/grpc/analytics_handler.go содержит реализацию пакета grpc.
package grpc

import (
	"context"

	"github.com/rusgainew/kkm-project-mks/analytics-server/internal/application/services"
	pb "github.com/rusgainew/kkm-project-mks/proto-lib/analytics"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// AnalyticsHandler gRPC handler для аналитики
type AnalyticsHandler struct {
	pb.UnimplementedAnalyticsServiceServer
	service *services.AnalyticsService
	logger  *zap.Logger
}

// NewAnalyticsHandler создает новый gRPC handler
func NewAnalyticsHandler(service *services.AnalyticsService, logger *zap.Logger) *AnalyticsHandler {
	return &AnalyticsHandler{
		service: service,
		logger:  logger,
	}
}

// GetDashboardStats возвращает статистику для дашборда
func (h *AnalyticsHandler) GetDashboardStats(ctx context.Context, req *pb.StatsRequest) (*pb.StatsResponse, error) {
	h.logger.Debug("GetDashboardStats called",
		zap.Time("start_date", req.StartDate.AsTime()),
		zap.Time("end_date", req.EndDate.AsTime()),
	)

	stats, err := h.service.GetStats(ctx, req.StartDate.AsTime(), req.EndDate.AsTime())
	if err != nil {
		h.logger.Error("Failed to get dashboard stats", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "failed to get stats: %v", err)
	}

	return &pb.StatsResponse{
		TotalInvoices:     int64(stats.TotalInvoices),
		TotalRevenue:      stats.TotalRevenue,
		AverageAmount:     stats.AverageAmount,
		UniqueContractors: int64(stats.ActiveContractors),
		PendingCount:      int64(stats.PendingInvoices),
		ApprovedCount:     int64(stats.ApprovedInvoices),
		RejectedCount:     int64(stats.RejectedInvoices),
	}, nil
}

// GetSalesChart возвращает данные для графика продаж
func (h *AnalyticsHandler) GetSalesChart(ctx context.Context, req *pb.SalesChartRequest) (*pb.SalesChartResponse, error) {
	h.logger.Debug("GetSalesChart called",
		zap.Time("start_date", req.StartDate.AsTime()),
		zap.Time("end_date", req.EndDate.AsTime()),
		zap.String("granularity", req.Granularity),
	)

	data, err := h.service.GetSalesData(ctx, req.StartDate.AsTime(), req.EndDate.AsTime(), req.Granularity)
	if err != nil {
		h.logger.Error("Failed to get sales chart data", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "failed to get sales data: %v", err)
	}

	dataPoints := make([]*pb.SalesDataPoint, len(data))
	for i, point := range data {
		dataPoints[i] = &pb.SalesDataPoint{
			Date:   point.Date.Format("2006-01-02"),
			Count:  int64(1), // В базовой структуре нет Count, используем 1
			Amount: point.Amount,
		}
	}

	return &pb.SalesChartResponse{
		DataPoints: dataPoints,
	}, nil
}

// GetStatusDistribution возвращает распределение по статусам
func (h *AnalyticsHandler) GetStatusDistribution(ctx context.Context, req *pb.StatusRequest) (*pb.StatusResponse, error) {
	h.logger.Debug("GetStatusDistribution called",
		zap.Time("start_date", req.StartDate.AsTime()),
		zap.Time("end_date", req.EndDate.AsTime()),
	)

	data, err := h.service.GetStatusDistribution(ctx, req.StartDate.AsTime(), req.EndDate.AsTime())
	if err != nil {
		h.logger.Error("Failed to get status distribution", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "failed to get status distribution: %v", err)
	}

	items := make([]*pb.StatusDistributionItem, len(data))
	for i, item := range data {
		// Вычисляем процент если сумма > 0
		percentage := float64(0)
		if item.Amount > 0 {
			totalAmount := float64(0)
			for _, d := range data {
				totalAmount += d.Amount
			}
			if totalAmount > 0 {
				percentage = (item.Amount / totalAmount) * 100
			}
		}

		items[i] = &pb.StatusDistributionItem{
			Status:     item.Status,
			Count:      int64(item.Count),
			Percentage: percentage,
		}
	}

	return &pb.StatusResponse{
		Items: items,
	}, nil
}

// GetOperationTypeDistribution возвращает распределение по типам операций
func (h *AnalyticsHandler) GetOperationTypeDistribution(ctx context.Context, req *pb.OperationTypeRequest) (*pb.OperationTypeResponse, error) {
	h.logger.Debug("GetOperationTypeDistribution called",
		zap.Time("start_date", req.StartDate.AsTime()),
		zap.Time("end_date", req.EndDate.AsTime()),
	)

	data, err := h.service.GetOperationTypeDistribution(ctx, req.StartDate.AsTime(), req.EndDate.AsTime())
	if err != nil {
		h.logger.Error("Failed to get operation type distribution", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "failed to get operation type distribution: %v", err)
	}

	items := make([]*pb.OperationTypeItem, len(data))
	for i, item := range data {
		// Вычисляем процент если сумма > 0
		percentage := float64(0)
		if item.Amount > 0 {
			totalAmount := float64(0)
			for _, d := range data {
				totalAmount += d.Amount
			}
			if totalAmount > 0 {
				percentage = (item.Amount / totalAmount) * 100
			}
		}

		items[i] = &pb.OperationTypeItem{
			OperationType: item.OperationType,
			Count:         int64(item.Count),
			Amount:        item.Amount,
			Percentage:    percentage,
		}
	}

	return &pb.OperationTypeResponse{
		Items: items,
	}, nil
}

// GetTopContractors возвращает топ контрагентов
func (h *AnalyticsHandler) GetTopContractors(ctx context.Context, req *pb.TopContractorsRequest) (*pb.TopContractorsResponse, error) {
	h.logger.Debug("GetTopContractors called",
		zap.Time("start_date", req.StartDate.AsTime()),
		zap.Time("end_date", req.EndDate.AsTime()),
		zap.Int32("limit", req.Limit),
	)

	data, err := h.service.GetTopContractors(ctx, req.StartDate.AsTime(), req.EndDate.AsTime(), req.Limit)
	if err != nil {
		h.logger.Error("Failed to get top contractors", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "failed to get top contractors: %v", err)
	}

	contractors := make([]*pb.TopContractorItem, len(data))
	for i, contractor := range data {
		contractors[i] = &pb.TopContractorItem{
			ContractorId:   contractor.ContractorID,
			ContractorName: contractor.ContractorName,
			InvoiceCount:   int64(contractor.InvoiceCount),
			TotalAmount:    contractor.TotalAmount,
		}
	}

	return &pb.TopContractorsResponse{
		Contractors: contractors,
	}, nil
}

// GetMonthlyRevenue возвращает месячную выручку
func (h *AnalyticsHandler) GetMonthlyRevenue(ctx context.Context, req *pb.MonthlyRevenueRequest) (*pb.MonthlyRevenueResponse, error) {
	h.logger.Debug("GetMonthlyRevenue called",
		zap.Time("start_date", req.StartDate.AsTime()),
		zap.Time("end_date", req.EndDate.AsTime()),
	)

	data, err := h.service.GetMonthlyRevenue(ctx, req.StartDate.AsTime(), req.EndDate.AsTime())
	if err != nil {
		h.logger.Error("Failed to get monthly revenue", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "failed to get monthly revenue: %v", err)
	}

	months := make([]*pb.MonthlyRevenueItem, len(data))
	for i, month := range data {
		months[i] = &pb.MonthlyRevenueItem{
			Month:        month.Month.Format("2006-01"),
			Revenue:      month.Amount,
			InvoiceCount: int64(1), // В базовой структуре нет InvoiceCount
		}
	}

	return &pb.MonthlyRevenueResponse{
		Months: months,
	}, nil
}

// RegisterAnalyticsServer регистрирует Analytics gRPC сервис
func RegisterAnalyticsServer(grpcServer *grpc.Server, service *services.AnalyticsService, logger *zap.Logger) {
	handler := NewAnalyticsHandler(service, logger)
	pb.RegisterAnalyticsServiceServer(grpcServer, handler)
	logger.Info("Analytics gRPC service registered successfully")
}
