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

// ForeignCompanyQueryHandler обработчик для чтения иностранных компаний
type ForeignCompanyQueryHandler struct {
	service *services.ForeignCompanyQueryService
	logger  *zap.Logger
}

// NewForeignCompanyQueryHandler создает новый ForeignCompanyQueryHandler
func NewForeignCompanyQueryHandler(service *services.ForeignCompanyQueryService, logger *zap.Logger) *ForeignCompanyQueryHandler {
	return &ForeignCompanyQueryHandler{
		service: service,
		logger:  logger,
	}
}

// ListForeignCompanies получает список иностранных компаний с пагинацией
// @Summary		List foreign companies
// @Description	Get a paginated list of foreign companies
// @Tags			ForeignCompanies
// @Accept			json
// @Produce		json
// @Param			page		query		int	false	"Page number"	default(0)
// @Param			page_size	query		int	false	"Page size"		default(20)
// @Success		200			{object}	models.APIResponse{data=[]models.ForeignCompany}
// @Failure		400			{object}	models.APIResponse	"Invalid parameters"
// @Failure		500			{object}	models.APIResponse	"Internal server error"
// @Router			/foreign-companies-query [get]
// @Security		BearerAuth
func (h *ForeignCompanyQueryHandler) ListForeignCompanies(c *gin.Context) {
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

	resp, err := h.service.ListForeignCompanies(c.Request.Context(), conversion.SafeInt64ToInt32WithDefault(page, 0), conversion.SafeInt64ToInt32WithDefault(pageSize, 20))
	if err != nil {
		h.logger.Error("Failed to list foreign companies", zap.Error(err))
		c.JSON(http.StatusInternalServerError, errors.NewInternalServerError())
		return
	}

	c.JSON(http.StatusOK, resp)
}

// SearchForeignCompanies осуществляет поиск иностранных компаний
// @Summary		Search foreign companies
// @Description	Search foreign companies by query
// @Tags			ForeignCompanies
// @Accept			json
// @Produce		json
// @Param			q			query		string	true	"Search query"
// @Success		200			{object}	models.APIResponse{data=[]models.ForeignCompany}
// @Failure		400			{object}	models.APIResponse	"Invalid parameters"
// @Failure		500			{object}	models.APIResponse	"Internal server error"
// @Router			/foreign-companies-query/search [get]
// @Security		BearerAuth
func (h *ForeignCompanyQueryHandler) SearchForeignCompanies(c *gin.Context) {
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

	resp, err := h.service.SearchForeignCompanies(c.Request.Context(), searchReq)
	if err != nil {
		h.logger.Error("Failed to search foreign companies", zap.Error(err))
		c.JSON(http.StatusInternalServerError, errors.NewInternalServerError())
		return
	}

	c.JSON(http.StatusOK, resp)
}

// FilterForeignCompanies фильтрует иностранные компании
// @Summary		Filter foreign companies
// @Description	Filter foreign companies with multiple criteria
// @Tags			ForeignCompanies
// @Accept			json
// @Produce		json
// @Param			name		query		string	false	"Company name"
// @Param			country		query		string	false	"Country"
// @Param			page		query		int		false	"Page number"	default(0)
// @Param			page_size	query		int		false	"Page size"		default(20)
// @Success		200			{object}	models.APIResponse{data=[]models.ForeignCompany}
// @Failure		400			{object}	models.APIResponse	"Invalid parameters"
// @Failure		500			{object}	models.APIResponse	"Internal server error"
// @Router			/foreign-companies-query/filter [get]
// @Security		BearerAuth
func (h *ForeignCompanyQueryHandler) FilterForeignCompanies(c *gin.Context) {
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

	resp, err := h.service.FilterForeignCompanies(c.Request.Context(), filterReq)
	if err != nil {
		h.logger.Error("Failed to filter foreign companies", zap.Error(err))
		c.JSON(http.StatusInternalServerError, errors.NewInternalServerError())
		return
	}

	c.JSON(http.StatusOK, resp)
}
