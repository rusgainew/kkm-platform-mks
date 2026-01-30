// Файл api-gateway/internal/interfaces/http/analytics_handler.go содержит реализацию пакета http.
package http

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rusgainew/kkm-project-mks/api-gateway/internal/infrastructure/client"
	"go.uber.org/zap"
)

// AnalyticsHandler обработчик для аналитики Dashboard
type AnalyticsHandler struct {
	analyticsClient *client.AnalyticsClient
	logger          *zap.Logger
}

// NewAnalyticsHandler создает новый AnalyticsHandler
func NewAnalyticsHandler(
	analyticsClient *client.AnalyticsClient,
	logger *zap.Logger,
) *AnalyticsHandler {
	return &AnalyticsHandler{
		analyticsClient: analyticsClient,
		logger:          logger,
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

	// Получаем данные из analytics-server через gRPC
	resp, err := h.analyticsClient.GetDashboardStats(c.Request.Context(), start, end)
	if err != nil {
		h.logger.Error("Failed to get analytics stats", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to retrieve analytics data",
		})
		return
	}

	// TODO: Вычислить изменения относительно предыдущего периода
	stats := DashboardStats{
		TotalRevenue:         resp.TotalRevenue,
		TotalInvoices:        int32(resp.TotalInvoices),
		AverageInvoiceAmount: resp.AverageAmount,
		ActiveContractors:    int32(resp.UniqueContractors),
		PendingInvoices:      int32(resp.PendingCount),
		ApprovedInvoices:     int32(resp.ApprovedCount),
		RejectedInvoices:     int32(resp.RejectedCount),
		RevenueChange:        0,
		InvoiceCountChange:   0,
		AverageAmountChange:  0,
		ContractorsChange:    0,
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

	// Определяем гранулярность на основе периода
	granularity := h.determineGranularity(start, end)

	h.logger.Info("Getting sales chart",
		zap.String("period", period),
		zap.String("granularity", granularity),
		zap.Time("start", start),
		zap.Time("end", end),
	)

	// Получаем данные из analytics-server через gRPC
	resp, err := h.analyticsClient.GetSalesChart(c.Request.Context(), start, end, granularity)
	if err != nil {
		h.logger.Error("Failed to get sales data", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to retrieve sales data",
		})
		return
	}

	// Конвертируем в формат ответа
	data := make([]ChartDataPoint, len(resp.DataPoints))
	for i, point := range resp.DataPoints {
		data[i] = ChartDataPoint{
			Date:  point.Date,
			Value: point.Amount,
		}
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

	// Получаем данные из analytics-server через gRPC
	resp, err := h.analyticsClient.GetStatusDistribution(c.Request.Context(), start, end)
	if err != nil {
		h.logger.Error("Failed to get status distribution", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to retrieve status statistics",
		})
		return
	}

	// Вычисляем общее количество для процентов
	var total float64
	for _, dist := range resp.Items {
		total += float64(dist.Count)
	}

	// Мапа для дружественных имен и цветов
	statusNames := map[string]string{
		"draft":    "Черновик",
		"sent":     "Отправлено",
		"signed":   "Подписано",
		"accepted": "Утверждено",
		"rejected": "Отклонено",
		"revoked":  "Аннулировано",
	}

	statusColors := map[string]string{
		"draft":    "#94a3b8",
		"sent":     "#f59e0b",
		"signed":   "#3b82f6",
		"accepted": "#10b981",
		"rejected": "#ef4444",
		"revoked":  "#6b7280",
	}

	// Конвертируем в формат ответа
	data := make([]PieChartDataPoint, len(resp.Items))
	for i, dist := range resp.Items {
		percentage := dist.Percentage
		if percentage == 0 && total > 0 {
			percentage = (float64(dist.Count) / total) * 100
		}

		name := statusNames[dist.Status]
		if name == "" {
			name = dist.Status
		}

		color := statusColors[dist.Status]

		data[i] = PieChartDataPoint{
			Name:       name,
			Value:      float64(dist.Count),
			Percentage: percentage,
			Color:      color,
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": PieChartData{
			Data: data,
		},
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

	// Получаем данные из analytics-server через gRPC
	resp, err := h.analyticsClient.GetOperationTypeDistribution(c.Request.Context(), start, end)
	if err != nil {
		h.logger.Error("Failed to get operation type distribution", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to retrieve operation type statistics",
		})
		return
	}

	// Вычисляем общее количество для процентов
	var total float64
	for _, dist := range resp.Items {
		total += float64(dist.Count)
	}

	// Мапа для дружественных имен и цветов
	opTypeNames := map[string]string{
		"local":  "Местные операции",
		"import": "Импортные операции",
	}

	opTypeColors := map[string]string{
		"local":  "#3b82f6",
		"import": "#8b5cf6",
	}

	// Конвертируем в формат ответа
	data := make([]PieChartDataPoint, len(resp.Items))
	for i, dist := range resp.Items {
		percentage := dist.Percentage
		if percentage == 0 && total > 0 {
			percentage = (float64(dist.Count) / total) * 100
		}

		name := opTypeNames[dist.OperationType]
		if name == "" {
			name = dist.OperationType
		}

		color := opTypeColors[dist.OperationType]

		data[i] = PieChartDataPoint{
			Name:       name,
			Value:      float64(dist.Count),
			Percentage: percentage,
			Color:      color,
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": PieChartData{
			Data: data,
		},
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
	limitStr := c.DefaultQuery("limit", "10")

	limit, err := strconv.ParseInt(limitStr, 10, 32)
	if err != nil || limit <= 0 {
		limit = 10
	}

	start, end := h.calculateDateRange(period, startDate, endDate)

	h.logger.Info("Getting top contractors",
		zap.String("period", period),
		zap.Int64("limit", limit),
		zap.Time("start", start),
		zap.Time("end", end),
	)

	// Получаем данные из analytics-server через gRPC
	resp, err := h.analyticsClient.GetTopContractors(c.Request.Context(), start, end, int32(limit))
	if err != nil {
		h.logger.Error("Failed to get top contractors", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to retrieve top contractors",
		})
		return
	}

	// Конвертируем в формат ответа
	contractors := make([]TopContractor, len(resp.Contractors))
	for i, data := range resp.Contractors {
		contractors[i] = TopContractor{
			ContractorID:   data.ContractorId,
			ContractorName: data.ContractorName,
			TotalAmount:    data.TotalAmount,
			InvoiceCount:   int32(data.InvoiceCount),
		}
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

	// Получаем данные из analytics-server через gRPC
	resp, err := h.analyticsClient.GetMonthlyRevenue(c.Request.Context(), start, end)
	if err != nil {
		h.logger.Error("Failed to get monthly revenue", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to retrieve monthly revenue",
		})
		return
	}

	// Мапа для названий месяцев на русском
	monthNames := []string{
		"Январь", "Февраль", "Март", "Апрель", "Май", "Июнь",
		"Июль", "Август", "Сентябрь", "Октябрь", "Ноябрь", "Декабрь",
	}

	// Конвертируем в формат ответа
	data := make([]ChartDataPoint, len(resp.Months))
	for i, month := range resp.Months {
		// Парсим месяц из формата "2006-01"
		monthTime, err := time.Parse("2006-01", month.Month)
		if err == nil {
			monthName := monthNames[monthTime.Month()-1]
			data[i] = ChartDataPoint{
				Date:  monthName,
				Value: month.Revenue,
			}
		} else {
			data[i] = ChartDataPoint{
				Date:  month.Month,
				Value: month.Revenue,
			}
		}
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

// determineGranularity определяет гранулярность данных на основе диапазона дат
func (h *AnalyticsHandler) determineGranularity(start, end time.Time) string {
	days := int(end.Sub(start).Hours() / 24)

	switch {
	case days <= 7:
		return "day"
	case days <= 90:
		return "week"
	default:
		return "month"
	}
}
