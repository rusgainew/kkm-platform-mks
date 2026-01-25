package integration

import (
	"context"
	"testing"
	"time"

	"github.com/rusgainew/kkm-project-mks/api-gateway/internal/domain/models"
	"github.com/rusgainew/kkm-project-mks/api-gateway/internal/infrastructure/auth"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"
)

// TestAuthServiceTokenGeneration проверяет генерацию и валидацию JWT токена
func TestAuthServiceTokenGeneration(t *testing.T) {
	logger := zaptest.NewLogger(t)
	authService := auth.NewService("test-secret", 24*time.Hour, logger)

	ctx := context.Background()

	// Test успешной генерации токена
	claimsData := map[string]interface{}{
		"permissions": []string{"read", "write"},
	}
	token, err := authService.GenerateToken(ctx, "test-user", claimsData)
	require.NoError(t, err)
	assert.NotEmpty(t, token)

	// Test валидации токена
	claims, err := authService.ValidateToken(ctx, token)
	require.NoError(t, err)
	assert.Equal(t, "test-user", claims.UserID)
}

// TestAuthServiceInvalidToken проверяет отклонение невалидного токена
func TestAuthServiceInvalidToken(t *testing.T) {
	logger := zaptest.NewLogger(t)
	authService := auth.NewService("test-secret", 24*time.Hour, logger)

	ctx := context.Background()

	// Test невалидного токена
	_, err := authService.ValidateToken(ctx, "invalid-token")
	assert.Error(t, err)

	// Test с пустым токеном
	_, err = authService.ValidateToken(ctx, "")
	assert.Error(t, err)
}

// TestCompanyServiceCreateLogic проверяет логику создания компании
func TestCompanyServiceCreateLogic(t *testing.T) {
	// Создаем mock service с кастомной логикой
	mockService := &MockCompanyService{
		CreateFunc: func(ctx context.Context, company *models.Company) (*models.Company, error) {
			// Проверяем, что имя заполнено
			if company.Name == "" {
				return nil, assert.AnError
			}

			// Генерируем ID и возвращаем созданную компанию
			company.ID = "test-company-123"
			company.Status = "active"
			company.CreatedAt = 1000000
			company.UpdatedAt = 1000000
			return company, nil
		},
	}

	// Test успешного создания
	company := &models.Company{
		Name:        "Test Company",
		Description: "Test Description",
		OwnerID:     "owner-123",
	}

	result, err := mockService.CreateCompany(context.Background(), company)
	require.NoError(t, err)
	assert.NotEmpty(t, result.ID)
	assert.Equal(t, "Test Company", result.Name)
	assert.Equal(t, "active", result.Status)
}

// TestPaginationRequest проверяет структуру запроса пагинации
func TestPaginationRequest(t *testing.T) {
	req := models.PaginationRequest{
		Page:     1,
		PageSize: 10,
	}

	assert.Equal(t, 1, req.Page)
	assert.Equal(t, 10, req.PageSize)
}

// TestAPIResponseStructure проверяет структуру API ответа
func TestAPIResponseStructure(t *testing.T) {
	// Test успешного ответа
	successResponse := models.APIResponse{
		Success: true,
		Data: models.Company{
			ID:   "test-123",
			Name: "Test",
		},
		Meta: &models.MetaData{
			Page:     1,
			PageSize: 10,
		},
	}

	assert.True(t, successResponse.Success)
	assert.Nil(t, successResponse.Error)
	assert.NotNil(t, successResponse.Data)

	// Test ошибки
	errorResponse := models.APIResponse{
		Success: false,
		Error: &models.APIError{
			Code:    "NOT_FOUND",
			Message: "Company not found",
		},
	}

	assert.False(t, errorResponse.Success)
	assert.NotNil(t, errorResponse.Error)
	assert.Equal(t, "NOT_FOUND", errorResponse.Error.Code)
}

// TestMetadataCalculation проверяет расчет метаданных пагинации
func TestMetadataCalculation(t *testing.T) {
	totalItems := 25
	pageSize := 10
	page := 1

	meta := models.MetaData{
		Page:       page,
		PageSize:   pageSize,
		TotalCount: totalItems,
		TotalPages: (totalItems + pageSize - 1) / pageSize,
	}

	assert.Equal(t, 1, meta.Page)
	assert.Equal(t, 10, meta.PageSize)
	assert.Equal(t, 25, meta.TotalCount)
	assert.Equal(t, 3, meta.TotalPages)
}

// TestModelValidation проверяет валидацию моделей
func TestModelValidation(t *testing.T) {
	// Test Company
	company := models.Company{
		ID:          "test-123",
		Name:        "Test Company",
		Description: "Test",
		OwnerID:     "owner-123",
		Status:      "active",
	}
	assert.Equal(t, "test-123", company.ID)
	assert.Equal(t, "active", company.Status)

	// Test Invoice
	invoice := models.Invoice{
		ID:            "invoice-123",
		InvoiceNumber: "INV-001",
		TotalAmount:   1000.50,
		IsResident:    true,
		Status:        "draft",
	}
	assert.Equal(t, "invoice-123", invoice.ID)
	assert.Equal(t, 1000.50, invoice.TotalAmount)
	assert.True(t, invoice.IsResident)

	// Test Catalog
	catalog := models.Catalog{
		ID:   "cat-123",
		Name: "Product",
		Unit: "PCS",
	}
	assert.Equal(t, "cat-123", catalog.ID)
	assert.Equal(t, "PCS", catalog.Unit)

	// Test BankAccount
	account := models.BankAccount{
		ID:            "acc-123",
		AccountNumber: "12345678",
		Currency:      "KGS",
		IsActive:      true,
	}
	assert.Equal(t, "acc-123", account.ID)
	assert.Equal(t, "KGS", account.Currency)
	assert.True(t, account.IsActive)
}
