// Файл analytics-server/cmd/main.go содержит реализацию пакета main.
package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/lib/pq"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"

	"github.com/rusgainew/kkm-project-mks/analytics-server/internal/application/services"
	"github.com/rusgainew/kkm-project-mks/analytics-server/internal/infrastructure/cache"
	"github.com/rusgainew/kkm-project-mks/analytics-server/internal/infrastructure/config"
	"github.com/rusgainew/kkm-project-mks/analytics-server/internal/infrastructure/repository"
	grpcHandler "github.com/rusgainew/kkm-project-mks/analytics-server/internal/interfaces/grpc"
)

func main() {
	// Создание логгера
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("Failed to create logger: %v", err)
	}
	defer logger.Sync()

	// Загрузка конфигурации
	cfg, err := config.Load()
	if err != nil {
		logger.Fatal("Failed to load config", zap.Error(err))
	}

	logger.Info("Starting analytics-server",
		zap.String("grpc_port", cfg.Server.GRPCPort),
		zap.String("metrics_port", cfg.Server.MetricsPort),
	)

	// Подключение к базе данных
	db, err := initDatabase(&cfg.Database, logger)
	if err != nil {
		logger.Fatal("Failed to connect to database", zap.Error(err))
	}
	defer db.Close()

	// Подключение к Redis
	redisCache := cache.NewRedisCache(
		cfg.Redis.Addr,
		cfg.Redis.Password,
		cfg.Redis.DB,
		cfg.Redis.TTL,
		logger,
	)
	defer redisCache.Close()

	// Проверка подключения к Redis
	ctx := context.Background()
	if err := redisCache.Ping(ctx); err != nil {
		logger.Warn("Redis is not available, caching will be disabled", zap.Error(err))
	} else {
		logger.Info("Connected to Redis successfully")
	}

	// Создание repository
	analyticsRepo := repository.NewPostgresAnalyticsRepository(db, logger)

	// Создание service
	analyticsService := services.NewAnalyticsService(analyticsRepo, redisCache, logger)

	// Создание gRPC сервера
	grpcServer := grpc.NewServer()

	// Регистрация сервисов
	grpcHandler.RegisterAnalyticsServer(grpcServer, analyticsService, logger)

	// Регистрация health check
	healthServer := health.NewServer()
	grpc_health_v1.RegisterHealthServer(grpcServer, healthServer)
	healthServer.SetServingStatus("", grpc_health_v1.HealthCheckResponse_SERVING)

	// Включение reflection для отладки
	reflection.Register(grpcServer)

	// Запуск gRPC сервера
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.Server.GRPCPort))
	if err != nil {
		logger.Fatal("Failed to listen", zap.Error(err))
	}

	// Graceful shutdown
	go func() {
		logger.Info("Analytics gRPC server is running", zap.String("port", cfg.Server.GRPCPort))
		if err := grpcServer.Serve(listener); err != nil {
			logger.Fatal("Failed to serve", zap.Error(err))
		}
	}()

	// Ожидание сигнала остановки
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down analytics server...")
	grpcServer.GracefulStop()
	logger.Info("Analytics server stopped")
}

// initDatabase инициализирует подключение к базе данных
func initDatabase(cfg *config.DatabaseConfig, logger *zap.Logger) (*sql.DB, error) {
	logger.Info("Connecting to database",
		zap.String("host", cfg.Host),
		zap.String("port", cfg.Port),
		zap.String("database", cfg.DBName),
	)

	db, err := sql.Open("postgres", cfg.GetDSN())
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	// Настройка пула соединений
	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetConnMaxLifetime(cfg.ConnMaxLifetime)

	// Проверка подключения
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	logger.Info("Database connection established successfully",
		zap.Int("max_open_conns", cfg.MaxOpenConns),
		zap.Int("max_idle_conns", cfg.MaxIdleConns),
	)

	return db, nil
}
