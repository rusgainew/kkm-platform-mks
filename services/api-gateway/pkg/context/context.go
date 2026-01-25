package context

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	// Context keys
	UserIDKey    = "user_id"
	UsernameKey  = "username"
	EmailKey     = "email"
	RolesKey     = "roles"
	TraceIDKey   = "trace_id"
	RequestIDKey = "request_id"
)

// GetUserID извлекает user_id из контекста Gin
func GetUserID(c *gin.Context) (string, bool) {
	userID, exists := c.Get(UserIDKey)
	if !exists {
		return "", false
	}
	id, ok := userID.(string)
	return id, ok
}

// GetUsername извлекает username из контекста Gin
func GetUsername(c *gin.Context) (string, bool) {
	username, exists := c.Get(UsernameKey)
	if !exists {
		return "", false
	}
	name, ok := username.(string)
	return name, ok
}

// GetEmail извлекает email из контекста Gin
func GetEmail(c *gin.Context) (string, bool) {
	email, exists := c.Get(EmailKey)
	if !exists {
		return "", false
	}
	em, ok := email.(string)
	return em, ok
}

// GetRoles извлекает roles из контекста Gin
func GetRoles(c *gin.Context) ([]string, bool) {
	roles, exists := c.Get(RolesKey)
	if !exists {
		return nil, false
	}
	r, ok := roles.([]string)
	return r, ok
}

// GetTraceID извлекает trace_id из контекста Gin
func GetTraceID(c *gin.Context) string {
	traceID, exists := c.Get(TraceIDKey)
	if !exists {
		return ""
	}
	id, ok := traceID.(string)
	if !ok {
		return ""
	}
	return id
}

// GetRequestID извлекает request_id из контекста Gin
func GetRequestID(c *gin.Context) string {
	requestID, exists := c.Get(RequestIDKey)
	if !exists {
		return ""
	}
	id, ok := requestID.(string)
	if !ok {
		return ""
	}
	return id
}

// SetUserID устанавливает user_id в контекст Gin
func SetUserID(c *gin.Context, userID string) {
	c.Set(UserIDKey, userID)
}

// SetUsername устанавливает username в контекст Gin
func SetUsername(c *gin.Context, username string) {
	c.Set(UsernameKey, username)
}

// SetEmail устанавливает email в контекст Gin
func SetEmail(c *gin.Context, email string) {
	c.Set(EmailKey, email)
}

// SetRoles устанавливает roles в контекст Gin
func SetRoles(c *gin.Context, roles []string) {
	c.Set(RolesKey, roles)
}

// SetTraceID устанавливает trace_id в контекст Gin
func SetTraceID(c *gin.Context, traceID string) {
	c.Set(TraceIDKey, traceID)
}

// SetRequestID устанавливает request_id в контекст Gin
func SetRequestID(c *gin.Context, requestID string) {
	c.Set(RequestIDKey, requestID)
}

// GenerateRequestID генерирует новый request ID
func GenerateRequestID() string {
	return uuid.New().String()
}

// PropagateToContext копирует значения из Gin context в Go context
func PropagateToContext(ginCtx *gin.Context, ctx context.Context) context.Context {
	if userID, ok := GetUserID(ginCtx); ok {
		ctx = context.WithValue(ctx, UserIDKey, userID)
	}
	if username, ok := GetUsername(ginCtx); ok {
		ctx = context.WithValue(ctx, UsernameKey, username)
	}
	if email, ok := GetEmail(ginCtx); ok {
		ctx = context.WithValue(ctx, EmailKey, email)
	}
	if roles, ok := GetRoles(ginCtx); ok {
		ctx = context.WithValue(ctx, RolesKey, roles)
	}
	if traceID := GetTraceID(ginCtx); traceID != "" {
		ctx = context.WithValue(ctx, TraceIDKey, traceID)
	}
	if requestID := GetRequestID(ginCtx); requestID != "" {
		ctx = context.WithValue(ctx, RequestIDKey, requestID)
	}
	return ctx
}
