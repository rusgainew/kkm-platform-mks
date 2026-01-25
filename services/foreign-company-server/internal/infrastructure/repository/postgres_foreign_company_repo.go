package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/rusgainew/kkm-project-mks/foreign-company-server/internal/domain"
)

// PostgresForeignCompanyRepository реализация репозитория иностранных компаний для PostgreSQL
type PostgresForeignCompanyRepository struct {
	db *sqlx.DB
}

// NewPostgresForeignCompanyRepository создает новый экземпляр репозитория
func NewPostgresForeignCompanyRepository(db *sqlx.DB) *PostgresForeignCompanyRepository {
	return &PostgresForeignCompanyRepository{db: db}
}

// Create создает новую иностранную компанию
func (r *PostgresForeignCompanyRepository) Create(ctx context.Context, company *domain.ForeignCompany) error {
	query := `
		INSERT INTO foreign_companies (
			pin, full_name, country_code, address, is_active,
			created_by, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id
	`
	err := r.db.QueryRowContext(ctx, query,
		company.PIN, company.FullName, company.CountryCode, company.Address,
		company.IsActive, company.CreatedBy, company.CreatedAt, company.UpdatedAt,
	).Scan(&company.ID)
	return err
}

// GetByID получает компанию по ID
func (r *PostgresForeignCompanyRepository) GetByID(ctx context.Context, id int64) (*domain.ForeignCompany, error) {
	var company domain.ForeignCompany
	query := `
		SELECT id, pin, full_name, country_code, address, is_active,
			   created_by, updated_by, created_at, updated_at
		FROM foreign_companies WHERE id = $1
	`

	err := r.db.GetContext(ctx, &company, query, id)
	if err == sql.ErrNoRows {
		return nil, domain.ErrForeignCompanyNotFound
	}
	if err != nil {
		return nil, err
	}

	return &company, nil
}

// Update обновляет иностранную компанию
func (r *PostgresForeignCompanyRepository) Update(ctx context.Context, company *domain.ForeignCompany) error {
	query := `
		UPDATE foreign_companies
		SET pin = $2, full_name = $3, country_code = $4, address = $5,
		    is_active = $6, updated_by = $7, updated_at = $8
		WHERE id = $1
	`
	result, err := r.db.ExecContext(ctx, query,
		company.ID, company.PIN, company.FullName, company.CountryCode,
		company.Address, company.IsActive, company.UpdatedBy, company.UpdatedAt,
	)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return domain.ErrForeignCompanyNotFound
	}

	return nil
}

// Delete удаляет иностранную компанию
func (r *PostgresForeignCompanyRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM foreign_companies WHERE id = $1`
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return domain.ErrForeignCompanyNotFound
	}

	return nil
}

// ExistsByPIN проверяет существование компании по PIN
func (r *PostgresForeignCompanyRepository) ExistsByPIN(ctx context.Context, pin string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM foreign_companies WHERE pin = $1)`
	err := r.db.GetContext(ctx, &exists, query, pin)
	return exists, err
}

// ListAll возвращает список всех компаний с пагинацией
func (r *PostgresForeignCompanyRepository) ListAll(ctx context.Context, limit, offset int) ([]*domain.ForeignCompany, int, error) {
	var companies []*domain.ForeignCompany
	query := `
		SELECT id, pin, full_name, country_code, address, is_active,
			   created_by, updated_by, created_at, updated_at
		FROM foreign_companies
		WHERE is_active = true
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	err := r.db.SelectContext(ctx, &companies, query, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	// Получение общего количества
	var total int
	countQuery := `SELECT COUNT(*) FROM foreign_companies WHERE is_active = true`
	err = r.db.GetContext(ctx, &total, countQuery)
	if err != nil {
		return nil, 0, err
	}

	return companies, total, nil
}

// Search ищет компании по запросу
func (r *PostgresForeignCompanyRepository) Search(ctx context.Context, query string, limit, offset int) ([]*domain.ForeignCompany, int, error) {
	var companies []*domain.ForeignCompany
	searchQuery := `
		SELECT id, pin, full_name, country_code, address, is_active,
			   created_by, updated_by, created_at, updated_at
		FROM foreign_companies
		WHERE is_active = true
		  AND (full_name ILIKE $1 OR pin ILIKE $1)
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`
	searchPattern := fmt.Sprintf("%%%s%%", query)

	err := r.db.SelectContext(ctx, &companies, searchQuery, searchPattern, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	// Получение общего количества
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM foreign_companies
		WHERE is_active = true
		  AND (full_name ILIKE $1 OR pin ILIKE $1)
	`
	err = r.db.GetContext(ctx, &total, countQuery, searchPattern)
	if err != nil {
		return nil, 0, err
	}

	return companies, total, nil
}

// GetByPIN получает компанию по PIN
func (r *PostgresForeignCompanyRepository) GetByPIN(ctx context.Context, pin string) (*domain.ForeignCompany, error) {
	var company domain.ForeignCompany
	query := `
		SELECT id, pin, full_name, country_code, address, is_active,
			   created_by, updated_by, created_at, updated_at
		FROM foreign_companies WHERE pin = $1
	`

	err := r.db.GetContext(ctx, &company, query, pin)
	if err == sql.ErrNoRows {
		return nil, domain.ErrForeignCompanyNotFound
	}
	if err != nil {
		return nil, err
	}

	return &company, nil
}

// List возвращает список всех компаний
func (r *PostgresForeignCompanyRepository) List(ctx context.Context, limit, offset int) ([]*domain.ForeignCompany, int, error) {
	return r.ListAll(ctx, limit, offset)
}

// ListByCountry возвращает компании из определенной страны
func (r *PostgresForeignCompanyRepository) ListByCountry(ctx context.Context, countryCode string, limit, offset int) ([]*domain.ForeignCompany, int, error) {
	var companies []*domain.ForeignCompany
	query := `
		SELECT id, pin, full_name, country_code, address, is_active,
			   created_by, updated_by, created_at, updated_at
		FROM foreign_companies
		WHERE is_active = true AND country_code = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	err := r.db.SelectContext(ctx, &companies, query, countryCode, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	// Получение общего количества
	var total int
	countQuery := `SELECT COUNT(*) FROM foreign_companies WHERE is_active = true AND country_code = $1`
	err = r.db.GetContext(ctx, &total, countQuery, countryCode)
	if err != nil {
		return nil, 0, err
	}

	return companies, total, nil
}
