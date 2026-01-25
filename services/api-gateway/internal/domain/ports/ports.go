package ports

import (
	"context"
)

// ServiceClient интерфейс для всех gRPC клиентов
type ServiceClient interface {
	Close() error
}

// TokenClaims структура для хранения данных из JWT токена
type TokenClaims struct {
	UserID   string
	Username string
	Email    string
	Roles    []string
	Claims   map[string]interface{}
}

// AuthService интерфейс для аутентификации
type AuthService interface {
	ValidateToken(ctx context.Context, token string) (*TokenClaims, error)
	GenerateToken(ctx context.Context, userID string, claims map[string]interface{}) (string, error)
	RefreshToken(ctx context.Context, token string) (string, error)
}
