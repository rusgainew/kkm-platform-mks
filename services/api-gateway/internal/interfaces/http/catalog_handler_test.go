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

// MockCatalogService мок для CatalogServiceInterface
type MockCatalogService struct {
	mock.Mock
}

func (m *MockCatalogService) GetCatalog(ctx context.Context, id string) (*models.Catalog, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Catalog), args.Error(1)
}

func (m *MockCatalogService) CreateCatalog(ctx context.Context, catalog *models.Catalog) (*models.Catalog, error) {
	args := m.Called(ctx, catalog)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Catalog), args.Error(1)
}

func (m *MockCatalogService) UpdateCatalog(ctx context.Context, catalog *models.Catalog) (*models.Catalog, error) {
	args := m.Called(ctx, catalog)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Catalog), args.Error(1)
}

func (m *MockCatalogService) DeleteCatalog(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// TestCatalogHandler_CreateCatalog_Success тест успешного создания
func TestCatalogHandler_CreateCatalog_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	logger := zap.NewNop()
	mockService := new(MockCatalogService)

	expectedCatalog := &models.Catalog{
		ID:          "550e8400-e29b-41d4-a716-446655440000",
		Name:        "Test Catalog",
		Description: "Test Description",
		Price:       99.99,
		Currency:    "USD",
	}

	mockService.On("CreateCatalog", mock.Anything, mock.AnythingOfType("*models.Catalog")).Return(expectedCatalog, nil)

	h := handler.NewCatalogHandler(mockService, logger)
	router := gin.New()
	router.POST("/catalog", h.CreateCatalog)

	body, _ := json.Marshal(models.Catalog{
		Name:        "Test Catalog",
		Description: "Test Description",
		Price:       99.99,
		Currency:    "USD",
	})
	req := httptest.NewRequest("POST", "/catalog", bytes.NewReader(body))
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

// TestCatalogHandler_CreateCatalog_ValidationError тесты валидации
func TestCatalogHandler_CreateCatalog_ValidationError(t *testing.T) {
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
			requestBody: models.Catalog{
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
			mockService := new(MockCatalogService)
			h := handler.NewCatalogHandler(mockService, logger)
			router := gin.New()
			router.POST("/catalog", h.CreateCatalog)

			var body []byte
			switch v := tt.requestBody.(type) {
			case string:
				body = []byte(v)
			default:
				body, _ = json.Marshal(v)
			}

			req := httptest.NewRequest("POST", "/catalog", bytes.NewReader(body))
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

// TestCatalogHandler_CreateCatalog_GRPCError тест обработки gRPC ошибок
func TestCatalogHandler_CreateCatalog_GRPCError(t *testing.T) {
	tests := []struct {
		name           string
		grpcErr        error
		expectedStatus int
		expectedCode   string
	}{
		{
			name:           "already_exists",
			grpcErr:        status.Error(codes.AlreadyExists, "catalog already exists"),
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
			mockService := new(MockCatalogService)
			mockService.On("CreateCatalog", mock.Anything, mock.AnythingOfType("*models.Catalog")).Return(nil, tt.grpcErr)

			h := handler.NewCatalogHandler(mockService, logger)
			router := gin.New()
			router.POST("/catalog", h.CreateCatalog)

			body, _ := json.Marshal(models.Catalog{
				Name:        "Test Catalog",
				Description: "Test",
			})
			req := httptest.NewRequest("POST", "/catalog", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			var response models.APIResponse
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.False(t, response.Success)
			assert.Equal(t, tt.expectedCode, response.Error.Code)
			mockService.AssertExpectations(t)
		})
	}
}

// TestCatalogHandler_GetCatalog_Success тест успешного получения
func TestCatalogHandler_GetCatalog_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	logger := zap.NewNop()
	mockService := new(MockCatalogService)

	expectedCatalog := &models.Catalog{
		ID:          "550e8400-e29b-41d4-a716-446655440000",
		Name:        "Test Catalog",
		Description: "Test Description",
	}

	mockService.On("GetCatalog", mock.Anything, "550e8400-e29b-41d4-a716-446655440000").Return(expectedCatalog, nil)

	h := handler.NewCatalogHandler(mockService, logger)
	router := gin.New()
	router.GET("/catalog/:id", h.GetCatalog)

	req := httptest.NewRequest("GET", "/catalog/550e8400-e29b-41d4-a716-446655440000", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var response models.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response.Success)
	mockService.AssertExpectations(t)
}

// TestCatalogHandler_GetCatalog_NotFound тест когда каталог не найден
func TestCatalogHandler_GetCatalog_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	logger := zap.NewNop()
	mockService := new(MockCatalogService)

	mockService.On("GetCatalog", mock.Anything, "550e8400-e29b-41d4-a716-446655440001").Return(nil, status.Error(codes.NotFound, "catalog not found"))

	h := handler.NewCatalogHandler(mockService, logger)
	router := gin.New()
	router.GET("/catalog/:id", h.GetCatalog)

	req := httptest.NewRequest("GET", "/catalog/550e8400-e29b-41d4-a716-446655440001", nil)
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

// TestCatalogHandler_GetCatalog_InvalidUUID тест невалидного UUID
func TestCatalogHandler_GetCatalog_InvalidUUID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	logger := zap.NewNop()
	mockService := new(MockCatalogService)

	h := handler.NewCatalogHandler(mockService, logger)
	router := gin.New()
	router.GET("/catalog/:id", h.GetCatalog)

	req := httptest.NewRequest("GET", "/catalog/not-a-uuid", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	var response models.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.False(t, response.Success)
	assert.Equal(t, "INVALID_INPUT", response.Error.Code)
}

// TestCatalogHandler_UpdateCatalog_Success тест успешного обновления
func TestCatalogHandler_UpdateCatalog_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	logger := zap.NewNop()
	mockService := new(MockCatalogService)

	expectedCatalog := &models.Catalog{
		ID:          "550e8400-e29b-41d4-a716-446655440000",
		Name:        "Updated Catalog",
		Description: "Updated Description",
	}

	mockService.On("UpdateCatalog", mock.Anything, mock.AnythingOfType("*models.Catalog")).Return(expectedCatalog, nil)

	h := handler.NewCatalogHandler(mockService, logger)
	router := gin.New()
	router.PUT("/catalog/:id", h.UpdateCatalog)

	body, _ := json.Marshal(models.Catalog{
		Name:        "Updated Catalog",
		Description: "Updated Description",
	})
	req := httptest.NewRequest("PUT", "/catalog/550e8400-e29b-41d4-a716-446655440000", bytes.NewReader(body))
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

// TestCatalogHandler_DeleteCatalog_Success тест успешного удаления
func TestCatalogHandler_DeleteCatalog_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	logger := zap.NewNop()
	mockService := new(MockCatalogService)

	mockService.On("DeleteCatalog", mock.Anything, "550e8400-e29b-41d4-a716-446655440000").Return(nil)

	h := handler.NewCatalogHandler(mockService, logger)
	router := gin.New()
	router.DELETE("/catalog/:id", h.DeleteCatalog)

	req := httptest.NewRequest("DELETE", "/catalog/550e8400-e29b-41d4-a716-446655440000", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var response models.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response.Success)
	mockService.AssertExpectations(t)
}

// TestCatalogHandler_DeleteCatalog_NotFound тест удаления несуществующего
func TestCatalogHandler_DeleteCatalog_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	logger := zap.NewNop()
	mockService := new(MockCatalogService)

	mockService.On("DeleteCatalog", mock.Anything, "550e8400-e29b-41d4-a716-446655440001").Return(status.Error(codes.NotFound, "catalog not found"))

	h := handler.NewCatalogHandler(mockService, logger)
	router := gin.New()
	router.DELETE("/catalog/:id", h.DeleteCatalog)

	req := httptest.NewRequest("DELETE", "/catalog/550e8400-e29b-41d4-a716-446655440001", nil)
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

// TestCatalogHandler_DeleteCatalog_InvalidUUID тест удаления с невалидным UUID
func TestCatalogHandler_DeleteCatalog_InvalidUUID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	logger := zap.NewNop()
	mockService := new(MockCatalogService)

	h := handler.NewCatalogHandler(mockService, logger)
	router := gin.New()
	router.DELETE("/catalog/:id", h.DeleteCatalog)

	req := httptest.NewRequest("DELETE", "/catalog/not-a-uuid", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	var response models.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.False(t, response.Success)
	assert.Equal(t, "INVALID_INPUT", response.Error.Code)
}
