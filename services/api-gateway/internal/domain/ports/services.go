// Файл api-gateway/internal/domain/ports/services.go содержит реализацию пакета ports.
package ports

import (
	"context"

	"github.com/rusgainew/kkm-project-mks/api-gateway/internal/domain/models"
)

// CompanyServiceInterface интерфейс для CompanyService
type CompanyServiceInterface interface {
	CreateCompany(ctx context.Context, company *models.Company) (*models.Company, error)
	GetCompany(ctx context.Context, id string) (*models.Company, error)
	UpdateCompany(ctx context.Context, company *models.Company) (*models.Company, error)
	DeleteCompany(ctx context.Context, id string) error
	ListCompanies(ctx context.Context, page, pageSize int) ([]*models.Company, int, error)
}

// InvoiceServiceInterface интерфейс для InvoiceService
type InvoiceServiceInterface interface {
	CreateInvoice(ctx context.Context, invoice *models.Invoice) (*models.Invoice, error)
	UpdateInvoice(ctx context.Context, invoice *models.Invoice) (*models.Invoice, error)
	SignInvoice(ctx context.Context, id, signatureData string) (*models.Invoice, error)
	RevokeInvoice(ctx context.Context, id, reason string) (*models.Invoice, error)
	AcceptOrRejectInvoice(ctx context.Context, id string, accept bool, reason string) (*models.Invoice, error)
}

// DocumentServiceInterface интерфейс для DocumentService
type DocumentServiceInterface interface {
	CreateDocument(ctx context.Context, doc *models.Document) (*models.Document, error)
	GetDocument(ctx context.Context, id string) (*models.Document, error)
	UpdateDocument(ctx context.Context, doc *models.Document) (*models.Document, error)
	SendDocument(ctx context.Context, id, recipientID, message string) (*models.Document, error)
	ApproveDocument(ctx context.Context, id, approvedBy, comments string) (*models.Document, error)
	RejectDocument(ctx context.Context, id, rejectedBy, reason string) (*models.Document, error)
	ArchiveDocument(ctx context.Context, id string) (*models.Document, error)
	ListDocuments(ctx context.Context, organizationID, status, createdBy string, page, perPage int32) ([]models.Document, models.PageInfo, error)
}

// CatalogServiceInterface интерфейс для CatalogService
type CatalogServiceInterface interface {
	CreateCatalog(ctx context.Context, catalog *models.Catalog) (*models.Catalog, error)
	GetCatalog(ctx context.Context, id string) (*models.Catalog, error)
	UpdateCatalog(ctx context.Context, catalog *models.Catalog) (*models.Catalog, error)
	DeleteCatalog(ctx context.Context, id string) error
}

// BankAccountServiceInterface интерфейс для BankAccountService
type BankAccountServiceInterface interface {
	CreateBankAccount(ctx context.Context, account *models.BankAccount) (*models.BankAccount, error)
	GetBankAccount(ctx context.Context, id string) (*models.BankAccount, error)
	UpdateBankAccount(ctx context.Context, account *models.BankAccount) (*models.BankAccount, error)
	DeleteBankAccount(ctx context.Context, id string) error
}
