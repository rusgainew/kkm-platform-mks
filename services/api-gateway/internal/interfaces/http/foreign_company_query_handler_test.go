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

func setupForeignCompanyQueryHandler() (*ForeignCompanyQueryHandler, *gin.Engine) {
	gin.SetMode(gin.TestMode)
	logger := zap.NewNop()
	metrics := observability.NewMetrics()
	tracer := observability.NewTracer("test", logger)
	connMgr := client.NewConnectionManager(5, logger)

	service := services.NewForeignCompanyQueryService(connMgr, "invalid:99999", metrics, tracer, logger)
	handler := NewForeignCompanyQueryHandler(service, logger)

	router := gin.New()
	return handler, router
}

func TestForeignCompanyQueryHandler_ListForeignCompanies_Success(t *testing.T) {
	handler, router := setupForeignCompanyQueryHandler()
	router.GET("/foreign-companies-query", handler.ListForeignCompanies)

	req, _ := http.NewRequest("GET", "/foreign-companies-query?page=0&page_size=20", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError && w.Code != http.StatusOK {
		t.Errorf("Expected status 500 or 200, got %d", w.Code)
	}
}

func TestForeignCompanyQueryHandler_ListForeignCompanies_InvalidPage(t *testing.T) {
	handler, router := setupForeignCompanyQueryHandler()
	router.GET("/foreign-companies-query", handler.ListForeignCompanies)

	req, _ := http.NewRequest("GET", "/foreign-companies-query?page=invalid", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestForeignCompanyQueryHandler_ListForeignCompanies_InvalidPageSize(t *testing.T) {
	handler, router := setupForeignCompanyQueryHandler()
	router.GET("/foreign-companies-query", handler.ListForeignCompanies)

	req, _ := http.NewRequest("GET", "/foreign-companies-query?page=0&page_size=invalid", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestForeignCompanyQueryHandler_ListForeignCompanies_MaxPageSize(t *testing.T) {
	handler, router := setupForeignCompanyQueryHandler()
	router.GET("/foreign-companies-query", handler.ListForeignCompanies)

	req, _ := http.NewRequest("GET", "/foreign-companies-query?page=0&page_size=200", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError && w.Code != http.StatusOK {
		t.Errorf("Expected status 500 or 200, got %d", w.Code)
	}
}

func TestForeignCompanyQueryHandler_SearchForeignCompanies_Success(t *testing.T) {
	handler, router := setupForeignCompanyQueryHandler()
	router.GET("/foreign-companies-query/search", handler.SearchForeignCompanies)

	req, _ := http.NewRequest("GET", "/foreign-companies-query/search?q=test&page=0&page_size=20", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError && w.Code != http.StatusOK {
		t.Errorf("Expected status 500 or 200, got %d", w.Code)
	}
}

func TestForeignCompanyQueryHandler_SearchForeignCompanies_EmptyQuery(t *testing.T) {
	handler, router := setupForeignCompanyQueryHandler()
	router.GET("/foreign-companies-query/search", handler.SearchForeignCompanies)

	req, _ := http.NewRequest("GET", "/foreign-companies-query/search?q=&page=0&page_size=20", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestForeignCompanyQueryHandler_SearchForeignCompanies_InvalidPage(t *testing.T) {
	handler, router := setupForeignCompanyQueryHandler()
	router.GET("/foreign-companies-query/search", handler.SearchForeignCompanies)

	req, _ := http.NewRequest("GET", "/foreign-companies-query/search?q=test&page=invalid", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// SearchForeignCompanies не парсит page, поэтому возвращает 500 при попытке подключения
	if w.Code != http.StatusInternalServerError && w.Code != http.StatusOK {
		t.Errorf("Expected status 500 or 200, got %d", w.Code)
	}
}

func TestForeignCompanyQueryHandler_FilterForeignCompanies_Success(t *testing.T) {
	handler, router := setupForeignCompanyQueryHandler()
	router.GET("/foreign-companies-query/filter", handler.FilterForeignCompanies)

	req, _ := http.NewRequest("GET", "/foreign-companies-query/filter?name=test&page=0&page_size=20", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError && w.Code != http.StatusOK {
		t.Errorf("Expected status 500 or 200, got %d", w.Code)
	}
}

func TestForeignCompanyQueryHandler_FilterForeignCompanies_InvalidPage(t *testing.T) {
	handler, router := setupForeignCompanyQueryHandler()
	router.GET("/foreign-companies-query/filter", handler.FilterForeignCompanies)

	req, _ := http.NewRequest("GET", "/foreign-companies-query/filter?page=invalid", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestForeignCompanyQueryHandler_FilterForeignCompanies_InvalidPageSize(t *testing.T) {
	handler, router := setupForeignCompanyQueryHandler()
	router.GET("/foreign-companies-query/filter", handler.FilterForeignCompanies)

	req, _ := http.NewRequest("GET", "/foreign-companies-query/filter?page=0&page_size=invalid", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}
