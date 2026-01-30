// Файл catalog-server/internal/infrastructure/repository/postgres_catalog_repo.go содержит реализацию пакета repository.
package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/rusgainew/kkm-project-mks/catalog-server/internal/domain"
)

// PostgresCatalogRepository реализация репозитория каталога для PostgreSQL
type PostgresCatalogRepository struct {
	db *sqlx.DB
}

// NewPostgresCatalogRepository создает новый экземпляр репозитория
func NewPostgresCatalogRepository(db *sqlx.DB) *PostgresCatalogRepository {
	return &PostgresCatalogRepository{db: db}
}

// Create создает новый элемент каталога
func (r *PostgresCatalogRepository) Create(ctx context.Context, item *domain.CatalogItem) error {
	query := `
		INSERT INTO catalog_items (
			id, organization_id, name, code, description, unit_type,
			price, vat_rate, created_by, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`
	_, err := r.db.ExecContext(ctx, query,
		item.ID, item.OrganizationID, item.Name, item.Code, item.Description,
		item.UnitType, item.Price, item.VATRate, item.CreatedBy,
		item.CreatedAt, item.UpdatedAt,
	)
	return err
}

// GetByID получает элемент по ID
func (r *PostgresCatalogRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.CatalogItem, error) {
	var item domain.CatalogItem
	query := `
		SELECT id, organization_id, name, code, description, unit_type,
			   price, vat_rate, created_by, created_at, updated_at
		FROM catalog_items WHERE id = $1
	`

	err := r.db.GetContext(ctx, &item, query, id)
	if err == sql.ErrNoRows {
		return nil, domain.ErrCatalogItemNotFound
	}
	if err != nil {
		return nil, err
	}

	return &item, nil
}

// GetByCode получает элемент по коду и организации
func (r *PostgresCatalogRepository) GetByCode(ctx context.Context, organizationID uuid.UUID, code string) (*domain.CatalogItem, error) {
	var item domain.CatalogItem
	query := `
		SELECT id, organization_id, name, code, description, unit_type,
			   price, vat_rate, created_by, created_at, updated_at
		FROM catalog_items WHERE organization_id = $1 AND code = $2
	`

	err := r.db.GetContext(ctx, &item, query, organizationID, code)
	if err == sql.ErrNoRows {
		return nil, domain.ErrCatalogItemNotFound
	}
	if err != nil {
		return nil, err
	}

	return &item, nil
}

// Update обновляет элемент каталога
func (r *PostgresCatalogRepository) Update(ctx context.Context, item *domain.CatalogItem) error {
	query := `
		UPDATE catalog_items
		SET name = $2, description = $3, unit_type = $4, price = $5,
		    vat_rate = $6, updated_at = $7
		WHERE id = $1
	`
	result, err := r.db.ExecContext(ctx, query,
		item.ID, item.Name, item.Description, item.UnitType,
		item.Price, item.VATRate, item.UpdatedAt,
	)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return domain.ErrCatalogItemNotFound
	}

	return nil
}

// Delete удаляет элемент каталога
func (r *PostgresCatalogRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM catalog_items WHERE id = $1`
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return domain.ErrCatalogItemNotFound
	}

	return nil
}

// ExistsByCode проверяет существование элемента по коду
func (r *PostgresCatalogRepository) ExistsByCode(ctx context.Context, organizationID uuid.UUID, code string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM catalog_items WHERE organization_id = $1 AND code = $2)`
	err := r.db.GetContext(ctx, &exists, query, organizationID, code)
	return exists, err
}

// ListByOrganization возвращает список элементов организации с пагинацией
func (r *PostgresCatalogRepository) ListByOrganization(ctx context.Context, organizationID uuid.UUID, limit, offset int) ([]*domain.CatalogItem, int, error) {
	var items []*domain.CatalogItem
	query := `
		SELECT id, organization_id, name, code, description, unit_type,
			   price, vat_rate, created_by, created_at, updated_at
		FROM catalog_items
		WHERE organization_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	err := r.db.SelectContext(ctx, &items, query, organizationID, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	// Получение общего количества
	var total int
	countQuery := `SELECT COUNT(*) FROM catalog_items WHERE organization_id = $1`
	err = r.db.GetContext(ctx, &total, countQuery, organizationID)
	if err != nil {
		return nil, 0, err
	}

	return items, total, nil
}

// Search ищет элементы по запросу
func (r *PostgresCatalogRepository) Search(ctx context.Context, organizationID uuid.UUID, query string, limit, offset int) ([]*domain.CatalogItem, int, error) {
	var items []*domain.CatalogItem
	searchQuery := `
		SELECT id, organization_id, name, code, description, unit_type,
			   price, vat_rate, created_by, created_at, updated_at
		FROM catalog_items
		WHERE organization_id = $1
		  AND (name ILIKE $2 OR code ILIKE $2 OR description ILIKE $2)
		ORDER BY created_at DESC
		LIMIT $3 OFFSET $4
	`
	searchPattern := fmt.Sprintf("%%%s%%", query)

	err := r.db.SelectContext(ctx, &items, searchQuery, organizationID, searchPattern, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	// Получение общего количества
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM catalog_items
		WHERE organization_id = $1
		  AND (name ILIKE $2 OR code ILIKE $2 OR description ILIKE $2)
	`
	err = r.db.GetContext(ctx, &total, countQuery, organizationID, searchPattern)
	if err != nil {
		return nil, 0, err
	}

	return items, total, nil
}
