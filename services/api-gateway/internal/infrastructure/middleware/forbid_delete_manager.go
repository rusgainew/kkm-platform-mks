package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	contextHelper "github.com/rusgainew/kkm-project-mks/api-gateway/pkg/context"
	"go.uber.org/zap"
)

// ForbidDeleteForManagerMiddleware forbids HTTP DELETE requests for users with role "manager" (unless also admin)
// Attach this middleware AFTER AuthMiddleware on protected routes or groups.
func ForbidDeleteForManagerMiddleware(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method != http.MethodDelete {
			c.Next()
			return
		}

		roles, ok := contextHelper.GetRoles(c)
		if !ok || len(roles) == 0 {
			// No roles in context (likely unauthenticated route) — do not enforce here
			c.Next()
			return
		}

		isManager := false
		isAdmin := false
		for _, r := range roles {
			if r == "manager" {
				isManager = true
			}
			if r == "admin" {
				isAdmin = true
			}
		}

		if isManager && !isAdmin {
			userID, _ := contextHelper.GetUserID(c)
			logger.Warn("DELETE forbidden for manager",
				zap.String("user_id", userID),
				zap.Strings("roles", roles),
				zap.String("path", c.Request.URL.Path),
			)
			c.JSON(http.StatusForbidden, gin.H{
				"error":   "forbidden",
				"message": "managers are not allowed to delete resources",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
