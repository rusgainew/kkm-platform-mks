package http

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rusgainew/kkm-project-mks/api-gateway/internal/application/services"
	"github.com/rusgainew/kkm-project-mks/api-gateway/internal/domain/errors"
	pb "github.com/rusgainew/kkm-project-mks/proto-lib/api"
	"go.uber.org/zap"
)

// BankAccountQueryHandler обработчик для чтения банковских счетов
type BankAccountQueryHandler struct {
	service *services.BankAccountQueryService
	logger  *zap.Logger
}

// NewBankAccountQueryHandler создает новый BankAccountQueryHandler
func NewBankAccountQueryHandler(service *services.BankAccountQueryService, logger *zap.Logger) *BankAccountQueryHandler {
	return &BankAccountQueryHandler{
		service: service,
		logger:  logger,
	}
}

// ListBankAccounts получает список банковских счетов с пагинацией
//
//	@Summary		List bank accounts
//	@Description	Get a paginated list of bank accounts
//	@Tags			Bank Accounts
//	@Accept			json
//	@Produce		json
//	@Param			page		query		int	false	"Page number"	default(0)
//	@Param			page_size	query		int	false	"Page size"		default(20)
//	@Success		200			{object}	models.APIResponse{data=[]models.BankAccount}
//	@Failure		400			{object}	models.APIResponse	"Invalid parameters"
//	@Failure		500			{object}	models.APIResponse	"Internal server error"
//	@Router			/bank-accounts-query [get]
//	@Security		BearerAuth
func (h *BankAccountQueryHandler) ListBankAccounts(c *gin.Context) {
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

	resp, err := h.service.ListBankAccounts(c.Request.Context(), int32(page), int32(pageSize))
	if err != nil {
		h.logger.Error("Failed to list bank accounts", zap.Error(err))
		c.JSON(http.StatusInternalServerError, errors.NewInternalServerError())
		return
	}

	c.JSON(http.StatusOK, resp)
}

// ListBankAccountsWithFilter получает список банковских счетов с фильтрацией
//
//	@Summary		List bank accounts with filter
//	@Description	Get a paginated list of bank accounts with filter criteria
//	@Tags			Bank Accounts
//	@Accept			json
//	@Produce		json
//	@Param			account_name	query		string	false	"Account/bank name filter"
//	@Param			bank_account	query		string	false	"Bank account number filter"
//	@Param			is_active		query		bool	false	"Is active filter"
//	@Param			search_text		query		string	false	"Search text"
//	@Param			sort_field		query		string	false	"Sort field (account_name, bank_account)"
//	@Param			sort_order		query		string	false	"Sort order (ASC/DESC)"
//	@Param			page			query		int		false	"Page number"	default(0)
//	@Param			page_size		query		int		false	"Page size"		default(20)
//	@Success		200				{object}	models.APIResponse{data=[]models.BankAccount}
//	@Failure		400				{object}	models.APIResponse	"Invalid parameters"
//	@Failure		500				{object}	models.APIResponse	"Internal server error"
//	@Router			/bank-accounts-query/filter [get]
//	@Security		BearerAuth
func (h *BankAccountQueryHandler) ListBankAccountsWithFilter(c *gin.Context) {
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

	// Parse is_active
	isActive := false
	if isActiveStr := c.Query("is_active"); isActiveStr != "" {
		isActive, _ = strconv.ParseBool(isActiveStr)
	}

	filter := &pb.BankAccountFilterRequest{
		AccountName: c.Query("account_name"),
		BankAccount: c.Query("bank_account"),
		IsActive:    isActive,
		SearchText:  c.Query("search_text"),
		SortField:   c.Query("sort_field"),
		SortOrder:   c.Query("sort_order"),
		Page:        int32(page),
		Size:        int32(pageSize),
	}

	resp, err := h.service.ListBankAccountsWithFilter(c.Request.Context(), filter)
	if err != nil {
		h.logger.Error("Failed to list bank accounts with filter", zap.Error(err))
		c.JSON(http.StatusInternalServerError, errors.NewInternalServerError())
		return
	}

	c.JSON(http.StatusOK, resp)
}
