package auth

import (
	"context"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/rusgainew/kkm-project-mks/api-gateway/internal/infrastructure/cache"
	"github.com/rusgainew/kkm-project-mks/api-gateway/internal/infrastructure/observability"
)

func init() {
	observability.ResetForTesting()
}

func TestIntegration_ValidateToken_WithBlacklist(t *testing.T) {
	// Setup Redis
	mr := miniredis.RunT(t)
	defer mr.Close()

	client := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})

	metrics := observability.NewMetrics()
	blacklist := cache.NewTokenBlacklist(client, metrics)
	logger := zap.NewNop()

	// Create service with blacklist
	service := NewServiceWithBlacklist("test-secret", 1*time.Hour, logger, blacklist, metrics)

	ctx := context.Background()
	token, err := service.GenerateToken(ctx, "user-123", map[string]interface{}{
		"username": "testuser",
		"email":    "test@example.com",
		"roles":    []string{"user"},
	})
	require.NoError(t, err)

	// Validate token initially
	claims, err := service.ValidateToken(ctx, token)
	require.NoError(t, err)
	assert.Equal(t, "user-123", claims.UserID)
}

func TestIntegration_MultipleTokens_InBlacklist(t *testing.T) {
	// Setup Redis
	mr := miniredis.RunT(t)
	defer mr.Close()

	client := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})

	metrics := observability.NewMetrics()
	blacklist := cache.NewTokenBlacklist(client, metrics)

	ctx := context.Background()

	// Add multiple tokens
	tokenIDs := []string{"token-1", "token-2", "token-3"}
	for _, id := range tokenIDs {
		err := blacklist.AddToken(ctx, id, 1*time.Hour)
		require.NoError(t, err)
	}

	// Verify all are blacklisted
	for _, id := range tokenIDs {
		isBlacklisted, err := blacklist.IsBlacklisted(ctx, id)
		require.NoError(t, err)
		assert.True(t, isBlacklisted)
	}

	// Verify non-existent token is not blacklisted
	isBlacklisted, err := blacklist.IsBlacklisted(ctx, "token-99")
	require.NoError(t, err)
	assert.False(t, isBlacklisted)
}
