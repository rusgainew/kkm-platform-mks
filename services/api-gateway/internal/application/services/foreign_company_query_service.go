// Файл api-gateway/internal/application/services/foreign_company_query_service.go содержит реализацию пакета services.
package services

import (
	"context"
	"fmt"

	"github.com/rusgainew/kkm-project-mks/api-gateway/internal/infrastructure/client"
	"github.com/rusgainew/kkm-project-mks/api-gateway/internal/infrastructure/observability"
	pb "github.com/rusgainew/kkm-project-mks/proto-lib/api"
	"go.uber.org/zap"
)

// ForeignCompanyQueryService сервис для чтения иностранных компаний (query side)
type ForeignCompanyQueryService struct {
	connManager *client.ConnectionManager
	serviceAddr string
	metrics     *observability.Metrics
	tracer      *observability.Tracer
	logger      *zap.Logger
}

// NewForeignCompanyQueryService создает новый ForeignCompanyQueryService
func NewForeignCompanyQueryService(
	connManager *client.ConnectionManager,
	serviceAddr string,
	metrics *observability.Metrics,
	tracer *observability.Tracer,
	logger *zap.Logger,
) *ForeignCompanyQueryService {
	return &ForeignCompanyQueryService{
		connManager: connManager,
		serviceAddr: serviceAddr,
		metrics:     metrics,
		tracer:      tracer,
		logger:      logger,
	}
}

// getClient возвращает gRPC клиент для foreign company query service
func (s *ForeignCompanyQueryService) getClient(ctx context.Context) (pb.CatalogQueryServiceClient, error) {
	conn, err := s.connManager.GetConnection(ctx, s.serviceAddr)
	if err != nil {
		return nil, fmt.Errorf("failed to get connection: %w", err)
	}
	return pb.NewCatalogQueryServiceClient(conn), nil
}

// ListForeignCompanies получает список всех иностранных компаний с пагинацией
func (s *ForeignCompanyQueryService) ListForeignCompanies(ctx context.Context, pageNum, pageSize int32) (*pb.APIResponse, error) {
	ctx, endSpan := s.tracer.StartSpan(ctx, "ForeignCompanyQueryService.ListForeignCompanies")
	defer endSpan()

	client, err := s.getClient(ctx)
	if err != nil {
		s.logger.Error("Failed to get foreign company query client", zap.Error(err))
		return nil, err
	}

	req := &pb.PageInfo{
		Page: pageNum,
		Size: pageSize,
	}

	resp, err := client.ListCatalogs(ctx, req)
	if err != nil {
		s.logger.Error("Failed to list foreign companies", zap.Error(err))
		s.metrics.IncrementErrorCount("foreign_company_query_service", "list_foreign_companies")
		return nil, fmt.Errorf("failed to list foreign companies: %w", err)
	}

	return resp, nil
}

// SearchForeignCompanies осуществляет поиск иностранных компаний по текстовому запросу
func (s *ForeignCompanyQueryService) SearchForeignCompanies(ctx context.Context, searchReq *pb.SearchRequest) (*pb.APIResponse, error) {
	ctx, endSpan := s.tracer.StartSpan(ctx, "ForeignCompanyQueryService.SearchForeignCompanies")
	defer endSpan()

	client, err := s.getClient(ctx)
	if err != nil {
		s.logger.Error("Failed to get foreign company query client", zap.Error(err))
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
		s.logger.Error("Failed to search foreign companies", zap.Error(err))
		s.metrics.IncrementErrorCount("foreign_company_query_service", "search_foreign_companies")
		return nil, fmt.Errorf("failed to search foreign companies: %w", err)
	}

	return resp, nil
}

// FilterForeignCompanies получает иностранные компании с фильтрацией
func (s *ForeignCompanyQueryService) FilterForeignCompanies(ctx context.Context, filterReq *pb.CatalogFilterRequest) (*pb.APIResponse, error) {
	ctx, endSpan := s.tracer.StartSpan(ctx, "ForeignCompanyQueryService.FilterForeignCompanies")
	defer endSpan()

	client, err := s.getClient(ctx)
	if err != nil {
		s.logger.Error("Failed to get foreign company query client", zap.Error(err))
		return nil, err
	}

	resp, err := client.ListCatalogsWithFilter(ctx, filterReq)
	if err != nil {
		s.logger.Error("Failed to filter foreign companies", zap.Error(err))
		s.metrics.IncrementErrorCount("foreign_company_query_service", "filter_foreign_companies")
		return nil, fmt.Errorf("failed to filter foreign companies: %w", err)
	}

	return resp, nil
}
