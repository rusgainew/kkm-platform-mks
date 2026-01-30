// Файл document-server/internal/infrastructure/cache/redis_cache.go содержит реализацию пакета cache.
package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

// RedisCache реализует интерфейс CacheLayer с использованием Redis
type RedisCache struct {
	client *redis.Client
	logger *zap.Logger
}

// NewRedisCache создаёт новый экземпляр Redis кэша
func NewRedisCache(addr string, logger *zap.Logger) (*RedisCache, error) {
	client := redis.NewClient(&redis.Options{
		Addr:         addr,
		MaxRetries:   3,
		PoolSize:     10,
		MinIdleConns: 5,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Проверяем соединение
	if err := client.Ping(ctx).Err(); err != nil {
		logger.Error("Failed to connect to Redis", zap.Error(err))
		return nil, err
	}

	logger.Info("Redis cache connected", zap.String("addr", addr))

	return &RedisCache{
		client: client,
		logger: logger,
	}, nil
}

// Get получает значение из кэша
func (rc *RedisCache) Get(ctx context.Context, key string) ([]byte, error) {
	val, err := rc.client.Get(ctx, key).Bytes()
	if err == redis.Nil {
		// Ключа нет в кэше
		rc.logger.Debug("Cache miss", zap.String("key", key))
		return nil, nil
	}
	if err != nil {
		rc.logger.Error("Failed to get from cache", zap.String("key", key), zap.Error(err))
		return nil, err
	}

	rc.logger.Debug("Cache hit", zap.String("key", key))
	return val, nil
}

// Set устанавливает значение в кэше
func (rc *RedisCache) Set(ctx context.Context, key string, value []byte, ttl time.Duration) error {
	if err := rc.client.Set(ctx, key, value, ttl).Err(); err != nil {
		rc.logger.Error("Failed to set cache",
			zap.String("key", key),
			zap.Duration("ttl", ttl),
			zap.Error(err),
		)
		return err
	}

	rc.logger.Debug("Cache set",
		zap.String("key", key),
		zap.Duration("ttl", ttl),
	)
	return nil
}

// Delete удаляет значение из кэша
func (rc *RedisCache) Delete(ctx context.Context, key string) error {
	if err := rc.client.Del(ctx, key).Err(); err != nil {
		rc.logger.Error("Failed to delete from cache", zap.String("key", key), zap.Error(err))
		return err
	}

	rc.logger.Debug("Cache delete", zap.String("key", key))
	return nil
}

// DeletePattern удаляет все ключи, совпадающие с паттерном
func (rc *RedisCache) DeletePattern(ctx context.Context, pattern string) error {
	iter := rc.client.Scan(ctx, 0, pattern, 0).Iterator()

	var keys []string
	for iter.Next(ctx) {
		keys = append(keys, iter.Val())
	}

	if err := iter.Err(); err != nil {
		rc.logger.Error("Failed to scan keys", zap.String("pattern", pattern), zap.Error(err))
		return err
	}

	if len(keys) == 0 {
		rc.logger.Debug("No keys found for pattern", zap.String("pattern", pattern))
		return nil
	}

	if err := rc.client.Del(ctx, keys...).Err(); err != nil {
		rc.logger.Error("Failed to delete keys by pattern",
			zap.String("pattern", pattern),
			zap.Int("count", len(keys)),
			zap.Error(err),
		)
		return err
	}

	rc.logger.Debug("Cache delete pattern",
		zap.String("pattern", pattern),
		zap.Int("count", len(keys)),
	)
	return nil
}

// Exists проверяет наличие ключа в кэше
func (rc *RedisCache) Exists(ctx context.Context, key string) (bool, error) {
	exists, err := rc.client.Exists(ctx, key).Result()
	if err != nil {
		rc.logger.Error("Failed to check key existence", zap.String("key", key), zap.Error(err))
		return false, err
	}

	return exists > 0, nil
}

// TTL получает время жизни ключа
func (rc *RedisCache) TTL(ctx context.Context, key string) (int64, error) {
	ttl, err := rc.client.TTL(ctx, key).Result()
	if err != nil {
		rc.logger.Error("Failed to get TTL", zap.String("key", key), zap.Error(err))
		return -2, err
	}

	return int64(ttl.Seconds()), nil
}

// Close закрывает соединение с Redis
func (rc *RedisCache) Close() error {
	if err := rc.client.Close(); err != nil {
		rc.logger.Error("Failed to close Redis connection", zap.Error(err))
		return err
	}

	rc.logger.Info("Redis connection closed")
	return nil
}

// SerializeDocument сериализует документ в JSON для кэша
func SerializeDocument(doc map[string]interface{}) ([]byte, error) {
	return json.Marshal(doc)
}

// DeserializeDocument десериализует JSON из кэша в документ
func DeserializeDocument(data []byte) (map[string]interface{}, error) {
	var doc map[string]interface{}
	if err := json.Unmarshal(data, &doc); err != nil {
		return nil, fmt.Errorf("failed to deserialize document: %w", err)
	}

	return doc, nil
}

// GetDocumentCacheKey возвращает ключ кэша для документа
func GetDocumentCacheKey(documentID string) string {
	return fmt.Sprintf("document:%s", documentID)
}

// GetDocumentListCacheKey возвращает ключ кэша для списка документов
func GetDocumentListCacheKey(organizationID, status string, page, perPage int) string {
	return fmt.Sprintf("document:list:%s:%s:%d:%d", organizationID, status, page, perPage)
}

// InvalidateDocumentCache инвалидирует кэш для документа
func InvalidateDocumentCache(ctx context.Context, cache interface {
	DeletePattern(context.Context, string) error
}, documentID string) error {
	// Удаляем кэш самого документа
	if err := cache.DeletePattern(ctx, fmt.Sprintf("document:%s", documentID)); err != nil {
		return err
	}

	// Удаляем кэш всех списков (так как документ мог быть в списках)
	if err := cache.DeletePattern(ctx, "document:list:*"); err != nil {
		return err
	}

	return nil
}

// InvalidateListCache инвалидирует кэш списков для организации
func InvalidateListCache(ctx context.Context, cache interface {
	DeletePattern(context.Context, string) error
}, organizationID string) error {
	pattern := fmt.Sprintf("document:list:%s:*", organizationID)
	return cache.DeletePattern(ctx, pattern)
}
