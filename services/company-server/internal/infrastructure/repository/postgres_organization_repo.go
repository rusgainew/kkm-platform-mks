// Файл company-server/internal/infrastructure/repository/postgres_organization_repo.go содержит реализацию пакета repository.
package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/rusgainew/kkm-project-mks/company-server/internal/domain"
)

// PostgresOrganizationRepository реализация репозитория организаций для PostgreSQL
type PostgresOrganizationRepository struct {
	db *sqlx.DB
}

// NewPostgresOrganizationRepository создает новый экземпляр репозитория
func NewPostgresOrganizationRepository(db *sqlx.DB) *PostgresOrganizationRepository {
	return &PostgresOrganizationRepository{db: db}
}

// Create создает новую организацию
func (r *PostgresOrganizationRepository) Create(ctx context.Context, org *domain.Organization) error {
	query := `
		INSERT INTO organizations (id, name, description, owner_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := r.db.ExecContext(ctx, query, org.ID, org.Name, org.Description, org.OwnerID, org.CreatedAt, org.UpdatedAt)
	return err
}

// GetByID получает организацию по ID
func (r *PostgresOrganizationRepository) GetByID(ctx context.Context, id string) (*domain.Organization, error) {
	var org domain.Organization
	query := `SELECT id, name, description, owner_id, created_at, updated_at FROM organizations WHERE id = $1`

	err := r.db.GetContext(ctx, &org, query, id)
	if err == sql.ErrNoRows {
		return nil, domain.ErrOrganizationNotFound
	}
	if err != nil {
		return nil, err
	}

	return &org, nil
}

// Update обновляет организацию
func (r *PostgresOrganizationRepository) Update(ctx context.Context, org *domain.Organization) error {
	query := `
		UPDATE organizations 
		SET name = $2, description = $3, updated_at = $4
		WHERE id = $1
	`
	result, err := r.db.ExecContext(ctx, query, org.ID, org.Name, org.Description, org.UpdatedAt)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return domain.ErrOrganizationNotFound
	}

	return nil
}

// Delete удаляет организацию
func (r *PostgresOrganizationRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM organizations WHERE id = $1`
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return domain.ErrOrganizationNotFound
	}

	return nil
}

// List возвращает список организаций с пагинацией
func (r *PostgresOrganizationRepository) List(ctx context.Context, page, perPage int32, ownerID string) ([]*domain.Organization, int32, error) {
	var orgs []*domain.Organization
	var total int32

	offset := (page - 1) * perPage

	// Построение запроса с фильтром
	baseQuery := `FROM organizations`
	whereClause := ""
	args := []interface{}{}
	argPosition := 1

	if ownerID != "" {
		whereClause = fmt.Sprintf(" WHERE owner_id = $%d", argPosition)
		args = append(args, ownerID)
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
	selectQuery := `SELECT id, name, description, owner_id, created_at, updated_at ` +
		baseQuery + whereClause +
		fmt.Sprintf(` ORDER BY created_at DESC LIMIT $%d OFFSET $%d`, argPosition, argPosition+1)

	err = r.db.SelectContext(ctx, &orgs, selectQuery, args...)
	if err != nil {
		return nil, 0, err
	}

	return orgs, total, nil
}

// ExistsByName проверяет существование организации по имени
func (r *PostgresOrganizationRepository) ExistsByName(ctx context.Context, name string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM organizations WHERE name = $1)`
	err := r.db.GetContext(ctx, &exists, query, name)
	return exists, err
}
