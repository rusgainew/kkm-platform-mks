package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/rusgainew/kkm-project-mks/catalog-query-server/internal/domain/ports"
	"github.com/rusgainew/kkm-project-mks/proto-lib/dictionaries"
)

// PostgresCatalogRepository реализует репозиторий для запросов к каталогу
type PostgresCatalogRepository struct {
	db *sql.DB
}

// NewPostgresCatalogRepository создает новый PostgresCatalogRepository
func NewPostgresCatalogRepository(db *sql.DB) ports.CatalogQueryRepository {
	return &PostgresCatalogRepository{db: db}
}

// ListCatalogs возвращает список элементов каталога с пагинацией
func (r *PostgresCatalogRepository) ListCatalogs(ctx context.Context, page, size int32) ([]*dictionaries.Catalog, int32, error) {
	// Валидация параметров пагинации
	if page < 1 {
		page = 1
	}
	if size < 1 || size > 100 {
		size = 10
	}

	offset := (page - 1) * size

	// Запрос данных с пагинацией из таблицы catalog_items
	query := `
		SELECT code, name, tnved, gked
		FROM catalog_items
		WHERE is_active = true
		ORDER BY code ASC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.QueryContext(ctx, query, size, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query catalogs: %w", err)
	}
	defer rows.Close()

	var catalogs []*dictionaries.Catalog
	for rows.Next() {
		var number, name, tnvedCode, gkedCode sql.NullString

		err := rows.Scan(&number, &name, &tnvedCode, &gkedCode)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan catalog: %w", err)
		}

		catalog := &dictionaries.Catalog{
			Number:    number.String,
			Name:      name.String,
			TnvedCode: tnvedCode.String,
			GkedCode:  gkedCode.String,
		}

		catalogs = append(catalogs, catalog)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error iterating catalog rows: %w", err)
	}

	// Подсчет общего количества записей
	var totalCount int32
	countQuery := "SELECT COUNT(*) FROM catalog_items WHERE is_active = true"
	err = r.db.QueryRowContext(ctx, countQuery).Scan(&totalCount)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count catalogs: %w", err)
	}

	return catalogs, totalCount, nil
}

// ListCatalogsWithFilter возвращает список элементов каталога с фильтрацией и сортировкой
func (r *PostgresCatalogRepository) ListCatalogsWithFilter(ctx context.Context, filter *ports.CatalogFilter, sort *ports.CatalogSort, page, size int32) ([]*dictionaries.Catalog, int32, error) {
	if page < 1 {
		page = 1
	}
	if size < 1 || size > 100 {
		size = 10
	}
	offset := (page - 1) * size

	// Построение WHERE clause
	whereClause := " WHERE is_active = true"
	args := []interface{}{}
	argIndex := 1

	if filter != nil {
		if filter.Name != "" {
			whereClause += fmt.Sprintf(" AND name ILIKE $%d", argIndex)
			args = append(args, "%"+filter.Name+"%")
			argIndex++
		}

		if filter.Number != "" {
			whereClause += fmt.Sprintf(" AND code = $%d", argIndex)
			args = append(args, filter.Number)
			argIndex++
		}

		if filter.TnvedCode != "" {
			whereClause += fmt.Sprintf(" AND tnved = $%d", argIndex)
			args = append(args, filter.TnvedCode)
			argIndex++
		}

		if filter.GkedCode != "" {
			whereClause += fmt.Sprintf(" AND gked = $%d", argIndex)
			args = append(args, filter.GkedCode)
			argIndex++
		}

		if filter.SearchText != "" {
			whereClause += fmt.Sprintf(" AND (name ILIKE $%d OR code ILIKE $%d)", argIndex, argIndex+1)
			searchPattern := "%" + filter.SearchText + "%"
			args = append(args, searchPattern, searchPattern)
			argIndex += 2
		}
	}

	// Построение ORDER BY clause с валидацией
	orderBy := "code ASC"
	if sort != nil && sort.Field != "" {
		validFields := map[string]string{
			"name":       "name",
			"number":     "code",
			"tnved_code": "tnved",
			"gked_code":  "gked",
		}
		if mappedField, ok := validFields[sort.Field]; ok {
			orderBy = mappedField + " " + string(sort.Order)
		}
	}

	// Запрос данных
	query := fmt.Sprintf(`
		SELECT code, name, tnved, gked
		FROM catalog_items
		%s
		ORDER BY %s
		LIMIT $%d OFFSET $%d
	`, whereClause, orderBy, argIndex, argIndex+1)

	args = append(args, size, offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query catalogs with filter: %w", err)
	}
	defer rows.Close()

	var catalogs []*dictionaries.Catalog
	for rows.Next() {
		var number, name, tnvedCode, gkedCode sql.NullString

		err := rows.Scan(&number, &name, &tnvedCode, &gkedCode)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan catalog: %w", err)
		}

		catalog := &dictionaries.Catalog{
			Number:    number.String,
			Name:      name.String,
			TnvedCode: tnvedCode.String,
			GkedCode:  gkedCode.String,
		}

		catalogs = append(catalogs, catalog)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error iterating catalog rows: %w", err)
	}

	// Подсчет общего количества с теми же фильтрами
	countQuery := "SELECT COUNT(*) FROM catalog_items" + whereClause
	countArgs := args[:len(args)-2] // Убираем LIMIT и OFFSET

	var totalCount int32
	err = r.db.QueryRowContext(ctx, countQuery, countArgs...).Scan(&totalCount)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count catalogs with filter: %w", err)
	}

	return catalogs, totalCount, nil
}

// SearchCatalogs выполняет полнотекстовый поиск по каталогу
func (r *PostgresCatalogRepository) SearchCatalogs(ctx context.Context, searchText string, page, size int32) ([]*dictionaries.Catalog, int32, error) {
	filter := &ports.CatalogFilter{
		SearchText: searchText,
	}
	sort := &ports.CatalogSort{
		Field: "name",
		Order: ports.SortOrderAsc,
	}
	return r.ListCatalogsWithFilter(ctx, filter, sort, page, size)
}

// GetCatalogByNumber возвращает элемент каталога по коду
func (r *PostgresCatalogRepository) GetCatalogByNumber(ctx context.Context, number string) (*dictionaries.Catalog, error) {
	query := `
		SELECT code, name, tnved, gked
		FROM catalog_items
		WHERE code = $1 AND is_active = true
		LIMIT 1
	`

	var catalogNumber, name, tnvedCode, gkedCode sql.NullString

	err := r.db.QueryRowContext(ctx, query, number).Scan(&catalogNumber, &name, &tnvedCode, &gkedCode)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("catalog not found with number: %s", number)
		}
		return nil, fmt.Errorf("failed to get catalog by number: %w", err)
	}

	catalog := &dictionaries.Catalog{
		Number:    catalogNumber.String,
		Name:      name.String,
		TnvedCode: tnvedCode.String,
		GkedCode:  gkedCode.String,
	}

	return catalog, nil
}

// GetCatalogsByTnvedCode возвращает элементы каталога по коду ТН ВЭД
func (r *PostgresCatalogRepository) GetCatalogsByTnvedCode(ctx context.Context, tnvedCode string, page, size int32) ([]*dictionaries.Catalog, int32, error) {
	filter := &ports.CatalogFilter{
		TnvedCode: tnvedCode,
	}
	sort := &ports.CatalogSort{
		Field: "name",
		Order: ports.SortOrderAsc,
	}
	return r.ListCatalogsWithFilter(ctx, filter, sort, page, size)
}
