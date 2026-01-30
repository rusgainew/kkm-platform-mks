// Файл document-server/internal/application/document/validation.go содержит реализацию пакета document.
package document

import (
	"fmt"
	"regexp"
)

var (
	// UUID v4 regex pattern
	uuidRegex = regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)
)

// ValidationError представляет ошибку валидации одного поля
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// ValidateCreateDocumentRequest валидирует параметры для создания документа
func ValidateCreateDocumentRequest(organizationID, title, content, createdBy string) []ValidationError {
	var errors []ValidationError

	// organizationID validation
	if organizationID == "" {
		errors = append(errors, ValidationError{
			Field:   "organization_id",
			Message: "required",
		})
	} else if !isValidUUID(organizationID) {
		errors = append(errors, ValidationError{
			Field:   "organization_id",
			Message: "must be a valid UUID",
		})
	}

	// title validation
	if title == "" {
		errors = append(errors, ValidationError{
			Field:   "title",
			Message: "required",
		})
	} else if len(title) < 3 {
		errors = append(errors, ValidationError{
			Field:   "title",
			Message: "must be at least 3 characters",
		})
	} else if len(title) > 255 {
		errors = append(errors, ValidationError{
			Field:   "title",
			Message: "must not exceed 255 characters",
		})
	}

	// content validation (optional but max length)
	if len(content) > 10000000 { // 10MB max
		errors = append(errors, ValidationError{
			Field:   "content",
			Message: "must not exceed 10MB",
		})
	}

	// createdBy validation
	if createdBy == "" {
		errors = append(errors, ValidationError{
			Field:   "created_by",
			Message: "required",
		})
	} else if !isValidUUID(createdBy) {
		errors = append(errors, ValidationError{
			Field:   "created_by",
			Message: "must be a valid UUID",
		})
	}

	return errors
}

// ValidateUpdateDocumentRequest валидирует параметры для обновления документа
func ValidateUpdateDocumentRequest(documentID, title, content string) []ValidationError {
	var errors []ValidationError

	// documentID validation
	if documentID == "" {
		errors = append(errors, ValidationError{
			Field:   "id",
			Message: "required",
		})
	} else if !isValidUUID(documentID) {
		errors = append(errors, ValidationError{
			Field:   "id",
			Message: "must be a valid UUID",
		})
	}

	// title validation
	if title == "" {
		errors = append(errors, ValidationError{
			Field:   "title",
			Message: "required",
		})
	} else if len(title) < 3 {
		errors = append(errors, ValidationError{
			Field:   "title",
			Message: "must be at least 3 characters",
		})
	} else if len(title) > 255 {
		errors = append(errors, ValidationError{
			Field:   "title",
			Message: "must not exceed 255 characters",
		})
	}

	// content validation (optional but max length)
	if len(content) > 10000000 { // 10MB max
		errors = append(errors, ValidationError{
			Field:   "content",
			Message: "must not exceed 10MB",
		})
	}

	return errors
}

// ValidateSendDocumentRequest валидирует параметры для отправки документа
func ValidateSendDocumentRequest(documentID, recipientID, message string) []ValidationError {
	var errors []ValidationError

	// documentID validation
	if documentID == "" {
		errors = append(errors, ValidationError{
			Field:   "document_id",
			Message: "required",
		})
	} else if !isValidUUID(documentID) {
		errors = append(errors, ValidationError{
			Field:   "document_id",
			Message: "must be a valid UUID",
		})
	}

	// recipientID validation
	if recipientID == "" {
		errors = append(errors, ValidationError{
			Field:   "recipient_id",
			Message: "required",
		})
	} else if !isValidUUID(recipientID) {
		errors = append(errors, ValidationError{
			Field:   "recipient_id",
			Message: "must be a valid UUID",
		})
	}

	// message validation (optional but max length)
	if len(message) > 5000 {
		errors = append(errors, ValidationError{
			Field:   "message",
			Message: "must not exceed 5000 characters",
		})
	}

	return errors
}

// ValidateApproveDocumentRequest валидирует параметры для одобрения документа
func ValidateApproveDocumentRequest(documentID, approvedBy string) []ValidationError {
	var errors []ValidationError

	// documentID validation
	if documentID == "" {
		errors = append(errors, ValidationError{
			Field:   "document_id",
			Message: "required",
		})
	} else if !isValidUUID(documentID) {
		errors = append(errors, ValidationError{
			Field:   "document_id",
			Message: "must be a valid UUID",
		})
	}

	// approvedBy validation
	if approvedBy == "" {
		errors = append(errors, ValidationError{
			Field:   "approved_by",
			Message: "required",
		})
	} else if !isValidUUID(approvedBy) {
		errors = append(errors, ValidationError{
			Field:   "approved_by",
			Message: "must be a valid UUID",
		})
	}

	return errors
}

// ValidateRejectDocumentRequest валидирует параметры для отклонения документа
func ValidateRejectDocumentRequest(documentID, rejectedBy, reason string) []ValidationError {
	var errors []ValidationError

	// documentID validation
	if documentID == "" {
		errors = append(errors, ValidationError{
			Field:   "document_id",
			Message: "required",
		})
	} else if !isValidUUID(documentID) {
		errors = append(errors, ValidationError{
			Field:   "document_id",
			Message: "must be a valid UUID",
		})
	}

	// rejectedBy validation
	if rejectedBy == "" {
		errors = append(errors, ValidationError{
			Field:   "rejected_by",
			Message: "required",
		})
	} else if !isValidUUID(rejectedBy) {
		errors = append(errors, ValidationError{
			Field:   "rejected_by",
			Message: "must be a valid UUID",
		})
	}

	// reason validation (optional but max length)
	if len(reason) > 1000 {
		errors = append(errors, ValidationError{
			Field:   "reason",
			Message: "must not exceed 1000 characters",
		})
	}

	return errors
}

// ValidateArchiveDocumentRequest валидирует параметры для архивирования документа
func ValidateArchiveDocumentRequest(documentID string) []ValidationError {
	var errors []ValidationError

	// documentID validation
	if documentID == "" {
		errors = append(errors, ValidationError{
			Field:   "document_id",
			Message: "required",
		})
	} else if !isValidUUID(documentID) {
		errors = append(errors, ValidationError{
			Field:   "document_id",
			Message: "must be a valid UUID",
		})
	}

	return errors
}

// ValidateListDocumentsRequest валидирует параметры для получения списка документов
func ValidateListDocumentsRequest(organizationID string, page, perPage int) []ValidationError {
	var errors []ValidationError

	// organizationID validation
	if organizationID == "" {
		errors = append(errors, ValidationError{
			Field:   "organization_id",
			Message: "required",
		})
	} else if !isValidUUID(organizationID) {
		errors = append(errors, ValidationError{
			Field:   "organization_id",
			Message: "must be a valid UUID",
		})
	}

	// page validation
	if page < 1 {
		errors = append(errors, ValidationError{
			Field:   "page",
			Message: "must be >= 1",
		})
	}

	// perPage validation
	if perPage < 1 {
		errors = append(errors, ValidationError{
			Field:   "per_page",
			Message: "must be >= 1",
		})
	} else if perPage > 100 {
		errors = append(errors, ValidationError{
			Field:   "per_page",
			Message: "must not exceed 100 (max allowed)",
		})
	}

	return errors
}

// ValidateStatus валидирует значение статуса документа
func ValidateStatus(status string) bool {
	validStatuses := map[string]bool{
		"draft":    true,
		"sent":     true,
		"approved": true,
		"rejected": true,
		"archived": true,
	}
	return validStatuses[status]
}

// isValidUUID проверяет, является ли строка валидным UUID
func isValidUUID(id string) bool {
	return uuidRegex.MatchString(id)
}

// FormatValidationErrors форматирует ошибки валидации в строку
func FormatValidationErrors(errors []ValidationError) string {
	if len(errors) == 0 {
		return ""
	}

	msg := "validation failed:\n"
	for _, err := range errors {
		msg += fmt.Sprintf("  - %s: %s\n", err.Field, err.Message)
	}
	return msg
}
