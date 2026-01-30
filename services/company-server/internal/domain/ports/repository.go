// Файл company-server/internal/domain/ports/repository.go содержит реализацию пакета ports.
package ports

import (
	"context"

	"github.com/rusgainew/kkm-project-mks/company-server/internal/domain"
)

// OrganizationRepository определяет интерфейс для работы с хранилищем организаций
type OrganizationRepository interface {
	// Create создает новую организацию
	Create(ctx context.Context, org *domain.Organization) error

	// GetByID получает организацию по ID
	GetByID(ctx context.Context, id string) (*domain.Organization, error)

	// Update обновляет организацию
	Update(ctx context.Context, org *domain.Organization) error

	// Delete удаляет организацию
	Delete(ctx context.Context, id string) error

	// List возвращает список организаций с пагинацией
	List(ctx context.Context, page, perPage int32, ownerID string) ([]*domain.Organization, int32, error)

	// ExistsByName проверяет существование организации по имени
	ExistsByName(ctx context.Context, name string) (bool, error)
}

// EmployeeRepository определяет интерфейс для работы с участниками
type EmployeeRepository interface {
	// Create добавляет участника в организацию
	Create(ctx context.Context, emp *domain.Employee) error

	// GetByID получает участника по ID
	GetByID(ctx context.Context, id string) (*domain.Employee, error)

	// GetByUserAndOrganization получает участника по user_id и organization_id
	GetByUserAndOrganization(ctx context.Context, userID, organizationID string) (*domain.Employee, error)

	// Delete удаляет участника
	Delete(ctx context.Context, id string) error

	// ListByOrganization возвращает список участников организации
	ListByOrganization(ctx context.Context, organizationID string, page, perPage int32) ([]*domain.Employee, int32, error)

	// ExistsInOrganization проверяет, есть ли пользователь в организации
	ExistsInOrganization(ctx context.Context, userID, organizationID string) (bool, error)
}
