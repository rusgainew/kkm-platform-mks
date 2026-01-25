package middleware

import (
	"context"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	grpchandler "github.com/rusgainew/kkm-project-mks/company-server/internal/interfaces/grpc"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// AuthMiddleware validates JWT and injects user-id into context
type AuthMiddleware struct {
	secret string
	logger *zap.Logger
}

func NewAuthMiddleware(secret string, logger *zap.Logger) *AuthMiddleware {
	return &AuthMiddleware{secret: secret, logger: logger}
}

func (a *AuthMiddleware) UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		md, _ := metadata.FromIncomingContext(ctx)
		var token string
		if md != nil {
			vals := md.Get("authorization")
			if len(vals) > 0 {
				h := vals[0]
				if strings.HasPrefix(strings.ToLower(h), "bearer ") {
					token = strings.TrimSpace(h[len("bearer "):])
				}
			}
		}

		if token == "" {
			// Allow requests without auth for health and reflection
			if strings.Contains(info.FullMethod, "grpc.health.v1.Health") || strings.Contains(info.FullMethod, "ServerReflection") {
				return handler(ctx, req)
			}
			a.logger.Warn("missing authorization header")
			return nil, status.Error(codes.Unauthenticated, "missing authorization header")
		}

		parsed, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
			return []byte(a.secret), nil
		})
		if err != nil || !parsed.Valid {
			a.logger.Warn("invalid token", zap.Error(err))
			return nil, status.Error(codes.Unauthenticated, "invalid token")
		}

		claims, ok := parsed.Claims.(jwt.MapClaims)
		if !ok {
			return nil, status.Error(codes.Unauthenticated, "invalid token claims")
		}

		// Try user_id first, fallback to sub (subject) claim
		userIDVal, _ := claims["user_id"].(string)
		if userIDVal == "" {
			userIDVal, _ = claims["sub"].(string)
		}
		if userIDVal == "" {
			return nil, status.Error(codes.Unauthenticated, "user_id claim missing")
		}

		// Inject user-id into context for handlers
		ctx = grpchandler.SetUserID(ctx, userIDVal)
		return handler(ctx, req)
	}
}
