// Файл api-gateway/internal/interfaces/http/base_handler.go содержит реализацию пакета http.
package http

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rusgainew/kkm-project-mks/api-gateway/internal/domain/errors"
	"github.com/rusgainew/kkm-project-mks/api-gateway/internal/domain/models"
	"github.com/rusgainew/kkm-project-mks/api-gateway/internal/infrastructure/validation"
	"go.uber.org/zap"
)

// BaseHandler предоставляет общие методы для всех handlers
type BaseHandler struct {
	logger    *zap.Logger
	validator *validation.Validator
}

// NewBaseHandler создает новый базовый handler
func NewBaseHandler(logger *zap.Logger) *BaseHandler {
	return &BaseHandler{
		logger:    logger,
		validator: validation.NewValidator(),
	}
}

// BindJSON пытается привязать JSON и возвращает true если успешно
func (h *BaseHandler) BindJSON(c *gin.Context, obj interface{}) bool {
	if err := c.ShouldBindJSON(obj); err != nil {
		h.logger.Warn("Invalid request body", zap.Error(err))
		h.RespondWithValidationError(c, err.Error())
		return false
	}
	return true
}

// RespondWithError возвращает ошибку с указанным статусом
func (h *BaseHandler) RespondWithError(c *gin.Context, statusCode int, apiErr *errors.APIError) {
	c.JSON(statusCode, models.APIResponse{
		Success: false,
		Error: &models.APIError{
			Code:    apiErr.Code,
			Message: apiErr.Message,
			Details: apiErr.Details,
		},
	})
}

// RespondWithValidationError возвращает ошибку валидации (400)
func (h *BaseHandler) RespondWithValidationError(c *gin.Context, details string) {
	c.JSON(http.StatusBadRequest, models.APIResponse{
		Success: false,
		Error: &models.APIError{
			Code:    "INVALID_INPUT",
			Message: "Invalid request body",
			Details: details,
		},
	})
}

// RespondWithGRPCError обрабатывает gRPC ошибку и возвращает соответствующий HTTP статус
func (h *BaseHandler) RespondWithGRPCError(c *gin.Context, err error, operation string) {
	statusCode, apiErr := errors.MapGRPCErrorToHTTP(err)
	h.logger.Error(operation+" failed",
		zap.Error(err),
		zap.Int("status_code", statusCode),
		zap.String("error_code", apiErr.Code),
	)
	h.RespondWithError(c, statusCode, apiErr)
}

// RespondWithSuccess возвращает успешный ответ (200)
func (h *BaseHandler) RespondWithSuccess(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data:    data,
	})
}

// RespondWithCreated возвращает успешный ответ с 201 Created
func (h *BaseHandler) RespondWithCreated(c *gin.Context, data interface{}) {
	c.JSON(http.StatusCreated, models.APIResponse{
		Success: true,
		Data:    data,
	})
}

// RespondWithNoContent возвращает 204 No Content
func (h *BaseHandler) RespondWithNoContent(c *gin.Context) {
	c.Status(http.StatusNoContent)
}

// RequirePathParam извлекает параметр пути и возвращает его, если он есть
func (h *BaseHandler) RequirePathParam(c *gin.Context, param, fieldName string) (string, bool) {
	value := c.Param(param)
	if value == "" {
		h.RespondWithValidationError(c, fieldName+" is required")
		return "", false
	}
	return value, true
}

// ValidateUUID проверяет UUID и возвращает true если валидный
func (h *BaseHandler) ValidateUUID(c *gin.Context, id, fieldName string) bool {
	if err := h.validator.ValidateUUID(id, fieldName); err != nil {
		h.logger.Warn("Invalid UUID",
			zap.String("id", id),
			zap.String("field", fieldName),
			zap.Error(err),
		)
		h.RespondWithValidationError(c, err.Error())
		return false
	}
	return true
}

// GetPagination извлекает параметры пагинации из query string
func (h *BaseHandler) GetPagination(c *gin.Context) (page, pageSize int, ok bool) {
	page = 1
	pageSize = 20

	if p := c.Query("page"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil && parsed > 0 {
			page = parsed
		} else if err != nil {
			h.RespondWithValidationError(c, "page must be a positive integer")
			return 0, 0, false
		}
	}

	if ps := c.Query("page_size"); ps != "" {
		if parsed, err := strconv.Atoi(ps); err == nil && parsed > 0 {
			pageSize = parsed
		} else if err != nil {
			h.RespondWithValidationError(c, "page_size must be a positive integer")
			return 0, 0, false
		}
	}

	// DOS protection: ограничение максимального размера страницы
	if pageSize > 100 {
		pageSize = 100
	}

	return page, pageSize, true
}

// GetIntPathParam извлекает целочисленный параметр пути
func (h *BaseHandler) GetIntPathParam(c *gin.Context, param, fieldName string) (int64, bool) {
	value := c.Param(param)
	if value == "" {
		h.RespondWithValidationError(c, fieldName+" is required")
		return 0, false
	}

	parsed, err := strconv.ParseInt(value, 10, 64)
	if err != nil || parsed <= 0 {
		h.RespondWithValidationError(c, fieldName+" must be a positive integer")
		return 0, false
	}

	return parsed, true
}

// Validator возвращает валидатор для использования в дочерних handlers
func (h *BaseHandler) Validator() *validation.Validator {
	return h.validator
}

// Logger возвращает логгер для использования в дочерних handlers
func (h *BaseHandler) Logger() *zap.Logger {
	return h.logger
}
