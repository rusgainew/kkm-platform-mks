// Файл api-gateway/internal/interfaces/http/catalog_handler.go содержит реализацию пакета http.
package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rusgainew/kkm-project-mks/api-gateway/internal/domain/models"
	"github.com/rusgainew/kkm-project-mks/api-gateway/internal/domain/ports"
	"github.com/rusgainew/kkm-project-mks/api-gateway/internal/infrastructure/validation"
	"go.uber.org/zap"
)

// CatalogHandler отвечает только за HTTP запросы каталога
type CatalogHandler struct {
	*BaseHandler
	service   ports.CatalogServiceInterface
	validator *validation.Validator
}

// NewCatalogHandler создает новый catalog handler
func NewCatalogHandler(service ports.CatalogServiceInterface, logger *zap.Logger) *CatalogHandler {
	return &CatalogHandler{
		BaseHandler: NewBaseHandler(logger),
		service:     service,
		validator:   validation.NewValidator(),
	}
}

// CreateCatalog создает новый элемент каталога
//
//	@Summary		Create catalog item
//	@Description	Create a new catalog item
//	@Tags			Catalog
//	@Accept			json
//	@Produce		json
//	@Param			catalog	body		models.Catalog	true	"Catalog data"
//	@Success		201		{object}	models.APIResponse{data=models.Catalog}
//	@Failure		400		{object}	models.APIResponse	"Invalid input"
//	@Failure		500		{object}	models.APIResponse	"Internal server error"
//	@Router			/catalog [post]
//	@Security		BearerAuth
func (h *CatalogHandler) CreateCatalog(c *gin.Context) {
	var catalog models.Catalog
	if !h.BindJSON(c, &catalog) {
		return
	}

	// Validate catalog data
	if err := h.validator.ValidateCatalog(&catalog); err != nil {
		h.Logger().Warn("Catalog validation failed", zap.Error(err))
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error: &models.APIError{
				Code:    "VALIDATION_FAILED",
				Message: "Catalog validation failed",
				Details: err.Error(),
			},
		})
		return
	}

	result, err := h.service.CreateCatalog(c.Request.Context(), &catalog)
	if err != nil {
		h.RespondWithGRPCError(c, err, "Create catalog")
		return
	}

	h.RespondWithCreated(c, result)
}

// GetCatalog возвращает элемент каталога по ID
//
//	@Summary		Get catalog item by ID
//	@Description	Retrieve a catalog item by its ID
//	@Tags			Catalog
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"Catalog ID"
//	@Success		200	{object}	models.APIResponse{data=models.Catalog}
//	@Failure		400	{object}	models.APIResponse	"Invalid ID"
//	@Failure		404	{object}	models.APIResponse	"Not found"
//	@Failure		500	{object}	models.APIResponse	"Internal server error"
//	@Router			/catalog/{id} [get]
//	@Security		BearerAuth
func (h *CatalogHandler) GetCatalog(c *gin.Context) {
	id := c.Param("id")
	if !h.ValidateUUID(c, id, "catalog_id") {
		return
	}

	result, err := h.service.GetCatalog(c.Request.Context(), id)
	if err != nil {
		h.RespondWithGRPCError(c, err, "Get catalog")
		return
	}

	h.RespondWithSuccess(c, result)
}

// UpdateCatalog обновляет элемент каталога
//
//	@Summary		Update catalog item
//	@Description	Update an existing catalog item
//	@Tags			Catalog
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string			true	"Catalog ID"
//	@Param			catalog	body		models.Catalog	true	"Updated catalog data"
//	@Success		200		{object}	models.APIResponse{data=models.Catalog}
//	@Failure		400		{object}	models.APIResponse	"Invalid input"
//	@Failure		404		{object}	models.APIResponse	"Not found"
//	@Failure		500		{object}	models.APIResponse	"Internal server error"
//	@Router			/catalog/{id} [put]
//	@Security		BearerAuth
func (h *CatalogHandler) UpdateCatalog(c *gin.Context) {
	id := c.Param("id")
	if !h.ValidateUUID(c, id, "catalog_id") {
		return
	}

	var catalog models.Catalog
	if !h.BindJSON(c, &catalog) {
		return
	}

	// Validate catalog data
	if err := h.validator.ValidateCatalog(&catalog); err != nil {
		h.Logger().Warn("Catalog validation failed", zap.Error(err))
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error: &models.APIError{
				Code:    "VALIDATION_FAILED",
				Message: "Catalog validation failed",
				Details: err.Error(),
			},
		})
		return
	}

	catalog.ID = id
	result, err := h.service.UpdateCatalog(c.Request.Context(), &catalog)
	if err != nil {
		h.RespondWithGRPCError(c, err, "Update catalog")
		return
	}

	h.RespondWithSuccess(c, result)
}

// DeleteCatalog удаляет элемент каталога
//
//	@Summary		Delete catalog item
//	@Description	Delete a catalog item by ID
//	@Tags			Catalog
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string				true	"Catalog ID"
//	@Success		200	{object}	models.APIResponse	"Catalog deleted successfully"
//	@Failure		400	{object}	models.APIResponse	"Invalid ID"
//	@Failure		404	{object}	models.APIResponse	"Not found"
//	@Failure		500	{object}	models.APIResponse	"Internal server error"
//	@Router			/catalog/{id} [delete]
//	@Security		BearerAuth
func (h *CatalogHandler) DeleteCatalog(c *gin.Context) {
	id := c.Param("id")
	if !h.ValidateUUID(c, id, "catalog_id") {
		return
	}

	err := h.service.DeleteCatalog(c.Request.Context(), id)
	if err != nil {
		h.RespondWithGRPCError(c, err, "Delete catalog")
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data:    map[string]string{"message": "Catalog deleted successfully"},
	})
}
