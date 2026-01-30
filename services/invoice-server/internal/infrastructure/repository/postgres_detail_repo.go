// Файл invoice-server/internal/infrastructure/repository/postgres_detail_repo.go содержит реализацию пакета repository.
package repository

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/rusgainew/kkm-project-mks/invoice-server/internal/domain"
)

// PostgresInvoiceDetailRepository реализация репозитория позиций счета
type PostgresInvoiceDetailRepository struct {
	db *sqlx.DB
}

// NewPostgresInvoiceDetailRepository создает новый экземпляр репозитория
func NewPostgresInvoiceDetailRepository(db *sqlx.DB) *PostgresInvoiceDetailRepository {
	return &PostgresInvoiceDetailRepository{db: db}
}

// Create создает позицию счета
func (r *PostgresInvoiceDetailRepository) Create(ctx context.Context, detail *domain.InvoiceDetail) error {
	query := `
		INSERT INTO invoice_details (
			invoice_uuid, base_count, price, amount, amount_without_vat,
			amount_vat, amount_st, goods_name, tnved_code, gked_code,
			fcd_number, catalog_id, unit_type_id
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
		RETURNING id
	`
	return r.db.QueryRowContext(ctx, query,
		detail.InvoiceUUID, detail.BaseCount, detail.Price, detail.Amount,
		detail.AmountWithoutVAT, detail.AmountVAT, detail.AmountST, detail.GoodsName,
		detail.TNVEDCode, detail.GKEDCode, detail.FCDNumber, detail.CatalogID, detail.UnitTypeID,
	).Scan(&detail.ID)
}

// CreateBatch создает несколько позиций счета
func (r *PostgresInvoiceDetailRepository) CreateBatch(ctx context.Context, details []*domain.InvoiceDetail) error {
	query := `
		INSERT INTO invoice_details (
			invoice_uuid, base_count, price, amount, amount_without_vat,
			amount_vat, amount_st, goods_name, tnved_code, gked_code,
			fcd_number, catalog_id, unit_type_id
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
	`

	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.PreparexContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, detail := range details {
		_, err = stmt.ExecContext(ctx,
			detail.InvoiceUUID, detail.BaseCount, detail.Price, detail.Amount,
			detail.AmountWithoutVAT, detail.AmountVAT, detail.AmountST, detail.GoodsName,
			detail.TNVEDCode, detail.GKEDCode, detail.FCDNumber, detail.CatalogID, detail.UnitTypeID,
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

// GetByInvoiceUUID получает позиции счета по UUID счета
func (r *PostgresInvoiceDetailRepository) GetByInvoiceUUID(ctx context.Context, invoiceUUID string) ([]*domain.InvoiceDetail, error) {
	var details []*domain.InvoiceDetail
	query := `
		SELECT id, invoice_uuid, base_count, price, amount, amount_without_vat,
			   amount_vat, amount_st, goods_name, tnved_code, gked_code,
			   fcd_number, catalog_id, unit_type_id
		FROM invoice_details 
		WHERE invoice_uuid = $1
		ORDER BY id
	`

	err := r.db.SelectContext(ctx, &details, query, invoiceUUID)
	return details, err
}

// ListByInvoiceUUID получает позиции счета по UUID счета с пагинацией
func (r *PostgresInvoiceDetailRepository) ListByInvoiceUUID(ctx context.Context, invoiceUUID string, page, perPage int32) ([]*domain.InvoiceDetail, int32, error) {
	if page < 1 {
		page = 1
	}
	if perPage <= 0 || perPage > 100 {
		perPage = 20
	}

	var total int32
	countQuery := `SELECT COUNT(*) FROM invoice_details WHERE invoice_uuid = $1`
	if err := r.db.QueryRowContext(ctx, countQuery, invoiceUUID).Scan(&total); err != nil {
		return nil, 0, err
	}

	var details []*domain.InvoiceDetail
	listQuery := `
		SELECT id, invoice_uuid, base_count, price, amount, amount_without_vat,
			   amount_vat, amount_st, goods_name, tnved_code, gked_code,
			   fcd_number, catalog_id, unit_type_id
		FROM invoice_details
		WHERE invoice_uuid = $1
		ORDER BY id
		LIMIT $2 OFFSET $3
	`
	offset := (page - 1) * perPage
	if err := r.db.SelectContext(ctx, &details, listQuery, invoiceUUID, perPage, offset); err != nil {
		return nil, 0, err
	}

	return details, total, nil
}

// DeleteByInvoiceUUID удаляет все позиции счета
func (r *PostgresInvoiceDetailRepository) DeleteByInvoiceUUID(ctx context.Context, invoiceUUID string) error {
	query := `DELETE FROM invoice_details WHERE invoice_uuid = $1`
	_, err := r.db.ExecContext(ctx, query, invoiceUUID)
	return err
}
