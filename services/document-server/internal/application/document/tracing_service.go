// Файл document-server/internal/application/document/tracing_service.go содержит реализацию пакета document.
package document

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/rusgainew/kkm-project-mks/document-server/internal/domain"
	"github.com/rusgainew/kkm-project-mks/document-server/internal/domain/events"
	"github.com/rusgainew/kkm-project-mks/document-server/internal/domain/ports"
	"github.com/rusgainew/kkm-project-mks/document-server/internal/infrastructure/observability"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

// ServiceWithTracing оборачивает Service с трассированием операций
type ServiceWithTracing struct {
	*Service
	tracing *observability.TracingProvider
}

// NewServiceWithTracing создаёт сервис с трассированием
func NewServiceWithTracing(
	repo ports.DocumentRepository,
	eventPublisher ports.EventPublisher,
	logger *zap.Logger,
	tracing *observability.TracingProvider,
) *ServiceWithTracing {
	return &ServiceWithTracing{
		Service: NewService(repo, eventPublisher, logger),
		tracing: tracing,
	}
}

// CreateDocumentWithTracing создаёт документ с трассированием
func (st *ServiceWithTracing) CreateDocumentWithTracing(ctx context.Context, organizationID, title, content, createdBy string) (string, error) {
	opSpan := st.tracing.NewOperationSpan(ctx, "create", "",
		attribute.String("organization_id", organizationID))

	// Валидация
	valSpan := st.tracing.NewOperationSpan(opSpan.Context(), "validate", "",
		attribute.String("field", "title"))
	if organizationID == "" || title == "" {
		valSpan.RecordError(domain.ErrInvalidInput)
		valSpan.End(domain.ErrInvalidInput)
		opSpan.End(domain.ErrInvalidInput)
		return "", domain.ErrInvalidInput
	}
	valSpan.End(nil)
	opSpan.AddEvent(observability.EventValidationPassed)

	documentID := uuid.New().String()
	now := time.Now().UnixMilli()

	// Создание в БД
	dbSpan := st.tracing.NewOperationSpan(opSpan.Context(), "db_create", documentID)
	if err := st.repo.Create(ctx, documentID, organizationID, title, content, createdBy, now); err != nil {
		st.logger.Error("Failed to create document", zap.Error(err), zap.String("document_id", documentID))
		dbSpan.RecordError(err)
		dbSpan.End(err)
		opSpan.End(err)
		return "", err
	}
	dbSpan.End(nil)
	opSpan.AddEvent(observability.EventDatabaseWrite,
		attribute.String("operation", "insert"))

	// Публикация события
	pubSpan := st.tracing.NewOperationSpan(opSpan.Context(), "publish_event", documentID)
	event := events.DocumentCreatedEvent{
		DocumentID:     documentID,
		OrganizationID: organizationID,
		Title:          title,
		CreatedBy:      createdBy,
		CreatedAt:      now,
	}
	_ = st.eventPublisher.Publish(opSpan.Context(), "document.created", event)
	pubSpan.End(nil)
	opSpan.AddEvent(observability.EventDocumentCreated)

	st.logger.Info("Document created", zap.String("document_id", documentID), zap.String("organization_id", organizationID))
	opSpan.End(nil)
	return documentID, nil
}

// GetDocumentWithTracing получает документ с трассированием
func (st *ServiceWithTracing) GetDocumentWithTracing(ctx context.Context, documentID string) (map[string]interface{}, error) {
	opSpan := st.tracing.NewOperationSpan(ctx, "get", documentID)

	if documentID == "" {
		opSpan.RecordError(domain.ErrInvalidDocumentID)
		opSpan.End(domain.ErrInvalidDocumentID)
		return nil, domain.ErrInvalidDocumentID
	}

	dbSpan := st.tracing.NewOperationSpan(opSpan.Context(), "db_get", documentID)
	doc, err := st.repo.Get(ctx, documentID)
	dbSpan.End(err)

	if err != nil {
		st.logger.Error("Failed to get document", zap.Error(err), zap.String("document_id", documentID))
		opSpan.End(err)
		return nil, err
	}

	opSpan.AddEvent(observability.EventDatabaseRead)
	opSpan.End(nil)
	return doc, nil
}

// UpdateDocumentWithTracing обновляет документ с трассированием
func (st *ServiceWithTracing) UpdateDocumentWithTracing(ctx context.Context, documentID, title, content string) error {
	opSpan := st.tracing.NewOperationSpan(ctx, "update", documentID)

	if documentID == "" {
		opSpan.RecordError(domain.ErrInvalidDocumentID)
		opSpan.End(domain.ErrInvalidDocumentID)
		return domain.ErrInvalidDocumentID
	}

	now := time.Now().UnixMilli()

	dbSpan := st.tracing.NewOperationSpan(opSpan.Context(), "db_update", documentID)
	if err := st.repo.Update(ctx, documentID, title, content, now); err != nil {
		st.logger.Error("Failed to update document", zap.Error(err), zap.String("document_id", documentID))
		dbSpan.RecordError(err)
		dbSpan.End(err)
		opSpan.End(err)
		return err
	}
	dbSpan.End(nil)
	opSpan.AddEvent(observability.EventDatabaseWrite, attribute.String("operation", "update"))

	// Публикация события
	pubSpan := st.tracing.NewOperationSpan(opSpan.Context(), "publish_event", documentID)
	event := events.DocumentUpdatedEvent{
		DocumentID: documentID,
		Title:      title,
		UpdatedAt:  now,
	}
	_ = st.eventPublisher.Publish(opSpan.Context(), "document.updated", event)
	pubSpan.End(nil)
	opSpan.AddEvent(observability.EventDocumentUpdated)

	st.logger.Info("Document updated", zap.String("document_id", documentID))
	opSpan.End(nil)
	return nil
}

// ListDocumentsWithTracing получает список документов с трассированием
func (st *ServiceWithTracing) ListDocumentsWithTracing(ctx context.Context, organizationID, status string, page, perPage int) ([]map[string]interface{}, int64, error) {
	opSpan := st.tracing.NewOperationSpan(ctx, "list", "",
		attribute.String("organization_id", organizationID),
		attribute.String("status", status),
		attribute.Int("page", page),
		attribute.Int("per_page", perPage))

	if organizationID == "" {
		opSpan.RecordError(domain.ErrInvalidOrganizationID)
		opSpan.End(domain.ErrInvalidOrganizationID)
		return nil, 0, domain.ErrInvalidOrganizationID
	}

	dbSpan := st.tracing.NewOperationSpan(opSpan.Context(), "db_list", "")
	docs, total, err := st.Service.ListDocuments(ctx, organizationID, status, page, perPage)
	dbSpan.End(err)

	if err != nil {
		st.logger.Error("Failed to list documents", zap.Error(err), zap.String("organization_id", organizationID))
		opSpan.End(err)
		return nil, 0, err
	}

	opSpan.SetAttribute(observability.ResultCountKey, len(docs))
	opSpan.AddEvent(observability.EventDatabaseRead, attribute.Int64("total", total))
	opSpan.End(nil)
	return docs, total, nil
}

// SendDocumentWithTracing отправляет документ с трассированием
func (st *ServiceWithTracing) SendDocumentWithTracing(ctx context.Context, documentID, recipientID, message string) error {
	opSpan := st.tracing.NewOperationSpan(ctx, "send", documentID,
		attribute.String("recipient_id", recipientID))

	if documentID == "" {
		opSpan.RecordError(domain.ErrInvalidDocumentID)
		opSpan.End(domain.ErrInvalidDocumentID)
		return domain.ErrInvalidDocumentID
	}

	now := time.Now().UnixMilli()

	dbSpan := st.tracing.NewOperationSpan(opSpan.Context(), "db_send", documentID)
	if err := st.repo.UpdateStatus(ctx, documentID, "sent", now); err != nil {
		st.logger.Error("Failed to send document", zap.Error(err), zap.String("document_id", documentID))
		dbSpan.RecordError(err)
		dbSpan.End(err)
		opSpan.End(err)
		return err
	}
	dbSpan.End(nil)

	pubSpan := st.tracing.NewOperationSpan(opSpan.Context(), "publish_event", documentID)
	event := events.DocumentSentEvent{
		DocumentID:  documentID,
		RecipientID: recipientID,
		Status:      "sent",
		SentAt:      now,
	}
	_ = st.eventPublisher.Publish(opSpan.Context(), "document.sent", event)
	pubSpan.End(nil)
	opSpan.AddEvent(observability.EventDocumentSent)

	st.logger.Info("Document sent", zap.String("document_id", documentID), zap.String("recipient_id", recipientID))
	opSpan.End(nil)
	return nil
}

// ApproveDocumentWithTracing одобряет документ с трассированием
func (st *ServiceWithTracing) ApproveDocumentWithTracing(ctx context.Context, documentID, approvedBy, comments string) error {
	opSpan := st.tracing.NewOperationSpan(ctx, "approve", documentID,
		attribute.String("approved_by", approvedBy))

	if documentID == "" {
		opSpan.RecordError(domain.ErrInvalidDocumentID)
		opSpan.End(domain.ErrInvalidDocumentID)
		return domain.ErrInvalidDocumentID
	}

	now := time.Now().UnixMilli()

	dbSpan := st.tracing.NewOperationSpan(opSpan.Context(), "db_approve", documentID)
	if err := st.repo.UpdateStatus(ctx, documentID, "approved", now); err != nil {
		st.logger.Error("Failed to approve document", zap.Error(err), zap.String("document_id", documentID))
		dbSpan.RecordError(err)
		dbSpan.End(err)
		opSpan.End(err)
		return err
	}
	dbSpan.End(nil)

	pubSpan := st.tracing.NewOperationSpan(opSpan.Context(), "publish_event", documentID)
	event := events.DocumentApprovedEvent{
		DocumentID: documentID,
		ApprovedBy: approvedBy,
		ApprovedAt: now,
	}
	_ = st.eventPublisher.Publish(opSpan.Context(), "document.approved", event)
	pubSpan.End(nil)
	opSpan.AddEvent(observability.EventDocumentApproved)

	st.logger.Info("Document approved", zap.String("document_id", documentID), zap.String("approved_by", approvedBy))
	opSpan.End(nil)
	return nil
}

// RejectDocumentWithTracing отклоняет документ с трассированием
func (st *ServiceWithTracing) RejectDocumentWithTracing(ctx context.Context, documentID, rejectedBy, reason string) error {
	opSpan := st.tracing.NewOperationSpan(ctx, "reject", documentID,
		attribute.String("rejected_by", rejectedBy))

	if documentID == "" {
		opSpan.RecordError(domain.ErrInvalidDocumentID)
		opSpan.End(domain.ErrInvalidDocumentID)
		return domain.ErrInvalidDocumentID
	}

	now := time.Now().UnixMilli()

	dbSpan := st.tracing.NewOperationSpan(opSpan.Context(), "db_reject", documentID)
	if err := st.repo.UpdateStatus(ctx, documentID, "rejected", now); err != nil {
		st.logger.Error("Failed to reject document", zap.Error(err), zap.String("document_id", documentID))
		dbSpan.RecordError(err)
		dbSpan.End(err)
		opSpan.End(err)
		return err
	}
	dbSpan.End(nil)

	pubSpan := st.tracing.NewOperationSpan(opSpan.Context(), "publish_event", documentID)
	event := events.DocumentRejectedEvent{
		DocumentID: documentID,
		RejectedBy: rejectedBy,
		Reason:     reason,
		RejectedAt: now,
	}
	_ = st.eventPublisher.Publish(opSpan.Context(), "document.rejected", event)
	pubSpan.End(nil)
	opSpan.AddEvent(observability.EventDocumentRejected)

	st.logger.Info("Document rejected", zap.String("document_id", documentID), zap.String("rejected_by", rejectedBy))
	opSpan.End(nil)
	return nil
}

// ArchiveDocumentWithTracing архивирует документ с трассированием
func (st *ServiceWithTracing) ArchiveDocumentWithTracing(ctx context.Context, documentID string) error {
	opSpan := st.tracing.NewOperationSpan(ctx, "archive", documentID)

	if documentID == "" {
		opSpan.RecordError(domain.ErrInvalidDocumentID)
		opSpan.End(domain.ErrInvalidDocumentID)
		return domain.ErrInvalidDocumentID
	}

	now := time.Now().UnixMilli()

	dbSpan := st.tracing.NewOperationSpan(opSpan.Context(), "db_archive", documentID)
	if err := st.repo.UpdateStatus(ctx, documentID, "archived", now); err != nil {
		st.logger.Error("Failed to archive document", zap.Error(err), zap.String("document_id", documentID))
		dbSpan.RecordError(err)
		dbSpan.End(err)
		opSpan.End(err)
		return err
	}
	dbSpan.End(nil)

	pubSpan := st.tracing.NewOperationSpan(opSpan.Context(), "publish_event", documentID)
	event := events.DocumentArchivedEvent{
		DocumentID: documentID,
		ArchivedAt: now,
	}
	_ = st.eventPublisher.Publish(opSpan.Context(), "document.archived", event)
	pubSpan.End(nil)
	opSpan.AddEvent(observability.EventDocumentArchived)

	st.logger.Info("Document archived", zap.String("document_id", documentID))
	opSpan.End(nil)
	return nil
}
