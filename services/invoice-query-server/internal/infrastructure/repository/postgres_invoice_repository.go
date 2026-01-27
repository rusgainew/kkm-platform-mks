package repository

import (
	"context"
	"database/sql"
	"fmt"
	"math"
	"strings"

	"github.com/rusgainew/kkm-project-mks/invoice-query-server/internal/domain/ports"
	"github.com/rusgainew/kkm-project-mks/proto-lib/entities"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type PostgresInvoiceRepository struct {
	db     *sql.DB
	logger *zap.Logger
}

func NewPostgresInvoiceRepository(db *sql.DB, logger *zap.Logger) ports.InvoiceQueryRepository {
	return &PostgresInvoiceRepository{
		db:     db,
		logger: logger,
	}
}

func (r *PostgresInvoiceRepository) ListInvoices(ctx context.Context, page, size int32) ([]*entities.Invoice, int32, error) {
	// Исправлена пагинация: offset = (page - 1) * size
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * size

	// SQL запрос для получения счетов с пагинацией
	query := `
		SELECT document_uuid, invoice_number, invoice_date, 
		       total_amount, note, created_date
		FROM invoices
		ORDER BY created_date DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.QueryContext(ctx, query, size, offset)
	if err != nil {
		r.logger.Error("Failed to query invoices", zap.Error(err))
		return nil, 0, fmt.Errorf("query invoices: %w", err)
	}
	defer rows.Close()

	var invoices []*entities.Invoice
	for rows.Next() {
		invoice := &entities.Invoice{}
		var invoiceDate, createdDate sql.NullString

		err := rows.Scan(
			&invoice.DocumentUuid,
			&invoice.InvoiceNumber,
			&invoiceDate,
			&invoice.TotalAmount,
			&invoice.Note,
			&createdDate,
		)
		if err != nil {
			r.logger.Error("Failed to scan invoice row", zap.Error(err))
			return nil, 0, fmt.Errorf("scan invoice: %w", err)
		}

		// Конвертируем sql.NullString в wrapperspb.StringValue
		if invoiceDate.Valid {
			invoice.InvoiceDate = &wrapperspb.StringValue{Value: invoiceDate.String}
		}
		if createdDate.Valid {
			invoice.CreatedDate = &wrapperspb.StringValue{Value: createdDate.String}
		}

		invoices = append(invoices, invoice)
	}

	if err := rows.Err(); err != nil {
		r.logger.Error("Rows iteration error", zap.Error(err))
		return nil, 0, fmt.Errorf("rows error: %w", err)
	}

	// Получаем общее количество записей
	var totalCount int32
	countQuery := "SELECT COUNT(*) FROM invoices"
	err = r.db.QueryRowContext(ctx, countQuery).Scan(&totalCount)
	if err != nil {
		r.logger.Error("Failed to count invoices", zap.Error(err))
		return nil, 0, fmt.Errorf("count invoices: %w", err)
	}

	r.logger.Info("Listed invoices",
		zap.Int32("page", page),
		zap.Int32("size", size),
		zap.Int("returned", len(invoices)),
		zap.Int32("total", totalCount),
	)

	return invoices, totalCount, nil
}

func (r *PostgresInvoiceRepository) ListInvoiceDetails(ctx context.Context, page, size int32) ([]*entities.InvoiceDetail, int32, error) {
	// Исправлена пагинация
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * size

	query := `
		SELECT id, invoice_uuid, base_count, 
		       price, amount, amount_without_vat, 
		       amount_vat, amount_st, goods_name
		FROM invoice_details
		ORDER BY id DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.QueryContext(ctx, query, size, offset)
	if err != nil {
		r.logger.Error("Failed to query invoice details", zap.Error(err))
		return nil, 0, fmt.Errorf("query invoice details: %w", err)
	}
	defer rows.Close()

	var details []*entities.InvoiceDetail
	for rows.Next() {
		detail := &entities.InvoiceDetail{}

		err := rows.Scan(
			&detail.Id,
			&detail.InvoiceUuid,
			&detail.BaseCount,
			&detail.Price,
			&detail.Amount,
			&detail.AmountWithoutVat,
			&detail.AmountVat,
			&detail.AmountSt,
			&detail.GoodsName,
		)
		if err != nil {
			r.logger.Error("Failed to scan detail row", zap.Error(err))
			return nil, 0, fmt.Errorf("scan detail: %w", err)
		}

		details = append(details, detail)
	}

	if err := rows.Err(); err != nil {
		r.logger.Error("Rows iteration error", zap.Error(err))
		return nil, 0, fmt.Errorf("rows error: %w", err)
	}

	var totalCount int32
	countQuery := "SELECT COUNT(*) FROM invoice_details"
	err = r.db.QueryRowContext(ctx, countQuery).Scan(&totalCount)
	if err != nil {
		r.logger.Error("Failed to count details", zap.Error(err))
		return nil, 0, fmt.Errorf("count details: %w", err)
	}

	r.logger.Info("Listed invoice details",
		zap.Int32("page", page),
		zap.Int32("size", size),
		zap.Int("returned", len(details)),
		zap.Int32("total", totalCount),
	)

	return details, totalCount, nil
}

// ListInvoicesWithFilter возвращает список счетов-фактур с фильтрацией, сортировкой и пагинацией
func (r *PostgresInvoiceRepository) ListInvoicesWithFilter(ctx context.Context, filter *ports.InvoiceFilter, sort *ports.InvoiceSort, page, size int32) ([]*entities.Invoice, int32, error) {
	// Validate input parameters
	if err := ports.ValidateInvoiceFilter(filter); err != nil {
		return nil, 0, fmt.Errorf("invalid filter: %w", err)
	}

	if err := ports.ValidateInvoiceSort(sort); err != nil {
		return nil, 0, fmt.Errorf("invalid sort: %w", err)
	}

	var err error
	page, size, err = ports.ValidatePagination(page, size)
	if err != nil {
		return nil, 0, fmt.Errorf("invalid pagination: %w", err)
	}

	offset := (page - 1) * size

	// Строим WHERE clause динамически
	whereClause := ""
	args := []interface{}{}
	argIndex := 1

	if filter != nil {
		conditions := []string{}

		if filter.InvoiceNumber != "" {
			conditions = append(conditions, fmt.Sprintf("invoice_number = $%d", argIndex))
			args = append(args, filter.InvoiceNumber)
			argIndex++
		}

		if filter.DateFrom != "" {
			conditions = append(conditions, fmt.Sprintf("invoice_date >= $%d", argIndex))
			args = append(args, filter.DateFrom)
			argIndex++
		}

		if filter.DateTo != "" {
			conditions = append(conditions, fmt.Sprintf("invoice_date <= $%d", argIndex))
			args = append(args, filter.DateTo)
			argIndex++
		}

		if filter.MinAmount > 0 {
			conditions = append(conditions, fmt.Sprintf("total_amount >= $%d", argIndex))
			args = append(args, filter.MinAmount)
			argIndex++
		}

		if filter.MaxAmount > 0 {
			conditions = append(conditions, fmt.Sprintf("total_amount <= $%d", argIndex))
			args = append(args, filter.MaxAmount)
			argIndex++
		}

		if filter.SearchText != "" {
			conditions = append(conditions, fmt.Sprintf("(invoice_number ILIKE $%d OR note ILIKE $%d)", argIndex, argIndex+1))
			searchPattern := "%" + filter.SearchText + "%"
			args = append(args, searchPattern, searchPattern)
			argIndex += 2
		}

		if len(conditions) > 0 {
			whereClause = " WHERE " + strings.Join(conditions, " AND ")
		}
	}

	// Определяем ORDER BY
	orderBy := "created_date DESC" // По умолчанию
	if sort != nil && sort.Field != "" {
		validFields := map[string]bool{
			"invoice_date":   true,
			"total_amount":   true,
			"invoice_number": true,
			"created_date":   true,
		}
		if validFields[sort.Field] {
			// Валидация sort.Order для предотвращения SQL injection
			sortOrder := "ASC"
			if sort.Order == "DESC" || sort.Order == "desc" {
				sortOrder = "DESC"
			}
			orderBy = sort.Field + " " + sortOrder
		}
	}

	// SQL запрос с фильтрацией и сортировкой
	query := fmt.Sprintf(`
		SELECT document_uuid, invoice_number, invoice_date, 
		       total_amount, note, created_date
		FROM invoices
		%s
		ORDER BY %s
		LIMIT $%d OFFSET $%d
	`, whereClause, orderBy, argIndex, argIndex+1)

	args = append(args, size, offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		r.logger.Error("Failed to query invoices with filter", zap.Error(err))
		return nil, 0, fmt.Errorf("query invoices with filter: %w", err)
	}
	defer rows.Close()

	var invoices []*entities.Invoice
	for rows.Next() {
		invoice := &entities.Invoice{}
		var invoiceDate, createdDate sql.NullString
		var note sql.NullString

		err := rows.Scan(
			&invoice.DocumentUuid,
			&invoice.InvoiceNumber,
			&invoiceDate,
			&invoice.TotalAmount,
			&note,
			&createdDate,
		)
		if err != nil {
			r.logger.Error("Failed to scan invoice", zap.Error(err))
			return nil, 0, fmt.Errorf("scan invoice: %w", err)
		}

		if invoiceDate.Valid {
			invoice.InvoiceDate = wrapperspb.String(invoiceDate.String)
		}
		if note.Valid {
			invoice.Note = note.String
		}
		if createdDate.Valid {
			invoice.CreatedDate = wrapperspb.String(createdDate.String)
		}

		invoices = append(invoices, invoice)
	}

	// Получаем общее количество записей с учетом фильтров
	countQuery := "SELECT COUNT(*) FROM invoices" + whereClause
	var totalCount int32
	countArgs := args[:len(args)-2] // Убираем LIMIT и OFFSET
	err = r.db.QueryRowContext(ctx, countQuery, countArgs...).Scan(&totalCount)
	if err != nil {
		r.logger.Error("Failed to count invoices", zap.Error(err))
		return nil, 0, fmt.Errorf("count invoices: %w", err)
	}

	r.logger.Info("Listed invoices with filter",
		zap.Int32("count", int32(len(invoices))),
		zap.Int32("total", totalCount),
		zap.Int32("page", page),
		zap.Int32("size", size))

	return invoices, totalCount, nil
}

// SearchInvoices выполняет полнотекстовый поиск по счетам-фактурам
func (r *PostgresInvoiceRepository) SearchInvoices(ctx context.Context, searchText string, page, size int32) ([]*entities.Invoice, int32, error) {
	filter := &ports.InvoiceFilter{
		SearchText: searchText,
	}
	return r.ListInvoicesWithFilter(ctx, filter, nil, page, size)
}

// GetInvoiceByNumber возвращает счет-фактуру по номеру
func (r *PostgresInvoiceRepository) GetInvoiceByNumber(ctx context.Context, invoiceNumber string) (*entities.Invoice, error) {
	query := `
		SELECT document_uuid, invoice_number, invoice_date, 
		       total_amount, note, created_date
		FROM invoices
		WHERE invoice_number = $1
		LIMIT 1
	`

	invoice := &entities.Invoice{}
	var invoiceDate, createdDate sql.NullString
	var note sql.NullString

	err := r.db.QueryRowContext(ctx, query, invoiceNumber).Scan(
		&invoice.DocumentUuid,
		&invoice.InvoiceNumber,
		&invoiceDate,
		&invoice.TotalAmount,
		&note,
		&createdDate,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("invoice not found: %s", invoiceNumber)
	}
	if err != nil {
		r.logger.Error("Failed to get invoice by number", zap.Error(err), zap.String("invoice_number", invoiceNumber))
		return nil, fmt.Errorf("get invoice by number: %w", err)
	}

	if invoiceDate.Valid {
		invoice.InvoiceDate = wrapperspb.String(invoiceDate.String)
	}
	if note.Valid {
		invoice.Note = note.String
	}
	if createdDate.Valid {
		invoice.CreatedDate = wrapperspb.String(createdDate.String)
	}

	r.logger.Info("Got invoice by number", zap.String("invoice_number", invoiceNumber))

	return invoice, nil
}

// GetInvoicesByDateRange возвращает счета-фактуры за указанный период
func (r *PostgresInvoiceRepository) GetInvoicesByDateRange(ctx context.Context, dateFrom, dateTo string, page, size int32) ([]*entities.Invoice, int32, error) {
	filter := &ports.InvoiceFilter{
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}
	sort := &ports.InvoiceSort{
		Field: "invoice_date",
		Order: ports.SortOrderDesc,
	}
	return r.ListInvoicesWithFilter(ctx, filter, sort, page, size)
}

// calculateTotalPages вычисляет общее количество страниц
func calculateTotalPages(totalElements, pageSize int32) int32 {
	if pageSize == 0 {
		return 0
	}
	return int32(math.Ceil(float64(totalElements) / float64(pageSize)))
}
