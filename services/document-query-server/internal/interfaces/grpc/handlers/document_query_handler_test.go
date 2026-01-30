// Файл document-query-server/internal/interfaces/grpc/handlers/document_query_handler_test.go содержит реализацию пакета handlers.
package handlers

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"

	pb "github.com/rusgainew/kkm-project-mks/proto-lib/api"
)

// MockDocumentRepository мок для repository
type MockDocumentRepository struct {
	mock.Mock
}

func (m *MockDocumentRepository) GetDocument(ctx context.Context, documentID string) (*pb.DocumentReadModel, error) {
	args := m.Called(ctx, documentID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*pb.DocumentReadModel), args.Error(1)
}

func (m *MockDocumentRepository) ListDocuments(ctx context.Context, offset, limit int32, status, docType, companyID, approvalStatus string) ([]*pb.DocumentReadModel, int64, error) {
	args := m.Called(ctx, offset, limit, status, docType, companyID, approvalStatus)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).([]*pb.DocumentReadModel), args.Get(1).(int64), args.Error(2)
}

func (m *MockDocumentRepository) SearchDocuments(ctx context.Context, query string, offset, limit int32, docType, companyID, approvalStatus string) ([]*pb.DocumentReadModel, int64, error) {
	args := m.Called(ctx, query, offset, limit, docType, companyID, approvalStatus)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).([]*pb.DocumentReadModel), args.Get(1).(int64), args.Error(2)
}

func (m *MockDocumentRepository) GetPendingApprovalDocuments(ctx context.Context, companyID string, offset, limit int32) ([]*pb.DocumentReadModel, int64, error) {
	args := m.Called(ctx, companyID, offset, limit)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).([]*pb.DocumentReadModel), args.Get(1).(int64), args.Error(2)
}

// MockDocumentRedisCache мок для cache
type MockDocumentRedisCache struct {
	mock.Mock
}

func (m *MockDocumentRedisCache) GetDocument(ctx context.Context, documentID string) (*pb.DocumentReadModel, error) {
	args := m.Called(ctx, documentID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*pb.DocumentReadModel), args.Error(1)
}

func (m *MockDocumentRedisCache) SetDocument(ctx context.Context, doc *pb.DocumentReadModel) error {
	args := m.Called(ctx, doc)
	return args.Error(0)
}

func (m *MockDocumentRedisCache) InvalidateDocument(ctx context.Context, documentID string) error {
	args := m.Called(ctx, documentID)
	return args.Error(0)
}

// Test GetDocument
func TestGetDocument_Success(t *testing.T) {
	logger, _ := zap.NewProduction()
	mockRepo := new(MockDocumentRepository)
	mockCache := new(MockDocumentRedisCache)

	expectedDoc := &pb.DocumentReadModel{
		Id:             "doc-123",
		DocumentNumber: "DOC-2026-001",
		Title:          "Invoice",
		DocumentType:   "INVOICE",
		CompanyId:      "comp-123",
		Status:         "ACTIVE",
		ApprovalStatus: "APPROVED",
	}

	mockCache.On("GetDocument", mock.Anything, "doc-123").Return(nil, nil)
	mockRepo.On("GetDocument", mock.Anything, "doc-123").Return(expectedDoc, nil)
	mockCache.On("SetDocument", mock.Anything, expectedDoc).Return(nil)

	handler := &DocumentQueryHandler{
		logger: logger,
		repo:   mockRepo,
		cache:  mockCache,
	}

	req := &pb.GetDocumentRequest{DocumentId: "doc-123"}
	resp, err := handler.GetDocument(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "doc-123", resp.Document.Id)
	assert.Equal(t, "DOC-2026-001", resp.Document.DocumentNumber)
	assert.Equal(t, "INVOICE", resp.Document.DocumentType)

	mockCache.AssertCalled(t, "GetDocument", mock.Anything, "doc-123")
	mockRepo.AssertCalled(t, "GetDocument", mock.Anything, "doc-123")
	mockCache.AssertCalled(t, "SetDocument", mock.Anything, expectedDoc)
}

func TestGetDocument_CacheHit(t *testing.T) {
	logger, _ := zap.NewProduction()
	mockRepo := new(MockDocumentRepository)
	mockCache := new(MockDocumentRedisCache)

	cachedDoc := &pb.DocumentReadModel{
		Id:             "doc-456",
		DocumentNumber: "DOC-2026-002",
		Status:         "ACTIVE",
	}

	mockCache.On("GetDocument", mock.Anything, "doc-456").Return(cachedDoc, nil)

	handler := &DocumentQueryHandler{
		logger: logger,
		repo:   mockRepo,
		cache:  mockCache,
	}

	req := &pb.GetDocumentRequest{DocumentId: "doc-456"}
	resp, err := handler.GetDocument(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "doc-456", resp.Document.Id)

	mockCache.AssertCalled(t, "GetDocument", mock.Anything, "doc-456")
	mockRepo.AssertNotCalled(t, "GetDocument")
}

func TestGetDocument_NotFound(t *testing.T) {
	logger, _ := zap.NewProduction()
	mockRepo := new(MockDocumentRepository)
	mockCache := new(MockDocumentRedisCache)

	mockCache.On("GetDocument", mock.Anything, "nonexistent").Return(nil, nil)
	mockRepo.On("GetDocument", mock.Anything, "nonexistent").Return(nil, nil)

	handler := &DocumentQueryHandler{
		logger: logger,
		repo:   mockRepo,
		cache:  mockCache,
	}

	req := &pb.GetDocumentRequest{DocumentId: "nonexistent"}
	resp, err := handler.GetDocument(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, resp)
}

func TestGetDocument_EmptyID(t *testing.T) {
	logger, _ := zap.NewProduction()
	mockRepo := new(MockDocumentRepository)
	mockCache := new(MockDocumentRedisCache)

	handler := &DocumentQueryHandler{
		logger: logger,
		repo:   mockRepo,
		cache:  mockCache,
	}

	req := &pb.GetDocumentRequest{DocumentId: ""}
	resp, err := handler.GetDocument(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, resp)

	mockCache.AssertNotCalled(t, "GetDocument")
	mockRepo.AssertNotCalled(t, "GetDocument")
}

// Test ListDocuments
func TestListDocuments_Success(t *testing.T) {
	logger, _ := zap.NewProduction()
	mockRepo := new(MockDocumentRepository)
	mockCache := new(MockDocumentRedisCache)

	expectedDocs := []*pb.DocumentReadModel{
		{Id: "doc-1", DocumentNumber: "DOC-001", Status: "ACTIVE"},
		{Id: "doc-2", DocumentNumber: "DOC-002", Status: "ACTIVE"},
	}

	mockRepo.On("ListDocuments", mock.Anything, int32(0), int32(10), "", "", "", "").Return(expectedDocs, int64(2), nil)

	handler := &DocumentQueryHandler{
		logger: logger,
		repo:   mockRepo,
		cache:  mockCache,
	}

	req := &pb.ListDocumentsRequest{
		Pagination: &pb.PageInfo{Page: 0, Size: 10},
	}
	resp, err := handler.ListDocuments(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Len(t, resp.Documents, 2)
	assert.Equal(t, int64(2), resp.TotalCount)

	mockRepo.AssertCalled(t, "ListDocuments", mock.Anything, int32(0), int32(10), "", "", "", "")
}

func TestListDocuments_WithFilters(t *testing.T) {
	logger, _ := zap.NewProduction()
	mockRepo := new(MockDocumentRepository)
	mockCache := new(MockDocumentRedisCache)

	expectedDocs := []*pb.DocumentReadModel{
		{Id: "doc-1", DocumentNumber: "INV-001", DocumentType: "INVOICE", Status: "ACTIVE"},
	}

	mockRepo.On("ListDocuments", mock.Anything, int32(0), int32(10), "ACTIVE", "INVOICE", "comp-1", "").
		Return(expectedDocs, int64(1), nil)

	handler := &DocumentQueryHandler{
		logger: logger,
		repo:   mockRepo,
		cache:  mockCache,
	}

	req := &pb.ListDocumentsRequest{
		Pagination:     &pb.PageInfo{Page: 0, Size: 10},
		Status:         "ACTIVE",
		DocumentType:   "INVOICE",
		CompanyId:      "comp-1",
		ApprovalStatus: "",
	}
	resp, err := handler.ListDocuments(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Len(t, resp.Documents, 1)
	assert.Equal(t, "INVOICE", resp.Documents[0].DocumentType)

	mockRepo.AssertCalled(t, "ListDocuments", mock.Anything, int32(0), int32(10), "ACTIVE", "INVOICE", "comp-1", "")
}

func TestListDocuments_DefaultPagination(t *testing.T) {
	logger, _ := zap.NewProduction()
	mockRepo := new(MockDocumentRepository)
	mockCache := new(MockDocumentRedisCache)

	mockRepo.On("ListDocuments", mock.Anything, int32(0), int32(10), "", "", "", "").
		Return([]*pb.DocumentReadModel{}, int64(0), nil)

	handler := &DocumentQueryHandler{
		logger: logger,
		repo:   mockRepo,
		cache:  mockCache,
	}

	req := &pb.ListDocumentsRequest{}
	resp, err := handler.ListDocuments(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)

	mockRepo.AssertCalled(t, "ListDocuments", mock.Anything, int32(0), int32(10), "", "", "", "")
}

// Test SearchDocuments
func TestSearchDocuments_Success(t *testing.T) {
	logger, _ := zap.NewProduction()
	mockRepo := new(MockDocumentRepository)
	mockCache := new(MockDocumentRedisCache)

	expectedDocs := []*pb.DocumentReadModel{
		{Id: "doc-123", DocumentNumber: "INV-001"},
	}

	mockRepo.On("SearchDocuments", mock.Anything, "invoice", int32(0), int32(10), "", "", "").
		Return(expectedDocs, int64(1), nil)

	handler := &DocumentQueryHandler{
		logger: logger,
		repo:   mockRepo,
		cache:  mockCache,
	}

	req := &pb.SearchDocumentsRequest{
		Query:      "invoice",
		Pagination: &pb.PageInfo{Page: 0, Size: 10},
	}
	resp, err := handler.SearchDocuments(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Len(t, resp.Documents, 1)

	mockRepo.AssertCalled(t, "SearchDocuments", mock.Anything, "invoice", int32(0), int32(10), "", "", "")
}

func TestSearchDocuments_EmptyQuery(t *testing.T) {
	logger, _ := zap.NewProduction()
	mockRepo := new(MockDocumentRepository)
	mockCache := new(MockDocumentRedisCache)

	handler := &DocumentQueryHandler{
		logger: logger,
		repo:   mockRepo,
		cache:  mockCache,
	}

	req := &pb.SearchDocumentsRequest{
		Query:      "",
		Pagination: &pb.PageInfo{Page: 0, Size: 10},
	}
	resp, err := handler.SearchDocuments(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, resp)

	mockRepo.AssertNotCalled(t, "SearchDocuments")
}

// Test GetPendingApproval
func TestGetPendingApproval_Success(t *testing.T) {
	logger, _ := zap.NewProduction()
	mockRepo := new(MockDocumentRepository)
	mockCache := new(MockDocumentRedisCache)

	expectedDocs := []*pb.DocumentReadModel{
		{Id: "doc-1", ApprovalStatus: "PENDING"},
		{Id: "doc-2", ApprovalStatus: "PENDING"},
	}

	mockRepo.On("GetPendingApprovalDocuments", mock.Anything, "comp-123", int32(0), int32(10)).
		Return(expectedDocs, int64(2), nil)

	handler := &DocumentQueryHandler{
		logger: logger,
		repo:   mockRepo,
		cache:  mockCache,
	}

	req := &pb.GetPendingApprovalRequest{
		CompanyId:  "comp-123",
		Pagination: &pb.PageInfo{Page: 0, Size: 10},
	}
	resp, err := handler.GetPendingApproval(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Len(t, resp.Documents, 2)
	assert.Equal(t, int64(2), resp.TotalCount)

	mockRepo.AssertCalled(t, "GetPendingApprovalDocuments", mock.Anything, "comp-123", int32(0), int32(10))
}

func TestGetPendingApproval_EmptyCompanyID(t *testing.T) {
	logger, _ := zap.NewProduction()
	mockRepo := new(MockDocumentRepository)
	mockCache := new(MockDocumentRedisCache)

	handler := &DocumentQueryHandler{
		logger: logger,
		repo:   mockRepo,
		cache:  mockCache,
	}

	req := &pb.GetPendingApprovalRequest{
		CompanyId:  "",
		Pagination: &pb.PageInfo{Page: 0, Size: 10},
	}
	resp, err := handler.GetPendingApproval(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, resp)

	mockRepo.AssertNotCalled(t, "GetPendingApprovalDocuments")
}
