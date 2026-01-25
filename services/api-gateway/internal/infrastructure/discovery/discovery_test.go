package discovery

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"

	"github.com/rusgainew/kkm-project-mks/api-gateway/internal/infrastructure/observability"
)

func init() {
	observability.ResetForTesting()
}

func TestStaticResolver_Register_Resolve(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	resolver := NewStaticResolver(logger)

	endpoint := &ServiceEndpoint{Host: "localhost", Port: 9001}
	resolver.Register("auth-service", endpoint)

	resolved, err := resolver.Resolve(context.Background(), "auth-service")
	assert.NoError(t, err)
	assert.Equal(t, endpoint, resolved)
	assert.Equal(t, "localhost:9001", resolved.String())
}

func TestStaticResolver_NotFound(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	resolver := NewStaticResolver(logger)

	_, err := resolver.Resolve(context.Background(), "nonexistent")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

func TestStaticResolver_ResolveAll(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	resolver := NewStaticResolver(logger)

	endpoint := &ServiceEndpoint{Host: "localhost", Port: 9001}
	resolver.Register("service", endpoint)

	endpoints, err := resolver.ResolveAll(context.Background(), "service")
	assert.NoError(t, err)
	assert.Len(t, endpoints, 1)
	assert.Equal(t, endpoint, endpoints[0])
}

func TestDNSResolver_Resolve(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	resolver := NewDNSResolver(0, logger)

	resolved, err := resolver.Resolve(context.Background(), "localhost:9001")
	assert.NoError(t, err)
	assert.Equal(t, "localhost", resolved.Host)
	assert.Equal(t, 9001, resolved.Port)
	assert.Equal(t, "localhost:9001", resolved.String())
}

func TestDNSResolver_InvalidFormat(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	resolver := NewDNSResolver(0, logger)

	_, err := resolver.Resolve(context.Background(), "invalid-format")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid service name format")
}

func TestRoundRobinLoadBalancer(t *testing.T) {
	endpoints := []*ServiceEndpoint{
		{Host: "host1", Port: 9001},
		{Host: "host2", Port: 9002},
		{Host: "host3", Port: 9003},
	}

	lb := &RoundRobinLoadBalancer{}

	selected1 := lb.SelectEndpoint(endpoints)
	assert.Equal(t, "host1", selected1.Host)

	selected2 := lb.SelectEndpoint(endpoints)
	assert.Equal(t, "host2", selected2.Host)

	selected3 := lb.SelectEndpoint(endpoints)
	assert.Equal(t, "host3", selected3.Host)

	selected4 := lb.SelectEndpoint(endpoints)
	assert.Equal(t, "host1", selected4.Host)
}

func TestServiceRegistry_GetServiceAddress(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	resolver := NewStaticResolver(logger)
	endpoint := &ServiceEndpoint{Host: "localhost", Port: 9001}
	resolver.Register("test-service", endpoint)

	metrics := observability.NewMetrics()
	registry := NewServiceRegistry(resolver, nil, logger, metrics)

	address, err := registry.GetServiceAddress(context.Background(), "test-service")
	assert.NoError(t, err)
	assert.Equal(t, "localhost:9001", address)
}

func TestServiceEndpoint_String(t *testing.T) {
	endpoint := &ServiceEndpoint{Host: "example.com", Port: 8080}
	assert.Equal(t, "example.com:8080", endpoint.String())
}
