// Файл api-gateway/internal/infrastructure/middleware/rbac.go содержит реализацию пакета middleware.
package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	contextHelper "github.com/rusgainew/kkm-project-mks/api-gateway/pkg/context"
	"go.uber.org/zap"
)

// RoleType определяет тип роли
type RoleType string

const (
	RoleAdmin   RoleType = "admin"
	RoleManager RoleType = "manager"
	RoleUser    RoleType = "user"
)

// RBACMiddleware проверяет права доступа на основе ролей
func RBACMiddleware(allowedRoles []RoleType, logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Получаем роли пользователя из контекста
		userRoles, exists := contextHelper.GetRoles(c)
		if !exists || len(userRoles) == 0 {
			logger.Warn("No roles found in context",
				zap.String("path", c.Request.URL.Path),
			)
			c.JSON(http.StatusForbidden, gin.H{
				"error":   "forbidden",
				"message": "insufficient permissions",
			})
			c.Abort()
			return
		}

		// Проверяем есть ли хотя бы одна разрешенная роль
		hasPermission := false
		for _, userRole := range userRoles {
			for _, allowedRole := range allowedRoles {
				if userRole == string(allowedRole) {
					hasPermission = true
					break
				}
			}
			if hasPermission {
				break
			}
		}

		if !hasPermission {
			userID, _ := contextHelper.GetUserID(c)
			logger.Warn("Access denied - insufficient permissions",
				zap.String("user_id", userID),
				zap.Strings("user_roles", userRoles),
				zap.String("path", c.Request.URL.Path),
			)
			c.JSON(http.StatusForbidden, gin.H{
				"error":   "forbidden",
				"message": "insufficient permissions",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequireAdmin требует роль admin
func RequireAdmin(logger *zap.Logger) gin.HandlerFunc {
	return RBACMiddleware([]RoleType{RoleAdmin}, logger)
}

// RequireAdminOrManager требует роль admin или manager
func RequireAdminOrManager(logger *zap.Logger) gin.HandlerFunc {
	return RBACMiddleware([]RoleType{RoleAdmin, RoleManager}, logger)
}

// RequireAnyRole требует любую роль (авторизованный пользователь)
func RequireAnyRole(logger *zap.Logger) gin.HandlerFunc {
	return RBACMiddleware([]RoleType{RoleAdmin, RoleManager, RoleUser}, logger)
}
