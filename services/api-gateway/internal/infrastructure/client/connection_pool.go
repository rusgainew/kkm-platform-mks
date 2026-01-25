package client

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/rusgainew/kkm-project-mks/api-gateway/internal/infrastructure/observability"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

// ConnectionPoolConfig конфигурация для ConnectionPool
type ConnectionPoolConfig struct {
	MaxConnections   int           // Максимум соединений (default 100)
	MaxIdleTime      time.Duration // Время жизни idle соединения (default 5 min)
	EvictionInterval time.Duration // Интервал проверки idle соединений (default 1 min)
	ConnectTimeout   time.Duration // Timeout для подключения
}

// DefaultConnectionPoolConfig возвращает дефолтную конфиг
func DefaultConnectionPoolConfig() *ConnectionPoolConfig {
	return &ConnectionPoolConfig{
		MaxConnections:   100,
		MaxIdleTime:      5 * time.Minute,
		EvictionInterval: 1 * time.Minute,
		ConnectTimeout:   10 * time.Second,
	}
}

// ConnectionStats статистика соединения
type ConnectionStats struct {
	Address    string
	CreatedAt  time.Time
	LastUsedAt time.Time
	UsageCount int64
	IsClosed   bool
}

// ConnectionPool управляет пулом gRPC соединений с ограничениями
type ConnectionPool struct {
	mu             sync.RWMutex
	connections    map[string]*pooledConnection
	config         *ConnectionPoolConfig
	logger         *zap.Logger
	metrics        *observability.Metrics
	activeCount    int
	totalCreated   int64
	evictionTicker *time.Ticker
	stopEviction   chan struct{}
}

type pooledConnection struct {
	conn       *grpc.ClientConn
	createdAt  time.Time
	lastUsedAt time.Time
	usageCount int64
}

// NewConnectionPool создает новый ConnectionPool
func NewConnectionPool(config *ConnectionPoolConfig, logger *zap.Logger, metrics *observability.Metrics) *ConnectionPool {
	if config == nil {
		config = DefaultConnectionPoolConfig()
	}

	pool := &ConnectionPool{
		connections:  make(map[string]*pooledConnection),
		config:       config,
		logger:       logger,
		metrics:      metrics,
		stopEviction: make(chan struct{}),
	}

	// Запускаем goroutine для очистки idle соединений
	pool.evictionTicker = time.NewTicker(config.EvictionInterval)
	go pool.evictIdleConnections()

	// Устанавливаем начальные метрики
	metrics.SetMaxConnections(int64(config.MaxConnections))

	return pool
}

// GetConnection возвращает соединение из пула или создает новое
func (cp *ConnectionPool) GetConnection(ctx context.Context, address string, dialFunc func(context.Context, string) (*grpc.ClientConn, error)) (*grpc.ClientConn, error) {
	start := time.Now()
	cp.mu.Lock()

	// Проверяем существующее соединение
	if pc, exists := cp.connections[address]; exists {
		pc.lastUsedAt = time.Now()
		pc.usageCount++
		cp.mu.Unlock()

		latency := time.Since(start).Seconds()
		cp.metrics.RecordConnectionPoolWaitLatency("hit", latency)
		cp.metrics.SetActiveConnectionsCount(int64(cp.activeCount))

		return pc.conn, nil
	}

	// Проверяем лимит соединений
	if cp.activeCount >= cp.config.MaxConnections {
		cp.mu.Unlock()
		latency := time.Since(start).Seconds()
		cp.metrics.RecordConnectionPoolWaitLatency("limit_reached", latency)
		return nil, fmt.Errorf("connection pool limit reached: %d/%d", cp.activeCount, cp.config.MaxConnections)
	}

	cp.mu.Unlock()

	// Создаем новое соединение
	conn, err := dialFunc(ctx, address)
	if err != nil {
		latency := time.Since(start).Seconds()
		cp.metrics.RecordConnectionPoolWaitLatency("error", latency)
		return nil, fmt.Errorf("failed to dial %s: %w", address, err)
	}

	cp.mu.Lock()
	defer cp.mu.Unlock()

	// Повторная проверка лимита (на случай race condition)
	if cp.activeCount >= cp.config.MaxConnections {
		conn.Close()
		latency := time.Since(start).Seconds()
		cp.metrics.RecordConnectionPoolWaitLatency("limit_reached", latency)
		return nil, fmt.Errorf("connection pool limit reached: %d/%d", cp.activeCount, cp.config.MaxConnections)
	}

	now := time.Now()
	cp.connections[address] = &pooledConnection{
		conn:       conn,
		createdAt:  now,
		lastUsedAt: now,
		usageCount: 1,
	}

	cp.activeCount++
	cp.totalCreated++

	latency := time.Since(start).Seconds()
	cp.metrics.RecordConnectionPoolWaitLatency("new", latency)
	cp.metrics.SetActiveConnectionsCount(int64(cp.activeCount))

	cp.logger.Info("Connection added to pool",
		zap.String("address", address),
		zap.Int("active_connections", cp.activeCount),
	)

	return conn, nil
}

// evictIdleConnections удаляет неиспользуемые соединения
func (cp *ConnectionPool) evictIdleConnections() {
	for {
		select {
		case <-cp.evictionTicker.C:
			cp.evict()
		case <-cp.stopEviction:
			cp.evictionTicker.Stop()
			return
		}
	}
}

func (cp *ConnectionPool) evict() {
	cp.mu.Lock()
	defer cp.mu.Unlock()

	now := time.Now()
	var toDelete []string

	for address, pc := range cp.connections {
		idleTime := now.Sub(pc.lastUsedAt)

		// Проверяем, прошло ли достаточно времени неиспользования
		if idleTime > cp.config.MaxIdleTime {
			toDelete = append(toDelete, address)
		}
	}

	// Удаляем idle соединения
	for _, address := range toDelete {
		pc := cp.connections[address]
		pc.conn.Close()
		delete(cp.connections, address)
		cp.activeCount--

		cp.metrics.RecordConnectionEviction("idle_timeout")
		cp.metrics.SetActiveConnectionsCount(int64(cp.activeCount))
		cp.metrics.SetIdleConnectionsCount(int64(cp.getIdleConnectionCount()))

		cp.logger.Debug("Evicted idle connection",
			zap.String("address", address),
			zap.Duration("idle_time", now.Sub(pc.lastUsedAt)),
		)
	}
}

func (cp *ConnectionPool) getIdleConnectionCount() int {
	now := time.Now()
	count := 0
	for _, pc := range cp.connections {
		if now.Sub(pc.lastUsedAt) > cp.config.MaxIdleTime {
			count++
		}
	}
	return count
}

// GetStats возвращает статистику соединения
func (cp *ConnectionPool) GetStats(address string) *ConnectionStats {
	cp.mu.RLock()
	defer cp.mu.RUnlock()

	if pc, exists := cp.connections[address]; exists {
		return &ConnectionStats{
			Address:    address,
			CreatedAt:  pc.createdAt,
			LastUsedAt: pc.lastUsedAt,
			UsageCount: pc.usageCount,
			IsClosed:   false,
		}
	}

	return nil
}

// GetPoolStats возвращает статистику всего пула
func (cp *ConnectionPool) GetPoolStats() map[string]*ConnectionStats {
	cp.mu.RLock()
	defer cp.mu.RUnlock()

	stats := make(map[string]*ConnectionStats)
	for address, pc := range cp.connections {
		stats[address] = &ConnectionStats{
			Address:    address,
			CreatedAt:  pc.createdAt,
			LastUsedAt: pc.lastUsedAt,
			UsageCount: pc.usageCount,
			IsClosed:   false,
		}
	}

	return stats
}

// Close закрывает все соединения в пуле
func (cp *ConnectionPool) Close() error {
	cp.mu.Lock()
	defer cp.mu.Unlock()

	close(cp.stopEviction)

	var lastErr error
	for address, pc := range cp.connections {
		if err := pc.conn.Close(); err != nil {
			cp.logger.Error("Failed to close connection",
				zap.String("address", address),
				zap.Error(err),
			)
			lastErr = err
		}
	}

	cp.connections = make(map[string]*pooledConnection)
	cp.activeCount = 0

	return lastErr
}

// GetActiveConnectionCount возвращает количество активных соединений
func (cp *ConnectionPool) GetActiveConnectionCount() int {
	cp.mu.RLock()
	defer cp.mu.RUnlock()
	return cp.activeCount
}

// GetMaxConnections возвращает максимально допустимое количество соединений
func (cp *ConnectionPool) GetMaxConnections() int {
	return cp.config.MaxConnections
}
