// Файл invoice-query-server/cmd/main.go содержит реализацию пакета main.
package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rusgainew/kkm-project-mks/invoice-query-server/internal/infrastructure/cache"
	"github.com/rusgainew/kkm-project-mks/invoice-query-server/internal/infrastructure/config"
	apphealth "github.com/rusgainew/kkm-project-mks/invoice-query-server/internal/infrastructure/health"
	"github.com/rusgainew/kkm-project-mks/invoice-query-server/internal/infrastructure/middleware"
	"github.com/rusgainew/kkm-project-mks/invoice-query-server/internal/infrastructure/observability"
	"github.com/rusgainew/kkm-project-mks/invoice-query-server/internal/infrastructure/repository"
	"github.com/rusgainew/kkm-project-mks/invoice-query-server/internal/interfaces/grpc/handlers"
	pb "github.com/rusgainew/kkm-project-mks/proto-lib/api"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
)

func main() {
	// Load and validate config
	cfg, err := config.LoadValidated()
	if err != nil {
		log.Fatalf("Invalid configuration: %v", err)
	}

	// Setup logger
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	// Initialize OpenTelemetry tracer
	tp, err := observability.InitTracer("invoice-query-service", cfg.JaegerEndpoint)
	if err != nil {
		logger.Warn("Failed to initialize tracer", zap.Error(err))
	} else {
		defer func() {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			if err := observability.ShutdownTracer(ctx, tp); err != nil {
				logger.Error("Failed to shutdown tracer", zap.Error(err))
			}
		}()
		logger.Info("OpenTelemetry tracer initialized", zap.String("jaeger_endpoint", cfg.JaegerEndpoint))
	}

	// Create in-memory repository (data will come from RabbitMQ events)
	repo := repository.NewInMemoryInvoiceRepository(logger)

	// Initialize metrics
	metrics := observability.NewMetricsCollector("invoice_query", "service")

	// Initialize Redis cache
	redisCache, err := cache.NewRedisCache(cfg.RedisURL, cfg.CacheTTL, logger)
	if err != nil {
		logger.Warn("Failed to initialize Redis cache, continuing without cache", zap.Error(err))
		redisCache = nil
	} else {
		defer redisCache.Close()
		logger.Info("Redis cache initialized", zap.String("url", cfg.RedisURL), zap.Duration("ttl", cfg.CacheTTL))
	}

	// JWT secret from environment (required)
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		logger.Fatal("JWT_SECRET environment variable is required")
	}
	authMiddleware := middleware.NewAuthMiddleware(jwtSecret, logger)

	// Create gRPC server with JWT interceptor
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(authMiddleware.UnaryServerInterceptor()),
	)

	// Health check
	healthServer := health.NewServer()
	grpc_health_v1.RegisterHealthServer(grpcServer, healthServer)
	healthServer.SetServingStatus("api.InvoiceQueryService", grpc_health_v1.HealthCheckResponse_SERVING)

	// Register handler with repository, metrics, and cache
	handler := handlers.NewInvoiceQueryHandler(logger, repo, metrics, redisCache)
	pb.RegisterInvoiceQueryServiceServer(grpcServer, handler)

	// Initialize health check manager
	healthManager := apphealth.NewManager(logger)

	// Get the concrete repository type for health checks
	inMemoryRepo, ok := repo.(*repository.InMemoryInvoiceRepository)
	if ok {
		// Register cache staleness check (consider stale if no updates in 1 hour)
		healthManager.RegisterCheck("cache_staleness",
			apphealth.CacheStalenessChecker(inMemoryRepo.GetLastUpdate, 1*time.Hour))

		// Register cache size check (expect at least 10 items, max 100000)
		healthManager.RegisterCheck("cache_size",
			apphealth.CacheSizeChecker(inMemoryRepo.GetCacheSize, 10, 100000))
	}

	// Listen on gRPC port
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.GRPCPort))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	// Health check endpoints
	healthMux := http.NewServeMux()
	healthMux.HandleFunc("/health", healthManager.HealthHandler())
	healthMux.HandleFunc("/health/live", healthManager.LivenessHandler())
	healthMux.HandleFunc("/health/ready", healthManager.ReadinessHandler())
	// Metrics endpoint
	healthMux.Handle("/metrics", promhttp.Handler())

	httpServer := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.HealthPort),
		Handler:      healthMux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  30 * time.Second,
	}

	go func() {
		logger.Info("Starting health check server", zap.String("addr", httpServer.Addr))
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("Health check server error", zap.Error(err))
		}
	}()

	logger.Info("Invoice Query Server starting", zap.Int("grpc_port", cfg.GRPCPort))

	// Run gRPC server in background
	go func() {
		if err := grpcServer.Serve(listener); err != nil {
			logger.Fatal("gRPC server error", zap.Error(err))
		}
	}()

	// Wait for interrupt
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	logger.Info("Shutdown signal received, gracefully stopping servers")

	// Create context with timeout for graceful shutdown
	shutdownCtx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
	defer cancel()

	// Shutdown gRPC server gracefully
	stopped := make(chan struct{})
	go func() {
		grpcServer.GracefulStop()
		close(stopped)
	}()

	// Wait for gRPC shutdown or timeout
	select {
	case <-stopped:
		logger.Info("gRPC server stopped gracefully")
	case <-shutdownCtx.Done():
		logger.Warn("gRPC graceful shutdown timeout exceeded, forcing stop")
		grpcServer.Stop()
	}

	// Shutdown HTTP server gracefully
	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		logger.Error("HTTP server shutdown error", zap.Error(err))
	} else {
		logger.Info("HTTP server stopped gracefully")
	}

	logger.Info("Invoice Query Server stopped")
}
