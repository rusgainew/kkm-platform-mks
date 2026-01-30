// Файл document-query-server/internal/interfaces/grpc/handlers/document_query_handler.go содержит реализацию пакета handlers.
package handlers

import (
	"context"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/rusgainew/kkm-project-mks/document-query-server/internal/application/ports"
	pb "github.com/rusgainew/kkm-project-mks/proto-lib/api"
)

type DocumentQueryHandler struct {
	pb.UnimplementedDocumentQueryServiceServer
	logger *zap.Logger
	repo   ports.DocumentRepository
	cache  ports.DocumentCache
}

func NewDocumentQueryHandler(
	logger *zap.Logger,
	repo ports.DocumentRepository,
	cache ports.DocumentCache,
) *DocumentQueryHandler {
	return &DocumentQueryHandler{
		logger: logger,
		repo:   repo,
		cache:  cache,
	}
}

// GetDocument retrieves a single document by ID
func (h *DocumentQueryHandler) GetDocument(ctx context.Context, req *pb.GetDocumentRequest) (*pb.GetDocumentResponse, error) {
	if req.DocumentId == "" {
		h.logger.Warn("GetDocument called with empty document_id")
		return nil, status.Error(codes.InvalidArgument, "document_id is required")
	}

	// Try cache first
	if h.cache != nil {
		cached, err := h.cache.GetDocument(ctx, req.DocumentId)
		if err == nil && cached != nil {
			h.logger.Debug("Document cache hit", zap.String("document_id", req.DocumentId))
			return &pb.GetDocumentResponse{Document: cached}, nil
		}
	}

	// Query from read-model
	doc, err := h.repo.GetDocument(ctx, req.DocumentId)
	if err != nil {
		h.logger.Error("Failed to get document",
			zap.String("document_id", req.DocumentId),
			zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to get document")
	}

	if doc == nil {
		h.logger.Warn("Document not found", zap.String("document_id", req.DocumentId))
		return nil, status.Error(codes.NotFound, "document not found")
	}

	// Cache result
	if h.cache != nil {
		_ = h.cache.SetDocument(ctx, doc)
	}

	return &pb.GetDocumentResponse{Document: doc}, nil
}

// ListDocuments retrieves documents with pagination and optional filters
func (h *DocumentQueryHandler) ListDocuments(ctx context.Context, req *pb.ListDocumentsRequest) (*pb.ListDocumentsResponse, error) {
	// Validate pagination
	if req.Pagination == nil {
		h.logger.Warn("ListDocuments called without pagination")
		req.Pagination = &pb.PageInfo{Page: 0, Size: 10}
	}

	size := req.Pagination.Size
	if size <= 0 {
		size = 10
	}

	if size > 100 {
		size = 100
	}

	page := req.Pagination.Page
	if page < 0 {
		page = 0
	}

	offset := page * size

	// Query from read-model
	documents, totalCount, err := h.repo.ListDocuments(ctx, offset, size,
		req.Status, req.DocumentType, req.CompanyId, req.ApprovalStatus)
	if err != nil {
		h.logger.Error("Failed to list documents",
			zap.Int32("page", page),
			zap.Int32("size", size),
			zap.String("status", req.Status),
			zap.String("type", req.DocumentType),
			zap.String("company_id", req.CompanyId),
			zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to list documents")
	}

	h.logger.Debug("Listed documents",
		zap.Int32("page", page),
		zap.Int32("size", size),
		zap.Int("count", len(documents)),
		zap.Int64("total", totalCount))

	return &pb.ListDocumentsResponse{
		Documents:  documents,
		TotalCount: totalCount,
		Page:       page,
		PerPage:    size,
	}, nil
}

// SearchDocuments performs search on documents
func (h *DocumentQueryHandler) SearchDocuments(ctx context.Context, req *pb.SearchDocumentsRequest) (*pb.SearchDocumentsResponse, error) {
	// Validate search query
	if req.Query == "" {
		h.logger.Warn("SearchDocuments called with empty query")
		return nil, status.Error(codes.InvalidArgument, "search query is required")
	}

	// Validate pagination
	if req.Pagination == nil {
		h.logger.Warn("SearchDocuments called without pagination")
		req.Pagination = &pb.PageInfo{Page: 0, Size: 10}
	}

	size := req.Pagination.Size
	if size <= 0 {
		size = 10
	}

	if size > 100 {
		size = 100
	}

	page := req.Pagination.Page
	if page < 0 {
		page = 0
	}

	offset := page * size

	// Query from read-model
	documents, totalCount, err := h.repo.SearchDocuments(ctx, req.Query, offset, size,
		req.DocumentType, req.CompanyId, req.ApprovalStatus)
	if err != nil {
		h.logger.Error("Failed to search documents",
			zap.String("query", req.Query),
			zap.Int32("page", page),
			zap.String("type", req.DocumentType),
			zap.String("company_id", req.CompanyId),
			zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to search documents")
	}

	h.logger.Debug("Searched documents",
		zap.String("query", req.Query),
		zap.Int32("page", page),
		zap.Int("count", len(documents)),
		zap.Int64("total", totalCount))

	return &pb.SearchDocumentsResponse{
		Documents:  documents,
		TotalCount: totalCount,
		Page:       page,
		PerPage:    size,
	}, nil
}

// GetPendingApproval retrieves documents pending approval for a company
func (h *DocumentQueryHandler) GetPendingApproval(ctx context.Context, req *pb.GetPendingApprovalRequest) (*pb.GetPendingApprovalResponse, error) {
	if req.CompanyId == "" {
		h.logger.Warn("GetPendingApproval called with empty company_id")
		return nil, status.Error(codes.InvalidArgument, "company_id is required")
	}

	// Validate pagination
	if req.Pagination == nil {
		h.logger.Warn("GetPendingApproval called without pagination")
		req.Pagination = &pb.PageInfo{Page: 0, Size: 10}
	}

	size := req.Pagination.Size
	if size <= 0 {
		size = 10
	}

	if size > 100 {
		size = 100
	}

	page := req.Pagination.Page
	if page < 0 {
		page = 0
	}

	offset := page * size

	// Query from read-model
	documents, totalCount, err := h.repo.GetPendingApprovalDocuments(ctx, req.CompanyId, offset, size)
	if err != nil {
		h.logger.Error("Failed to get pending approval documents",
			zap.String("company_id", req.CompanyId),
			zap.Int32("page", page),
			zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to get pending documents")
	}

	h.logger.Debug("Got pending approval documents",
		zap.String("company_id", req.CompanyId),
		zap.Int32("page", page),
		zap.Int("count", len(documents)),
		zap.Int64("total", totalCount))

	return &pb.GetPendingApprovalResponse{
		Documents:  documents,
		TotalCount: totalCount,
		Page:       page,
		PerPage:    size,
	}, nil
}
