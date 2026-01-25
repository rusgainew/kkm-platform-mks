package services

import (
	"context"
	"fmt"

	"github.com/rusgainew/kkm-project-mks/api-gateway/internal/infrastructure/client"
	"github.com/rusgainew/kkm-project-mks/api-gateway/internal/infrastructure/observability"
	pb "github.com/rusgainew/kkm-project-mks/proto-lib/api"
	"go.uber.org/zap"
)

// CatalogQueryService сервис для чтения каталога (query side)
type CatalogQueryService struct {
	connManager *client.ConnectionManager
	serviceAddr string
	metrics     *observability.Metrics
	tracer      *observability.Tracer
	logger      *zap.Logger
}

// NewCatalogQueryService создает новый CatalogQueryService
func NewCatalogQueryService(
	connManager *client.ConnectionManager,
	serviceAddr string,
	metrics *observability.Metrics,
	tracer *observability.Tracer,
	logger *zap.Logger,
) *CatalogQueryService {
	return &CatalogQueryService{
		connManager: connManager,
		serviceAddr: serviceAddr,
		metrics:     metrics,
		tracer:      tracer,
		logger:      logger,
	}
}

// getClient возвращает gRPC клиент для catalog query service
func (s *CatalogQueryService) getClient(ctx context.Context) (pb.CatalogQueryServiceClient, error) {
	conn, err := s.connManager.GetConnection(ctx, s.serviceAddr)
	if err != nil {
		return nil, fmt.Errorf("failed to get connection: %w", err)
	}
	return pb.NewCatalogQueryServiceClient(conn), nil
}

// ListCatalogs получает список всех элементов каталога с пагинацией
func (s *CatalogQueryService) ListCatalogs(ctx context.Context, pageNum, pageSize int32) (*pb.APIResponse, error) {
	ctx, endSpan := s.tracer.StartSpan(ctx, "CatalogQueryService.ListCatalogs")
	defer endSpan()

	client, err := s.getClient(ctx)
	if err != nil {
		s.logger.Error("Failed to get catalog query client", zap.Error(err))
		return nil, err
	}

	req := &pb.PageInfo{
		Page: pageNum,
		Size: pageSize,
	}

	resp, err := client.ListCatalogs(ctx, req)
	if err != nil {
		s.logger.Error("Failed to list catalogs", zap.Error(err))
		s.metrics.IncrementErrorCount("catalog_query_service", "list_catalogs")
		return nil, fmt.Errorf("failed to list catalogs: %w", err)
	}

	return resp, nil
}

// ListCatalogsWithFilter получает список каталогов с фильтрацией
func (s *CatalogQueryService) ListCatalogsWithFilter(ctx context.Context, filter *pb.CatalogFilterRequest) (*pb.APIResponse, error) {
	ctx, endSpan := s.tracer.StartSpan(ctx, "CatalogQueryService.ListCatalogsWithFilter")
	defer endSpan()

	client, err := s.getClient(ctx)
	if err != nil {
		s.logger.Error("Failed to get catalog query client", zap.Error(err))
		return nil, err
	}

	resp, err := client.ListCatalogsWithFilter(ctx, filter)
	if err != nil {
		s.logger.Error("Failed to list catalogs with filter", zap.Error(err))
		s.metrics.IncrementErrorCount("catalog_query_service", "list_catalogs_with_filter")
		return nil, fmt.Errorf("failed to list catalogs with filter: %w", err)
	}

	return resp, nil
}
