// Файл company-server/internal/interfaces/grpc/company_handler.go содержит реализацию пакета grpc.
package grpc

import (
	"context"

	"github.com/rusgainew/kkm-project-mks/company-server/internal/application/company"
	"github.com/rusgainew/kkm-project-mks/company-server/internal/domain"
	"github.com/rusgainew/kkm-project-mks/proto-lib/common"
	pb "github.com/rusgainew/kkm-project-mks/proto-lib/company"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// CompanyHandler обрабатывает gRPC запросы для CompanyService
type CompanyHandler struct {
	pb.UnimplementedCompanyServiceServer
	service *company.Service
	logger  *zap.Logger
}

// NewCompanyHandler создает новый handler
func NewCompanyHandler(service *company.Service, logger *zap.Logger) *CompanyHandler {
	return &CompanyHandler{
		service: service,
		logger:  logger,
	}
}

// GetOrganization получает организацию по ID
func (h *CompanyHandler) GetOrganization(ctx context.Context, req *pb.GetOrganizationRequest) (*pb.Organization, error) {
	if err := ValidateGetOrganizationRequest(req.GetOrganizationId()); err != nil {
		h.logger.Warn("GetOrganization validation failed", zap.Error(err))
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	h.logger.Info("GetOrganization request", zap.String("id", req.GetOrganizationId()))

	org, err := h.service.GetOrganization(ctx, req.GetOrganizationId())
	if err != nil {
		return nil, mapError(err)
	}

	return organizationToProto(org), nil
}

// CreateOrganization создает новую организацию
func (h *CompanyHandler) CreateOrganization(ctx context.Context, req *pb.CreateOrganizationRequest) (*pb.Organization, error) {
	if err := ValidateCreateOrganizationRequest(req.GetName(), req.GetDescription(), req.GetOwnerId()); err != nil {
		h.logger.Warn("CreateOrganization validation failed", zap.Error(err))
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	// Извлекаем userID из контекста (установлен перехватчиком auth)
	userID, err := ExtractUserID(ctx)
	if err != nil {
		h.logger.Warn("CreateOrganization: failed to extract user id", zap.Error(err))
		return nil, status.Error(codes.Unauthenticated, "user authentication required")
	}

	// Проверяем что пользователь создает организацию только для себя
	if userID != req.GetOwnerId() {
		h.logger.Warn("CreateOrganization: user cannot create org for another user",
			zap.String("user_id", userID),
			zap.String("owner_id", req.GetOwnerId()))
		return nil, status.Error(codes.PermissionDenied, "user can only create organization for themselves")
	}

	h.logger.Info("CreateOrganization request", zap.String("name", req.GetName()), zap.String("user_id", userID))

	org, err := h.service.CreateOrganization(ctx, req.GetName(), req.GetDescription(), req.GetOwnerId())
	if err != nil {
		return nil, mapError(err)
	}

	return organizationToProto(org), nil
}

// UpdateOrganization обновляет организацию
func (h *CompanyHandler) UpdateOrganization(ctx context.Context, req *pb.UpdateOrganizationRequest) (*pb.Organization, error) {
	if err := ValidateUpdateOrganizationRequest(req.GetId(), req.GetName(), req.GetDescription()); err != nil {
		h.logger.Warn("UpdateOrganization validation failed", zap.Error(err))
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	// Извлекаем userID из контекста
	userID, err := ExtractUserID(ctx)
	if err != nil {
		h.logger.Warn("UpdateOrganization: failed to extract user id", zap.Error(err))
		return nil, status.Error(codes.Unauthenticated, "user authentication required")
	}

	// Проверяем право на обновление (должен быть владельцем)
	if err := h.service.CheckUserIsOwner(ctx, userID, req.GetId()); err != nil {
		h.logger.Warn("UpdateOrganization: unauthorized", zap.Error(err), zap.String("user_id", userID), zap.String("org_id", req.GetId()))
		return nil, status.Error(codes.PermissionDenied, "user is not the owner of this organization")
	}

	h.logger.Info("UpdateOrganization request", zap.String("id", req.GetId()), zap.String("user_id", userID))

	org, err := h.service.UpdateOrganization(ctx, req.GetId(), req.GetName(), req.GetDescription())
	if err != nil {
		return nil, mapError(err)
	}

	return organizationToProto(org), nil
}

// DeleteOrganization удаляет организацию
func (h *CompanyHandler) DeleteOrganization(ctx context.Context, req *pb.DeleteOrganizationRequest) (*common.Empty, error) {
	if err := ValidateDeleteOrganizationRequest(req.GetOrganizationId()); err != nil {
		h.logger.Warn("DeleteOrganization validation failed", zap.Error(err))
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	// Извлекаем userID из контекста
	userID, err := ExtractUserID(ctx)
	if err != nil {
		h.logger.Warn("DeleteOrganization: failed to extract user id", zap.Error(err))
		return nil, status.Error(codes.Unauthenticated, "user authentication required")
	}

	// Проверяем право на удаление (должен быть владельцем)
	if err := h.service.CheckUserIsOwner(ctx, userID, req.GetOrganizationId()); err != nil {
		h.logger.Warn("DeleteOrganization: unauthorized", zap.Error(err), zap.String("user_id", userID), zap.String("org_id", req.GetOrganizationId()))
		return nil, status.Error(codes.PermissionDenied, "user is not the owner of this organization")
	}

	h.logger.Info("DeleteOrganization request", zap.String("id", req.GetOrganizationId()), zap.String("user_id", userID))

	err = h.service.DeleteOrganization(ctx, req.GetOrganizationId())
	if err != nil {
		return nil, mapError(err)
	}

	return &common.Empty{}, nil
}

// ListOrganizations возвращает список организаций
func (h *CompanyHandler) ListOrganizations(ctx context.Context, req *pb.ListOrganizationsRequest) (*pb.ListOrganizationsResponse, error) {
	if err := ValidateListOrganizationsRequest(req.GetPage(), req.GetPerPage(), req.GetOwnerId()); err != nil {
		h.logger.Warn("ListOrganizations validation failed", zap.Error(err))
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	h.logger.Info("ListOrganizations request", zap.Int32("page", req.GetPage()), zap.Int32("per_page", req.GetPerPage()))

	orgs, total, err := h.service.ListOrganizations(ctx, req.GetPage(), req.GetPerPage(), req.GetOwnerId())
	if err != nil {
		return nil, mapError(err)
	}

	protoOrgs := make([]*pb.Organization, len(orgs))
	for i, org := range orgs {
		protoOrgs[i] = organizationToProto(org)
	}

	return &pb.ListOrganizationsResponse{
		Organizations: protoOrgs,
		PageInfo: &common.PageInfo{
			Page:    req.GetPage(),
			PerPage: req.GetPerPage(),
			Total:   int64(total),
		},
	}, nil
}

// GetOrganizationMembers получает участников организации
func (h *CompanyHandler) GetOrganizationMembers(ctx context.Context, req *pb.GetOrganizationMembersRequest) (*pb.ListEmployeesResponse, error) {
	if err := ValidateGetOrganizationMembersRequest(req.GetOrganizationId(), req.GetPage(), req.GetPerPage()); err != nil {
		h.logger.Warn("GetOrganizationMembers validation failed", zap.Error(err))
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	h.logger.Info("GetOrganizationMembers request", zap.String("org_id", req.GetOrganizationId()))

	employees, total, err := h.service.GetOrganizationMembers(ctx, req.GetOrganizationId(), req.GetPage(), req.GetPerPage())
	if err != nil {
		return nil, mapError(err)
	}

	protoEmployees := make([]*pb.Employee, len(employees))
	for i, emp := range employees {
		protoEmployees[i] = employeeToProto(emp)
	}

	return &pb.ListEmployeesResponse{
		Employees: protoEmployees,
		PageInfo: &common.PageInfo{
			Page:    req.GetPage(),
			PerPage: req.GetPerPage(),
			Total:   int64(total),
		},
	}, nil
}

// AddMember добавляет участника в организацию
func (h *CompanyHandler) AddMember(ctx context.Context, req *pb.AddMemberRequest) (*pb.Employee, error) {
	if err := ValidateAddMemberRequest(req.GetOrganizationId(), req.GetUserId(), req.GetRole()); err != nil {
		h.logger.Warn("AddMember validation failed", zap.Error(err))
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	// Извлекаем userID из контекста
	userID, err := ExtractUserID(ctx)
	if err != nil {
		h.logger.Warn("AddMember: failed to extract user id", zap.Error(err))
		return nil, status.Error(codes.Unauthenticated, "user authentication required")
	}

	// Проверяем право добавления участника (должен быть владельцем или админом)
	if err := h.service.CheckUserIsOwner(ctx, userID, req.GetOrganizationId()); err != nil {
		h.logger.Warn("AddMember: unauthorized", zap.Error(err), zap.String("user_id", userID), zap.String("org_id", req.GetOrganizationId()))
		return nil, status.Error(codes.PermissionDenied, "user does not have permission to add members")
	}

	h.logger.Info("AddMember request", zap.String("org_id", req.GetOrganizationId()), zap.String("user_id", req.GetUserId()), zap.String("requester_id", userID))

	emp, err := h.service.AddMember(ctx, req.GetOrganizationId(), req.GetUserId(), req.GetRole())
	if err != nil {
		return nil, mapError(err)
	}

	return employeeToProto(emp), nil
}

// RemoveMember удаляет участника из организации
func (h *CompanyHandler) RemoveMember(ctx context.Context, req *pb.RemoveMemberRequest) (*common.Empty, error) {
	if err := ValidateRemoveMemberRequest(req.GetOrganizationId(), req.GetEmployeeId()); err != nil {
		h.logger.Warn("RemoveMember validation failed", zap.Error(err))
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	// Извлекаем userID из контекста
	userID, err := ExtractUserID(ctx)
	if err != nil {
		h.logger.Warn("RemoveMember: failed to extract user id", zap.Error(err))
		return nil, status.Error(codes.Unauthenticated, "user authentication required")
	}

	// Проверяем право удаления участника (должен быть владельцем)
	if err := h.service.CheckUserIsOwner(ctx, userID, req.GetOrganizationId()); err != nil {
		h.logger.Warn("RemoveMember: unauthorized", zap.Error(err), zap.String("user_id", userID), zap.String("org_id", req.GetOrganizationId()))
		return nil, status.Error(codes.PermissionDenied, "user does not have permission to remove members")
	}

	h.logger.Info("RemoveMember request", zap.String("org_id", req.GetOrganizationId()), zap.String("employee_id", req.GetEmployeeId()), zap.String("requester_id", userID))

	err = h.service.RemoveMember(ctx, req.GetOrganizationId(), req.GetEmployeeId())
	if err != nil {
		return nil, mapError(err)
	}

	return &common.Empty{}, nil
}

// organizationToProto преобразует domain.Organization в proto
func organizationToProto(org *domain.Organization) *pb.Organization {
	var desc string
	if org.Description != nil {
		desc = *org.Description
	}
	return &pb.Organization{
		Id:          org.ID,
		Name:        org.Name,
		Description: desc,
		OwnerId:     org.OwnerID,
		CreatedAt:   org.CreatedAt.UnixMilli(),
		UpdatedAt:   org.UpdatedAt.UnixMilli(),
	}
}

// employeeToProto преобразует domain.Employee в proto
func employeeToProto(emp *domain.Employee) *pb.Employee {
	return &pb.Employee{
		Id:             emp.ID,
		OrganizationId: emp.OrganizationID,
		UserId:         emp.UserID,
		Role:           emp.Role,
		JoinedAt:       emp.JoinedAt.UnixMilli(),
	}
}

// mapError преобразует domain ошибки в gRPC статусы
func mapError(err error) error {
	switch err {
	case domain.ErrOrganizationNotFound, domain.ErrEmployeeNotFound:
		return status.Error(codes.NotFound, err.Error())
	case domain.ErrOrganizationExists, domain.ErrEmployeeExists:
		return status.Error(codes.AlreadyExists, err.Error())
	case domain.ErrInvalidOrganizationID, domain.ErrInvalidOrganizationName,
		domain.ErrOrganizationNameRequired, domain.ErrInvalidEmployeeID,
		domain.ErrInvalidUserID, domain.ErrInvalidRole, domain.ErrInvalidInput,
		domain.ErrMissingRequiredField:
		return status.Error(codes.InvalidArgument, err.Error())
	case domain.ErrUnauthorized, domain.ErrPermissionDenied, company.ErrUnauthorized:
		return status.Error(codes.PermissionDenied, err.Error())
	default:
		return status.Error(codes.Internal, "internal server error")
	}
}
