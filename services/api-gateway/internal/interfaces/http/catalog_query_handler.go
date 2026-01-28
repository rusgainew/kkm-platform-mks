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

// CatalogQueryHandler обработчик для чтения каталога
type CatalogQueryHandler struct {
	service *services.CatalogQueryService
	logger  *zap.Logger
}

// NewCatalogQueryHandler создает новый CatalogQueryHandler
func NewCatalogQueryHandler(service *services.CatalogQueryService, logger *zap.Logger) *CatalogQueryHandler {
	return &CatalogQueryHandler{
		service: service,
		logger:  logger,
	}
}

// ListCatalogs получает список каталогов с пагинацией
//
//	@Summary		List catalogs
//	@Description	Get a paginated list of catalogs
//	@Tags			Catalogs
//	@Accept			json
//	@Produce		json
//	@Param			page		query		int	false	"Page number"	default(0)
//	@Param			page_size	query		int	false	"Page size"		default(20)
//	@Success		200			{object}	models.APIResponse{data=[]models.Catalog}
//	@Failure		400			{object}	models.APIResponse	"Invalid parameters"
//	@Failure		500			{object}	models.APIResponse	"Internal server error"
//	@Router			/catalogs-query [get]
//	@Security		BearerAuth
func (h *CatalogQueryHandler) ListCatalogs(c *gin.Context) {
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

	resp, err := h.service.ListCatalogs(c.Request.Context(), conversion.SafeInt64ToInt32WithDefault(page, 0), conversion.SafeInt64ToInt32WithDefault(pageSize, 20))
	if err != nil {
		h.logger.Error("Failed to list catalogs", zap.Error(err))
		c.JSON(http.StatusInternalServerError, errors.NewInternalServerError())
		return
	}

	c.JSON(http.StatusOK, resp)
}

// ListCatalogsWithFilter получает список каталогов с фильтрацией
//
//	@Summary		List catalogs with filter
//	@Description	Get a paginated list of catalogs with filter criteria
//	@Tags			Catalogs
//	@Accept			json
//	@Produce		json
//	@Param			name		query		string	false	"Name filter"
//	@Param			number		query		string	false	"Number/code filter"
//	@Param			tnved_code	query		string	false	"TNVED code filter"
//	@Param			gked_code	query		string	false	"GKED code filter"
//	@Param			search_text	query		string	false	"Search text"
//	@Param			sort_field	query		string	false	"Sort field (name, number)"
//	@Param			sort_order	query		string	false	"Sort order (ASC/DESC)"
//	@Param			page		query		int		false	"Page number"	default(0)
//	@Param			page_size	query		int		false	"Page size"		default(20)
//	@Success		200			{object}	models.APIResponse{data=[]models.Catalog}
//	@Failure		400			{object}	models.APIResponse	"Invalid parameters"
//	@Failure		500			{object}	models.APIResponse	"Internal server error"
//	@Router			/catalogs-query/filter [get]
//	@Security		BearerAuth
func (h *CatalogQueryHandler) ListCatalogsWithFilter(c *gin.Context) {
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

	filter := &pb.CatalogFilterRequest{
		Name:       c.Query("name"),
		Number:     c.Query("number"),
		TnvedCode:  c.Query("tnved_code"),
		GkedCode:   c.Query("gked_code"),
		SearchText: c.Query("search_text"),
		SortField:  c.Query("sort_field"),
		SortOrder:  c.Query("sort_order"),
		Page:       conversion.SafeInt64ToInt32WithDefault(page, 0),
		Size:       conversion.SafeInt64ToInt32WithDefault(pageSize, 20),
	}

	resp, err := h.service.ListCatalogsWithFilter(c.Request.Context(), filter)
	if err != nil {
		h.logger.Error("Failed to list catalogs with filter", zap.Error(err))
		c.JSON(http.StatusInternalServerError, errors.NewInternalServerError())
		return
	}

	c.JSON(http.StatusOK, resp)
}
