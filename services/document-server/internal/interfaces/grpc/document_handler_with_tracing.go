// Файл document-server/internal/interfaces/grpc/document_handler_with_tracing.go содержит реализацию пакета grpc.
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

// DocumentHandlerWithTracing обрабатывает запросы с трассированием
type DocumentHandlerWithTracing struct {
	pb.UnimplementedDocumentServiceServer
	service *document.ServiceWithTracing
	logger  *zap.Logger
	metrics *observability.Metrics
}

// NewDocumentHandlerWithTracing создает новый handler с трассированием
func NewDocumentHandlerWithTracing(service *document.ServiceWithTracing, logger *zap.Logger, metrics *observability.Metrics) *DocumentHandlerWithTracing {
	return &DocumentHandlerWithTracing{
		service: service,
		logger:  logger,
		metrics: metrics,
	}
}

// GetDocument получает документ по ID с трассированием
func (h *DocumentHandlerWithTracing) GetDocument(ctx context.Context, req *pb.GetDocumentRequest) (*pb.Document, error) {
	start := time.Now()
	h.metrics.IncActiveRequests()
	defer h.metrics.DecActiveRequests()
	defer func() { h.metrics.RecordRequestDuration(time.Since(start)) }()

	if req.DocumentId == "" {
		h.metrics.RecordValidationError()
		return nil, status.Error(codes.InvalidArgument, "document_id is required")
	}

	doc, err := h.service.GetDocumentWithTracing(ctx, req.DocumentId)
	if err != nil {
		h.logger.Error("Failed to get document",
			zap.Error(err),
			zap.String("document_id", req.DocumentId))
		h.metrics.RecordRepositoryError()
		return nil, MapRepositoryErrorToGRPC(err)
	}

	return mapToDocumentPb(doc), nil
}

// CreateDocument создает новый документ с трассированием
func (h *DocumentHandlerWithTracing) CreateDocument(ctx context.Context, req *pb.CreateDocumentRequest) (*pb.Document, error) {
	start := time.Now()
	h.metrics.IncActiveRequests()
	defer h.metrics.DecActiveRequests()
	defer func() { h.metrics.RecordRequestDuration(time.Since(start)) }()

	// Валидация входных параметров
	validationErrors := document.ValidateCreateDocumentRequest(
		req.OrganizationId, req.Title, req.Content, req.CreatedBy)

	if len(validationErrors) > 0 {
		h.logger.Warn("Validation failed for CreateDocument",
			zap.Any("errors", validationErrors))
		h.metrics.RecordValidationError()
		return nil, MapValidationErrorsToGRPC(validationErrors)
	}

	docID, err := h.service.CreateDocumentWithTracing(ctx, req.OrganizationId, req.Title, req.Content, req.CreatedBy)
	if err != nil {
		h.logger.Error("Failed to create document",
			zap.Error(err),
			zap.String("organization_id", req.OrganizationId),
			zap.String("title", req.Title))
		h.metrics.RecordRepositoryError()
		return nil, MapDomainErrorToGRPC(err)
	}

	// Получаем созданный документ
	doc, err := h.service.GetDocumentWithTracing(ctx, docID)
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

// UpdateDocument обновляет документ с трассированием
func (h *DocumentHandlerWithTracing) UpdateDocument(ctx context.Context, req *pb.UpdateDocumentRequest) (*pb.Document, error) {
	start := time.Now()
	h.metrics.IncActiveRequests()
	defer h.metrics.DecActiveRequests()
	defer func() { h.metrics.RecordRequestDuration(time.Since(start)) }()

	validationErrors := document.ValidateUpdateDocumentRequest(req.Id, req.Title, req.Content)
	if len(validationErrors) > 0 {
		h.metrics.RecordValidationError()
		return nil, MapValidationErrorsToGRPC(validationErrors)
	}

	if err := h.service.UpdateDocumentWithTracing(ctx, req.Id, req.Title, req.Content); err != nil {
		h.logger.Error("Failed to update document",
			zap.Error(err),
			zap.String("document_id", req.Id))
		h.metrics.RecordRepositoryError()
		return nil, MapDomainErrorToGRPC(err)
	}

	doc, err := h.service.GetDocumentWithTracing(ctx, req.Id)
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

// SendDocument отправляет документ с трассированием
func (h *DocumentHandlerWithTracing) SendDocument(ctx context.Context, req *pb.SendDocumentRequest) (*pb.Document, error) {
	start := time.Now()
	h.metrics.IncActiveRequests()
	defer h.metrics.DecActiveRequests()
	defer func() { h.metrics.RecordRequestDuration(time.Since(start)) }()

	validationErrors := document.ValidateSendDocumentRequest(req.DocumentId, req.RecipientId, req.Message)
	if len(validationErrors) > 0 {
		h.metrics.RecordValidationError()
		return nil, MapValidationErrorsToGRPC(validationErrors)
	}

	if err := h.service.SendDocumentWithTracing(ctx, req.DocumentId, req.RecipientId, req.Message); err != nil {
		h.logger.Error("Failed to send document",
			zap.Error(err),
			zap.String("document_id", req.DocumentId))
		h.metrics.RecordRepositoryError()
		return nil, MapDomainErrorToGRPC(err)
	}

	doc, err := h.service.GetDocumentWithTracing(ctx, req.DocumentId)
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

	return mapToDocumentPb(doc), nil
}

// ListDocuments получает список документов с трассированием
func (h *DocumentHandlerWithTracing) ListDocuments(ctx context.Context, req *pb.ListDocumentsRequest) (*pb.ListDocumentsResponse, error) {
	start := time.Now()
	h.metrics.IncActiveRequests()
	defer h.metrics.DecActiveRequests()
	defer func() { h.metrics.RecordRequestDuration(time.Since(start)) }()

	validationErrors := document.ValidateListDocumentsRequest(req.OrganizationId, int(req.Page), int(req.PerPage))
	if len(validationErrors) > 0 {
		h.metrics.RecordValidationError()
		return nil, MapValidationErrorsToGRPC(validationErrors)
	}

	docs, total, err := h.service.ListDocumentsWithTracing(ctx, req.OrganizationId, req.Status, int(req.Page), int(req.PerPage))
	if err != nil {
		h.logger.Error("Failed to list documents",
			zap.Error(err),
			zap.String("organization_id", req.OrganizationId))
		h.metrics.RecordRepositoryError()
		return nil, MapRepositoryErrorToGRPC(err)
	}

	h.metrics.RecordPaginationRequest(int(req.PerPage))

	var pbDocs []*pb.Document
	for _, doc := range docs {
		pbDocs = append(pbDocs, mapToDocumentPb(doc))
	}

	return &pb.ListDocumentsResponse{
		Documents: pbDocs,
		PageInfo: &common.PageInfo{
			Page:    req.Page,
			PerPage: req.PerPage,
			Total:   total,
		},
	}, nil
}

// ApproveDocument одобряет документ с трассированием
func (h *DocumentHandlerWithTracing) ApproveDocument(ctx context.Context, req *pb.ApproveDocumentRequest) (*pb.Document, error) {
	start := time.Now()
	h.metrics.IncActiveRequests()
	defer h.metrics.DecActiveRequests()
	defer func() { h.metrics.RecordRequestDuration(time.Since(start)) }()

	validationErrors := document.ValidateApproveDocumentRequest(req.DocumentId, req.ApprovedBy)
	if len(validationErrors) > 0 {
		h.metrics.RecordValidationError()
		return nil, MapValidationErrorsToGRPC(validationErrors)
	}

	if err := h.service.ApproveDocumentWithTracing(ctx, req.DocumentId, req.ApprovedBy, req.Comments); err != nil {
		h.logger.Error("Failed to approve document",
			zap.Error(err),
			zap.String("document_id", req.DocumentId))
		h.metrics.RecordRepositoryError()
		return nil, MapDomainErrorToGRPC(err)
	}

	doc, err := h.service.GetDocumentWithTracing(ctx, req.DocumentId)
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

// RejectDocument отклоняет документ с трассированием
func (h *DocumentHandlerWithTracing) RejectDocument(ctx context.Context, req *pb.RejectDocumentRequest) (*pb.Document, error) {
	start := time.Now()
	h.metrics.IncActiveRequests()
	defer h.metrics.DecActiveRequests()
	defer func() { h.metrics.RecordRequestDuration(time.Since(start)) }()

	validationErrors := document.ValidateRejectDocumentRequest(req.DocumentId, req.RejectedBy, req.Reason)
	if len(validationErrors) > 0 {
		h.metrics.RecordValidationError()
		return nil, MapValidationErrorsToGRPC(validationErrors)
	}

	if err := h.service.RejectDocumentWithTracing(ctx, req.DocumentId, req.RejectedBy, req.Reason); err != nil {
		h.logger.Error("Failed to reject document",
			zap.Error(err),
			zap.String("document_id", req.DocumentId))
		h.metrics.RecordRepositoryError()
		return nil, MapDomainErrorToGRPC(err)
	}

	doc, err := h.service.GetDocumentWithTracing(ctx, req.DocumentId)
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

// ArchiveDocument архивирует документ с трассированием
func (h *DocumentHandlerWithTracing) ArchiveDocument(ctx context.Context, req *pb.ArchiveDocumentRequest) (*pb.Document, error) {
	start := time.Now()
	h.metrics.IncActiveRequests()
	defer h.metrics.DecActiveRequests()
	defer func() { h.metrics.RecordRequestDuration(time.Since(start)) }()

	if err := h.service.ArchiveDocumentWithTracing(ctx, req.DocumentId); err != nil {
		h.logger.Error("Failed to archive document",
			zap.Error(err),
			zap.String("document_id", req.DocumentId))
		h.metrics.RecordRepositoryError()
		return nil, MapDomainErrorToGRPC(err)
	}

	doc, err := h.service.GetDocumentWithTracing(ctx, req.DocumentId)
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
