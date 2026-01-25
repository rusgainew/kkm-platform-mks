package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// SecurityHeadersMiddleware добавляет security headers к ответам
func SecurityHeadersMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Предотвращает MIME type sniffing
		c.Header("X-Content-Type-Options", "nosniff")

		// Предотвращает clickjacking
		c.Header("X-Frame-Options", "DENY")

		// Включает XSS фильтр браузера
		c.Header("X-XSS-Protection", "1; mode=block")

		// Контролирует какие ресурсы могут загружаться
		c.Header("Content-Security-Policy", "default-src 'self'; script-src 'self' 'unsafe-inline'; style-src 'self' 'unsafe-inline'")

		// Контролирует информацию в Referer header
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")

		// Требует HTTPS (в production)
		// c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains")

		// Запрещает кеширование для API endpoints
		c.Header("Cache-Control", "no-store, no-cache, must-revalidate, proxy-revalidate")
		c.Header("Pragma", "no-cache")
		c.Header("Expires", "0")

		// Permissions Policy (бывший Feature-Policy)
		c.Header("Permissions-Policy", "geolocation=(), microphone=(), camera=()")

		c.Next()
	}
}

// ContentTypeMiddleware проверяет Content-Type для запросов с телом
func ContentTypeMiddleware(requiredContentType string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Проверяем только для методов с телом
		if c.Request.Method == "POST" || c.Request.Method == "PUT" || c.Request.Method == "PATCH" {
			contentType := c.GetHeader("Content-Type")
			if contentType == "" {
				c.JSON(http.StatusUnsupportedMediaType, gin.H{
					"error":   "unsupported_media_type",
					"message": "Content-Type header is required",
				})
				c.Abort()
				return
			}

			// Проверяем что Content-Type содержит требуемый тип
			// Используем Contains для поддержки charset и других параметров
			if len(contentType) < len(requiredContentType) ||
				contentType[:len(requiredContentType)] != requiredContentType {
				c.JSON(http.StatusUnsupportedMediaType, gin.H{
					"error":   "unsupported_media_type",
					"message": "Content-Type must be " + requiredContentType,
				})
				c.Abort()
				return
			}
		}
		c.Next()
	}
}

// RequestSizeLimitMiddleware ограничивает размер тела запроса
func RequestSizeLimitMiddleware(maxBytes int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.ContentLength > maxBytes {
			c.JSON(http.StatusRequestEntityTooLarge, gin.H{
				"error":   "request_too_large",
				"message": "Request body too large",
			})
			c.Abort()
			return
		}

		// Устанавливаем лимит на reader
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxBytes)
		c.Next()
	}
}

// ValidateRequestMiddleware проверяет базовые параметры запроса
func ValidateRequestMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Проверяем что для POST/PUT/PATCH есть тело
		if c.Request.Method == "POST" || c.Request.Method == "PUT" || c.Request.Method == "PATCH" {
			if c.Request.ContentLength == 0 && c.GetHeader("Content-Length") != "" {
				c.JSON(http.StatusBadRequest, gin.H{
					"error":   "empty_body",
					"message": "Request body cannot be empty",
				})
				c.Abort()
				return
			}
		}
		c.Next()
	}
}
