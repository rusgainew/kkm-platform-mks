package messaging

import (
	"fmt"
	"time"

	"go.uber.org/zap"
)

// DocumentValidationError represents validation error with details
type DocumentValidationError struct {
	Field   string
	Message string
}

// DocumentEventValidator validates document domain events
type DocumentEventValidator struct {
	logger *zap.Logger
}

// NewDocumentEventValidator creates a new document event validator
func NewDocumentEventValidator(logger *zap.Logger) *DocumentEventValidator {
	return &DocumentEventValidator{
		logger: logger,
	}
}

// ValidateDocumentEvent validates a document domain event
func (v *DocumentEventValidator) ValidateDocumentEvent(event DocumentEvent) []DocumentValidationError {
	var errs []DocumentValidationError

	// Validate required fields
	if event.ID == "" {
		errs = append(errs, DocumentValidationError{
			Field:   "id",
			Message: "event id is required",
		})
	}

	if event.Type == "" {
		errs = append(errs, DocumentValidationError{
			Field:   "type",
			Message: "event type is required",
		})
	}

	if event.Timestamp == "" {
		errs = append(errs, DocumentValidationError{
			Field:   "timestamp",
			Message: "timestamp is required",
		})
	} else {
		// Validate timestamp format (RFC3339 or similar)
		if _, err := time.Parse(time.RFC3339, event.Timestamp); err != nil {
			// Try ISO8601-like format
			if _, err := time.Parse("2006-01-02T15:04:05Z07:00", event.Timestamp); err != nil {
				errs = append(errs, DocumentValidationError{
					Field:   "timestamp",
					Message: fmt.Sprintf("invalid timestamp format: %v", err),
				})
			}
		}
	}

	// Validate event type is known
	validTypes := map[string]bool{
		"document.created":  true,
		"document.sent":     true,
		"document.approved": true,
		"document.rejected": true,
		"document.archived": true,
	}
	if !validTypes[event.Type] && event.Type != "" {
		errs = append(errs, DocumentValidationError{
			Field:   "type",
			Message: fmt.Sprintf("unknown event type: %s", event.Type),
		})
	}

	// Validate data payload exists
	if event.Data == nil || len(event.Data) == 0 {
		errs = append(errs, DocumentValidationError{
			Field:   "data",
			Message: "event data payload is required",
		})
	} else {
		// Validate data based on event type
		switch event.Type {
		case "document.created":
			errs = append(errs, v.validateDocumentCreatedData(event.Data)...)
		case "document.sent":
			errs = append(errs, v.validateDocumentSentData(event.Data)...)
		case "document.approved":
			errs = append(errs, v.validateDocumentApprovedData(event.Data)...)
		case "document.rejected":
			errs = append(errs, v.validateDocumentRejectedData(event.Data)...)
		case "document.archived":
			errs = append(errs, v.validateDocumentArchivedData(event.Data)...)
		}
	}

	return errs
}

// validateDocumentCreatedData validates data for document.created event
func (v *DocumentEventValidator) validateDocumentCreatedData(data map[string]interface{}) []DocumentValidationError {
	var errs []DocumentValidationError

	// Required: document_id, document_number, user_id, created_at
	if _, exists := data["document_id"]; !exists {
		errs = append(errs, DocumentValidationError{
			Field:   "data.document_id",
			Message: "document_id is required for document.created event",
		})
	} else if docID, ok := data["document_id"].(string); !ok || docID == "" {
		errs = append(errs, DocumentValidationError{
			Field:   "data.document_id",
			Message: "document_id must be a non-empty string",
		})
	}

	if _, exists := data["document_number"]; !exists {
		errs = append(errs, DocumentValidationError{
			Field:   "data.document_number",
			Message: "document_number is required for document.created event",
		})
	}

	if _, exists := data["user_id"]; !exists {
		errs = append(errs, DocumentValidationError{
			Field:   "data.user_id",
			Message: "user_id is required for document.created event",
		})
	} else if userID, ok := data["user_id"].(string); !ok || userID == "" {
		errs = append(errs, DocumentValidationError{
			Field:   "data.user_id",
			Message: "user_id must be a non-empty string",
		})
	}

	if _, exists := data["created_at"]; !exists {
		errs = append(errs, DocumentValidationError{
			Field:   "data.created_at",
			Message: "created_at is required for document.created event",
		})
	}

	return errs
}

// validateDocumentSentData validates data for document.sent event
func (v *DocumentEventValidator) validateDocumentSentData(data map[string]interface{}) []DocumentValidationError {
	var errs []DocumentValidationError

	// Required: document_id, sent_at, recipient
	if _, exists := data["document_id"]; !exists {
		errs = append(errs, DocumentValidationError{
			Field:   "data.document_id",
			Message: "document_id is required for document.sent event",
		})
	} else if docID, ok := data["document_id"].(string); !ok || docID == "" {
		errs = append(errs, DocumentValidationError{
			Field:   "data.document_id",
			Message: "document_id must be a non-empty string",
		})
	}

	if _, exists := data["sent_at"]; !exists {
		errs = append(errs, DocumentValidationError{
			Field:   "data.sent_at",
			Message: "sent_at is required for document.sent event",
		})
	}

	if _, exists := data["recipient"]; !exists {
		errs = append(errs, DocumentValidationError{
			Field:   "data.recipient",
			Message: "recipient is required for document.sent event",
		})
	}

	return errs
}

// validateDocumentApprovedData validates data for document.approved event
func (v *DocumentEventValidator) validateDocumentApprovedData(data map[string]interface{}) []DocumentValidationError {
	var errs []DocumentValidationError

	// Required: document_id, approved_at, approver_id
	if _, exists := data["document_id"]; !exists {
		errs = append(errs, DocumentValidationError{
			Field:   "data.document_id",
			Message: "document_id is required for document.approved event",
		})
	} else if docID, ok := data["document_id"].(string); !ok || docID == "" {
		errs = append(errs, DocumentValidationError{
			Field:   "data.document_id",
			Message: "document_id must be a non-empty string",
		})
	}

	if _, exists := data["approved_at"]; !exists {
		errs = append(errs, DocumentValidationError{
			Field:   "data.approved_at",
			Message: "approved_at is required for document.approved event",
		})
	}

	if _, exists := data["approver_id"]; !exists {
		errs = append(errs, DocumentValidationError{
			Field:   "data.approver_id",
			Message: "approver_id is required for document.approved event",
		})
	}

	return errs
}

// validateDocumentRejectedData validates data for document.rejected event
func (v *DocumentEventValidator) validateDocumentRejectedData(data map[string]interface{}) []DocumentValidationError {
	var errs []DocumentValidationError

	// Required: document_id, rejected_at, rejection_reason
	if _, exists := data["document_id"]; !exists {
		errs = append(errs, DocumentValidationError{
			Field:   "data.document_id",
			Message: "document_id is required for document.rejected event",
		})
	} else if docID, ok := data["document_id"].(string); !ok || docID == "" {
		errs = append(errs, DocumentValidationError{
			Field:   "data.document_id",
			Message: "document_id must be a non-empty string",
		})
	}

	if _, exists := data["rejected_at"]; !exists {
		errs = append(errs, DocumentValidationError{
			Field:   "data.rejected_at",
			Message: "rejected_at is required for document.rejected event",
		})
	}

	if _, exists := data["rejection_reason"]; !exists {
		errs = append(errs, DocumentValidationError{
			Field:   "data.rejection_reason",
			Message: "rejection_reason is required for document.rejected event",
		})
	}

	return errs
}

// validateDocumentArchivedData validates data for document.archived event
func (v *DocumentEventValidator) validateDocumentArchivedData(data map[string]interface{}) []DocumentValidationError {
	var errs []DocumentValidationError

	// Required: document_id, archived_at
	if _, exists := data["document_id"]; !exists {
		errs = append(errs, DocumentValidationError{
			Field:   "data.document_id",
			Message: "document_id is required for document.archived event",
		})
	} else if docID, ok := data["document_id"].(string); !ok || docID == "" {
		errs = append(errs, DocumentValidationError{
			Field:   "data.document_id",
			Message: "document_id must be a non-empty string",
		})
	}

	if _, exists := data["archived_at"]; !exists {
		errs = append(errs, DocumentValidationError{
			Field:   "data.archived_at",
			Message: "archived_at is required for document.archived event",
		})
	}

	return errs
}

// LogValidationErrors logs validation errors with appropriate level
func (v *DocumentEventValidator) LogValidationErrors(eventID string, errs []DocumentValidationError) {
	if len(errs) == 0 {
		return
	}

	fields := make([]zap.Field, len(errs))
	for i, err := range errs {
		fields[i] = zap.String(err.Field, err.Message)
	}

	v.logger.Warn("Document event validation errors",
		append([]zap.Field{
			zap.String("event_id", eventID),
			zap.Int("error_count", len(errs)),
		}, fields...)...,
	)
}
