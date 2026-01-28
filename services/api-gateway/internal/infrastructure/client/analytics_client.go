package client

import (
	"context"
	"fmt"
	"time"

	pb "github.com/rusgainew/kkm-project-mks/proto-lib/analytics"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// AnalyticsClient клиент для взаимодействия с analytics-server
type AnalyticsClient struct {
	client pb.AnalyticsServiceClient
	conn   *grpc.ClientConn
	logger *zap.Logger
}

// NewAnalyticsClient создает новый клиент для analytics-server
func NewAnalyticsClient(address string, logger *zap.Logger) (*AnalyticsClient, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(ctx, address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to analytics-server at %s: %w", address, err)
	}

	client := pb.NewAnalyticsServiceClient(conn)
	logger.Info("Connected to analytics-server successfully", zap.String("address", address))

	return &AnalyticsClient{
		client: client,
		conn:   conn,
		logger: logger,
	}, nil
}

// GetDashboardStats получает статистику дашборда
func (c *AnalyticsClient) GetDashboardStats(ctx context.Context, startDate, endDate time.Time) (*pb.StatsResponse, error) {
	req := &pb.StatsRequest{
		StartDate: timestamppb.New(startDate),
		EndDate:   timestamppb.New(endDate),
	}

	resp, err := c.client.GetDashboardStats(ctx, req)
	if err != nil {
		c.logger.Error("Failed to get dashboard stats", zap.Error(err))
		return nil, err
	}

	return resp, nil
}

// GetSalesChart получает данные для графика продаж
func (c *AnalyticsClient) GetSalesChart(ctx context.Context, startDate, endDate time.Time, granularity string) (*pb.SalesChartResponse, error) {
	req := &pb.SalesChartRequest{
		StartDate:   timestamppb.New(startDate),
		EndDate:     timestamppb.New(endDate),
		Granularity: granularity,
	}

	resp, err := c.client.GetSalesChart(ctx, req)
	if err != nil {
		c.logger.Error("Failed to get sales chart", zap.Error(err))
		return nil, err
	}

	return resp, nil
}

// GetStatusDistribution получает распределение по статусам
func (c *AnalyticsClient) GetStatusDistribution(ctx context.Context, startDate, endDate time.Time) (*pb.StatusResponse, error) {
	req := &pb.StatusRequest{
		StartDate: timestamppb.New(startDate),
		EndDate:   timestamppb.New(endDate),
	}

	resp, err := c.client.GetStatusDistribution(ctx, req)
	if err != nil {
		c.logger.Error("Failed to get status distribution", zap.Error(err))
		return nil, err
	}

	return resp, nil
}

// GetOperationTypeDistribution получает распределение по типам операций
func (c *AnalyticsClient) GetOperationTypeDistribution(ctx context.Context, startDate, endDate time.Time) (*pb.OperationTypeResponse, error) {
	req := &pb.OperationTypeRequest{
		StartDate: timestamppb.New(startDate),
		EndDate:   timestamppb.New(endDate),
	}

	resp, err := c.client.GetOperationTypeDistribution(ctx, req)
	if err != nil {
		c.logger.Error("Failed to get operation type distribution", zap.Error(err))
		return nil, err
	}

	return resp, nil
}

// GetTopContractors получает топ контрагентов
func (c *AnalyticsClient) GetTopContractors(ctx context.Context, startDate, endDate time.Time, limit int32) (*pb.TopContractorsResponse, error) {
	req := &pb.TopContractorsRequest{
		StartDate: timestamppb.New(startDate),
		EndDate:   timestamppb.New(endDate),
		Limit:     limit,
	}

	resp, err := c.client.GetTopContractors(ctx, req)
	if err != nil {
		c.logger.Error("Failed to get top contractors", zap.Error(err))
		return nil, err
	}

	return resp, nil
}

// GetMonthlyRevenue получает месячную выручку
func (c *AnalyticsClient) GetMonthlyRevenue(ctx context.Context, startDate, endDate time.Time) (*pb.MonthlyRevenueResponse, error) {
	req := &pb.MonthlyRevenueRequest{
		StartDate: timestamppb.New(startDate),
		EndDate:   timestamppb.New(endDate),
	}

	resp, err := c.client.GetMonthlyRevenue(ctx, req)
	if err != nil {
		c.logger.Error("Failed to get monthly revenue", zap.Error(err))
		return nil, err
	}

	return resp, nil
}

// Close закрывает соединение
func (c *AnalyticsClient) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}
