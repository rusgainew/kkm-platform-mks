package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"

	"github.com/rusgainew/kkm-project-mks/document-query-server/internal/application/handlers"
	"github.com/rusgainew/kkm-project-mks/document-query-server/internal/infrastructure/cache"
	"github.com/rusgainew/kkm-project-mks/document-query-server/internal/infrastructure/config"
	"github.com/rusgainew/kkm-project-mks/document-query-server/internal/infrastructure/messaging"
	"github.com/rusgainew/kkm-project-mks/document-query-server/internal/infrastructure/repository"
	grpchandlers "github.com/rusgainew/kkm-project-mks/document-query-server/internal/interfaces/grpc/handlers"
	pb "github.com/rusgainew/kkm-project-mks/proto-lib/api"
)

func main() {
	// Load configuration
	cfg, err := config.LoadValidated()
	if err != nil {
		fmt.Printf("Failed to load config: %v\n", err)
		os.Exit(1)
	}

	// Initialize logger
	logger, err := initLogger(cfg.LogLevel)
	if err != nil {
		fmt.Printf("Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer logger.Sync()

	logger.Info("Starting Document Query Server",
		zap.Int("grpc_port", cfg.GRPCPort),
		zap.Int("metrics_port", cfg.MetricsPort),
	)

	// Initialize in-memory repository (data will come from RabbitMQ events)
	documentRepo := repository.NewInMemoryDocumentRepository(logger)
	logger.Info("In-memory document repository initialized")

	// Initialize Redis cache (optional)
	// Create application lifecycle context
	appCtx, appCancel := context.WithCancel(context.Background())
	defer appCancel()

	var redisCache *cache.RedisCache
	if cfg.RedisURL != "" {
		redisClient := redis.NewClient(&redis.Options{
			Addr: cfg.RedisURL,
		})

		// Use context with timeout for ping
		pingCtx, pingCancel := context.WithTimeout(appCtx, 5*time.Second)
		if err := redisClient.Ping(pingCtx).Err(); err == nil {
			redisCache = cache.NewRedisCache(redisClient, logger)
			logger.Info("Redis cache initialized")
		} else {
			logger.Warn("Failed to connect to Redis, continuing without caching", zap.Error(err))
		}
		pingCancel()
	}

	// Initialize gRPC handler
	documentQueryHandler := grpchandlers.NewDocumentQueryHandler(logger, documentRepo, redisCache)

	// Initialize RabbitMQ consumer for event-driven updates
	var consumer *messaging.RabbitMQDocumentConsumer
	if cfg.RabbitMQEnabled {
		eventHandler := handlers.NewDocumentEventHandler(documentRepo, logger)
		consumer, err = messaging.NewRabbitMQDocumentConsumer(
			cfg.RabbitMQURL,
			"document-query-server",
			cfg.RabbitMQExchange,
			logger,
			eventHandler,
		)
		if err != nil {
			logger.Error("Failed to create RabbitMQ consumer", zap.Error(err))
			logger.Info("Continuing without event consumer")
		} else {
			defer consumer.Close()
			// Start listening for events in background with app context
			go func() {
				if err := consumer.Start(appCtx); err != nil && err != context.Canceled {
					logger.Error("Failed to start RabbitMQ consumer", zap.Error(err))
				}
			}()
		}
	} else {
		logger.Info("RabbitMQ event consumer disabled")
	}

	// Set up gRPC server
	grpcListener, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.GRPCPort))
	if err != nil {
		logger.Fatal("Failed to listen on gRPC port", zap.Error(err))
	}
	defer grpcListener.Close()

	grpcServer := grpc.NewServer()

	// Register health check
	healthService := health.NewServer()
	grpc_health_v1.RegisterHealthServer(grpcServer, healthService)
	healthService.SetServingStatus("api.DocumentQueryService", grpc_health_v1.HealthCheckResponse_SERVING)

	// Register Document Query Service
	pb.RegisterDocumentQueryServiceServer(grpcServer, documentQueryHandler)

	reflection.Register(grpcServer)

	logger.Info("gRPC server registered", zap.Int("port", cfg.GRPCPort))

	// Set up metrics HTTP server
	metricsMux := http.NewServeMux()
	metricsMux.Handle("/metrics", promhttp.Handler())
	metricsServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.MetricsPort),
		Handler: metricsMux,
	}

	// Start servers
	go func() {
		if err := grpcServer.Serve(grpcListener); err != nil {
			logger.Error("gRPC server error", zap.Error(err))
		}
	}()

	go func() {
		if err := metricsServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("Metrics server error", zap.Error(err))
		}
	}()

	logger.Info("Servers started successfully",
		zap.Int("grpc_port", cfg.GRPCPort),
		zap.Int("metrics_port", cfg.MetricsPort),
	)

	// Wait for shutdown signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	logger.Info("Shutdown signal received, gracefully stopping servers...")

	// Cancel application context to stop consumers
	appCancel()

	// Shutdown gRPC server
	grpcServer.GracefulStop()

	// Shutdown metrics server
	ctx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
	if err := metricsServer.Shutdown(ctx); err != nil {
		logger.Error("Error shutting down metrics server", zap.Error(err))
	}
	cancel()

	logger.Info("Document Query Server stopped")
}

func initLogger(logLevel string) (*zap.Logger, error) {
	switch logLevel {
	case "debug":
		return zap.NewDevelopment()
	case "info":
		return zap.NewProduction()
	default:
		return zap.NewProduction()
	}
}
