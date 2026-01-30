// Файл user-query-server/internal/interfaces/grpc/handlers/user_query_handler_test.go содержит реализацию пакета handlers.
package handlers

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"

	pb "github.com/rusgainew/kkm-project-mks/proto-lib/api"
)

// MockUserRepository мок для repository
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) GetUser(ctx context.Context, userID string) (*pb.UserReadModel, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*pb.UserReadModel), args.Error(1)
}

func (m *MockUserRepository) ListUsers(ctx context.Context, offset, limit int32, status, role string) ([]*pb.UserReadModel, int64, error) {
	args := m.Called(ctx, offset, limit, status, role)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).([]*pb.UserReadModel), args.Get(1).(int64), args.Error(2)
}

func (m *MockUserRepository) SearchUsers(ctx context.Context, query string, offset, limit int32, status string) ([]*pb.UserReadModel, int64, error) {
	args := m.Called(ctx, query, offset, limit, status)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).([]*pb.UserReadModel), args.Get(1).(int64), args.Error(2)
}

func (m *MockUserRepository) UpsertUser(ctx context.Context, user *pb.UserReadModel) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

// MockRedisCache мок для cache
type MockRedisCache struct {
	mock.Mock
}

func (m *MockRedisCache) GetUser(ctx context.Context, userID string) (*pb.UserReadModel, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*pb.UserReadModel), args.Error(1)
}

func (m *MockRedisCache) SetUser(ctx context.Context, user *pb.UserReadModel) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockRedisCache) InvalidateUser(ctx context.Context, userID string) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

// Test GetUser
func TestGetUser_Success(t *testing.T) {
	logger, _ := zap.NewProduction()
	mockRepo := new(MockUserRepository)
	mockCache := new(MockRedisCache)

	expectedUser := &pb.UserReadModel{
		Id:        "user-123",
		Email:     "user@example.com",
		FirstName: "John",
		LastName:  "Doe",
		Status:    "active",
		Role:      "admin",
	}

	mockCache.On("GetUser", mock.Anything, "user-123").Return(nil, nil)
	mockRepo.On("GetUser", mock.Anything, "user-123").Return(expectedUser, nil)
	mockCache.On("SetUser", mock.Anything, expectedUser).Return(nil)

	handler := &UserQueryHandler{
		logger: logger,
		repo:   mockRepo,
		cache:  mockCache,
	}

	req := &pb.GetUserRequest{UserId: "user-123"}
	resp, err := handler.GetUser(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "user-123", resp.User.Id)
	assert.Equal(t, "user@example.com", resp.User.Email)
	assert.Equal(t, "admin", resp.User.Role)

	mockCache.AssertCalled(t, "GetUser", mock.Anything, "user-123")
	mockRepo.AssertCalled(t, "GetUser", mock.Anything, "user-123")
	mockCache.AssertCalled(t, "SetUser", mock.Anything, expectedUser)
}

func TestGetUser_CacheHit(t *testing.T) {
	logger, _ := zap.NewProduction()
	mockRepo := new(MockUserRepository)
	mockCache := new(MockRedisCache)

	cachedUser := &pb.UserReadModel{
		Id:     "user-456",
		Email:  "cached@example.com",
		Role:   "user",
		Status: "active",
	}

	mockCache.On("GetUser", mock.Anything, "user-456").Return(cachedUser, nil)

	handler := &UserQueryHandler{
		logger: logger,
		repo:   mockRepo,
		cache:  mockCache,
	}

	req := &pb.GetUserRequest{UserId: "user-456"}
	resp, err := handler.GetUser(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "user-456", resp.User.Id)

	mockCache.AssertCalled(t, "GetUser", mock.Anything, "user-456")
	mockRepo.AssertNotCalled(t, "GetUser")
}

func TestGetUser_NotFound(t *testing.T) {
	logger, _ := zap.NewProduction()
	mockRepo := new(MockUserRepository)
	mockCache := new(MockRedisCache)

	mockCache.On("GetUser", mock.Anything, "nonexistent").Return(nil, nil)
	mockRepo.On("GetUser", mock.Anything, "nonexistent").Return(nil, nil)

	handler := &UserQueryHandler{
		logger: logger,
		repo:   mockRepo,
		cache:  mockCache,
	}

	req := &pb.GetUserRequest{UserId: "nonexistent"}
	resp, err := handler.GetUser(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, resp)
}

func TestGetUser_EmptyID(t *testing.T) {
	logger, _ := zap.NewProduction()
	mockRepo := new(MockUserRepository)
	mockCache := new(MockRedisCache)

	handler := &UserQueryHandler{
		logger: logger,
		repo:   mockRepo,
		cache:  mockCache,
	}

	req := &pb.GetUserRequest{UserId: ""}
	resp, err := handler.GetUser(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, resp)

	mockCache.AssertNotCalled(t, "GetUser")
	mockRepo.AssertNotCalled(t, "GetUser")
}

// Test ListUsers
func TestListUsers_Success(t *testing.T) {
	logger, _ := zap.NewProduction()
	mockRepo := new(MockUserRepository)
	mockCache := new(MockRedisCache)

	expectedUsers := []*pb.UserReadModel{
		{Id: "user-1", Email: "user1@example.com", Status: "active"},
		{Id: "user-2", Email: "user2@example.com", Status: "active"},
	}

	mockRepo.On("ListUsers", mock.Anything, int32(0), int32(10), "", "").Return(expectedUsers, int64(2), nil)

	handler := &UserQueryHandler{
		logger: logger,
		repo:   mockRepo,
		cache:  mockCache,
	}

	req := &pb.ListUsersRequest{
		Pagination: &pb.PageInfo{Page: 0, Size: 10},
		Status:     "",
		Role:       "",
	}
	resp, err := handler.ListUsers(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Len(t, resp.Users, 2)
	assert.Equal(t, int64(2), resp.TotalCount)

	mockRepo.AssertCalled(t, "ListUsers", mock.Anything, int32(0), int32(10), "", "")
}

func TestListUsers_WithFilters(t *testing.T) {
	logger, _ := zap.NewProduction()
	mockRepo := new(MockUserRepository)
	mockCache := new(MockRedisCache)

	expectedUsers := []*pb.UserReadModel{
		{Id: "admin-1", Email: "admin@example.com", Role: "admin", Status: "active"},
	}

	mockRepo.On("ListUsers", mock.Anything, int32(0), int32(10), "active", "admin").Return(expectedUsers, int64(1), nil)

	handler := &UserQueryHandler{
		logger: logger,
		repo:   mockRepo,
		cache:  mockCache,
	}

	req := &pb.ListUsersRequest{
		Pagination: &pb.PageInfo{Page: 0, Size: 10},
		Status:     "active",
		Role:       "admin",
	}
	resp, err := handler.ListUsers(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Len(t, resp.Users, 1)
	assert.Equal(t, "admin", resp.Users[0].Role)

	mockRepo.AssertCalled(t, "ListUsers", mock.Anything, int32(0), int32(10), "active", "admin")
}

func TestListUsers_Pagination(t *testing.T) {
	logger, _ := zap.NewProduction()
	mockRepo := new(MockUserRepository)
	mockCache := new(MockRedisCache)

	tests := []struct {
		page        int32
		size        int32
		expectedOff int32
		expectedSz  int32
	}{
		{0, 10, 0, 10},   // Page 0 → offset 0, size 10
		{1, 10, 10, 10},  // Page 1 → offset 10, size 10
		{2, 20, 40, 20},  // Page 2, size 20 → offset 40, size 20
		{0, 0, 0, 10},    // size 0 defaults to 10 → offset 0, size 10
		{0, 101, 0, 100}, // size 101 capped to 100 → offset 0, size 100
	}

	for _, tc := range tests {
		mockRepo.On("ListUsers", mock.Anything, tc.expectedOff, tc.expectedSz, "", "").
			Return([]*pb.UserReadModel{}, int64(0), nil).
			Times(1)

		handler := &UserQueryHandler{
			logger: logger,
			repo:   mockRepo,
			cache:  mockCache,
		}

		req := &pb.ListUsersRequest{
			Pagination: &pb.PageInfo{Page: tc.page, Size: tc.size},
		}
		_, err := handler.ListUsers(context.Background(), req)

		assert.NoError(t, err)
		mockRepo.AssertCalled(t, "ListUsers", mock.Anything, tc.expectedOff, tc.expectedSz, "", "")
		mockRepo.ExpectedCalls = nil
		mockRepo.Calls = nil
	}
}

// Test SearchUsers
func TestSearchUsers_Success(t *testing.T) {
	logger, _ := zap.NewProduction()
	mockRepo := new(MockUserRepository)
	mockCache := new(MockRedisCache)

	expectedUsers := []*pb.UserReadModel{
		{Id: "user-123", Email: "john@example.com"},
	}

	mockRepo.On("SearchUsers", mock.Anything, "john", int32(0), int32(10), "").Return(expectedUsers, int64(1), nil)

	handler := &UserQueryHandler{
		logger: logger,
		repo:   mockRepo,
		cache:  mockCache,
	}

	req := &pb.SearchUsersRequest{
		Query:      "john",
		Pagination: &pb.PageInfo{Page: 0, Size: 10},
		Status:     "",
	}
	resp, err := handler.SearchUsers(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Len(t, resp.Users, 1)

	mockRepo.AssertCalled(t, "SearchUsers", mock.Anything, "john", int32(0), int32(10), "")
}

func TestSearchUsers_EmptyQuery(t *testing.T) {
	logger, _ := zap.NewProduction()
	mockRepo := new(MockUserRepository)
	mockCache := new(MockRedisCache)

	handler := &UserQueryHandler{
		logger: logger,
		repo:   mockRepo,
		cache:  mockCache,
	}

	req := &pb.SearchUsersRequest{
		Query:      "",
		Pagination: &pb.PageInfo{Page: 0, Size: 10},
	}
	resp, err := handler.SearchUsers(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, resp)

	mockRepo.AssertNotCalled(t, "SearchUsers")
}
