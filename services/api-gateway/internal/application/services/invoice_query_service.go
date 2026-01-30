// Файл api-gateway/internal/application/services/invoice_query_service.go содержит реализацию пакета services.
package services

import (
	"context"
	"fmt"

	"github.com/rusgainew/kkm-project-mks/api-gateway/internal/infrastructure/client"
	"github.com/rusgainew/kkm-project-mks/api-gateway/internal/infrastructure/observability"
	pb "github.com/rusgainew/kkm-project-mks/proto-lib/api"
	"go.uber.org/zap"
)

// InvoiceQueryService сервис для чтения счетов (query side)
type InvoiceQueryService struct {
	connManager *client.ConnectionManager
	serviceAddr string
	metrics     *observability.Metrics
	tracer      *observability.Tracer
	logger      *zap.Logger
}

// NewInvoiceQueryService создает новый InvoiceQueryService
func NewInvoiceQueryService(
	connManager *client.ConnectionManager,
	serviceAddr string,
	metrics *observability.Metrics,
	tracer *observability.Tracer,
	logger *zap.Logger,
) *InvoiceQueryService {
	return &InvoiceQueryService{
		connManager: connManager,
		serviceAddr: serviceAddr,
		metrics:     metrics,
		tracer:      tracer,
		logger:      logger,
	}
}

// getClient возвращает gRPC клиент для invoice query service
func (s *InvoiceQueryService) getClient(ctx context.Context) (pb.InvoiceQueryServiceClient, error) {
	conn, err := s.connManager.GetConnection(ctx, s.serviceAddr)
	if err != nil {
		return nil, fmt.Errorf("failed to get connection: %w", err)
	}
	return pb.NewInvoiceQueryServiceClient(conn), nil
}

// ListInvoices получает список всех счетов с пагинацией
func (s *InvoiceQueryService) ListInvoices(ctx context.Context, pageNum, pageSize int32) (*pb.APIResponse, error) {
	ctx, endSpan := s.tracer.StartSpan(ctx, "InvoiceQueryService.ListInvoices")
	defer endSpan()

	client, err := s.getClient(ctx)
	if err != nil {
		s.logger.Error("Failed to get invoice query client", zap.Error(err))
		return nil, err
	}

	req := &pb.PageInfo{
		Page: pageNum,
		Size: pageSize,
	}

	resp, err := client.ListInvoices(ctx, req)
	if err != nil {
		s.logger.Error("Failed to list invoices", zap.Error(err))
		s.metrics.IncrementErrorCount("invoice_query_service", "list_invoices")
		return nil, fmt.Errorf("failed to list invoices: %w", err)
	}

	return resp, nil
}

// ListInvoicesWithFilter получает список счетов с фильтрацией
func (s *InvoiceQueryService) ListInvoicesWithFilter(ctx context.Context, filter *pb.InvoiceFilterRequest) (*pb.APIResponse, error) {
	ctx, endSpan := s.tracer.StartSpan(ctx, "InvoiceQueryService.ListInvoicesWithFilter")
	defer endSpan()

	client, err := s.getClient(ctx)
	if err != nil {
		s.logger.Error("Failed to get invoice query client", zap.Error(err))
		return nil, err
	}

	resp, err := client.ListInvoicesWithFilter(ctx, filter)
	if err != nil {
		s.logger.Error("Failed to list invoices with filter", zap.Error(err))
		s.metrics.IncrementErrorCount("invoice_query_service", "list_invoices_with_filter")
		return nil, fmt.Errorf("failed to list invoices with filter: %w", err)
	}

	return resp, nil
}

// SearchInvoices осуществляет поиск счетов по текстовому запросу
func (s *InvoiceQueryService) SearchInvoices(ctx context.Context, searchReq *pb.SearchRequest) (*pb.APIResponse, error) {
	ctx, endSpan := s.tracer.StartSpan(ctx, "InvoiceQueryService.SearchInvoices")
	defer endSpan()

	client, err := s.getClient(ctx)
	if err != nil {
		s.logger.Error("Failed to get invoice query client", zap.Error(err))
		return nil, err
	}

	resp, err := client.SearchInvoices(ctx, searchReq)
	if err != nil {
		s.logger.Error("Failed to search invoices", zap.Error(err))
		s.metrics.IncrementErrorCount("invoice_query_service", "search_invoices")
		return nil, fmt.Errorf("failed to search invoices: %w", err)
	}

	return resp, nil
}

// GetInvoiceByNumber получает счет по номеру
func (s *InvoiceQueryService) GetInvoiceByNumber(ctx context.Context, invoiceNum *pb.InvoiceNumberRequest) (*pb.APIResponse, error) {
	ctx, endSpan := s.tracer.StartSpan(ctx, "InvoiceQueryService.GetInvoiceByNumber")
	defer endSpan()

	client, err := s.getClient(ctx)
	if err != nil {
		s.logger.Error("Failed to get invoice query client", zap.Error(err))
		return nil, err
	}

	resp, err := client.GetInvoiceByNumber(ctx, invoiceNum)
	if err != nil {
		s.logger.Error("Failed to get invoice by number", zap.Error(err))
		s.metrics.IncrementErrorCount("invoice_query_service", "get_invoice_by_number")
		return nil, fmt.Errorf("failed to get invoice by number: %w", err)
	}

	return resp, nil
}

// ListInvoiceDetails получает список деталей счета
// TODO: требует реализации в proto и invoice-query-server
func (s *InvoiceQueryService) ListInvoiceDetails(ctx context.Context, invoiceUUID string, page, pageSize int32) (*pb.APIResponse, error) {
	// Временная заглушка - метод требует добавления в proto definitions
	return &pb.APIResponse{
		// Method not implemented
		RequestId: nil,
		Data:      nil,
	}, nil
}

// GetInvoicesByDateRange получает счета по диапазону дат
// TODO: требует реализации в proto и invoice-query-server
func (s *InvoiceQueryService) GetInvoicesByDateRange(ctx context.Context, startDate, endDate string, page, pageSize int32) (*pb.APIResponse, error) {
	// Временная заглушка - метод требует добавления в proto definitions
	return &pb.APIResponse{
		// Method not implemented
		RequestId: nil,
		Data:      nil,
	}, nil
}
