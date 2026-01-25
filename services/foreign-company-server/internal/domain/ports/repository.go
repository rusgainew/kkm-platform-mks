package ports

import (
	"context"

	"github.com/rusgainew/kkm-project-mks/foreign-company-server/internal/domain"
)

// ForeignCompanyRepository определяет контракт для работы с хранилищем иностранных компаний
type ForeignCompanyRepository interface {
	// Create создает новую иностранную компанию
	Create(ctx context.Context, company *domain.ForeignCompany) error

	// GetByID возвращает компанию по ID
	GetByID(ctx context.Context, id int64) (*domain.ForeignCompany, error)

	// GetByPIN возвращает компанию по PIN
	GetByPIN(ctx context.Context, pin string) (*domain.ForeignCompany, error)

	// Update обновляет данные компании
	Update(ctx context.Context, company *domain.ForeignCompany) error

	// Delete удаляет компанию
	Delete(ctx context.Context, id int64) error

	// List возвращает список компаний с пагинацией
	List(ctx context.Context, limit, offset int) ([]*domain.ForeignCompany, int, error)

	// Search ищет компании по названию
	Search(ctx context.Context, query string, limit, offset int) ([]*domain.ForeignCompany, int, error)

	// ExistsByPIN проверяет существование компании с таким PIN
	ExistsByPIN(ctx context.Context, pin string) (bool, error)

	// ListByCountry возвращает компании из определенной страны
	ListByCountry(ctx context.Context, countryCode string, limit, offset int) ([]*domain.ForeignCompany, int, error)
}
