// Файл api-gateway/internal/infrastructure/cache/redis_cache.go содержит реализацию пакета cache.
package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/rusgainew/kkm-project-mks/api-gateway/internal/infrastructure/observability"
)

// RedisCache предоставляет кеширование через Redis
type RedisCache struct {
	client  *redis.Client
	ttl     time.Duration
	metrics *observability.Metrics
	name    string
}

// NewRedisCache создает новый Redis кеш
func NewRedisCache(ctx context.Context, addr string, password string, ttl time.Duration, metrics *observability.Metrics) (*RedisCache, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       0,
	})

	// Проверяем подключение
	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := client.Ping(pingCtx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to redis: %w", err)
	}

	return &RedisCache{
		client:  client,
		ttl:     ttl,
		metrics: metrics,
		name:    "default",
	}, nil
}

// Set сохраняет значение в кеш
func (rc *RedisCache) Set(ctx context.Context, key string, value interface{}) error {
	start := time.Now()
	defer func() {
		latency := time.Since(start).Seconds()
		rc.metrics.RecordCacheLatency("set", rc.name, latency)
	}()

	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return rc.client.Set(ctx, key, data, rc.ttl).Err()
}

// Get получает значение из кеша
func (rc *RedisCache) Get(ctx context.Context, key string, dest interface{}) error {
	start := time.Now()

	val, err := rc.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			rc.metrics.RecordCacheMiss(rc.name)
		}
		return err
	}

	latency := time.Since(start).Seconds()
	rc.metrics.RecordCacheLatency("get", rc.name, latency)
	rc.metrics.RecordCacheHit(rc.name)

	return json.Unmarshal([]byte(val), dest)
}

// Delete удаляет значение из кеша
func (rc *RedisCache) Delete(ctx context.Context, keys ...string) error {
	if len(keys) == 0 {
		return nil
	}

	start := time.Now()
	defer func() {
		latency := time.Since(start).Seconds()
		rc.metrics.RecordCacheLatency("delete", rc.name, latency)
	}()

	return rc.client.Del(ctx, keys...).Err()
}

// DeletePattern удаляет все ключи по паттерну
func (rc *RedisCache) DeletePattern(ctx context.Context, pattern string) error {
	iter := rc.client.Scan(ctx, 0, pattern, 0).Iterator()
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

	start := time.Now()
	defer func() {
		latency := time.Since(start).Seconds()
		rc.metrics.RecordCacheLatency("delete_pattern", rc.name, latency)
		rc.metrics.RecordCacheEviction(rc.name)
	}()

	return rc.client.Del(ctx, keys...).Err()
}

// Close закрывает подключение к Redis
func (rc *RedisCache) Close() error {
	return rc.client.Close()
}

// IsAvailable проверяет доступность Redis
func (rc *RedisCache) IsAvailable(ctx context.Context) bool {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	return rc.client.Ping(ctx).Err() == nil
}
