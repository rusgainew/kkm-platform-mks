package services

import (
	"context"
	"fmt"

	"github.com/rusgainew/kkm-project-mks/api-gateway/internal/domain/models"
	"github.com/rusgainew/kkm-project-mks/api-gateway/internal/infrastructure/cache"
	"github.com/rusgainew/kkm-project-mks/api-gateway/internal/infrastructure/client"
	"github.com/rusgainew/kkm-project-mks/api-gateway/internal/infrastructure/observability"
	"github.com/rusgainew/kkm-project-mks/lib/conversion"
	pb "github.com/rusgainew/kkm-project-mks/proto-lib/company"
	"go.uber.org/zap"
)

// CompanyService сервис для работы с компаниями
type CompanyService struct {
	connManager *client.ConnectionManager
	serviceAddr string
	metrics     *observability.Metrics
	tracer      *observability.Tracer
	logger      *zap.Logger
	cache       *cache.RedisCache
}

// NewCompanyService создает новый CompanyService
func NewCompanyService(
	connManager *client.ConnectionManager,
	serviceAddr string,
	metrics *observability.Metrics,
	tracer *observability.Tracer,
	logger *zap.Logger,
	cache *cache.RedisCache,
) *CompanyService {
	return &CompanyService{
		connManager: connManager,
		serviceAddr: serviceAddr,
		metrics:     metrics,
		tracer:      tracer,
		logger:      logger,
		cache:       cache,
	}
}

// getClient возвращает gRPC клиент для company service
func (s *CompanyService) getClient(ctx context.Context) (pb.CompanyServiceClient, error) {
	conn, err := s.connManager.GetConnection(ctx, s.serviceAddr)
	if err != nil {
		return nil, fmt.Errorf("failed to get connection: %w", err)
	}
	return pb.NewCompanyServiceClient(conn), nil
}

// CreateCompany создает новую компанию
func (s *CompanyService) CreateCompany(ctx context.Context, company *models.Company) (*models.Company, error) {
	ctx, endSpan := s.tracer.StartSpan(ctx, "CompanyService.CreateCompany")
	defer endSpan()

	client, err := s.getClient(ctx)
	if err != nil {
		s.logger.Error("Failed to get company client", zap.Error(err))
		return nil, err
	}

	req := &pb.CreateOrganizationRequest{
		Name:        company.Name,
		Description: company.Description,
		OwnerId:     company.OwnerID,
	}

	resp, err := client.CreateOrganization(ctx, req)
	if err != nil {
		s.logger.Error("Failed to create company", zap.Error(err))
		s.metrics.IncrementErrorCount("company_service", "create_company")
		return nil, fmt.Errorf("failed to create company: %w", err)
	}

	s.logger.Info("Company created successfully", zap.String("id", resp.Id))

	return &models.Company{
		ID:          resp.Id,
		Name:        resp.Name,
		Description: resp.Description,
		OwnerID:     resp.OwnerId,
		MemberCount: int(resp.MemberCount),
		CreatedAt:   resp.CreatedAt,
		UpdatedAt:   resp.UpdatedAt,
		Status:      resp.Status,
	}, nil
}

// GetCompany получает компанию по ID
func (s *CompanyService) GetCompany(ctx context.Context, id string) (*models.Company, error) {
	ctx, endSpan := s.tracer.StartSpan(ctx, "CompanyService.GetCompany")
	defer endSpan()

	// Проверяем кэш перед запросом к backend
	if s.cache != nil && s.cache.IsAvailable(ctx) {
		var cachedCompany models.Company
		cacheKey := fmt.Sprintf("company:%s", id)
		if err := s.cache.Get(ctx, cacheKey, &cachedCompany); err == nil {
			s.logger.Debug("Company retrieved from cache", zap.String("id", id))
			return &cachedCompany, nil
		}
		// Если ошибка кэша - продолжаем, не падаем
		s.logger.Debug("Cache miss for company", zap.String("id", id))
	}

	client, err := s.getClient(ctx)
	if err != nil {
		s.logger.Error("Failed to get company client", zap.Error(err))
		return nil, err
	}

	req := &pb.GetOrganizationRequest{
		OrganizationId: id,
	}

	resp, err := client.GetOrganization(ctx, req)
	if err != nil {
		s.logger.Error("Failed to get company", zap.String("id", id), zap.Error(err))
		s.metrics.IncrementErrorCount("company_service", "get_company")
		return nil, fmt.Errorf("failed to get company: %w", err)
	}

	s.logger.Info("Company retrieved successfully", zap.String("id", resp.Id))

	company := &models.Company{
		ID:          resp.Id,
		Name:        resp.Name,
		Description: resp.Description,
		OwnerID:     resp.OwnerId,
		MemberCount: int(resp.MemberCount),
		CreatedAt:   resp.CreatedAt,
		UpdatedAt:   resp.UpdatedAt,
		Status:      resp.Status,
	}

	// Сохраняем в кэш
	if s.cache != nil && s.cache.IsAvailable(ctx) {
		cacheKey := fmt.Sprintf("company:%s", id)
		if err := s.cache.Set(ctx, cacheKey, company); err != nil {
			s.logger.Warn("Failed to cache company", zap.String("id", id), zap.Error(err))
		}
	}

	return company, nil
}

// UpdateCompany обновляет компанию
func (s *CompanyService) UpdateCompany(ctx context.Context, company *models.Company) (*models.Company, error) {
	ctx, endSpan := s.tracer.StartSpan(ctx, "CompanyService.UpdateCompany")
	defer endSpan()

	client, err := s.getClient(ctx)
	if err != nil {
		s.logger.Error("Failed to get company client", zap.Error(err))
		return nil, err
	}

	req := &pb.UpdateOrganizationRequest{
		Id:          company.ID,
		Name:        company.Name,
		Description: company.Description,
	}

	resp, err := client.UpdateOrganization(ctx, req)
	if err != nil {
		s.logger.Error("Failed to update company", zap.String("id", company.ID), zap.Error(err))
		s.metrics.IncrementErrorCount("company_service", "update_company")
		return nil, fmt.Errorf("failed to update company: %w", err)
	}

	s.logger.Info("Company updated successfully", zap.String("id", resp.Id))

	updatedCompany := &models.Company{
		ID:          resp.Id,
		Name:        resp.Name,
		Description: resp.Description,
		OwnerID:     resp.OwnerId,
		MemberCount: int(resp.MemberCount),
		CreatedAt:   resp.CreatedAt,
		UpdatedAt:   resp.UpdatedAt,
		Status:      resp.Status,
	}

	// Инвалидируем кэш
	if s.cache != nil {
		cacheKey := fmt.Sprintf("company:%s", company.ID)
		if err := s.cache.Delete(ctx, cacheKey); err != nil {
			s.logger.Warn("Failed to invalidate cache", zap.String("id", company.ID), zap.Error(err))
		}
	}

	return updatedCompany, nil
}

// DeleteCompany удаляет компанию
func (s *CompanyService) DeleteCompany(ctx context.Context, id string) error {
	ctx, endSpan := s.tracer.StartSpan(ctx, "CompanyService.DeleteCompany")
	defer endSpan()

	client, err := s.getClient(ctx)
	if err != nil {
		s.logger.Error("Failed to get company client", zap.Error(err))
		return err
	}

	req := &pb.DeleteOrganizationRequest{
		OrganizationId: id,
	}

	_, err = client.DeleteOrganization(ctx, req)
	if err != nil {
		s.logger.Error("Failed to delete company", zap.String("id", id), zap.Error(err))
		s.metrics.IncrementErrorCount("company_service", "delete_company")
		return fmt.Errorf("failed to delete company: %w", err)
	}

	s.logger.Info("Company deleted successfully", zap.String("id", id))

	// Инвалидируем кэш
	if s.cache != nil {
		cacheKey := fmt.Sprintf("company:%s", id)
		if err := s.cache.Delete(ctx, cacheKey); err != nil {
			s.logger.Warn("Failed to invalidate cache", zap.String("id", id), zap.Error(err))
		}
	}

	return nil
}

// ListCompanies получает список компаний
func (s *CompanyService) ListCompanies(ctx context.Context, page, pageSize int) ([]*models.Company, int, error) {
	ctx, endSpan := s.tracer.StartSpan(ctx, "CompanyService.ListCompanies")
	defer endSpan()

	client, err := s.getClient(ctx)
	if err != nil {
		s.logger.Error("Failed to get company client", zap.Error(err))
		return nil, 0, err
	}

	req := &pb.ListOrganizationsRequest{
		Page:    conversion.SafeIntToInt32WithDefault(page, 0),
		PerPage: conversion.SafeIntToInt32WithDefault(pageSize, 20),
	}

	resp, err := client.ListOrganizations(ctx, req)
	if err != nil {
		s.logger.Error("Failed to list companies", zap.Error(err))
		s.metrics.IncrementErrorCount("company_service", "list_companies")
		return nil, 0, fmt.Errorf("failed to list companies: %w", err)
	}

	s.logger.Info("Companies retrieved successfully", zap.Int("count", len(resp.Organizations)))

	companies := make([]*models.Company, 0, len(resp.Organizations))
	for _, org := range resp.Organizations {
		companies = append(companies, &models.Company{
			ID:          org.Id,
			Name:        org.Name,
			Description: org.Description,
			OwnerID:     org.OwnerId,
			MemberCount: int(org.MemberCount),
			CreatedAt:   org.CreatedAt,
			UpdatedAt:   org.UpdatedAt,
			Status:      org.Status,
		})
	}

	return companies, int(resp.PageInfo.Total), nil
}

// GetOrganizationMembers получает список членов организации
func (s *CompanyService) GetOrganizationMembers(ctx context.Context, orgID string, page, pageSize int32) ([]*models.Employee, int32, error) {
	ctx, endSpan := s.tracer.StartSpan(ctx, "CompanyService.GetOrganizationMembers")
	defer endSpan()

	client, err := s.getClient(ctx)
	if err != nil {
		s.logger.Error("Failed to get company client", zap.Error(err))
		return nil, 0, err
	}

	req := &pb.GetOrganizationMembersRequest{
		OrganizationId: orgID,
		Page:           page,
		PerPage:        pageSize,
	}

	resp, err := client.GetOrganizationMembers(ctx, req)
	if err != nil {
		s.logger.Error("Failed to get organization members", zap.String("org_id", orgID), zap.Error(err))
		s.metrics.IncrementErrorCount("company_service", "get_members")
		return nil, 0, fmt.Errorf("failed to get organization members: %w", err)
	}

	s.logger.Info("Organization members retrieved", zap.String("org_id", orgID), zap.Int("count", len(resp.Employees)))

	members := make([]*models.Employee, 0, len(resp.Employees))
	for _, emp := range resp.Employees {
		members = append(members, &models.Employee{
			ID:             emp.Id,
			UserID:         emp.UserId,
			OrganizationID: emp.OrganizationId,
			Role:           emp.Role,
			Status:         emp.Status,
			Position:       emp.Position,
			Department:     emp.Department,
			JoinedAt:       emp.JoinedAt,
			LastActiveAt:   emp.LastActiveAt,
		})
	}

	return members, conversion.SafeInt64ToInt32WithDefault(resp.PageInfo.Total, 0), nil
}

// AddMember добавляет члена в организацию
func (s *CompanyService) AddMember(ctx context.Context, req *models.AddMemberRequest) (*models.Employee, error) {
	ctx, endSpan := s.tracer.StartSpan(ctx, "CompanyService.AddMember")
	defer endSpan()

	client, err := s.getClient(ctx)
	if err != nil {
		s.logger.Error("Failed to get company client", zap.Error(err))
		return nil, err
	}

	grpcReq := &pb.AddMemberRequest{
		OrganizationId: req.OrganizationID,
		UserId:         req.UserID,
		Role:           req.Role,
	}

	resp, err := client.AddMember(ctx, grpcReq)
	if err != nil {
		s.logger.Error("Failed to add member",
			zap.String("org_id", req.OrganizationID),
			zap.String("user_id", req.UserID),
			zap.Error(err))
		s.metrics.IncrementErrorCount("company_service", "add_member")
		return nil, fmt.Errorf("failed to add member: %w", err)
	}

	s.logger.Info("Member added successfully",
		zap.String("org_id", req.OrganizationID),
		zap.String("user_id", req.UserID))

	return &models.Employee{
		ID:             resp.Id,
		UserID:         resp.UserId,
		OrganizationID: resp.OrganizationId,
		Role:           resp.Role,
		Status:         resp.Status,
		Position:       resp.Position,
		Department:     resp.Department,
		JoinedAt:       resp.JoinedAt,
		LastActiveAt:   resp.LastActiveAt,
	}, nil
}

// RemoveMember удаляет члена из организации
func (s *CompanyService) RemoveMember(ctx context.Context, orgID, memberID string) error {
	ctx, endSpan := s.tracer.StartSpan(ctx, "CompanyService.RemoveMember")
	defer endSpan()

	client, err := s.getClient(ctx)
	if err != nil {
		s.logger.Error("Failed to get company client", zap.Error(err))
		return err
	}

	req := &pb.RemoveMemberRequest{
		OrganizationId: orgID,
		EmployeeId:     memberID,
	}

	_, err = client.RemoveMember(ctx, req)
	if err != nil {
		s.logger.Error("Failed to remove member",
			zap.String("org_id", orgID),
			zap.String("member_id", memberID),
			zap.Error(err))
		s.metrics.IncrementErrorCount("company_service", "remove_member")
		return fmt.Errorf("failed to remove member: %w", err)
	}

	s.logger.Info("Member removed successfully",
		zap.String("org_id", orgID),
		zap.String("member_id", memberID))

	return nil
}
