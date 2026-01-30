// Файл invoice-server/internal/infrastructure/repository/postgres_invoice_repo.go содержит реализацию пакета repository.
package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/rusgainew/kkm-project-mks/invoice-server/internal/domain"
)

// PostgresInvoiceRepository реализация репозитория счетов для PostgreSQL
type PostgresInvoiceRepository struct {
	db *sqlx.DB
}

// NewPostgresInvoiceRepository создает новый экземпляр репозитория
func NewPostgresInvoiceRepository(db *sqlx.DB) *PostgresInvoiceRepository {
	return &PostgresInvoiceRepository{db: db}
}

// Create создает новый счет-фактуру
func (r *PostgresInvoiceRepository) Create(ctx context.Context, invoice *domain.Invoice) error {
	query := `
		INSERT INTO invoices (
			id, document_uuid, invoice_number, number, corrected_receipt_uuid,
			invoice_date, created_date, delivery_date, corrected_receipt_creation_date,
			total_amount, is_resident, note, status, legal_person_id, contractor_id,
			created_by, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17)
	`
	_, err := r.db.ExecContext(ctx, query,
		invoice.ID, invoice.DocumentUUID, invoice.InvoiceNumber, invoice.Number,
		invoice.CorrectedReceiptUUID, invoice.InvoiceDate, invoice.CreatedDate,
		invoice.DeliveryDate, invoice.CorrectedReceiptCreationDate, invoice.TotalAmount,
		invoice.IsResident, invoice.Note, invoice.Status, invoice.LegalPersonID,
		invoice.ContractorID, invoice.CreatedBy, invoice.UpdatedAt,
	)
	return err
}

// GetByID получает счет по ID
func (r *PostgresInvoiceRepository) GetByID(ctx context.Context, id string) (*domain.Invoice, error) {
	var invoice domain.Invoice
	query := `
		SELECT id, document_uuid, invoice_number, number, corrected_receipt_uuid,
			   invoice_date, created_date, delivery_date, corrected_receipt_creation_date,
			   total_amount, is_resident, note, status, legal_person_id, contractor_id,
			   created_by, updated_at, signed_at, signed_by
		FROM invoices WHERE id = $1
	`

	err := r.db.GetContext(ctx, &invoice, query, id)
	if err == sql.ErrNoRows {
		return nil, domain.ErrInvoiceNotFound
	}
	if err != nil {
		return nil, err
	}

	return &invoice, nil
}

// GetByDocumentUUID получает счет по DocumentUUID
func (r *PostgresInvoiceRepository) GetByDocumentUUID(ctx context.Context, uuid string) (*domain.Invoice, error) {
	var invoice domain.Invoice
	query := `
		SELECT id, document_uuid, invoice_number, number, corrected_receipt_uuid,
			   invoice_date, created_date, delivery_date, corrected_receipt_creation_date,
			   total_amount, is_resident, note, status, legal_person_id, contractor_id,
			   created_by, updated_at, signed_at, signed_by
		FROM invoices WHERE document_uuid = $1
	`

	err := r.db.GetContext(ctx, &invoice, query, uuid)
	if err == sql.ErrNoRows {
		return nil, domain.ErrInvoiceNotFound
	}
	if err != nil {
		return nil, err
	}

	return &invoice, nil
}

// Update обновляет счет
func (r *PostgresInvoiceRepository) Update(ctx context.Context, invoice *domain.Invoice) error {
	query := `
		UPDATE invoices 
		SET invoice_number = $2, number = $3, total_amount = $4, note = $5,
			status = $6, updated_at = $7, signed_at = $8, signed_by = $9
		WHERE id = $1
	`
	result, err := r.db.ExecContext(ctx, query,
		invoice.ID, invoice.InvoiceNumber, invoice.Number, invoice.TotalAmount,
		invoice.Note, invoice.Status, invoice.UpdatedAt, invoice.SignedAt, invoice.SignedBy,
	)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return domain.ErrInvoiceNotFound
	}

	return nil
}

// List возвращает список счетов с пагинацией
func (r *PostgresInvoiceRepository) List(ctx context.Context, page, perPage int32, status string) ([]*domain.Invoice, int32, error) {
	var invoices []*domain.Invoice
	var total int32

	offset := (page - 1) * perPage

	baseQuery := `FROM invoices`
	whereClause := ""
	args := []interface{}{}
	argPosition := 1

	if status != "" {
		whereClause = fmt.Sprintf(" WHERE status = $%d", argPosition)
		args = append(args, status)
		argPosition++
	}

	// Получение общего количества
	countQuery := `SELECT COUNT(*) ` + baseQuery + whereClause
	err := r.db.GetContext(ctx, &total, countQuery, args...)
	if err != nil {
		return nil, 0, err
	}

	// Получение списка
	args = append(args, perPage, offset)
	selectQuery := `
		SELECT id, document_uuid, invoice_number, number, corrected_receipt_uuid,
			   invoice_date, created_date, delivery_date, corrected_receipt_creation_date,
			   total_amount, is_resident, note, status, legal_person_id, contractor_id,
			   created_by, updated_at, signed_at, signed_by
	` + baseQuery + whereClause +
		fmt.Sprintf(` ORDER BY created_date DESC LIMIT $%d OFFSET $%d`, argPosition, argPosition+1)

	err = r.db.SelectContext(ctx, &invoices, selectQuery, args...)
	if err != nil {
		return nil, 0, err
	}

	return invoices, total, nil
}

// ExistsByNumber проверяет существование счета по номеру
func (r *PostgresInvoiceRepository) ExistsByNumber(ctx context.Context, invoiceNumber string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM invoices WHERE invoice_number = $1)`
	err := r.db.GetContext(ctx, &exists, query, invoiceNumber)
	return exists, err
}

// GetByInvoiceNumber получает счет по номеру счета
func (r *PostgresInvoiceRepository) GetByInvoiceNumber(ctx context.Context, invoiceNumber string) (*domain.Invoice, error) {
	var invoice domain.Invoice
	query := `
		SELECT id, document_uuid, invoice_number, number, corrected_receipt_uuid,
			   invoice_date, created_date, delivery_date, corrected_receipt_creation_date,
			   total_amount, is_resident, note, status, legal_person_id, contractor_id,
			   created_by, updated_at, signed_at, signed_by
		FROM invoices WHERE invoice_number = $1
	`

	err := r.db.GetContext(ctx, &invoice, query, invoiceNumber)
	if err == sql.ErrNoRows {
		return nil, domain.ErrInvoiceNotFound
	}
	if err != nil {
		return nil, err
	}

	return &invoice, nil
}

// ListByDateRange возвращает список счетов в указанном диапазоне дат с пагинацией
func (r *PostgresInvoiceRepository) ListByDateRange(ctx context.Context, startDate, endDate string, page, perPage int32) ([]*domain.Invoice, int32, error) {
	var invoices []*domain.Invoice
	var total int32

	offset := (page - 1) * perPage

	baseQuery := `FROM invoices`
	whereClause := ""
	args := []interface{}{}
	argPosition := 1

	if startDate != "" && endDate != "" {
		whereClause = fmt.Sprintf(" WHERE created_date >= $%d AND created_date <= $%d", argPosition, argPosition+1)
		args = append(args, startDate, endDate)
		argPosition += 2
	}

	// Получение общего количества
	countQuery := `SELECT COUNT(*) ` + baseQuery + whereClause
	err := r.db.GetContext(ctx, &total, countQuery, args...)
	if err != nil {
		return nil, 0, err
	}

	// Получение списка
	args = append(args, perPage, offset)
	selectQuery := `
		SELECT id, document_uuid, invoice_number, number, corrected_receipt_uuid,
			   invoice_date, created_date, delivery_date, corrected_receipt_creation_date,
			   total_amount, is_resident, note, status, legal_person_id, contractor_id,
			   created_by, updated_at, signed_at, signed_by
	` + baseQuery + whereClause +
		fmt.Sprintf(` ORDER BY created_date DESC LIMIT $%d OFFSET $%d`, argPosition, argPosition+1)

	err = r.db.SelectContext(ctx, &invoices, selectQuery, args...)
	if err != nil {
		return nil, 0, err
	}

	return invoices, total, nil
}
