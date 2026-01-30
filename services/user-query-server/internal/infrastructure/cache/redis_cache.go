// Файл user-query-server/internal/infrastructure/cache/redis_cache.go содержит реализацию пакета cache.
package cache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"

	pb "github.com/rusgainew/kkm-project-mks/proto-lib/api"
)

const (
	userCacheKeyPrefix = "user:"
	userCacheTTL       = 5 * time.Minute
)

type RedisCache struct {
	client *redis.Client
	logger *zap.Logger
}

func NewRedisCache(client *redis.Client, logger *zap.Logger) *RedisCache {
	return &RedisCache{
		client: client,
		logger: logger,
	}
}

// GetUser retrieves a user from cache
func (c *RedisCache) GetUser(ctx context.Context, userID string) (*pb.UserReadModel, error) {
	key := userCacheKeyPrefix + userID

	val, err := c.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil // Cache miss
		}
		c.logger.Warn("Failed to get user from cache",
			zap.String("user_id", userID),
			zap.Error(err))
		return nil, nil // Treat errors as cache miss
	}

	var user pb.UserReadModel
	if err := json.Unmarshal([]byte(val), &user); err != nil {
		c.logger.Warn("Failed to unmarshal cached user",
			zap.String("user_id", userID),
			zap.Error(err))
		return nil, nil
	}

	return &user, nil
}

// SetUser stores a user in cache
func (c *RedisCache) SetUser(ctx context.Context, user *pb.UserReadModel) error {
	key := userCacheKeyPrefix + user.Id

	data, err := json.Marshal(user)
	if err != nil {
		c.logger.Warn("Failed to marshal user for cache",
			zap.String("user_id", user.Id),
			zap.Error(err))
		return nil // Fail silently, don't block request
	}

	err = c.client.Set(ctx, key, data, userCacheTTL).Err()
	if err != nil {
		c.logger.Warn("Failed to set user cache",
			zap.String("user_id", user.Id),
			zap.Error(err))
		return nil // Fail silently
	}

	return nil
}

// InvalidateUser removes a user from cache
func (c *RedisCache) InvalidateUser(ctx context.Context, userID string) error {
	key := userCacheKeyPrefix + userID
	return c.client.Del(ctx, key).Err()
}
