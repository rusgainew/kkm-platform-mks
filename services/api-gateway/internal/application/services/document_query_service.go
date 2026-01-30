// Файл api-gateway/internal/application/services/document_query_service.go содержит реализацию пакета services.
package services

import (
	"context"
	"fmt"

	"github.com/rusgainew/kkm-project-mks/api-gateway/internal/infrastructure/client"
	"github.com/rusgainew/kkm-project-mks/api-gateway/internal/infrastructure/observability"
	pb "github.com/rusgainew/kkm-project-mks/proto-lib/api"
	"go.uber.org/zap"
)

// DocumentQueryService сервис для чтения документов (query side)
type DocumentQueryService struct {
	connManager *client.ConnectionManager
	serviceAddr string
	metrics     *observability.Metrics
	tracer      *observability.Tracer
	logger      *zap.Logger
}

// NewDocumentQueryService создает новый DocumentQueryService
func NewDocumentQueryService(
	connManager *client.ConnectionManager,
	serviceAddr string,
	metrics *observability.Metrics,
	tracer *observability.Tracer,
	logger *zap.Logger,
) *DocumentQueryService {
	return &DocumentQueryService{
		connManager: connManager,
		serviceAddr: serviceAddr,
		metrics:     metrics,
		tracer:      tracer,
		logger:      logger,
	}
}

// getClient возвращает gRPC клиент для document query service
func (s *DocumentQueryService) getClient(ctx context.Context) (pb.DocumentQueryServiceClient, error) {
	conn, err := s.connManager.GetConnection(ctx, s.serviceAddr)
	if err != nil {
		return nil, fmt.Errorf("failed to get connection: %w", err)
	}
	return pb.NewDocumentQueryServiceClient(conn), nil
}

// GetDocument получает документ по ID
func (s *DocumentQueryService) GetDocument(ctx context.Context, documentID string) (*pb.DocumentReadModel, error) {
	ctx, endSpan := s.tracer.StartSpan(ctx, "DocumentQueryService.GetDocument")
	defer endSpan()

	client, err := s.getClient(ctx)
	if err != nil {
		s.logger.Error("Failed to get document query client", zap.Error(err))
		s.metrics.IncrementErrorCount("document_query_service", "get_document")
		return nil, fmt.Errorf("failed to get client: %w", err)
	}

	req := &pb.GetDocumentRequest{
		DocumentId: documentID,
	}

	resp, err := client.GetDocument(ctx, req)
	if err != nil {
		s.logger.Error("Failed to get document", zap.String("document_id", documentID), zap.Error(err))
		s.metrics.IncrementErrorCount("document_query_service", "get_document")
		return nil, fmt.Errorf("failed to get document: %w", err)
	}

	return resp.Document, nil
}

// ListDocuments получает список документов с пагинацией и фильтрацией
func (s *DocumentQueryService) ListDocuments(ctx context.Context, page, pageSize int32, status, docType, companyID, approvalStatus string) (*pb.ListDocumentsResponse, error) {
	ctx, endSpan := s.tracer.StartSpan(ctx, "DocumentQueryService.ListDocuments")
	defer endSpan()

	client, err := s.getClient(ctx)
	if err != nil {
		s.logger.Error("Failed to get document query client", zap.Error(err))
		s.metrics.IncrementErrorCount("document_query_service", "list_documents")
		return nil, fmt.Errorf("failed to get client: %w", err)
	}

	req := &pb.ListDocumentsRequest{
		Pagination: &pb.PageInfo{
			Page: page,
			Size: pageSize,
		},
		Status:         status,
		DocumentType:   docType,
		CompanyId:      companyID,
		ApprovalStatus: approvalStatus,
	}

	resp, err := client.ListDocuments(ctx, req)
	if err != nil {
		s.logger.Error("Failed to list documents",
			zap.Int32("page", page),
			zap.Int32("page_size", pageSize),
			zap.Error(err))
		s.metrics.IncrementErrorCount("document_query_service", "list_documents")
		return nil, fmt.Errorf("failed to list documents: %w", err)
	}

	return resp, nil
}

// SearchDocuments выполняет поиск документов
func (s *DocumentQueryService) SearchDocuments(ctx context.Context, query string, page, pageSize int32, docType, companyID, approvalStatus string) (*pb.SearchDocumentsResponse, error) {
	ctx, endSpan := s.tracer.StartSpan(ctx, "DocumentQueryService.SearchDocuments")
	defer endSpan()

	client, err := s.getClient(ctx)
	if err != nil {
		s.logger.Error("Failed to get document query client", zap.Error(err))
		s.metrics.IncrementErrorCount("document_query_service", "search_documents")
		return nil, fmt.Errorf("failed to get client: %w", err)
	}

	req := &pb.SearchDocumentsRequest{
		Query: query,
		Pagination: &pb.PageInfo{
			Page: page,
			Size: pageSize,
		},
		DocumentType:   docType,
		CompanyId:      companyID,
		ApprovalStatus: approvalStatus,
	}

	resp, err := client.SearchDocuments(ctx, req)
	if err != nil {
		s.logger.Error("Failed to search documents",
			zap.String("query", query),
			zap.Int32("page", page),
			zap.Error(err))
		s.metrics.IncrementErrorCount("document_query_service", "search_documents")
		return nil, fmt.Errorf("failed to search documents: %w", err)
	}

	return resp, nil
}

// GetPendingApprovalDocuments получает документы ожидающие одобрения
func (s *DocumentQueryService) GetPendingApprovalDocuments(ctx context.Context, companyID string, page, pageSize int32) (*pb.GetPendingApprovalResponse, error) {
	ctx, endSpan := s.tracer.StartSpan(ctx, "DocumentQueryService.GetPendingApprovalDocuments")
	defer endSpan()

	client, err := s.getClient(ctx)
	if err != nil {
		s.logger.Error("Failed to get document query client", zap.Error(err))
		s.metrics.IncrementErrorCount("document_query_service", "get_pending_approval")
		return nil, fmt.Errorf("failed to get client: %w", err)
	}

	req := &pb.GetPendingApprovalRequest{
		CompanyId: companyID,
		Pagination: &pb.PageInfo{
			Page: page,
			Size: pageSize,
		},
	}

	resp, err := client.GetPendingApproval(ctx, req)
	if err != nil {
		s.logger.Error("Failed to get pending approval documents",
			zap.String("company_id", companyID),
			zap.Int32("page", page),
			zap.Error(err))
		s.metrics.IncrementErrorCount("document_query_service", "get_pending_approval")
		return nil, fmt.Errorf("failed to get pending approval documents: %w", err)
	}

	return resp, nil
}
