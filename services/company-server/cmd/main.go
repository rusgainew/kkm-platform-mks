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

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rusgainew/kkm-project-mks/company-server/internal/application/company"
	"github.com/rusgainew/kkm-project-mks/company-server/internal/domain/ports"
	"github.com/rusgainew/kkm-project-mks/company-server/internal/infrastructure/config"
	"github.com/rusgainew/kkm-project-mks/company-server/internal/infrastructure/messaging"
	"github.com/rusgainew/kkm-project-mks/company-server/internal/infrastructure/middleware"
	"github.com/rusgainew/kkm-project-mks/company-server/internal/infrastructure/migration"
	"github.com/rusgainew/kkm-project-mks/company-server/internal/infrastructure/observability"
	"github.com/rusgainew/kkm-project-mks/company-server/internal/infrastructure/repository"
	grpchandler "github.com/rusgainew/kkm-project-mks/company-server/internal/interfaces/grpc"
	pb "github.com/rusgainew/kkm-project-mks/proto-lib/company"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

func main() {
	// Загрузка конфигурации
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Инициализация логгера
	logger, err := initLogger(cfg.Observability.LogLevel)
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer logger.Sync()

	logger.Info("Starting company-server",
		zap.String("service", cfg.Observability.ServiceName),
		zap.String("port", cfg.Server.Port),
	)

	// Инициализация OpenTelemetry трассировки
	ctx := context.Background()
	if cfg.Observability.EnableTracing {
		tp, shutdownTracer, err := observability.InitializeTracer(
			ctx,
			cfg.Observability.ServiceName,
			cfg.Observability.JaegerEndpoint,
			logger,
		)
		if err != nil {
			logger.Warn("Failed to initialize tracer", zap.Error(err))
		} else {
			defer shutdownTracer()
			logger.Info("OpenTelemetry tracer initialized successfully")
			_ = tp
		}
	}

	// Подключение к базе данных
	db, err := initDatabase(cfg.Database, logger)
	if err != nil {
		logger.Fatal("Failed to connect to database", zap.Error(err))
	}
	defer db.Close()

	// Выполнение автоматических миграций
	migrator := migration.NewMigrator(db.DB, logger)
	if err := migrator.Run(); err != nil {
		logger.Fatal("Failed to run database migrations", zap.Error(err))
	}

	// Инициализация event publisher
	var eventPublisher ports.EventPublisher
	if cfg.RabbitMQ.Enabled {
		eventPublisher, err = messaging.NewRabbitMQPublisher(
			cfg.RabbitMQ.URL,
			cfg.RabbitMQ.Exchange,
			cfg.RabbitMQ.ExchangeType,
			logger,
		)
		if err != nil {
			logger.Warn("Failed to initialize RabbitMQ publisher, using NoOp", zap.Error(err))
			eventPublisher = messaging.NewNoOpPublisher(logger)
		}
	} else {
		eventPublisher = messaging.NewNoOpPublisher(logger)
	}
	defer eventPublisher.Close()

	// Инициализация репозиториев
	orgRepo := repository.NewPostgresOrganizationRepository(db)
	empRepo := repository.NewPostgresEmployeeRepository(db)

	// Инициализация сервиса
	companyService := company.NewService(orgRepo, empRepo, eventPublisher, logger)

	// Инициализация gRPC handlers
	companyHandler := grpchandler.NewCompanyHandler(companyService, logger)
	healthHandler := grpchandler.NewHealthHandlerWithDeps(db, logger)

	// Инициализация JWT middleware (секрет из окружения, required)
	jwtSecret := os.Getenv("COMPANY_SERVER_JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = os.Getenv("JWT_SECRET")
	}
	if jwtSecret == "" {
		logger.Fatal("Переменная COMPANY_SERVER_JWT_SECRET или JWT_SECRET должна быть установлена")
	}
	authMiddleware := middleware.NewAuthMiddleware(jwtSecret, logger)

	// Запуск gRPC сервера с аутентификационным интерсептором
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(authMiddleware.UnaryServerInterceptor()),
	)
	pb.RegisterCompanyServiceServer(grpcServer, companyHandler)
	grpc_health_v1.RegisterHealthServer(grpcServer, healthHandler)
	reflection.Register(grpcServer)

	// Запуск metrics сервера
	go startMetricsServer(cfg.Server.MetricsPort, logger)

	// Запуск gRPC сервера
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.Server.Port))
	if err != nil {
		logger.Fatal("Failed to listen", zap.Error(err))
	}

	go func() {
		logger.Info("gRPC server started", zap.String("port", cfg.Server.Port))
		if err := grpcServer.Serve(lis); err != nil {
			logger.Fatal("Failed to serve", zap.Error(err))
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	// Graceful stop с таймаутом
	stopped := make(chan struct{})
	go func() {
		grpcServer.GracefulStop()
		close(stopped)
	}()

	select {
	case <-stopped:
		logger.Info("Server stopped gracefully")
	case <-time.After(30 * time.Second):
		logger.Warn("Server stop timeout, forcing shutdown")
		grpcServer.Stop()
	}
}

// initLogger инициализирует zap logger
func initLogger(level string) (*zap.Logger, error) {
	var zapConfig zap.Config

	if level == "debug" {
		zapConfig = zap.NewDevelopmentConfig()
	} else {
		zapConfig = zap.NewProductionConfig()
	}

	// Установка уровня логирования
	switch level {
	case "debug":
		zapConfig.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	case "info":
		zapConfig.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	case "warn":
		zapConfig.Level = zap.NewAtomicLevelAt(zap.WarnLevel)
	case "error":
		zapConfig.Level = zap.NewAtomicLevelAt(zap.ErrorLevel)
	default:
		zapConfig.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	}

	return zapConfig.Build()
}

// initDatabase инициализирует подключение к базе данных
func initDatabase(cfg config.DatabaseConfig, logger *zap.Logger) (*sqlx.DB, error) {
	logger.Info("Connecting to database", zap.String("host", cfg.Host), zap.String("database", cfg.DBName))

	db, err := sqlx.Connect("postgres", cfg.GetDSN())
	if err != nil {
		return nil, err
	}

	// Настройка connection pool
	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetConnMaxLifetime(cfg.ConnMaxLifetime)

	// Проверка подключения
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return nil, err
	}

	logger.Info("Database connection established")
	return db, nil
}

// startMetricsServer запускает HTTP сервер для метрик Prometheus
func startMetricsServer(port string, logger *zap.Logger) {
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	server := &http.Server{
		Addr:         fmt.Sprintf(":%s", port),
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	logger.Info("Metrics server started", zap.String("port", port))
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Error("Metrics server error", zap.Error(err))
	}
}
