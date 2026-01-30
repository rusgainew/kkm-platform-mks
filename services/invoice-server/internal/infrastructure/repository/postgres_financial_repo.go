// Файл invoice-server/internal/infrastructure/repository/postgres_financial_repo.go содержит реализацию пакета repository.
package repository

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
	"github.com/rusgainew/kkm-project-mks/invoice-server/internal/domain"
)

// PostgresFinancialDataRepository реализация репозитория финансовых данных
type PostgresFinancialDataRepository struct {
	db *sqlx.DB
}

// NewPostgresFinancialDataRepository создает новый экземпляр репозитория
func NewPostgresFinancialDataRepository(db *sqlx.DB) *PostgresFinancialDataRepository {
	return &PostgresFinancialDataRepository{db: db}
}

// Create создает финансовые данные
func (r *PostgresFinancialDataRepository) Create(ctx context.Context, data *domain.FinancialData) error {
	query := `
		INSERT INTO invoice_financials (
			invoice_uuid, total_amount, opening_balances, assessed_contributions_amount,
			paid_amount, penalties_amount, fines_amount, closing_balances,
			amount_to_be_paid, personal_account_number, legal_person_bank_account,
			contractor_bank_account
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
	`
	_, err := r.db.ExecContext(ctx, query,
		data.InvoiceUUID, data.TotalAmount, data.OpeningBalances, data.AssessedContributionsAmount,
		data.PaidAmount, data.PenaltiesAmount, data.FinesAmount, data.ClosingBalances,
		data.AmountToBePaid, data.PersonalAccountNumber, data.LegalPersonBankAccount,
		data.ContractorBankAccount,
	)
	return err
}

// GetByInvoiceUUID получает финансовые данные по UUID счета
func (r *PostgresFinancialDataRepository) GetByInvoiceUUID(ctx context.Context, invoiceUUID string) (*domain.FinancialData, error) {
	var data domain.FinancialData
	query := `
		SELECT invoice_uuid, total_amount, opening_balances, assessed_contributions_amount,
			   paid_amount, penalties_amount, fines_amount, closing_balances,
			   amount_to_be_paid, personal_account_number, legal_person_bank_account,
			   contractor_bank_account
		FROM invoice_financials 
		WHERE invoice_uuid = $1
	`

	err := r.db.GetContext(ctx, &data, query, invoiceUUID)
	if err == sql.ErrNoRows {
		return nil, nil // Financial data is optional
	}
	if err != nil {
		return nil, err
	}

	return &data, nil
}

// Update обновляет финансовые данные
func (r *PostgresFinancialDataRepository) Update(ctx context.Context, data *domain.FinancialData) error {
	query := `
		UPDATE invoice_financials 
		SET total_amount = $2, opening_balances = $3, assessed_contributions_amount = $4,
			paid_amount = $5, penalties_amount = $6, fines_amount = $7, closing_balances = $8,
			amount_to_be_paid = $9, personal_account_number = $10, legal_person_bank_account = $11,
			contractor_bank_account = $12
		WHERE invoice_uuid = $1
	`
	_, err := r.db.ExecContext(ctx, query,
		data.InvoiceUUID, data.TotalAmount, data.OpeningBalances, data.AssessedContributionsAmount,
		data.PaidAmount, data.PenaltiesAmount, data.FinesAmount, data.ClosingBalances,
		data.AmountToBePaid, data.PersonalAccountNumber, data.LegalPersonBankAccount,
		data.ContractorBankAccount,
	)
	return err
}
