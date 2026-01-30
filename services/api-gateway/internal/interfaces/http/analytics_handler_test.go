// Файл api-gateway/internal/interfaces/http/analytics_handler_test.go содержит реализацию пакета http.
package http

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestAnalyticsHandler_GetDashboardStats(t *testing.T) {
	gin.SetMode(gin.TestMode)
	logger := zap.NewNop()

	handler := NewAnalyticsHandler(nil, logger)

	tests := []struct {
		name       string
		period     string
		startDate  string
		endDate    string
		wantStatus int
	}{
		{
			name:       "default period (month)",
			period:     "",
			wantStatus: http.StatusOK,
		},
		{
			name:       "today period",
			period:     "today",
			wantStatus: http.StatusOK,
		},
		{
			name:       "week period",
			period:     "week",
			wantStatus: http.StatusOK,
		},
		{
			name:       "custom period",
			period:     "custom",
			startDate:  "2024-01-01",
			endDate:    "2024-01-31",
			wantStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			req, _ := http.NewRequest(http.MethodGet, "/api/analytics/stats", nil)
			q := req.URL.Query()
			if tt.period != "" {
				q.Add("period", tt.period)
			}
			if tt.startDate != "" {
				q.Add("startDate", tt.startDate)
			}
			if tt.endDate != "" {
				q.Add("endDate", tt.endDate)
			}
			req.URL.RawQuery = q.Encode()
			c.Request = req

			handler.GetDashboardStats(c)

			assert.Equal(t, tt.wantStatus, w.Code)

			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.True(t, response["success"].(bool))

			data := response["data"].(map[string]interface{})
			assert.NotNil(t, data["totalRevenue"])
			assert.NotNil(t, data["totalInvoices"])
			assert.NotNil(t, data["period"])
		})
	}
}

func TestAnalyticsHandler_GetSalesChart(t *testing.T) {
	gin.SetMode(gin.TestMode)
	logger := zap.NewNop()

	handler := NewAnalyticsHandler(nil, logger)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	req, _ := http.NewRequest(http.MethodGet, "/api/analytics/sales-chart?period=week", nil)
	c.Request = req

	handler.GetSalesChart(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response["success"].(bool))

	data := response["data"].(map[string]interface{})
	assert.NotNil(t, data["data"])
	assert.NotNil(t, data["period"])
}

func TestAnalyticsHandler_GetStatusStats(t *testing.T) {
	gin.SetMode(gin.TestMode)
	logger := zap.NewNop()

	handler := NewAnalyticsHandler(nil, logger)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	req, _ := http.NewRequest(http.MethodGet, "/api/analytics/status-stats", nil)
	c.Request = req

	handler.GetStatusStats(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response["success"].(bool))

	data := response["data"].(map[string]interface{})
	chartData := data["data"].([]interface{})
	assert.Greater(t, len(chartData), 0)
}

func TestAnalyticsHandler_GetOperationTypeStats(t *testing.T) {
	gin.SetMode(gin.TestMode)
	logger := zap.NewNop()

	handler := NewAnalyticsHandler(nil, logger)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	req, _ := http.NewRequest(http.MethodGet, "/api/analytics/operation-stats", nil)
	c.Request = req

	handler.GetOperationTypeStats(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response["success"].(bool))
}

func TestAnalyticsHandler_GetTopContractors(t *testing.T) {
	gin.SetMode(gin.TestMode)
	logger := zap.NewNop()

	handler := NewAnalyticsHandler(nil, logger)

	tests := []struct {
		name       string
		limit      string
		wantStatus int
	}{
		{
			name:       "default limit",
			limit:      "",
			wantStatus: http.StatusOK,
		},
		{
			name:       "custom limit",
			limit:      "5",
			wantStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			req, _ := http.NewRequest(http.MethodGet, "/api/analytics/top-contractors", nil)
			if tt.limit != "" {
				q := req.URL.Query()
				q.Add("limit", tt.limit)
				req.URL.RawQuery = q.Encode()
			}
			c.Request = req

			handler.GetTopContractors(c)

			assert.Equal(t, tt.wantStatus, w.Code)

			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.True(t, response["success"].(bool))

			contractors := response["data"].([]interface{})
			assert.Greater(t, len(contractors), 0)
		})
	}
}

func TestAnalyticsHandler_GetRevenueByMonth(t *testing.T) {
	gin.SetMode(gin.TestMode)
	logger := zap.NewNop()

	handler := NewAnalyticsHandler(nil, logger)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	req, _ := http.NewRequest(http.MethodGet, "/api/analytics/revenue-by-month", nil)
	c.Request = req

	handler.GetRevenueByMonth(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response["success"].(bool))

	data := response["data"].([]interface{})
	assert.Greater(t, len(data), 0)
}

func TestAnalyticsHandler_calculateDateRange(t *testing.T) {
	logger := zap.NewNop()
	handler := NewAnalyticsHandler(nil, logger)

	tests := []struct {
		name      string
		period    string
		startDate string
		endDate   string
	}{
		{
			name:   "today period",
			period: "today",
		},
		{
			name:   "week period",
			period: "week",
		},
		{
			name:   "month period",
			period: "month",
		},
		{
			name:   "quarter period",
			period: "quarter",
		},
		{
			name:   "year period",
			period: "year",
		},
		{
			name:      "custom period",
			period:    "custom",
			startDate: "2024-01-01",
			endDate:   "2024-01-31",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			start, end := handler.calculateDateRange(tt.period, tt.startDate, tt.endDate)

			assert.False(t, start.IsZero())
			assert.False(t, end.IsZero())
			assert.True(t, start.Before(end) || start.Equal(end))
		})
	}
}
