package di_container

import (
	"context"
	"fmt"
	"time"

	"github.com/rusgainew/kkm-project-mks/api-gateway/internal/application/services"
	"github.com/rusgainew/kkm-project-mks/api-gateway/internal/infrastructure/cache"
	"github.com/rusgainew/kkm-project-mks/api-gateway/internal/infrastructure/client"
	"github.com/rusgainew/kkm-project-mks/api-gateway/internal/infrastructure/config"
	"github.com/rusgainew/kkm-project-mks/api-gateway/internal/infrastructure/observability"
	"go.uber.org/zap"
)

// Initializer отвечает за инициализацию всех компонентов
type Initializer struct {
	container *Container
}

// NewInitializer создает новый инициализатор
func NewInitializer(container *Container) *Initializer {
	return &Initializer{
		container: container,
	}
}

// InitializeConfig инициализирует конфигурацию
func (i *Initializer) InitializeConfig() error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}
	i.container.SetConfig(cfg)
	return nil
}

// InitializeLogger инициализирует логгер
func (i *Initializer) InitializeLogger() error {
	cfg := i.container.Config()
	logger, err := createLogger(cfg.Observability.LogLevel)
	if err != nil {
		return err
	}
	i.container.SetLogger(logger)
	return nil
}

// InitializeObservability инициализирует наблюдаемость (метрики, трейсинг)
func (i *Initializer) InitializeObservability(ctx context.Context) error {
	cfg := i.container.Config()
	logger := i.container.Logger()

	// Инициализация метрик
	metrics := observability.NewMetrics()
	i.container.SetMetrics(metrics)

	// Инициализация трейсинга
	if cfg.Observability.EnableTracing {
		_, shutdownTracer, err := observability.InitTracer(
			ctx,
			cfg.Observability.ServiceName,
			cfg.Observability.JaegerEndpoint,
			logger,
		)
		if err != nil {
			logger.Warn("Failed to initialize tracer", zap.Error(err))
		} else {
			i.container.RegisterCleanup(func() error {
				shutdownTracer()
				return nil
			})
		}
	}

	tracer := observability.NewTracer(cfg.Observability.ServiceName, logger)
	i.container.SetTracer(tracer)

	return nil
}

// InitializeInfrastructure инициализирует инфраструктурные компоненты
func (i *Initializer) InitializeInfrastructure(ctx context.Context) error {
	cfg := i.container.Config()
	logger := i.container.Logger()

	// Инициализация Connection Manager
	connManager := createConnectionManager(cfg, logger)
	i.container.SetConnectionManager(connManager)

	// Инициализация Auth Service
	authService := createAuthService(cfg, logger)
	i.container.SetAuthService(authService)

	// Инициализация Redis Cache (опционально, если доступен)
	if cfg.Redis != nil && cfg.Redis.Addr != "" {
		redisCache, err := cache.NewRedisCache(
			ctx,
			cfg.Redis.Addr,
			cfg.Redis.Password,
			time.Duration(cfg.Redis.TTL)*time.Second,
			i.container.Metrics(),
		)
		if err != nil {
			logger.Warn("Failed to initialize Redis cache, proceeding without caching", zap.Error(err))
		} else {
			i.container.SetRedisCache(redisCache)
			i.container.RegisterCleanup(func() error {
				return redisCache.Close()
			})
			logger.Info("Redis cache initialized successfully")
		}
	}

	return nil
}

// InitializeApplicationServices инициализирует бизнес-логику сервисов
func (i *Initializer) InitializeApplicationServices() error {
	cfg := i.container.Config()
	logger := i.container.Logger()
	tracer := i.container.Tracer()
	metrics := i.container.Metrics()
	connManager := i.container.ConnectionManager()
	redisCache := i.container.RedisCache()

	// Инициализация Application Services
	companyService := createCompanyService(cfg, connManager, metrics, tracer, logger, redisCache)
	i.container.SetCompanyService(companyService)

	invoiceService := createInvoiceService(cfg, connManager, metrics, tracer, logger)
	i.container.SetInvoiceService(invoiceService)

	catalogService := createCatalogService(cfg, connManager, metrics, tracer, logger)
	i.container.SetCatalogService(catalogService)

	bankAccountService := createBankAccountService(cfg, connManager, metrics, tracer, logger)
	i.container.SetBankAccountService(bankAccountService)

	userService := createUserService(cfg, connManager, metrics, tracer, logger)
	i.container.SetUserService(userService)

	foreignCompanyService := createForeignCompanyService(cfg, connManager, metrics, tracer, logger)
	i.container.SetForeignCompanyService(foreignCompanyService)

	documentService := createDocumentService(cfg, connManager, metrics, tracer, logger)
	i.container.SetDocumentService(documentService)

	// Инициализация Query Services
	invoiceQueryService := createInvoiceQueryService(cfg, connManager, metrics, tracer, logger)
	i.container.SetInvoiceQueryService(invoiceQueryService)

	catalogQueryService := createCatalogQueryService(cfg, connManager, metrics, tracer, logger)
	i.container.SetCatalogQueryService(catalogQueryService)

	bankAccountQueryService := createBankAccountQueryService(cfg, connManager, metrics, tracer, logger)
	i.container.SetBankAccountQueryService(bankAccountQueryService)

	userQueryService := createUserQueryService(cfg, connManager, metrics, tracer, logger)
	i.container.SetUserQueryService(userQueryService)

	documentQueryService := createDocumentQueryService(cfg, connManager, metrics, tracer, logger)
	i.container.SetDocumentQueryService(documentQueryService)

	companyQueryService := createCompanyQueryService(cfg, connManager, metrics, tracer, logger)
	i.container.SetCompanyQueryService(companyQueryService)

	foreignCompanyQueryService := createForeignCompanyQueryService(cfg, connManager, metrics, tracer, logger)
	i.container.SetForeignCompanyQueryService(foreignCompanyQueryService)

	// Инициализация Analytics Service
	analyticsService := createAnalyticsService(cfg, redisCache, metrics, logger)
	i.container.SetAnalyticsService(analyticsService)

	return nil
}

// InitializeAll инициализирует все компоненты в правильном порядке
func (i *Initializer) InitializeAll(ctx context.Context) error {
	steps := []struct {
		name string
		fn   func() error
	}{
		{"config", i.InitializeConfig},
		{"logger", i.InitializeLogger},
		{"observability", func() error { return i.InitializeObservability(ctx) }},
		{"infrastructure", func() error { return i.InitializeInfrastructure(ctx) }},
		{"application services", i.InitializeApplicationServices},
	}

	for _, step := range steps {
		if err := step.fn(); err != nil {
			logger := i.container.Logger()
			if logger != nil {
				logger.Error("Failed to initialize "+step.name, zap.Error(err))
			} else {
				fmt.Printf("Failed to initialize %s: %v\n", step.name, err)
			}
			return err
		}
		logger := i.container.Logger()
		if logger != nil {
			logger.Info("Initialized " + step.name)
		}
	}

	return nil
}

// Вспомогательные функции для создания компонентов

func createLogger(level string) (*zap.Logger, error) {
	var cfg zap.Config

	if level == "debug" {
		cfg = zap.NewDevelopmentConfig()
	} else {
		cfg = zap.NewProductionConfig()
	}

	switch level {
	case "debug":
		cfg.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	case "info":
		cfg.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	case "warn":
		cfg.Level = zap.NewAtomicLevelAt(zap.WarnLevel)
	case "error":
		cfg.Level = zap.NewAtomicLevelAt(zap.ErrorLevel)
	default:
		cfg.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	}

	return cfg.Build()
}
func createCompanyQueryService(
	cfg *config.Config,
	connManager *client.ConnectionManager,
	metrics *observability.Metrics,
	tracer *observability.Tracer,
	logger *zap.Logger,
) *services.CompanyQueryService {
	return services.NewCompanyQueryService(
		connManager,
		cfg.Services.CompanyService,
		metrics,
		tracer,
		logger,
	)
}

func createForeignCompanyQueryService(
	cfg *config.Config,
	connManager *client.ConnectionManager,
	metrics *observability.Metrics,
	tracer *observability.Tracer,
	logger *zap.Logger,
) *services.ForeignCompanyQueryService {
	return services.NewForeignCompanyQueryService(
		connManager,
		cfg.Services.ForeignCompanyService,
		metrics,
		tracer,
		logger,
	)
}

func createAnalyticsService(
	cfg *config.Config,
	redisCache *cache.RedisCache,
	metrics *observability.Metrics,
	logger *zap.Logger,
) *services.AnalyticsService {
	// TODO: Подключить PostgreSQL для аналитики
	// Пока создаем сервис с nil repository - это вызовет панику при запросе
	// Необходимо:
	// 1. Добавить DB connection в Infrastructure
	// 2. Создать PostgresAnalyticsRepository
	// 3. Применить миграции (001_create_analytics_views.up.sql)
	// 4. Обернуть в metrics wrapper
	//
	// Пример полной инициализации:
	// db := getDBConnection(cfg.Database.InvoiceDB) // Нужно добавить DB connection
	// repo := repository.NewPostgresAnalyticsRepository(db, logger)
	// repoWithMetrics := repository.NewAnalyticsRepositoryWithMetrics(repo, metrics, logger)
	// return services.NewAnalyticsService(repoWithMetrics, redisCache, logger)

	logger.Warn("AnalyticsService created WITHOUT repository - analytics endpoints will fail until DB connection is configured")
	return services.NewAnalyticsService(nil, redisCache, logger)
}
