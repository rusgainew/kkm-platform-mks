// Файл document-server/internal/application/document/service_test.go содержит реализацию пакета document.
package document

import (
	"context"
	"testing"

	"github.com/rusgainew/kkm-project-mks/document-server/internal/domain"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

// MockRepository для тестирования Service
type MockRepository struct {
	CreateCalled bool
	GetCalled    bool
	UpdateCalled bool
	DeleteCalled bool
	ListCalled   bool
	GetCalls     int
	ListCalls    int
	CreateCalls  int
	UpdateCalls  int
	CreateErr    error
	GetErr       error
	UpdateErr    error
	ListErr      error
	GetResult    map[string]interface{}
	ListResult   []map[string]interface{}
	ListTotal    int64
}

// NewMockRepository создаёт новый экземпляр MockRepository
func NewMockRepository() *MockRepository {
	return &MockRepository{
		GetResult:  make(map[string]interface{}),
		ListResult: make([]map[string]interface{}, 0),
	}
}

func (m *MockRepository) Create(ctx context.Context, id, organizationID, title, content, createdBy string, createdAt int64) error {
	m.CreateCalled = true
	m.CreateCalls++
	return m.CreateErr
}

func (m *MockRepository) Get(ctx context.Context, documentID string) (map[string]interface{}, error) {
	m.GetCalled = true
	m.GetCalls++
	if m.GetErr != nil {
		return nil, m.GetErr
	}
	if m.GetResult == nil {
		m.GetResult = map[string]interface{}{
			"id":                documentID,
			"organization_id":   "org-123",
			"title":             "Test Document",
			"content":           "Test Content",
			"status":            "draft",
			"created_by":        "user-456",
			"created_at":        int64(1000000),
			"updated_at":        int64(1000000),
			"status_changed_at": int64(1000000),
			"version":           int32(1),
		}
	}
	return m.GetResult, nil
}

func (m *MockRepository) Update(ctx context.Context, documentID, title, content string, updatedAt int64) error {
	m.UpdateCalled = true
	m.UpdateCalls++
	return m.UpdateErr
}

func (m *MockRepository) UpdateStatus(ctx context.Context, documentID, status string, statusChangedAt int64) error {
	return nil
}

func (m *MockRepository) List(ctx context.Context, organizationID string, status string, page, perPage int) ([]map[string]interface{}, int64, error) {
	m.ListCalled = true
	m.ListCalls++
	if m.ListErr != nil {
		return nil, 0, m.ListErr
	}
	if m.ListResult == nil {
		m.ListResult = []map[string]interface{}{}
	}
	return m.ListResult, m.ListTotal, nil
}

func (m *MockRepository) Delete(ctx context.Context, documentID string) error {
	m.DeleteCalled = true
	return nil
}

// GetWithVersion реализует интерфейс DocumentRepository
func (m *MockRepository) GetWithVersion(ctx context.Context, documentID string) (map[string]interface{}, error) {
	m.GetCalled = true
	if m.GetErr != nil {
		return nil, m.GetErr
	}
	if m.GetResult == nil {
		m.GetResult = map[string]interface{}{
			"id":                documentID,
			"organization_id":   "org-123",
			"title":             "Test Document",
			"content":           "Test Content",
			"status":            "draft",
			"created_by":        "user-456",
			"created_at":        int64(1000000),
			"updated_at":        int64(1000000),
			"status_changed_at": int64(1000000),
			"version":           1,
		}
	}
	return m.GetResult, nil
}

// UpdateWithVersion реализует интерфейс DocumentRepository
func (m *MockRepository) UpdateWithVersion(ctx context.Context, documentID, title, content string, expectedVersion int) error {
	m.UpdateCalled = true
	if m.UpdateErr != nil {
		return m.UpdateErr
	}
	// Обновляем версию в результате
	if m.GetResult != nil {
		if v, ok := m.GetResult["version"].(int); ok {
			if v == expectedVersion {
				m.GetResult["version"] = v + 1
			} else {
				return domain.ErrVersionConflict
			}
		}
	}
	return nil
}

// MockPublisher для тестирования Service
type MockPublisher struct {
	PublishCalled bool
	PublishErr    error
	Events        []map[string]interface{}
}

func (m *MockPublisher) Publish(ctx context.Context, topic string, event interface{}) error {
	m.PublishCalled = true
	m.Events = append(m.Events, map[string]interface{}{
		"topic": topic,
		"event": event,
	})
	return m.PublishErr
}

func (m *MockPublisher) Close() error {
	return nil
}

// TestCreateDocument тестирует создание документа
func TestCreateDocument(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	mockRepo := &MockRepository{}
	mockPublisher := &MockPublisher{}
	service := NewService(mockRepo, mockPublisher, logger)

	ctx := context.Background()
	organizationID := "550e8400-e29b-41d4-a716-446655440001"
	title := "Test Document"
	content := "Test Content"
	createdBy := "550e8400-e29b-41d4-a716-446655440002"

	docID, err := service.CreateDocument(ctx, organizationID, title, content, createdBy)

	assert.NoError(t, err)
	assert.NotEmpty(t, docID)
	assert.True(t, mockRepo.CreateCalled, "Repository Create should be called")
	assert.True(t, mockPublisher.PublishCalled, "Publisher should publish event")
	assert.Len(t, mockPublisher.Events, 1)
	assert.Equal(t, "document.created", mockPublisher.Events[0]["topic"])
}

// TestCreateDocumentWithInvalidInput тестирует создание документа с невалидными параметрами
func TestCreateDocumentWithInvalidInput(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	mockRepo := &MockRepository{}
	mockPublisher := &MockPublisher{}
	service := NewService(mockRepo, mockPublisher, logger)

	ctx := context.Background()

	tests := []struct {
		name            string
		organizationID  string
		title           string
		content         string
		createdBy       string
		shouldReturnErr bool
	}{
		{
			name:            "Empty organizationID",
			organizationID:  "",
			title:           "Test",
			content:         "Content",
			createdBy:       "550e8400-e29b-41d4-a716-446655440002",
			shouldReturnErr: true,
		},
		{
			name:            "Empty title",
			organizationID:  "550e8400-e29b-41d4-a716-446655440001",
			title:           "",
			content:         "Content",
			createdBy:       "550e8400-e29b-41d4-a716-446655440002",
			shouldReturnErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := service.CreateDocument(ctx, tt.organizationID, tt.title, tt.content, tt.createdBy)
			if tt.shouldReturnErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// TestGetDocument тестирует получение документа
func TestGetDocument(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	mockRepo := &MockRepository{}
	mockPublisher := &MockPublisher{}
	service := NewService(mockRepo, mockPublisher, logger)

	ctx := context.Background()
	documentID := "550e8400-e29b-41d4-a716-446655440003"

	doc, err := service.GetDocument(ctx, documentID)

	assert.NoError(t, err)
	assert.NotNil(t, doc)
	assert.True(t, mockRepo.GetCalled, "Repository Get should be called")
	assert.Equal(t, documentID, doc["id"])
}

// TestGetDocumentNotFound тестирует получение несуществующего документа
func TestGetDocumentNotFound(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	mockRepo := &MockRepository{
		GetErr: domain.ErrInvalidDocumentID,
	}
	mockPublisher := &MockPublisher{}
	service := NewService(mockRepo, mockPublisher, logger)

	ctx := context.Background()
	documentID := "non-existent"

	doc, err := service.GetDocument(ctx, documentID)

	assert.Error(t, err)
	assert.Nil(t, doc)
}

// TestUpdateDocument тестирует обновление документа
func TestUpdateDocument(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	mockRepo := &MockRepository{}
	mockPublisher := &MockPublisher{}
	service := NewService(mockRepo, mockPublisher, logger)

	ctx := context.Background()
	documentID := "550e8400-e29b-41d4-a716-446655440003"
	title := "Updated Title"
	content := "Updated Content"

	err := service.UpdateDocument(ctx, documentID, title, content)

	assert.NoError(t, err)
	assert.True(t, mockRepo.UpdateCalled, "Repository Update should be called")
	assert.True(t, mockPublisher.PublishCalled, "Publisher should publish event")
}

// TestListDocuments тестирует получение списка документов
func TestListDocuments(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	mockRepo := &MockRepository{
		ListResult: []map[string]interface{}{
			{
				"id":    "550e8400-e29b-41d4-a716-446655440001",
				"title": "Document 1",
			},
			{
				"id":    "550e8400-e29b-41d4-a716-446655440002",
				"title": "Document 2",
			},
		},
		ListTotal: 2,
	}
	mockPublisher := &MockPublisher{}
	service := NewService(mockRepo, mockPublisher, logger)

	ctx := context.Background()

	docs, total, err := service.ListDocuments(ctx, "550e8400-e29b-41d4-a716-446655440001", "", 1, 20)

	assert.NoError(t, err)
	assert.Equal(t, int64(2), total)
	assert.Len(t, docs, 2)
	assert.True(t, mockRepo.ListCalled, "Repository List should be called")
}

// TestSendDocument тестирует отправку документа на согласование
func TestSendDocument(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	mockRepo := &MockRepository{}
	mockPublisher := &MockPublisher{}
	service := NewService(mockRepo, mockPublisher, logger)

	ctx := context.Background()
	documentID := "550e8400-e29b-41d4-a716-446655440001"
	recipientID := "550e8400-e29b-41d4-a716-446655440002"
	message := "Please review"

	err := service.SendDocument(ctx, documentID, recipientID, message)

	assert.NoError(t, err)
	assert.True(t, mockPublisher.PublishCalled, "Publisher should publish event")
	assert.Equal(t, "document.sent", mockPublisher.Events[0]["topic"])
}

// TestApproveDocument тестирует одобрение документа
func TestApproveDocument(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	mockRepo := &MockRepository{}
	mockPublisher := &MockPublisher{}
	service := NewService(mockRepo, mockPublisher, logger)

	ctx := context.Background()
	documentID := "550e8400-e29b-41d4-a716-446655440001"
	approvedBy := "550e8400-e29b-41d4-a716-446655440002"

	err := service.ApproveDocument(ctx, documentID, approvedBy, "Approved")

	assert.NoError(t, err)
	assert.True(t, mockPublisher.PublishCalled, "Publisher should publish event")
	assert.Equal(t, "document.approved", mockPublisher.Events[0]["topic"])
}

// TestRejectDocument тестирует отклонение документа
func TestRejectDocument(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	mockRepo := &MockRepository{}
	mockPublisher := &MockPublisher{}
	service := NewService(mockRepo, mockPublisher, logger)

	ctx := context.Background()
	documentID := "550e8400-e29b-41d4-a716-446655440001"
	rejectedBy := "550e8400-e29b-41d4-a716-446655440002"
	reason := "Needs revision"

	err := service.RejectDocument(ctx, documentID, rejectedBy, reason)

	assert.NoError(t, err)
	assert.True(t, mockPublisher.PublishCalled, "Publisher should publish event")
	assert.Equal(t, "document.rejected", mockPublisher.Events[0]["topic"])
}

// TestArchiveDocument тестирует архивирование документа
func TestArchiveDocument(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	mockRepo := &MockRepository{}
	mockPublisher := &MockPublisher{}
	service := NewService(mockRepo, mockPublisher, logger)

	ctx := context.Background()
	documentID := "550e8400-e29b-41d4-a716-446655440001"

	err := service.ArchiveDocument(ctx, documentID)

	assert.NoError(t, err)
	assert.True(t, mockPublisher.PublishCalled, "Publisher should publish event")
	assert.Equal(t, "document.archived", mockPublisher.Events[0]["topic"])
}
