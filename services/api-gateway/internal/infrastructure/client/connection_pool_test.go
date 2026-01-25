package client

import (
	"context"
	"testing"
	"time"

	"github.com/rusgainew/kkm-project-mks/api-gateway/internal/infrastructure/observability"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func init() {
	observability.ResetForTesting()
}

func TestConnectionPool_MaxConnections_Enforced(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	metrics := observability.NewMetrics()
	config := &ConnectionPoolConfig{
		MaxConnections:   2,
		MaxIdleTime:      1 * time.Minute,
		EvictionInterval: 10 * time.Second,
		ConnectTimeout:   1 * time.Second,
	}

	pool := NewConnectionPool(config, logger, metrics)

	callCount := 0
	dialFunc := func(ctx context.Context, addr string) (*grpc.ClientConn, error) {
		callCount++
		return &grpc.ClientConn{}, nil
	}

	ctx := context.Background()

	// First connection OK
	conn1, err := pool.GetConnection(ctx, "host1:9001", dialFunc)
	assert.NoError(t, err)
	assert.NotNil(t, conn1)
	assert.Equal(t, 1, pool.GetActiveConnectionCount())

	// Second connection OK
	conn2, err := pool.GetConnection(ctx, "host2:9002", dialFunc)
	assert.NoError(t, err)
	assert.NotNil(t, conn2)
	assert.Equal(t, 2, pool.GetActiveConnectionCount())

	// Third connection should fail (limit reached)
	_, err = pool.GetConnection(ctx, "host3:9003", dialFunc)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "connection pool limit reached")
	assert.Equal(t, 2, callCount) // Only 2 dials were made
}

func TestConnectionPool_Reuse_Connection(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	metrics := observability.NewMetrics()
	config := DefaultConnectionPoolConfig()
	pool := NewConnectionPool(config, logger, metrics)

	callCount := 0
	dialFunc := func(ctx context.Context, addr string) (*grpc.ClientConn, error) {
		callCount++
		return &grpc.ClientConn{}, nil
	}

	ctx := context.Background()

	conn1, err := pool.GetConnection(ctx, "host:9001", dialFunc)
	require.NoError(t, err)
	assert.Equal(t, 1, callCount)

	// Second call should reuse
	conn2, err := pool.GetConnection(ctx, "host:9001", dialFunc)
	require.NoError(t, err)
	assert.Equal(t, 1, callCount) // Dial not called again
	assert.Equal(t, conn1, conn2)
}

func TestConnectionPool_GetStats(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	metrics := observability.NewMetrics()
	config := DefaultConnectionPoolConfig()
	pool := NewConnectionPool(config, logger, metrics)

	dialFunc := func(ctx context.Context, addr string) (*grpc.ClientConn, error) {
		return &grpc.ClientConn{}, nil
	}

	ctx := context.Background()

	_, err := pool.GetConnection(ctx, "host:9001", dialFunc)
	require.NoError(t, err)

	stats := pool.GetStats("host:9001")
	assert.NotNil(t, stats)
	assert.Equal(t, "host:9001", stats.Address)
	assert.Equal(t, int64(1), stats.UsageCount)
	assert.False(t, stats.IsClosed)
}

func TestConnectionPool_PoolStats(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	metrics := observability.NewMetrics()
	config := DefaultConnectionPoolConfig()
	pool := NewConnectionPool(config, logger, metrics)

	dialFunc := func(ctx context.Context, addr string) (*grpc.ClientConn, error) {
		return &grpc.ClientConn{}, nil
	}

	ctx := context.Background()

	_, err := pool.GetConnection(ctx, "host1:9001", dialFunc)
	require.NoError(t, err)

	_, err = pool.GetConnection(ctx, "host2:9002", dialFunc)
	require.NoError(t, err)

	stats := pool.GetPoolStats()
	assert.Len(t, stats, 2)
	assert.Contains(t, stats, "host1:9001")
	assert.Contains(t, stats, "host2:9002")
}

func TestConnectionPool_ActiveCount(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	metrics := observability.NewMetrics()
	config := DefaultConnectionPoolConfig()
	pool := NewConnectionPool(config, logger, metrics)

	dialFunc := func(ctx context.Context, addr string) (*grpc.ClientConn, error) {
		return &grpc.ClientConn{}, nil
	}

	ctx := context.Background()

	assert.Equal(t, 0, pool.GetActiveConnectionCount())

	_, err := pool.GetConnection(ctx, "host1:9001", dialFunc)
	require.NoError(t, err)
	assert.Equal(t, 1, pool.GetActiveConnectionCount())

	_, err = pool.GetConnection(ctx, "host2:9002", dialFunc)
	require.NoError(t, err)
	assert.Equal(t, 2, pool.GetActiveConnectionCount())
}

func TestConnectionPool_GetMaxConnections(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	metrics := observability.NewMetrics()
	config := &ConnectionPoolConfig{
		MaxConnections:   50,
		MaxIdleTime:      5 * time.Minute,
		EvictionInterval: 1 * time.Minute,
	}
	pool := NewConnectionPool(config, logger, metrics)
	assert.Equal(t, 50, pool.GetMaxConnections())
}

func TestDefaultConnectionPoolConfig(t *testing.T) {
	config := DefaultConnectionPoolConfig()
	assert.Equal(t, 100, config.MaxConnections)
	assert.Equal(t, 5*time.Minute, config.MaxIdleTime)
	assert.Equal(t, 1*time.Minute, config.EvictionInterval)
	assert.Equal(t, 10*time.Second, config.ConnectTimeout)
}
