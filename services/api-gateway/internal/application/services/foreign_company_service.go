// Файл api-gateway/internal/application/services/foreign_company_service.go содержит реализацию пакета services.
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

// ForeignCompanyService сервис для работы с иностранными компаниями
type ForeignCompanyService struct {
	connManager *client.ConnectionManager
	serviceAddr string
	metrics     *observability.Metrics
	tracer      *observability.Tracer
	logger      *zap.Logger
}

// NewForeignCompanyService создает новый ForeignCompanyService
func NewForeignCompanyService(
	connManager *client.ConnectionManager,
	serviceAddr string,
	metrics *observability.Metrics,
	tracer *observability.Tracer,
	logger *zap.Logger,
) *ForeignCompanyService {
	return &ForeignCompanyService{
		connManager: connManager,
		serviceAddr: serviceAddr,
		metrics:     metrics,
		tracer:      tracer,
		logger:      logger,
	}
}

// getClient возвращает gRPC клиент для foreign company service
func (s *ForeignCompanyService) getClient(ctx context.Context) (pb.ForeignCompanyCommandServiceClient, error) {
	conn, err := s.connManager.GetConnection(ctx, s.serviceAddr)
	if err != nil {
		return nil, fmt.Errorf("failed to get connection: %w", err)
	}
	return pb.NewForeignCompanyCommandServiceClient(conn), nil
}

// CreateForeignCompany создает новую иностранную компанию
func (s *ForeignCompanyService) CreateForeignCompany(ctx context.Context, pin, fullName, countryCode, address string) (*models.ForeignCompany, error) {
	ctx, endSpan := s.tracer.StartSpan(ctx, "ForeignCompanyService.CreateForeignCompany")
	defer endSpan()

	client, err := s.getClient(ctx)
	if err != nil {
		s.logger.Error("Failed to get foreign company client", zap.Error(err))
		return nil, err
	}

	req := &pb.CreateForeignCompanyRequest{
		Pin:         pin,
		FullName:    fullName,
		CountryCode: countryCode,
		Address:     address,
	}

	resp, err := client.CreateForeignCompany(ctx, req)
	if err != nil {
		s.logger.Error("Failed to create foreign company", zap.Error(err))
		s.metrics.IncrementErrorCount("foreign_company_service", "create")
		return nil, fmt.Errorf("failed to create foreign company: %w", err)
	}

	s.logger.Info("Foreign company created successfully", zap.Int64("id", resp.Id))

	return &models.ForeignCompany{
		ID:          resp.Id,
		PIN:         resp.Pin,
		FullName:    resp.FullName,
		CountryCode: countryCode,
		Address:     address,
	}, nil
}

// UpdateForeignCompany обновляет иностранную компанию
func (s *ForeignCompanyService) UpdateForeignCompany(ctx context.Context, id int64, pin, fullName, countryCode, address string) (*models.ForeignCompany, error) {
	ctx, endSpan := s.tracer.StartSpan(ctx, "ForeignCompanyService.UpdateForeignCompany")
	defer endSpan()

	client, err := s.getClient(ctx)
	if err != nil {
		s.logger.Error("Failed to get foreign company client", zap.Error(err))
		return nil, err
	}

	req := &pb.UpdateForeignCompanyRequest{
		Id:          id,
		Pin:         pin,
		FullName:    fullName,
		CountryCode: countryCode,
		Address:     address,
	}

	resp, err := client.UpdateForeignCompany(ctx, req)
	if err != nil {
		s.logger.Error("Failed to update foreign company", zap.Error(err))
		s.metrics.IncrementErrorCount("foreign_company_service", "update")
		return nil, fmt.Errorf("failed to update foreign company: %w", err)
	}

	s.logger.Info("Foreign company updated successfully", zap.Int64("id", resp.Id))

	return &models.ForeignCompany{
		ID:          resp.Id,
		PIN:         resp.Pin,
		FullName:    resp.FullName,
		CountryCode: countryCode,
		Address:     address,
	}, nil
}

// GetForeignCompany возвращает иностранную компанию по ID
func (s *ForeignCompanyService) GetForeignCompany(ctx context.Context, id int64) (*models.ForeignCompany, error) {
	ctx, endSpan := s.tracer.StartSpan(ctx, "ForeignCompanyService.GetForeignCompany")
	defer endSpan()

	client, err := s.getClient(ctx)
	if err != nil {
		s.logger.Error("Failed to get foreign company client", zap.Error(err))
		return nil, err
	}

	req := &pb.GetForeignCompanyRequest{
		Id: id,
	}

	resp, err := client.GetForeignCompany(ctx, req)
	if err != nil {
		s.logger.Error("Failed to get foreign company", zap.Int64("id", id), zap.Error(err))
		s.metrics.IncrementErrorCount("foreign_company_service", "get")
		return nil, fmt.Errorf("failed to get foreign company: %w", err)
	}

	s.logger.Info("Foreign company retrieved successfully", zap.Int64("id", resp.Id))

	return &models.ForeignCompany{
		ID:       resp.Id,
		PIN:      resp.Pin,
		FullName: resp.FullName,
	}, nil
}

// DeleteForeignCompany удаляет иностранную компанию
func (s *ForeignCompanyService) DeleteForeignCompany(ctx context.Context, id int64) error {
	ctx, endSpan := s.tracer.StartSpan(ctx, "ForeignCompanyService.DeleteForeignCompany")
	defer endSpan()

	client, err := s.getClient(ctx)
	if err != nil {
		s.logger.Error("Failed to get foreign company client", zap.Error(err))
		return err
	}

	req := &pb.DeleteForeignCompanyRequest{
		Id: id,
	}

	_, err = client.DeleteForeignCompany(ctx, req)
	if err != nil {
		s.logger.Error("Failed to delete foreign company", zap.Int64("id", id), zap.Error(err))
		s.metrics.IncrementErrorCount("foreign_company_service", "delete")
		return fmt.Errorf("failed to delete foreign company: %w", err)
	}

	s.logger.Info("Foreign company deleted successfully", zap.Int64("id", id))

	return nil
}
