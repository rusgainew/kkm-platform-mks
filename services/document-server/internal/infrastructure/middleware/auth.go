// Файл document-server/internal/infrastructure/middleware/auth.go содержит реализацию пакета middleware.
package middleware

import (
	"context"
	"fmt"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type AuthMiddleware struct {
	jwtSecret string
	logger    *zap.Logger
}

// NewAuthMiddleware создает новый middleware аутентификации
func NewAuthMiddleware(jwtSecret string, logger *zap.Logger) *AuthMiddleware {
	return &AuthMiddleware{
		jwtSecret: jwtSecret,
		logger:    logger,
	}
}

// UnaryServerInterceptor возвращает unary interceptor для аутентификации
func (m *AuthMiddleware) UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		// Извлекаем токен из метаданных
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			m.logger.Warn("No metadata in context", zap.String("method", info.FullMethod))
			return nil, status.Error(codes.Unauthenticated, "missing metadata")
		}

		tokens := md.Get("authorization")
		if len(tokens) == 0 {
			m.logger.Warn("No authorization header", zap.String("method", info.FullMethod))
			return nil, status.Error(codes.Unauthenticated, "missing authorization header")
		}

		token := strings.TrimPrefix(tokens[0], "Bearer ")

		// Проверяем JWT токен
		claims := jwt.MapClaims{}
		_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(m.jwtSecret), nil
		})

		if err != nil {
			m.logger.Warn("Invalid token", zap.Error(err), zap.String("method", info.FullMethod))
			return nil, status.Error(codes.Unauthenticated, "invalid token")
		}

		// Добавляем claims в контекст
		newCtx := context.WithValue(ctx, "jwt_claims", claims)

		return handler(newCtx, req)
	}
}

// GetUserIDFromContext извлекает user_id из контекста
func GetUserIDFromContext(ctx context.Context) (string, error) {
	claims, ok := ctx.Value("jwt_claims").(jwt.MapClaims)
	if !ok {
		return "", fmt.Errorf("claims not found in context")
	}

	userID, ok := claims["user_id"].(string)
	if !ok {
		return "", fmt.Errorf("user_id not found in claims")
	}

	return userID, nil
}
