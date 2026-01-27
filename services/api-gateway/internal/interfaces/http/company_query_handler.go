package http

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rusgainew/kkm-project-mks/api-gateway/internal/application/services"
	"github.com/rusgainew/kkm-project-mks/api-gateway/internal/domain/errors"
	pb "github.com/rusgainew/kkm-project-mks/proto-lib/api"
	"github.com/rusgainew/kkm-project-mks/services/pkg/conversion"
	"go.uber.org/zap"
)

// CompanyQueryHandler обработчик для чтения компаний
type CompanyQueryHandler struct {
	service *services.CompanyQueryService
	logger  *zap.Logger
}

// NewCompanyQueryHandler создает новый CompanyQueryHandler
func NewCompanyQueryHandler(service *services.CompanyQueryService, logger *zap.Logger) *CompanyQueryHandler {
	return &CompanyQueryHandler{
		service: service,
		logger:  logger,
	}
}

// ListCompanies получает список компаний с пагинацией
// @Summary		List companies
// @Description	Get a paginated list of companies
// @Tags			Companies
// @Accept			json
// @Produce		json
// @Param			page		query		int	false	"Page number"	default(0)
// @Param			page_size	query		int	false	"Page size"		default(20)
// @Success		200			{object}	models.APIResponse{data=[]models.Company}
// @Failure		400			{object}	models.APIResponse	"Invalid parameters"
// @Failure		500			{object}	models.APIResponse	"Internal server error"
// @Router			/companies-query [get]
// @Security		BearerAuth
func (h *CompanyQueryHandler) ListCompanies(c *gin.Context) {
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

	resp, err := h.service.ListCompanies(c.Request.Context(), conversion.SafeInt64ToInt32WithDefault(page, 0), conversion.SafeInt64ToInt32WithDefault(pageSize, 20))
	if err != nil {
		h.logger.Error("Failed to list companies", zap.Error(err))
		c.JSON(http.StatusInternalServerError, errors.NewInternalServerError())
		return
	}

	c.JSON(http.StatusOK, resp)
}

// SearchCompanies осуществляет поиск компаний
// @Summary		Search companies
// @Description	Search companies by query
// @Tags			Companies
// @Accept			json
// @Produce		json
// @Param			q			query		string	true	"Search query"
// @Success		200			{object}	models.APIResponse{data=[]models.Company}
// @Failure		400			{object}	models.APIResponse	"Invalid parameters"
// @Failure		500			{object}	models.APIResponse	"Internal server error"
// @Router			/companies-query/search [get]
// @Security		BearerAuth
func (h *CompanyQueryHandler) SearchCompanies(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, errors.NewValidationError("search query is required"))
		return
	}

	searchReq := &pb.SearchRequest{
		SearchText: query,
		Page:       0,
		Size:       20,
	}

	resp, err := h.service.SearchCompanies(c.Request.Context(), searchReq)
	if err != nil {
		h.logger.Error("Failed to search companies", zap.Error(err))
		c.JSON(http.StatusInternalServerError, errors.NewInternalServerError())
		return
	}

	c.JSON(http.StatusOK, resp)
}

// FilterCompanies фильтрует компании
// @Summary		Filter companies
// @Description	Filter companies with multiple criteria
// @Tags			Companies
// @Accept			json
// @Produce		json
// @Param			name		query		string	false	"Company name"
// @Param			country		query		string	false	"Country"
// @Param			page		query		int		false	"Page number"	default(0)
// @Param			page_size	query		int		false	"Page size"		default(20)
// @Success		200			{object}	models.APIResponse{data=[]models.Company}
// @Failure		400			{object}	models.APIResponse	"Invalid parameters"
// @Failure		500			{object}	models.APIResponse	"Internal server error"
// @Router			/companies-query/filter [get]
// @Security		BearerAuth
func (h *CompanyQueryHandler) FilterCompanies(c *gin.Context) {
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

	filterReq := &pb.CatalogFilterRequest{
		Name: c.Query("name"),
		Page: conversion.SafeInt64ToInt32WithDefault(page, 0),
		Size: conversion.SafeInt64ToInt32WithDefault(pageSize, 20),
	}

	resp, err := h.service.FilterCompanies(c.Request.Context(), filterReq)
	if err != nil {
		h.logger.Error("Failed to filter companies", zap.Error(err))
		c.JSON(http.StatusInternalServerError, errors.NewInternalServerError())
		return
	}

	c.JSON(http.StatusOK, resp)
}
