// Файл invoice-query-server/internal/infrastructure/cache/redis_cache.go содержит реализацию пакета cache.
package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

// RedisCache provides caching functionality for query results
type RedisCache struct {
	client *redis.Client
	logger *zap.Logger
	ttl    time.Duration
}

// NewRedisCache creates a new Redis cache instance
func NewRedisCache(redisURL string, ttl time.Duration, logger *zap.Logger) (*RedisCache, error) {
	opts, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse Redis URL: %w", err)
	}

	client := redis.NewClient(opts)

	// Ping to verify connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	logger.Info("Redis cache initialized", zap.String("url", redisURL), zap.Duration("ttl", ttl))

	return &RedisCache{
		client: client,
		logger: logger,
		ttl:    ttl,
	}, nil
}

// Get retrieves a value from cache
func (c *RedisCache) Get(ctx context.Context, key string, dest interface{}) error {
	val, err := c.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return fmt.Errorf("cache miss")
	}
	if err != nil {
		c.logger.Error("Redis GET error", zap.String("key", key), zap.Error(err))
		return err
	}

	if err := json.Unmarshal([]byte(val), dest); err != nil {
		// Invalid cache entry - delete it and treat as cache miss
		c.logger.Warn("Cache entry corrupted, removing", zap.String("key", key))
		c.client.Del(ctx, key)
		return fmt.Errorf("cache miss: corrupted entry")
	}

	c.logger.Debug("Cache hit", zap.String("key", key))
	return nil
}

// Set stores a value in cache with TTL
func (c *RedisCache) Set(ctx context.Context, key string, value interface{}) error {
	data, err := json.Marshal(value)
	if err != nil {
		c.logger.Error("Failed to marshal value", zap.String("key", key), zap.Error(err))
		return err
	}

	if err := c.client.Set(ctx, key, data, c.ttl).Err(); err != nil {
		c.logger.Error("Redis SET error", zap.String("key", key), zap.Error(err))
		return err
	}

	c.logger.Debug("Cache set", zap.String("key", key), zap.Duration("ttl", c.ttl))
	return nil
}

// Delete removes a key from cache
func (c *RedisCache) Delete(ctx context.Context, key string) error {
	if err := c.client.Del(ctx, key).Err(); err != nil {
		c.logger.Error("Redis DEL error", zap.String("key", key), zap.Error(err))
		return err
	}
	return nil
}

// Close closes the Redis connection
func (c *RedisCache) Close() error {
	return c.client.Close()
}

// GenerateCacheKey creates a cache key from method name and parameters
func GenerateCacheKey(service, method string, params ...interface{}) string {
	key := fmt.Sprintf("%s:%s", service, method)
	for _, param := range params {
		key += fmt.Sprintf(":%v", param)
	}
	return key
}
