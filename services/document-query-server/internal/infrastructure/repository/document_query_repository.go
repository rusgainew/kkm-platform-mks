package repository

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"

	pb "github.com/rusgainew/kkm-project-mks/proto-lib/api"
)

type DocumentQueryRepository struct {
	db     *sqlx.DB
	logger *zap.Logger
}

func NewDocumentQueryRepository(db *sqlx.DB, logger *zap.Logger) *DocumentQueryRepository {
	return &DocumentQueryRepository{
		db:     db,
		logger: logger,
	}
}

// GetDocument retrieves a single document by ID
func (r *DocumentQueryRepository) GetDocument(ctx context.Context, documentID string) (*pb.DocumentReadModel, error) {
	query := `
		SELECT id, document_number, title, description, document_type, status, company_id, company_name,
		       created_by_user_id, created_by_user_name, approval_status, approved_by, approved_at,
		       sent_at, rejected_by, rejected_at, rejection_reason, created_at, updated_at, archived_at
		FROM document_read_model
		WHERE id = $1 AND deleted_at IS NULL
	`

	var doc pb.DocumentReadModel
	err := r.db.QueryRowContext(ctx, query, documentID).Scan(
		&doc.Id,
		&doc.DocumentNumber,
		&doc.Title,
		&doc.Description,
		&doc.DocumentType,
		&doc.Status,
		&doc.CompanyId,
		&doc.CompanyName,
		&doc.CreatedByUserId,
		&doc.CreatedByUserName,
		&doc.ApprovalStatus,
		&doc.ApprovedBy,
		&doc.ApprovedAt,
		&doc.SentAt,
		&doc.RejectedBy,
		&doc.RejectedAt,
		&doc.RejectionReason,
		&doc.CreatedAt,
		&doc.UpdatedAt,
		&doc.ArchivedAt,
	)

	if err != nil {
		r.logger.Error("Failed to get document",
			zap.String("document_id", documentID),
			zap.Error(err))
		return nil, fmt.Errorf("failed to get document: %w", err)
	}

	return &doc, nil
}

// ListDocuments retrieves documents with pagination and optional filters
func (r *DocumentQueryRepository) ListDocuments(ctx context.Context, offset, limit int32, status, docType, companyID, approvalStatus string) ([]*pb.DocumentReadModel, int64, error) {
	query := `
		SELECT id, document_number, title, description, document_type, status, company_id, company_name,
		       created_by_user_id, created_by_user_name, approval_status, approved_by, approved_at,
		       sent_at, rejected_by, rejected_at, rejection_reason, created_at, updated_at, archived_at
		FROM document_read_model
		WHERE deleted_at IS NULL
	`

	args := []interface{}{}
	argCount := 1

	if status != "" {
		query += fmt.Sprintf(" AND status = $%d", argCount)
		args = append(args, status)
		argCount++
	}

	if docType != "" {
		query += fmt.Sprintf(" AND document_type = $%d", argCount)
		args = append(args, docType)
		argCount++
	}

	if companyID != "" {
		query += fmt.Sprintf(" AND company_id = $%d", argCount)
		args = append(args, companyID)
		argCount++
	}

	if approvalStatus != "" {
		query += fmt.Sprintf(" AND approval_status = $%d", argCount)
		args = append(args, approvalStatus)
		argCount++
	}

	// Count total
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM (%s) t", query)
	var totalCount int64
	if err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&totalCount); err != nil {
		r.logger.Error("Failed to count documents", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to count documents: %w", err)
	}

	// Add ordering and pagination
	query += fmt.Sprintf(" ORDER BY created_at DESC LIMIT $%d OFFSET $%d", argCount, argCount+1)
	args = append(args, limit, offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		r.logger.Error("Failed to list documents", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list documents: %w", err)
	}
	defer rows.Close()

	var documents []*pb.DocumentReadModel
	for rows.Next() {
		var doc pb.DocumentReadModel
		if err := rows.Scan(
			&doc.Id,
			&doc.DocumentNumber,
			&doc.Title,
			&doc.Description,
			&doc.DocumentType,
			&doc.Status,
			&doc.CompanyId,
			&doc.CompanyName,
			&doc.CreatedByUserId,
			&doc.CreatedByUserName,
			&doc.ApprovalStatus,
			&doc.ApprovedBy,
			&doc.ApprovedAt,
			&doc.SentAt,
			&doc.RejectedBy,
			&doc.RejectedAt,
			&doc.RejectionReason,
			&doc.CreatedAt,
			&doc.UpdatedAt,
			&doc.ArchivedAt,
		); err != nil {
			r.logger.Error("Failed to scan document", zap.Error(err))
			return nil, 0, fmt.Errorf("failed to scan document: %w", err)
		}
		documents = append(documents, &doc)
	}

	if err = rows.Err(); err != nil {
		r.logger.Error("Row iteration error", zap.Error(err))
		return nil, 0, fmt.Errorf("row iteration error: %w", err)
	}

	return documents, totalCount, nil
}

// SearchDocuments performs search on document_number and title
func (r *DocumentQueryRepository) SearchDocuments(ctx context.Context, searchQuery string, offset, limit int32, docType, companyID, approvalStatus string) ([]*pb.DocumentReadModel, int64, error) {
	query := `
		SELECT id, document_number, title, description, document_type, status, company_id, company_name,
		       created_by_user_id, created_by_user_name, approval_status, approved_by, approved_at,
		       sent_at, rejected_by, rejected_at, rejection_reason, created_at, updated_at, archived_at
		FROM document_read_model
		WHERE deleted_at IS NULL
		  AND (
			document_number ILIKE $1
			OR title ILIKE $1
		  )
	`

	args := []interface{}{fmt.Sprintf("%%%s%%", searchQuery)}
	argCount := 2

	if docType != "" {
		query += fmt.Sprintf(" AND document_type = $%d", argCount)
		args = append(args, docType)
		argCount++
	}

	if companyID != "" {
		query += fmt.Sprintf(" AND company_id = $%d", argCount)
		args = append(args, companyID)
		argCount++
	}

	if approvalStatus != "" {
		query += fmt.Sprintf(" AND approval_status = $%d", argCount)
		args = append(args, approvalStatus)
		argCount++
	}

	// Count total
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM (%s) t", query)
	var totalCount int64
	if err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&totalCount); err != nil {
		r.logger.Error("Failed to count documents", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to count documents: %w", err)
	}

	// Add ordering and pagination
	query += fmt.Sprintf(" ORDER BY created_at DESC LIMIT $%d OFFSET $%d", argCount, argCount+1)
	args = append(args, limit, offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		r.logger.Error("Failed to search documents", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to search documents: %w", err)
	}
	defer rows.Close()

	var documents []*pb.DocumentReadModel
	for rows.Next() {
		var doc pb.DocumentReadModel
		if err := rows.Scan(
			&doc.Id,
			&doc.DocumentNumber,
			&doc.Title,
			&doc.Description,
			&doc.DocumentType,
			&doc.Status,
			&doc.CompanyId,
			&doc.CompanyName,
			&doc.CreatedByUserId,
			&doc.CreatedByUserName,
			&doc.ApprovalStatus,
			&doc.ApprovedBy,
			&doc.ApprovedAt,
			&doc.SentAt,
			&doc.RejectedBy,
			&doc.RejectedAt,
			&doc.RejectionReason,
			&doc.CreatedAt,
			&doc.UpdatedAt,
			&doc.ArchivedAt,
		); err != nil {
			r.logger.Error("Failed to scan document", zap.Error(err))
			return nil, 0, fmt.Errorf("failed to scan document: %w", err)
		}
		documents = append(documents, &doc)
	}

	if err = rows.Err(); err != nil {
		r.logger.Error("Row iteration error", zap.Error(err))
		return nil, 0, fmt.Errorf("row iteration error: %w", err)
	}

	return documents, totalCount, nil
}

// GetPendingApprovalDocuments retrieves documents pending approval for a company
func (r *DocumentQueryRepository) GetPendingApprovalDocuments(ctx context.Context, companyID string, offset, limit int32) ([]*pb.DocumentReadModel, int64, error) {
	query := `
		SELECT id, document_number, title, description, document_type, status, company_id, company_name,
		       created_by_user_id, created_by_user_name, approval_status, approved_by, approved_at,
		       sent_at, rejected_by, rejected_at, rejection_reason, created_at, updated_at, archived_at
		FROM document_read_model
		WHERE deleted_at IS NULL
		  AND company_id = $1
		  AND approval_status = 'pending'
	`

	args := []interface{}{companyID}

	// Count total
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM (%s) t", query)
	var totalCount int64
	if err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&totalCount); err != nil {
		r.logger.Error("Failed to count pending documents", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to count documents: %w", err)
	}

	// Add ordering and pagination
	query += fmt.Sprintf(" ORDER BY created_at DESC LIMIT $2 OFFSET $3")
	args = append(args, limit, offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		r.logger.Error("Failed to get pending documents", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to get pending documents: %w", err)
	}
	defer rows.Close()

	var documents []*pb.DocumentReadModel
	for rows.Next() {
		var doc pb.DocumentReadModel
		if err := rows.Scan(
			&doc.Id,
			&doc.DocumentNumber,
			&doc.Title,
			&doc.Description,
			&doc.DocumentType,
			&doc.Status,
			&doc.CompanyId,
			&doc.CompanyName,
			&doc.CreatedByUserId,
			&doc.CreatedByUserName,
			&doc.ApprovalStatus,
			&doc.ApprovedBy,
			&doc.ApprovedAt,
			&doc.SentAt,
			&doc.RejectedBy,
			&doc.RejectedAt,
			&doc.RejectionReason,
			&doc.CreatedAt,
			&doc.UpdatedAt,
			&doc.ArchivedAt,
		); err != nil {
			r.logger.Error("Failed to scan document", zap.Error(err))
			return nil, 0, fmt.Errorf("failed to scan document: %w", err)
		}
		documents = append(documents, &doc)
	}

	if err = rows.Err(); err != nil {
		r.logger.Error("Row iteration error", zap.Error(err))
		return nil, 0, fmt.Errorf("row iteration error: %w", err)
	}

	return documents, totalCount, nil
}
