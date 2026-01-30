// Файл api-gateway/internal/application/services/catalog_service.go содержит реализацию пакета services.
package services

import (
	"context"
	"fmt"

	"github.com/rusgainew/kkm-project-mks/api-gateway/internal/domain/models"
	"github.com/rusgainew/kkm-project-mks/api-gateway/internal/infrastructure/client"
	"github.com/rusgainew/kkm-project-mks/api-gateway/internal/infrastructure/observability"
	pb "github.com/rusgainew/kkm-project-mks/proto-lib/api"
	"go.uber.org/zap"
)

// CatalogService сервис для работы с каталогом
type CatalogService struct {
	connManager *client.ConnectionManager
	serviceAddr string
	metrics     *observability.Metrics
	tracer      *observability.Tracer
	logger      *zap.Logger
}

// NewCatalogService создает новый CatalogService
func NewCatalogService(
	connManager *client.ConnectionManager,
	serviceAddr string,
	metrics *observability.Metrics,
	tracer *observability.Tracer,
	logger *zap.Logger,
) *CatalogService {
	return &CatalogService{
		connManager: connManager,
		serviceAddr: serviceAddr,
		metrics:     metrics,
		tracer:      tracer,
		logger:      logger,
	}
}

// getClient возвращает gRPC клиент для catalog service
func (s *CatalogService) getClient(ctx context.Context) (pb.CatalogCommandServiceClient, error) {
	conn, err := s.connManager.GetConnection(ctx, s.serviceAddr)
	if err != nil {
		return nil, fmt.Errorf("failed to get connection: %w", err)
	}
	return pb.NewCatalogCommandServiceClient(conn), nil
}

// CreateCatalog создает новый элемент каталога
func (s *CatalogService) CreateCatalog(ctx context.Context, catalog *models.Catalog) (*models.Catalog, error) {
	ctx, endSpan := s.tracer.StartSpan(ctx, "CatalogService.CreateCatalog")
	defer endSpan()

	client, err := s.getClient(ctx)
	if err != nil {
		s.logger.Error("Failed to get catalog client", zap.Error(err))
		return nil, err
	}

	req := &pb.CreateCatalogRequest{
		Name:      catalog.Name,
		Number:    catalog.Number,
		TnvedCode: catalog.TnvedCode,
	}

	resp, err := client.CreateCatalog(ctx, req)
	if err != nil {
		s.logger.Error("Failed to create catalog", zap.Error(err))
		s.metrics.IncrementErrorCount("catalog_service", "create_catalog")
		return nil, fmt.Errorf("failed to create catalog: %w", err)
	}

	s.logger.Info("Catalog created successfully", zap.String("id", resp.Id))

	return &models.Catalog{
		ID:          resp.Id,
		Name:        resp.Name,
		Description: resp.Number,
		Category:    resp.TnvedCode,
		Price:       0,
		Currency:    "",
		Unit:        resp.UnitClassification.GetCode(),
		CreatedAt:   0,
		UpdatedAt:   0,
	}, nil
}

// UpdateCatalog обновляет элемент каталога
func (s *CatalogService) UpdateCatalog(ctx context.Context, catalog *models.Catalog) (*models.Catalog, error) {
	ctx, endSpan := s.tracer.StartSpan(ctx, "CatalogService.UpdateCatalog")
	defer endSpan()

	client, err := s.getClient(ctx)
	if err != nil {
		s.logger.Error("Failed to get catalog client", zap.Error(err))
		return nil, err
	}

	req := &pb.UpdateCatalogRequest{
		Id:        catalog.ID,
		Name:      catalog.Name,
		Number:    catalog.Number,
		TnvedCode: catalog.TnvedCode,
	}

	resp, err := client.UpdateCatalog(ctx, req)
	if err != nil {
		s.logger.Error("Failed to update catalog", zap.String("id", catalog.ID), zap.Error(err))
		s.metrics.IncrementErrorCount("catalog_service", "update_catalog")
		return nil, fmt.Errorf("failed to update catalog: %w", err)
	}

	s.logger.Info("Catalog updated successfully", zap.String("id", resp.Id))

	return &models.Catalog{
		ID:          resp.Id,
		Name:        resp.Name,
		Number:      resp.Number,
		Description: "", // TODO: добавить в proto если нужно
		TnvedCode:   resp.TnvedCode,
		Category:    resp.TnvedCode, // deprecated, дублируем для обратной совместимости
		Price:       0,
		Currency:    "",
		Unit:        resp.UnitClassification.GetCode(),
		CreatedAt:   0,
		UpdatedAt:   0,
	}, nil
}

// GetCatalog возвращает элемент каталога по ID
func (s *CatalogService) GetCatalog(ctx context.Context, id string) (*models.Catalog, error) {
	ctx, endSpan := s.tracer.StartSpan(ctx, "CatalogService.GetCatalog")
	defer endSpan()

	client, err := s.getClient(ctx)
	if err != nil {
		s.logger.Error("Failed to get catalog client", zap.Error(err))
		return nil, err
	}

	req := &pb.GetCatalogRequest{
		Id: id,
	}

	resp, err := client.GetCatalog(ctx, req)
	if err != nil {
		s.logger.Error("Failed to get catalog", zap.String("id", id), zap.Error(err))
		s.metrics.IncrementErrorCount("catalog_service", "get_catalog")
		return nil, fmt.Errorf("failed to get catalog: %w", err)
	}

	s.logger.Info("Catalog retrieved successfully", zap.String("id", resp.Id))

	return &models.Catalog{
		ID:          resp.Id,
		Name:        resp.Name,
		Description: resp.Number,
		Category:    resp.TnvedCode,
		Price:       0,
		Currency:    "",
		Unit:        resp.UnitClassification.GetCode(),
		CreatedAt:   0,
		UpdatedAt:   0,
	}, nil
}

// DeleteCatalog удаляет элемент каталога
func (s *CatalogService) DeleteCatalog(ctx context.Context, id string) error {
	ctx, endSpan := s.tracer.StartSpan(ctx, "CatalogService.DeleteCatalog")
	defer endSpan()

	client, err := s.getClient(ctx)
	if err != nil {
		s.logger.Error("Failed to get catalog client", zap.Error(err))
		return err
	}

	req := &pb.DeleteCatalogRequest{
		Id: id,
	}

	_, err = client.DeleteCatalog(ctx, req)
	if err != nil {
		s.logger.Error("Failed to delete catalog", zap.String("id", id), zap.Error(err))
		s.metrics.IncrementErrorCount("catalog_service", "delete_catalog")
		return fmt.Errorf("failed to delete catalog: %w", err)
	}

	s.logger.Info("Catalog deleted successfully", zap.String("id", id))

	return nil
}
