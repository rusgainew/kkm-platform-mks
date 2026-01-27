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

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rusgainew/kkm-project-mks/invoice-server/internal/application/invoice"
	"github.com/rusgainew/kkm-project-mks/invoice-server/internal/domain/ports"
	"github.com/rusgainew/kkm-project-mks/invoice-server/internal/infrastructure/config"
	"github.com/rusgainew/kkm-project-mks/invoice-server/internal/infrastructure/messaging"
	"github.com/rusgainew/kkm-project-mks/invoice-server/internal/infrastructure/middleware"
	"github.com/rusgainew/kkm-project-mks/invoice-server/internal/infrastructure/migration"
	"github.com/rusgainew/kkm-project-mks/invoice-server/internal/infrastructure/observability"
	"github.com/rusgainew/kkm-project-mks/invoice-server/internal/infrastructure/repository"
	grpcHandlers "github.com/rusgainew/kkm-project-mks/invoice-server/internal/interfaces/grpc"
	pb "github.com/rusgainew/kkm-project-mks/proto-lib/api"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

func main() {
	// Создание логгера
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("failed to create logger: %v", err)
	}
	defer logger.Sync()

	// Загрузка конфигурации
	cfg, err := config.Load()
	if err != nil {
		logger.Fatal("Failed to load config", zap.Error(err))
	}

	logger.Info("Starting invoice-server",
		zap.String("grpc_port", cfg.Server.Port),
		zap.String("metrics_port", cfg.Server.MetricsPort),
	)

	// Подключение к базе данных
	db, err := sqlx.Connect("pgx", cfg.Database.GetDSN())
	if err != nil {
		logger.Fatal("Failed to connect to database", zap.Error(err))
	}
	defer db.Close()

	// Настройка пула соединений
	db.SetMaxOpenConns(cfg.Database.MaxOpenConns)
	db.SetMaxIdleConns(cfg.Database.MaxIdleConns)
	db.SetConnMaxLifetime(cfg.Database.ConnMaxLifetime)

	// Проверка соединения с БД
	if err := db.Ping(); err != nil {
		logger.Fatal("Failed to ping database", zap.Error(err))
	}
	logger.Info("Database connection established")

	// Выполнение автоматических миграций
	migrator := migration.NewMigrator(db.DB, logger)
	if err := migrator.Run(); err != nil {
		logger.Fatal("Failed to run database migrations", zap.Error(err))
	}

	// Создание метрик
	metrics := observability.NewMetrics()

	// Создание репозиториев
	invoiceRepo := repository.NewPostgresInvoiceRepository(db)
	detailRepo := repository.NewPostgresInvoiceDetailRepository(db)
	financialRepo := repository.NewPostgresFinancialDataRepository(db)

	// Создание event publisher
	var eventPublisher ports.EventPublisher

	if cfg.RabbitMQ.Enabled {
		rabbitPublisher, err := messaging.NewRabbitMQPublisher(
			cfg.RabbitMQ.URL,
			cfg.RabbitMQ.Exchange,
			cfg.RabbitMQ.ExchangeType,
			logger,
		)
		if err != nil {
			logger.Warn("Failed to create RabbitMQ publisher, using NoOp", zap.Error(err))
			eventPublisher = messaging.NewNoOpPublisher(logger)
		} else {
			eventPublisher = rabbitPublisher
			defer rabbitPublisher.Close()
			logger.Info("RabbitMQ publisher initialized")
		}
	} else {
		eventPublisher = messaging.NewNoOpPublisher(logger)
	}

	// Создание сервисов
	invoiceService := invoice.NewService(
		invoiceRepo,
		detailRepo,
		financialRepo,
		eventPublisher,
		logger,
	)

	// Создание JWT authentication middleware
	authMiddleware := middleware.NewAuthMiddleware(cfg.Auth.JWTSecret, logger)

	// Создание gRPC обработчиков
	healthHandler := grpcHandlers.NewHealthHandler()
	invoiceCommandHandler := grpcHandlers.NewInvoiceCommandHandler(invoiceService, metrics, logger)
	invoiceQueryHandler := grpcHandlers.NewInvoiceQueryHandler(invoiceService, metrics, logger)

	// Создание gRPC сервера с authentication interceptor
	grpcServer := grpc.NewServer(
		grpc.MaxRecvMsgSize(10*1024*1024), // 10MB
		grpc.MaxSendMsgSize(10*1024*1024), // 10MB
		grpc.UnaryInterceptor(authMiddleware.UnaryServerInterceptor()),
	)

	// Регистрация сервисов
	grpc_health_v1.RegisterHealthServer(grpcServer, healthHandler)
	pb.RegisterInvoiceCommandServiceServer(grpcServer, invoiceCommandHandler)
	pb.RegisterInvoiceQueryServiceServer(grpcServer, invoiceQueryHandler)

	// Включение reflection для grpcurl
	reflection.Register(grpcServer)

	// Запуск gRPC сервера
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.Server.Port))
	if err != nil {
		logger.Fatal("Failed to listen", zap.Error(err))
	}

	go func() {
		logger.Info("gRPC server listening", zap.String("port", cfg.Server.Port))
		if err := grpcServer.Serve(listener); err != nil {
			logger.Fatal("Failed to serve gRPC", zap.Error(err))
		}
	}()

	// Запуск metrics сервера
	metricsServer := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.Server.MetricsPort),
		Handler: promhttp.Handler(),
	}

	go func() {
		logger.Info("Metrics server listening", zap.String("port", cfg.Server.MetricsPort))
		if err := metricsServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("Metrics server error", zap.Error(err))
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down servers...")

	// Остановка gRPC сервера
	grpcServer.GracefulStop()
	logger.Info("gRPC server stopped")

	// Остановка metrics сервера
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := metricsServer.Shutdown(ctx); err != nil {
		logger.Error("Metrics server forced to shutdown", zap.Error(err))
	}
	logger.Info("Metrics server stopped")

	logger.Info("invoice-server shutdown completed")
}
