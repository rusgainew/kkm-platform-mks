package http

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rusgainew/kkm-project-mks/api-gateway/internal/infrastructure/config"
	"github.com/rusgainew/kkm-project-mks/api-gateway/internal/infrastructure/middleware"
	"github.com/rusgainew/kkm-project-mks/api-gateway/pkg/di_container"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"
)

// RouteConfigurator отвечает за конфигурацию маршрутов
type RouteConfigurator struct {
	container       *di_container.Container
	logger          *zap.Logger
	cfg             *config.Config
	authRateLimiter *middleware.IPBasedRateLimiter
}

// NewRouteConfigurator создает новый конфигуратор маршрутов
func NewRouteConfigurator(
	container *di_container.Container,
	logger *zap.Logger,
	cfg *config.Config,
) *RouteConfigurator {
	// Rate limiter для auth endpoints: 5 запросов в секунду, burst 10
	authRateLimiter := middleware.NewIPBasedRateLimiter(5, 10, logger)
	// Запуск периодической очистки каждые 10 минут
	authRateLimiter.Cleanup(10 * 60 * 1000000000) // 10 minutes in nanoseconds

	return &RouteConfigurator{
		container:       container,
		logger:          logger,
		cfg:             cfg,
		authRateLimiter: authRateLimiter,
	}
}

// Configure конфигурирует маршруты API
func (rc *RouteConfigurator) Configure(router *gin.Engine) error {
	// Глобальные middleware
	router.Use(middleware.SecurityHeadersMiddleware())
	router.Use(middleware.RequestSizeLimitMiddleware(10 * 1024 * 1024)) // 10MB max

	// Swagger 2.0 documentation (legacy)
	router.GET("/api/v1/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// OpenAPI 3.0 documentation
	router.StaticFile("/openapi.yaml", "./docs/openapi.yaml")
	router.StaticFile("/openapi.json", "./docs/openapi.json")
	router.GET("/api/v1/openapi.yaml", func(c *gin.Context) {
		c.File("./docs/openapi.yaml")
	})
	router.GET("/api/v1/openapi.json", func(c *gin.Context) {
		c.File("./docs/openapi.json")
	})

	// Prometheus metrics (без аутентификации)
	rc.configureMetricsRoute(router)

	api := router.Group("/api/v1")

	// Health checks (без аутентификации)
	rc.configureHealthRoutes(api)

	// User routes (регистрация и логин без аутентификации, но с rate limiting)
	rc.configureUserRoutes(api)

	// Protected routes (с аутентификацией)
	protected := api.Group("")
	protected.Use(middleware.AuthMiddleware(rc.container.AuthService(), rc.logger))
	// Forbid DELETE for managers (unless also admin) across protected routes
	protected.Use(middleware.ForbidDeleteForManagerMiddleware(rc.logger))
	{
		rc.configureProtectedUserRoutes(protected)
		rc.configureUserQueryRoutes(protected)
		rc.configureCompanyRoutes(protected)
		rc.configureCompanyQueryRoutes(protected)
		rc.configureInvoiceRoutes(protected)
		rc.configureCatalogRoutes(protected)
		rc.configureBankAccountRoutes(protected)
		rc.configureForeignCompanyRoutes(protected)
		rc.configureForeignCompanyQueryRoutes(protected)
		rc.configureInvoiceQueryRoutes(protected)
		rc.configureCatalogQueryRoutes(protected)
		rc.configureBankAccountQueryRoutes(protected)
		rc.configureDocumentRoutes(protected)
		rc.configureDocumentQueryRoutes(protected)
	}

	return nil
}

// configureHealthRoutes конфигурирует health check маршруты
func (rc *RouteConfigurator) configureHealthRoutes(api *gin.RouterGroup) {
	healthHandler := NewHealthHandler(rc.logger, rc.container)
	api.GET("/health", healthHandler.Health)
	api.GET("/ready", healthHandler.Ready)
}

// configureCompanyRoutes конфигурирует маршруты компаний
func (rc *RouteConfigurator) configureCompanyRoutes(protected *gin.RouterGroup) {
	companyHandler := NewCompanyHandler(rc.container.CompanyService(), rc.logger)
	companies := protected.Group("/companies")
	{
		// Создание доступно всем авторизованным пользователям (они создают компании для себя)
		companies.POST("", companyHandler.CreateCompany)
		companies.DELETE("/:id", middleware.RequireAdmin(rc.logger), companyHandler.DeleteCompany)
		// Чтение и обновление доступны всем авторизованным
		companies.GET("/:id", companyHandler.GetCompany)
		companies.PUT("/:id", companyHandler.UpdateCompany)
		companies.GET("", companyHandler.ListCompanies)
	}
}

// configureInvoiceRoutes конфигурирует маршруты счетов
func (rc *RouteConfigurator) configureInvoiceRoutes(protected *gin.RouterGroup) {
	invoiceHandler := NewInvoiceHandler(rc.container.InvoiceService(), rc.logger)
	invoices := protected.Group("/invoices")
	{
		// Создание требует admin/manager
		invoices.POST("", middleware.RequireAdminOrManager(rc.logger), invoiceHandler.CreateInvoice)
		// Подпись, отзыв, принятие/отклонение требуют admin/manager
		invoices.POST("/:id/sign", middleware.RequireAdminOrManager(rc.logger), invoiceHandler.SignInvoice)
		invoices.POST("/:id/revoke", middleware.RequireAdminOrManager(rc.logger), invoiceHandler.RevokeInvoice)
		invoices.POST("/:id/accept-reject", middleware.RequireAdminOrManager(rc.logger), invoiceHandler.AcceptOrRejectInvoice)
		// Обновление доступно всем авторизованным
		invoices.PUT("/:id", invoiceHandler.UpdateInvoice)
	}
}

// configureCatalogRoutes конфигурирует маршруты каталога
func (rc *RouteConfigurator) configureCatalogRoutes(protected *gin.RouterGroup) {
	catalogHandler := NewCatalogHandler(rc.container.CatalogService(), rc.logger)
	catalog := protected.Group("/catalog")
	{
		// Создание и удаление требуют admin/manager
		catalog.POST("", middleware.RequireAdminOrManager(rc.logger), catalogHandler.CreateCatalog)
		catalog.DELETE("/:id", middleware.RequireAdmin(rc.logger), catalogHandler.DeleteCatalog)
		// Чтение и обновление доступны всем
		catalog.GET("/:id", catalogHandler.GetCatalog)
		catalog.PUT("/:id", catalogHandler.UpdateCatalog)
	}
}

// configureBankAccountRoutes конфигурирует маршруты банковских счетов
func (rc *RouteConfigurator) configureBankAccountRoutes(protected *gin.RouterGroup) {
	bankAccountHandler := NewBankAccountHandler(rc.container.BankAccountService(), rc.logger)
	bankAccounts := protected.Group("/bank-accounts")
	{
		// Создание и удаление требуют admin/manager
		bankAccounts.POST("", middleware.RequireAdminOrManager(rc.logger), bankAccountHandler.CreateBankAccount)
		bankAccounts.DELETE("/:id", middleware.RequireAdmin(rc.logger), bankAccountHandler.DeleteBankAccount)
		// Чтение и обновление доступны всем
		bankAccounts.GET("/:id", bankAccountHandler.GetBankAccount)
		bankAccounts.PUT("/:id", bankAccountHandler.UpdateBankAccount)
	}
}

// configureUserRoutes конфигурирует маршруты пользователей
func (rc *RouteConfigurator) configureUserRoutes(api *gin.RouterGroup) {
	userHandler := NewUserHandler(rc.container.UserService(), rc.logger)

	// Public routes с rate limiting для защиты от brute-force
	api.POST("/users/register", rc.authRateLimiter.Middleware(), userHandler.Register)
	api.POST("/users/login", rc.authRateLimiter.Middleware(), userHandler.Login)
	api.POST("/users/refresh", userHandler.RefreshToken)
	api.POST("/users/forgot-password", userHandler.ForgotPassword)
	api.POST("/users/reset-password", userHandler.ResetPassword)

	// Protected user routes будут добавлены в protected group
	// Используем middleware для /users/me и /users/:id
}

// configureProtectedUserRoutes конфигурирует защищенные маршруты пользователей
func (rc *RouteConfigurator) configureProtectedUserRoutes(protected *gin.RouterGroup) {
	userHandler := NewUserHandler(rc.container.UserService(), rc.logger)
	users := protected.Group("/users")
	{
		users.GET("/me", userHandler.GetMe)
		users.GET("/:id", userHandler.GetUser)
		users.PUT("/:id/role", middleware.RequireAdmin(rc.logger), userHandler.AssignRole)
		users.POST("/logout", userHandler.Logout)
		users.PUT("/profile", userHandler.UpdateProfile)
		users.PUT("/password", userHandler.ChangePassword)
		users.GET("", middleware.RequireAdmin(rc.logger), userHandler.ListUsers)
	}
}

// configureForeignCompanyRoutes конфигурирует маршруты иностранных компаний
func (rc *RouteConfigurator) configureForeignCompanyRoutes(protected *gin.RouterGroup) {
	foreignCompanyHandler := NewForeignCompanyHandler(rc.container.ForeignCompanyService(), rc.logger)
	foreignCompanies := protected.Group("/foreign-companies")
	{
		// Создание и удаление требуют admin/manager
		foreignCompanies.POST("", middleware.RequireAdminOrManager(rc.logger), foreignCompanyHandler.CreateForeignCompany)
		foreignCompanies.DELETE("/:id", middleware.RequireAdmin(rc.logger), foreignCompanyHandler.DeleteForeignCompany)
		// Чтение и обновление доступны всем
		foreignCompanies.GET("/:id", foreignCompanyHandler.GetForeignCompany)
		foreignCompanies.PUT("/:id", foreignCompanyHandler.UpdateForeignCompany)
	}
}

// configureInvoiceQueryRoutes конфигурирует маршруты для чтения счетов
func (rc *RouteConfigurator) configureInvoiceQueryRoutes(protected *gin.RouterGroup) {
	invoiceQueryHandler := NewInvoiceQueryHandler(rc.container.InvoiceQueryService(), rc.logger)
	invoiceQueries := protected.Group("/invoices-query")
	{
		invoiceQueries.GET("", invoiceQueryHandler.ListInvoices)
		invoiceQueries.GET("/search", invoiceQueryHandler.SearchInvoices)
		invoiceQueries.GET("/filter", invoiceQueryHandler.ListInvoicesWithFilter)
		invoiceQueries.GET("/by-number/:number", invoiceQueryHandler.GetInvoiceByNumber)
	}
}

// configureCatalogQueryRoutes конфигурирует маршруты для чтения каталога
func (rc *RouteConfigurator) configureCatalogQueryRoutes(protected *gin.RouterGroup) {
	catalogQueryHandler := NewCatalogQueryHandler(rc.container.CatalogQueryService(), rc.logger)
	catalogQueries := protected.Group("/catalogs-query")
	{
		catalogQueries.GET("", catalogQueryHandler.ListCatalogs)
		catalogQueries.GET("/filter", catalogQueryHandler.ListCatalogsWithFilter)
	}
}

// configureBankAccountQueryRoutes конфигурирует маршруты для чтения банковских счетов
func (rc *RouteConfigurator) configureBankAccountQueryRoutes(protected *gin.RouterGroup) {
	bankAccountQueryHandler := NewBankAccountQueryHandler(rc.container.BankAccountQueryService(), rc.logger)
	bankAccountQueries := protected.Group("/bank-accounts-query")
	{
		bankAccountQueries.GET("", bankAccountQueryHandler.ListBankAccounts)
		bankAccountQueries.GET("/filter", bankAccountQueryHandler.ListBankAccountsWithFilter)
	}
}

// configureDocumentRoutes конфигурирует маршруты документов
func (rc *RouteConfigurator) configureDocumentRoutes(protected *gin.RouterGroup) {
	documentHandler := NewDocumentHandler(rc.container.DocumentService(), rc.logger)
	docs := protected.Group("/documents")
	{
		docs.POST("", documentHandler.CreateDocument)
		docs.GET(":id", documentHandler.GetDocument)
		docs.PUT(":id", documentHandler.UpdateDocument)
		docs.POST(":id/send", documentHandler.SendDocument)
		docs.POST(":id/approve", documentHandler.ApproveDocument)
		docs.POST(":id/reject", documentHandler.RejectDocument)
		docs.POST(":id/archive", documentHandler.ArchiveDocument)
		docs.GET("", documentHandler.ListDocuments)
	}
}

// configureUserQueryRoutes конфигурирует маршруты для чтения пользователей
func (rc *RouteConfigurator) configureUserQueryRoutes(protected *gin.RouterGroup) {
	userQueryHandler := NewUserQueryHandler(rc.container.UserQueryService(), rc.logger)
	userQueries := protected.Group("/users-query")
	{
		userQueries.GET("/:id", userQueryHandler.GetUser)
		userQueries.GET("", userQueryHandler.ListUsers)
		userQueries.GET("/search", userQueryHandler.SearchUsers)
	}
}

// configureDocumentQueryRoutes конфигурирует маршруты для чтения документов
func (rc *RouteConfigurator) configureDocumentQueryRoutes(protected *gin.RouterGroup) {
	documentQueryHandler := NewDocumentQueryHandler(rc.container.DocumentQueryService(), rc.logger)
	documentQueries := protected.Group("/documents-query")
	{
		documentQueries.GET("/:id", documentQueryHandler.GetDocument)
		documentQueries.GET("", documentQueryHandler.ListDocuments)
		documentQueries.GET("/search", documentQueryHandler.SearchDocuments)
		documentQueries.GET("/pending-approval", documentQueryHandler.GetPendingApprovalDocuments)
	}
}

// configureCompanyQueryRoutes конфигурирует маршруты для чтения компаний
func (rc *RouteConfigurator) configureCompanyQueryRoutes(protected *gin.RouterGroup) {
	companyQueryHandler := NewCompanyQueryHandler(rc.container.CompanyQueryService(), rc.logger)
	companyQueries := protected.Group("/companies-query")
	{
		companyQueries.GET("", companyQueryHandler.ListCompanies)
		companyQueries.GET("/search", companyQueryHandler.SearchCompanies)
		companyQueries.GET("/filter", companyQueryHandler.FilterCompanies)
	}
}

// configureForeignCompanyQueryRoutes конфигурирует маршруты для чтения иностранных компаний
func (rc *RouteConfigurator) configureForeignCompanyQueryRoutes(protected *gin.RouterGroup) {
	foreignCompanyQueryHandler := NewForeignCompanyQueryHandler(rc.container.ForeignCompanyQueryService(), rc.logger)
	foreignCompanyQueries := protected.Group("/foreign-companies-query")
	{
		foreignCompanyQueries.GET("", foreignCompanyQueryHandler.ListForeignCompanies)
		foreignCompanyQueries.GET("/search", foreignCompanyQueryHandler.SearchForeignCompanies)
		foreignCompanyQueries.GET("/filter", foreignCompanyQueryHandler.FilterForeignCompanies)
	}
}

// configureMetricsRoute конфигурирует маршрут для Prometheus метрик
func (rc *RouteConfigurator) configureMetricsRoute(router *gin.Engine) {
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))
}
