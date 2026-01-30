// Файл document-server/internal/infrastructure/repository/postgres_document_repo.go содержит реализацию пакета repository.
package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type PostgresDocumentRepository struct {
	db *sqlx.DB
}

// NewPostgresDocumentRepository создает новый экземпляр репозитория
func NewPostgresDocumentRepository(db *sqlx.DB) *PostgresDocumentRepository {
	return &PostgresDocumentRepository{db: db}
}

// Create создает новый документ
func (r *PostgresDocumentRepository) Create(ctx context.Context, id, organizationID, title, content, createdBy string, createdAt int64) error {
	query := `
		INSERT INTO documents (id, organization_id, title, content, status, created_by, created_at, updated_at, status_changed_at, version)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $7, $7, $8)
	`
	_, err := r.db.ExecContext(ctx, query, id, organizationID, title, content, "draft", createdBy, createdAt, 1)
	return err
}

// Get получает документ по ID
func (r *PostgresDocumentRepository) Get(ctx context.Context, documentID string) (map[string]interface{}, error) {
	var id, organizationID, title, content, status, createdBy string
	var createdAt, updatedAt, statusChangedAt int64
	var version int32

	query := `
		SELECT id, organization_id, title, content, status, created_by, created_at, updated_at, status_changed_at, version
		FROM documents
		WHERE id = $1
	`
	err := r.db.QueryRowContext(ctx, query, documentID).Scan(
		&id, &organizationID, &title, &content, &status, &createdBy, &createdAt, &updatedAt, &statusChangedAt, &version,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("document not found")
		}
		return nil, err
	}

	return map[string]interface{}{
		"id":                id,
		"organization_id":   organizationID,
		"title":             title,
		"content":           content,
		"status":            status,
		"created_by":        createdBy,
		"created_at":        createdAt,
		"updated_at":        updatedAt,
		"status_changed_at": statusChangedAt,
		"version":           version,
	}, nil
}

// Update обновляет документ
func (r *PostgresDocumentRepository) Update(ctx context.Context, documentID, title, content string, updatedAt int64) error {
	query := `
		UPDATE documents
		SET title = $1, content = $2, updated_at = $3, version = version + 1
		WHERE id = $4
	`
	result, err := r.db.ExecContext(ctx, query, title, content, updatedAt, documentID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("document not found")
	}

	return nil
}

// UpdateStatus обновляет статус документа
func (r *PostgresDocumentRepository) UpdateStatus(ctx context.Context, documentID, status string, statusChangedAt int64) error {
	query := `
		UPDATE documents
		SET status = $1, status_changed_at = $2, updated_at = $2, version = version + 1
		WHERE id = $3
	`
	result, err := r.db.ExecContext(ctx, query, status, statusChangedAt, documentID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("document not found")
	}

	return nil
}

// List получает список документов по организации с пагинацией
func (r *PostgresDocumentRepository) List(ctx context.Context, organizationID string, status string, page, perPage int) ([]map[string]interface{}, int64, error) {
	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 20
	}

	offset := (page - 1) * perPage

	// Получаем общее количество документов
	var total int64
	countQuery := `SELECT COUNT(*) FROM documents WHERE organization_id = $1`
	countArgs := []interface{}{organizationID}

	if status != "" {
		countQuery += ` AND status = $2`
		countArgs = append(countArgs, status)
	}

	err := r.db.QueryRowContext(ctx, countQuery, countArgs...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Получаем документы - используем parameterized queries для безопасности
	listQuery := `
		SELECT id, organization_id, title, content, status, created_by, created_at, updated_at, status_changed_at, version
		FROM documents
		WHERE organization_id = $1
	`
	listArgs := []interface{}{organizationID}
	paramIndex := 2

	if status != "" {
		listQuery += ` AND status = $` + fmt.Sprint(paramIndex)
		listArgs = append(listArgs, status)
		paramIndex++
	}

	// LIMIT и OFFSET используют параметризованные placeholders
	listQuery += ` ORDER BY created_at DESC LIMIT $` + fmt.Sprint(paramIndex) + ` OFFSET $` + fmt.Sprint(paramIndex+1)
	listArgs = append(listArgs, perPage, offset)

	rows, err := r.db.QueryContext(ctx, listQuery, listArgs...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var documents []map[string]interface{}
	for rows.Next() {
		var id, organizationID, title, content, statusVal, createdBy string
		var createdAt, updatedAt, statusChangedAt int64
		var version int32

		if err := rows.Scan(&id, &organizationID, &title, &content, &statusVal, &createdBy, &createdAt, &updatedAt, &statusChangedAt, &version); err != nil {
			return nil, 0, err
		}

		documents = append(documents, map[string]interface{}{
			"id":                id,
			"organization_id":   organizationID,
			"title":             title,
			"content":           content,
			"status":            statusVal,
			"created_by":        createdBy,
			"created_at":        createdAt,
			"updated_at":        updatedAt,
			"status_changed_at": statusChangedAt,
			"version":           version,
		})
	}

	if err = rows.Err(); err != nil {
		return nil, 0, err
	}

	return documents, total, nil
}

// Delete удаляет документ
func (r *PostgresDocumentRepository) Delete(ctx context.Context, documentID string) error {
	query := `DELETE FROM documents WHERE id = $1`
	result, err := r.db.ExecContext(ctx, query, documentID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("document not found")
	}

	return nil
}
