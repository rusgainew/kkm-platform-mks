// Файл user-query-server/internal/interfaces/grpc/handlers/user_query_handler.go содержит реализацию пакета handlers.
package handlers

import (
	"context"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/rusgainew/kkm-project-mks/proto-lib/api"
	"github.com/rusgainew/kkm-project-mks/user-query-server/internal/application/ports"
)

type UserQueryHandler struct {
	pb.UnimplementedUserQueryServiceServer
	logger *zap.Logger
	repo   ports.UserRepository
	cache  ports.UserCache
}

func NewUserQueryHandler(
	logger *zap.Logger,
	repo ports.UserRepository,
	cache ports.UserCache,
) *UserQueryHandler {
	return &UserQueryHandler{
		logger: logger,
		repo:   repo,
		cache:  cache,
	}
}

// GetUser retrieves a single user by ID
func (h *UserQueryHandler) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	if req.UserId == "" {
		h.logger.Warn("GetUser called with empty user_id")
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}

	// Try cache first
	if h.cache != nil {
		cached, err := h.cache.GetUser(ctx, req.UserId)
		if err == nil && cached != nil {
			h.logger.Debug("User cache hit", zap.String("user_id", req.UserId))
			return &pb.GetUserResponse{User: cached}, nil
		}
	}

	// Query from read-model
	user, err := h.repo.GetUser(ctx, req.UserId)
	if err != nil {
		h.logger.Error("Failed to get user",
			zap.String("user_id", req.UserId),
			zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to get user")
	}

	if user == nil {
		h.logger.Warn("User not found", zap.String("user_id", req.UserId))
		return nil, status.Error(codes.NotFound, "user not found")
	}

	// Cache result
	if h.cache != nil {
		_ = h.cache.SetUser(ctx, user)
	}

	return &pb.GetUserResponse{User: user}, nil
}

// ListUsers retrieves users with pagination and optional filters
func (h *UserQueryHandler) ListUsers(ctx context.Context, req *pb.ListUsersRequest) (*pb.ListUsersResponse, error) {
	// Validate pagination
	if req.Pagination == nil {
		h.logger.Warn("ListUsers called without pagination")
		return nil, status.Error(codes.InvalidArgument, "pagination is required")
	}

	size := req.Pagination.Size
	if size <= 0 {
		size = 10
	}

	if size > 100 {
		size = 100
	}

	page := req.Pagination.Page
	if page < 0 {
		page = 0
	}

	offset := page * size

	// Query from read-model
	users, totalCount, err := h.repo.ListUsers(ctx, offset, size, req.Status, req.Role)
	if err != nil {
		h.logger.Error("Failed to list users",
			zap.Int32("page", page),
			zap.Int32("size", size),
			zap.String("status", req.Status),
			zap.String("role", req.Role),
			zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to list users")
	}

	h.logger.Debug("Listed users",
		zap.Int32("page", page),
		zap.Int32("size", size),
		zap.Int("count", len(users)),
		zap.Int64("total", totalCount))

	return &pb.ListUsersResponse{
		Users:      users,
		TotalCount: totalCount,
		Page:       page,
		PerPage:    size,
	}, nil
}

// SearchUsers performs search on users
func (h *UserQueryHandler) SearchUsers(ctx context.Context, req *pb.SearchUsersRequest) (*pb.SearchUsersResponse, error) {
	// Validate search query
	if req.Query == "" {
		h.logger.Warn("SearchUsers called with empty query")
		return nil, status.Error(codes.InvalidArgument, "search query is required")
	}

	// Validate pagination
	if req.Pagination == nil {
		h.logger.Warn("SearchUsers called without pagination")
		return nil, status.Error(codes.InvalidArgument, "pagination is required")
	}

	size := req.Pagination.Size
	if size <= 0 {
		size = 10
	}

	if size > 100 {
		size = 100
	}

	page := req.Pagination.Page
	if page < 0 {
		page = 0
	}

	offset := page * size

	// Query from read-model
	users, totalCount, err := h.repo.SearchUsers(ctx, req.Query, offset, size, req.Status)
	if err != nil {
		h.logger.Error("Failed to search users",
			zap.String("query", req.Query),
			zap.Int32("page", page),
			zap.String("status", req.Status),
			zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to search users")
	}

	h.logger.Debug("Searched users",
		zap.String("query", req.Query),
		zap.Int32("page", page),
		zap.Int("count", len(users)),
		zap.Int64("total", totalCount))

	return &pb.SearchUsersResponse{
		Users:      users,
		TotalCount: totalCount,
		Page:       page,
		PerPage:    size,
	}, nil
}
