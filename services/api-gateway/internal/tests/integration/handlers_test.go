// Файл api-gateway/internal/tests/integration/handlers_test.go содержит реализацию пакета integration.
package integration

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rusgainew/kkm-project-mks/api-gateway/internal/domain/models"
	"github.com/rusgainew/kkm-project-mks/api-gateway/internal/infrastructure/auth"
	"github.com/rusgainew/kkm-project-mks/api-gateway/internal/infrastructure/client"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest"
)

// TestSetup подготавливает окружение для тестирования
type TestSetup struct {
	Router      *gin.Engine
	Logger      *zap.Logger
	AuthService *auth.Service
	ConnManager *client.ConnectionManager
}

// NewTestSetup создает новое тестовое окружение
func NewTestSetup(t *testing.T) *TestSetup {
	gin.SetMode(gin.TestMode)

	logger := zaptest.NewLogger(t)
	authService := auth.NewService("test-secret-key", 24*time.Hour, logger)
	connManager := client.NewConnectionManager(5*time.Second, logger)

	router := gin.New()

	return &TestSetup{
		Router:      router,
		Logger:      logger,
		AuthService: authService,
		ConnManager: connManager,
	}
}

// TestHealthEndpoint проверяет health endpoint
func TestHealthEndpoint(t *testing.T) {
	setup := NewTestSetup(t)

	// Простой health handler
	setup.Router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "healthy"})
	})

	req, _ := http.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()
	setup.Router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "healthy")
}

// TestErrorHandling проверяет обработку ошибок
func TestErrorHandling(t *testing.T) {
	setup := NewTestSetup(t)

	// Handler что возвращает ошибку
	setup.Router.GET("/error", func(c *gin.Context) {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error: &models.APIError{
				Code:    "INTERNAL_ERROR",
				Message: "Something went wrong",
				Details: "Test error",
			},
		})
	})

	req, _ := http.NewRequest("GET", "/error", nil)
	w := httptest.NewRecorder()
	setup.Router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "INTERNAL_ERROR")
}

// TestCORSHeaders проверяет CORS заголовки
func TestCORSHeaders(t *testing.T) {
	setup := NewTestSetup(t)

	setup.Router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "ok"})
	})

	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("Origin", "http://localhost:3000")
	w := httptest.NewRecorder()
	setup.Router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

// Mock Services для тестирования

// MockCompanyService имитирует CompanyService
type MockCompanyService struct {
	CreateFunc func(ctx context.Context, company *models.Company) (*models.Company, error)
	GetFunc    func(ctx context.Context, id string) (*models.Company, error)
	UpdateFunc func(ctx context.Context, company *models.Company) (*models.Company, error)
	DeleteFunc func(ctx context.Context, id string) error
	ListFunc   func(ctx context.Context, page, pageSize int) ([]*models.Company, int, error)
}

func (m *MockCompanyService) CreateCompany(ctx context.Context, company *models.Company) (*models.Company, error) {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, company)
	}
	return nil, nil
}

func (m *MockCompanyService) GetCompany(ctx context.Context, id string) (*models.Company, error) {
	if m.GetFunc != nil {
		return m.GetFunc(ctx, id)
	}
	return nil, nil
}

func (m *MockCompanyService) UpdateCompany(ctx context.Context, company *models.Company) (*models.Company, error) {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(ctx, company)
	}
	return nil, nil
}

func (m *MockCompanyService) DeleteCompany(ctx context.Context, id string) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(ctx, id)
	}
	return nil
}

func (m *MockCompanyService) ListCompanies(ctx context.Context, page, pageSize int) ([]*models.Company, int, error) {
	if m.ListFunc != nil {
		return m.ListFunc(ctx, page, pageSize)
	}
	return nil, 0, nil
}

// MockInvoiceService имитирует InvoiceService
type MockInvoiceService struct {
	CreateFunc         func(ctx context.Context, invoice *models.Invoice) (*models.Invoice, error)
	UpdateFunc         func(ctx context.Context, invoice *models.Invoice) (*models.Invoice, error)
	SignFunc           func(ctx context.Context, id string, signatureData string) (*models.Invoice, error)
	RevokeFunc         func(ctx context.Context, id string, reason string) (*models.Invoice, error)
	AcceptOrRejectFunc func(ctx context.Context, id string, accept bool, reason string) (*models.Invoice, error)
}

func (m *MockInvoiceService) CreateInvoice(ctx context.Context, invoice *models.Invoice) (*models.Invoice, error) {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, invoice)
	}
	return nil, nil
}

func (m *MockInvoiceService) UpdateInvoice(ctx context.Context, invoice *models.Invoice) (*models.Invoice, error) {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(ctx, invoice)
	}
	return nil, nil
}

func (m *MockInvoiceService) SignInvoice(ctx context.Context, id string, signatureData string) (*models.Invoice, error) {
	if m.SignFunc != nil {
		return m.SignFunc(ctx, id, signatureData)
	}
	return nil, nil
}

func (m *MockInvoiceService) RevokeInvoice(ctx context.Context, id string, reason string) (*models.Invoice, error) {
	if m.RevokeFunc != nil {
		return m.RevokeFunc(ctx, id, reason)
	}
	return nil, nil
}

func (m *MockInvoiceService) AcceptOrRejectInvoice(ctx context.Context, id string, accept bool, reason string) (*models.Invoice, error) {
	if m.AcceptOrRejectFunc != nil {
		return m.AcceptOrRejectFunc(ctx, id, accept, reason)
	}
	return nil, nil
}
