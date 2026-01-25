package ports

import (
	"context"

	"github.com/rusgainew/kkm-project-mks/proto-lib/dictionaries"
)

// BankAccountFilter содержит параметры фильтрации для банковских счетов
type BankAccountFilter struct {
	AccountName   string // Поиск по названию счета (ILIKE)
	BankAccount   string // Точное совпадение по номеру счета
	ContractorTin string // Точное совпадение по ИНН контрагента
	IsActive      *bool  // Фильтр по активности (nil = все, true = активные, false = неактивные)
	SearchText    string // Полнотекстовый поиск по account_name и bank_account
}

// SortOrder определяет направление сортировки
type SortOrder string

const (
	SortOrderAsc  SortOrder = "ASC"
	SortOrderDesc SortOrder = "DESC"
)

// BankAccountSort содержит параметры сортировки
type BankAccountSort struct {
	Field string    // Поле для сортировки: account_name, bank_account, contractor_tin
	Order SortOrder // Направление: ASC/DESC
}

// BankAccountQueryRepository определяет контракт для запросов к банковским счетам
type BankAccountQueryRepository interface {
	// ListBankAccounts возвращает список банковских счетов с пагинацией
	ListBankAccounts(ctx context.Context, page, size int32) ([]*dictionaries.BankAccount, int32, error)

	// ListBankAccountsWithFilter возвращает список банковских счетов с фильтрацией и сортировкой
	ListBankAccountsWithFilter(ctx context.Context, filter *BankAccountFilter, sort *BankAccountSort, page, size int32) ([]*dictionaries.BankAccount, int32, error)

	// SearchBankAccounts выполняет полнотекстовый поиск по банковским счетам
	SearchBankAccounts(ctx context.Context, searchText string, page, size int32) ([]*dictionaries.BankAccount, int32, error)

	// GetBankAccountByNumber возвращает банковский счет по номеру
	GetBankAccountByNumber(ctx context.Context, bankAccount string) (*dictionaries.BankAccount, error)

	// GetActiveBankAccounts возвращает только активные банковские счета
	GetActiveBankAccounts(ctx context.Context, page, size int32) ([]*dictionaries.BankAccount, int32, error)
}
