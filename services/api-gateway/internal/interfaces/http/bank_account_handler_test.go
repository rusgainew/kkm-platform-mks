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

// MockBankAccountService мок для BankAccountServiceInterface
type MockBankAccountService struct {
	mock.Mock
}

func (m *MockBankAccountService) GetBankAccount(ctx context.Context, id string) (*models.BankAccount, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.BankAccount), args.Error(1)
}

func (m *MockBankAccountService) CreateBankAccount(ctx context.Context, account *models.BankAccount) (*models.BankAccount, error) {
	args := m.Called(ctx, account)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.BankAccount), args.Error(1)
}

func (m *MockBankAccountService) UpdateBankAccount(ctx context.Context, account *models.BankAccount) (*models.BankAccount, error) {
	args := m.Called(ctx, account)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.BankAccount), args.Error(1)
}

func (m *MockBankAccountService) DeleteBankAccount(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// TestBankAccountHandler_CreateBankAccount_Success тест успешного создания
func TestBankAccountHandler_CreateBankAccount_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	logger := zap.NewNop()
	mockService := new(MockBankAccountService)

	expectedAccount := &models.BankAccount{
		ID:            "550e8400-e29b-41d4-a716-446655440000",
		AccountNumber: "40817810099910004321",
		BankName:      "Test Bank",
		BankCode:      "044525225",
		Currency:      "RUB",
		OwnerID:       "550e8400-e29b-41d4-a716-446655440001",
	}

	mockService.On("CreateBankAccount", mock.Anything, mock.AnythingOfType("*models.BankAccount")).Return(expectedAccount, nil)

	h := handler.NewBankAccountHandler(mockService, logger)
	router := gin.New()
	router.POST("/bank-accounts", h.CreateBankAccount)

	body, _ := json.Marshal(models.BankAccount{
		AccountNumber: "40817810099910004321",
		BankName:      "Test Bank",
		BankCode:      "044525225",
		Currency:      "RUB",
		OwnerID:       "550e8400-e29b-41d4-a716-446655440001",
	})
	req := httptest.NewRequest("POST", "/bank-accounts", bytes.NewReader(body))
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

// TestBankAccountHandler_CreateBankAccount_ValidationError тесты валидации
func TestBankAccountHandler_CreateBankAccount_ValidationError(t *testing.T) {
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
			name: "empty_account_number",
			requestBody: models.BankAccount{
				AccountNumber: "",
				BankName:      "Test Bank",
				BankCode:      "044525225",
			},
			expectedStatus: http.StatusBadRequest,
			expectedCode:   "VALIDATION_FAILED",
		},
		{
			name: "empty_bank_name",
			requestBody: models.BankAccount{
				AccountNumber: "40817810099910004321",
				BankName:      "",
				BankCode:      "044525225",
			},
			expectedStatus: http.StatusBadRequest,
			expectedCode:   "VALIDATION_FAILED",
		},
		{
			name: "empty_bank_code",
			requestBody: models.BankAccount{
				AccountNumber: "40817810099910004321",
				BankName:      "Test Bank",
				BankCode:      "",
			},
			expectedStatus: http.StatusBadRequest,
			expectedCode:   "VALIDATION_FAILED",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)
			logger := zap.NewNop()
			mockService := new(MockBankAccountService)
			h := handler.NewBankAccountHandler(mockService, logger)
			router := gin.New()
			router.POST("/bank-accounts", h.CreateBankAccount)

			var body []byte
			switch v := tt.requestBody.(type) {
			case string:
				body = []byte(v)
			default:
				body, _ = json.Marshal(v)
			}

			req := httptest.NewRequest("POST", "/bank-accounts", bytes.NewReader(body))
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

// TestBankAccountHandler_CreateBankAccount_GRPCError тест обработки gRPC ошибок
func TestBankAccountHandler_CreateBankAccount_GRPCError(t *testing.T) {
	tests := []struct {
		name           string
		grpcErr        error
		expectedStatus int
		expectedCode   string
	}{
		{
			name:           "already_exists",
			grpcErr:        status.Error(codes.AlreadyExists, "bank account already exists"),
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
			mockService := new(MockBankAccountService)
			mockService.On("CreateBankAccount", mock.Anything, mock.AnythingOfType("*models.BankAccount")).Return(nil, tt.grpcErr)

			h := handler.NewBankAccountHandler(mockService, logger)
			router := gin.New()
			router.POST("/bank-accounts", h.CreateBankAccount)

			body, _ := json.Marshal(models.BankAccount{
				AccountNumber: "40817810099910004321",
				BankName:      "Test Bank",
				BankCode:      "044525225",
			})
			req := httptest.NewRequest("POST", "/bank-accounts", bytes.NewReader(body))
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

// TestBankAccountHandler_GetBankAccount_Success тест успешного получения
func TestBankAccountHandler_GetBankAccount_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	logger := zap.NewNop()
	mockService := new(MockBankAccountService)

	expectedAccount := &models.BankAccount{
		ID:            "550e8400-e29b-41d4-a716-446655440000",
		AccountNumber: "40817810099910004321",
		BankName:      "Test Bank",
		BankCode:      "044525225",
	}

	mockService.On("GetBankAccount", mock.Anything, "550e8400-e29b-41d4-a716-446655440000").Return(expectedAccount, nil)

	h := handler.NewBankAccountHandler(mockService, logger)
	router := gin.New()
	router.GET("/bank-accounts/:id", h.GetBankAccount)

	req := httptest.NewRequest("GET", "/bank-accounts/550e8400-e29b-41d4-a716-446655440000", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var response models.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response.Success)
	mockService.AssertExpectations(t)
}

// TestBankAccountHandler_GetBankAccount_NotFound тест когда счёт не найден
func TestBankAccountHandler_GetBankAccount_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	logger := zap.NewNop()
	mockService := new(MockBankAccountService)

	mockService.On("GetBankAccount", mock.Anything, "550e8400-e29b-41d4-a716-446655440001").Return(nil, status.Error(codes.NotFound, "bank account not found"))

	h := handler.NewBankAccountHandler(mockService, logger)
	router := gin.New()
	router.GET("/bank-accounts/:id", h.GetBankAccount)

	req := httptest.NewRequest("GET", "/bank-accounts/550e8400-e29b-41d4-a716-446655440001", nil)
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

// TestBankAccountHandler_GetBankAccount_InvalidUUID тест невалидного UUID
func TestBankAccountHandler_GetBankAccount_InvalidUUID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	logger := zap.NewNop()
	mockService := new(MockBankAccountService)

	h := handler.NewBankAccountHandler(mockService, logger)
	router := gin.New()
	router.GET("/bank-accounts/:id", h.GetBankAccount)

	req := httptest.NewRequest("GET", "/bank-accounts/not-a-uuid", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	var response models.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.False(t, response.Success)
	assert.Equal(t, "INVALID_INPUT", response.Error.Code)
}

// TestBankAccountHandler_UpdateBankAccount_Success тест успешного обновления
func TestBankAccountHandler_UpdateBankAccount_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	logger := zap.NewNop()
	mockService := new(MockBankAccountService)

	expectedAccount := &models.BankAccount{
		ID:            "550e8400-e29b-41d4-a716-446655440000",
		AccountNumber: "40817810099910004322",
		BankName:      "Updated Bank",
		BankCode:      "044525999",
	}

	mockService.On("UpdateBankAccount", mock.Anything, mock.AnythingOfType("*models.BankAccount")).Return(expectedAccount, nil)

	h := handler.NewBankAccountHandler(mockService, logger)
	router := gin.New()
	router.PUT("/bank-accounts/:id", h.UpdateBankAccount)

	body, _ := json.Marshal(models.BankAccount{
		AccountNumber: "40817810099910004322",
		BankName:      "Updated Bank",
		BankCode:      "044525999",
	})
	req := httptest.NewRequest("PUT", "/bank-accounts/550e8400-e29b-41d4-a716-446655440000", bytes.NewReader(body))
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

// TestBankAccountHandler_DeleteBankAccount_Success тест успешного удаления
func TestBankAccountHandler_DeleteBankAccount_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	logger := zap.NewNop()
	mockService := new(MockBankAccountService)

	mockService.On("DeleteBankAccount", mock.Anything, "550e8400-e29b-41d4-a716-446655440000").Return(nil)

	h := handler.NewBankAccountHandler(mockService, logger)
	router := gin.New()
	router.DELETE("/bank-accounts/:id", h.DeleteBankAccount)

	req := httptest.NewRequest("DELETE", "/bank-accounts/550e8400-e29b-41d4-a716-446655440000", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var response models.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response.Success)
	mockService.AssertExpectations(t)
}

// TestBankAccountHandler_DeleteBankAccount_NotFound тест удаления несуществующего
func TestBankAccountHandler_DeleteBankAccount_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	logger := zap.NewNop()
	mockService := new(MockBankAccountService)

	mockService.On("DeleteBankAccount", mock.Anything, "550e8400-e29b-41d4-a716-446655440001").Return(status.Error(codes.NotFound, "bank account not found"))

	h := handler.NewBankAccountHandler(mockService, logger)
	router := gin.New()
	router.DELETE("/bank-accounts/:id", h.DeleteBankAccount)

	req := httptest.NewRequest("DELETE", "/bank-accounts/550e8400-e29b-41d4-a716-446655440001", nil)
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

// TestBankAccountHandler_DeleteBankAccount_InvalidUUID тест удаления с невалидным UUID
func TestBankAccountHandler_DeleteBankAccount_InvalidUUID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	logger := zap.NewNop()
	mockService := new(MockBankAccountService)

	h := handler.NewBankAccountHandler(mockService, logger)
	router := gin.New()
	router.DELETE("/bank-accounts/:id", h.DeleteBankAccount)

	req := httptest.NewRequest("DELETE", "/bank-accounts/not-a-uuid", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	var response models.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.False(t, response.Success)
	assert.Equal(t, "INVALID_INPUT", response.Error.Code)
}
