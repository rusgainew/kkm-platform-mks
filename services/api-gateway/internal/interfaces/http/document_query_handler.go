package http

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rusgainew/kkm-project-mks/api-gateway/internal/application/services"
	"github.com/rusgainew/kkm-project-mks/api-gateway/internal/domain/errors"
	"github.com/rusgainew/kkm-project-mks/services/pkg/conversion"
	"go.uber.org/zap"
)

// DocumentQueryHandler обработчик для чтения документов
type DocumentQueryHandler struct {
	service *services.DocumentQueryService
	logger  *zap.Logger
}

// NewDocumentQueryHandler создает новый DocumentQueryHandler
func NewDocumentQueryHandler(service *services.DocumentQueryService, logger *zap.Logger) *DocumentQueryHandler {
	return &DocumentQueryHandler{
		service: service,
		logger:  logger,
	}
}

// GetDocument получает документ по ID
//
//	@Summary		Get document by ID
//	@Description	Retrieve a single document by its ID
//	@Tags			Documents Query
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string	true	"Document ID"
//	@Success		200		{object}	models.Document
//	@Failure		404		{object}	models.APIResponse	"Document not found"
//	@Failure		500		{object}	models.APIResponse	"Internal server error"
//	@Router			/documents-query/{id} [get]
//	@Security		BearerAuth
func (h *DocumentQueryHandler) GetDocument(c *gin.Context) {
	documentID := c.Param("id")
	if documentID == "" {
		c.JSON(http.StatusBadRequest, errors.NewValidationError("document_id is required"))
		return
	}

	doc, err := h.service.GetDocument(c.Request.Context(), documentID)
	if err != nil {
		h.logger.Error("Failed to get document", zap.String("document_id", documentID), zap.Error(err))
		c.JSON(http.StatusInternalServerError, errors.NewInternalServerError())
		return
	}

	if doc == nil {
		c.JSON(http.StatusNotFound, errors.NewNotFoundError("document not found"))
		return
	}

	c.JSON(http.StatusOK, doc)
}

// ListDocuments получает список документов с пагинацией и фильтрацией
//
//	@Summary		List documents
//	@Description	Get a paginated list of documents with optional filters
//	@Tags			Documents Query
//	@Accept			json
//	@Produce		json
//	@Param			page				query		int		false	"Page number"	default(1)
//	@Param			per_page			query		int		false	"Items per page"	default(10)
//	@Param			status				query		string	false	"Document status filter"
//	@Param			type				query		string	false	"Document type filter"
//	@Param			company_id			query		string	false	"Company ID filter"
//	@Param			approval_status		query		string	false	"Approval status filter"
//	@Success		200		{object}	models.DocumentListResponse
//	@Failure		400		{object}	models.APIResponse	"Invalid parameters"
//	@Failure		500		{object}	models.APIResponse	"Internal server error"
//	@Router			/documents-query [get]
//	@Security		BearerAuth
func (h *DocumentQueryHandler) ListDocuments(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	perPageStr := c.DefaultQuery("per_page", "10")
	status := c.Query("status")
	docType := c.Query("type")
	companyID := c.Query("company_id")
	approvalStatus := c.Query("approval_status")

	page, err := strconv.ParseInt(pageStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, errors.NewValidationError("invalid page number"))
		return
	}

	perPage, err := strconv.ParseInt(perPageStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, errors.NewValidationError("invalid per_page number"))
		return
	}

	// Limit page size to 100
	if perPage > 100 {
		perPage = 100
	}

	if page < 1 {
		page = 1
	}

	resp, err := h.service.ListDocuments(c.Request.Context(), conversion.SafeInt64ToInt32WithDefault(page, 1), conversion.SafeInt64ToInt32WithDefault(perPage, 10), status, docType, companyID, approvalStatus)
	if err != nil {
		h.logger.Error("Failed to list documents",
			zap.Int64("page", page),
			zap.Int64("per_page", perPage),
			zap.Error(err))
		c.JSON(http.StatusInternalServerError, errors.NewInternalServerError())
		return
	}

	c.JSON(http.StatusOK, resp)
}

// SearchDocuments выполняет поиск документов
//
//	@Summary		Search documents
//	@Description	Search documents by query with optional filters
//	@Tags			Documents Query
//	@Accept			json
//	@Produce		json
//	@Param			q					query		string	true	"Search query"
//	@Param			page				query		int		false	"Page number"	default(1)
//	@Param			per_page			query		int		false	"Items per page"	default(10)
//	@Param			type				query		string	false	"Document type filter"
//	@Param			company_id			query		string	false	"Company ID filter"
//	@Param			approval_status		query		string	false	"Approval status filter"
//	@Success		200		{object}	models.DocumentListResponse
//	@Failure		400		{object}	models.APIResponse	"Invalid parameters"
//	@Failure		500		{object}	models.APIResponse	"Internal server error"
//	@Router			/documents-query/search [get]
//	@Security		BearerAuth
func (h *DocumentQueryHandler) SearchDocuments(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, errors.NewValidationError("search query 'q' is required"))
		return
	}

	pageStr := c.DefaultQuery("page", "1")
	perPageStr := c.DefaultQuery("per_page", "10")
	docType := c.Query("type")
	companyID := c.Query("company_id")
	approvalStatus := c.Query("approval_status")

	page, err := strconv.ParseInt(pageStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, errors.NewValidationError("invalid page number"))
		return
	}

	perPage, err := strconv.ParseInt(perPageStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, errors.NewValidationError("invalid per_page number"))
		return
	}

	// Limit page size to 100
	if perPage > 100 {
		perPage = 100
	}

	if page < 1 {
		page = 1
	}

	resp, err := h.service.SearchDocuments(c.Request.Context(), query, conversion.SafeInt64ToInt32WithDefault(page, 1), conversion.SafeInt64ToInt32WithDefault(perPage, 10), docType, companyID, approvalStatus)
	if err != nil {
		h.logger.Error("Failed to search documents",
			zap.String("query", query),
			zap.Int64("page", page),
			zap.Error(err))
		c.JSON(http.StatusInternalServerError, errors.NewInternalServerError())
		return
	}

	c.JSON(http.StatusOK, resp)
}

// GetPendingApprovalDocuments получает документы ожидающие одобрения
//
//	@Summary		Get pending approval documents
//	@Description	Get documents pending approval for a company
//	@Tags			Documents Query
//	@Accept			json
//	@Produce		json
//	@Param			company_id		query		string	true	"Company ID"
//	@Param			page			query		int		false	"Page number"	default(1)
//	@Param			per_page		query		int		false	"Items per page"	default(10)
//	@Success		200		{object}	models.DocumentListResponse
//	@Failure		400		{object}	models.APIResponse	"Invalid parameters"
//	@Failure		500		{object}	models.APIResponse	"Internal server error"
//	@Router			/documents-query/pending-approval [get]
//	@Security		BearerAuth
func (h *DocumentQueryHandler) GetPendingApprovalDocuments(c *gin.Context) {
	companyID := c.Query("company_id")
	if companyID == "" {
		c.JSON(http.StatusBadRequest, errors.NewValidationError("company_id is required"))
		return
	}

	pageStr := c.DefaultQuery("page", "1")
	perPageStr := c.DefaultQuery("per_page", "10")

	page, err := strconv.ParseInt(pageStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, errors.NewValidationError("invalid page number"))
		return
	}

	perPage, err := strconv.ParseInt(perPageStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, errors.NewValidationError("invalid per_page number"))
		return
	}

	// Limit page size to 100
	if perPage > 100 {
		perPage = 100
	}

	if page < 1 {
		page = 1
	}

	resp, err := h.service.GetPendingApprovalDocuments(c.Request.Context(), companyID, conversion.SafeInt64ToInt32WithDefault(page, 1), conversion.SafeInt64ToInt32WithDefault(perPage, 10))
	if err != nil {
		h.logger.Error("Failed to get pending approval documents",
			zap.String("company_id", companyID),
			zap.Int64("page", page),
			zap.Error(err))
		c.JSON(http.StatusInternalServerError, errors.NewInternalServerError())
		return
	}

	c.JSON(http.StatusOK, resp)
}
