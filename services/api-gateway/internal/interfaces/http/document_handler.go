package http

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rusgainew/kkm-project-mks/api-gateway/internal/domain/errors"
	"github.com/rusgainew/kkm-project-mks/api-gateway/internal/domain/models"
	"github.com/rusgainew/kkm-project-mks/api-gateway/internal/domain/ports"
	"github.com/rusgainew/kkm-project-mks/api-gateway/internal/infrastructure/validation"
	"github.com/rusgainew/kkm-project-mks/services/pkg/conversion"
	"go.uber.org/zap"
)

// DocumentHandler обрабатывает HTTP запросы к документ-сервису
// Он валидирует входные данные и проксирует вызовы в gRPC документ-сервис.
type DocumentHandler struct {
	service   ports.DocumentServiceInterface
	logger    *zap.Logger
	validator *validation.Validator
}

// NewDocumentHandler создаёт обработчик документов
func NewDocumentHandler(service ports.DocumentServiceInterface, logger *zap.Logger) *DocumentHandler {
	return &DocumentHandler{
		service:   service,
		logger:    logger,
		validator: validation.NewValidator(),
	}
}

// CreateDocument создаёт документ
//
//	@Summary	Create document
//	@Description	Create a new document in organization
//	@Tags		Documents
//	@Accept		json
//	@Produce	json
//	@Param		document	body		models.Document	true	"Document payload"
//	@Success	201	{object}	models.APIResponse{data=models.Document}
//	@Failure	400	{object}	models.APIResponse
//	@Failure	500	{object}	models.APIResponse
//	@Router		/documents [post]
//	@Security	BearerAuth
func (h *DocumentHandler) CreateDocument(c *gin.Context) {
	var doc models.Document
	if err := c.ShouldBindJSON(&doc); err != nil {
		h.logger.Warn("Invalid request body", zap.Error(err))
		c.JSON(http.StatusBadRequest, models.APIResponse{Success: false, Error: &models.APIError{Code: "INVALID_INPUT", Message: "Invalid request body", Details: err.Error()}})
		return
	}

	if err := h.validator.ValidateDocument(&doc); err != nil {
		h.logger.Warn("Document validation failed", zap.Error(err))
		c.JSON(http.StatusBadRequest, models.APIResponse{Success: false, Error: &models.APIError{Code: "VALIDATION_FAILED", Message: err.Error()}})
		return
	}

	res, err := h.service.CreateDocument(c.Request.Context(), &doc)
	if err != nil {
		statusCode, apiErr := errors.MapGRPCErrorToHTTP(err)
		c.JSON(statusCode, models.APIResponse{Success: false, Error: &models.APIError{Code: apiErr.Code, Message: apiErr.Message, Details: apiErr.Details}})
		return
	}

	c.JSON(http.StatusCreated, models.APIResponse{Success: true, Data: res})
}

// GetDocument возвращает документ по ID
//
//	@Summary	Get document
//	@Description	Get document by ID
//	@Tags		Documents
//	@Produce	json
//	@Param		id	path	string	true	"Document ID"
//	@Success	200	{object}	models.APIResponse{data=models.Document}
//	@Failure	400	{object}	models.APIResponse
//	@Failure	404	{object}	models.APIResponse
//	@Router		/documents/{id} [get]
//	@Security	BearerAuth
func (h *DocumentHandler) GetDocument(c *gin.Context) {
	id := c.Param("id")
	if err := h.validator.ValidateUUID(id, "document_id"); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{Success: false, Error: &models.APIError{Code: "INVALID_ID", Message: err.Error()}})
		return
	}

	res, err := h.service.GetDocument(c.Request.Context(), id)
	if err != nil {
		statusCode, apiErr := errors.MapGRPCErrorToHTTP(err)
		c.JSON(statusCode, models.APIResponse{Success: false, Error: &models.APIError{Code: apiErr.Code, Message: apiErr.Message, Details: apiErr.Details}})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{Success: true, Data: res})
}

// UpdateDocument обновляет title/content
//
//	@Summary	Update document
//	@Description	Update document title/content
//	@Tags		Documents
//	@Accept		json
//	@Produce	json
//	@Param		id	path	string	true	"Document ID"
//	@Param		document	body	models.Document	true	"Document payload"
//	@Success	200	{object}	models.APIResponse{data=models.Document}
//	@Failure	400	{object}	models.APIResponse
//	@Failure	500	{object}	models.APIResponse
//	@Router		/documents/{id} [put]
//	@Security	BearerAuth
func (h *DocumentHandler) UpdateDocument(c *gin.Context) {
	id := c.Param("id")
	if err := h.validator.ValidateUUID(id, "document_id"); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{Success: false, Error: &models.APIError{Code: "INVALID_ID", Message: err.Error()}})
		return
	}

	var doc models.Document
	if err := c.ShouldBindJSON(&doc); err != nil {
		h.logger.Warn("Invalid request body", zap.Error(err))
		c.JSON(http.StatusBadRequest, models.APIResponse{Success: false, Error: &models.APIError{Code: "INVALID_INPUT", Message: "Invalid request body", Details: err.Error()}})
		return
	}
	doc.ID = id

	if err := h.validator.ValidateString(doc.Title, "title", 1, 255); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{Success: false, Error: &models.APIError{Code: "VALIDATION_FAILED", Message: err.Error()}})
		return
	}
	if err := h.validator.ValidateString(doc.Content, "content", 1, 5000); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{Success: false, Error: &models.APIError{Code: "VALIDATION_FAILED", Message: err.Error()}})
		return
	}

	res, err := h.service.UpdateDocument(c.Request.Context(), &doc)
	if err != nil {
		statusCode, apiErr := errors.MapGRPCErrorToHTTP(err)
		c.JSON(statusCode, models.APIResponse{Success: false, Error: &models.APIError{Code: apiErr.Code, Message: apiErr.Message, Details: apiErr.Details}})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{Success: true, Data: res})
}

// SendDocument отправляет документ на согласование
//
//	@Summary	Send document
//	@Description	Send document for approval
//	@Tags		Documents
//	@Accept		json
//	@Produce	json
//	@Param		id	path	string	true	"Document ID"
//	@Param		body	body	object{recipient_id=string,message=string}	true	"Recipient info"
//	@Success	200	{object}	models.APIResponse{data=models.Document}
//	@Failure	400	{object}	models.APIResponse
//	@Router		/documents/{id}/send [post]
//	@Security	BearerAuth
func (h *DocumentHandler) SendDocument(c *gin.Context) {
	id := c.Param("id")
	if err := h.validator.ValidateUUID(id, "document_id"); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{Success: false, Error: &models.APIError{Code: "INVALID_ID", Message: err.Error()}})
		return
	}

	var req struct {
		RecipientID string `json:"recipient_id" binding:"required"`
		Message     string `json:"message"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{Success: false, Error: &models.APIError{Code: "INVALID_INPUT", Message: "Invalid request body", Details: err.Error()}})
		return
	}
	if err := h.validator.ValidateUUID(req.RecipientID, "recipient_id"); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{Success: false, Error: &models.APIError{Code: "VALIDATION_FAILED", Message: err.Error()}})
		return
	}

	res, err := h.service.SendDocument(c.Request.Context(), id, req.RecipientID, req.Message)
	if err != nil {
		statusCode, apiErr := errors.MapGRPCErrorToHTTP(err)
		c.JSON(statusCode, models.APIResponse{Success: false, Error: &models.APIError{Code: apiErr.Code, Message: apiErr.Message, Details: apiErr.Details}})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{Success: true, Data: res})
}

// ApproveDocument одобряет документ
//
//	@Summary	Approve document
//	@Description	Approve document
//	@Tags		Documents
//	@Accept		json
//	@Produce	json
//	@Param		id	path	string	true	"Document ID"
//	@Param		body	body	object{approved_by=string,comments=string}	true	"Approve payload"
//	@Success	200	{object}	models.APIResponse{data=models.Document}
//	@Failure	400	{object}	models.APIResponse
//	@Router		/documents/{id}/approve [post]
//	@Security	BearerAuth
func (h *DocumentHandler) ApproveDocument(c *gin.Context) {
	id := c.Param("id")
	if err := h.validator.ValidateUUID(id, "document_id"); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{Success: false, Error: &models.APIError{Code: "INVALID_ID", Message: err.Error()}})
		return
	}

	var req struct {
		ApprovedBy string `json:"approved_by" binding:"required"`
		Comments   string `json:"comments"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{Success: false, Error: &models.APIError{Code: "INVALID_INPUT", Message: "Invalid request body", Details: err.Error()}})
		return
	}
	if err := h.validator.ValidateUUID(req.ApprovedBy, "approved_by"); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{Success: false, Error: &models.APIError{Code: "VALIDATION_FAILED", Message: err.Error()}})
		return
	}

	res, err := h.service.ApproveDocument(c.Request.Context(), id, req.ApprovedBy, req.Comments)
	if err != nil {
		statusCode, apiErr := errors.MapGRPCErrorToHTTP(err)
		c.JSON(statusCode, models.APIResponse{Success: false, Error: &models.APIError{Code: apiErr.Code, Message: apiErr.Message, Details: apiErr.Details}})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{Success: true, Data: res})
}

// RejectDocument отклоняет документ
//
//	@Summary	Reject document
//	@Description	Reject document
//	@Tags		Documents
//	@Accept		json
//	@Produce	json
//	@Param		id	path	string	true	"Document ID"
//	@Param		body	body	object{rejected_by=string,reason=string}	true	"Reject payload"
//	@Success	200	{object}	models.APIResponse{data=models.Document}
//	@Failure	400	{object}	models.APIResponse
//	@Router		/documents/{id}/reject [post]
//	@Security	BearerAuth
func (h *DocumentHandler) RejectDocument(c *gin.Context) {
	id := c.Param("id")
	if err := h.validator.ValidateUUID(id, "document_id"); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{Success: false, Error: &models.APIError{Code: "INVALID_ID", Message: err.Error()}})
		return
	}

	var req struct {
		RejectedBy string `json:"rejected_by" binding:"required"`
		Reason     string `json:"reason"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{Success: false, Error: &models.APIError{Code: "INVALID_INPUT", Message: "Invalid request body", Details: err.Error()}})
		return
	}
	if err := h.validator.ValidateUUID(req.RejectedBy, "rejected_by"); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{Success: false, Error: &models.APIError{Code: "VALIDATION_FAILED", Message: err.Error()}})
		return
	}

	res, err := h.service.RejectDocument(c.Request.Context(), id, req.RejectedBy, req.Reason)
	if err != nil {
		statusCode, apiErr := errors.MapGRPCErrorToHTTP(err)
		c.JSON(statusCode, models.APIResponse{Success: false, Error: &models.APIError{Code: apiErr.Code, Message: apiErr.Message, Details: apiErr.Details}})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{Success: true, Data: res})
}

// ArchiveDocument архивирует документ
//
//	@Summary	Archive document
//	@Description	Archive document
//	@Tags		Documents
//	@Produce	json
//	@Param		id	path	string	true	"Document ID"
//	@Success	200	{object}	models.APIResponse{data=models.Document}
//	@Failure	400	{object}	models.APIResponse
//	@Router		/documents/{id}/archive [post]
//	@Security	BearerAuth
func (h *DocumentHandler) ArchiveDocument(c *gin.Context) {
	id := c.Param("id")
	if err := h.validator.ValidateUUID(id, "document_id"); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{Success: false, Error: &models.APIError{Code: "INVALID_ID", Message: err.Error()}})
		return
	}

	res, err := h.service.ArchiveDocument(c.Request.Context(), id)
	if err != nil {
		statusCode, apiErr := errors.MapGRPCErrorToHTTP(err)
		c.JSON(statusCode, models.APIResponse{Success: false, Error: &models.APIError{Code: apiErr.Code, Message: apiErr.Message, Details: apiErr.Details}})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{Success: true, Data: res})
}

// ListDocuments возвращает список документов с пагинацией
//
//	@Summary	List documents
//	@Description	List documents with pagination
//	@Tags		Documents
//	@Produce	json
//	@Param		organization_id	query	string	true	"Organization ID"
//	@Param		status	query	string	false	"Status filter"
//	@Param		created_by	query	string	false	"Created by filter"
//	@Param		page	query	int	false	"Page"
//	@Param		per_page	query	int	false	"Per page"
//	@Success	200	{object}	models.APIResponse{data=models.DocumentListResponse}
//	@Failure	400	{object}	models.APIResponse
//	@Router		/documents [get]
//	@Security	BearerAuth
func (h *DocumentHandler) ListDocuments(c *gin.Context) {
	organizationID := c.Query("organization_id")
	if err := h.validator.ValidateUUID(organizationID, "organization_id"); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{Success: false, Error: &models.APIError{Code: "VALIDATION_FAILED", Message: err.Error()}})
		return
	}

	status := c.Query("status")
	createdBy := c.Query("created_by")

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "20"))
	if err := h.validator.ValidatePagination(page, perPage); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{Success: false, Error: &models.APIError{Code: "VALIDATION_FAILED", Message: err.Error()}})
		return
	}

	docs, pageInfo, err := h.service.ListDocuments(c.Request.Context(), organizationID, status, createdBy, conversion.SafeIntToInt32WithDefault(page, 1), conversion.SafeIntToInt32WithDefault(perPage, 20))
	if err != nil {
		statusCode, apiErr := errors.MapGRPCErrorToHTTP(err)
		c.JSON(statusCode, models.APIResponse{Success: false, Error: &models.APIError{Code: apiErr.Code, Message: apiErr.Message, Details: apiErr.Details}})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{Success: true, Data: models.DocumentListResponse{Documents: docs, Page: pageInfo}})
}
