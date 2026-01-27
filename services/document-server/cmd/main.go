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
	"github.com/rusgainew/kkm-project-mks/document-server/internal/application/document"
	"github.com/rusgainew/kkm-project-mks/document-server/internal/domain/ports"
	"github.com/rusgainew/kkm-project-mks/document-server/internal/infrastructure/config"
	"github.com/rusgainew/kkm-project-mks/document-server/internal/infrastructure/messaging"
	"github.com/rusgainew/kkm-project-mks/document-server/internal/infrastructure/middleware"
	"github.com/rusgainew/kkm-project-mks/document-server/internal/infrastructure/observability"
	"github.com/rusgainew/kkm-project-mks/document-server/internal/infrastructure/repository"
	grpchandler "github.com/rusgainew/kkm-project-mks/document-server/internal/interfaces/grpc"
	"github.com/rusgainew/kkm-project-mks/document-server/migration"
	pb "github.com/rusgainew/kkm-project-mks/proto-lib/document"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
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

	// Инициализация трассирования
	tracing := observability.InitTracingProvider(cfg.Observability.ServiceName)
	logger.Info("Distributed tracing initialized", zap.String("service", cfg.Observability.ServiceName))

	logger.Info("Starting document-server",
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
	migrationDSN := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.Database.User, cfg.Database.Password, cfg.Database.Host,
		cfg.Database.Port, cfg.Database.DBName, cfg.Database.SSLMode)
	migrator := migration.NewMigrator(db.DB, migrationDSN, logger)
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
	docRepo := repository.NewPostgresDocumentRepository(db)

	// Инициализация сервиса с трассированием
	docService := document.NewServiceWithTracing(docRepo, eventPublisher, logger, tracing)

	// Инициализация metrics
	metrics := observability.InitMetrics()

	// Инициализация gRPC handlers с трассированием
	docHandler := grpchandler.NewDocumentHandlerWithTracing(docService, logger, metrics)
	healthHandler := grpchandler.NewHealthHandler()

	// Инициализация JWT middleware (секрет из окружения, required)
	jwtSecret := os.Getenv("DOCUMENT_SERVER_JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = os.Getenv("JWT_SECRET")
	}
	if jwtSecret == "" {
		logger.Fatal("Переменная DOCUMENT_SERVER_JWT_SECRET или JWT_SECRET должна быть установлена")
	}
	authMiddleware := middleware.NewAuthMiddleware(jwtSecret, logger)

	// Инициализация rate limiting (1000 requests per second)
	rateLimiter := middleware.NewRateLimiter(1000, logger)

	// Инициализация timeout interceptor
	timeoutInterceptor := middleware.NewTimeoutInterceptor(logger)

	// Запуск gRPC сервера с цепочкой интерсепторов
	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			authMiddleware.UnaryServerInterceptor(),
			timeoutInterceptor.UnaryServerInterceptor(),
			rateLimiter.UnaryServerInterceptor(),
		),
	)
	pb.RegisterDocumentServiceServer(grpcServer, docHandler)
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
	grpcServer.GracefulStop()
	logger.Info("Server stopped")
}

func initDatabase(cfg config.DatabaseConfig, logger *zap.Logger) (*sqlx.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, cfg.SSLMode)

	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		logger.Error("Failed to connect to database", zap.Error(err))
		return nil, err
	}

	// Configure connection pool
	db.SetMaxOpenConns(25)                 // Maximum 25 open connections
	db.SetMaxIdleConns(5)                  // Keep 5 idle connections
	db.SetConnMaxLifetime(5 * time.Minute) // Reconnect every 5 minutes
	db.SetConnMaxIdleTime(2 * time.Minute) // Close idle connections after 2 minutes

	// Проверяем подключение
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		logger.Error("Failed to ping database", zap.Error(err))
		return nil, err
	}

	logger.Info("Database connected successfully",
		zap.String("host", cfg.Host),
		zap.String("port", cfg.Port),
		zap.String("dbname", cfg.DBName),
	)

	return db, nil
}

func initLogger(logLevel string) (*zap.Logger, error) {
	atom := zap.NewAtomicLevel()

	switch logLevel {
	case "debug":
		atom.SetLevel(zapcore.DebugLevel)
	case "info":
		atom.SetLevel(zapcore.InfoLevel)
	case "warn":
		atom.SetLevel(zapcore.WarnLevel)
	case "error":
		atom.SetLevel(zapcore.ErrorLevel)
	default:
		atom.SetLevel(zapcore.InfoLevel)
	}

	config := zap.NewProductionConfig()
	config.Level = atom

	logger, err := config.Build()
	if err != nil {
		return nil, err
	}

	return logger, nil
}

func startMetricsServer(port string, logger *zap.Logger) {
	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	addr := fmt.Sprintf(":%s", port)
	logger.Info("Metrics server started", zap.String("port", port))

	server := &http.Server{
		Addr:              addr,
		Handler:           nil,
		ReadHeaderTimeout: 10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	if err := server.ListenAndServe(); err != nil {
		logger.Error("Failed to start metrics server", zap.Error(err))
	}
}
