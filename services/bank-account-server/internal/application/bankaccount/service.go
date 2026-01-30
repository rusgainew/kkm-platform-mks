// Файл bank-account-server/internal/application/bankaccount/service.go содержит реализацию пакета bankaccount.
package bankaccount

import (
	"context"
	"fmt"
	"regexp"

	"github.com/google/uuid"
	"github.com/rusgainew/kkm-project-mks/bank-account-server/internal/domain"
	"github.com/rusgainew/kkm-project-mks/bank-account-server/internal/domain/events"
	"github.com/rusgainew/kkm-project-mks/bank-account-server/internal/domain/ports"
	"go.uber.org/zap"
)

// Регулярные выражения для валидации
var (
	accountNumberRegex = regexp.MustCompile(`^[A-Z0-9]{5,34}$`)                                  // Гибкий формат для разных стран
	ibanRegex          = regexp.MustCompile(`^[A-Z]{2}[0-9]{2}[A-Z0-9]+$`)                       // Международный IBAN
	bicRegex           = regexp.MustCompile(`^([A-Z]{6}[A-Z0-9]{2}([A-Z0-9]{3})?|[0-9]{5,11})$`) // SWIFT или БИК (цифровой)
	currencyRegex      = regexp.MustCompile(`^[A-Z]{3}$`)
)

// Service реализует бизнес-логику управления банковскими счетами
type Service struct {
	repo      ports.BankAccountRepository
	publisher ports.EventPublisher
	logger    *zap.Logger
}

// NewService создает новый сервис банковских счетов
func NewService(repo ports.BankAccountRepository, publisher ports.EventPublisher, logger *zap.Logger) *Service {
	return &Service{
		repo:      repo,
		publisher: publisher,
		logger:    logger,
	}
}

// CreateAccount создает новый банковский счет
func (s *Service) CreateAccount(ctx context.Context, organizationID, bankID, createdBy uuid.UUID, accountNumber, iban, currency, bic, bankName string) (*domain.BankAccount, error) {
	// Устанавливаем значение по умолчанию для currency
	if currency == "" {
		currency = "KGS"
	}

	// Валидация
	if err := s.validateAccountData(accountNumber, iban, currency, bic); err != nil {
		return nil, err
	}

	// Проверка существования
	exists, err := s.repo.ExistsByAccountNumber(ctx, organizationID, accountNumber)
	if err != nil {
		s.logger.Error("failed to check account existence", zap.Error(err))
		return nil, fmt.Errorf("failed to check account existence: %w", err)
	}
	if exists {
		return nil, domain.ErrBankAccountAlreadyExists
	}

	// Создание счета
	account := domain.NewBankAccount(organizationID, bankID, createdBy, accountNumber, iban, currency, bic, bankName)

	if err := s.repo.Create(ctx, account); err != nil {
		s.logger.Error("failed to create bank account", zap.Error(err))
		return nil, fmt.Errorf("failed to create bank account: %w", err)
	}

	// Публикация события
	event := events.NewBankAccountCreatedEvent(account.ID, organizationID, createdBy, accountNumber, currency, bankName)
	if err := s.publisher.Publish(ctx, event); err != nil {
		s.logger.Warn("failed to publish bank account created event", zap.Error(err))
	}

	s.logger.Info("bank account created",
		zap.String("account_id", account.ID.String()),
		zap.String("organization_id", organizationID.String()),
	)

	return account, nil
}

// GetAccount возвращает банковский счет по ID
func (s *Service) GetAccount(ctx context.Context, id uuid.UUID) (*domain.BankAccount, error) {
	account, err := s.repo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error("failed to get bank account", zap.Error(err), zap.String("account_id", id.String()))
		return nil, err
	}
	return account, nil
}

// UpdateAccount обновляет банковский счет
func (s *Service) UpdateAccount(ctx context.Context, id, updatedBy uuid.UUID, accountNumber, iban, currency, bic, bankName string) (*domain.BankAccount, error) {
	// Получение счета
	account, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Валидация
	if err := s.validateAccountData(accountNumber, iban, currency, bic); err != nil {
		return nil, err
	}

	// Обновление
	account.Update(accountNumber, iban, currency, bic, bankName, updatedBy)

	if err := s.repo.Update(ctx, account); err != nil {
		s.logger.Error("failed to update bank account", zap.Error(err))
		return nil, fmt.Errorf("failed to update bank account: %w", err)
	}

	// Публикация события
	event := events.NewBankAccountUpdatedEvent(account.ID, account.OrganizationID, updatedBy, accountNumber, currency)
	if err := s.publisher.Publish(ctx, event); err != nil {
		s.logger.Warn("failed to publish bank account updated event", zap.Error(err))
	}

	s.logger.Info("bank account updated", zap.String("account_id", account.ID.String()))

	return account, nil
}

// DeleteAccount удаляет банковский счет
func (s *Service) DeleteAccount(ctx context.Context, id, deletedBy uuid.UUID) error {
	// Получение счета для проверки
	account, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	// Проверка: нельзя удалить счет по умолчанию, если есть другие активные счета
	if account.IsDefault {
		activeCount, err := s.repo.CountActiveAccounts(ctx, account.OrganizationID)
		if err != nil {
			return err
		}
		if activeCount > 1 {
			return domain.ErrCannotDeleteDefaultAccount
		}
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		s.logger.Error("failed to delete bank account", zap.Error(err))
		return fmt.Errorf("failed to delete bank account: %w", err)
	}

	// Публикация события
	event := events.NewBankAccountDeletedEvent(id, account.OrganizationID, deletedBy)
	if err := s.publisher.Publish(ctx, event); err != nil {
		s.logger.Warn("failed to publish bank account deleted event", zap.Error(err))
	}

	s.logger.Info("bank account deleted", zap.String("account_id", id.String()))

	return nil
}

// ListAccounts возвращает список банковских счетов организации
func (s *Service) ListAccounts(ctx context.Context, organizationID uuid.UUID, limit, offset int) ([]*domain.BankAccount, int, error) {
	accounts, total, err := s.repo.ListByOrganization(ctx, organizationID, limit, offset)
	if err != nil {
		s.logger.Error("failed to list bank accounts", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list bank accounts: %w", err)
	}
	return accounts, total, nil
}

// SetDefaultAccount устанавливает счет как счет по умолчанию
func (s *Service) SetDefaultAccount(ctx context.Context, accountID, updatedBy uuid.UUID) error {
	// Получение счета
	account, err := s.repo.GetByID(ctx, accountID)
	if err != nil {
		return err
	}

	// Установка счета по умолчанию (снимет флаг с других счетов)
	if err := s.repo.SetDefaultAccount(ctx, account.OrganizationID, accountID); err != nil {
		s.logger.Error("failed to set default bank account", zap.Error(err))
		return fmt.Errorf("failed to set default bank account: %w", err)
	}

	// Публикация события
	event := events.NewBankAccountSetDefaultEvent(accountID, account.OrganizationID, updatedBy)
	if err := s.publisher.Publish(ctx, event); err != nil {
		s.logger.Warn("failed to publish bank account set default event", zap.Error(err))
	}

	s.logger.Info("bank account set as default", zap.String("account_id", accountID.String()))

	return nil
}

// GetDefaultAccount возвращает счет по умолчанию для организации
func (s *Service) GetDefaultAccount(ctx context.Context, organizationID uuid.UUID) (*domain.BankAccount, error) {
	account, err := s.repo.GetDefaultAccount(ctx, organizationID)
	if err != nil {
		s.logger.Error("failed to get default bank account", zap.Error(err))
		return nil, err
	}
	return account, nil
}

// validateAccountData валидирует данные банковского счета
func (s *Service) validateAccountData(accountNumber, iban, currency, bic string) error {
	if accountNumber != "" && !accountNumberRegex.MatchString(accountNumber) {
		return domain.ErrInvalidAccountNumber
	}

	if iban != "" && !ibanRegex.MatchString(iban) {
		return domain.ErrInvalidIBAN
	}

	// Currency опционален - если не указан или не соответствует формату, пропускаем проверку
	// Значение по умолчанию устанавливается в CreateAccount
	if currency != "" && !currencyRegex.MatchString(currency) {
		return domain.ErrInvalidCurrency
	}

	if bic != "" && !bicRegex.MatchString(bic) {
		return domain.ErrInvalidBIC
	}

	return nil
}
