// Файл api-gateway/internal/interfaces/http/health_handler.go содержит реализацию пакета http.
package http

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rusgainew/kkm-project-mks/api-gateway/pkg/di_container"
	"go.uber.org/zap"
)

// HealthHandler отвечает только за health check endpoints
type HealthHandler struct {
	logger    *zap.Logger
	container *di_container.Container
}

// NewHealthHandler создает новый HealthHandler
func NewHealthHandler(logger *zap.Logger, container *di_container.Container) *HealthHandler {
	return &HealthHandler{
		logger:    logger,
		container: container,
	}
}

// Health проверяет работоспособность API (liveness probe)
//
//	@Summary		Health check (liveness)
//	@Description	Returns API health status for Kubernetes liveness probe
//	@Tags			Health
//	@Produce		json
//	@Success		200	{object}	map[string]interface{}
//	@Router			/health [get]
func (h *HealthHandler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "healthy",
		"time":   time.Now().Unix(),
	})
}

// Ready проверяет готовность API к обработке запросов (readiness probe)
//
//	@Summary		Readiness check
//	@Description	Returns API readiness status with backend services check
//	@Tags			Health
//	@Produce		json
//	@Success		200	{object}	map[string]interface{}
//	@Failure		503	{object}	map[string]interface{}	"Service unavailable"
//	@Router			/ready [get]
func (h *HealthHandler) Ready(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	// Проверяем доступность backend сервисов
	services := make(map[string]string)

	// Проверяем company service
	if err := h.checkServiceHealth(ctx, "company"); err != nil {
		services["company_service"] = "unhealthy"
		h.logger.Warn("Company service unhealthy", zap.Error(err))
	} else {
		services["company_service"] = "healthy"
	}

	// Проверяем invoice service
	if err := h.checkServiceHealth(ctx, "invoice"); err != nil {
		services["invoice_service"] = "unhealthy"
		h.logger.Warn("Invoice service unhealthy", zap.Error(err))
	} else {
		services["invoice_service"] = "healthy"
	}

	// Проверяем catalog service
	if err := h.checkServiceHealth(ctx, "catalog"); err != nil {
		services["catalog_service"] = "unhealthy"
		h.logger.Warn("Catalog service unhealthy", zap.Error(err))
	} else {
		services["catalog_service"] = "healthy"
	}

	// Проверяем bank account service
	if err := h.checkServiceHealth(ctx, "bank-account"); err != nil {
		services["bank_account_service"] = "unhealthy"
		h.logger.Warn("Bank account service unhealthy", zap.Error(err))
	} else {
		services["bank_account_service"] = "healthy"
	}

	// Проверяем user service
	if err := h.checkServiceHealth(ctx, "user"); err != nil {
		services["user_service"] = "unhealthy"
		h.logger.Warn("User service unhealthy", zap.Error(err))
	} else {
		services["user_service"] = "healthy"
	}

	// Проверяем foreign company service
	if err := h.checkServiceHealth(ctx, "foreign-company"); err != nil {
		services["foreign_company_service"] = "unhealthy"
		h.logger.Warn("Foreign company service unhealthy", zap.Error(err))
	} else {
		services["foreign_company_service"] = "healthy"
	}

	// Проверяем invoice query service
	if err := h.checkServiceHealth(ctx, "invoice-query"); err != nil {
		services["invoice_query_service"] = "unhealthy"
		h.logger.Warn("Invoice query service unhealthy", zap.Error(err))
	} else {
		services["invoice_query_service"] = "healthy"
	}

	// Проверяем catalog query service
	if err := h.checkServiceHealth(ctx, "catalog-query"); err != nil {
		services["catalog_query_service"] = "unhealthy"
		h.logger.Warn("Catalog query service unhealthy", zap.Error(err))
	} else {
		services["catalog_query_service"] = "healthy"
	}

	// Проверяем bank account query service
	if err := h.checkServiceHealth(ctx, "bank-account-query"); err != nil {
		services["bank_account_query_service"] = "unhealthy"
		h.logger.Warn("Bank account query service unhealthy", zap.Error(err))
	} else {
		services["bank_account_query_service"] = "healthy"
	}

	// Определяем общий статус
	allHealthy := true
	for _, status := range services {
		if status != "healthy" {
			allHealthy = false
			break
		}
	}

	statusCode := http.StatusOK
	overallStatus := "ready"
	if !allHealthy {
		statusCode = http.StatusServiceUnavailable
		overallStatus = "not_ready"
	}

	c.JSON(statusCode, gin.H{
		"status":   overallStatus,
		"services": services,
		"time":     time.Now().Unix(),
	})
}

// checkServiceHealth проверяет здоровье конкретного сервиса
func (h *HealthHandler) checkServiceHealth(ctx context.Context, serviceName string) error {
	// Получаем соединение через container
	connMgr := h.container.ConnectionManager()
	if connMgr == nil {
		h.logger.Debug("Connection manager not available", zap.String("service", serviceName))
		return nil // Не блокируем если нет connection manager
	}

	// Map имён сервисов на адреса (дефолтные)
	serviceAddresses := map[string]string{
		"company":            "localhost:9001",
		"invoice":            "localhost:9002",
		"catalog":            "localhost:9003",
		"bank-account":       "localhost:9004",
		"user":               "localhost:9005",
		"foreign-company":    "localhost:9006",
		"invoice-query":      "localhost:9007",
		"catalog-query":      "localhost:9008",
		"bank-account-query": "localhost:9009",
	}

	address, exists := serviceAddresses[serviceName]
	if !exists || address == "" {
		h.logger.Debug("Service address not found", zap.String("service", serviceName))
		return nil
	}

	// Пытаемся получить соединение с таймаутом
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	conn, err := connMgr.GetConnection(ctx, address)
	if err != nil {
		h.logger.Debug("Failed to get connection to service",
			zap.String("service", serviceName),
			zap.String("address", address),
			zap.Error(err))
		return err
	}

	// Проверяем что соединение доступно через State
	state := conn.GetState()
	if state.String() == "SHUTDOWN" {
		return context.DeadlineExceeded
	}

	h.logger.Debug("Service health check passed",
		zap.String("service", serviceName),
		zap.String("state", state.String()))

	return nil
}
