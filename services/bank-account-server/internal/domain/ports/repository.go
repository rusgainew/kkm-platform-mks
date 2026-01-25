package ports

import (
	"context"

	"github.com/google/uuid"
	"github.com/rusgainew/kkm-project-mks/bank-account-server/internal/domain"
)

// BankAccountRepository определяет контракт для работы с хранилищем банковских счетов
type BankAccountRepository interface {
	// Create создает новый банковский счет
	Create(ctx context.Context, account *domain.BankAccount) error

	// GetByID возвращает счет по ID
	GetByID(ctx context.Context, id uuid.UUID) (*domain.BankAccount, error)

	// GetByAccountNumber возвращает счет по номеру и организации
	GetByAccountNumber(ctx context.Context, organizationID uuid.UUID, accountNumber string) (*domain.BankAccount, error)

	// Update обновляет банковский счет
	Update(ctx context.Context, account *domain.BankAccount) error

	// Delete удаляет банковский счет
	Delete(ctx context.Context, id uuid.UUID) error

	// ListByOrganization возвращает список счетов организации
	ListByOrganization(ctx context.Context, organizationID uuid.UUID, limit, offset int) ([]*domain.BankAccount, int, error)

	// GetDefaultAccount возвращает счет по умолчанию для организации
	GetDefaultAccount(ctx context.Context, organizationID uuid.UUID) (*domain.BankAccount, error)

	// SetDefaultAccount устанавливает счет как счет по умолчанию
	SetDefaultAccount(ctx context.Context, organizationID, accountID uuid.UUID) error

	// ExistsByAccountNumber проверяет существование счета с номером
	ExistsByAccountNumber(ctx context.Context, organizationID uuid.UUID, accountNumber string) (bool, error)

	// CountActiveAccounts подсчитывает активные счета организации
	CountActiveAccounts(ctx context.Context, organizationID uuid.UUID) (int, error)
}
