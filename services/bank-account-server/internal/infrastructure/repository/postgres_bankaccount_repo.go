package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/rusgainew/kkm-project-mks/bank-account-server/internal/domain"
)

// PostgresBankAccountRepository реализация репозитория банковских счетов для PostgreSQL
type PostgresBankAccountRepository struct {
	db *sqlx.DB
}

// NewPostgresBankAccountRepository создает новый экземпляр репозитория
func NewPostgresBankAccountRepository(db *sqlx.DB) *PostgresBankAccountRepository {
	return &PostgresBankAccountRepository{db: db}
}

// Create создает новый банковский счет
func (r *PostgresBankAccountRepository) Create(ctx context.Context, account *domain.BankAccount) error {
	query := `
		INSERT INTO bank_accounts (
			id, organization_id, account_number, bank_id, iban, currency, bic, bank_name,
			is_active, is_default, created_by, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
	`
	_, err := r.db.ExecContext(ctx, query,
		account.ID, account.OrganizationID, account.AccountNumber,
		account.BankID, account.IBAN, account.Currency, account.BIC, account.BankName,
		account.IsActive, account.IsDefault, account.CreatedBy, account.CreatedAt, account.UpdatedAt,
	)
	return err
}

// GetByID получает счет по ID
func (r *PostgresBankAccountRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.BankAccount, error) {
	var account domain.BankAccount
	query := `
		SELECT id, organization_id, account_number, bank_id, iban, currency, bic, bank_name,
			   is_active, is_default, created_by, updated_by, created_at, updated_at
		FROM bank_accounts WHERE id = $1
	`

	err := r.db.GetContext(ctx, &account, query, id)
	if err == sql.ErrNoRows {
		return nil, domain.ErrBankAccountNotFound
	}
	if err != nil {
		return nil, err
	}

	return &account, nil
}

// Update обновляет банковский счет
func (r *PostgresBankAccountRepository) Update(ctx context.Context, account *domain.BankAccount) error {
	query := `
		UPDATE bank_accounts
		SET account_number = $2, bank_id = $3, iban = $4, currency = $5, bic = $6, bank_name = $7,
			is_active = $8, is_default = $9, updated_by = $10, updated_at = $11
		WHERE id = $1
	`
	result, err := r.db.ExecContext(ctx, query,
		account.ID, account.AccountNumber, account.BankID, account.IBAN, account.Currency, account.BIC, account.BankName,
		account.IsActive, account.IsDefault, account.UpdatedBy, account.UpdatedAt,
	)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return domain.ErrBankAccountNotFound
	}

	return nil
}

// Delete удаляет банковский счет
func (r *PostgresBankAccountRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM bank_accounts WHERE id = $1`
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return domain.ErrBankAccountNotFound
	}

	return nil
}

// ExistsByAccountNumber проверяет существование счета по номеру
func (r *PostgresBankAccountRepository) ExistsByAccountNumber(ctx context.Context, organizationID uuid.UUID, accountNumber string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM bank_accounts WHERE organization_id = $1 AND account_number = $2)`
	err := r.db.GetContext(ctx, &exists, query, organizationID, accountNumber)
	return exists, err
}

// ListByOrganization возвращает список счетов организации с пагинацией
func (r *PostgresBankAccountRepository) ListByOrganization(ctx context.Context, organizationID uuid.UUID, limit, offset int) ([]*domain.BankAccount, int, error) {
	var accounts []*domain.BankAccount
	query := `
		SELECT id, organization_id, account_number, bank_id, iban, currency, bic, bank_name,
			   is_active, is_default, created_by, updated_by, created_at, updated_at
		FROM bank_accounts
		WHERE organization_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	err := r.db.SelectContext(ctx, &accounts, query, organizationID, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	// Получение общего количества
	var total int
	countQuery := `SELECT COUNT(*) FROM bank_accounts WHERE organization_id = $1`
	err = r.db.GetContext(ctx, &total, countQuery, organizationID)
	if err != nil {
		return nil, 0, err
	}

	return accounts, total, nil
}

// GetByAccountNumber получает счет по номеру и организации
func (r *PostgresBankAccountRepository) GetByAccountNumber(ctx context.Context, organizationID uuid.UUID, accountNumber string) (*domain.BankAccount, error) {
	var account domain.BankAccount
	query := `
		SELECT id, organization_id, account_number, bank_id, iban, currency, bic, bank_name,
			   is_active, is_default, created_by, updated_by, created_at, updated_at
		FROM bank_accounts WHERE organization_id = $1 AND account_number = $2
	`

	err := r.db.GetContext(ctx, &account, query, organizationID, accountNumber)
	if err == sql.ErrNoRows {
		return nil, domain.ErrBankAccountNotFound
	}
	if err != nil {
		return nil, err
	}

	return &account, nil
}

// GetDefaultAccount получает счет по умолчанию для организации
func (r *PostgresBankAccountRepository) GetDefaultAccount(ctx context.Context, organizationID uuid.UUID) (*domain.BankAccount, error) {
	// TODO: Implement default account logic
	return nil, domain.ErrBankAccountNotFound
}

// SetDefaultAccount устанавливает счет как счет по умолчанию
func (r *PostgresBankAccountRepository) SetDefaultAccount(ctx context.Context, organizationID, accountID uuid.UUID) error {
	// TODO: Implement default account logic
	return nil
}

// CountActiveAccounts подсчитывает активные счета организации
func (r *PostgresBankAccountRepository) CountActiveAccounts(ctx context.Context, organizationID uuid.UUID) (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM bank_accounts WHERE organization_id = $1`
	err := r.db.GetContext(ctx, &count, query, organizationID)
	return count, err
}

// Search ищет счета по запросу
func (r *PostgresBankAccountRepository) Search(ctx context.Context, organizationID uuid.UUID, query string, limit, offset int) ([]*domain.BankAccount, int, error) {
	var accounts []*domain.BankAccount
	searchQuery := `
		SELECT id, organization_id, account_number, bank_id, iban, currency, bic, bank_name,
			   is_active, is_default, created_by, updated_by, created_at, updated_at
		FROM bank_accounts
		WHERE organization_id = $1
		  AND (account_number ILIKE $2 OR bank_name ILIKE $2)
		ORDER BY created_at DESC
		LIMIT $3 OFFSET $4
	`
	searchPattern := fmt.Sprintf("%%%s%%", query)

	err := r.db.SelectContext(ctx, &accounts, searchQuery, organizationID, searchPattern, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	// Получение общего количества
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM bank_accounts
		WHERE organization_id = $1
		  AND (account_number ILIKE $2 OR bank_name ILIKE $2)
	`
	err = r.db.GetContext(ctx, &total, countQuery, organizationID, searchPattern)
	if err != nil {
		return nil, 0, err
	}

	return accounts, total, nil
}
