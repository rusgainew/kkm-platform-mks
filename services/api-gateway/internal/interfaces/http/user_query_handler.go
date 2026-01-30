// Файл api-gateway/internal/interfaces/http/user_query_handler.go содержит реализацию пакета http.
package http

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rusgainew/kkm-project-mks/api-gateway/internal/application/services"
	"github.com/rusgainew/kkm-project-mks/api-gateway/internal/domain/errors"
	"github.com/rusgainew/kkm-project-mks/lib/conversion"
	"go.uber.org/zap"
)

// UserQueryHandler обработчик для чтения пользователей
type UserQueryHandler struct {
	service *services.UserQueryService
	logger  *zap.Logger
}

// NewUserQueryHandler создает новый UserQueryHandler
func NewUserQueryHandler(service *services.UserQueryService, logger *zap.Logger) *UserQueryHandler {
	return &UserQueryHandler{
		service: service,
		logger:  logger,
	}
}

// GetUser получает пользователя по ID
//
//	@Summary		Get user by ID
//	@Description	Retrieve a single user by their ID
//	@Tags			Users Query
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string	true	"User ID"
//	@Success		200		{object}	models.User
//	@Failure		404		{object}	models.APIResponse	"User not found"
//	@Failure		500		{object}	models.APIResponse	"Internal server error"
//	@Router			/users-query/{id} [get]
//	@Security		BearerAuth
func (h *UserQueryHandler) GetUser(c *gin.Context) {
	userID := c.Param("id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, errors.NewValidationError("user_id is required"))
		return
	}

	user, err := h.service.GetUser(c.Request.Context(), userID)
	if err != nil {
		h.logger.Error("Failed to get user", zap.String("user_id", userID), zap.Error(err))
		c.JSON(http.StatusInternalServerError, errors.NewInternalServerError())
		return
	}

	if user == nil {
		c.JSON(http.StatusNotFound, errors.NewNotFoundError("user not found"))
		return
	}

	c.JSON(http.StatusOK, user)
}

// ListUsers получает список пользователей с пагинацией
//
//	@Summary		List users
//	@Description	Get a paginated list of users
//	@Tags			Users Query
//	@Accept			json
//	@Produce		json
//	@Param			page		query		int		false	"Page number"	default(1)
//	@Param			per_page	query		int		false	"Items per page"	default(10)
//	@Param			status		query		string	false	"User status filter"
//	@Param			role		query		string	false	"User role filter"
//	@Success		200		{object}	models.APIResponse{data=models.ListUsersResponse}
//	@Failure		400		{object}	models.APIResponse	"Invalid parameters"
//	@Failure		500		{object}	models.APIResponse	"Internal server error"
//	@Router			/users-query [get]
//	@Security		BearerAuth
func (h *UserQueryHandler) ListUsers(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	perPageStr := c.DefaultQuery("per_page", "10")
	status := c.Query("status")
	role := c.Query("role")

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

	resp, err := h.service.ListUsers(c.Request.Context(), conversion.SafeInt64ToInt32WithDefault(page, 1), conversion.SafeInt64ToInt32WithDefault(perPage, 10), status, role)
	if err != nil {
		h.logger.Error("Failed to list users",
			zap.Int64("page", page),
			zap.Int64("per_page", perPage),
			zap.Error(err))
		c.JSON(http.StatusInternalServerError, errors.NewInternalServerError())
		return
	}

	c.JSON(http.StatusOK, resp)
}

// SearchUsers выполняет поиск пользователей
//
//	@Summary		Search users
//	@Description	Search users by query
//	@Tags			Users Query
//	@Accept			json
//	@Produce		json
//	@Param			q			query		string	true	"Search query"
//	@Param			page		query		int		false	"Page number"	default(1)
//	@Param			per_page	query		int		false	"Items per page"	default(10)
//	@Param			status		query		string	false	"User status filter"
//	@Success		200		{object}	models.ListUsersResponse
//	@Failure		400		{object}	models.APIResponse	"Invalid parameters"
//	@Failure		500		{object}	models.APIResponse	"Internal server error"
//	@Router			/users-query/search [get]
//	@Security		BearerAuth
func (h *UserQueryHandler) SearchUsers(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, errors.NewValidationError("search query 'q' is required"))
		return
	}

	pageStr := c.DefaultQuery("page", "1")
	perPageStr := c.DefaultQuery("per_page", "10")
	status := c.Query("status")

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

	resp, err := h.service.SearchUsers(c.Request.Context(), query, conversion.SafeInt64ToInt32WithDefault(page, 1), conversion.SafeInt64ToInt32WithDefault(perPage, 10), status)
	if err != nil {
		h.logger.Error("Failed to search users",
			zap.String("query", query),
			zap.Int64("page", page),
			zap.Error(err))
		c.JSON(http.StatusInternalServerError, errors.NewInternalServerError())
		return
	}

	c.JSON(http.StatusOK, resp)
}
