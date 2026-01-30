// Файл document-query-server/internal/application/handlers/document_event_handler.go содержит реализацию пакета handlers.
package handlers

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"

	"github.com/rusgainew/kkm-project-mks/document-query-server/internal/infrastructure/messaging"
	"github.com/rusgainew/kkm-project-mks/document-query-server/internal/infrastructure/repository"
	pb "github.com/rusgainew/kkm-project-mks/proto-lib/api"
)

// DocumentEventHandler implements messaging.DocumentEventHandler for document events
type DocumentEventHandler struct {
	repo   *repository.InMemoryDocumentRepository
	logger *zap.Logger
}

// NewDocumentEventHandler creates a new document event handler
func NewDocumentEventHandler(repo *repository.InMemoryDocumentRepository, logger *zap.Logger) *DocumentEventHandler {
	return &DocumentEventHandler{
		repo:   repo,
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
	description, _ := event.Data["description"].(string)
	documentType, _ := event.Data["document_type"].(string)
	companyID, _ := event.Data["company_id"].(string)
	companyName, _ := event.Data["company_name"].(string)
	createdByUserID, _ := event.Data["created_by_user_id"].(string)
	createdByUserName, _ := event.Data["created_by_user_name"].(string)

	// Create document read model
	doc := &pb.DocumentReadModel{
		Id:                documentID,
		DocumentNumber:    documentNumber,
		Title:             title,
		Description:       description,
		DocumentType:      documentType,
		Status:            "draft",
		CompanyId:         companyID,
		CompanyName:       companyName,
		CreatedByUserId:   createdByUserID,
		CreatedByUserName: createdByUserName,
		ApprovalStatus:    "pending",
		CreatedAt:         time.Now().Format(time.RFC3339),
		UpdatedAt:         time.Now().Format(time.RFC3339),
	}

	// Upsert into in-memory repository
	if err := h.repo.UpsertDocument(ctx, doc); err != nil {
		h.logger.Error("Failed to upsert document",
			zap.String("document_id", documentID),
			zap.String("title", title),
			zap.Error(err))
		return fmt.Errorf("failed to upsert document: %w", err)
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

	// Get existing document
	doc, err := h.repo.GetDocument(ctx, documentID)
	if err != nil {
		h.logger.Error("Failed to get document",
			zap.String("document_id", documentID),
			zap.Error(err))
		return fmt.Errorf("failed to get document: %w", err)
	}

	if doc == nil {
		h.logger.Warn("Document not found for sent update",
			zap.String("document_id", documentID))
		return nil
	}

	// Update status
	doc.Status = "sent"
	doc.SentAt = time.Now().Format(time.RFC3339)
	doc.UpdatedAt = time.Now().Format(time.RFC3339)

	// Update in repository
	if err := h.repo.UpsertDocument(ctx, doc); err != nil {
		h.logger.Error("Failed to update document sent status",
			zap.String("document_id", documentID),
			zap.Error(err))
		return fmt.Errorf("failed to update document: %w", err)
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

	// Get existing document
	doc, err := h.repo.GetDocument(ctx, documentID)
	if err != nil {
		h.logger.Error("Failed to get document",
			zap.String("document_id", documentID),
			zap.Error(err))
		return fmt.Errorf("failed to get document: %w", err)
	}

	if doc == nil {
		h.logger.Warn("Document not found for approval update",
			zap.String("document_id", documentID))
		return nil
	}

	// Update approval status
	doc.ApprovalStatus = "approved"
	doc.ApprovedBy = approvedBy
	doc.ApprovedAt = time.Now().Format(time.RFC3339)
	doc.UpdatedAt = time.Now().Format(time.RFC3339)

	// Update in repository
	if err := h.repo.UpsertDocument(ctx, doc); err != nil {
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

	// Get existing document
	doc, err := h.repo.GetDocument(ctx, documentID)
	if err != nil {
		h.logger.Error("Failed to get document",
			zap.String("document_id", documentID),
			zap.Error(err))
		return fmt.Errorf("failed to get document: %w", err)
	}

	if doc == nil {
		h.logger.Warn("Document not found for rejection update",
			zap.String("document_id", documentID))
		return nil
	}

	// Update rejection status
	doc.ApprovalStatus = "rejected"
	doc.RejectedBy = rejectedBy
	doc.RejectionReason = rejectionReason
	doc.RejectedAt = time.Now().Format(time.RFC3339)
	doc.UpdatedAt = time.Now().Format(time.RFC3339)

	// Update in repository
	if err := h.repo.UpsertDocument(ctx, doc); err != nil {
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

	// Get existing document
	doc, err := h.repo.GetDocument(ctx, documentID)
	if err != nil {
		h.logger.Error("Failed to get document",
			zap.String("document_id", documentID),
			zap.Error(err))
		return fmt.Errorf("failed to get document: %w", err)
	}

	if doc == nil {
		h.logger.Warn("Document not found for archive update",
			zap.String("document_id", documentID))
		return nil
	}

	// Update archived status
	doc.ArchivedAt = time.Now().Format(time.RFC3339)
	doc.UpdatedAt = time.Now().Format(time.RFC3339)

	// Update in repository
	if err := h.repo.UpsertDocument(ctx, doc); err != nil {
		h.logger.Error("Failed to update document archived status",
			zap.String("document_id", documentID),
			zap.Error(err))
		return fmt.Errorf("failed to update document: %w", err)
	}

	h.logger.Info("Document archived event processed",
		zap.String("document_id", documentID))

	return nil
}
