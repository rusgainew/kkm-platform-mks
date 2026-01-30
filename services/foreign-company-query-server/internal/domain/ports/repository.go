// Файл foreign-company-query-server/internal/domain/ports/repository.go содержит реализацию пакета ports.
package ports

import (
	"context"

	"github.com/rusgainew/kkm-project-mks/proto-lib/dictionaries"
)

// ForeignCompanyFilter содержит параметры фильтрации для иностранных компаний
type ForeignCompanyFilter struct {
	PIN         string // Поиск по PIN (точное совпадение)
	FullName    string // Поиск по названию (ILIKE)
	CountryCode string // Фильтр по коду страны
	SearchText  string // Полнотекстовый поиск по PIN и FullName
}

// SortOrder определяет направление сортировки
type SortOrder string

const (
	SortOrderAsc  SortOrder = "ASC"
	SortOrderDesc SortOrder = "DESC"
)

// ForeignCompanySort содержит параметры сортировки
type ForeignCompanySort struct {
	Field string    // Поле для сортировки: pin, full_name, country_code
	Order SortOrder // Направление: ASC/DESC
}

// ForeignCompanyQueryRepository определяет контракт для запросов к иностранным компаниям
type ForeignCompanyQueryRepository interface {
	// ListForeignCompanies возвращает список иностранных компаний с пагинацией
	ListForeignCompanies(ctx context.Context, page, size int32) ([]*dictionaries.ForeignCompany, int32, error)

	// ListForeignCompaniesWithFilter возвращает список иностранных компаний с фильтрацией и сортировкой
	ListForeignCompaniesWithFilter(ctx context.Context, filter *ForeignCompanyFilter, sort *ForeignCompanySort, page, size int32) ([]*dictionaries.ForeignCompany, int32, error)

	// SearchForeignCompanies выполняет полнотекстовый поиск по иностранным компаниям
	SearchForeignCompanies(ctx context.Context, searchText string, page, size int32) ([]*dictionaries.ForeignCompany, int32, error)

	// GetForeignCompanyByPIN возвращает иностранную компанию по PIN
	GetForeignCompanyByPIN(ctx context.Context, pin string) (*dictionaries.ForeignCompany, error)

	// GetForeignCompaniesByCountry возвращает компании по коду страны
	GetForeignCompaniesByCountry(ctx context.Context, countryCode string, page, size int32) ([]*dictionaries.ForeignCompany, int32, error)
}
