// Файл bank-account-query-server/internal/infrastructure/middleware/auth.go содержит реализацию пакета middleware.
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

func NewAuthMiddleware(jwtSecret string, logger *zap.Logger) *AuthMiddleware {
	return &AuthMiddleware{
		jwtSecret: jwtSecret,
		logger:    logger,
	}
}

// UnaryServerInterceptor returns a new unary server interceptor that performs JWT authentication
func (m *AuthMiddleware) UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		// Skip authentication for health checks
		if strings.HasPrefix(info.FullMethod, "/grpc.health.v1.Health/") {
			return handler(ctx, req)
		}

		// Get metadata from context
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			m.logger.Warn("no metadata in request")
			return nil, status.Error(codes.Unauthenticated, "missing metadata")
		}

		// Get authorization header
		authHeader := md.Get("authorization")
		if len(authHeader) == 0 {
			m.logger.Warn("no authorization header")
			return nil, status.Error(codes.Unauthenticated, "missing authorization header")
		}

		// Extract token from "Bearer <token>"
		parts := strings.Split(authHeader[0], " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			m.logger.Warn("invalid authorization header format")
			return nil, status.Error(codes.Unauthenticated, "invalid authorization header format")
		}
		tokenString := parts[1]

		// Parse and validate JWT token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Validate signing method
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(m.jwtSecret), nil
		})

		if err != nil {
			m.logger.Warn("failed to parse JWT token", zap.Error(err))
			return nil, status.Error(codes.Unauthenticated, "invalid token")
		}

		if !token.Valid {
			m.logger.Warn("invalid JWT token")
			return nil, status.Error(codes.Unauthenticated, "invalid token")
		}

		// Extract claims
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			m.logger.Warn("failed to extract JWT claims")
			return nil, status.Error(codes.Unauthenticated, "invalid token claims")
		}

		// Get user_id from claims, fallback to sub (subject)
		userID, _ := claims["user_id"].(string)
		if userID == "" {
			userID, _ = claims["sub"].(string)
		}
		if userID == "" {
			m.logger.Warn("user_id not found in JWT claims")
			return nil, status.Error(codes.Unauthenticated, "missing user_id in token")
		}

		// Add user_id to context
		ctx = context.WithValue(ctx, "user_id", userID)

		// Optionally add other claims to context
		if email, ok := claims["email"].(string); ok {
			ctx = context.WithValue(ctx, "email", email)
		}

		m.logger.Debug("authenticated request",
			zap.String("user_id", userID),
			zap.String("method", info.FullMethod))

		return handler(ctx, req)
	}
}
