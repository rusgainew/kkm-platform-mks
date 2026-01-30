// Файл api-gateway/pkg/di_container/factories.go содержит реализацию пакета di_container.
package di_container

import (
	"github.com/rusgainew/kkm-project-mks/api-gateway/internal/application/services"
	"github.com/rusgainew/kkm-project-mks/api-gateway/internal/infrastructure/auth"
	"github.com/rusgainew/kkm-project-mks/api-gateway/internal/infrastructure/cache"
	"github.com/rusgainew/kkm-project-mks/api-gateway/internal/infrastructure/client"
	"github.com/rusgainew/kkm-project-mks/api-gateway/internal/infrastructure/config"
	"github.com/rusgainew/kkm-project-mks/api-gateway/internal/infrastructure/observability"
	"go.uber.org/zap"
)

// createConnectionManager создает connection manager для gRPC клиентов
func createConnectionManager(cfg *config.Config, logger *zap.Logger) *client.ConnectionManager {
	// Если TLS включен, используем NewConnectionManagerWithTLS
	if cfg.TLS.Enabled {
		tlsConfig := &client.GRPCTLSConfig{
			Enabled:            cfg.TLS.Enabled,
			CertFile:           cfg.TLS.CertFile,
			KeyFile:            cfg.TLS.KeyFile,
			CAFile:             cfg.TLS.CAFile,
			InsecureSkipVerify: cfg.TLS.InsecureSkipVerify,
		}

		// Валидация TLS конфигурации
		if err := tlsConfig.Validate(); err != nil {
			logger.Error("Invalid TLS configuration", zap.Error(err))
			// Продолжаем с insecure соединением как fallback
			return client.NewConnectionManager(cfg.Timeouts.GRPCDial, logger)
		}

		logger.Info("TLS enabled for gRPC connections")
		return client.NewConnectionManagerWithTLS(cfg.Timeouts.GRPCDial, logger, tlsConfig)
	}

	// Используем insecure соединение (для разработки)
	logger.Warn("TLS disabled for gRPC connections - use for development only")
	return client.NewConnectionManager(cfg.Timeouts.GRPCDial, logger)
}

// createAuthService создает auth сервис
func createAuthService(cfg *config.Config, logger *zap.Logger) *auth.Service {
	return auth.NewService(cfg.JWT.Secret, cfg.JWT.Expiration, logger)
}

// createCompanyService создает company сервис
func createCompanyService(
	cfg *config.Config,
	connManager *client.ConnectionManager,
	metrics *observability.Metrics,
	tracer *observability.Tracer,
	logger *zap.Logger,
	redisCache *cache.RedisCache,
) *services.CompanyService {
	return services.NewCompanyService(
		connManager,
		cfg.Services.CompanyService,
		metrics,
		tracer,
		logger,
		redisCache,
	)
}

// createInvoiceService создает invoice сервис
func createInvoiceService(
	cfg *config.Config,
	connManager *client.ConnectionManager,
	metrics *observability.Metrics,
	tracer *observability.Tracer,
	logger *zap.Logger,
) *services.InvoiceService {
	return services.NewInvoiceService(
		connManager,
		cfg.Services.InvoiceService,
		metrics,
		tracer,
		logger,
	)
}

// createCatalogService создает catalog сервис
func createCatalogService(
	cfg *config.Config,
	connManager *client.ConnectionManager,
	metrics *observability.Metrics,
	tracer *observability.Tracer,
	logger *zap.Logger,
) *services.CatalogService {
	return services.NewCatalogService(
		connManager,
		cfg.Services.CatalogService,
		metrics,
		tracer,
		logger,
	)
}

// createBankAccountService создает bank account сервис
func createBankAccountService(
	cfg *config.Config,
	connManager *client.ConnectionManager,
	metrics *observability.Metrics,
	tracer *observability.Tracer,
	logger *zap.Logger,
) *services.BankAccountService {
	return services.NewBankAccountService(
		connManager,
		cfg.Services.BankAccountService,
		metrics,
		tracer,
		logger,
	)
}

// createUserService создает user сервис
func createUserService(
	cfg *config.Config,
	connManager *client.ConnectionManager,
	metrics *observability.Metrics,
	tracer *observability.Tracer,
	logger *zap.Logger,
) *services.UserService {
	return services.NewUserService(
		connManager,
		cfg.Services.UserService,
		metrics,
		tracer,
		logger,
	)
}

// createForeignCompanyService создает foreign company сервис
func createForeignCompanyService(
	cfg *config.Config,
	connManager *client.ConnectionManager,
	metrics *observability.Metrics,
	tracer *observability.Tracer,
	logger *zap.Logger,
) *services.ForeignCompanyService {
	return services.NewForeignCompanyService(
		connManager,
		cfg.Services.ForeignCompanyService,
		metrics,
		tracer,
		logger,
	)
}

// createDocumentService создает document сервис
func createDocumentService(
	cfg *config.Config,
	connManager *client.ConnectionManager,
	metrics *observability.Metrics,
	tracer *observability.Tracer,
	logger *zap.Logger,
) *services.DocumentService {
	return services.NewDocumentService(
		connManager,
		cfg.Services.DocumentService,
		metrics,
		tracer,
		logger,
	)
}

// createInvoiceQueryService создает invoice query сервис
func createInvoiceQueryService(
	cfg *config.Config,
	connManager *client.ConnectionManager,
	metrics *observability.Metrics,
	tracer *observability.Tracer,
	logger *zap.Logger,
) *services.InvoiceQueryService {
	return services.NewInvoiceQueryService(
		connManager,
		cfg.Services.InvoiceQueryService,
		metrics,
		tracer,
		logger,
	)
}

// createCatalogQueryService создает catalog query сервис
func createCatalogQueryService(
	cfg *config.Config,
	connManager *client.ConnectionManager,
	metrics *observability.Metrics,
	tracer *observability.Tracer,
	logger *zap.Logger,
) *services.CatalogQueryService {
	return services.NewCatalogQueryService(
		connManager,
		cfg.Services.CatalogQueryService,
		metrics,
		tracer,
		logger,
	)
}

// createBankAccountQueryService создает bank account query сервис
func createBankAccountQueryService(
	cfg *config.Config,
	connManager *client.ConnectionManager,
	metrics *observability.Metrics,
	tracer *observability.Tracer,
	logger *zap.Logger,
) *services.BankAccountQueryService {
	return services.NewBankAccountQueryService(
		connManager,
		cfg.Services.BankAccountQueryService,
		metrics,
		tracer,
		logger,
	)
}

// createUserQueryService создает user query сервис
func createUserQueryService(
	cfg *config.Config,
	connManager *client.ConnectionManager,
	metrics *observability.Metrics,
	tracer *observability.Tracer,
	logger *zap.Logger,
) *services.UserQueryService {
	return services.NewUserQueryService(
		connManager,
		cfg.Services.UserQueryService,
		metrics,
		tracer,
		logger,
	)
}

// createDocumentQueryService создает document query сервис
func createDocumentQueryService(
	cfg *config.Config,
	connManager *client.ConnectionManager,
	metrics *observability.Metrics,
	tracer *observability.Tracer,
	logger *zap.Logger,
) *services.DocumentQueryService {
	return services.NewDocumentQueryService(
		connManager,
		cfg.Services.DocumentQueryService,
		metrics,
		tracer,
		logger,
	)
}
