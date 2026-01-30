// Файл foreign-company-query-server/internal/infrastructure/repository/postgres_foreign_company_repository.go содержит реализацию пакета repository.
package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/rusgainew/kkm-project-mks/foreign-company-query-server/internal/domain/ports"
	"github.com/rusgainew/kkm-project-mks/proto-lib/dictionaries"
	"go.uber.org/zap"
)

type PostgresForeignCompanyRepository struct {
	db     *sql.DB
	logger *zap.Logger
}

func NewPostgresForeignCompanyRepository(db *sql.DB, logger *zap.Logger) ports.ForeignCompanyQueryRepository {
	return &PostgresForeignCompanyRepository{
		db:     db,
		logger: logger,
	}
}

func (r *PostgresForeignCompanyRepository) ListForeignCompanies(ctx context.Context, page, size int32) ([]*dictionaries.ForeignCompany, int32, error) {
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * size

	query := `
		SELECT id, pin, full_name
		FROM foreign_companies
		WHERE is_active = true
		ORDER BY full_name
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.QueryContext(ctx, query, size, offset)
	if err != nil {
		r.logger.Error("Failed to query foreign companies", zap.Error(err))
		return nil, 0, fmt.Errorf("query foreign companies: %w", err)
	}
	defer rows.Close()

	var companies []*dictionaries.ForeignCompany
	for rows.Next() {
		company := &dictionaries.ForeignCompany{}

		err := rows.Scan(
			&company.Id,
			&company.Pin,
			&company.FullName,
		)
		if err != nil {
			r.logger.Error("Failed to scan foreign company row", zap.Error(err))
			return nil, 0, fmt.Errorf("scan foreign company: %w", err)
		}

		companies = append(companies, company)
	}

	if err := rows.Err(); err != nil {
		r.logger.Error("Rows iteration error", zap.Error(err))
		return nil, 0, fmt.Errorf("rows error: %w", err)
	}

	var totalCount int32
	countQuery := "SELECT COUNT(*) FROM foreign_companies WHERE is_active = true"
	err = r.db.QueryRowContext(ctx, countQuery).Scan(&totalCount)
	if err != nil {
		r.logger.Error("Failed to count foreign companies", zap.Error(err))
		return nil, 0, fmt.Errorf("count foreign companies: %w", err)
	}

	r.logger.Info("Listed foreign companies",
		zap.Int32("page", page),
		zap.Int32("size", size),
		zap.Int("returned", len(companies)),
		zap.Int32("total", totalCount),
	)

	return companies, totalCount, nil
}

// ListForeignCompaniesWithFilter возвращает список иностранных компаний с фильтрацией и сортировкой
func (r *PostgresForeignCompanyRepository) ListForeignCompaniesWithFilter(ctx context.Context, filter *ports.ForeignCompanyFilter, sort *ports.ForeignCompanySort, page, size int32) ([]*dictionaries.ForeignCompany, int32, error) {
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * size

	// Строим WHERE clause динамически
	whereClause := " WHERE is_active = true"
	args := []interface{}{}
	argIndex := 1

	if filter != nil {
		if filter.PIN != "" {
			whereClause += fmt.Sprintf(" AND pin = $%d", argIndex)
			args = append(args, filter.PIN)
			argIndex++
		}

		if filter.FullName != "" {
			whereClause += fmt.Sprintf(" AND full_name ILIKE $%d", argIndex)
			args = append(args, "%"+filter.FullName+"%")
			argIndex++
		}

		if filter.CountryCode != "" {
			whereClause += fmt.Sprintf(" AND country_code = $%d", argIndex)
			args = append(args, filter.CountryCode)
			argIndex++
		}

		if filter.SearchText != "" {
			whereClause += fmt.Sprintf(" AND (pin ILIKE $%d OR full_name ILIKE $%d)", argIndex, argIndex+1)
			searchPattern := "%" + filter.SearchText + "%"
			args = append(args, searchPattern, searchPattern)
			argIndex += 2
		}
	}

	// Определяем ORDER BY
	orderBy := "full_name ASC" // По умолчанию
	if sort != nil && sort.Field != "" {
		validFields := map[string]string{
			"pin":          "pin",
			"full_name":    "full_name",
			"country_code": "country_code",
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
		SELECT id, pin, full_name
		FROM foreign_companies
		` + whereClause + `
		ORDER BY ` + orderBy + `
		LIMIT $` + fmt.Sprintf("%d", argIndex) + ` OFFSET $` + fmt.Sprintf("%d", argIndex+1)

	args = append(args, size, offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		r.logger.Error("Failed to query foreign companies with filter", zap.Error(err))
		return nil, 0, fmt.Errorf("query foreign companies with filter: %w", err)
	}
	defer rows.Close()

	var companies []*dictionaries.ForeignCompany
	for rows.Next() {
		company := &dictionaries.ForeignCompany{}
		err := rows.Scan(
			&company.Id,
			&company.Pin,
			&company.FullName,
		)
		if err != nil {
			r.logger.Error("Failed to scan foreign company", zap.Error(err))
			return nil, 0, fmt.Errorf("scan foreign company: %w", err)
		}
		companies = append(companies, company)
	}

	// Получаем общее количество записей с учетом фильтров
	countQuery := "SELECT COUNT(*) FROM foreign_companies" + whereClause
	var totalCount int32
	countArgs := args[:len(args)-2] // Убираем LIMIT и OFFSET
	err = r.db.QueryRowContext(ctx, countQuery, countArgs...).Scan(&totalCount)
	if err != nil {
		r.logger.Error("Failed to count foreign companies", zap.Error(err))
		return nil, 0, fmt.Errorf("count foreign companies: %w", err)
	}

	r.logger.Info("Listed foreign companies with filter",
		zap.Int("count", len(companies)),
		zap.Int32("total", totalCount),
		zap.Int32("page", page),
		zap.Int32("size", size))

	return companies, totalCount, nil
}

// SearchForeignCompanies выполняет полнотекстовый поиск по иностранным компаниям
func (r *PostgresForeignCompanyRepository) SearchForeignCompanies(ctx context.Context, searchText string, page, size int32) ([]*dictionaries.ForeignCompany, int32, error) {
	filter := &ports.ForeignCompanyFilter{
		SearchText: searchText,
	}
	return r.ListForeignCompaniesWithFilter(ctx, filter, nil, page, size)
}

// GetForeignCompanyByPIN возвращает иностранную компанию по PIN
func (r *PostgresForeignCompanyRepository) GetForeignCompanyByPIN(ctx context.Context, pin string) (*dictionaries.ForeignCompany, error) {
	query := `
		SELECT id, pin, full_name
		FROM foreign_companies
		WHERE pin = $1 AND is_active = true
		LIMIT 1
	`

	company := &dictionaries.ForeignCompany{}
	err := r.db.QueryRowContext(ctx, query, pin).Scan(
		&company.Id,
		&company.Pin,
		&company.FullName,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("foreign company not found: %s", pin)
	}
	if err != nil {
		r.logger.Error("Failed to get foreign company by PIN", zap.Error(err), zap.String("pin", pin))
		return nil, fmt.Errorf("get foreign company by PIN: %w", err)
	}

	r.logger.Info("Got foreign company by PIN", zap.String("pin", pin))

	return company, nil
}

// GetForeignCompaniesByCountry возвращает компании по коду страны
func (r *PostgresForeignCompanyRepository) GetForeignCompaniesByCountry(ctx context.Context, countryCode string, page, size int32) ([]*dictionaries.ForeignCompany, int32, error) {
	filter := &ports.ForeignCompanyFilter{
		CountryCode: countryCode,
	}
	sort := &ports.ForeignCompanySort{
		Field: "full_name",
		Order: ports.SortOrderAsc,
	}
	return r.ListForeignCompaniesWithFilter(ctx, filter, sort, page, size)
}
