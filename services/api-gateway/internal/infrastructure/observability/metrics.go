// Файл api-gateway/internal/infrastructure/observability/metrics.go содержит реализацию пакета observability.
package observability

import (
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	once    sync.Once
	metrics *Metrics
)

// Metrics содержит все метрики для мониторинга
type Metrics struct {
	// HTTP metrics
	httpRequestsTotal   *prometheus.CounterVec
	httpRequestDuration *prometheus.HistogramVec
	activeConnections   prometheus.Gauge

	// gRPC metrics
	grpcCallsTotal   *prometheus.CounterVec
	grpcCallDuration *prometheus.HistogramVec
	errorsTotal      *prometheus.CounterVec

	// Redis cache metrics
	cacheHitsTotal       *prometheus.CounterVec
	cacheMissesTotal     *prometheus.CounterVec
	cacheLatency         *prometheus.HistogramVec
	cacheEvictionsTotal  *prometheus.CounterVec
	cacheEntriesCount    *prometheus.GaugeVec
	cacheMemoryBytesUsed prometheus.Gauge

	// JWT key rotation metrics
	keyRotationsTotal    *prometheus.CounterVec
	activeKeysCount      prometheus.Gauge
	keyValidationLatency *prometheus.HistogramVec

	// Token blacklist metrics
	revokedTokensTotal    *prometheus.CounterVec
	blacklistSizeGauge    prometheus.Gauge
	blacklistCheckLatency *prometheus.HistogramVec

	// Service discovery metrics
	discoveryResolutionsTotal *prometheus.CounterVec
	availableEndpointsGauge   *prometheus.GaugeVec
	selectedEndpointsTotal    *prometheus.CounterVec
	discoveryLatency          *prometheus.HistogramVec

	// Data masking metrics
	maskedFieldsTotal *prometheus.CounterVec
	maskingLatency    *prometheus.HistogramVec

	// Connection pool metrics
	activeConnectionsCount    prometheus.Gauge
	idleConnectionsCount      prometheus.Gauge
	maxConnectionsGauge       prometheus.Gauge
	connectionEvictionsTotal  *prometheus.CounterVec
	connectionPoolWaitLatency *prometheus.HistogramVec
}

// NewMetrics создает новый экземпляр метрик или возвращает существующий (singleton для тестов)
func NewMetrics() *Metrics {
	once.Do(func() {
		metrics = createMetrics()
	})
	return metrics
}

// createMetrics создает и регистрирует все метрики
func createMetrics() *Metrics {
	return &Metrics{
		// HTTP metrics
		httpRequestsTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "api_gateway_http_requests_total",
				Help: "Total number of HTTP requests",
			},
			[]string{"method", "path", "status"},
		),
		httpRequestDuration: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "api_gateway_http_request_duration_seconds",
				Help:    "HTTP request duration in seconds",
				Buckets: prometheus.DefBuckets,
			},
			[]string{"method", "path"},
		),
		activeConnections: promauto.NewGauge(
			prometheus.GaugeOpts{
				Name: "api_gateway_active_connections",
				Help: "Number of active connections",
			},
		),

		// gRPC metrics
		grpcCallsTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "api_gateway_grpc_calls_total",
				Help: "Total number of gRPC calls",
			},
			[]string{"service", "method", "status"},
		),
		grpcCallDuration: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "api_gateway_grpc_call_duration_seconds",
				Help:    "gRPC call duration in seconds",
				Buckets: prometheus.DefBuckets,
			},
			[]string{"service", "method"},
		),
		errorsTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "api_gateway_errors_total",
				Help: "Total number of errors",
			},
			[]string{"type", "service"},
		),

		// Redis cache metrics
		cacheHitsTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "api_gateway_cache_hits_total",
				Help: "Total number of cache hits",
			},
			[]string{"cache_name"},
		),
		cacheMissesTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "api_gateway_cache_misses_total",
				Help: "Total number of cache misses",
			},
			[]string{"cache_name"},
		),
		cacheLatency: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "api_gateway_cache_latency_seconds",
				Help:    "Cache operation latency in seconds",
				Buckets: []float64{0.001, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1},
			},
			[]string{"operation", "cache_name"},
		),
		cacheEvictionsTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "api_gateway_cache_evictions_total",
				Help: "Total number of cache evictions",
			},
			[]string{"cache_name"},
		),
		cacheEntriesCount: promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "api_gateway_cache_entries_count",
				Help: "Current number of entries in cache",
			},
			[]string{"cache_name"},
		),
		cacheMemoryBytesUsed: promauto.NewGauge(
			prometheus.GaugeOpts{
				Name: "api_gateway_cache_memory_bytes_used",
				Help: "Memory used by cache in bytes",
			},
		),

		// JWT key rotation metrics
		keyRotationsTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "api_gateway_jwt_key_rotations_total",
				Help: "Total number of JWT key rotations",
			},
			[]string{"status"},
		),
		activeKeysCount: promauto.NewGauge(
			prometheus.GaugeOpts{
				Name: "api_gateway_jwt_active_keys_count",
				Help: "Number of currently active JWT keys",
			},
		),
		keyValidationLatency: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "api_gateway_jwt_validation_latency_seconds",
				Help:    "JWT validation latency in seconds",
				Buckets: []float64{0.0001, 0.0005, 0.001, 0.005, 0.01, 0.05, 0.1},
			},
			[]string{"result"},
		),

		// Token blacklist metrics
		revokedTokensTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "api_gateway_revoked_tokens_total",
				Help: "Total number of revoked tokens",
			},
			[]string{"reason"},
		),
		blacklistSizeGauge: promauto.NewGauge(
			prometheus.GaugeOpts{
				Name: "api_gateway_blacklist_size",
				Help: "Current size of token blacklist",
			},
		),
		blacklistCheckLatency: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "api_gateway_blacklist_check_latency_seconds",
				Help:    "Token blacklist check latency in seconds",
				Buckets: []float64{0.0001, 0.0005, 0.001, 0.005, 0.01, 0.05},
			},
			[]string{"result"},
		),

		// Service discovery metrics
		discoveryResolutionsTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "api_gateway_discovery_resolutions_total",
				Help: "Total number of service discovery resolutions",
			},
			[]string{"service", "status"},
		),
		availableEndpointsGauge: promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "api_gateway_discovery_available_endpoints",
				Help: "Number of available endpoints for service",
			},
			[]string{"service"},
		),
		selectedEndpointsTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "api_gateway_discovery_selected_endpoints_total",
				Help: "Total number of selected endpoints",
			},
			[]string{"service", "endpoint"},
		),
		discoveryLatency: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "api_gateway_discovery_latency_seconds",
				Help:    "Service discovery latency in seconds",
				Buckets: []float64{0.001, 0.005, 0.01, 0.05, 0.1, 0.5, 1},
			},
			[]string{"service"},
		),

		// Data masking metrics
		maskedFieldsTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "api_gateway_masked_fields_total",
				Help: "Total number of masked fields",
			},
			[]string{"field_type"},
		),
		maskingLatency: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "api_gateway_masking_latency_seconds",
				Help:    "Data masking latency in seconds",
				Buckets: []float64{0.00001, 0.0001, 0.0005, 0.001, 0.005, 0.01},
			},
			[]string{"field_type"},
		),

		// Connection pool metrics
		activeConnectionsCount: promauto.NewGauge(
			prometheus.GaugeOpts{
				Name: "api_gateway_pool_active_connections",
				Help: "Number of active connections in pool",
			},
		),
		idleConnectionsCount: promauto.NewGauge(
			prometheus.GaugeOpts{
				Name: "api_gateway_pool_idle_connections",
				Help: "Number of idle connections in pool",
			},
		),
		maxConnectionsGauge: promauto.NewGauge(
			prometheus.GaugeOpts{
				Name: "api_gateway_pool_max_connections",
				Help: "Maximum connections allowed in pool",
			},
		),
		connectionEvictionsTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "api_gateway_pool_evictions_total",
				Help: "Total number of connection evictions",
			},
			[]string{"reason"},
		),
		connectionPoolWaitLatency: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "api_gateway_pool_wait_latency_seconds",
				Help:    "Connection pool wait latency in seconds",
				Buckets: []float64{0.001, 0.005, 0.01, 0.05, 0.1, 0.5},
			},
			[]string{"status"},
		),
	}
}

// RecordHTTPRequest записывает метрику HTTP запроса
func (m *Metrics) RecordHTTPRequest(method, path, status string, duration float64) {
	m.httpRequestsTotal.WithLabelValues(method, path, status).Inc()
	m.httpRequestDuration.WithLabelValues(method, path).Observe(duration)
}

// RecordGRPCRequest записывает метрику gRPC запроса
func (m *Metrics) RecordGRPCRequest(service, method string, duration float64, err error) {
	status := "success"
	if err != nil {
		status = "error"
	}
	m.grpcCallsTotal.WithLabelValues(service, method, status).Inc()
	m.grpcCallDuration.WithLabelValues(service, method).Observe(duration)
}

// IncrementErrorCount увеличивает счетчик ошибок
func (m *Metrics) IncrementErrorCount(errorType, service string) {
	m.errorsTotal.WithLabelValues(errorType, service).Inc()
}

// IncrementActiveConnections увеличивает счетчик активных подключений
func (m *Metrics) IncrementActiveConnections() {
	m.activeConnections.Inc()
}

// DecrementActiveConnections уменьшает счетчик активных подключений
func (m *Metrics) DecrementActiveConnections() {
	m.activeConnections.Dec()
}

// Cache metrics

// RecordCacheHit записывает попадание в кеш
func (m *Metrics) RecordCacheHit(cacheName string) {
	m.cacheHitsTotal.WithLabelValues(cacheName).Inc()
}

// RecordCacheMiss записывает промах кеша
func (m *Metrics) RecordCacheMiss(cacheName string) {
	m.cacheMissesTotal.WithLabelValues(cacheName).Inc()
}

// RecordCacheLatency записывает задержку кеша
func (m *Metrics) RecordCacheLatency(operation, cacheName string, latency float64) {
	m.cacheLatency.WithLabelValues(operation, cacheName).Observe(latency)
}

// RecordCacheEviction записывает вытеснение из кеша
func (m *Metrics) RecordCacheEviction(cacheName string) {
	m.cacheEvictionsTotal.WithLabelValues(cacheName).Inc()
}

// SetCacheEntriesCount устанавливает количество элементов в кеше
func (m *Metrics) SetCacheEntriesCount(cacheName string, count int64) {
	m.cacheEntriesCount.WithLabelValues(cacheName).Set(float64(count))
}

// SetCacheMemoryUsage устанавливает использованную памяти кешем
func (m *Metrics) SetCacheMemoryUsage(bytes int64) {
	m.cacheMemoryBytesUsed.Set(float64(bytes))
}

// JWT Key Rotation metrics

// RecordKeyRotation записывает ротацию ключа
func (m *Metrics) RecordKeyRotation(status string) {
	m.keyRotationsTotal.WithLabelValues(status).Inc()
}

// SetActiveKeysCount устанавливает количество активных ключей
func (m *Metrics) SetActiveKeysCount(count int64) {
	m.activeKeysCount.Set(float64(count))
}

// RecordKeyValidationLatency записывает задержку валидации ключа
func (m *Metrics) RecordKeyValidationLatency(result string, latency float64) {
	m.keyValidationLatency.WithLabelValues(result).Observe(latency)
}

// Token Blacklist metrics

// RecordRevokedToken записывает отозванный токен
func (m *Metrics) RecordRevokedToken(reason string) {
	m.revokedTokensTotal.WithLabelValues(reason).Inc()
}

// SetBlacklistSize устанавливает размер черного списка
func (m *Metrics) SetBlacklistSize(size int64) {
	m.blacklistSizeGauge.Set(float64(size))
}

// RecordBlacklistCheckLatency записывает задержку проверки черного списка
func (m *Metrics) RecordBlacklistCheckLatency(result string, latency float64) {
	m.blacklistCheckLatency.WithLabelValues(result).Observe(latency)
}

// Service Discovery metrics

// RecordDiscoveryResolution записывает разрешение сервиса
func (m *Metrics) RecordDiscoveryResolution(service, status string) {
	m.discoveryResolutionsTotal.WithLabelValues(service, status).Inc()
}

// SetAvailableEndpoints устанавливает количество доступных эндпоинтов
func (m *Metrics) SetAvailableEndpoints(service string, count int64) {
	m.availableEndpointsGauge.WithLabelValues(service).Set(float64(count))
}

// RecordSelectedEndpoint записывает выбранный эндпоинт
func (m *Metrics) RecordSelectedEndpoint(service, endpoint string) {
	m.selectedEndpointsTotal.WithLabelValues(service, endpoint).Inc()
}

// RecordDiscoveryLatency записывает задержку обнаружения сервиса
func (m *Metrics) RecordDiscoveryLatency(service string, latency float64) {
	m.discoveryLatency.WithLabelValues(service).Observe(latency)
}

// Data Masking metrics

// RecordMaskedField записывает замаскированное поле
func (m *Metrics) RecordMaskedField(fieldType string) {
	m.maskedFieldsTotal.WithLabelValues(fieldType).Inc()
}

// RecordMaskingLatency записывает задержку маскирования
func (m *Metrics) RecordMaskingLatency(fieldType string, latency float64) {
	m.maskingLatency.WithLabelValues(fieldType).Observe(latency)
}

// Connection Pool metrics

// SetActiveConnectionsCount устанавливает количество активных соединений
func (m *Metrics) SetActiveConnectionsCount(count int64) {
	m.activeConnectionsCount.Set(float64(count))
}

// SetIdleConnectionsCount устанавливает количество неактивных соединений
func (m *Metrics) SetIdleConnectionsCount(count int64) {
	m.idleConnectionsCount.Set(float64(count))
}

// ResetForTesting сбрасывает состояние для тестирования
func ResetForTesting() {
	once = sync.Once{}
	metrics = nil
} // SetMaxConnections устанавливает максимальное количество соединений
func (m *Metrics) SetMaxConnections(max int64) {
	m.maxConnectionsGauge.Set(float64(max))
}

// RecordConnectionEviction записывает вытеснение соединения
func (m *Metrics) RecordConnectionEviction(reason string) {
	m.connectionEvictionsTotal.WithLabelValues(reason).Inc()
}

// RecordConnectionPoolWaitLatency записывает задержку ожидания в пуле
func (m *Metrics) RecordConnectionPoolWaitLatency(status string, latency float64) {
	m.connectionPoolWaitLatency.WithLabelValues(status).Observe(latency)
}

// Analytics Repository metrics

// RecordQueryLatency записывает задержку SQL запроса для аналитики
func (m *Metrics) RecordQueryLatency(operation string, latency float64) {
	// Используем существующую метрику gRPC для database запросов
	m.grpcCallDuration.WithLabelValues("analytics-repository", operation).Observe(latency)
}

// RecordAnalyticsError увеличивает счетчик ошибок для аналитики
func (m *Metrics) RecordAnalyticsError(operation, errorType string) {
	// Используем существующую метрику ошибок
	m.errorsTotal.WithLabelValues("analytics-repository", operation, errorType).Inc()
}
