package ports

import (
	"context"

	"github.com/google/uuid"
	"github.com/rusgainew/kkm-project-mks/catalog-server/internal/domain"
)

// CatalogRepository определяет контракт для работы с хранилищем каталога
type CatalogRepository interface {
	// Create создает новый элемент каталога
	Create(ctx context.Context, item *domain.CatalogItem) error

	// GetByID возвращает элемент по ID
	GetByID(ctx context.Context, id uuid.UUID) (*domain.CatalogItem, error)

	// GetByCode возвращает элемент по коду и организации
	GetByCode(ctx context.Context, organizationID uuid.UUID, code string) (*domain.CatalogItem, error)

	// Update обновляет элемент каталога
	Update(ctx context.Context, item *domain.CatalogItem) error

	// Delete удаляет элемент из каталога
	Delete(ctx context.Context, id uuid.UUID) error

	// ListByOrganization возвращает список элементов организации
	ListByOrganization(ctx context.Context, organizationID uuid.UUID, limit, offset int) ([]*domain.CatalogItem, int, error)

	// Search ищет элементы по имени или коду
	Search(ctx context.Context, organizationID uuid.UUID, query string, limit, offset int) ([]*domain.CatalogItem, int, error)

	// ExistsByCode проверяет существование элемента с кодом
	ExistsByCode(ctx context.Context, organizationID uuid.UUID, code string) (bool, error)
}
