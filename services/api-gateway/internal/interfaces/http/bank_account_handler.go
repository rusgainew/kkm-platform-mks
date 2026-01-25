package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rusgainew/kkm-project-mks/api-gateway/internal/domain/models"
	"github.com/rusgainew/kkm-project-mks/api-gateway/internal/domain/ports"
	"github.com/rusgainew/kkm-project-mks/api-gateway/internal/infrastructure/validation"
	"go.uber.org/zap"
)

// BankAccountHandler отвечает только за HTTP запросы банковских счетов
type BankAccountHandler struct {
	*BaseHandler
	service   ports.BankAccountServiceInterface
	validator *validation.Validator
}

// NewBankAccountHandler создает новый bank account handler
func NewBankAccountHandler(service ports.BankAccountServiceInterface, logger *zap.Logger) *BankAccountHandler {
	return &BankAccountHandler{
		BaseHandler: NewBaseHandler(logger),
		service:     service,
		validator:   validation.NewValidator(),
	}
}

// CreateBankAccount создает новый банковский счет
//
//	@Summary		Create bank account
//	@Description	Create a new bank account
//	@Tags			Bank Accounts
//	@Accept			json
//	@Produce		json
//	@Param			bankAccount	body		models.BankAccount	true	"Bank account data"
//	@Success		201			{object}	models.APIResponse{data=models.BankAccount}
//	@Failure		400			{object}	models.APIResponse	"Invalid input"
//	@Failure		500			{object}	models.APIResponse	"Internal server error"
//	@Router			/bank-accounts [post]
//	@Security		BearerAuth
func (h *BankAccountHandler) CreateBankAccount(c *gin.Context) {
	var bankAccount models.BankAccount
	if !h.BindJSON(c, &bankAccount) {
		return
	}

	// Validate bank account data
	if err := h.validator.ValidateBankAccount(&bankAccount); err != nil {
		h.Logger().Warn("Bank account validation failed", zap.Error(err))
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error: &models.APIError{
				Code:    "VALIDATION_FAILED",
				Message: "Bank account validation failed",
				Details: err.Error(),
			},
		})
		return
	}

	result, err := h.service.CreateBankAccount(c.Request.Context(), &bankAccount)
	if err != nil {
		h.RespondWithGRPCError(c, err, "Create bank account")
		return
	}

	h.RespondWithCreated(c, result)
}

// GetBankAccount возвращает банковский счет по ID
//
//	@Summary		Get bank account by ID
//	@Description	Retrieve a bank account by its ID
//	@Tags			Bank Accounts
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"Bank Account ID"
//	@Success		200	{object}	models.APIResponse{data=models.BankAccount}
//	@Failure		400	{object}	models.APIResponse	"Invalid ID"
//	@Failure		404	{object}	models.APIResponse	"Not found"
//	@Failure		500	{object}	models.APIResponse	"Internal server error"
//	@Router			/bank-accounts/{id} [get]
//	@Security		BearerAuth
func (h *BankAccountHandler) GetBankAccount(c *gin.Context) {
	id := c.Param("id")
	if !h.ValidateUUID(c, id, "bank_account_id") {
		return
	}

	result, err := h.service.GetBankAccount(c.Request.Context(), id)
	if err != nil {
		h.RespondWithGRPCError(c, err, "Get bank account")
		return
	}

	h.RespondWithSuccess(c, result)
}

// UpdateBankAccount обновляет банковский счет
//
//	@Summary		Update bank account
//	@Description	Update an existing bank account
//	@Tags			Bank Accounts
//	@Accept			json
//	@Produce		json
//	@Param			id			path		string				true	"Bank Account ID"
//	@Param			bankAccount	body		models.BankAccount	true	"Updated bank account data"
//	@Success		200			{object}	models.APIResponse{data=models.BankAccount}
//	@Failure		400			{object}	models.APIResponse	"Invalid input"
//	@Failure		404			{object}	models.APIResponse	"Not found"
//	@Failure		500			{object}	models.APIResponse	"Internal server error"
//	@Router			/bank-accounts/{id} [put]
//	@Security		BearerAuth
func (h *BankAccountHandler) UpdateBankAccount(c *gin.Context) {
	id := c.Param("id")
	if !h.ValidateUUID(c, id, "bank_account_id") {
		return
	}

	var bankAccount models.BankAccount
	if !h.BindJSON(c, &bankAccount) {
		return
	}

	// Validate bank account data
	if err := h.validator.ValidateBankAccount(&bankAccount); err != nil {
		h.Logger().Warn("Bank account validation failed", zap.Error(err))
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error: &models.APIError{
				Code:    "VALIDATION_FAILED",
				Message: "Bank account validation failed",
				Details: err.Error(),
			},
		})
		return
	}

	bankAccount.ID = id
	result, err := h.service.UpdateBankAccount(c.Request.Context(), &bankAccount)
	if err != nil {
		h.RespondWithGRPCError(c, err, "Update bank account")
		return
	}

	h.RespondWithSuccess(c, result)
}

// DeleteBankAccount удаляет банковский счет
//
//	@Summary		Delete bank account
//	@Description	Delete a bank account by ID
//	@Tags			Bank Accounts
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string				true	"Bank Account ID"
//	@Success		200	{object}	models.APIResponse	"Bank account deleted successfully"
//	@Failure		400	{object}	models.APIResponse	"Invalid ID"
//	@Failure		404	{object}	models.APIResponse	"Not found"
//	@Failure		500	{object}	models.APIResponse	"Internal server error"
//	@Router			/bank-accounts/{id} [delete]
//	@Security		BearerAuth
func (h *BankAccountHandler) DeleteBankAccount(c *gin.Context) {
	id := c.Param("id")
	if !h.ValidateUUID(c, id, "bank_account_id") {
		return
	}

	err := h.service.DeleteBankAccount(c.Request.Context(), id)
	if err != nil {
		h.RespondWithGRPCError(c, err, "Delete bank account")
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data:    map[string]string{"message": "Bank account deleted successfully"},
	})
}
