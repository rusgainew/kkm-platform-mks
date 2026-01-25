package services

import (
	"context"
	"fmt"

	"github.com/rusgainew/kkm-project-mks/api-gateway/internal/infrastructure/client"
	"github.com/rusgainew/kkm-project-mks/api-gateway/internal/infrastructure/observability"
	pb "github.com/rusgainew/kkm-project-mks/proto-lib/api"
	"go.uber.org/zap"
)

// CompanyQueryService сервис для чтения компаний (query side)
type CompanyQueryService struct {
	connManager *client.ConnectionManager
	serviceAddr string
	metrics     *observability.Metrics
	tracer      *observability.Tracer
	logger      *zap.Logger
}

// NewCompanyQueryService создает новый CompanyQueryService
func NewCompanyQueryService(
	connManager *client.ConnectionManager,
	serviceAddr string,
	metrics *observability.Metrics,
	tracer *observability.Tracer,
	logger *zap.Logger,
) *CompanyQueryService {
	return &CompanyQueryService{
		connManager: connManager,
		serviceAddr: serviceAddr,
		metrics:     metrics,
		tracer:      tracer,
		logger:      logger,
	}
}

// getClient возвращает gRPC клиент для company query service
func (s *CompanyQueryService) getClient(ctx context.Context) (pb.CatalogQueryServiceClient, error) {
	conn, err := s.connManager.GetConnection(ctx, s.serviceAddr)
	if err != nil {
		return nil, fmt.Errorf("failed to get connection: %w", err)
	}
	return pb.NewCatalogQueryServiceClient(conn), nil
}

// ListCompanies получает список всех компаний с пагинацией
func (s *CompanyQueryService) ListCompanies(ctx context.Context, pageNum, pageSize int32) (*pb.APIResponse, error) {
	ctx, endSpan := s.tracer.StartSpan(ctx, "CompanyQueryService.ListCompanies")
	defer endSpan()

	client, err := s.getClient(ctx)
	if err != nil {
		s.logger.Error("Failed to get company query client", zap.Error(err))
		return nil, err
	}

	req := &pb.PageInfo{
		Page: pageNum,
		Size: pageSize,
	}

	resp, err := client.ListCatalogs(ctx, req)
	if err != nil {
		s.logger.Error("Failed to list companies", zap.Error(err))
		s.metrics.IncrementErrorCount("company_query_service", "list_companies")
		return nil, fmt.Errorf("failed to list companies: %w", err)
	}

	return resp, nil
}

// SearchCompanies осуществляет поиск компаний по текстовому запросу
func (s *CompanyQueryService) SearchCompanies(ctx context.Context, searchReq *pb.SearchRequest) (*pb.APIResponse, error) {
	ctx, endSpan := s.tracer.StartSpan(ctx, "CompanyQueryService.SearchCompanies")
	defer endSpan()

	client, err := s.getClient(ctx)
	if err != nil {
		s.logger.Error("Failed to get company query client", zap.Error(err))
		return nil, err
	}

	// Преобразуем SearchRequest в CatalogFilterRequest для использования ListCatalogsWithFilter
	filterReq := &pb.CatalogFilterRequest{
		Page: searchReq.GetPage(),
		Size: searchReq.GetSize(),
		Name: searchReq.GetSearchText(), // Используем SearchText как Name фильтра
	}

	resp, err := client.ListCatalogsWithFilter(ctx, filterReq)
	if err != nil {
		s.logger.Error("Failed to search companies", zap.Error(err))
		s.metrics.IncrementErrorCount("company_query_service", "search_companies")
		return nil, fmt.Errorf("failed to search companies: %w", err)
	}

	return resp, nil
}

// FilterCompanies получает компании с фильтрацией
func (s *CompanyQueryService) FilterCompanies(ctx context.Context, filterReq *pb.CatalogFilterRequest) (*pb.APIResponse, error) {
	ctx, endSpan := s.tracer.StartSpan(ctx, "CompanyQueryService.FilterCompanies")
	defer endSpan()

	client, err := s.getClient(ctx)
	if err != nil {
		s.logger.Error("Failed to get company query client", zap.Error(err))
		return nil, err
	}

	resp, err := client.ListCatalogsWithFilter(ctx, filterReq)
	if err != nil {
		s.logger.Error("Failed to filter companies", zap.Error(err))
		s.metrics.IncrementErrorCount("company_query_service", "filter_companies")
		return nil, fmt.Errorf("failed to filter companies: %w", err)
	}

	return resp, nil
}
