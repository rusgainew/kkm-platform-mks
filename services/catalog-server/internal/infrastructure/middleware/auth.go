// Файл catalog-server/internal/infrastructure/middleware/auth.go содержит реализацию пакета middleware.
package middleware

import (
	"context"
	"fmt"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type contextKey string

const (
	userIDKey    contextKey = "user_id"
	userEmailKey contextKey = "user_email"
)

// AuthMiddleware обрабатывает JWT аутентификацию для gRPC запросов
type AuthMiddleware struct {
	jwtSecret []byte
	logger    *zap.Logger
}

// NewAuthMiddleware создает новый middleware для аутентификации
func NewAuthMiddleware(jwtSecret string, logger *zap.Logger) *AuthMiddleware {
	return &AuthMiddleware{
		jwtSecret: []byte(jwtSecret),
		logger:    logger,
	}
}

// UnaryServerInterceptor возвращает gRPC interceptor для унарных вызовов
func (m *AuthMiddleware) UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		// Пропускаем аутентификацию для health checks
		if strings.HasSuffix(info.FullMethod, "/grpc.health.v1.Health/Check") ||
			strings.HasSuffix(info.FullMethod, "/grpc.health.v1.Health/Watch") {
			return handler(ctx, req)
		}

		// Извлекаем metadata из контекста
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			m.logger.Warn("missing metadata in request")
			return nil, status.Error(codes.Unauthenticated, "missing metadata")
		}

		// Получаем токен из authorization header
		authHeaders := md.Get("authorization")
		if len(authHeaders) == 0 {
			m.logger.Warn("missing authorization header")
			return nil, status.Error(codes.Unauthenticated, "missing authorization header")
		}

		// Извлекаем токен (формат: "Bearer <token>")
		tokenString := strings.TrimPrefix(authHeaders[0], "Bearer ")
		if tokenString == authHeaders[0] {
			m.logger.Warn("invalid authorization header format")
			return nil, status.Error(codes.Unauthenticated, "invalid authorization header format")
		}

		// Парсим и валидируем JWT токен
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Проверяем метод подписи
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return m.jwtSecret, nil
		})

		if err != nil {
			m.logger.Warn("invalid token", zap.Error(err))
			return nil, status.Error(codes.Unauthenticated, "invalid token")
		}

		if !token.Valid {
			m.logger.Warn("token is not valid")
			return nil, status.Error(codes.Unauthenticated, "invalid token")
		}

		// Извлекаем claims
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			m.logger.Warn("invalid token claims")
			return nil, status.Error(codes.Unauthenticated, "invalid token claims")
		}

		// Извлекаем user_id (с fallback на sub claim для совместимости)
		userIDStr, ok := claims["user_id"].(string)
		if !ok {
			// Fallback на "sub" claim (стандартный JWT claim для subject/user ID)
			userIDStr, ok = claims["sub"].(string)
			if !ok {
				m.logger.Warn("missing user_id/sub in token")
				return nil, status.Error(codes.Unauthenticated, "missing user_id in token")
			}
		}

		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			m.logger.Warn("invalid user_id format", zap.Error(err))
			return nil, status.Error(codes.Unauthenticated, "invalid user_id format")
		}

		// Извлекаем email (опционально)
		email, _ := claims["email"].(string)

		// Добавляем user_id и email в контекст
		ctx = context.WithValue(ctx, userIDKey, userID)
		if email != "" {
			ctx = context.WithValue(ctx, userEmailKey, email)
		}

		m.logger.Debug("authenticated request",
			zap.String("user_id", userID.String()),
			zap.String("email", email),
			zap.String("method", info.FullMethod),
		)

		// Вызываем handler с обновленным контекстом
		return handler(ctx, req)
	}
}

// GetUserIDFromContext извлекает user_id из контекста
func GetUserIDFromContext(ctx context.Context) (uuid.UUID, error) {
	userID, ok := ctx.Value(userIDKey).(uuid.UUID)
	if !ok {
		return uuid.Nil, status.Error(codes.Unauthenticated, "user_id not found in context")
	}
	return userID, nil
}

// GetUserEmailFromContext извлекает email из контекста
func GetUserEmailFromContext(ctx context.Context) string {
	email, _ := ctx.Value(userEmailKey).(string)
	return email
}
