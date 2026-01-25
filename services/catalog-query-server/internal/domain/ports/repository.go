package ports

import (
	"context"

	"github.com/rusgainew/kkm-project-mks/proto-lib/dictionaries"
)

// CatalogFilter содержит параметры фильтрации для каталога
type CatalogFilter struct {
	Name       string // Поиск по названию (ILIKE)
	Number     string // Точное совпадение по номеру
	TnvedCode  string // Точное совпадение по коду ТН ВЭД
	GkedCode   string // Точное совпадение по коду ГКЭД
	SearchText string // Полнотекстовый поиск по name и number
}

// SortOrder определяет направление сортировки
type SortOrder string

const (
	SortOrderAsc  SortOrder = "ASC"
	SortOrderDesc SortOrder = "DESC"
)

// CatalogSort содержит параметры сортировки
type CatalogSort struct {
	Field string    // Поле для сортировки: name, number, tnved_code, gked_code
	Order SortOrder // Направление: ASC/DESC
}

// CatalogQueryRepository определяет контракт для запросов к каталогу
type CatalogQueryRepository interface {
	// ListCatalogs возвращает список элементов каталога с пагинацией
	ListCatalogs(ctx context.Context, page, size int32) ([]*dictionaries.Catalog, int32, error)

	// ListCatalogsWithFilter возвращает список элементов каталога с фильтрацией и сортировкой
	ListCatalogsWithFilter(ctx context.Context, filter *CatalogFilter, sort *CatalogSort, page, size int32) ([]*dictionaries.Catalog, int32, error)

	// SearchCatalogs выполняет полнотекстовый поиск по каталогу
	SearchCatalogs(ctx context.Context, searchText string, page, size int32) ([]*dictionaries.Catalog, int32, error)

	// GetCatalogByNumber возвращает элемент каталога по номеру
	GetCatalogByNumber(ctx context.Context, number string) (*dictionaries.Catalog, error)

	// GetCatalogsByTnvedCode возвращает элементы каталога по коду ТН ВЭД
	GetCatalogsByTnvedCode(ctx context.Context, tnvedCode string, page, size int32) ([]*dictionaries.Catalog, int32, error)
}
