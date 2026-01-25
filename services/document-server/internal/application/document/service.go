package document

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/rusgainew/kkm-project-mks/document-server/internal/domain"
	"github.com/rusgainew/kkm-project-mks/document-server/internal/domain/events"
	"github.com/rusgainew/kkm-project-mks/document-server/internal/domain/ports"
	"go.uber.org/zap"
)

type Service struct {
	repo           ports.DocumentRepository
	eventPublisher ports.EventPublisher
	logger         *zap.Logger
}

// NewService создает новый сервис документов
func NewService(repo ports.DocumentRepository, eventPublisher ports.EventPublisher, logger *zap.Logger) *Service {
	return &Service{
		repo:           repo,
		eventPublisher: eventPublisher,
		logger:         logger,
	}
}

// CreateDocument создает новый документ
func (s *Service) CreateDocument(ctx context.Context, organizationID, title, content, createdBy string) (string, error) {
	if organizationID == "" || title == "" {
		return "", domain.ErrInvalidInput
	}

	documentID := uuid.New().String()
	now := time.Now().UnixMilli()

	if err := s.repo.Create(ctx, documentID, organizationID, title, content, createdBy, now); err != nil {
		s.logger.Error("Failed to create document", zap.Error(err), zap.String("document_id", documentID))
		return "", err
	}

	// Публикуем событие
	event := events.DocumentCreatedEvent{
		DocumentID:     documentID,
		OrganizationID: organizationID,
		Title:          title,
		CreatedBy:      createdBy,
		CreatedAt:      now,
	}
	_ = s.eventPublisher.Publish(ctx, "document.created", event)

	s.logger.Info("Document created", zap.String("document_id", documentID), zap.String("organization_id", organizationID))
	return documentID, nil
}

// GetDocument получает документ по ID
func (s *Service) GetDocument(ctx context.Context, documentID string) (map[string]interface{}, error) {
	if documentID == "" {
		return nil, domain.ErrInvalidDocumentID
	}

	doc, err := s.repo.Get(ctx, documentID)
	if err != nil {
		s.logger.Error("Failed to get document", zap.Error(err), zap.String("document_id", documentID))
		return nil, err
	}

	return doc, nil
}

// UpdateDocument обновляет документ
func (s *Service) UpdateDocument(ctx context.Context, documentID, title, content string) error {
	if documentID == "" {
		return domain.ErrInvalidDocumentID
	}

	now := time.Now().UnixMilli()

	if err := s.repo.Update(ctx, documentID, title, content, now); err != nil {
		s.logger.Error("Failed to update document", zap.Error(err), zap.String("document_id", documentID))
		return err
	}

	// Публикуем событие
	event := events.DocumentUpdatedEvent{
		DocumentID: documentID,
		Title:      title,
		UpdatedAt:  now,
	}
	_ = s.eventPublisher.Publish(ctx, "document.updated", event)

	s.logger.Info("Document updated", zap.String("document_id", documentID))
	return nil
}

// SendDocument отправляет документ на согласование
func (s *Service) SendDocument(ctx context.Context, documentID, recipientID, message string) error {
	if documentID == "" {
		return domain.ErrInvalidDocumentID
	}

	now := time.Now().UnixMilli()

	if err := s.repo.UpdateStatus(ctx, documentID, "sent", now); err != nil {
		s.logger.Error("Failed to send document", zap.Error(err), zap.String("document_id", documentID))
		return err
	}

	// Публикуем событие
	event := events.DocumentSentEvent{
		DocumentID:  documentID,
		RecipientID: recipientID,
		Status:      "sent",
		SentAt:      now,
	}
	_ = s.eventPublisher.Publish(ctx, "document.sent", event)

	s.logger.Info("Document sent", zap.String("document_id", documentID), zap.String("recipient_id", recipientID))
	return nil
}

// ApproveDocument одобряет документ
func (s *Service) ApproveDocument(ctx context.Context, documentID, approvedBy, comments string) error {
	if documentID == "" {
		return domain.ErrInvalidDocumentID
	}

	now := time.Now().UnixMilli()

	if err := s.repo.UpdateStatus(ctx, documentID, "approved", now); err != nil {
		s.logger.Error("Failed to approve document", zap.Error(err), zap.String("document_id", documentID))
		return err
	}

	// Публикуем событие
	event := events.DocumentApprovedEvent{
		DocumentID: documentID,
		ApprovedBy: approvedBy,
		ApprovedAt: now,
	}
	_ = s.eventPublisher.Publish(ctx, "document.approved", event)

	s.logger.Info("Document approved", zap.String("document_id", documentID), zap.String("approved_by", approvedBy))
	return nil
}

// RejectDocument отклоняет документ
func (s *Service) RejectDocument(ctx context.Context, documentID, rejectedBy, reason string) error {
	if documentID == "" {
		return domain.ErrInvalidDocumentID
	}

	now := time.Now().UnixMilli()

	if err := s.repo.UpdateStatus(ctx, documentID, "rejected", now); err != nil {
		s.logger.Error("Failed to reject document", zap.Error(err), zap.String("document_id", documentID))
		return err
	}

	// Публикуем событие
	event := events.DocumentRejectedEvent{
		DocumentID: documentID,
		RejectedBy: rejectedBy,
		Reason:     reason,
		RejectedAt: now,
	}
	_ = s.eventPublisher.Publish(ctx, "document.rejected", event)

	s.logger.Info("Document rejected", zap.String("document_id", documentID), zap.String("rejected_by", rejectedBy))
	return nil
}

// ArchiveDocument архивирует документ
func (s *Service) ArchiveDocument(ctx context.Context, documentID string) error {
	if documentID == "" {
		return domain.ErrInvalidDocumentID
	}

	now := time.Now().UnixMilli()

	if err := s.repo.UpdateStatus(ctx, documentID, "archived", now); err != nil {
		s.logger.Error("Failed to archive document", zap.Error(err), zap.String("document_id", documentID))
		return err
	}

	// Публикуем событие
	event := events.DocumentArchivedEvent{
		DocumentID: documentID,
		ArchivedAt: now,
	}
	_ = s.eventPublisher.Publish(ctx, "document.archived", event)

	s.logger.Info("Document archived", zap.String("document_id", documentID))
	return nil
}

// ListDocuments получает список документов по организации
func (s *Service) ListDocuments(ctx context.Context, organizationID, status string, page, perPage int) ([]map[string]interface{}, int64, error) {
	if organizationID == "" {
		return nil, 0, domain.ErrInvalidOrganizationID
	}

	docs, total, err := s.repo.List(ctx, organizationID, status, page, perPage)
	if err != nil {
		s.logger.Error("Failed to list documents", zap.Error(err), zap.String("organization_id", organizationID))
		return nil, 0, err
	}

	return docs, total, nil
}
