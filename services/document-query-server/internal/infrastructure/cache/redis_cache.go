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
	documentCacheKeyPrefix = "document:"
	documentCacheTTL       = 10 * time.Minute
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

// GetDocument retrieves a document from cache
func (c *RedisCache) GetDocument(ctx context.Context, documentID string) (*pb.DocumentReadModel, error) {
	key := documentCacheKeyPrefix + documentID

	val, err := c.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil // Cache miss
		}
		c.logger.Warn("Failed to get document from cache",
			zap.String("document_id", documentID),
			zap.Error(err))
		return nil, nil // Treat errors as cache miss
	}

	var doc pb.DocumentReadModel
	if err := json.Unmarshal([]byte(val), &doc); err != nil {
		c.logger.Warn("Failed to unmarshal cached document",
			zap.String("document_id", documentID),
			zap.Error(err))
		return nil, nil
	}

	return &doc, nil
}

// SetDocument stores a document in cache
func (c *RedisCache) SetDocument(ctx context.Context, doc *pb.DocumentReadModel) error {
	key := documentCacheKeyPrefix + doc.Id

	data, err := json.Marshal(doc)
	if err != nil {
		c.logger.Warn("Failed to marshal document for cache",
			zap.String("document_id", doc.Id),
			zap.Error(err))
		return nil // Fail silently, don't block request
	}

	err = c.client.Set(ctx, key, data, documentCacheTTL).Err()
	if err != nil {
		c.logger.Warn("Failed to set document cache",
			zap.String("document_id", doc.Id),
			zap.Error(err))
		return nil // Fail silently
	}

	return nil
}

// InvalidateDocument removes a document from cache
func (c *RedisCache) InvalidateDocument(ctx context.Context, documentID string) error {
	key := documentCacheKeyPrefix + documentID
	return c.client.Del(ctx, key).Err()
}
