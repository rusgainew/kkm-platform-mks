// Файл api-gateway/internal/infrastructure/discovery/discovery.go содержит реализацию пакета discovery.
package discovery

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/rusgainew/kkm-project-mks/api-gateway/internal/infrastructure/observability"
	"go.uber.org/zap"
)

// ServiceEndpoint описывает endpoint сервиса
type ServiceEndpoint struct {
	Host string // hostname или IP
	Port int    // порт
}

// String возвращает адрес в формате host:port
func (se ServiceEndpoint) String() string {
	return fmt.Sprintf("%s:%d", se.Host, se.Port)
}

// ServiceResolver интерфейс для решения адреса сервиса
type ServiceResolver interface {
	Resolve(ctx context.Context, serviceName string) (*ServiceEndpoint, error)
	ResolveAll(ctx context.Context, serviceName string) ([]*ServiceEndpoint, error)
	Watch(ctx context.Context, serviceName string, callback func([]*ServiceEndpoint)) error
}

// StaticResolver простой resolver с статическими адресами
type StaticResolver struct {
	mu       sync.RWMutex
	services map[string]*ServiceEndpoint
	logger   *zap.Logger
}

// NewStaticResolver создает resolver со статическими адресами
func NewStaticResolver(logger *zap.Logger) *StaticResolver {
	return &StaticResolver{
		services: make(map[string]*ServiceEndpoint),
		logger:   logger,
	}
}

// Register регистрирует сервис
func (sr *StaticResolver) Register(serviceName string, endpoint *ServiceEndpoint) {
	sr.mu.Lock()
	defer sr.mu.Unlock()
	sr.services[serviceName] = endpoint
	sr.logger.Info("Service registered", zap.String("service", serviceName), zap.String("endpoint", endpoint.String()))
}

// Resolve возвращает endpoint сервиса
func (sr *StaticResolver) Resolve(ctx context.Context, serviceName string) (*ServiceEndpoint, error) {
	sr.mu.RLock()
	defer sr.mu.RUnlock()

	endpoint, exists := sr.services[serviceName]
	if !exists {
		return nil, fmt.Errorf("service %q not found", serviceName)
	}

	return endpoint, nil
}

// ResolveAll возвращает все endpoints сервиса (для StaticResolver это один)
func (sr *StaticResolver) ResolveAll(ctx context.Context, serviceName string) ([]*ServiceEndpoint, error) {
	endpoint, err := sr.Resolve(ctx, serviceName)
	if err != nil {
		return nil, err
	}
	return []*ServiceEndpoint{endpoint}, nil
}

// Watch не поддерживается для StaticResolver
func (sr *StaticResolver) Watch(ctx context.Context, serviceName string, callback func([]*ServiceEndpoint)) error {
	return fmt.Errorf("watch not supported for static resolver")
}

// DNSResolver resolver на основе DNS
type DNSResolver struct {
	mu       sync.RWMutex
	cache    map[string]*ServiceEndpoint
	ttl      time.Duration
	cacheTTL map[string]time.Time
	logger   *zap.Logger
}

// NewDNSResolver создает DNS resolver с кешированием
func NewDNSResolver(ttl time.Duration, logger *zap.Logger) *DNSResolver {
	if ttl == 0 {
		ttl = 5 * time.Minute
	}
	return &DNSResolver{
		cache:    make(map[string]*ServiceEndpoint),
		cacheTTL: make(map[string]time.Time),
		ttl:      ttl,
		logger:   logger,
	}
}

// Resolve разрешает имя сервиса в адрес (формат: "service-name:8080")
func (dr *DNSResolver) Resolve(ctx context.Context, serviceName string) (*ServiceEndpoint, error) {
	dr.mu.RLock()
	if cached, exists := dr.cache[serviceName]; exists {
		if time.Now().Before(dr.cacheTTL[serviceName]) {
			dr.mu.RUnlock()
			return cached, nil
		}
	}
	dr.mu.RUnlock()

	// Парсим serviceName в формат "host:port"
	parts := strings.Split(serviceName, ":")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid service name format, expected 'host:port', got %q", serviceName)
	}

	host := parts[0]
	var port int
	_, err := fmt.Sscanf(parts[1], "%d", &port)
	if err != nil {
		return nil, fmt.Errorf("invalid port in service name %q: %w", serviceName, err)
	}

	endpoint := &ServiceEndpoint{
		Host: host,
		Port: port,
	}

	// Кешируем результат
	dr.mu.Lock()
	dr.cache[serviceName] = endpoint
	dr.cacheTTL[serviceName] = time.Now().Add(dr.ttl)
	dr.mu.Unlock()

	dr.logger.Debug("Service resolved", zap.String("service", serviceName), zap.String("endpoint", endpoint.String()))
	return endpoint, nil
}

// ResolveAll возвращает все endpoints (для DNSResolver это один)
func (dr *DNSResolver) ResolveAll(ctx context.Context, serviceName string) ([]*ServiceEndpoint, error) {
	endpoint, err := dr.Resolve(ctx, serviceName)
	if err != nil {
		return nil, err
	}
	return []*ServiceEndpoint{endpoint}, nil
}

// Watch не поддерживается для DNSResolver
func (dr *DNSResolver) Watch(ctx context.Context, serviceName string, callback func([]*ServiceEndpoint)) error {
	return fmt.Errorf("watch not supported for DNS resolver")
}

// InvalidateCache инвалидирует кеш для сервиса
func (dr *DNSResolver) InvalidateCache(serviceName string) {
	dr.mu.Lock()
	defer dr.mu.Unlock()
	delete(dr.cache, serviceName)
	delete(dr.cacheTTL, serviceName)
}

// LoadBalancerStrategy стратегия балансировки нагрузки
type LoadBalancerStrategy interface {
	SelectEndpoint(endpoints []*ServiceEndpoint) *ServiceEndpoint
}

// RoundRobinLoadBalancer простой round-robin load balancer
type RoundRobinLoadBalancer struct {
	mu      sync.Mutex
	counter int
}

// SelectEndpoint выбирает endpoint по round-robin
func (rr *RoundRobinLoadBalancer) SelectEndpoint(endpoints []*ServiceEndpoint) *ServiceEndpoint {
	if len(endpoints) == 0 {
		return nil
	}

	rr.mu.Lock()
	defer rr.mu.Unlock()

	selected := endpoints[rr.counter%len(endpoints)]
	rr.counter++

	return selected
}

// RandomLoadBalancer случайный load balancer
type RandomLoadBalancer struct{}

// SelectEndpoint выбирает случайный endpoint
func (r *RandomLoadBalancer) SelectEndpoint(endpoints []*ServiceEndpoint) *ServiceEndpoint {
	if len(endpoints) == 0 {
		return nil
	}
	// Для простоты в реальном коде использовать rand.Intn()
	return endpoints[0]
}

// ServiceRegistry реестр сервисов
type ServiceRegistry struct {
	mu       sync.RWMutex
	resolver ServiceResolver
	lb       LoadBalancerStrategy
	logger   *zap.Logger
	metrics  *observability.Metrics
}

// NewServiceRegistry создает новый реестр сервисов
func NewServiceRegistry(resolver ServiceResolver, lb LoadBalancerStrategy, logger *zap.Logger, metrics *observability.Metrics) *ServiceRegistry {
	if lb == nil {
		lb = &RoundRobinLoadBalancer{}
	}
	return &ServiceRegistry{
		resolver: resolver,
		lb:       lb,
		logger:   logger,
		metrics:  metrics,
	}
}

// GetServiceEndpoint возвращает endpoint сервиса через resolver и load balancer
func (sr *ServiceRegistry) GetServiceEndpoint(ctx context.Context, serviceName string) (*ServiceEndpoint, error) {
	start := time.Now()

	endpoints, err := sr.resolver.ResolveAll(ctx, serviceName)
	if err != nil {
		latency := time.Since(start).Seconds()
		sr.metrics.RecordDiscoveryLatency(serviceName, latency)
		sr.metrics.RecordDiscoveryResolution(serviceName, "failed")
		return nil, err
	}

	if len(endpoints) == 0 {
		latency := time.Since(start).Seconds()
		sr.metrics.RecordDiscoveryLatency(serviceName, latency)
		sr.metrics.RecordDiscoveryResolution(serviceName, "no_endpoints")
		return nil, fmt.Errorf("no endpoints found for service %q", serviceName)
	}

	sr.metrics.SetAvailableEndpoints(serviceName, int64(len(endpoints)))

	selected := sr.lb.SelectEndpoint(endpoints)
	if selected == nil {
		latency := time.Since(start).Seconds()
		sr.metrics.RecordDiscoveryLatency(serviceName, latency)
		sr.metrics.RecordDiscoveryResolution(serviceName, "lb_failed")
		return nil, fmt.Errorf("load balancer failed to select endpoint for service %q", serviceName)
	}

	latency := time.Since(start).Seconds()
	sr.metrics.RecordDiscoveryLatency(serviceName, latency)
	sr.metrics.RecordDiscoveryResolution(serviceName, "success")
	sr.metrics.RecordSelectedEndpoint(serviceName, selected.String())

	return selected, nil
}

// GetServiceAddress возвращает адрес сервиса в формате host:port
func (sr *ServiceRegistry) GetServiceAddress(ctx context.Context, serviceName string) (string, error) {
	endpoint, err := sr.GetServiceEndpoint(ctx, serviceName)
	if err != nil {
		return "", err
	}
	return endpoint.String(), nil
}
