// Файл company-server/internal/infrastructure/repository/postgres_employee_repo.go содержит реализацию пакета repository.
package repository

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
	"github.com/rusgainew/kkm-project-mks/company-server/internal/domain"
)

// PostgresEmployeeRepository реализация репозитория участников для PostgreSQL
type PostgresEmployeeRepository struct {
	db *sqlx.DB
}

// NewPostgresEmployeeRepository создает новый экземпляр репозитория
func NewPostgresEmployeeRepository(db *sqlx.DB) *PostgresEmployeeRepository {
	return &PostgresEmployeeRepository{db: db}
}

// Create добавляет участника в организацию
func (r *PostgresEmployeeRepository) Create(ctx context.Context, emp *domain.Employee) error {
	query := `
		INSERT INTO employees (id, organization_id, user_id, role, joined_at)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err := r.db.ExecContext(ctx, query, emp.ID, emp.OrganizationID, emp.UserID, emp.Role, emp.JoinedAt)
	return err
}

// GetByID получает участника по ID
func (r *PostgresEmployeeRepository) GetByID(ctx context.Context, id string) (*domain.Employee, error) {
	var emp domain.Employee
	query := `SELECT id, organization_id, user_id, role, joined_at FROM employees WHERE id = $1`

	err := r.db.GetContext(ctx, &emp, query, id)
	if err == sql.ErrNoRows {
		return nil, domain.ErrEmployeeNotFound
	}
	if err != nil {
		return nil, err
	}

	return &emp, nil
}

// GetByUserAndOrganization получает участника по user_id и organization_id
func (r *PostgresEmployeeRepository) GetByUserAndOrganization(ctx context.Context, userID, organizationID string) (*domain.Employee, error) {
	var emp domain.Employee
	query := `SELECT id, organization_id, user_id, role, joined_at FROM employees WHERE user_id = $1 AND organization_id = $2`

	err := r.db.GetContext(ctx, &emp, query, userID, organizationID)
	if err == sql.ErrNoRows {
		return nil, domain.ErrEmployeeNotFound
	}
	if err != nil {
		return nil, err
	}

	return &emp, nil
}

// Delete удаляет участника
func (r *PostgresEmployeeRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM employees WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return domain.ErrEmployeeNotFound
	}

	return nil
}

// ListByOrganization возвращает список участников организации
func (r *PostgresEmployeeRepository) ListByOrganization(ctx context.Context, organizationID string, page, perPage int32) ([]*domain.Employee, int32, error) {
	var employees []*domain.Employee
	var total int32

	offset := (page - 1) * perPage

	// Получение общего количества
	countQuery := `SELECT COUNT(*) FROM employees WHERE organization_id = $1`
	err := r.db.GetContext(ctx, &total, countQuery, organizationID)
	if err != nil {
		return nil, 0, err
	}

	// Получение списка
	selectQuery := `
		SELECT id, organization_id, user_id, role, joined_at 
		FROM employees 
		WHERE organization_id = $1
		ORDER BY joined_at DESC
		LIMIT $2 OFFSET $3
	`
	err = r.db.SelectContext(ctx, &employees, selectQuery, organizationID, perPage, offset)
	if err != nil {
		return nil, 0, err
	}

	return employees, total, nil
}

// ExistsInOrganization проверяет, есть ли пользователь в организации
func (r *PostgresEmployeeRepository) ExistsInOrganization(ctx context.Context, userID, organizationID string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM employees WHERE user_id = $1 AND organization_id = $2)`
	err := r.db.GetContext(ctx, &exists, query, userID, organizationID)
	return exists, err
}
