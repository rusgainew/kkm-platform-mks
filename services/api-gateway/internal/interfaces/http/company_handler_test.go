package http_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/rusgainew/kkm-project-mks/api-gateway/internal/domain/models"
	handler "github.com/rusgainew/kkm-project-mks/api-gateway/internal/interfaces/http"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// MockCompanyService мок для CompanyService
type MockCompanyService struct {
	mock.Mock
}

func (m *MockCompanyService) CreateCompany(ctx context.Context, company *models.Company) (*models.Company, error) {
	args := m.Called(ctx, company)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Company), args.Error(1)
}

func (m *MockCompanyService) GetCompany(ctx context.Context, id string) (*models.Company, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Company), args.Error(1)
}

func (m *MockCompanyService) UpdateCompany(ctx context.Context, company *models.Company) (*models.Company, error) {
	args := m.Called(ctx, company)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Company), args.Error(1)
}

func (m *MockCompanyService) DeleteCompany(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockCompanyService) ListCompanies(ctx context.Context, page, pageSize int) ([]*models.Company, int, error) {
	args := m.Called(ctx, page, pageSize)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).([]*models.Company), args.Int(1), args.Error(2)
}

// TestCompanyHandler_CreateCompany_Success тест успешного создания
func TestCompanyHandler_CreateCompany_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	logger := zap.NewNop()

	mockService := new(MockCompanyService)
	expectedCompany := &models.Company{
		ID:          "550e8400-e29b-41d4-a716-446655440000",
		Name:        "Test Company",
		Description: "Test Description",
	}

	mockService.On("CreateCompany", mock.Anything, mock.AnythingOfType("*models.Company")).Return(expectedCompany, nil)

	h := handler.NewCompanyHandler(mockService, logger)
	router := gin.New()
	router.POST("/companies", h.CreateCompany)

	body, _ := json.Marshal(models.Company{Name: "Test Company", Description: "Test Description"})
	req := httptest.NewRequest("POST", "/companies", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response models.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response.Success)

	mockService.AssertExpectations(t)
}

// TestCompanyHandler_CreateCompany_GRPCError тест обработки gRPC ошибок
func TestCompanyHandler_CreateCompany_GRPCError(t *testing.T) {
	tests := []struct {
		name           string
		grpcErr        error
		expectedStatus int
		expectedCode   string
	}{
		{
			name:           "already_exists",
			grpcErr:        status.Error(codes.AlreadyExists, "company already exists"),
			expectedStatus: http.StatusConflict,
			expectedCode:   "ALREADY_EXISTS",
		},
		{
			name:           "internal_error",
			grpcErr:        status.Error(codes.Internal, "internal error"),
			expectedStatus: http.StatusInternalServerError,
			expectedCode:   "INTERNAL_SERVER_ERROR",
		},
		{
			name:           "unavailable",
			grpcErr:        status.Error(codes.Unavailable, "service unavailable"),
			expectedStatus: http.StatusServiceUnavailable,
			expectedCode:   "SERVICE_UNAVAILABLE",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)
			logger := zap.NewNop()

			mockService := new(MockCompanyService)
			mockService.On("CreateCompany", mock.Anything, mock.AnythingOfType("*models.Company")).Return(nil, tt.grpcErr)

			h := handler.NewCompanyHandler(mockService, logger)
			router := gin.New()
			router.POST("/companies", h.CreateCompany)

			body, _ := json.Marshal(models.Company{Name: "Test Company", Description: "Test"})
			req := httptest.NewRequest("POST", "/companies", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			var response models.APIResponse
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.False(t, response.Success)
			assert.NotNil(t, response.Error)
			assert.Equal(t, tt.expectedCode, response.Error.Code)

			mockService.AssertExpectations(t)
		})
	}
}

// TestCompanyHandler_CreateCompany_ValidationError тесты валидации
func TestCompanyHandler_CreateCompany_ValidationError(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    interface{}
		expectedStatus int
		expectedCode   string
	}{
		{
			name:           "invalid_json",
			requestBody:    "invalid json",
			expectedStatus: http.StatusBadRequest,
			expectedCode:   "INVALID_INPUT",
		},
		{
			name: "empty_name",
			requestBody: models.Company{
				Name:        "",
				Description: "Test Description",
			},
			expectedStatus: http.StatusBadRequest,
			expectedCode:   "VALIDATION_FAILED",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)
			logger := zap.NewNop()

			mockService := new(MockCompanyService)
			h := handler.NewCompanyHandler(mockService, logger)
			router := gin.New()
			router.POST("/companies", h.CreateCompany)

			var body []byte
			switch v := tt.requestBody.(type) {
			case string:
				body = []byte(v)
			default:
				body, _ = json.Marshal(v)
			}

			req := httptest.NewRequest("POST", "/companies", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			var response models.APIResponse
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.False(t, response.Success)
			assert.Equal(t, tt.expectedCode, response.Error.Code)
		})
	}
}

// TestCompanyHandler_GetCompany_Success тест успешного получения
func TestCompanyHandler_GetCompany_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	logger := zap.NewNop()

	mockService := new(MockCompanyService)
	expectedCompany := &models.Company{
		ID:          "550e8400-e29b-41d4-a716-446655440000",
		Name:        "Test Company",
		Description: "Test Description",
	}

	mockService.On("GetCompany", mock.Anything, "550e8400-e29b-41d4-a716-446655440000").Return(expectedCompany, nil)

	h := handler.NewCompanyHandler(mockService, logger)
	router := gin.New()
	router.GET("/companies/:id", h.GetCompany)

	req := httptest.NewRequest("GET", "/companies/550e8400-e29b-41d4-a716-446655440000", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response models.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response.Success)

	mockService.AssertExpectations(t)
}

// TestCompanyHandler_GetCompany_NotFound тест когда компания не найдена
func TestCompanyHandler_GetCompany_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	logger := zap.NewNop()

	mockService := new(MockCompanyService)
	mockService.On("GetCompany", mock.Anything, "550e8400-e29b-41d4-a716-446655440001").Return(nil, status.Error(codes.NotFound, "company not found"))

	h := handler.NewCompanyHandler(mockService, logger)
	router := gin.New()
	router.GET("/companies/:id", h.GetCompany)

	req := httptest.NewRequest("GET", "/companies/550e8400-e29b-41d4-a716-446655440001", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)

	var response models.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.False(t, response.Success)
	assert.Equal(t, "NOT_FOUND", response.Error.Code)

	mockService.AssertExpectations(t)
}

// TestCompanyHandler_GetCompany_InvalidUUID тест невалидного UUID
func TestCompanyHandler_GetCompany_InvalidUUID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	logger := zap.NewNop()

	mockService := new(MockCompanyService)
	h := handler.NewCompanyHandler(mockService, logger)
	router := gin.New()
	router.GET("/companies/:id", h.GetCompany)

	req := httptest.NewRequest("GET", "/companies/not-a-uuid", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response models.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.False(t, response.Success)
	assert.Equal(t, "INVALID_ID", response.Error.Code)
}

// TestCompanyHandler_UpdateCompany_Success тест успешного обновления
func TestCompanyHandler_UpdateCompany_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	logger := zap.NewNop()

	mockService := new(MockCompanyService)
	expectedCompany := &models.Company{
		ID:          "550e8400-e29b-41d4-a716-446655440000",
		Name:        "Updated Company",
		Description: "Updated Description",
	}

	mockService.On("UpdateCompany", mock.Anything, mock.AnythingOfType("*models.Company")).Return(expectedCompany, nil)

	h := handler.NewCompanyHandler(mockService, logger)
	router := gin.New()
	router.PUT("/companies/:id", h.UpdateCompany)

	body, _ := json.Marshal(models.Company{Name: "Updated Company", Description: "Updated Description"})
	req := httptest.NewRequest("PUT", "/companies/550e8400-e29b-41d4-a716-446655440000", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response models.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response.Success)

	mockService.AssertExpectations(t)
}

// TestCompanyHandler_DeleteCompany_Success тест успешного удаления
func TestCompanyHandler_DeleteCompany_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	logger := zap.NewNop()

	mockService := new(MockCompanyService)
	mockService.On("DeleteCompany", mock.Anything, "550e8400-e29b-41d4-a716-446655440000").Return(nil)

	h := handler.NewCompanyHandler(mockService, logger)
	router := gin.New()
	router.DELETE("/companies/:id", h.DeleteCompany)

	req := httptest.NewRequest("DELETE", "/companies/550e8400-e29b-41d4-a716-446655440000", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response models.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response.Success)

	mockService.AssertExpectations(t)
}

// TestCompanyHandler_ListCompanies_Success тест списка компаний
func TestCompanyHandler_ListCompanies_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	logger := zap.NewNop()

	mockService := new(MockCompanyService)
	companies := []*models.Company{
		{ID: "550e8400-e29b-41d4-a716-446655440000", Name: "Company 1"},
		{ID: "550e8400-e29b-41d4-a716-446655440001", Name: "Company 2"},
	}

	mockService.On("ListCompanies", mock.Anything, 1, 10).Return(companies, 2, nil)

	h := handler.NewCompanyHandler(mockService, logger)
	router := gin.New()
	router.GET("/companies", h.ListCompanies)

	req := httptest.NewRequest("GET", "/companies?page=1&page_size=10", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response models.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response.Success)
	assert.NotNil(t, response.Meta)

	mockService.AssertExpectations(t)
}

// TestCompanyHandler_ListCompanies_DefaultPagination тест дефолтной пагинации
func TestCompanyHandler_ListCompanies_DefaultPagination(t *testing.T) {
	gin.SetMode(gin.TestMode)
	logger := zap.NewNop()

	mockService := new(MockCompanyService)
	companies := []*models.Company{}

	// Default page=1, page_size=10
	mockService.On("ListCompanies", mock.Anything, 1, 10).Return(companies, 0, nil)

	h := handler.NewCompanyHandler(mockService, logger)
	router := gin.New()
	router.GET("/companies", h.ListCompanies)

	req := httptest.NewRequest("GET", "/companies", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}
