package health

import (
	"context"
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"go.uber.org/zap"
)

var (
	healthStatus = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "service_health_status",
		Help: "Health status of the service (1=healthy, 0=unhealthy)",
	}, []string{"check_type"})
)

// Status represents the health status
type Status string

const (
	StatusHealthy   Status = "healthy"
	StatusUnhealthy Status = "unhealthy"
	StatusDegraded  Status = "degraded"
)

// CheckResult represents a single health check result
type CheckResult struct {
	Status      Status                 `json:"status"`
	Message     string                 `json:"message,omitempty"`
	LastChecked time.Time              `json:"last_checked"`
	Details     map[string]interface{} `json:"details,omitempty"`
}

// HealthResponse represents the overall health response
type HealthResponse struct {
	Status    Status                 `json:"status"`
	Timestamp time.Time              `json:"timestamp"`
	Checks    map[string]CheckResult `json:"checks"`
	Uptime    time.Duration          `json:"uptime"`
}

// Checker defines the interface for health checks
type Checker interface {
	Check(ctx context.Context) CheckResult
}

// CheckFunc is a function adapter for Checker interface
type CheckFunc func(ctx context.Context) CheckResult

func (f CheckFunc) Check(ctx context.Context) CheckResult {
	return f(ctx)
}

// Manager manages health checks
type Manager struct {
	mu        sync.RWMutex
	checks    map[string]Checker
	startTime time.Time
	logger    *zap.Logger
}

// NewManager creates a new health check manager
func NewManager(logger *zap.Logger) *Manager {
	return &Manager{
		checks:    make(map[string]Checker),
		startTime: time.Now(),
		logger:    logger,
	}
}

// RegisterCheck registers a health check
func (m *Manager) RegisterCheck(name string, checker Checker) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.checks[name] = checker
	m.logger.Info("Registered health check", zap.String("name", name))
}

// RunChecks executes all health checks
func (m *Manager) RunChecks(ctx context.Context) HealthResponse {
	m.mu.RLock()
	checks := make(map[string]Checker, len(m.checks))
	for name, checker := range m.checks {
		checks[name] = checker
	}
	m.mu.RUnlock()

	results := make(map[string]CheckResult)
	overallStatus := StatusHealthy

	for name, checker := range checks {
		result := checker.Check(ctx)
		results[name] = result

		// Update Prometheus metrics
		statusValue := 1.0
		if result.Status != StatusHealthy {
			statusValue = 0.0
			if result.Status == StatusUnhealthy {
				overallStatus = StatusUnhealthy
			} else if overallStatus == StatusHealthy {
				overallStatus = StatusDegraded
			}
		}
		healthStatus.WithLabelValues(name).Set(statusValue)
	}

	return HealthResponse{
		Status:    overallStatus,
		Timestamp: time.Now(),
		Checks:    results,
		Uptime:    time.Since(m.startTime),
	}
}

// LivenessHandler returns an HTTP handler for liveness probe
// Liveness checks if the service is running (basic check)
func (m *Manager) LivenessHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status": "alive",
			"uptime": time.Since(m.startTime).Seconds(),
		})
	}
}

// ReadinessHandler returns an HTTP handler for readiness probe
// Readiness checks if the service is ready to accept traffic
func (m *Manager) ReadinessHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		response := m.RunChecks(ctx)

		w.Header().Set("Content-Type", "application/json")
		if response.Status == StatusHealthy {
			w.WriteHeader(http.StatusOK)
		} else if response.Status == StatusDegraded {
			w.WriteHeader(http.StatusOK) // Degraded but still serving
		} else {
			w.WriteHeader(http.StatusServiceUnavailable)
		}

		json.NewEncoder(w).Encode(response)
	}
}

// HealthHandler returns an HTTP handler for detailed health check
func (m *Manager) HealthHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
		defer cancel()

		response := m.RunChecks(ctx)

		w.Header().Set("Content-Type", "application/json")
		if response.Status == StatusHealthy || response.Status == StatusDegraded {
			w.WriteHeader(http.StatusOK)
		} else {
			w.WriteHeader(http.StatusServiceUnavailable)
		}

		json.NewEncoder(w).Encode(response)
	}
}

// CacheStalenessChecker creates a checker for cache staleness
func CacheStalenessChecker(getLastUpdate func() time.Time, maxAge time.Duration) Checker {
	return CheckFunc(func(ctx context.Context) CheckResult {
		lastUpdate := getLastUpdate()

		if lastUpdate.IsZero() {
			return CheckResult{
				Status:      StatusDegraded,
				Message:     "Cache has not been initialized",
				LastChecked: time.Now(),
				Details: map[string]interface{}{
					"last_update": "never",
				},
			}
		}

		age := time.Since(lastUpdate)

		if age > maxAge {
			return CheckResult{
				Status:      StatusDegraded,
				Message:     "Cache data is stale",
				LastChecked: time.Now(),
				Details: map[string]interface{}{
					"last_update":     lastUpdate.Format(time.RFC3339),
					"age_seconds":     age.Seconds(),
					"max_age_seconds": maxAge.Seconds(),
				},
			}
		}

		return CheckResult{
			Status:      StatusHealthy,
			Message:     "Cache is up to date",
			LastChecked: time.Now(),
			Details: map[string]interface{}{
				"last_update": lastUpdate.Format(time.RFC3339),
				"age_seconds": age.Seconds(),
			},
		}
	})
}

// CacheSizeChecker creates a checker for cache size
func CacheSizeChecker(getSize func() int, minSize, maxSize int) Checker {
	return CheckFunc(func(ctx context.Context) CheckResult {
		size := getSize()

		if size == 0 {
			return CheckResult{
				Status:      StatusDegraded,
				Message:     "Cache is empty",
				LastChecked: time.Now(),
				Details: map[string]interface{}{
					"size": size,
				},
			}
		}

		if size > maxSize {
			return CheckResult{
				Status:      StatusDegraded,
				Message:     "Cache size exceeds maximum",
				LastChecked: time.Now(),
				Details: map[string]interface{}{
					"size":     size,
					"max_size": maxSize,
				},
			}
		}

		if size < minSize {
			return CheckResult{
				Status:      StatusDegraded,
				Message:     "Cache size below minimum",
				LastChecked: time.Now(),
				Details: map[string]interface{}{
					"size":     size,
					"min_size": minSize,
				},
			}
		}

		return CheckResult{
			Status:      StatusHealthy,
			Message:     "Cache size is within acceptable range",
			LastChecked: time.Now(),
			Details: map[string]interface{}{
				"size": size,
			},
		}
	})
}
