// Файл bank-account-query-server/internal/infrastructure/repository/postgres_bank_account_repository.go содержит реализацию пакета repository.
package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/rusgainew/kkm-project-mks/bank-account-query-server/internal/domain/ports"
	"github.com/rusgainew/kkm-project-mks/proto-lib/dictionaries"
	"go.uber.org/zap"
)

type PostgresBankAccountRepository struct {
	db     *sql.DB
	logger *zap.Logger
}

func NewPostgresBankAccountRepository(db *sql.DB, logger *zap.Logger) ports.BankAccountQueryRepository {
	return &PostgresBankAccountRepository{
		db:     db,
		logger: logger,
	}
}

func (r *PostgresBankAccountRepository) ListBankAccounts(ctx context.Context, page, size int32) ([]*dictionaries.BankAccount, int32, error) {
	// Исправлена пагинация: offset = (page - 1) * size
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * size

	query := `
		SELECT id::text, bank_name, account_number, 
		       organization_id::text, is_active
		FROM bank_accounts
		WHERE is_active = true
		ORDER BY bank_name
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.QueryContext(ctx, query, size, offset)
	if err != nil {
		r.logger.Error("Failed to query bank accounts", zap.Error(err))
		return nil, 0, fmt.Errorf("query bank accounts: %w", err)
	}
	defer rows.Close()

	var accounts []*dictionaries.BankAccount
	for rows.Next() {
		account := &dictionaries.BankAccount{}

		err := rows.Scan(
			&account.Id,
			&account.AccountName,
			&account.BankAccount,
			&account.ContractorTin,
			&account.IsActive,
		)
		if err != nil {
			r.logger.Error("Failed to scan bank account row", zap.Error(err))
			return nil, 0, fmt.Errorf("scan bank account: %w", err)
		}

		accounts = append(accounts, account)
	}

	if err := rows.Err(); err != nil {
		r.logger.Error("Rows iteration error", zap.Error(err))
		return nil, 0, fmt.Errorf("rows error: %w", err)
	}

	var totalCount int32
	countQuery := "SELECT COUNT(*) FROM bank_accounts WHERE is_active = true"
	err = r.db.QueryRowContext(ctx, countQuery).Scan(&totalCount)
	if err != nil {
		r.logger.Error("Failed to count bank accounts", zap.Error(err))
		return nil, 0, fmt.Errorf("count bank accounts: %w", err)
	}

	r.logger.Info("Listed bank accounts",
		zap.Int32("page", page),
		zap.Int32("size", size),
		zap.Int("returned", len(accounts)),
		zap.Int32("total", totalCount),
	)

	return accounts, totalCount, nil
}

// ListBankAccountsWithFilter возвращает список банковских счетов с фильтрацией и сортировкой
func (r *PostgresBankAccountRepository) ListBankAccountsWithFilter(ctx context.Context, filter *ports.BankAccountFilter, sort *ports.BankAccountSort, page, size int32) ([]*dictionaries.BankAccount, int32, error) {
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * size

	// Строим WHERE clause динамически
	whereClause := " WHERE is_active = true"
	args := []interface{}{}
	argIndex := 1

	if filter != nil {
		if filter.AccountName != "" {
			whereClause += fmt.Sprintf(" AND bank_name ILIKE $%d", argIndex)
			args = append(args, "%"+filter.AccountName+"%")
			argIndex++
		}

		if filter.BankAccount != "" {
			whereClause += fmt.Sprintf(" AND account_number = $%d", argIndex)
			args = append(args, filter.BankAccount)
			argIndex++
		}

		if filter.ContractorTin != "" {
			whereClause += fmt.Sprintf(" AND organization_id::text = $%d", argIndex)
			args = append(args, filter.ContractorTin)
			argIndex++
		}

		if filter.IsActive != nil {
			// Перезаписываем условие is_active если явно указано
			if !*filter.IsActive {
				whereClause = strings.Replace(whereClause, "is_active = true", "is_active = false", 1)
			}
		}

		if filter.SearchText != "" {
			whereClause += fmt.Sprintf(" AND (bank_name ILIKE $%d OR account_number ILIKE $%d)", argIndex, argIndex+1)
			searchPattern := "%" + filter.SearchText + "%"
			args = append(args, searchPattern, searchPattern)
			argIndex += 2
		}
	}

	// Определяем ORDER BY
	orderBy := "bank_name ASC" // По умолчанию
	if sort != nil && sort.Field != "" {
		validFields := map[string]string{
			"account_name":   "bank_name",
			"bank_account":   "account_number",
			"contractor_tin": "organization_id",
			"is_active":      "is_active",
		}
		if mappedField, ok := validFields[sort.Field]; ok {
			// Валидация sort.Order для предотвращения SQL injection
			sortOrder := "ASC"
			if sort.Order == "DESC" || sort.Order == "desc" {
				sortOrder = "DESC"
			}
			orderBy = mappedField + " " + sortOrder
		}
	}

	// SQL запрос с фильтрацией и сортировкой
	query := `
		SELECT id::text, bank_name, account_number, organization_id::text, is_active
		FROM bank_accounts
		` + whereClause + `
		ORDER BY ` + orderBy + `
		LIMIT $` + fmt.Sprintf("%d", argIndex) + ` OFFSET $` + fmt.Sprintf("%d", argIndex+1)

	args = append(args, size, offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		r.logger.Error("Failed to query bank accounts with filter", zap.Error(err))
		return nil, 0, fmt.Errorf("query bank accounts with filter: %w", err)
	}
	defer rows.Close()

	var accounts []*dictionaries.BankAccount
	for rows.Next() {
		account := &dictionaries.BankAccount{}
		err := rows.Scan(
			&account.Id,
			&account.AccountName,
			&account.BankAccount,
			&account.ContractorTin,
			&account.IsActive,
		)
		if err != nil {
			r.logger.Error("Failed to scan bank account", zap.Error(err))
			return nil, 0, fmt.Errorf("scan bank account: %w", err)
		}
		accounts = append(accounts, account)
	}

	// Получаем общее количество записей с учетом фильтров
	countQuery := "SELECT COUNT(*) FROM bank_accounts" + whereClause
	var totalCount int32
	countArgs := args[:len(args)-2] // Убираем LIMIT и OFFSET
	err = r.db.QueryRowContext(ctx, countQuery, countArgs...).Scan(&totalCount)
	if err != nil {
		r.logger.Error("Failed to count bank accounts", zap.Error(err))
		return nil, 0, fmt.Errorf("count bank accounts: %w", err)
	}

	r.logger.Info("Listed bank accounts with filter",
		zap.Int("count", len(accounts)),
		zap.Int32("total", totalCount),
		zap.Int32("page", page),
		zap.Int32("size", size))

	return accounts, totalCount, nil
}

// SearchBankAccounts выполняет полнотекстовый поиск по банковским счетам
func (r *PostgresBankAccountRepository) SearchBankAccounts(ctx context.Context, searchText string, page, size int32) ([]*dictionaries.BankAccount, int32, error) {
	filter := &ports.BankAccountFilter{
		SearchText: searchText,
	}
	return r.ListBankAccountsWithFilter(ctx, filter, nil, page, size)
}

// GetBankAccountByNumber возвращает банковский счет по номеру
func (r *PostgresBankAccountRepository) GetBankAccountByNumber(ctx context.Context, bankAccount string) (*dictionaries.BankAccount, error) {
	query := `
		SELECT id::text, bank_name, account_number, organization_id::text, is_active
		FROM bank_accounts
		WHERE account_number = $1
		LIMIT 1
	`

	account := &dictionaries.BankAccount{}
	err := r.db.QueryRowContext(ctx, query, bankAccount).Scan(
		&account.Id,
		&account.AccountName,
		&account.BankAccount,
		&account.ContractorTin,
		&account.IsActive,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("bank account not found: %s", bankAccount)
	}
	if err != nil {
		r.logger.Error("Failed to get bank account by number", zap.Error(err), zap.String("bank_account", bankAccount))
		return nil, fmt.Errorf("get bank account by number: %w", err)
	}

	r.logger.Info("Got bank account by number", zap.String("bank_account", bankAccount))

	return account, nil
}

// GetActiveBankAccounts возвращает только активные банковские счета
func (r *PostgresBankAccountRepository) GetActiveBankAccounts(ctx context.Context, page, size int32) ([]*dictionaries.BankAccount, int32, error) {
	isActive := true
	filter := &ports.BankAccountFilter{
		IsActive: &isActive,
	}
	sort := &ports.BankAccountSort{
		Field: "account_name",
		Order: ports.SortOrderAsc,
	}
	return r.ListBankAccountsWithFilter(ctx, filter, sort, page, size)
}
