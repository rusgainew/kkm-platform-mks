// Файл api-gateway/internal/interfaces/http/invoice_handler.go содержит реализацию пакета http.
package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rusgainew/kkm-project-mks/api-gateway/internal/domain/errors"
	"github.com/rusgainew/kkm-project-mks/api-gateway/internal/domain/models"
	"github.com/rusgainew/kkm-project-mks/api-gateway/internal/domain/ports"
	"github.com/rusgainew/kkm-project-mks/api-gateway/internal/infrastructure/validation"
	"go.uber.org/zap"
)

// InvoiceHandler отвечает только за HTTP запросы счетов
type InvoiceHandler struct {
	service   ports.InvoiceServiceInterface
	logger    *zap.Logger
	validator *validation.Validator
}

// NewInvoiceHandler создает новый invoice handler
func NewInvoiceHandler(service ports.InvoiceServiceInterface, logger *zap.Logger) *InvoiceHandler {
	return &InvoiceHandler{
		service:   service,
		logger:    logger,
		validator: validation.NewValidator(),
	}
}

// CreateInvoice создает новый счет
//
//	@Summary		Create invoice
//	@Description	Create a new invoice. Requires admin or manager role.
//	@Tags			Invoices
//	@Accept			json
//	@Produce		json
//	@Param			invoice	body		models.Invoice							true	"Invoice data with all required fields"
//	@Success		201		{object}	models.APIResponse{data=models.Invoice}	"Invoice successfully created"
//	@Failure		400		{object}	models.APIResponse						"Invalid input - validation errors"
//	@Failure		401		{object}	models.APIResponse						"Unauthorized"
//	@Failure		403		{object}	models.APIResponse						"Forbidden - requires admin or manager role"
//	@Failure		500		{object}	models.APIResponse						"Internal server error"
//	@Router			/invoices [post]
//	@Security		BearerAuth
func (h *InvoiceHandler) CreateInvoice(c *gin.Context) {
	var invoice models.Invoice
	if err := c.ShouldBindJSON(&invoice); err != nil {
		h.logger.Warn("Invalid request body", zap.Error(err))
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error: &models.APIError{
				Code:    "INVALID_INPUT",
				Message: "Invalid request body",
				Details: err.Error(),
			},
		})
		return
	}

	// Validate invoice data
	if err := h.validator.ValidateInvoice(&invoice); err != nil {
		h.logger.Warn("Invoice validation failed", zap.Error(err))
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error: &models.APIError{
				Code:    "VALIDATION_FAILED",
				Message: "Invoice validation failed",
				Details: err.Error(),
			},
		})
		return
	}

	result, err := h.service.CreateInvoice(c.Request.Context(), &invoice)
	if err != nil {
		statusCode, apiErr := errors.MapGRPCErrorToHTTP(err)
		h.logger.Error("Failed to create invoice", zap.Error(err))
		c.JSON(statusCode, models.APIResponse{
			Success: false,
			Error: &models.APIError{
				Code:    apiErr.Code,
				Message: apiErr.Message,
				Details: apiErr.Details,
			},
		})
		return
	}

	c.JSON(http.StatusCreated, models.APIResponse{
		Success: true,
		Data:    result,
	})
}

// UpdateInvoice обновляет счет
//
//	@Summary		Update invoice
//	@Description	Update an existing invoice
//	@Tags			Invoices
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string			true	"Invoice ID"
//	@Param			invoice	body		models.Invoice	true	"Updated invoice data"
//	@Success		200		{object}	models.APIResponse{data=models.Invoice}
//	@Failure		400		{object}	models.APIResponse	"Invalid input"
//	@Failure		500		{object}	models.APIResponse	"Internal server error"
//	@Router			/invoices/{id} [put]
//	@Security		BearerAuth
func (h *InvoiceHandler) UpdateInvoice(c *gin.Context) {
	id := c.Param("id")

	// Validate UUID
	if err := h.validator.ValidateUUID(id, "invoice_id"); err != nil {
		h.logger.Warn("Invalid invoice ID", zap.String("id", id), zap.Error(err))
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error: &models.APIError{
				Code:    "INVALID_ID",
				Message: err.Error(),
			},
		})
		return
	}

	var invoice models.Invoice
	if err := c.ShouldBindJSON(&invoice); err != nil {
		h.logger.Warn("Invalid request body", zap.Error(err))
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error: &models.APIError{
				Code:    "INVALID_INPUT",
				Message: "Invalid request body",
				Details: err.Error(),
			},
		})
		return
	}

	// Validate invoice data
	if err := h.validator.ValidateInvoice(&invoice); err != nil {
		h.logger.Warn("Invoice validation failed", zap.Error(err))
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error: &models.APIError{
				Code:    "VALIDATION_FAILED",
				Message: "Invoice validation failed",
				Details: err.Error(),
			},
		})
		return
	}

	invoice.ID = id
	result, err := h.service.UpdateInvoice(c.Request.Context(), &invoice)
	if err != nil {
		statusCode, apiErr := errors.MapGRPCErrorToHTTP(err)
		h.logger.Error("Failed to update invoice", zap.String("id", id), zap.Error(err))
		c.JSON(statusCode, models.APIResponse{
			Success: false,
			Error: &models.APIError{
				Code:    apiErr.Code,
				Message: apiErr.Message,
				Details: apiErr.Details,
			},
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data:    result,
	})
}

// SignInvoice подписывает счет
//
//	@Summary		Sign invoice
//	@Description	Sign an existing invoice with digital signature
//	@Tags			Invoices
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string									true	"Invoice ID (UUID)"
//	@Param			body	body		object{signature_data=string}			true	"Signature data"
//	@Success		200		{object}	models.APIResponse{data=models.Invoice}	"Invoice signed successfully"
//	@Failure		400		{object}	models.APIResponse						"Invalid input"
//	@Failure		500		{object}	models.APIResponse						"Internal server error"
//	@Router			/invoices/{id}/sign [post]
//	@Security		BearerAuth
func (h *InvoiceHandler) SignInvoice(c *gin.Context) {
	id := c.Param("id")

	// Validate UUID
	if err := h.validator.ValidateUUID(id, "invoice_id"); err != nil {
		h.logger.Warn("Invalid invoice ID", zap.String("id", id), zap.Error(err))
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error: &models.APIError{
				Code:    "INVALID_ID",
				Message: err.Error(),
			},
		})
		return
	}

	var signRequest struct {
		SignatureData string `json:"signature_data" binding:"required"`
	}

	if err := c.ShouldBindJSON(&signRequest); err != nil {
		h.logger.Warn("Invalid request body", zap.Error(err))
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error: &models.APIError{
				Code:    "INVALID_INPUT",
				Message: "Invalid request body",
				Details: err.Error(),
			},
		})
		return
	}

	result, err := h.service.SignInvoice(c.Request.Context(), id, signRequest.SignatureData)
	if err != nil {
		statusCode, apiErr := errors.MapGRPCErrorToHTTP(err)
		h.logger.Error("Failed to sign invoice", zap.String("id", id), zap.Error(err))
		c.JSON(statusCode, models.APIResponse{
			Success: false,
			Error: &models.APIError{
				Code:    apiErr.Code,
				Message: apiErr.Message,
				Details: apiErr.Details,
			},
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data:    result,
	})
}

// RevokeInvoice отзывает счет
//
//	@Summary		Revoke invoice
//	@Description	Revoke an existing invoice
//	@Tags			Invoices
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string					true	"Invoice ID"
//	@Param			body	body		object{reason=string}	false	"Revocation reason"
//	@Success		200		{object}	models.APIResponse{data=models.Invoice}
//	@Failure		400		{object}	models.APIResponse	"Invalid input"
//	@Failure		500		{object}	models.APIResponse	"Internal server error"
//	@Router			/invoices/{id}/revoke [post]
//	@Security		BearerAuth
func (h *InvoiceHandler) RevokeInvoice(c *gin.Context) {
	id := c.Param("id")

	// Validate UUID
	if err := h.validator.ValidateUUID(id, "invoice_id"); err != nil {
		h.logger.Warn("Invalid invoice ID", zap.String("id", id), zap.Error(err))
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error: &models.APIError{
				Code:    "INVALID_ID",
				Message: err.Error(),
			},
		})
		return
	}

	var revokeRequest struct {
		Reason string `json:"reason"`
	}

	// Bind is optional - reason may be empty
	_ = c.ShouldBindJSON(&revokeRequest)

	result, err := h.service.RevokeInvoice(c.Request.Context(), id, revokeRequest.Reason)
	if err != nil {
		statusCode, apiErr := errors.MapGRPCErrorToHTTP(err)
		h.logger.Error("Failed to revoke invoice", zap.String("id", id), zap.Error(err))
		c.JSON(statusCode, models.APIResponse{
			Success: false,
			Error: &models.APIError{
				Code:    apiErr.Code,
				Message: apiErr.Message,
				Details: apiErr.Details,
			},
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data:    result,
	})
}

// AcceptOrRejectInvoice принимает или отклоняет счет
//
//	@Summary		Accept or reject invoice
//	@Description	Accept or reject an existing invoice
//	@Tags			Invoices
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string								true	"Invoice ID"
//	@Param			body	body		object{accept=bool,reason=string}	true	"Accept/reject decision"
//	@Success		200		{object}	models.APIResponse{data=models.Invoice}
//	@Failure		400		{object}	models.APIResponse	"Invalid input"
//	@Failure		500		{object}	models.APIResponse	"Internal server error"
//	@Router			/invoices/{id}/accept-reject [post]
//	@Security		BearerAuth
func (h *InvoiceHandler) AcceptOrRejectInvoice(c *gin.Context) {
	id := c.Param("id")

	// Validate UUID
	if err := h.validator.ValidateUUID(id, "invoice_id"); err != nil {
		h.logger.Warn("Invalid invoice ID", zap.String("id", id), zap.Error(err))
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error: &models.APIError{
				Code:    "INVALID_ID",
				Message: err.Error(),
			},
		})
		return
	}

	var request struct {
		Accept bool   `json:"accept" binding:"required"`
		Reason string `json:"reason"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		h.logger.Warn("Invalid request body", zap.Error(err))
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error: &models.APIError{
				Code:    "INVALID_INPUT",
				Message: "Invalid request body",
				Details: err.Error(),
			},
		})
		return
	}

	result, err := h.service.AcceptOrRejectInvoice(c.Request.Context(), id, request.Accept, request.Reason)
	if err != nil {
		statusCode, apiErr := errors.MapGRPCErrorToHTTP(err)
		h.logger.Error("Failed to accept/reject invoice", zap.String("id", id), zap.Bool("accept", request.Accept), zap.Error(err))
		c.JSON(statusCode, models.APIResponse{
			Success: false,
			Error: &models.APIError{
				Code:    apiErr.Code,
				Message: apiErr.Message,
				Details: apiErr.Details,
			},
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data:    result,
	})
}
