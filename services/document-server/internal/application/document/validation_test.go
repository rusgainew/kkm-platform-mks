package document

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestValidateCreateDocumentRequest тестирует валидацию параметров для создания документа
func TestValidateCreateDocumentRequest(t *testing.T) {
	tests := []struct {
		name               string
		organizationID     string
		title              string
		content            string
		createdBy          string
		expectedErrorCount int
	}{
		{
			name:               "Valid request",
			organizationID:     "550e8400-e29b-41d4-a716-446655440001",
			title:              "Valid Title",
			content:            "Valid content",
			createdBy:          "550e8400-e29b-41d4-a716-446655440002",
			expectedErrorCount: 0,
		},
		{
			name:               "Empty organization ID",
			organizationID:     "",
			title:              "Title",
			content:            "Content",
			createdBy:          "550e8400-e29b-41d4-a716-446655440002",
			expectedErrorCount: 1,
		},
		{
			name:               "Invalid organization ID",
			organizationID:     "invalid-uuid",
			title:              "Title",
			content:            "Content",
			createdBy:          "550e8400-e29b-41d4-a716-446655440002",
			expectedErrorCount: 1,
		},
		{
			name:               "Empty title",
			organizationID:     "550e8400-e29b-41d4-a716-446655440001",
			title:              "",
			content:            "Content",
			createdBy:          "550e8400-e29b-41d4-a716-446655440002",
			expectedErrorCount: 1,
		},
		{
			name:               "Title too short",
			organizationID:     "550e8400-e29b-41d4-a716-446655440001",
			title:              "AB",
			content:            "Content",
			createdBy:          "550e8400-e29b-41d4-a716-446655440002",
			expectedErrorCount: 1,
		},
		{
			name:               "Title too long",
			organizationID:     "550e8400-e29b-41d4-a716-446655440001",
			title:              strings.Repeat("A", 256),
			content:            "Content",
			createdBy:          "550e8400-e29b-41d4-a716-446655440002",
			expectedErrorCount: 1,
		},
		{
			name:               "Empty created by",
			organizationID:     "550e8400-e29b-41d4-a716-446655440001",
			title:              "Title",
			content:            "Content",
			createdBy:          "",
			expectedErrorCount: 1,
		},
		{
			name:               "Multiple errors",
			organizationID:     "invalid",
			title:              "",
			content:            "Content",
			createdBy:          "invalid",
			expectedErrorCount: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errors := ValidateCreateDocumentRequest(
				tt.organizationID, tt.title, tt.content, tt.createdBy)

			assert.Equal(t, tt.expectedErrorCount, len(errors),
				"Expected %d validation errors, got %d", tt.expectedErrorCount, len(errors))
		})
	}
}

// TestValidateUpdateDocumentRequest тестирует валидацию параметров для обновления документа
func TestValidateUpdateDocumentRequest(t *testing.T) {
	tests := []struct {
		name               string
		documentID         string
		title              string
		content            string
		expectedErrorCount int
	}{
		{
			name:               "Valid request",
			documentID:         "550e8400-e29b-41d4-a716-446655440001",
			title:              "Updated Title",
			content:            "Updated content",
			expectedErrorCount: 0,
		},
		{
			name:               "Empty document ID",
			documentID:         "",
			title:              "Title",
			content:            "Content",
			expectedErrorCount: 1,
		},
		{
			name:               "Invalid document ID",
			documentID:         "not-a-uuid",
			title:              "Title",
			content:            "Content",
			expectedErrorCount: 1,
		},
		{
			name:               "Empty title",
			documentID:         "550e8400-e29b-41d4-a716-446655440001",
			title:              "",
			content:            "Content",
			expectedErrorCount: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errors := ValidateUpdateDocumentRequest(tt.documentID, tt.title, tt.content)
			assert.Equal(t, tt.expectedErrorCount, len(errors))
		})
	}
}

// TestValidateSendDocumentRequest тестирует валидацию параметров для отправки документа
func TestValidateSendDocumentRequest(t *testing.T) {
	tests := []struct {
		name               string
		documentID         string
		recipientID        string
		message            string
		expectedErrorCount int
	}{
		{
			name:               "Valid request",
			documentID:         "550e8400-e29b-41d4-a716-446655440001",
			recipientID:        "550e8400-e29b-41d4-a716-446655440002",
			message:            "Please review",
			expectedErrorCount: 0,
		},
		{
			name:               "Empty document ID",
			documentID:         "",
			recipientID:        "550e8400-e29b-41d4-a716-446655440002",
			message:            "Message",
			expectedErrorCount: 1,
		},
		{
			name:               "Invalid recipient ID",
			documentID:         "550e8400-e29b-41d4-a716-446655440001",
			recipientID:        "invalid",
			message:            "Message",
			expectedErrorCount: 1,
		},
		{
			name:               "Message too long",
			documentID:         "550e8400-e29b-41d4-a716-446655440001",
			recipientID:        "550e8400-e29b-41d4-a716-446655440002",
			message:            strings.Repeat("A", 5001),
			expectedErrorCount: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errors := ValidateSendDocumentRequest(tt.documentID, tt.recipientID, tt.message)
			assert.Equal(t, tt.expectedErrorCount, len(errors))
		})
	}
}

// TestValidateListDocumentsRequest тестирует валидацию параметров для получения списка документов
func TestValidateListDocumentsRequest(t *testing.T) {
	tests := []struct {
		name               string
		organizationID     string
		page               int
		perPage            int
		expectedErrorCount int
	}{
		{
			name:               "Valid request",
			organizationID:     "550e8400-e29b-41d4-a716-446655440001",
			page:               1,
			perPage:            20,
			expectedErrorCount: 0,
		},
		{
			name:               "Empty organization ID",
			organizationID:     "",
			page:               1,
			perPage:            20,
			expectedErrorCount: 1,
		},
		{
			name:               "Invalid page number",
			organizationID:     "550e8400-e29b-41d4-a716-446655440001",
			page:               0,
			perPage:            20,
			expectedErrorCount: 1,
		},
		{
			name:               "Invalid per page",
			organizationID:     "550e8400-e29b-41d4-a716-446655440001",
			page:               1,
			perPage:            101,
			expectedErrorCount: 1,
		},
		{
			name:               "Per page zero",
			organizationID:     "550e8400-e29b-41d4-a716-446655440001",
			page:               1,
			perPage:            0,
			expectedErrorCount: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errors := ValidateListDocumentsRequest(tt.organizationID, tt.page, tt.perPage)
			assert.Equal(t, tt.expectedErrorCount, len(errors))
		})
	}
}

// TestValidateStatus тестирует валидацию значения статуса
func TestValidateStatus(t *testing.T) {
	validStatuses := []string{"draft", "sent", "approved", "rejected", "archived"}
	invalidStatuses := []string{"unknown", "pending", "completed", ""}

	for _, status := range validStatuses {
		assert.True(t, ValidateStatus(status), "Status %s should be valid", status)
	}

	for _, status := range invalidStatuses {
		assert.False(t, ValidateStatus(status), "Status %s should be invalid", status)
	}
}

// TestFormatValidationErrors тестирует форматирование ошибок валидации
func TestFormatValidationErrors(t *testing.T) {
	tests := []struct {
		name     string
		errors   []ValidationError
		expected string
	}{
		{
			name:     "No errors",
			errors:   []ValidationError{},
			expected: "",
		},
		{
			name: "Single error",
			errors: []ValidationError{
				{Field: "title", Message: "required"},
			},
			expected: "validation failed:\n  - title: required\n",
		},
		{
			name: "Multiple errors",
			errors: []ValidationError{
				{Field: "title", Message: "required"},
				{Field: "content", Message: "too long"},
			},
			expected: "validation failed:\n  - title: required\n  - content: too long\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatValidationErrors(tt.errors)
			assert.Equal(t, tt.expected, result)
		})
	}
}
