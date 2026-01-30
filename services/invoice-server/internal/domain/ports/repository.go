// Файл invoice-server/internal/domain/ports/repository.go содержит реализацию пакета ports.
package ports

import (
	"context"

	"github.com/rusgainew/kkm-project-mks/invoice-server/internal/domain"
)

// InvoiceRepository определяет интерфейс для работы с хранилищем счетов-фактур
type InvoiceRepository interface {
	// Create создает новый счет-фактуру
	Create(ctx context.Context, invoice *domain.Invoice) error

	// GetByID получает счет-фактуру по ID
	GetByID(ctx context.Context, id string) (*domain.Invoice, error)

	// GetByDocumentUUID получает счет-фактуру по DocumentUUID
	GetByDocumentUUID(ctx context.Context, uuid string) (*domain.Invoice, error)

	// GetByInvoiceNumber получает счет-фактуру по номеру счета
	GetByInvoiceNumber(ctx context.Context, invoiceNumber string) (*domain.Invoice, error)

	// Update обновляет счет-фактуру
	Update(ctx context.Context, invoice *domain.Invoice) error

	// List возвращает список счетов-фактур с пагинацией
	List(ctx context.Context, page, perPage int32, status string) ([]*domain.Invoice, int32, error)

	// ListByDateRange возвращает список счетов-фактур по диапазону дат с пагинацией
	ListByDateRange(ctx context.Context, startDate, endDate string, page, perPage int32) ([]*domain.Invoice, int32, error)

	// ExistsByNumber проверяет существование счета по номеру
	ExistsByNumber(ctx context.Context, invoiceNumber string) (bool, error)
}

// InvoiceDetailRepository определяет интерфейс для работы с позициями счета
type InvoiceDetailRepository interface {
	// Create создает позицию счета
	Create(ctx context.Context, detail *domain.InvoiceDetail) error

	// CreateBatch создает несколько позиций счета
	CreateBatch(ctx context.Context, details []*domain.InvoiceDetail) error

	// GetByInvoiceUUID получает позиции счета по UUID счета
	GetByInvoiceUUID(ctx context.Context, invoiceUUID string) ([]*domain.InvoiceDetail, error)

	// ListByInvoiceUUID получает позиции счета по UUID счета с пагинацией и общим количеством
	ListByInvoiceUUID(ctx context.Context, invoiceUUID string, page, perPage int32) ([]*domain.InvoiceDetail, int32, error)

	// DeleteByInvoiceUUID удаляет все позиции счета
	DeleteByInvoiceUUID(ctx context.Context, invoiceUUID string) error
}

// FinancialDataRepository определяет интерфейс для работы с финансовыми данными
type FinancialDataRepository interface {
	// Create создает финансовые данные
	Create(ctx context.Context, data *domain.FinancialData) error

	// GetByInvoiceUUID получает финансовые данные по UUID счета
	GetByInvoiceUUID(ctx context.Context, invoiceUUID string) (*domain.FinancialData, error)

	// Update обновляет финансовые данные
	Update(ctx context.Context, data *domain.FinancialData) error
}
