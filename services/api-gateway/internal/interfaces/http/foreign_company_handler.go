// Файл api-gateway/internal/interfaces/http/foreign_company_handler.go содержит реализацию пакета http.
package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rusgainew/kkm-project-mks/api-gateway/internal/application/services"
	"github.com/rusgainew/kkm-project-mks/api-gateway/internal/domain/models"
	"go.uber.org/zap"
)

// ForeignCompanyHandler обработчик HTTP запросов для иностранных компаний
type ForeignCompanyHandler struct {
	*BaseHandler
	service *services.ForeignCompanyService
}

// NewForeignCompanyHandler создает новый ForeignCompanyHandler
func NewForeignCompanyHandler(service *services.ForeignCompanyService, logger *zap.Logger) *ForeignCompanyHandler {
	return &ForeignCompanyHandler{
		BaseHandler: NewBaseHandler(logger),
		service:     service,
	}
}

// CreateForeignCompanyRequest запрос на создание иностранной компании
type CreateForeignCompanyRequest struct {
	PIN         string `json:"pin" binding:"required"`
	FullName    string `json:"full_name" binding:"required"`
	CountryCode string `json:"country_code" binding:"required,len=2"`
	Address     string `json:"address"`
}

// UpdateForeignCompanyRequest запрос на обновление иностранной компании
type UpdateForeignCompanyRequest struct {
	PIN         string `json:"pin" binding:"required"`
	FullName    string `json:"full_name" binding:"required"`
	CountryCode string `json:"country_code" binding:"required,len=2"`
	Address     string `json:"address"`
}

// CreateForeignCompany создает новую иностранную компанию
//
//	@Summary		Create foreign company
//	@Description	Create a new foreign company. Requires admin or manager role.
//	@Tags			Foreign Companies
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			request	body		CreateForeignCompanyRequest						true	"Foreign company data including PIN, full name, country code (ISO 3166-1 alpha-2)"
//	@Success		201		{object}	models.APIResponse{data=models.ForeignCompany}	"Foreign company successfully created"
//	@Failure		400		{object}	models.APIResponse								"Invalid input - validation errors"
//	@Failure		401		{object}	models.APIResponse								"Unauthorized"
//	@Failure		403		{object}	models.APIResponse								"Forbidden - requires admin or manager role"
//	@Failure		500		{object}	models.APIResponse								"Internal server error"
//	@Router			/foreign-companies [post]
func (h *ForeignCompanyHandler) CreateForeignCompany(c *gin.Context) {
	var req CreateForeignCompanyRequest
	if !h.BindJSON(c, &req) {
		return
	}

	// Валидация country code (ISO 3166-1 alpha-2)
	if len(req.CountryCode) != 2 {
		h.RespondWithValidationError(c, "country_code must be 2 characters (ISO 3166-1 alpha-2)")
		return
	}

	company, err := h.service.CreateForeignCompany(
		c.Request.Context(),
		req.PIN,
		req.FullName,
		req.CountryCode,
		req.Address,
	)
	if err != nil {
		h.RespondWithGRPCError(c, err, "Create foreign company")
		return
	}

	c.JSON(http.StatusCreated, models.APIResponse{
		Success: true,
		Data:    company,
	})
}

// GetForeignCompany возвращает иностранную компанию по ID
//
//	@Summary		Get foreign company by ID
//	@Description	Retrieve detailed information about a foreign company by its ID
//	@Tags			Foreign Companies
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		int												true	"Foreign Company ID"
//	@Success		200	{object}	models.APIResponse{data=models.ForeignCompany}	"Foreign company information"
//	@Failure		400	{object}	models.APIResponse								"Invalid ID"
//	@Failure		401	{object}	models.APIResponse								"Unauthorized"
//	@Failure		404	{object}	models.APIResponse								"Foreign company not found"
//	@Failure		500	{object}	models.APIResponse								"Internal server error"
//	@Router			/foreign-companies/{id} [get]
func (h *ForeignCompanyHandler) GetForeignCompany(c *gin.Context) {
	id, ok := h.GetIntPathParam(c, "id", "foreign_company_id")
	if !ok {
		return
	}

	company, err := h.service.GetForeignCompany(c.Request.Context(), id)
	if err != nil {
		h.RespondWithGRPCError(c, err, "Get foreign company")
		return
	}

	h.RespondWithSuccess(c, company)
}

// UpdateForeignCompany обновляет иностранную компанию
//
//	@Summary	Обновление иностранной компании
//	@Tags		foreign-companies
//	@Accept		json
//	@Produce	json
//	@Security	BearerAuth
//	@Param		id		path		int							true	"Foreign Company ID"
//	@Param		request	body		UpdateForeignCompanyRequest	true	"Данные для обновления"
//	@Success	200		{object}	models.APIResponse{data=models.ForeignCompany}
//	@Failure	400		{object}	models.APIResponse	"Invalid input"
//	@Failure	401		{object}	models.APIResponse	"Unauthorized"
//	@Failure	404		{object}	models.APIResponse	"Not found"
//	@Failure	500		{object}	models.APIResponse	"Internal server error"
//	@Router		/api/v1/foreign-companies/{id} [put]
func (h *ForeignCompanyHandler) UpdateForeignCompany(c *gin.Context) {
	id, ok := h.GetIntPathParam(c, "id", "foreign_company_id")
	if !ok {
		return
	}

	var req UpdateForeignCompanyRequest
	if !h.BindJSON(c, &req) {
		return
	}

	// Валидация country code
	if len(req.CountryCode) != 2 {
		h.RespondWithValidationError(c, "country_code must be 2 characters (ISO 3166-1 alpha-2)")
		return
	}

	company, err := h.service.UpdateForeignCompany(
		c.Request.Context(),
		id,
		req.PIN,
		req.FullName,
		req.CountryCode,
		req.Address,
	)
	if err != nil {
		h.RespondWithGRPCError(c, err, "Update foreign company")
		return
	}

	h.RespondWithSuccess(c, company)
}

// DeleteForeignCompany удаляет иностранную компанию
//
//	@Summary	Удаление иностранной компании
//	@Tags		foreign-companies
//	@Produce	json
//	@Security	BearerAuth
//	@Param		id	path		int					true	"Foreign Company ID"
//	@Success	200	{object}	models.APIResponse	"Foreign company deleted successfully"
//	@Failure	400	{object}	models.APIResponse	"Invalid ID"
//	@Failure	404	{object}	models.APIResponse	"Not found"
//	@Failure	500	{object}	models.APIResponse	"Internal server error"
//	@Router		/api/v1/foreign-companies/{id} [delete]
func (h *ForeignCompanyHandler) DeleteForeignCompany(c *gin.Context) {
	id, ok := h.GetIntPathParam(c, "id", "foreign_company_id")
	if !ok {
		return
	}

	err := h.service.DeleteForeignCompany(c.Request.Context(), id)
	if err != nil {
		h.RespondWithGRPCError(c, err, "Delete foreign company")
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data:    map[string]string{"message": "Foreign company deleted successfully"},
	})
}
