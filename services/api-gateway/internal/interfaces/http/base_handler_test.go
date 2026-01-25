package http_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	handler "github.com/rusgainew/kkm-project-mks/api-gateway/internal/interfaces/http"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestBaseHandler_ValidateUUID(t *testing.T) {
	tests := []struct {
		name     string
		id       string
		wantOk   bool
		wantCode int
	}{
		{
			name:     "valid_uuid_v4",
			id:       "550e8400-e29b-41d4-a716-446655440000",
			wantOk:   true,
			wantCode: 0,
		},
		{
			name:     "valid_uuid_lowercase",
			id:       "550e8400-e29b-41d4-a716-446655440000",
			wantOk:   true,
			wantCode: 0,
		},
		{
			name:     "empty_string",
			id:       "",
			wantOk:   false,
			wantCode: http.StatusBadRequest,
		},
		{
			name:     "invalid_format",
			id:       "not-a-uuid",
			wantOk:   false,
			wantCode: http.StatusBadRequest,
		},
		{
			name:     "too_short",
			id:       "550e8400",
			wantOk:   false,
			wantCode: http.StatusBadRequest,
		},
		{
			name:     "missing_hyphens",
			id:       "550e8400e29b41d4a716446655440000",
			wantOk:   false,
			wantCode: http.StatusBadRequest,
		},
		{
			name:     "invalid_version",
			id:       "550e8400-e29b-51d4-a716-446655440000", // v5 instead of v4
			wantOk:   false,
			wantCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)
			logger := zap.NewNop()
			h := handler.NewBaseHandler(logger)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			ok := h.ValidateUUID(c, tt.id, "test_id")

			assert.Equal(t, tt.wantOk, ok)
			if !tt.wantOk {
				assert.Equal(t, tt.wantCode, w.Code)
			}
		})
	}
}

func TestBaseHandler_BindJSON(t *testing.T) {
	type TestStruct struct {
		Name  string `json:"name" binding:"required"`
		Value int    `json:"value"`
	}

	tests := []struct {
		name       string
		body       string
		wantOk     bool
		wantStatus int
	}{
		{
			name:       "valid_json",
			body:       `{"name": "test", "value": 123}`,
			wantOk:     true,
			wantStatus: 0,
		},
		{
			name:       "empty_body",
			body:       "",
			wantOk:     false,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "invalid_json",
			body:       "not json",
			wantOk:     false,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "missing_required_field",
			body:       `{"value": 123}`,
			wantOk:     false,
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)
			logger := zap.NewNop()
			h := handler.NewBaseHandler(logger)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest("POST", "/test", strings.NewReader(tt.body))
			c.Request.Header.Set("Content-Type", "application/json")

			var obj TestStruct
			ok := h.BindJSON(c, &obj)

			assert.Equal(t, tt.wantOk, ok)
			if !tt.wantOk {
				assert.Equal(t, tt.wantStatus, w.Code)
			}
		})
	}
}

func TestBaseHandler_GetPagination(t *testing.T) {
	tests := []struct {
		name         string
		queryParams  string
		wantPage     int
		wantPageSize int
		wantOk       bool
	}{
		{
			name:         "default_values",
			queryParams:  "",
			wantPage:     1,
			wantPageSize: 20,
			wantOk:       true,
		},
		{
			name:         "custom_values",
			queryParams:  "page=5&page_size=50",
			wantPage:     5,
			wantPageSize: 50,
			wantOk:       true,
		},
		{
			name:         "exceeds_max_page_size",
			queryParams:  "page_size=200",
			wantPage:     1,
			wantPageSize: 100, // capped at 100
			wantOk:       true,
		},
		{
			name:         "invalid_page_not_number",
			queryParams:  "page=abc",
			wantPage:     0,
			wantPageSize: 0,
			wantOk:       false,
		},
		{
			name:         "invalid_page_size_not_number",
			queryParams:  "page_size=xyz",
			wantPage:     0,
			wantPageSize: 0,
			wantOk:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)
			logger := zap.NewNop()
			h := handler.NewBaseHandler(logger)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest("GET", "/test?"+tt.queryParams, nil)

			page, pageSize, ok := h.GetPagination(c)

			assert.Equal(t, tt.wantOk, ok)
			if tt.wantOk {
				assert.Equal(t, tt.wantPage, page)
				assert.Equal(t, tt.wantPageSize, pageSize)
			}
		})
	}
}

func TestBaseHandler_GetIntPathParam(t *testing.T) {
	tests := []struct {
		name       string
		paramValue string
		wantValue  int64
		wantOk     bool
	}{
		{
			name:       "valid_positive_int",
			paramValue: "123",
			wantValue:  123,
			wantOk:     true,
		},
		{
			name:       "valid_large_int",
			paramValue: "9999999",
			wantValue:  9999999,
			wantOk:     true,
		},
		{
			name:       "zero",
			paramValue: "0",
			wantValue:  0,
			wantOk:     false,
		},
		{
			name:       "negative",
			paramValue: "-5",
			wantValue:  0,
			wantOk:     false,
		},
		{
			name:       "not_a_number",
			paramValue: "abc",
			wantValue:  0,
			wantOk:     false,
		},
		{
			name:       "empty",
			paramValue: "",
			wantValue:  0,
			wantOk:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)
			logger := zap.NewNop()
			h := handler.NewBaseHandler(logger)

			router := gin.New()
			var resultValue int64
			var resultOk bool

			router.GET("/test/:id", func(c *gin.Context) {
				resultValue, resultOk = h.GetIntPathParam(c, "id", "test_id")
			})

			path := "/test/" + tt.paramValue
			if tt.paramValue == "" {
				path = "/test/"
			}

			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", path, nil)
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.wantOk, resultOk)
			if tt.wantOk {
				assert.Equal(t, tt.wantValue, resultValue)
			}
		})
	}
}
