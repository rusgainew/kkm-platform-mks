package handlers

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"

	"github.com/rusgainew/kkm-project-mks/document-query-server/internal/infrastructure/messaging"
)

// DocumentEventHandler implements messaging.DocumentEventHandler for document events
type DocumentEventHandler struct {
	db     *sqlx.DB
	logger *zap.Logger
}

// NewDocumentEventHandler creates a new document event handler
func NewDocumentEventHandler(db *sqlx.DB, logger *zap.Logger) *DocumentEventHandler {
	return &DocumentEventHandler{
		db:     db,
		logger: logger,
	}
}

// HandleDocumentCreated handles document.created events
func (h *DocumentEventHandler) HandleDocumentCreated(ctx context.Context, event messaging.DocumentEvent) error {
	// Extract data from event
	documentID, ok := event.Data["document_id"].(string)
	if !ok || documentID == "" {
		return fmt.Errorf("invalid document_id in event")
	}

	documentNumber, _ := event.Data["document_number"].(string)
	title, _ := event.Data["title"].(string)
	documentType, _ := event.Data["document_type"].(string)
	companyID, _ := event.Data["company_id"].(string)

	// Insert into read-model
	query := `
		INSERT INTO document_read_model (
			id, document_number, title, document_type, status, company_id, 
			approval_status, created_at
		)
		VALUES ($1, $2, $3, $4, 'draft', $5, 'pending', NOW())
		ON CONFLICT (id) DO UPDATE SET
			document_number = EXCLUDED.document_number,
			title = EXCLUDED.title,
			document_type = EXCLUDED.document_type,
			company_id = EXCLUDED.company_id
	`

	if _, err := h.db.ExecContext(ctx, query,
		documentID, documentNumber, title, documentType, companyID); err != nil {
		h.logger.Error("Failed to insert document",
			zap.String("document_id", documentID),
			zap.String("title", title),
			zap.Error(err))
		return fmt.Errorf("failed to insert document: %w", err)
	}

	h.logger.Info("Document created event processed",
		zap.String("document_id", documentID),
		zap.String("title", title))

	return nil
}

// HandleDocumentSent handles document.sent events
func (h *DocumentEventHandler) HandleDocumentSent(ctx context.Context, event messaging.DocumentEvent) error {
	documentID, ok := event.Data["document_id"].(string)
	if !ok || documentID == "" {
		return fmt.Errorf("invalid document_id in event")
	}

	// Update status to sent
	query := `
		UPDATE document_read_model 
		SET status = 'sent', sent_at = NOW()
		WHERE id = $1
	`

	result, err := h.db.ExecContext(ctx, query, documentID)
	if err != nil {
		h.logger.Error("Failed to update document sent status",
			zap.String("document_id", documentID),
			zap.Error(err))
		return fmt.Errorf("failed to update document: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		h.logger.Warn("Document not found for sent update",
			zap.String("document_id", documentID))
	}

	h.logger.Info("Document sent event processed",
		zap.String("document_id", documentID))

	return nil
}

// HandleDocumentApproved handles document.approved events
func (h *DocumentEventHandler) HandleDocumentApproved(ctx context.Context, event messaging.DocumentEvent) error {
	documentID, ok := event.Data["document_id"].(string)
	if !ok || documentID == "" {
		return fmt.Errorf("invalid document_id in event")
	}

	approvedBy, _ := event.Data["approved_by"].(string)

	// Update approval status
	query := `
		UPDATE document_read_model 
		SET approval_status = 'approved', approved_at = NOW(), approved_by = $2
		WHERE id = $1
	`

	if _, err := h.db.ExecContext(ctx, query, documentID, approvedBy); err != nil {
		h.logger.Error("Failed to update document approval status",
			zap.String("document_id", documentID),
			zap.Error(err))
		return fmt.Errorf("failed to update document: %w", err)
	}

	h.logger.Info("Document approved event processed",
		zap.String("document_id", documentID),
		zap.String("approved_by", approvedBy))

	return nil
}

// HandleDocumentRejected handles document.rejected events
func (h *DocumentEventHandler) HandleDocumentRejected(ctx context.Context, event messaging.DocumentEvent) error {
	documentID, ok := event.Data["document_id"].(string)
	if !ok || documentID == "" {
		return fmt.Errorf("invalid document_id in event")
	}

	rejectedBy, _ := event.Data["rejected_by"].(string)
	rejectionReason, _ := event.Data["rejection_reason"].(string)

	// Update rejection status
	query := `
		UPDATE document_read_model 
		SET approval_status = 'rejected', rejected_at = NOW(), 
		    rejected_by = $2, rejection_reason = $3
		WHERE id = $1
	`

	if _, err := h.db.ExecContext(ctx, query, documentID, rejectedBy, rejectionReason); err != nil {
		h.logger.Error("Failed to update document rejection status",
			zap.String("document_id", documentID),
			zap.Error(err))
		return fmt.Errorf("failed to update document: %w", err)
	}

	h.logger.Info("Document rejected event processed",
		zap.String("document_id", documentID),
		zap.String("rejected_by", rejectedBy))

	return nil
}

// HandleDocumentArchived handles document.archived events
func (h *DocumentEventHandler) HandleDocumentArchived(ctx context.Context, event messaging.DocumentEvent) error {
	documentID, ok := event.Data["document_id"].(string)
	if !ok || documentID == "" {
		return fmt.Errorf("invalid document_id in event")
	}

	// Update archived status
	query := `
		UPDATE document_read_model 
		SET archived_at = NOW()
		WHERE id = $1
	`

	if _, err := h.db.ExecContext(ctx, query, documentID); err != nil {
		h.logger.Error("Failed to update document archived status",
			zap.String("document_id", documentID),
			zap.Error(err))
		return fmt.Errorf("failed to update document: %w", err)
	}

	h.logger.Info("Document archived event processed",
		zap.String("document_id", documentID))

	return nil
}
