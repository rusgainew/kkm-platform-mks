package services

import (
	"context"
	"fmt"

	"github.com/rusgainew/kkm-project-mks/api-gateway/internal/domain/models"
	"github.com/rusgainew/kkm-project-mks/api-gateway/internal/infrastructure/client"
	"github.com/rusgainew/kkm-project-mks/api-gateway/internal/infrastructure/observability"
	docpb "github.com/rusgainew/kkm-project-mks/proto-lib/document"
	"go.uber.org/zap"
)

// DocumentService интегрирует документ-сервис через gRPC
// Он прозрачно пробрасывает контекст, метрики и трейсинг.
type DocumentService struct {
	connManager *client.ConnectionManager
	serviceAddr string
	metrics     *observability.Metrics
	tracer      *observability.Tracer
	logger      *zap.Logger
}

// NewDocumentService создает DocumentService
func NewDocumentService(
	connManager *client.ConnectionManager,
	serviceAddr string,
	metrics *observability.Metrics,
	tracer *observability.Tracer,
	logger *zap.Logger,
) *DocumentService {
	return &DocumentService{
		connManager: connManager,
		serviceAddr: serviceAddr,
		metrics:     metrics,
		tracer:      tracer,
		logger:      logger,
	}
}

func (s *DocumentService) getClient(ctx context.Context) (docpb.DocumentServiceClient, error) {
	conn, err := s.connManager.GetConnection(ctx, s.serviceAddr)
	if err != nil {
		return nil, fmt.Errorf("failed to get connection: %w", err)
	}
	return docpb.NewDocumentServiceClient(conn), nil
}

// CreateDocument создает документ
func (s *DocumentService) CreateDocument(ctx context.Context, doc *models.Document) (*models.Document, error) {
	ctx, endSpan := s.tracer.StartSpan(ctx, "DocumentService.CreateDocument")
	defer endSpan()

	client, err := s.getClient(ctx)
	if err != nil {
		s.logger.Error("Failed to get document client", zap.Error(err))
		return nil, err
	}

	req := &docpb.CreateDocumentRequest{
		OrganizationId: doc.OrganizationID,
		Title:          doc.Title,
		Content:        doc.Content,
		CreatedBy:      doc.CreatedBy,
	}

	resp, err := client.CreateDocument(ctx, req)
	if err != nil {
		s.metrics.IncrementErrorCount("document_service", "create_document")
		s.logger.Error("CreateDocument failed", zap.Error(err))
		return nil, fmt.Errorf("create document: %w", err)
	}

	return mapDocument(resp), nil
}

// GetDocument получает документ по ID
func (s *DocumentService) GetDocument(ctx context.Context, id string) (*models.Document, error) {
	ctx, endSpan := s.tracer.StartSpan(ctx, "DocumentService.GetDocument")
	defer endSpan()

	client, err := s.getClient(ctx)
	if err != nil {
		s.logger.Error("Failed to get document client", zap.Error(err))
		return nil, err
	}

	resp, err := client.GetDocument(ctx, &docpb.GetDocumentRequest{DocumentId: id})
	if err != nil {
		s.metrics.IncrementErrorCount("document_service", "get_document")
		s.logger.Error("GetDocument failed", zap.Error(err), zap.String("id", id))
		return nil, fmt.Errorf("get document: %w", err)
	}

	return mapDocument(resp), nil
}

// UpdateDocument обновляет заголовок/контент
func (s *DocumentService) UpdateDocument(ctx context.Context, doc *models.Document) (*models.Document, error) {
	ctx, endSpan := s.tracer.StartSpan(ctx, "DocumentService.UpdateDocument")
	defer endSpan()

	client, err := s.getClient(ctx)
	if err != nil {
		s.logger.Error("Failed to get document client", zap.Error(err))
		return nil, err
	}

	req := &docpb.UpdateDocumentRequest{
		Id:      doc.ID,
		Title:   doc.Title,
		Content: doc.Content,
	}

	resp, err := client.UpdateDocument(ctx, req)
	if err != nil {
		s.metrics.IncrementErrorCount("document_service", "update_document")
		s.logger.Error("UpdateDocument failed", zap.Error(err), zap.String("id", doc.ID))
		return nil, fmt.Errorf("update document: %w", err)
	}

	return mapDocument(resp), nil
}

// SendDocument отправляет документ на согласование
func (s *DocumentService) SendDocument(ctx context.Context, id, recipientID, message string) (*models.Document, error) {
	ctx, endSpan := s.tracer.StartSpan(ctx, "DocumentService.SendDocument")
	defer endSpan()

	client, err := s.getClient(ctx)
	if err != nil {
		s.logger.Error("Failed to get document client", zap.Error(err))
		return nil, err
	}

	req := &docpb.SendDocumentRequest{
		DocumentId:  id,
		RecipientId: recipientID,
		Message:     message,
	}

	resp, err := client.SendDocument(ctx, req)
	if err != nil {
		s.metrics.IncrementErrorCount("document_service", "send_document")
		s.logger.Error("SendDocument failed", zap.Error(err), zap.String("id", id))
		return nil, fmt.Errorf("send document: %w", err)
	}

	return mapDocument(resp), nil
}

// ApproveDocument одобряет документ
func (s *DocumentService) ApproveDocument(ctx context.Context, id, approvedBy, comments string) (*models.Document, error) {
	ctx, endSpan := s.tracer.StartSpan(ctx, "DocumentService.ApproveDocument")
	defer endSpan()

	client, err := s.getClient(ctx)
	if err != nil {
		s.logger.Error("Failed to get document client", zap.Error(err))
		return nil, err
	}

	req := &docpb.ApproveDocumentRequest{
		DocumentId: id,
		ApprovedBy: approvedBy,
		Comments:   comments,
	}

	resp, err := client.ApproveDocument(ctx, req)
	if err != nil {
		s.metrics.IncrementErrorCount("document_service", "approve_document")
		s.logger.Error("ApproveDocument failed", zap.Error(err), zap.String("id", id))
		return nil, fmt.Errorf("approve document: %w", err)
	}

	return mapDocument(resp), nil
}

// RejectDocument отклоняет документ
func (s *DocumentService) RejectDocument(ctx context.Context, id, rejectedBy, reason string) (*models.Document, error) {
	ctx, endSpan := s.tracer.StartSpan(ctx, "DocumentService.RejectDocument")
	defer endSpan()

	client, err := s.getClient(ctx)
	if err != nil {
		s.logger.Error("Failed to get document client", zap.Error(err))
		return nil, err
	}

	req := &docpb.RejectDocumentRequest{
		DocumentId: id,
		RejectedBy: rejectedBy,
		Reason:     reason,
	}

	resp, err := client.RejectDocument(ctx, req)
	if err != nil {
		s.metrics.IncrementErrorCount("document_service", "reject_document")
		s.logger.Error("RejectDocument failed", zap.Error(err), zap.String("id", id))
		return nil, fmt.Errorf("reject document: %w", err)
	}

	return mapDocument(resp), nil
}

// ArchiveDocument архивирует документ
func (s *DocumentService) ArchiveDocument(ctx context.Context, id string) (*models.Document, error) {
	ctx, endSpan := s.tracer.StartSpan(ctx, "DocumentService.ArchiveDocument")
	defer endSpan()

	client, err := s.getClient(ctx)
	if err != nil {
		s.logger.Error("Failed to get document client", zap.Error(err))
		return nil, err
	}

	resp, err := client.ArchiveDocument(ctx, &docpb.ArchiveDocumentRequest{DocumentId: id})
	if err != nil {
		s.metrics.IncrementErrorCount("document_service", "archive_document")
		s.logger.Error("ArchiveDocument failed", zap.Error(err), zap.String("id", id))
		return nil, fmt.Errorf("archive document: %w", err)
	}

	return mapDocument(resp), nil
}

// ListDocuments возвращает список документов с пагинацией
func (s *DocumentService) ListDocuments(ctx context.Context, organizationID, status, createdBy string, page, perPage int32) ([]models.Document, models.PageInfo, error) {
	ctx, endSpan := s.tracer.StartSpan(ctx, "DocumentService.ListDocuments")
	defer endSpan()

	client, err := s.getClient(ctx)
	if err != nil {
		s.logger.Error("Failed to get document client", zap.Error(err))
		return nil, models.PageInfo{}, err
	}

	req := &docpb.ListDocumentsRequest{
		OrganizationId: organizationID,
		Status:         status,
		CreatedBy:      createdBy,
		Page:           page,
		PerPage:        perPage,
	}

	resp, err := client.ListDocuments(ctx, req)
	if err != nil {
		s.metrics.IncrementErrorCount("document_service", "list_documents")
		s.logger.Error("ListDocuments failed", zap.Error(err), zap.String("organization_id", organizationID))
		return nil, models.PageInfo{}, fmt.Errorf("list documents: %w", err)
	}

	docs := make([]models.Document, 0, len(resp.GetDocuments()))
	for _, d := range resp.GetDocuments() {
		docs = append(docs, *mapDocument(d))
	}

	var pageInfo models.PageInfo
	if resp.PageInfo != nil {
		pageInfo = models.PageInfo{
			Page:       resp.PageInfo.Page,
			Size:       resp.PageInfo.PerPage,
			TotalCount: int32(resp.PageInfo.Total),
		}
	}

	return docs, pageInfo, nil
}

// mapDocument преобразует protobuf Document в доменную модель
func mapDocument(src *docpb.Document) *models.Document {
	if src == nil {
		return nil
	}
	entries := make([]models.DocumentEntry, 0, len(src.Entries))
	for _, e := range src.Entries {
		entries = append(entries, models.DocumentEntry{
			ID:         e.GetId(),
			DocumentID: e.GetDocumentId(),
			Key:        e.GetKey(),
			Value:      e.GetValue(),
			CreatedAt:  e.GetCreatedAt(),
			UpdatedAt:  e.GetUpdatedAt(),
		})
	}
	return &models.Document{
		ID:              src.GetId(),
		OrganizationID:  src.GetOrganizationId(),
		Title:           src.GetTitle(),
		Content:         src.GetContent(),
		Status:          src.GetStatus(),
		CreatedBy:       src.GetCreatedBy(),
		AssignedTo:      src.GetAssignedTo(),
		CreatedAt:       src.GetCreatedAt(),
		UpdatedAt:       src.GetUpdatedAt(),
		StatusChangedAt: src.GetStatusChangedAt(),
		Version:         src.GetVersion(),
		Entries:         entries,
	}
}
