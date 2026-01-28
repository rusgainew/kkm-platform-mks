package di_container

import (
	"sync"

	"github.com/rusgainew/kkm-project-mks/api-gateway/internal/application/services"
	"github.com/rusgainew/kkm-project-mks/api-gateway/internal/infrastructure/auth"
	"github.com/rusgainew/kkm-project-mks/api-gateway/internal/infrastructure/cache"
	"github.com/rusgainew/kkm-project-mks/api-gateway/internal/infrastructure/client"
	"github.com/rusgainew/kkm-project-mks/api-gateway/internal/infrastructure/config"

	"github.com/rusgainew/kkm-project-mks/api-gateway/internal/infrastructure/observability"
	"go.uber.org/zap"
)

// Container - DI контейнер для управления зависимостями
type Container struct {
	mu sync.RWMutex

	// Config
	config *config.Config

	// Observability
	logger  *zap.Logger
	tracer  *observability.Tracer
	metrics *observability.Metrics

	// Authentication
	authService *auth.Service

	// Infrastructure
	connManager *client.ConnectionManager
	redisCache  *cache.RedisCache

	// Application Services
	companyService        *services.CompanyService
	invoiceService        *services.InvoiceService
	catalogService        *services.CatalogService
	bankAccountService    *services.BankAccountService
	userService           *services.UserService
	foreignCompanyService *services.ForeignCompanyService
	documentService       *services.DocumentService

	// Query Services (read-side)
	invoiceQueryService        *services.InvoiceQueryService
	catalogQueryService        *services.CatalogQueryService
	bankAccountQueryService    *services.BankAccountQueryService
	userQueryService           *services.UserQueryService
	documentQueryService       *services.DocumentQueryService
	companyQueryService        *services.CompanyQueryService
	foreignCompanyQueryService *services.ForeignCompanyQueryService
	analyticsClient            *client.AnalyticsClient

	// Cleanup functions
	cleanupFuncs []func() error
}

// NewContainer создает новый DI контейнер
func NewContainer() *Container {
	return &Container{
		cleanupFuncs: make([]func() error, 0),
	}
}

// SetConfig устанавливает конфигурацию
func (c *Container) SetConfig(cfg *config.Config) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.config = cfg
}

// Config возвращает конфигурацию
func (c *Container) Config() *config.Config {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.config
}

// SetLogger устанавливает логгер
func (c *Container) SetLogger(logger *zap.Logger) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.logger = logger
}

// Logger возвращает логгер
func (c *Container) Logger() *zap.Logger {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.logger
}

// SetTracer устанавливает tracer
func (c *Container) SetTracer(tracer *observability.Tracer) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.tracer = tracer
}

// Tracer возвращает tracer
func (c *Container) Tracer() *observability.Tracer {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.tracer
}

// SetMetrics устанавливает metrics
func (c *Container) SetMetrics(metrics *observability.Metrics) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.metrics = metrics
}

// Metrics возвращает metrics
func (c *Container) Metrics() *observability.Metrics {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.metrics
}

// SetAuthService устанавливает auth сервис
func (c *Container) SetAuthService(service *auth.Service) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.authService = service
}

// AuthService возвращает auth сервис
func (c *Container) AuthService() *auth.Service {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.authService
}

// SetConnectionManager устанавливает connection manager
func (c *Container) SetConnectionManager(cm *client.ConnectionManager) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.connManager = cm
	// Регистрируем cleanup функцию
	c.cleanupFuncs = append(c.cleanupFuncs, func() error {
		cm.CloseAll()
		return nil
	})
}

// ConnectionManager возвращает connection manager
func (c *Container) ConnectionManager() *client.ConnectionManager {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.connManager
}

// SetCompanyService устанавливает company сервис
func (c *Container) SetCompanyService(service *services.CompanyService) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.companyService = service
}

// CompanyService возвращает company сервис
func (c *Container) CompanyService() *services.CompanyService {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.companyService
}

// SetInvoiceService устанавливает invoice сервис
func (c *Container) SetInvoiceService(service *services.InvoiceService) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.invoiceService = service
}

// InvoiceService возвращает invoice сервис
func (c *Container) InvoiceService() *services.InvoiceService {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.invoiceService
}

// SetCatalogService устанавливает catalog сервис
func (c *Container) SetCatalogService(service *services.CatalogService) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.catalogService = service
}

// CatalogService возвращает catalog сервис
func (c *Container) CatalogService() *services.CatalogService {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.catalogService
}

// SetBankAccountService устанавливает bank account сервис
func (c *Container) SetBankAccountService(service *services.BankAccountService) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.bankAccountService = service
}

// BankAccountService возвращает bank account сервис
func (c *Container) BankAccountService() *services.BankAccountService {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.bankAccountService
}

// SetUserService устанавливает user сервис
func (c *Container) SetUserService(service *services.UserService) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.userService = service
}

// UserService возвращает user сервис
func (c *Container) UserService() *services.UserService {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.userService
}

// SetForeignCompanyService устанавливает foreign company сервис
func (c *Container) SetForeignCompanyService(service *services.ForeignCompanyService) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.foreignCompanyService = service
}

// ForeignCompanyService возвращает foreign company сервис
func (c *Container) ForeignCompanyService() *services.ForeignCompanyService {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.foreignCompanyService
}

// SetDocumentService устанавливает document сервис
func (c *Container) SetDocumentService(service *services.DocumentService) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.documentService = service
}

// DocumentService возвращает document сервис
func (c *Container) DocumentService() *services.DocumentService {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.documentService
}

// SetInvoiceQueryService устанавливает invoice query сервис
func (c *Container) SetInvoiceQueryService(service *services.InvoiceQueryService) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.invoiceQueryService = service
}

// InvoiceQueryService возвращает invoice query сервис
func (c *Container) InvoiceQueryService() *services.InvoiceQueryService {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.invoiceQueryService
}

// SetCatalogQueryService устанавливает catalog query сервис
func (c *Container) SetCatalogQueryService(service *services.CatalogQueryService) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.catalogQueryService = service
}

// CatalogQueryService возвращает catalog query сервис
func (c *Container) CatalogQueryService() *services.CatalogQueryService {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.catalogQueryService
}

// SetBankAccountQueryService устанавливает bank account query сервис
func (c *Container) SetBankAccountQueryService(service *services.BankAccountQueryService) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.bankAccountQueryService = service
}

// BankAccountQueryService возвращает bank account query сервис
func (c *Container) BankAccountQueryService() *services.BankAccountQueryService {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.bankAccountQueryService
}

// SetUserQueryService устанавливает user query сервис
func (c *Container) SetUserQueryService(service *services.UserQueryService) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.userQueryService = service
}

// UserQueryService возвращает user query сервис
func (c *Container) UserQueryService() *services.UserQueryService {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.userQueryService
}

// SetDocumentQueryService устанавливает document query сервис
func (c *Container) SetDocumentQueryService(service *services.DocumentQueryService) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.documentQueryService = service
}

// DocumentQueryService возвращает document query сервис
func (c *Container) DocumentQueryService() *services.DocumentQueryService {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.documentQueryService
}

// SetCompanyQueryService устанавливает company query сервис
func (c *Container) SetCompanyQueryService(service *services.CompanyQueryService) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.companyQueryService = service
}

// CompanyQueryService возвращает company query сервис
func (c *Container) CompanyQueryService() *services.CompanyQueryService {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.companyQueryService
}

// SetForeignCompanyQueryService устанавливает foreign company query сервис
func (c *Container) SetForeignCompanyQueryService(service *services.ForeignCompanyQueryService) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.foreignCompanyQueryService = service
}

// ForeignCompanyQueryService возвращает foreign company query сервис
func (c *Container) ForeignCompanyQueryService() *services.ForeignCompanyQueryService {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.foreignCompanyQueryService
}

// SetAnalyticsClient устанавливает analytics gRPC клиент
func (c *Container) SetAnalyticsClient(cl *client.AnalyticsClient) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.analyticsClient = cl
}

// AnalyticsClient возвращает analytics gRPC клиент
func (c *Container) AnalyticsClient() *client.AnalyticsClient {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.analyticsClient
}

// SetRedisCache устанавливает Redis кеш
func (c *Container) SetRedisCache(cache *cache.RedisCache) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.redisCache = cache
}

// RedisCache возвращает Redis кеш
func (c *Container) RedisCache() *cache.RedisCache {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.redisCache
}

// RegisterCleanup регистрирует cleanup функцию
func (c *Container) RegisterCleanup(fn func() error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.cleanupFuncs = append(c.cleanupFuncs, fn)
}

// Cleanup вызывает все cleanup функции
func (c *Container) Cleanup() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Вызываем cleanup функции в обратном порядке
	for i := len(c.cleanupFuncs) - 1; i >= 0; i-- {
		if err := c.cleanupFuncs[i](); err != nil {
			c.logger.Error("Cleanup error", zap.Error(err))
		}
	}

	return nil
}
