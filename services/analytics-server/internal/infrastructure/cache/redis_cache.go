// Файл analytics-server/internal/infrastructure/cache/redis_cache.go содержит реализацию пакета cache.
package cache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

// RedisCache обертка для работы с Redis
type RedisCache struct {
	client *redis.Client
	logger *zap.Logger
	ttl    time.Duration
}

// NewRedisCache создает новый экземпляр Redis cache
func NewRedisCache(addr, password string, db int, ttl int, logger *zap.Logger) *RedisCache {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	return &RedisCache{
		client: client,
		logger: logger,
		ttl:    time.Duration(ttl) * time.Second,
	}
}

// Get получает значение по ключу и десериализует в value
func (c *RedisCache) Get(ctx context.Context, key string, value interface{}) error {
	data, err := c.client.Get(ctx, key).Bytes()
	if err == redis.Nil {
		return redis.Nil
	}
	if err != nil {
		c.logger.Error("Failed to get from cache", zap.String("key", key), zap.Error(err))
		return err
	}

	if err := json.Unmarshal(data, value); err != nil {
		c.logger.Error("Failed to unmarshal cached value", zap.Error(err))
		return err
	}

	return nil
}

// Set устанавливает значение с TTL и сериализацией в JSON
func (c *RedisCache) Set(ctx context.Context, key string, value interface{}) error {
	data, err := json.Marshal(value)
	if err != nil {
		c.logger.Error("Failed to marshal value", zap.Error(err))
		return err
	}

	err = c.client.Set(ctx, key, data, c.ttl).Err()
	if err != nil {
		c.logger.Error("Failed to set cache", zap.String("key", key), zap.Error(err))
		return err
	}
	return nil
}

// Delete удаляет значение по ключу
func (c *RedisCache) Delete(ctx context.Context, key string) error {
	err := c.client.Del(ctx, key).Err()
	if err != nil {
		c.logger.Error("Failed to delete from cache", zap.String("key", key), zap.Error(err))
		return err
	}
	return nil
}

// DeleteByPattern удаляет все ключи, соответствующие паттерну
func (c *RedisCache) DeleteByPattern(ctx context.Context, pattern string) error {
	iter := c.client.Scan(ctx, 0, pattern, 0).Iterator()
	for iter.Next(ctx) {
		err := c.client.Del(ctx, iter.Val()).Err()
		if err != nil {
			c.logger.Error("Failed to delete key", zap.String("key", iter.Val()), zap.Error(err))
		}
	}
	if err := iter.Err(); err != nil {
		c.logger.Error("Failed to scan keys", zap.String("pattern", pattern), zap.Error(err))
		return err
	}
	return nil
}

// Ping проверяет подключение к Redis
func (c *RedisCache) Ping(ctx context.Context) error {
	return c.client.Ping(ctx).Err()
}

// Close закрывает соединение с Redis
func (c *RedisCache) Close() error {
	return c.client.Close()
}
