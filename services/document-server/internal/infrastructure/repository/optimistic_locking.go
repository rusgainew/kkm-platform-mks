package repository

import (
	"context"
	"database/sql"
	"fmt"
)

// GetWithVersion получает документ с информацией о версии (для optimistic locking)
func (r *PostgresDocumentRepository) GetWithVersion(ctx context.Context, documentID string) (map[string]interface{}, error) {
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
			return nil, nil
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

// UpdateWithVersion обновляет документ проверяя версию (optimistic locking)
// Возвращает ошибку если версия не совпадает
func (r *PostgresDocumentRepository) UpdateWithVersion(ctx context.Context, documentID, title, content string, expectedVersion int) error {
	query := `
		UPDATE documents
		SET title = $1, content = $2, updated_at = EXTRACT(EPOCH FROM NOW())::BIGINT, version = version + 1
		WHERE id = $3 AND version = $4
	`
	result, err := r.db.ExecContext(ctx, query, title, content, documentID, expectedVersion)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		// Проверяем существует ли документ
		var exists bool
		checkQuery := `SELECT EXISTS(SELECT 1 FROM documents WHERE id = $1)`
		err := r.db.QueryRowContext(ctx, checkQuery, documentID).Scan(&exists)
		if err != nil {
			return err
		}

		if !exists {
			return fmt.Errorf("document not found")
		}

		// Документ существует но версия не совпадает
		return fmt.Errorf("version conflict: document has been modified by another request")
	}

	return nil
}
