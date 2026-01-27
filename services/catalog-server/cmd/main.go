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

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"

	"github.com/rusgainew/kkm-project-mks/catalog-server/internal/application/catalog"
	"github.com/rusgainew/kkm-project-mks/catalog-server/internal/infrastructure/config"
	"github.com/rusgainew/kkm-project-mks/catalog-server/internal/infrastructure/messaging"
	"github.com/rusgainew/kkm-project-mks/catalog-server/internal/infrastructure/middleware"
	"github.com/rusgainew/kkm-project-mks/catalog-server/internal/infrastructure/migration"
	"github.com/rusgainew/kkm-project-mks/catalog-server/internal/infrastructure/observability"
	"github.com/rusgainew/kkm-project-mks/catalog-server/internal/infrastructure/repository"
	grpchandlers "github.com/rusgainew/kkm-project-mks/catalog-server/internal/interfaces/grpc"
	pb "github.com/rusgainew/kkm-project-mks/proto-lib/api"
)

func main() {
	// Загрузка конфигурации
	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("Failed to load config: %v\n", err)
		os.Exit(1)
	}

	// Настройка логирования
	logger, err := observability.NewLogger(cfg.Observability.LogLevel, cfg.Observability.Environment)
	if err != nil {
		fmt.Printf("Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer logger.Sync()

	logger.Info("starting catalog server",
		zap.String("version", "1.0.0"),
		zap.String("environment", cfg.Observability.Environment),
	)

	// Подключение к базе данных
	db, err := sqlx.Connect(cfg.Database.Driver, cfg.Database.URL)
	if err != nil {
		logger.Fatal("failed to connect to database", zap.Error(err))
	}
	defer db.Close()

	// Настройка пула соединений
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	logger.Info("database connected")

	// Выполнение автоматических миграций
	migrator := migration.NewMigrator(db.DB, logger)
	if err := migrator.Run(); err != nil {
		logger.Fatal("failed to run database migrations", zap.Error(err))
	}

	// Инициализация RabbitMQ publisher
	publisher, err := messaging.NewRabbitMQPublisher(cfg.RabbitMQ.URL, logger)
	if err != nil {
		logger.Fatal("failed to initialize RabbitMQ publisher", zap.Error(err))
	}
	defer publisher.Close()

	// Создание репозитория и сервиса
	catalogRepo := repository.NewPostgresCatalogRepository(db)
	catalogService := catalog.NewService(catalogRepo, publisher, logger)

	// Создание gRPC обработчиков
	catalogHandler := grpchandlers.NewCatalogHandler(catalogService, logger)
	healthHandler := grpchandlers.NewHealthHandler(logger)

	// Создание JWT authentication middleware
	authMiddleware := middleware.NewAuthMiddleware(cfg.Auth.JWTSecret, logger)

	// Запуск Prometheus metrics endpoint
	go func() {
		metricsAddr := fmt.Sprintf(":%s", cfg.Server.MetricsPort)
		logger.Info("starting metrics server", zap.String("address", metricsAddr))
		http.Handle("/metrics", promhttp.Handler())
		server := &http.Server{
			Addr:              metricsAddr,
			Handler:           nil,
			ReadHeaderTimeout: 10 * time.Second,
			WriteTimeout:      10 * time.Second,
			IdleTimeout:       60 * time.Second,
		}
		if err := server.ListenAndServe(); err != nil {
			logger.Error("metrics server failed", zap.Error(err))
		}
	}()

	// Создание gRPC сервера с authentication interceptor
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(authMiddleware.UnaryServerInterceptor()),
	)

	// Регистрация сервисов
	pb.RegisterCatalogCommandServiceServer(grpcServer, catalogHandler)
	grpc_health_v1.RegisterHealthServer(grpcServer, healthHandler)

	// Включение gRPC reflection для отладки
	reflection.Register(grpcServer)

	// Запуск gRPC сервера
	grpcAddr := fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port)
	listener, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		logger.Fatal("failed to listen", zap.Error(err), zap.String("address", grpcAddr))
	}

	logger.Info("catalog server listening", zap.String("address", grpcAddr))

	// Graceful shutdown
	go func() {
		if err := grpcServer.Serve(listener); err != nil {
			logger.Fatal("failed to serve", zap.Error(err))
		}
	}()

	// Ожидание сигнала завершения
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("shutting down catalog server...")

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_ = ctx // Используется для будущего расширения
	grpcServer.GracefulStop()

	logger.Info("catalog server stopped")
}
