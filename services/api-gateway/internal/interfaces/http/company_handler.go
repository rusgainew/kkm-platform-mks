package http

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rusgainew/kkm-project-mks/api-gateway/internal/domain/errors"
	"github.com/rusgainew/kkm-project-mks/api-gateway/internal/domain/models"
	"github.com/rusgainew/kkm-project-mks/api-gateway/internal/domain/ports"
	"github.com/rusgainew/kkm-project-mks/api-gateway/internal/infrastructure/validation"
	"go.uber.org/zap"
)

// CompanyHandler отвечает только за HTTP запросы компаний
type CompanyHandler struct {
	service   ports.CompanyServiceInterface
	logger    *zap.Logger
	validator *validation.Validator
}

// NewCompanyHandler создает новый компания handler
func NewCompanyHandler(service ports.CompanyServiceInterface, logger *zap.Logger) *CompanyHandler {
	return &CompanyHandler{
		service:   service,
		logger:    logger,
		validator: validation.NewValidator(),
	}
}

// CreateCompany создает новую компанию
//
//	@Summary		Create a new company
//	@Description	Create a new company with provided details. Available to all authenticated users.
//	@Tags			Companies
//	@Accept			json
//	@Produce		json
//	@Param			company	body		models.Company							true	"Company data including PIN, full name, address"
//	@Success		201		{object}	models.APIResponse{data=models.Company}	"Company successfully created"
//	@Failure		400		{object}	models.APIResponse						"Invalid input - validation errors"
//	@Failure		401		{object}	models.APIResponse						"Unauthorized"
//	@Failure		500		{object}	models.APIResponse						"Internal server error"
//	@Router			/companies [post]
//	@Security		BearerAuth
func (h *CompanyHandler) CreateCompany(c *gin.Context) {
	var company models.Company
	if err := c.ShouldBindJSON(&company); err != nil {
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

	// Validate company data
	if err := h.validator.ValidateCompany(&company); err != nil {
		h.logger.Warn("Company validation failed", zap.Error(err))
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error: &models.APIError{
				Code:    "VALIDATION_FAILED",
				Message: "Company validation failed",
				Details: err.Error(),
			},
		})
		return
	}

	result, err := h.service.CreateCompany(c.Request.Context(), &company)
	if err != nil {
		statusCode, apiErr := errors.MapGRPCErrorToHTTP(err)
		h.logger.Error("Failed to create company", zap.Error(err))
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

// GetCompany получает компанию по ID
//
//	@Summary		Get company by ID
//	@Description	Retrieve a company by its ID
//	@Tags			Companies
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"Company ID"
//	@Success		200	{object}	models.APIResponse{data=models.Company}
//	@Failure		400	{object}	models.APIResponse	"Invalid ID"
//	@Failure		500	{object}	models.APIResponse	"Internal server error"
//	@Router			/companies/{id} [get]
//	@Security		BearerAuth
func (h *CompanyHandler) GetCompany(c *gin.Context) {
	id := c.Param("id")

	// Validate UUID
	if err := h.validator.ValidateUUID(id, "company_id"); err != nil {
		h.logger.Warn("Invalid company ID", zap.String("id", id), zap.Error(err))
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error: &models.APIError{
				Code:    "INVALID_ID",
				Message: err.Error(),
			},
		})
		return
	}

	result, err := h.service.GetCompany(c.Request.Context(), id)
	if err != nil {
		statusCode, apiErr := errors.MapGRPCErrorToHTTP(err)
		h.logger.Error("Failed to get company", zap.String("id", id), zap.Error(err))
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

// UpdateCompany обновляет компанию
//
//	@Summary		Update company
//	@Description	Update company information
//	@Tags			Companies
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string			true	"Company ID"
//	@Param			company	body		models.Company	true	"Updated company data"
//	@Success		200		{object}	models.APIResponse{data=models.Company}
//	@Failure		400		{object}	models.APIResponse	"Invalid input"
//	@Failure		500		{object}	models.APIResponse	"Internal server error"
//	@Router			/companies/{id} [put]
//	@Security		BearerAuth
func (h *CompanyHandler) UpdateCompany(c *gin.Context) {
	id := c.Param("id")

	// Validate UUID
	if err := h.validator.ValidateUUID(id, "company_id"); err != nil {
		h.logger.Warn("Invalid company ID", zap.String("id", id), zap.Error(err))
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error: &models.APIError{
				Code:    "INVALID_ID",
				Message: err.Error(),
			},
		})
		return
	}

	var company models.Company
	if err := c.ShouldBindJSON(&company); err != nil {
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

	// Validate company data
	if err := h.validator.ValidateCompany(&company); err != nil {
		h.logger.Warn("Company validation failed", zap.Error(err))
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error: &models.APIError{
				Code:    "VALIDATION_FAILED",
				Message: "Company validation failed",
				Details: err.Error(),
			},
		})
		return
	}

	company.ID = id
	result, err := h.service.UpdateCompany(c.Request.Context(), &company)
	if err != nil {
		statusCode, apiErr := errors.MapGRPCErrorToHTTP(err)
		h.logger.Error("Failed to update company", zap.String("id", id), zap.Error(err))
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

// DeleteCompany удаляет компанию
//
//	@Summary		Delete company
//	@Description	Delete a company by ID
//	@Tags			Companies
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string				true	"Company ID"
//	@Success		200	{object}	models.APIResponse	"Company deleted successfully"
//	@Failure		400	{object}	models.APIResponse	"Invalid ID"
//	@Failure		500	{object}	models.APIResponse	"Internal server error"
//	@Router			/companies/{id} [delete]
//	@Security		BearerAuth
func (h *CompanyHandler) DeleteCompany(c *gin.Context) {
	id := c.Param("id")

	// Validate UUID
	if err := h.validator.ValidateUUID(id, "company_id"); err != nil {
		h.logger.Warn("Invalid company ID", zap.String("id", id), zap.Error(err))
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error: &models.APIError{
				Code:    "INVALID_ID",
				Message: err.Error(),
			},
		})
		return
	}

	if err := h.service.DeleteCompany(c.Request.Context(), id); err != nil {
		statusCode, apiErr := errors.MapGRPCErrorToHTTP(err)
		h.logger.Error("Failed to delete company", zap.String("id", id), zap.Error(err))
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
	})
}

// ListCompanies получает список компаний с пагинацией
//
//	@Summary		List companies
//	@Description	Get a paginated list of companies
//	@Tags			Companies
//	@Accept			json
//	@Produce		json
//	@Param			page		query		int	false	"Page number"	default(1)
//	@Param			pageSize	query		int	false	"Page size"		default(10)
//	@Success		200			{object}	models.APIResponse{data=[]models.Company}
//	@Failure		500			{object}	models.APIResponse	"Internal server error"
//	@Router			/companies [get]
//	@Security		BearerAuth
func (h *CompanyHandler) ListCompanies(c *gin.Context) {
	page := 1
	if p := c.Query("page"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil && parsed > 0 {
			page = parsed
		}
	}

	pageSize := 10
	if ps := c.Query("page_size"); ps != "" {
		if parsed, err := strconv.Atoi(ps); err == nil && parsed > 0 {
			pageSize = parsed
		}
	}

	// Validate pagination
	if err := h.validator.ValidatePagination(page, pageSize); err != nil {
		h.logger.Warn("Invalid pagination", zap.Int("page", page), zap.Int("page_size", pageSize), zap.Error(err))
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error: &models.APIError{
				Code:    "INVALID_PAGINATION",
				Message: err.Error(),
			},
		})
		return
	}

	result, total, err := h.service.ListCompanies(c.Request.Context(), page, pageSize)
	if err != nil {
		statusCode, apiErr := errors.MapGRPCErrorToHTTP(err)
		h.logger.Error("Failed to list companies", zap.Error(err))
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
		Meta: &models.MetaData{
			Page:       page,
			PageSize:   pageSize,
			TotalCount: total,
			TotalPages: (total + pageSize - 1) / pageSize,
		},
	})
}

// GetOrganizationMembers получает список членов организации
//
//	@Summary		Get organization members
//	@Description	Retrieve all members of a specific organization
//	@Tags			Companies
//	@Accept			json
//	@Produce		json
//	@Param			id			path		string													true	"Organization ID"
//	@Param			page		query		int														false	"Page number"	default(1)
//	@Param			pageSize	query		int														false	"Page size"		default(10)
//	@Success		200			{object}	models.APIResponse{data=[]models.Employee}
//	@Failure		404			{object}	models.APIResponse	"Organization not found"
//	@Failure		500			{object}	models.APIResponse	"Internal server error"
//	@Router			/companies/{id}/members [get]
//	@Security		BearerAuth
func (h *CompanyHandler) GetOrganizationMembers(c *gin.Context) {
	orgID := c.Param("id")
	if orgID == "" {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error: &models.APIError{
				Code:    "INVALID_INPUT",
				Message: "Organization ID is required",
			},
		})
		return
	}

	// TODO: Implement GetOrganizationMembers method in CompanyService
	c.JSON(http.StatusNotImplemented, models.APIResponse{
		Success: false,
		Error: &models.APIError{
			Code:    "NOT_IMPLEMENTED",
			Message: "GetOrganizationMembers method is not implemented yet",
		},
	})

	/* Original implementation - commented out until service method is implemented
	page := 1
	if p := c.Query("page"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil && parsed > 0 {
			page = parsed
		}
	}

	pageSize := 10
	if ps := c.Query("page_size"); ps != "" {
		if parsed, err := strconv.Atoi(ps); err == nil && parsed > 0 {
			pageSize = parsed
		}
	}

	members, total, err := h.service.GetOrganizationMembers(c.Request.Context(), orgID, int32(page), int32(pageSize))
	if err != nil {
		statusCode, apiErr := errors.MapGRPCErrorToHTTP(err)
		h.logger.Error("Failed to get organization members", zap.String("org_id", orgID), zap.Error(err))
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
		Data:    members,
		Meta: &models.MetaData{
			Page:       page,
			PageSize:   pageSize,
			TotalCount: int(total),
			TotalPages: (int(total) + pageSize - 1) / pageSize,
		},
	})
	*/
}

// AddMember добавляет члена в организацию
//
//	@Summary		Add member to organization
//	@Description	Add a user as a member of an organization
//	@Tags			Companies
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string											true	"Organization ID"
//	@Param			member	body		models.AddMemberRequest							true	"Member data"
//	@Success		201		{object}	models.APIResponse{data=models.Employee}		"Member successfully added"
//	@Failure		400		{object}	models.APIResponse								"Invalid input"
//	@Failure		404		{object}	models.APIResponse								"Organization not found"
//	@Failure		500		{object}	models.APIResponse								"Internal server error"
//	@Router			/companies/{id}/members [post]
//	@Security		BearerAuth
func (h *CompanyHandler) AddMember(c *gin.Context) {
	orgID := c.Param("id")
	if orgID == "" {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error: &models.APIError{
				Code:    "INVALID_INPUT",
				Message: "Organization ID is required",
			},
		})
		return
	}

	// TODO: Implement AddMember method in CompanyService
	c.JSON(http.StatusNotImplemented, models.APIResponse{
		Success: false,
		Error: &models.APIError{
			Code:    "NOT_IMPLEMENTED",
			Message: "AddMember method is not implemented yet",
		},
	})

	/* Original implementation - commented out until service method is implemented
	var req models.AddMemberRequest
	if err := c.ShouldBindJSON(&req); err != nil {
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

	req.OrganizationID = orgID
	member, err := h.service.AddMember(c.Request.Context(), &req)
	if err != nil {
		statusCode, apiErr := errors.MapGRPCErrorToHTTP(err)
		h.logger.Error("Failed to add member", zap.String("org_id", orgID), zap.Error(err))
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
		Data:    member,
	})
	*/
}

// RemoveMember удаляет члена из организации
//
//	@Summary		Remove member from organization
//	@Description	Remove a user from an organization
//	@Tags			Companies
//	@Accept			json
//	@Produce		json
//	@Param			id			path		string				true	"Organization ID"
//	@Param			memberId	path		string				true	"Member/Employee ID"
//	@Success		200			{object}	models.APIResponse	"Member successfully removed"
//	@Failure		404			{object}	models.APIResponse	"Organization or member not found"
//	@Failure		500			{object}	models.APIResponse	"Internal server error"
//	@Router			/companies/{id}/members/{memberId} [delete]
//	@Security		BearerAuth
func (h *CompanyHandler) RemoveMember(c *gin.Context) {
	orgID := c.Param("id")
	memberID := c.Param("memberId")

	if orgID == "" || memberID == "" {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error: &models.APIError{
				Code:    "INVALID_INPUT",
				Message: "Organization ID and Member ID are required",
			},
		})
		return
	}

	// TODO: Implement RemoveMember method in CompanyService
	c.JSON(http.StatusNotImplemented, models.APIResponse{
		Success: false,
		Error: &models.APIError{
			Code:    "NOT_IMPLEMENTED",
			Message: "RemoveMember method is not implemented yet",
		},
	})

	/* Original implementation - commented out until service method is implemented
	err := h.service.RemoveMember(c.Request.Context(), orgID, memberID)
	if err != nil {
		statusCode, apiErr := errors.MapGRPCErrorToHTTP(err)
		h.logger.Error("Failed to remove member",
			zap.String("org_id", orgID),
			zap.String("member_id", memberID),
			zap.Error(err))
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
	})
	*/
}
