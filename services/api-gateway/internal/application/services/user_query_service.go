// Файл api-gateway/internal/application/services/user_query_service.go содержит реализацию пакета services.
package services

import (
	"context"
	"fmt"

	"github.com/rusgainew/kkm-project-mks/api-gateway/internal/infrastructure/client"
	"github.com/rusgainew/kkm-project-mks/api-gateway/internal/infrastructure/observability"
	pb "github.com/rusgainew/kkm-project-mks/proto-lib/api"
	"go.uber.org/zap"
)

// UserQueryService сервис для чтения пользователей (query side)
type UserQueryService struct {
	connManager *client.ConnectionManager
	serviceAddr string
	metrics     *observability.Metrics
	tracer      *observability.Tracer
	logger      *zap.Logger
}

// NewUserQueryService создает новый UserQueryService
func NewUserQueryService(
	connManager *client.ConnectionManager,
	serviceAddr string,
	metrics *observability.Metrics,
	tracer *observability.Tracer,
	logger *zap.Logger,
) *UserQueryService {
	return &UserQueryService{
		connManager: connManager,
		serviceAddr: serviceAddr,
		metrics:     metrics,
		tracer:      tracer,
		logger:      logger,
	}
}

// getClient возвращает gRPC клиент для user query service
func (s *UserQueryService) getClient(ctx context.Context) (pb.UserQueryServiceClient, error) {
	conn, err := s.connManager.GetConnection(ctx, s.serviceAddr)
	if err != nil {
		return nil, fmt.Errorf("failed to get connection: %w", err)
	}
	return pb.NewUserQueryServiceClient(conn), nil
}

// GetUser получает пользователя по ID
func (s *UserQueryService) GetUser(ctx context.Context, userID string) (*pb.UserReadModel, error) {
	ctx, endSpan := s.tracer.StartSpan(ctx, "UserQueryService.GetUser")
	defer endSpan()

	client, err := s.getClient(ctx)
	if err != nil {
		s.logger.Error("Failed to get user query client", zap.Error(err))
		s.metrics.IncrementErrorCount("user_query_service", "get_user")
		return nil, fmt.Errorf("failed to get client: %w", err)
	}

	req := &pb.GetUserRequest{
		UserId: userID,
	}

	resp, err := client.GetUser(ctx, req)
	if err != nil {
		s.logger.Error("Failed to get user", zap.String("user_id", userID), zap.Error(err))
		s.metrics.IncrementErrorCount("user_query_service", "get_user")
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return resp.User, nil
}

// ListUsers получает список пользователей с пагинацией
func (s *UserQueryService) ListUsers(ctx context.Context, page, pageSize int32, status, role string) (*pb.ListUsersResponse, error) {
	ctx, endSpan := s.tracer.StartSpan(ctx, "UserQueryService.ListUsers")
	defer endSpan()

	client, err := s.getClient(ctx)
	if err != nil {
		s.logger.Error("Failed to get user query client", zap.Error(err))
		s.metrics.IncrementErrorCount("user_query_service", "list_users")
		return nil, fmt.Errorf("failed to get client: %w", err)
	}

	req := &pb.ListUsersRequest{
		Pagination: &pb.PageInfo{
			Page: page,
			Size: pageSize,
		},
		Status: status,
		Role:   role,
	}

	resp, err := client.ListUsers(ctx, req)
	if err != nil {
		s.logger.Error("Failed to list users",
			zap.Int32("page", page),
			zap.Int32("page_size", pageSize),
			zap.Error(err))
		s.metrics.IncrementErrorCount("user_query_service", "list_users")
		return nil, fmt.Errorf("failed to list users: %w", err)
	}

	return resp, nil
}

// SearchUsers выполняет поиск пользователей
func (s *UserQueryService) SearchUsers(ctx context.Context, query string, page, pageSize int32, status string) (*pb.SearchUsersResponse, error) {
	ctx, endSpan := s.tracer.StartSpan(ctx, "UserQueryService.SearchUsers")
	defer endSpan()

	client, err := s.getClient(ctx)
	if err != nil {
		s.logger.Error("Failed to get user query client", zap.Error(err))
		s.metrics.IncrementErrorCount("user_query_service", "search_users")
		return nil, fmt.Errorf("failed to get client: %w", err)
	}

	req := &pb.SearchUsersRequest{
		Query: query,
		Pagination: &pb.PageInfo{
			Page: page,
			Size: pageSize,
		},
		Status: status,
	}

	resp, err := client.SearchUsers(ctx, req)
	if err != nil {
		s.logger.Error("Failed to search users",
			zap.String("query", query),
			zap.Int32("page", page),
			zap.Error(err))
		s.metrics.IncrementErrorCount("user_query_service", "search_users")
		return nil, fmt.Errorf("failed to search users: %w", err)
	}

	return resp, nil
}
