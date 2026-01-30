// Файл api-gateway/internal/infrastructure/cache/token_blacklist_test.go содержит реализацию пакета cache.
package cache

import (
	"context"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/rusgainew/kkm-project-mks/api-gateway/internal/infrastructure/observability"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func init() {
	observability.ResetForTesting()
}

// TestTokenBlacklist_AddToken проверяет добавление токена в blacklist
func TestTokenBlacklist_AddToken(t *testing.T) {
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	defer client.Close()

	if err := client.Ping(context.Background()).Err(); err != nil {
		t.Skip("Redis not available for testing")
	}

	metrics := observability.NewMetrics()
	tb := NewTokenBlacklist(client, metrics)
	ctx := context.Background()

	tb.ClearBlacklist(ctx)

	tokenJTI := "test-token-123"
	ttl := 24 * time.Hour

	err := tb.AddToken(ctx, tokenJTI, ttl)
	assert.NoError(t, err)

	isBlacklisted, err := tb.IsBlacklisted(ctx, tokenJTI)
	assert.NoError(t, err)
	assert.True(t, isBlacklisted)

	tb.ClearBlacklist(ctx)
}

// TestTokenBlacklist_IsBlacklisted проверяет проверку наличия токена в blacklist
func TestTokenBlacklist_IsBlacklisted(t *testing.T) {
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	defer client.Close()

	if err := client.Ping(context.Background()).Err(); err != nil {
		t.Skip("Redis not available for testing")
	}

	metrics := observability.NewMetrics()
	tb := NewTokenBlacklist(client, metrics)
	ctx := context.Background()

	tb.ClearBlacklist(ctx)

	tokenJTI := "test-token-456"

	isBlacklisted, err := tb.IsBlacklisted(ctx, tokenJTI)
	assert.NoError(t, err)
	assert.False(t, isBlacklisted)

	err = tb.AddToken(ctx, tokenJTI, 24*time.Hour)
	assert.NoError(t, err)

	isBlacklisted, err = tb.IsBlacklisted(ctx, tokenJTI)
	assert.NoError(t, err)
	assert.True(t, isBlacklisted)

	tb.ClearBlacklist(ctx)
}

// TestTokenBlacklist_RemoveToken проверяет удаление токена из blacklist
func TestTokenBlacklist_RemoveToken(t *testing.T) {
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	defer client.Close()

	if err := client.Ping(context.Background()).Err(); err != nil {
		t.Skip("Redis not available for testing")
	}

	metrics := observability.NewMetrics()
	tb := NewTokenBlacklist(client, metrics)
	ctx := context.Background()

	tb.ClearBlacklist(ctx)

	tokenJTI := "test-token-789"

	err := tb.AddToken(ctx, tokenJTI, 24*time.Hour)
	require.NoError(t, err)

	isBlacklisted, err := tb.IsBlacklisted(ctx, tokenJTI)
	require.NoError(t, err)
	require.True(t, isBlacklisted)

	err = tb.RemoveToken(ctx, tokenJTI)
	assert.NoError(t, err)

	isBlacklisted, err = tb.IsBlacklisted(ctx, tokenJTI)
	assert.NoError(t, err)
	assert.False(t, isBlacklisted)

	tb.ClearBlacklist(ctx)
}

// TestTokenBlacklist_InvalidJTI проверяет ошибку с пустым JTI
func TestTokenBlacklist_InvalidJTI(t *testing.T) {
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	defer client.Close()

	if err := client.Ping(context.Background()).Err(); err != nil {
		t.Skip("Redis not available for testing")
	}

	metrics := observability.NewMetrics()
	tb := NewTokenBlacklist(client, metrics)
	ctx := context.Background()

	err := tb.AddToken(ctx, "", 24*time.Hour)
	assert.Error(t, err)

	_, err = tb.IsBlacklisted(ctx, "")
	assert.Error(t, err)

	err = tb.RemoveToken(ctx, "")
	assert.Error(t, err)
}
