// Файл invoice-query-server/internal/domain/ports/repository.go содержит реализацию пакета ports.
package ports

import (
	"context"

	"github.com/rusgainew/kkm-project-mks/proto-lib/entities"
)

// InvoiceFilter содержит параметры фильтрации для счетов-фактур
type InvoiceFilter struct {
	InvoiceNumber string  // Точное совпадение по номеру счета
	DateFrom      string  // Дата начала диапазона (YYYY-MM-DD)
	DateTo        string  // Дата окончания диапазона (YYYY-MM-DD)
	MinAmount     float64 // Минимальная сумма
	MaxAmount     float64 // Максимальная сумма
	SearchText    string  // Поиск по invoice_number и note (ILIKE)
}

// SortOrder определяет направление сортировки
type SortOrder string

const (
	SortOrderAsc  SortOrder = "ASC"
	SortOrderDesc SortOrder = "DESC"
)

// InvoiceSort содержит параметры сортировки
type InvoiceSort struct {
	Field string    // Поле для сортировки: invoice_date, total_amount, invoice_number, created_date
	Order SortOrder // Направление: ASC/DESC
}

// InvoiceQueryRepository определяет контракт для запросов к счетам-фактурам
type InvoiceQueryRepository interface {
	// ListInvoices возвращает список счетов-фактур с пагинацией
	// page - номер страницы (начиная с 0)
	// size - количество элементов на странице
	// Возвращает: список счетов, общее количество элементов, ошибку
	ListInvoices(ctx context.Context, page, size int32) ([]*entities.Invoice, int32, error)

	// ListInvoiceDetails возвращает список деталей счетов-фактур с пагинацией
	ListInvoiceDetails(ctx context.Context, page, size int32) ([]*entities.InvoiceDetail, int32, error)

	// ListInvoicesWithFilter возвращает список счетов-фактур с фильтрацией, сортировкой и пагинацией
	ListInvoicesWithFilter(ctx context.Context, filter *InvoiceFilter, sort *InvoiceSort, page, size int32) ([]*entities.Invoice, int32, error)

	// SearchInvoices выполняет полнотекстовый поиск по счетам-фактурам
	SearchInvoices(ctx context.Context, searchText string, page, size int32) ([]*entities.Invoice, int32, error)

	// GetInvoiceByNumber возвращает счет-фактуру по номеру
	GetInvoiceByNumber(ctx context.Context, invoiceNumber string) (*entities.Invoice, error)

	// GetInvoicesByDateRange возвращает счета-фактуры за указанный период
	GetInvoicesByDateRange(ctx context.Context, dateFrom, dateTo string, page, size int32) ([]*entities.Invoice, int32, error)
}
