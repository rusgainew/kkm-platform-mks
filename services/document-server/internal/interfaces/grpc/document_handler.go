// Файл document-server/internal/interfaces/grpc/document_handler.go содержит реализацию пакета grpc.
package grpc

import (
	"context"
	"time"

	"github.com/rusgainew/kkm-project-mks/document-server/internal/application/document"
	"github.com/rusgainew/kkm-project-mks/document-server/internal/infrastructure/observability"
	common "github.com/rusgainew/kkm-project-mks/proto-lib/common"
	pb "github.com/rusgainew/kkm-project-mks/proto-lib/document"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type DocumentHandler struct {
	pb.UnimplementedDocumentServiceServer
	service *document.Service
	logger  *zap.Logger
	metrics *observability.Metrics
}

// NewDocumentHandler создает новый handler для документов
func NewDocumentHandler(service *document.Service, logger *zap.Logger, metrics *observability.Metrics) *DocumentHandler {
	return &DocumentHandler{
		service: service,
		logger:  logger,
		metrics: metrics,
	}
}

// GetDocument получает документ по ID
func (h *DocumentHandler) GetDocument(ctx context.Context, req *pb.GetDocumentRequest) (*pb.Document, error) {
	start := time.Now()
	h.metrics.IncActiveRequests()
	defer h.metrics.DecActiveRequests()
	defer func() {
		h.metrics.RecordRequestDuration(time.Since(start))
	}()

	if req.DocumentId == "" {
		h.metrics.RecordValidationError()
		return nil, status.Error(codes.InvalidArgument, "document_id is required")
	}

	doc, err := h.service.GetDocument(ctx, req.DocumentId)
	if err != nil {
		h.logger.Error("Failed to get document",
			zap.Error(err),
			zap.String("document_id", req.DocumentId))
		h.metrics.RecordRepositoryError()
		return nil, MapRepositoryErrorToGRPC(err)
	}

	return mapToDocumentPb(doc), nil
}

// CreateDocument создает новый документ
func (h *DocumentHandler) CreateDocument(ctx context.Context, req *pb.CreateDocumentRequest) (*pb.Document, error) {
	start := time.Now()
	h.metrics.IncActiveRequests()
	defer h.metrics.DecActiveRequests()
	defer func() {
		h.metrics.RecordRequestDuration(time.Since(start))
	}()

	// Валидация входных параметров
	validationErrors := document.ValidateCreateDocumentRequest(
		req.OrganizationId, req.Title, req.Content, req.CreatedBy)

	if len(validationErrors) > 0 {
		h.logger.Warn("Validation failed for CreateDocument",
			zap.Any("errors", validationErrors))
		h.metrics.RecordValidationError()
		return nil, MapValidationErrorsToGRPC(validationErrors)
	}

	docID, err := h.service.CreateDocument(ctx, req.OrganizationId, req.Title, req.Content, req.CreatedBy)
	if err != nil {
		h.logger.Error("Failed to create document",
			zap.Error(err),
			zap.String("organization_id", req.OrganizationId),
			zap.String("title", req.Title))
		h.metrics.RecordRepositoryError()
		return nil, MapDomainErrorToGRPC(err)
	}

	// Получаем созданный документ
	doc, err := h.service.GetDocument(ctx, docID)
	if err != nil {
		h.logger.Error("Failed to get created document",
			zap.Error(err),
			zap.String("document_id", docID))
		h.metrics.RecordRepositoryError()
		return nil, MapRepositoryErrorToGRPC(err)
	}

	h.logger.Info("Document created successfully",
		zap.String("document_id", docID),
		zap.String("organization_id", req.OrganizationId))

	h.metrics.RecordDocumentCreated()
	return mapToDocumentPb(doc), nil
}

// UpdateDocument обновляет документ
func (h *DocumentHandler) UpdateDocument(ctx context.Context, req *pb.UpdateDocumentRequest) (*pb.Document, error) {
	start := time.Now()
	h.metrics.IncActiveRequests()
	defer h.metrics.DecActiveRequests()
	defer func() {
		h.metrics.RecordRequestDuration(time.Since(start))
	}()

	// Валидация входных параметров
	validationErrors := document.ValidateUpdateDocumentRequest(
		req.Id, req.Title, req.Content)

	if len(validationErrors) > 0 {
		h.logger.Warn("Validation failed for UpdateDocument",
			zap.Any("errors", validationErrors))
		h.metrics.RecordValidationError()
		return nil, MapValidationErrorsToGRPC(validationErrors)
	}

	if err := h.service.UpdateDocument(ctx, req.Id, req.Title, req.Content); err != nil {
		h.logger.Error("Failed to update document",
			zap.Error(err),
			zap.String("document_id", req.Id))
		h.metrics.RecordRepositoryError()
		return nil, MapRepositoryErrorToGRPC(err)
	}

	doc, err := h.service.GetDocument(ctx, req.Id)
	if err != nil {
		h.logger.Error("Failed to get updated document",
			zap.Error(err),
			zap.String("document_id", req.Id))
		h.metrics.RecordRepositoryError()
		return nil, MapRepositoryErrorToGRPC(err)
	}

	h.logger.Info("Document updated successfully",
		zap.String("document_id", req.Id))

	h.metrics.RecordDocumentUpdated()
	return mapToDocumentPb(doc), nil
}

// SendDocument отправляет документ на согласование
func (h *DocumentHandler) SendDocument(ctx context.Context, req *pb.SendDocumentRequest) (*pb.Document, error) {
	start := time.Now()
	h.metrics.IncActiveRequests()
	defer h.metrics.DecActiveRequests()
	defer func() {
		h.metrics.RecordRequestDuration(time.Since(start))
	}()

	// Валидация входных параметров
	validationErrors := document.ValidateSendDocumentRequest(
		req.DocumentId, req.RecipientId, req.Message)

	if len(validationErrors) > 0 {
		h.logger.Warn("Validation failed for SendDocument",
			zap.Any("errors", validationErrors))
		h.metrics.RecordValidationError()
		return nil, MapValidationErrorsToGRPC(validationErrors)
	}

	if err := h.service.SendDocument(ctx, req.DocumentId, req.RecipientId, req.Message); err != nil {
		h.logger.Error("Failed to send document",
			zap.Error(err),
			zap.String("document_id", req.DocumentId))
		h.metrics.RecordRepositoryError()
		return nil, MapRepositoryErrorToGRPC(err)
	}

	doc, err := h.service.GetDocument(ctx, req.DocumentId)
	if err != nil {
		h.logger.Error("Failed to get sent document",
			zap.Error(err),
			zap.String("document_id", req.DocumentId))
		h.metrics.RecordRepositoryError()
		return nil, MapRepositoryErrorToGRPC(err)
	}

	h.logger.Info("Document sent successfully",
		zap.String("document_id", req.DocumentId),
		zap.String("recipient_id", req.RecipientId))

	h.metrics.RecordDocumentSent()
	return mapToDocumentPb(doc), nil
}

// ApproveDocument одобряет документ
func (h *DocumentHandler) ApproveDocument(ctx context.Context, req *pb.ApproveDocumentRequest) (*pb.Document, error) {
	start := time.Now()
	h.metrics.IncActiveRequests()
	defer h.metrics.DecActiveRequests()
	defer func() { h.metrics.RecordRequestDuration(time.Since(start)) }()

	// Валидация входных параметров
	validationErrors := document.ValidateApproveDocumentRequest(
		req.DocumentId, req.ApprovedBy)

	if len(validationErrors) > 0 {
		h.logger.Warn("Validation failed for ApproveDocument",
			zap.Any("errors", validationErrors))
		h.metrics.RecordValidationError()
		return nil, MapValidationErrorsToGRPC(validationErrors)
	}

	if err := h.service.ApproveDocument(ctx, req.DocumentId, req.ApprovedBy, req.Comments); err != nil {
		h.logger.Error("Failed to approve document",
			zap.Error(err),
			zap.String("document_id", req.DocumentId))
		h.metrics.RecordRepositoryError()
		return nil, MapRepositoryErrorToGRPC(err)
	}

	doc, err := h.service.GetDocument(ctx, req.DocumentId)
	if err != nil {
		h.logger.Error("Failed to get approved document",
			zap.Error(err),
			zap.String("document_id", req.DocumentId))
		h.metrics.RecordRepositoryError()
		return nil, MapRepositoryErrorToGRPC(err)
	}

	h.logger.Info("Document approved successfully",
		zap.String("document_id", req.DocumentId),
		zap.String("approved_by", req.ApprovedBy))

	h.metrics.RecordDocumentApproved()
	return mapToDocumentPb(doc), nil
}

// RejectDocument отклоняет документ
func (h *DocumentHandler) RejectDocument(ctx context.Context, req *pb.RejectDocumentRequest) (*pb.Document, error) {
	start := time.Now()
	h.metrics.IncActiveRequests()
	defer h.metrics.DecActiveRequests()
	defer func() { h.metrics.RecordRequestDuration(time.Since(start)) }()

	// Валидация входных параметров
	validationErrors := document.ValidateRejectDocumentRequest(
		req.DocumentId, req.RejectedBy, req.Reason)

	if len(validationErrors) > 0 {
		h.logger.Warn("Validation failed for RejectDocument",
			zap.Any("errors", validationErrors))
		h.metrics.RecordValidationError()
		return nil, MapValidationErrorsToGRPC(validationErrors)
	}

	if err := h.service.RejectDocument(ctx, req.DocumentId, req.RejectedBy, req.Reason); err != nil {
		h.logger.Error("Failed to reject document",
			zap.Error(err),
			zap.String("document_id", req.DocumentId))
		h.metrics.RecordRepositoryError()
		return nil, MapRepositoryErrorToGRPC(err)
	}

	doc, err := h.service.GetDocument(ctx, req.DocumentId)
	if err != nil {
		h.logger.Error("Failed to get rejected document",
			zap.Error(err),
			zap.String("document_id", req.DocumentId))
		h.metrics.RecordRepositoryError()
		return nil, MapRepositoryErrorToGRPC(err)
	}

	h.logger.Info("Document rejected successfully",
		zap.String("document_id", req.DocumentId),
		zap.String("rejected_by", req.RejectedBy))

	h.metrics.RecordDocumentRejected()
	return mapToDocumentPb(doc), nil
}

// ArchiveDocument архивирует документ
func (h *DocumentHandler) ArchiveDocument(ctx context.Context, req *pb.ArchiveDocumentRequest) (*pb.Document, error) {
	start := time.Now()
	h.metrics.IncActiveRequests()
	defer h.metrics.DecActiveRequests()
	defer func() { h.metrics.RecordRequestDuration(time.Since(start)) }()

	// Валидация входных параметров
	validationErrors := document.ValidateArchiveDocumentRequest(req.DocumentId)

	if len(validationErrors) > 0 {
		h.logger.Warn("Validation failed for ArchiveDocument",
			zap.Any("errors", validationErrors))
		h.metrics.RecordValidationError()
		return nil, MapValidationErrorsToGRPC(validationErrors)
	}

	if err := h.service.ArchiveDocument(ctx, req.DocumentId); err != nil {
		h.logger.Error("Failed to archive document",
			zap.Error(err),
			zap.String("document_id", req.DocumentId))
		h.metrics.RecordRepositoryError()
		return nil, MapRepositoryErrorToGRPC(err)
	}

	doc, err := h.service.GetDocument(ctx, req.DocumentId)
	if err != nil {
		h.logger.Error("Failed to get archived document",
			zap.Error(err),
			zap.String("document_id", req.DocumentId))
		h.metrics.RecordRepositoryError()
		return nil, MapRepositoryErrorToGRPC(err)
	}

	h.logger.Info("Document archived successfully",
		zap.String("document_id", req.DocumentId))

	h.metrics.RecordDocumentArchived()
	return mapToDocumentPb(doc), nil
}

// ListDocuments получает список документов
func (h *DocumentHandler) ListDocuments(ctx context.Context, req *pb.ListDocumentsRequest) (*pb.ListDocumentsResponse, error) {
	start := time.Now()
	h.metrics.IncActiveRequests()
	defer h.metrics.DecActiveRequests()
	defer func() { h.metrics.RecordRequestDuration(time.Since(start)) }()

	// Валидация входных параметров
	validationErrors := document.ValidateListDocumentsRequest(
		req.OrganizationId, int(req.Page), int(req.PerPage))

	if len(validationErrors) > 0 {
		h.logger.Warn("Validation failed for ListDocuments",
			zap.Any("errors", validationErrors))
		h.metrics.RecordValidationError()
		return nil, MapValidationErrorsToGRPC(validationErrors)
	}

	docs, total, err := h.service.ListDocuments(ctx, req.OrganizationId, req.Status, int(req.Page), int(req.PerPage))
	if err != nil {
		h.logger.Error("Failed to list documents",
			zap.Error(err),
			zap.String("organization_id", req.OrganizationId))
		h.metrics.RecordRepositoryError()
		return nil, MapRepositoryErrorToGRPC(err)
	}

	var pbDocs []*pb.Document
	for _, doc := range docs {
		pbDocs = append(pbDocs, mapToDocumentPb(doc))
	}

	h.logger.Info("Documents listed successfully",
		zap.String("organization_id", req.OrganizationId),
		zap.Int64("total", total))

	h.metrics.RecordPaginationRequest(int(req.PerPage))
	return &pb.ListDocumentsResponse{
		Documents: pbDocs,
		PageInfo: &common.PageInfo{
			Page:    req.Page,
			PerPage: req.PerPage,
			Total:   total,
		},
	}, nil
}

// Helper функция для маппинга документа
func mapToDocumentPb(doc map[string]interface{}) *pb.Document {
	var entries []*pb.DocumentEntry
	// TODO: загрузить entries из БД если необходимо

	return &pb.Document{
		Id:              doc["id"].(string),
		OrganizationId:  doc["organization_id"].(string),
		Title:           doc["title"].(string),
		Content:         doc["content"].(string),
		Status:          doc["status"].(string),
		CreatedBy:       doc["created_by"].(string),
		CreatedAt:       doc["created_at"].(int64),
		UpdatedAt:       doc["updated_at"].(int64),
		StatusChangedAt: doc["status_changed_at"].(int64),
		Version:         int32(doc["version"].(int32)),
		Entries:         entries,
	}
}
