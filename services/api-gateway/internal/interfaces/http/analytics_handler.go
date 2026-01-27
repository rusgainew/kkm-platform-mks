package http

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rusgainew/kkm-project-mks/api-gateway/internal/application/services"
	"go.uber.org/zap"
)

// AnalyticsHandler обработчик для аналитики Dashboard
type AnalyticsHandler struct {
	invoiceQueryService *services.InvoiceQueryService
	logger              *zap.Logger
}

// NewAnalyticsHandler создает новый AnalyticsHandler
func NewAnalyticsHandler(
	invoiceQueryService *services.InvoiceQueryService,
	logger *zap.Logger,
) *AnalyticsHandler {
	return &AnalyticsHandler{
		invoiceQueryService: invoiceQueryService,
		logger:              logger,
	}
}

// DashboardStats ответ со статистикой для Dashboard
type DashboardStats struct {
	TotalRevenue         float64 `json:"totalRevenue"`
	TotalInvoices        int32   `json:"totalInvoices"`
	AverageInvoiceAmount float64 `json:"averageInvoiceAmount"`
	ActiveContractors    int32   `json:"activeContractors"`
	PendingInvoices      int32   `json:"pendingInvoices"`
	ApprovedInvoices     int32   `json:"approvedInvoices"`
	RejectedInvoices     int32   `json:"rejectedInvoices"`
	RevenueChange        float64 `json:"revenueChange,omitempty"`
	InvoiceCountChange   float64 `json:"invoiceCountChange,omitempty"`
	AverageAmountChange  float64 `json:"averageAmountChange,omitempty"`
	ContractorsChange    float64 `json:"contractorsChange,omitempty"`
	Period               Period  `json:"period"`
}

// Period временной период
type Period struct {
	StartDate string `json:"startDate"`
	EndDate   string `json:"endDate"`
}

// ChartDataPoint точка данных для графика
type ChartDataPoint struct {
	Date  string  `json:"date"`
	Value float64 `json:"value"`
}

// SalesChartData данные графика продаж
type SalesChartData struct {
	Data   []ChartDataPoint `json:"data"`
	Period Period           `json:"period"`
}

// PieChartDataPoint точка данных для круговой диаграммы
type PieChartDataPoint struct {
	Name       string  `json:"name"`
	Value      float64 `json:"value"`
	Percentage float64 `json:"percentage"`
	Color      string  `json:"color,omitempty"`
}

// PieChartData данные круговой диаграммы
type PieChartData struct {
	Data []PieChartDataPoint `json:"data"`
}

// TopContractor топ контрагент
type TopContractor struct {
	ContractorID   string  `json:"contractorId"`
	ContractorName string  `json:"contractorName"`
	TotalAmount    float64 `json:"totalAmount"`
	InvoiceCount   int32   `json:"invoiceCount"`
}

// GetDashboardStats получение статистики для Dashboard
//
//	@Summary		Get dashboard statistics
//	@Description	Get aggregated statistics for dashboard
//	@Tags			Analytics
//	@Accept			json
//	@Produce		json
//	@Param			period		query		string	false	"Time period: today, week, month, quarter, year, custom"
//	@Param			startDate	query		string	false	"Start date for custom period (YYYY-MM-DD)"
//	@Param			endDate		query		string	false	"End date for custom period (YYYY-MM-DD)"
//	@Success		200			{object}	DashboardStats
//	@Failure		400			{object}	map[string]string
//	@Failure		500			{object}	map[string]string
//	@Security		Bearer
//	@Router			/api/analytics/stats [get]
func (h *AnalyticsHandler) GetDashboardStats(c *gin.Context) {
	period := c.DefaultQuery("period", "month")
	startDate := c.Query("startDate")
	endDate := c.Query("endDate")

	// Вычисляем диапазон дат на основе периода
	start, end := h.calculateDateRange(period, startDate, endDate)

	h.logger.Info("Getting dashboard stats",
		zap.String("period", period),
		zap.Time("start", start),
		zap.Time("end", end),
	)

	// TODO: Implement actual database queries
	// For now, return mock data
	stats := DashboardStats{
		TotalRevenue:         1500000.50,
		TotalInvoices:        245,
		AverageInvoiceAmount: 6122.45,
		ActiveContractors:    38,
		PendingInvoices:      12,
		ApprovedInvoices:     220,
		RejectedInvoices:     13,
		RevenueChange:        15.7,
		InvoiceCountChange:   8.3,
		AverageAmountChange:  3.2,
		ContractorsChange:    5.5,
		Period: Period{
			StartDate: start.Format("2006-01-02"),
			EndDate:   end.Format("2006-01-02"),
		},
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    stats,
	})
}

// GetSalesChart получение данных графика продаж
//
//	@Summary		Get sales chart data
//	@Description	Get time series data for sales chart
//	@Tags			Analytics
//	@Accept			json
//	@Produce		json
//	@Param			period		query		string	false	"Time period"
//	@Param			startDate	query		string	false	"Start date (YYYY-MM-DD)"
//	@Param			endDate		query		string	false	"End date (YYYY-MM-DD)"
//	@Success		200			{object}	SalesChartData
//	@Failure		500			{object}	map[string]string
//	@Security		Bearer
//	@Router			/api/analytics/sales-chart [get]
func (h *AnalyticsHandler) GetSalesChart(c *gin.Context) {
	period := c.DefaultQuery("period", "month")
	startDate := c.Query("startDate")
	endDate := c.Query("endDate")

	start, end := h.calculateDateRange(period, startDate, endDate)

	h.logger.Info("Getting sales chart",
		zap.String("period", period),
		zap.Time("start", start),
		zap.Time("end", end),
	)

	// TODO: Implement actual database queries
	// Generate mock data points
	data := []ChartDataPoint{}
	days := int(end.Sub(start).Hours() / 24)
	for i := 0; i <= days; i++ {
		date := start.AddDate(0, 0, i)
		data = append(data, ChartDataPoint{
			Date:  date.Format("2006-01-02"),
			Value: 45000 + float64(i)*1500,
		})
	}

	response := SalesChartData{
		Data: data,
		Period: Period{
			StartDate: start.Format("2006-01-02"),
			EndDate:   end.Format("2006-01-02"),
		},
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    response,
	})
}

// GetStatusStats получение статистики по статусам накладных
//
//	@Summary		Get invoice status statistics
//	@Description	Get distribution of invoices by status
//	@Tags			Analytics
//	@Accept			json
//	@Produce		json
//	@Param			period		query		string	false	"Time period"
//	@Param			startDate	query		string	false	"Start date (YYYY-MM-DD)"
//	@Param			endDate		query		string	false	"End date (YYYY-MM-DD)"
//	@Success		200			{object}	PieChartData
//	@Failure		500			{object}	map[string]string
//	@Security		Bearer
//	@Router			/api/analytics/status-stats [get]
func (h *AnalyticsHandler) GetStatusStats(c *gin.Context) {
	period := c.DefaultQuery("period", "month")
	startDate := c.Query("startDate")
	endDate := c.Query("endDate")

	start, end := h.calculateDateRange(period, startDate, endDate)

	h.logger.Info("Getting status stats",
		zap.String("period", period),
		zap.Time("start", start),
		zap.Time("end", end),
	)

	// TODO: Implement actual database queries
	data := PieChartData{
		Data: []PieChartDataPoint{
			{Name: "Утверждено", Value: 220, Percentage: 89.8, Color: "#10b981"},
			{Name: "На рассмотрении", Value: 12, Percentage: 4.9, Color: "#f59e0b"},
			{Name: "Отклонено", Value: 13, Percentage: 5.3, Color: "#ef4444"},
		},
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    data,
	})
}

// GetOperationTypeStats получение статистики по типам операций
//
//	@Summary		Get operation type statistics
//	@Description	Get distribution of invoices by operation type
//	@Tags			Analytics
//	@Accept			json
//	@Produce		json
//	@Param			period		query		string	false	"Time period"
//	@Param			startDate	query		string	false	"Start date (YYYY-MM-DD)"
//	@Param			endDate		query		string	false	"End date (YYYY-MM-DD)"
//	@Success		200			{object}	PieChartData
//	@Failure		500			{object}	map[string]string
//	@Security		Bearer
//	@Router			/api/analytics/operation-stats [get]
func (h *AnalyticsHandler) GetOperationTypeStats(c *gin.Context) {
	period := c.DefaultQuery("period", "month")
	startDate := c.Query("startDate")
	endDate := c.Query("endDate")

	start, end := h.calculateDateRange(period, startDate, endDate)

	h.logger.Info("Getting operation type stats",
		zap.String("period", period),
		zap.Time("start", start),
		zap.Time("end", end),
	)

	// TODO: Implement actual database queries
	data := PieChartData{
		Data: []PieChartDataPoint{
			{Name: "Продажа", Value: 180, Percentage: 73.5, Color: "#3b82f6"},
			{Name: "Возврат", Value: 45, Percentage: 18.4, Color: "#8b5cf6"},
			{Name: "Корректировка", Value: 20, Percentage: 8.1, Color: "#ec4899"},
		},
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    data,
	})
}

// GetTopContractors получение топ контрагентов
//
//	@Summary		Get top contractors
//	@Description	Get list of top contractors by revenue
//	@Tags			Analytics
//	@Accept			json
//	@Produce		json
//	@Param			period		query		string	false	"Time period"
//	@Param			startDate	query		string	false	"Start date (YYYY-MM-DD)"
//	@Param			endDate		query		string	false	"End date (YYYY-MM-DD)"
//	@Param			limit		query		int		false	"Number of contractors to return"	default(10)
//	@Success		200			{array}		TopContractor
//	@Failure		500			{object}	map[string]string
//	@Security		Bearer
//	@Router			/api/analytics/top-contractors [get]
func (h *AnalyticsHandler) GetTopContractors(c *gin.Context) {
	period := c.DefaultQuery("period", "month")
	startDate := c.Query("startDate")
	endDate := c.Query("endDate")
	limit := c.DefaultQuery("limit", "10")

	start, end := h.calculateDateRange(period, startDate, endDate)

	h.logger.Info("Getting top contractors",
		zap.String("period", period),
		zap.String("limit", limit),
		zap.Time("start", start),
		zap.Time("end", end),
	)

	// TODO: Implement actual database queries
	contractors := []TopContractor{
		{ContractorID: "1", ContractorName: "ТОО \"Рога и копыта\"", TotalAmount: 450000, InvoiceCount: 45},
		{ContractorID: "2", ContractorName: "АО \"Тех-Снаб\"", TotalAmount: 380000, InvoiceCount: 38},
		{ContractorID: "3", ContractorName: "ИП Иванов И.И.", TotalAmount: 290000, InvoiceCount: 29},
		{ContractorID: "4", ContractorName: "ТОО \"СтройМастер\"", TotalAmount: 215000, InvoiceCount: 22},
		{ContractorID: "5", ContractorName: "ООО \"Альфа-Трейд\"", TotalAmount: 165000, InvoiceCount: 16},
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    contractors,
	})
}

// GetRevenueByMonth получение выручки по месяцам
//
//	@Summary		Get revenue by month
//	@Description	Get monthly revenue data
//	@Tags			Analytics
//	@Accept			json
//	@Produce		json
//	@Param			period		query		string	false	"Time period"
//	@Param			startDate	query		string	false	"Start date (YYYY-MM-DD)"
//	@Param			endDate		query		string	false	"End date (YYYY-MM-DD)"
//	@Success		200			{object}	map[string]interface{}
//	@Failure		500			{object}	map[string]string
//	@Security		Bearer
//	@Router			/api/analytics/revenue-by-month [get]
func (h *AnalyticsHandler) GetRevenueByMonth(c *gin.Context) {
	period := c.DefaultQuery("period", "year")
	startDate := c.Query("startDate")
	endDate := c.Query("endDate")

	start, end := h.calculateDateRange(period, startDate, endDate)

	h.logger.Info("Getting revenue by month",
		zap.String("period", period),
		zap.Time("start", start),
		zap.Time("end", end),
	)

	// TODO: Implement actual database queries
	data := []ChartDataPoint{
		{Date: "Январь", Value: 120000},
		{Date: "Февраль", Value: 135000},
		{Date: "Март", Value: 142000},
		{Date: "Апрель", Value: 128000},
		{Date: "Май", Value: 155000},
		{Date: "Июнь", Value: 148000},
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    data,
	})
}

// calculateDateRange вычисляет диапазон дат на основе периода
func (h *AnalyticsHandler) calculateDateRange(period, startDate, endDate string) (time.Time, time.Time) {
	now := time.Now()
	var start, end time.Time

	switch period {
	case "today":
		start = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		end = time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 999999999, now.Location())

	case "week":
		// Начало недели (понедельник)
		weekday := int(now.Weekday())
		if weekday == 0 {
			weekday = 7 // Воскресенье = 7
		}
		start = now.AddDate(0, 0, -(weekday - 1))
		start = time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, start.Location())
		end = now

	case "month":
		start = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
		end = now

	case "quarter":
		quarter := ((int(now.Month()) - 1) / 3)
		startMonth := time.Month(quarter*3 + 1)
		start = time.Date(now.Year(), startMonth, 1, 0, 0, 0, 0, now.Location())
		end = now

	case "year":
		start = time.Date(now.Year(), 1, 1, 0, 0, 0, 0, now.Location())
		end = now

	case "custom":
		if startDate != "" && endDate != "" {
			if s, err := time.Parse("2006-01-02", startDate); err == nil {
				start = s
			} else {
				start = now.AddDate(0, -1, 0) // Default to last month
			}

			if e, err := time.Parse("2006-01-02", endDate); err == nil {
				end = e
			} else {
				end = now
			}
		} else {
			start = now.AddDate(0, -1, 0)
			end = now
		}

	default:
		// По умолчанию - месяц
		start = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
		end = now
	}

	return start, end
}
