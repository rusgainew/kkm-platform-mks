package grpc_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/rusgainew/kkm-project-mks/foreign-company-server/internal/domain"
	grpcHandler "github.com/rusgainew/kkm-project-mks/foreign-company-server/internal/interfaces/grpc"
	pb "github.com/rusgainew/kkm-project-mks/proto-lib/api"
)

// MockService implements the service interface for testing
type MockService struct {
	CreateFunc func(ctx context.Context, pin, fullName, countryCode string, createdBy uuid.UUID) (*domain.ForeignCompany, error)
	UpdateFunc func(ctx context.Context, id int64, pin, fullName, countryCode, address string, updatedBy uuid.UUID) (*domain.ForeignCompany, error)
	GetFunc    func(ctx context.Context, id int64) (*domain.ForeignCompany, error)
	DeleteFunc func(ctx context.Context, id int64, deletedBy uuid.UUID) error
}

func (m *MockService) CreateForeignCompany(ctx context.Context, pin, fullName, countryCode string, createdBy uuid.UUID) (*domain.ForeignCompany, error) {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, pin, fullName, countryCode, createdBy)
	}
	return &domain.ForeignCompany{ID: 1, PIN: pin, FullName: fullName}, nil
}

func (m *MockService) UpdateForeignCompany(ctx context.Context, id int64, pin, fullName, countryCode, address string, updatedBy uuid.UUID) (*domain.ForeignCompany, error) {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(ctx, id, pin, fullName, countryCode, address, updatedBy)
	}
	return &domain.ForeignCompany{ID: id, PIN: pin, FullName: fullName}, nil
}

func (m *MockService) GetForeignCompany(ctx context.Context, id int64) (*domain.ForeignCompany, error) {
	if m.GetFunc != nil {
		return m.GetFunc(ctx, id)
	}
	return &domain.ForeignCompany{ID: id, PIN: "123456789", FullName: "Test Company"}, nil
}

func (m *MockService) DeleteForeignCompany(ctx context.Context, id int64, deletedBy uuid.UUID) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(ctx, id, deletedBy)
	}
	return nil
}

func TestCreateForeignCompany_Success(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	userID := uuid.New()

	mockService := &MockService{
		CreateFunc: func(ctx context.Context, pin, fullName, countryCode string, createdBy uuid.UUID) (*domain.ForeignCompany, error) {
			assert.Equal(t, "123456789", pin)
			assert.Equal(t, "Test Company", fullName)
			assert.Equal(t, userID, createdBy)
			return &domain.ForeignCompany{ID: 1, PIN: pin, FullName: fullName}, nil
		},
	}

	handler := grpcHandler.NewForeignCompanyHandler(mockService, logger)
	ctx := context.WithValue(context.Background(), "user_id", userID.String())

	req := &pb.CreateForeignCompanyRequest{
		Pin:      "123456789",
		FullName: "Test Company",
	}

	resp, err := handler.CreateForeignCompany(ctx, req)

	require.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, int64(1), resp.Id)
	assert.Equal(t, "123456789", resp.Pin)
	assert.Equal(t, "Test Company", resp.FullName)
}

func TestCreateForeignCompany_NoAuth(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	mockService := &MockService{}
	handler := grpcHandler.NewForeignCompanyHandler(mockService, logger)

	req := &pb.CreateForeignCompanyRequest{
		Pin:      "123456789",
		FullName: "Test",
	}

	resp, err := handler.CreateForeignCompany(context.Background(), req)

	require.Error(t, err)
	assert.Nil(t, resp)
	st, ok := status.FromError(err)
	require.True(t, ok)
	assert.Equal(t, codes.Unauthenticated, st.Code())
}

func TestUpdateForeignCompany_Success(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	userID := uuid.New()

	mockService := &MockService{
		UpdateFunc: func(ctx context.Context, id int64, pin, fullName, countryCode, address string, updatedBy uuid.UUID) (*domain.ForeignCompany, error) {
			assert.Equal(t, int64(1), id)
			assert.Equal(t, "987654321", pin)
			assert.Equal(t, userID, updatedBy)
			return &domain.ForeignCompany{ID: id, PIN: pin, FullName: fullName}, nil
		},
	}

	handler := grpcHandler.NewForeignCompanyHandler(mockService, logger)
	ctx := context.WithValue(context.Background(), "user_id", userID.String())

	req := &pb.UpdateForeignCompanyRequest{
		Id:       1,
		Pin:      "987654321",
		FullName: "Updated",
	}

	resp, err := handler.UpdateForeignCompany(ctx, req)

	require.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, int64(1), resp.Id)
}
