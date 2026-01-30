// Файл api-gateway/internal/infrastructure/cache/token_blacklist.go содержит реализацию пакета cache.
package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/rusgainew/kkm-project-mks/api-gateway/internal/infrastructure/observability"
)

// TokenBlacklist управляет списком отозванных токенов через Redis
type TokenBlacklist struct {
	client  *redis.Client
	metrics *observability.Metrics
}

// NewTokenBlacklist создает новый TokenBlacklist
func NewTokenBlacklist(client *redis.Client, metrics *observability.Metrics) *TokenBlacklist {
	return &TokenBlacklist{
		client:  client,
		metrics: metrics,
	}
}

// AddToken добавляет токен в blacklist с TTL (Time To Live)
// TTL должен соответствовать времени истечения токена
func (tb *TokenBlacklist) AddToken(ctx context.Context, tokenJTI string, ttl time.Duration) error {
	if tokenJTI == "" {
		return fmt.Errorf("token JTI cannot be empty")
	}

	// Используем формат ключа: blacklist:<jti>
	key := fmt.Sprintf("blacklist:%s", tokenJTI)

	// Сохраняем с TTL - ключ автоматически удалится по истечении срока
	if err := tb.client.Set(ctx, key, "revoked", ttl).Err(); err != nil {
		tb.metrics.RecordRevokedToken("failed")
		return err
	}

	tb.metrics.RecordRevokedToken("success")
	return nil
}

// IsBlacklisted проверяет, находится ли токен в blacklist
func (tb *TokenBlacklist) IsBlacklisted(ctx context.Context, tokenJTI string) (bool, error) {
	start := time.Now()

	if tokenJTI == "" {
		latency := time.Since(start).Seconds()
		tb.metrics.RecordBlacklistCheckLatency("error", latency)
		return false, fmt.Errorf("token JTI cannot be empty")
	}

	key := fmt.Sprintf("blacklist:%s", tokenJTI)

	result, err := tb.client.Get(ctx, key).Result()
	if err == redis.Nil {
		// Ключ не существует = токен не в blacklist
		latency := time.Since(start).Seconds()
		tb.metrics.RecordBlacklistCheckLatency("not_found", latency)
		return false, nil
	}
	if err != nil {
		latency := time.Since(start).Seconds()
		tb.metrics.RecordBlacklistCheckLatency("error", latency)
		return false, fmt.Errorf("failed to check blacklist: %w", err)
	}

	latency := time.Since(start).Seconds()
	tb.metrics.RecordBlacklistCheckLatency("found", latency)

	return result == "revoked", nil
}

// RemoveToken удаляет токен из blacklist (например, при восстановлении)
func (tb *TokenBlacklist) RemoveToken(ctx context.Context, tokenJTI string) error {
	if tokenJTI == "" {
		return fmt.Errorf("token JTI cannot be empty")
	}

	key := fmt.Sprintf("blacklist:%s", tokenJTI)
	return tb.client.Del(ctx, key).Err()
}

// GetBlacklistedTokensCount возвращает количество токенов в blacklist
func (tb *TokenBlacklist) GetBlacklistedTokensCount(ctx context.Context) (int64, error) {
	// Сканируем все ключи с префиксом blacklist:
	iter := tb.client.Scan(ctx, 0, "blacklist:*", 0).Iterator()

	var count int64
	for iter.Next(ctx) {
		count++
	}

	if err := iter.Err(); err != nil {
		return 0, fmt.Errorf("failed to scan blacklist: %w", err)
	}

	tb.metrics.SetBlacklistSize(count)

	return count, nil
}

// ClearBlacklist удаляет все токены из blacklist
func (tb *TokenBlacklist) ClearBlacklist(ctx context.Context) error {
	iter := tb.client.Scan(ctx, 0, "blacklist:*", 0).Iterator()
	var keys []string

	for iter.Next(ctx) {
		keys = append(keys, iter.Val())
	}

	if err := iter.Err(); err != nil {
		return err
	}

	if len(keys) == 0 {
		return nil
	}

	return tb.client.Del(ctx, keys...).Err()
}
