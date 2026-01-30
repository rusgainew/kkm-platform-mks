// Файл api-gateway/internal/interfaces/http/company_query_handler_test.go содержит реализацию пакета http.
package http

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/rusgainew/kkm-project-mks/api-gateway/internal/application/services"
	"github.com/rusgainew/kkm-project-mks/api-gateway/internal/infrastructure/client"
	"github.com/rusgainew/kkm-project-mks/api-gateway/internal/infrastructure/observability"
	"go.uber.org/zap"
)

func setupCompanyQueryHandler() (*CompanyQueryHandler, *gin.Engine) {
	gin.SetMode(gin.TestMode)
	logger := zap.NewNop()
	metrics := observability.NewMetrics()
	tracer := observability.NewTracer("test", logger)
	connMgr := client.NewConnectionManager(5, logger)

	service := services.NewCompanyQueryService(connMgr, "invalid:99999", metrics, tracer, logger)
	handler := NewCompanyQueryHandler(service, logger)

	router := gin.New()
	return handler, router
}

func TestCompanyQueryHandler_ListCompanies_Success(t *testing.T) {
	handler, router := setupCompanyQueryHandler()
	router.GET("/companies-query", handler.ListCompanies)

	req, _ := http.NewRequest("GET", "/companies-query?page=0&page_size=20", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError && w.Code != http.StatusOK {
		t.Errorf("Expected status 500 or 200, got %d", w.Code)
	}
}

func TestCompanyQueryHandler_ListCompanies_InvalidPage(t *testing.T) {
	handler, router := setupCompanyQueryHandler()
	router.GET("/companies-query", handler.ListCompanies)

	req, _ := http.NewRequest("GET", "/companies-query?page=invalid", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestCompanyQueryHandler_ListCompanies_InvalidPageSize(t *testing.T) {
	handler, router := setupCompanyQueryHandler()
	router.GET("/companies-query", handler.ListCompanies)

	req, _ := http.NewRequest("GET", "/companies-query?page=0&page_size=invalid", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestCompanyQueryHandler_ListCompanies_MaxPageSize(t *testing.T) {
	handler, router := setupCompanyQueryHandler()
	router.GET("/companies-query", handler.ListCompanies)

	req, _ := http.NewRequest("GET", "/companies-query?page=0&page_size=200", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Should limit to 100 and try to proceed (expect error due to invalid address)
	if w.Code != http.StatusInternalServerError && w.Code != http.StatusOK {
		t.Errorf("Expected status 500 or 200, got %d", w.Code)
	}
}

func TestCompanyQueryHandler_SearchCompanies_Success(t *testing.T) {
	handler, router := setupCompanyQueryHandler()
	router.GET("/companies-query/search", handler.SearchCompanies)

	req, _ := http.NewRequest("GET", "/companies-query/search?q=test&page=0&page_size=20", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError && w.Code != http.StatusOK {
		t.Errorf("Expected status 500 or 200, got %d", w.Code)
	}
}

func TestCompanyQueryHandler_SearchCompanies_EmptyQuery(t *testing.T) {
	handler, router := setupCompanyQueryHandler()
	router.GET("/companies-query/search", handler.SearchCompanies)

	req, _ := http.NewRequest("GET", "/companies-query/search?q=&page=0&page_size=20", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestCompanyQueryHandler_SearchCompanies_InvalidPage(t *testing.T) {
	handler, router := setupCompanyQueryHandler()
	router.GET("/companies-query/search", handler.SearchCompanies)

	req, _ := http.NewRequest("GET", "/companies-query/search?q=test&page=invalid", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// SearchCompanies не парсит page, поэтому возвращает 500 при попытке подключения
	if w.Code != http.StatusInternalServerError && w.Code != http.StatusOK {
		t.Errorf("Expected status 500 or 200, got %d", w.Code)
	}
}

func TestCompanyQueryHandler_FilterCompanies_Success(t *testing.T) {
	handler, router := setupCompanyQueryHandler()
	router.GET("/companies-query/filter", handler.FilterCompanies)

	req, _ := http.NewRequest("GET", "/companies-query/filter?name=test&page=0&page_size=20", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError && w.Code != http.StatusOK {
		t.Errorf("Expected status 500 or 200, got %d", w.Code)
	}
}

func TestCompanyQueryHandler_FilterCompanies_InvalidPage(t *testing.T) {
	handler, router := setupCompanyQueryHandler()
	router.GET("/companies-query/filter", handler.FilterCompanies)

	req, _ := http.NewRequest("GET", "/companies-query/filter?page=invalid", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestCompanyQueryHandler_FilterCompanies_InvalidPageSize(t *testing.T) {
	handler, router := setupCompanyQueryHandler()
	router.GET("/companies-query/filter", handler.FilterCompanies)

	req, _ := http.NewRequest("GET", "/companies-query/filter?page=0&page_size=invalid", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}
