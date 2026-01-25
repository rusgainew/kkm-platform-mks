package grpc

import (
	"context"
	"errors"

	"google.golang.org/grpc/metadata"
)

const (
	// userIDKey - ключ для хранения userID в контексте
	userIDKey = "user-id"
	// authorizationKey - ключ для JWT токена в метаданных
	authorizationKey = "authorization"
)

var (
	// ErrNoUserID - ошибка когда userID не найден в контексте
	ErrNoUserID = errors.New("user id not found in context")
	// ErrUnauthorized - ошибка авторизации
	ErrUnauthorized = errors.New("user not authorized to perform this action")
)

// ExtractUserID извлекает ID пользователя из контекста
// Данные могут быть передены как через метаданные (из перехватчика auth),
// так и непосредственно установлены в контексте
func ExtractUserID(ctx context.Context) (string, error) {
	// Попытка получить userID из значения контекста
	if userID, ok := ctx.Value(userIDKey).(string); ok && userID != "" {
		return userID, nil
	}

	// Попытка получить userID из метаданных (для совместимости)
	// В реальном приложении это должно быть установлено перехватчиком auth
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", ErrNoUserID
	}

	userIDs := md.Get(userIDKey)
	if len(userIDs) > 0 && userIDs[0] != "" {
		return userIDs[0], nil
	}

	return "", ErrNoUserID
}

// SetUserID устанавливает ID пользователя в контекст
// Используется в перехватчиках аутентификации
func SetUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, userIDKey, userID)
}

// ValidateOwnership проверяет, что пользователь имеет право на действие
// Возвращает ошибку если userID не совпадает с ownerID
func ValidateOwnership(userID, ownerID string) error {
	if userID != ownerID {
		return ErrUnauthorized
	}
	return nil
}
