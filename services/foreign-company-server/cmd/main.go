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

	"github.com/rusgainew/kkm-project-mks/foreign-company-server/internal/application/foreigncompany"
	"github.com/rusgainew/kkm-project-mks/foreign-company-server/internal/infrastructure/config"
	"github.com/rusgainew/kkm-project-mks/foreign-company-server/internal/infrastructure/messaging"
	"github.com/rusgainew/kkm-project-mks/foreign-company-server/internal/infrastructure/middleware"
	"github.com/rusgainew/kkm-project-mks/foreign-company-server/internal/infrastructure/migration"
	"github.com/rusgainew/kkm-project-mks/foreign-company-server/internal/infrastructure/observability"
	"github.com/rusgainew/kkm-project-mks/foreign-company-server/internal/infrastructure/repository"
	grpchandlers "github.com/rusgainew/kkm-project-mks/foreign-company-server/internal/interfaces/grpc"
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

	logger.Info("starting foreign company server",
		zap.String("version", "1.0.0"),
		zap.String("environment", cfg.Observability.Environment),
	)

	// Подключение к базе данных
	db, err := sqlx.Connect(cfg.Database.Driver, cfg.Database.URL)
	if err != nil {
		logger.Fatal("failed to connect to database", zap.Error(err))
	}
	defer db.Close()

	// Настройка connection pool
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	logger.Info("database connection established")

	// Выполнение автоматических миграций
	migrator := migration.NewMigrator(db.DB, logger)
	if err := migrator.Run(); err != nil {
		logger.Fatal("failed to run database migrations", zap.Error(err))
	}

	// Инициализация RabbitMQ издателя
	publisher, err := messaging.NewRabbitMQPublisher(cfg.RabbitMQ.URL, logger)
	if err != nil {
		logger.Warn("failed to initialize RabbitMQ publisher (continuing without events)", zap.Error(err))
		publisher = nil
	}
	if publisher != nil {
		defer publisher.Close()
	}

	// Инициализация репозитория
	repo := repository.NewPostgresForeignCompanyRepository(db)

	// Инициализация сервиса
	service := foreigncompany.NewService(repo, publisher, logger)

	// Создание JWT authentication middleware
	authMiddleware := middleware.NewAuthMiddleware(cfg.Auth.JWTSecret, logger)

	// Создание gRPC сервера с authentication interceptor
	grpcServer := grpc.NewServer(
		grpc.MaxConcurrentStreams(1000),
		grpc.UnaryInterceptor(authMiddleware.UnaryServerInterceptor()),
	)

	// Регистрация обработчиков
	fcHandler := grpchandlers.NewForeignCompanyHandler(service, logger)
	pb.RegisterForeignCompanyCommandServiceServer(grpcServer, fcHandler)

	// Регистрация health check
	healthHandler := grpchandlers.NewHealthHandler()
	grpc_health_v1.RegisterHealthServer(grpcServer, healthHandler)

	// Регистрация reflection (для grpcurl)
	reflection.Register(grpcServer)

	logger.Info("gRPC handlers registered")

	// Запуск metrics endpoint
	go func() {
		metricsAddr := fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.MetricsPort)
		mux := http.NewServeMux()
		mux.Handle("/metrics", promhttp.Handler())
		mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("OK"))
		})

		metricsServer := &http.Server{
			Addr:         metricsAddr,
			Handler:      mux,
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 10 * time.Second,
		}

		logger.Info("metrics server starting", zap.String("address", metricsAddr))
		if err := metricsServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("metrics server failed", zap.Error(err))
		}
	}()

	// Запуск gRPC сервера
	grpcAddr := fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port)
	listener, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		logger.Fatal("failed to create listener", zap.Error(err))
	}

	// Канал для graceful shutdown
	done := make(chan bool, 1)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Запуск сервера в отдельной горутине
	go func() {
		logger.Info("gRPC server starting", zap.String("address", grpcAddr))
		if err := grpcServer.Serve(listener); err != nil {
			logger.Error("gRPC server failed", zap.Error(err))
		}
	}()

	logger.Info("foreign company server is ready",
		zap.String("grpc_address", grpcAddr),
		zap.String("metrics_address", fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.MetricsPort)),
	)

	// Ожидание сигнала завершения
	<-quit
	logger.Info("shutting down server...")

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	_ = ctx // Используется для будущего расширения
	defer cancel()

	// Остановка gRPC сервера
	stopped := make(chan struct{})
	go func() {
		grpcServer.GracefulStop()
		close(stopped)
	}()

	select {
	case <-stopped:
		logger.Info("gRPC server stopped gracefully")
	case <-ctx.Done():
		logger.Warn("forcing gRPC server shutdown")
		grpcServer.Stop()
	}

	done <- true
	logger.Info("server shutdown complete")
}
