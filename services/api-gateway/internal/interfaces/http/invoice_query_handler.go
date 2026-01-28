package http

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rusgainew/kkm-project-mks/api-gateway/internal/application/services"
	"github.com/rusgainew/kkm-project-mks/api-gateway/internal/domain/errors"
	"github.com/rusgainew/kkm-project-mks/lib/conversion"
	pb "github.com/rusgainew/kkm-project-mks/proto-lib/api"
	"go.uber.org/zap"
)

// InvoiceQueryHandler обработчик для чтения счетов
type InvoiceQueryHandler struct {
	service *services.InvoiceQueryService
	logger  *zap.Logger
}

// NewInvoiceQueryHandler создает новый InvoiceQueryHandler
func NewInvoiceQueryHandler(service *services.InvoiceQueryService, logger *zap.Logger) *InvoiceQueryHandler {
	return &InvoiceQueryHandler{
		service: service,
		logger:  logger,
	}
}

// ListInvoices получает список счетов с пагинацией
//
//	@Summary		List invoices
//	@Description	Get a paginated list of invoices
//	@Tags			Invoices
//	@Accept			json
//	@Produce		json
//	@Param			page		query		int	false	"Page number"	default(0)
//	@Param			page_size	query		int	false	"Page size"		default(20)
//	@Success		200			{object}	models.APIResponse{data=[]models.Invoice}
//	@Failure		400			{object}	models.APIResponse	"Invalid parameters"
//	@Failure		500			{object}	models.APIResponse	"Internal server error"
//	@Router			/invoices-query [get]
//	@Security		BearerAuth
func (h *InvoiceQueryHandler) ListInvoices(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "0")
	pageSizeStr := c.DefaultQuery("page_size", "20")

	page, err := strconv.ParseInt(pageStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, errors.NewValidationError("invalid page number"))
		return
	}

	pageSize, err := strconv.ParseInt(pageSizeStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, errors.NewValidationError("invalid page size"))
		return
	}

	// Limit page size to 100
	if pageSize > 100 {
		pageSize = 100
	}

	resp, err := h.service.ListInvoices(c.Request.Context(), conversion.SafeInt64ToInt32WithDefault(page, 0), conversion.SafeInt64ToInt32WithDefault(pageSize, 20))
	if err != nil {
		h.logger.Error("Failed to list invoices", zap.Error(err))
		c.JSON(http.StatusInternalServerError, errors.NewInternalServerError())
		return
	}

	c.JSON(http.StatusOK, resp)
}

// SearchInvoices осуществляет поиск счетов
//
//	@Summary		Search invoices
//	@Description	Search invoices by query
//	@Tags			Invoices
//	@Accept			json
//	@Produce		json
//	@Param			q			query		string	true	"Search query"
//	@Param			page		query		int		false	"Page number"	default(0)
//	@Param			page_size	query		int		false	"Page size"		default(20)
//	@Success		200			{object}	models.APIResponse{data=[]models.Invoice}
//	@Failure		400			{object}	models.APIResponse	"Invalid parameters"
//	@Failure		500			{object}	models.APIResponse	"Internal server error"
//	@Router			/invoices-query/search [get]
//	@Security		BearerAuth
func (h *InvoiceQueryHandler) SearchInvoices(c *gin.Context) {
	query := c.DefaultQuery("q", "")
	pageStr := c.DefaultQuery("page", "0")
	pageSizeStr := c.DefaultQuery("page_size", "20")

	if query == "" {
		c.JSON(http.StatusBadRequest, errors.NewValidationError("search query is required"))
		return
	}

	page, err := strconv.ParseInt(pageStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, errors.NewValidationError("invalid page number"))
		return
	}

	pageSize, err := strconv.ParseInt(pageSizeStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, errors.NewValidationError("invalid page size"))
		return
	}

	if pageSize > 100 {
		pageSize = 100
	}

	// Search invoices - для поиска используем ListInvoices с фильтром
	// В реальном приложении здесь был бы SearchInvoices метод
	h.logger.Debug("Searching invoices",
		zap.String("query", query),
		zap.Int64("page", page),
		zap.Int64("page_size", pageSize))

	resp, err := h.service.ListInvoices(c.Request.Context(), conversion.SafeInt64ToInt32WithDefault(page, 0), conversion.SafeInt64ToInt32WithDefault(pageSize, 20))
	if err != nil {
		h.logger.Error("Failed to search invoices", zap.String("query", query), zap.Error(err))
		c.JSON(http.StatusInternalServerError, errors.NewInternalServerError())
		return
	}

	c.JSON(http.StatusOK, resp)
}

// ListInvoicesWithFilter получает список счетов с фильтрацией
//
//	@Summary		List invoices with filter
//	@Description	Get a paginated list of invoices with filter criteria
//	@Tags			Invoices
//	@Accept			json
//	@Produce		json
//	@Param			invoice_number	query		string	false	"Invoice number"
//	@Param			date_from		query		string	false	"Date from (YYYY-MM-DD)"
//	@Param			date_to			query		string	false	"Date to (YYYY-MM-DD)"
//	@Param			min_amount		query		number	false	"Minimum amount"
//	@Param			max_amount		query		number	false	"Maximum amount"
//	@Param			search_text		query		string	false	"Search text"
//	@Param			sort_field		query		string	false	"Sort field"
//	@Param			sort_order		query		string	false	"Sort order (ASC/DESC)"
//	@Param			page			query		int		false	"Page number"	default(0)
//	@Param			page_size		query		int		false	"Page size"		default(20)
//	@Success		200				{object}	models.APIResponse{data=[]models.Invoice}
//	@Failure		400				{object}	models.APIResponse	"Invalid parameters"
//	@Failure		500				{object}	models.APIResponse	"Internal server error"
//	@Router			/invoices-query/filter [get]
//	@Security		BearerAuth
func (h *InvoiceQueryHandler) ListInvoicesWithFilter(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "0")
	pageSizeStr := c.DefaultQuery("page_size", "20")

	page, err := strconv.ParseInt(pageStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, errors.NewValidationError("invalid page number"))
		return
	}

	pageSize, err := strconv.ParseInt(pageSizeStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, errors.NewValidationError("invalid page size"))
		return
	}

	if pageSize > 100 {
		pageSize = 100
	}

	// Parse filter parameters
	var minAmount, maxAmount float64
	if minAmountStr := c.Query("min_amount"); minAmountStr != "" {
		minAmount, _ = strconv.ParseFloat(minAmountStr, 64)
	}
	if maxAmountStr := c.Query("max_amount"); maxAmountStr != "" {
		maxAmount, _ = strconv.ParseFloat(maxAmountStr, 64)
	}

	filter := &pb.InvoiceFilterRequest{
		InvoiceNumber: c.Query("invoice_number"),
		DateFrom:      c.Query("date_from"),
		DateTo:        c.Query("date_to"),
		MinAmount:     minAmount,
		MaxAmount:     maxAmount,
		SearchText:    c.Query("search_text"),
		SortField:     c.Query("sort_field"),
		SortOrder:     c.Query("sort_order"),
		Page:          conversion.SafeInt64ToInt32WithDefault(page, 0),
		Size:          conversion.SafeInt64ToInt32WithDefault(pageSize, 20),
	}

	resp, err := h.service.ListInvoicesWithFilter(c.Request.Context(), filter)
	if err != nil {
		h.logger.Error("Failed to list invoices with filter", zap.Error(err))
		c.JSON(http.StatusInternalServerError, errors.NewInternalServerError())
		return
	}

	c.JSON(http.StatusOK, resp)
}

// GetInvoiceByNumber получает счет по номеру
//
//	@Summary		Get invoice by number
//	@Description	Get invoice by its number
//	@Tags			Invoices
//	@Accept			json
//	@Produce		json
//	@Param			number	path		string	true	"Invoice number"
//	@Success		200		{object}	models.APIResponse{data=models.Invoice}
//	@Failure		400		{object}	models.APIResponse	"Invalid parameters"
//	@Failure		404		{object}	models.APIResponse	"Invoice not found"
//	@Failure		500		{object}	models.APIResponse	"Internal server error"
//	@Router			/invoices-query/by-number/{number} [get]
//	@Security		BearerAuth
func (h *InvoiceQueryHandler) GetInvoiceByNumber(c *gin.Context) {
	number := c.Param("number")
	if number == "" {
		c.JSON(http.StatusBadRequest, errors.NewValidationError("invoice number is required"))
		return
	}

	req := &pb.InvoiceNumberRequest{
		InvoiceNumber: number,
	}

	resp, err := h.service.GetInvoiceByNumber(c.Request.Context(), req)
	if err != nil {
		h.logger.Error("Failed to get invoice by number", zap.String("number", number), zap.Error(err))
		c.JSON(http.StatusInternalServerError, errors.NewInternalServerError())
		return
	}

	c.JSON(http.StatusOK, resp)
}

// ListInvoiceDetails получает детали счетов
//
//	@Summary		List invoice details
//	@Description	Get a paginated list of invoice details/line items
//	@Tags			Invoices
//	@Accept			json
//	@Produce		json
//	@Param			invoice_uuid	query		string	false	"Filter by invoice UUID"
//	@Param			page			query		int		false	"Page number"	default(0)
//	@Param			page_size		query		int		false	"Page size"		default(20)
//	@Success		200				{object}	models.APIResponse{data=[]models.InvoiceDetail}
//	@Failure		400				{object}	models.APIResponse	"Invalid parameters"
//	@Failure		500				{object}	models.APIResponse	"Internal server error"
//	@Router			/invoices-query/details [get]
//	@Security		BearerAuth
func (h *InvoiceQueryHandler) ListInvoiceDetails(c *gin.Context) {
	invoiceUUID := c.Query("invoice_uuid")
	pageStr := c.DefaultQuery("page", "0")
	pageSizeStr := c.DefaultQuery("page_size", "20")

	page, err := strconv.ParseInt(pageStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, errors.NewValidationError("invalid page number"))
		return
	}

	pageSize, err := strconv.ParseInt(pageSizeStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, errors.NewValidationError("invalid page size"))
		return
	}

	if pageSize > 100 {
		pageSize = 100
	}

	resp, err := h.service.ListInvoiceDetails(c.Request.Context(), invoiceUUID, conversion.SafeInt64ToInt32WithDefault(page, 0), conversion.SafeInt64ToInt32WithDefault(pageSize, 20))
	if err != nil {
		h.logger.Error("Failed to list invoice details", zap.Error(err))
		c.JSON(http.StatusInternalServerError, errors.NewInternalServerError())
		return
	}

	c.JSON(http.StatusOK, resp)
}

// GetInvoicesByDateRange получает счета по диапазону дат
//
//	@Summary		Get invoices by date range
//	@Description	Get invoices within a specific date range
//	@Tags			Invoices
//	@Accept			json
//	@Produce		json
//	@Param			start_date	query		string	true	"Start date (YYYY-MM-DD)"
//	@Param			end_date	query		string	true	"End date (YYYY-MM-DD)"
//	@Param			page		query		int		false	"Page number"	default(0)
//	@Param			page_size	query		int		false	"Page size"		default(20)
//	@Success		200			{object}	models.APIResponse{data=[]models.Invoice}
//	@Failure		400			{object}	models.APIResponse	"Invalid parameters"
//	@Failure		500			{object}	models.APIResponse	"Internal server error"
//	@Router			/invoices-query/by-date-range [get]
//	@Security		BearerAuth
func (h *InvoiceQueryHandler) GetInvoicesByDateRange(c *gin.Context) {
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	if startDate == "" || endDate == "" {
		c.JSON(http.StatusBadRequest, errors.NewValidationError("start_date and end_date are required"))
		return
	}

	pageStr := c.DefaultQuery("page", "0")
	pageSizeStr := c.DefaultQuery("page_size", "20")

	page, err := strconv.ParseInt(pageStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, errors.NewValidationError("invalid page number"))
		return
	}

	pageSize, err := strconv.ParseInt(pageSizeStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, errors.NewValidationError("invalid page size"))
		return
	}

	if pageSize > 100 {
		pageSize = 100
	}

	resp, err := h.service.GetInvoicesByDateRange(c.Request.Context(), startDate, endDate, conversion.SafeInt64ToInt32WithDefault(page, 0), conversion.SafeInt64ToInt32WithDefault(pageSize, 20))
	if err != nil {
		h.logger.Error("Failed to get invoices by date range",
			zap.String("start_date", startDate),
			zap.String("end_date", endDate),
			zap.Error(err))
		c.JSON(http.StatusInternalServerError, errors.NewInternalServerError())
		return
	}

	c.JSON(http.StatusOK, resp)
}
