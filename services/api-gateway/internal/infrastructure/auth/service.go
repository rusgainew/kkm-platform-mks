package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/rusgainew/kkm-project-mks/api-gateway/internal/domain/ports"
	"github.com/rusgainew/kkm-project-mks/api-gateway/internal/infrastructure/cache"
	"github.com/rusgainew/kkm-project-mks/api-gateway/internal/infrastructure/observability"
	"go.uber.org/zap"
)

var (
	ErrInvalidToken  = errors.New("invalid token")
	ErrExpiredToken  = errors.New("token has expired")
	ErrInvalidClaims = errors.New("invalid token claims")
)

// Service реализация AuthService
type Service struct {
	jwtSecret   []byte
	expiration  time.Duration
	logger      *zap.Logger
	metrics     *observability.Metrics
	keyProvider *KeyProvider          // ← Для ротации ключей
	blacklist   *cache.TokenBlacklist // ← Для отозванных токенов
}

// NewService создает новый AuthService
func NewService(jwtSecret string, expiration time.Duration, logger *zap.Logger) *Service {
	return &Service{
		jwtSecret:   []byte(jwtSecret),
		expiration:  expiration,
		logger:      logger,
		metrics:     &observability.Metrics{},
		keyProvider: nil,
		blacklist:   nil,
	}
}

// NewServiceWithKeyRotation создает AuthService с ротацией ключей
func NewServiceWithKeyRotation(
	jwtSecret string,
	expiration time.Duration,
	logger *zap.Logger,
	rotationConfig *KeyRotationConfig,
	metrics *observability.Metrics,
) (*Service, error) {
	keyProvider, err := NewKeyProvider(rotationConfig, logger, metrics)
	if err != nil {
		return nil, fmt.Errorf("failed to create key provider: %w", err)
	}

	return &Service{
		jwtSecret:   []byte(jwtSecret),
		expiration:  expiration,
		logger:      logger,
		metrics:     metrics,
		keyProvider: keyProvider,
		blacklist:   nil,
	}, nil
}

// NewServiceWithBlacklist создает AuthService с поддержкой Token Blacklist
func NewServiceWithBlacklist(
	jwtSecret string,
	expiration time.Duration,
	logger *zap.Logger,
	blacklist *cache.TokenBlacklist,
	metrics *observability.Metrics,
) *Service {
	return &Service{
		jwtSecret:   []byte(jwtSecret),
		expiration:  expiration,
		logger:      logger,
		metrics:     metrics,
		keyProvider: nil,
		blacklist:   blacklist,
	}
}

// NewServiceWithKeyRotationAndBlacklist создает AuthService с обеими фичами
func NewServiceWithKeyRotationAndBlacklist(
	jwtSecret string,
	expiration time.Duration,
	logger *zap.Logger,
	rotationConfig *KeyRotationConfig,
	blacklist *cache.TokenBlacklist,
	metrics *observability.Metrics,
) (*Service, error) {
	keyProvider, err := NewKeyProvider(rotationConfig, logger, metrics)
	if err != nil {
		return nil, fmt.Errorf("failed to create key provider: %w", err)
	}

	return &Service{
		jwtSecret:   []byte(jwtSecret),
		expiration:  expiration,
		logger:      logger,
		metrics:     metrics,
		keyProvider: keyProvider,
		blacklist:   blacklist,
	}, nil
}

// Claims структура JWT claims
type Claims struct {
	UserID    string                 `json:"user_id"`
	Username  string                 `json:"username"`
	Email     string                 `json:"email"`
	Roles     []string               `json:"roles"`
	Role      string                 `json:"role"` // Single role string (user-server format)
	TokenType string                 `json:"typ"`  // Token type: access/refresh
	Extra     map[string]interface{} `json:"extra,omitempty"`
	jwt.RegisteredClaims
}

// ValidateToken проверяет JWT токен
func (s *Service) ValidateToken(ctx context.Context, tokenString string) (*ports.TokenClaims, error) {
	var keysToTry [][]byte

	// Если используется KeyProvider, собираем все доступные ключи
	if s.keyProvider != nil {
		validKeys := s.keyProvider.GetValidKeys()
		for _, keyVersion := range validKeys {
			keysToTry = append(keysToTry, []byte(keyVersion.Secret))
		}
	} else {
		keysToTry = append(keysToTry, s.jwtSecret)
	}

	// Если нет ключей для валидации, ошибка
	if len(keysToTry) == 0 {
		s.logger.Error("No keys available for token validation")
		return nil, ErrInvalidToken
	}

	var lastErr error
	var token *jwt.Token

	// Пытаемся валидировать с каждым ключом
	for i, signingKey := range keysToTry {
		parsedToken, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(t *jwt.Token) (interface{}, error) {
			// Проверяем метод подписи
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return signingKey, nil
		})

		if err == nil && parsedToken.Valid {
			token = parsedToken
			break
		}

		lastErr = err
		// Если это не последний ключ, продолжаем
		if i < len(keysToTry)-1 {
			continue
		}
	}

	if token == nil || !token.Valid {
		s.logger.Error("Failed to parse token", zap.Error(lastErr))
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		s.logger.Error("Invalid token claims")
		return nil, ErrInvalidClaims
	}

	// Проверяем срок действия
	if claims.ExpiresAt != nil && claims.ExpiresAt.Before(time.Now()) {
		return nil, ErrExpiredToken
	}

	// Проверяем, находится ли токен в blacklist
	if s.blacklist != nil {
		// Используем JTI (JWT ID) из claims или ID из ключа как уникальный идентификатор
		tokenID := claims.ID
		if tokenID == "" {
			tokenID = claims.UserID // Fallback если JTI не установлен
		}

		isBlacklisted, err := s.blacklist.IsBlacklisted(ctx, tokenID)
		if err != nil {
			s.logger.Warn("Failed to check token blacklist", zap.Error(err))
			// Продолжаем, не блокируем на ошибке Redis
		} else if isBlacklisted {
			s.logger.Warn("Token is blacklisted", zap.String("token_id", tokenID))
			return nil, ErrInvalidToken
		}
	}

	// Fallback: если user_id пуст, используем sub (Subject) из RegisteredClaims
	userID := claims.UserID
	if userID == "" && claims.Subject != "" {
		userID = claims.Subject
	}

	// Normalize roles: merge roles array and single role string (user-server compat)
	normalizedRoles := claims.Roles
	if len(normalizedRoles) == 0 && claims.Role != "" {
		normalizedRoles = []string{claims.Role}
	}

	return &ports.TokenClaims{
		UserID:   userID,
		Username: claims.Username,
		Email:    claims.Email,
		Roles:    normalizedRoles,
		Claims:   claims.Extra,
	}, nil
}

// GenerateToken генерирует JWT токен
func (s *Service) GenerateToken(ctx context.Context, userID string, claimsData map[string]interface{}) (string, error) {
	now := time.Now()
	expiresAt := now.Add(s.expiration)

	claims := &Claims{
		UserID: userID,
		Extra:  claimsData,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
		},
	}

	// Извлекаем дополнительные поля из claimsData
	if username, ok := claimsData["username"].(string); ok {
		claims.Username = username
	}
	if email, ok := claimsData["email"].(string); ok {
		claims.Email = email
	}
	if roles, ok := claimsData["roles"].([]string); ok {
		claims.Roles = roles
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Подписываем активным ключом из KeyProvider или обычным секретом
	var signingKey []byte
	if s.keyProvider != nil {
		activeKey, err := s.keyProvider.GetActiveKey()
		if err != nil {
			s.logger.Error("Failed to get active key", zap.Error(err))
			signingKey = s.jwtSecret // Fallback
		} else {
			signingKey = []byte(activeKey.Secret)
			claims.ID = activeKey.ID // Сохраняем ID ключа в токен
		}
	} else {
		signingKey = s.jwtSecret
	}

	tokenString, err := token.SignedString(signingKey)
	if err != nil {
		s.logger.Error("Failed to sign token", zap.Error(err))
		return "", fmt.Errorf("failed to generate token: %w", err)
	}

	s.logger.Debug("Token generated successfully",
		zap.String("user_id", userID),
		zap.Time("expires_at", expiresAt),
	)

	return tokenString, nil
}

// RefreshToken обновляет JWT токен
func (s *Service) RefreshToken(ctx context.Context, tokenString string) (string, error) {
	// Валидируем текущий токен
	claims, err := s.ValidateToken(ctx, tokenString)
	if err != nil && !errors.Is(err, ErrExpiredToken) {
		return "", err
	}

	// Генерируем новый токен с теми же claims
	claimsData := map[string]interface{}{
		"username": claims.Username,
		"email":    claims.Email,
		"roles":    claims.Roles,
	}

	// Добавляем дополнительные claims
	for k, v := range claims.Claims {
		claimsData[k] = v
	}

	return s.GenerateToken(ctx, claims.UserID, claimsData)
}

// ExtractTokenFromHeader извлекает токен из Authorization header
func ExtractTokenFromHeader(authHeader string) (string, error) {
	if authHeader == "" {
		return "", errors.New("authorization header is empty")
	}

	// Ожидаем формат: "Bearer <token>"
	const bearerPrefix = "Bearer "
	if len(authHeader) <= len(bearerPrefix) {
		return "", errors.New("invalid authorization header format")
	}

	if authHeader[:len(bearerPrefix)] != bearerPrefix {
		return "", errors.New("authorization header must start with 'Bearer '")
	}

	return authHeader[len(bearerPrefix):], nil
}

// Stop останавливает сервис и очищает ресурсы
func (s *Service) Stop() {
	if s.keyProvider != nil {
		s.keyProvider.Stop()
		s.logger.Info("Auth service stopped, key rotation disabled")
	}
}

// GetKeyProvider возвращает KeyProvider (для управления ключами)
func (s *Service) GetKeyProvider() *KeyProvider {
	return s.keyProvider
}

// GetBlacklist возвращает TokenBlacklist (для управления отозванными токенами)
func (s *Service) GetBlacklist() *cache.TokenBlacklist {
	return s.blacklist
}

// RevokeToken добавляет токен в blacklist (отзывает его)
func (s *Service) RevokeToken(ctx context.Context, tokenString string) error {
	if s.blacklist == nil {
		return fmt.Errorf("token blacklist is not enabled")
	}

	// Парсим токен без проверки подписи (нужен только ID)
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte("dummy"), nil // Подпись не проверяем
	})

	if err != nil && !errors.Is(err, jwt.ErrSignatureInvalid) {
		// Если ошибка не в подписи, это реальная ошибка
		return fmt.Errorf("failed to parse token: %w", err)
	}

	// Извлекаем Claims (даже если подпись невалидна, структура токена может быть OK)
	var tokenID string
	if token != nil && token.Claims != nil {
		if c, ok := token.Claims.(*Claims); ok {
			tokenID = c.ID
			if tokenID == "" {
				tokenID = c.UserID
			}
		}
	}

	if tokenID == "" {
		return fmt.Errorf("failed to extract token ID")
	}

	// Получаем время истечения токена
	var expiration time.Duration
	if token != nil && token.Claims != nil {
		if c, ok := token.Claims.(*Claims); ok {
			if c.ExpiresAt != nil && c.ExpiresAt.After(time.Now()) {
				expiration = time.Until(c.ExpiresAt.Time)
			} else {
				// Если токен уже истек, добавляем в blacklist на 1 день для логирования
				expiration = 24 * time.Hour
			}
		}
	}

	if expiration == 0 {
		expiration = 24 * time.Hour // Default TTL для неизвестных токенов
	}

	return s.blacklist.AddToken(ctx, tokenID, expiration)
}
