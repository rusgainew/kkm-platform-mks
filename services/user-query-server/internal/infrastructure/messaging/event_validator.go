package messaging

import (
	"fmt"
	"time"

	"go.uber.org/zap"
)

// EventValidationError represents validation error with details
type EventValidationError struct {
	Field   string
	Message string
}

// EventValidator validates received events
type EventValidator struct {
	logger *zap.Logger
}

// NewEventValidator creates a new event validator
func NewEventValidator(logger *zap.Logger) *EventValidator {
	return &EventValidator{
		logger: logger,
	}
}

// ValidateUserEvent validates a user domain event
func (v *EventValidator) ValidateUserEvent(event UserEvent) []EventValidationError {
	var errs []EventValidationError

	// Validate required fields
	if event.ID == "" {
		errs = append(errs, EventValidationError{
			Field:   "id",
			Message: "event id is required",
		})
	}

	if event.Type == "" {
		errs = append(errs, EventValidationError{
			Field:   "type",
			Message: "event type is required",
		})
	}

	if event.Timestamp == "" {
		errs = append(errs, EventValidationError{
			Field:   "timestamp",
			Message: "timestamp is required",
		})
	} else {
		// Validate timestamp format (RFC3339 or similar)
		if _, err := time.Parse(time.RFC3339, event.Timestamp); err != nil {
			// Try ISO8601-like format
			if _, err := time.Parse("2006-01-02T15:04:05Z07:00", event.Timestamp); err != nil {
				errs = append(errs, EventValidationError{
					Field:   "timestamp",
					Message: fmt.Sprintf("invalid timestamp format: %v", err),
				})
			}
		}
	}

	// Validate event type is known
	validTypes := map[string]bool{
		"user.registered": true,
		"user.logged_in":  true,
		"token.refreshed": true,
	}
	if !validTypes[event.Type] && event.Type != "" {
		errs = append(errs, EventValidationError{
			Field:   "type",
			Message: fmt.Sprintf("unknown event type: %s", event.Type),
		})
	}

	// Validate data payload exists
	if event.Data == nil || len(event.Data) == 0 {
		errs = append(errs, EventValidationError{
			Field:   "data",
			Message: "event data payload is required",
		})
	} else {
		// Validate data based on event type
		switch event.Type {
		case "user.registered":
			errs = append(errs, v.validateUserRegisteredData(event.Data)...)
		case "user.logged_in":
			errs = append(errs, v.validateUserLoggedInData(event.Data)...)
		case "token.refreshed":
			errs = append(errs, v.validateTokenRefreshedData(event.Data)...)
		}
	}

	return errs
}

// validateUserRegisteredData validates data for user.registered event
func (v *EventValidator) validateUserRegisteredData(data map[string]interface{}) []EventValidationError {
	var errs []EventValidationError

	// Required: user_id, email, created_at
	if _, exists := data["user_id"]; !exists {
		errs = append(errs, EventValidationError{
			Field:   "data.user_id",
			Message: "user_id is required for user.registered event",
		})
	} else if userID, ok := data["user_id"].(string); !ok || userID == "" {
		errs = append(errs, EventValidationError{
			Field:   "data.user_id",
			Message: "user_id must be a non-empty string",
		})
	}

	if _, exists := data["email"]; !exists {
		errs = append(errs, EventValidationError{
			Field:   "data.email",
			Message: "email is required for user.registered event",
		})
	} else if email, ok := data["email"].(string); !ok || email == "" {
		errs = append(errs, EventValidationError{
			Field:   "data.email",
			Message: "email must be a non-empty string",
		})
	} else if !v.isValidEmail(email) {
		errs = append(errs, EventValidationError{
			Field:   "data.email",
			Message: "invalid email format",
		})
	}

	if _, exists := data["created_at"]; !exists {
		errs = append(errs, EventValidationError{
			Field:   "data.created_at",
			Message: "created_at is required for user.registered event",
		})
	}

	return errs
}

// validateUserLoggedInData validates data for user.logged_in event
func (v *EventValidator) validateUserLoggedInData(data map[string]interface{}) []EventValidationError {
	var errs []EventValidationError

	// Required: user_id, login_at
	if _, exists := data["user_id"]; !exists {
		errs = append(errs, EventValidationError{
			Field:   "data.user_id",
			Message: "user_id is required for user.logged_in event",
		})
	} else if userID, ok := data["user_id"].(string); !ok || userID == "" {
		errs = append(errs, EventValidationError{
			Field:   "data.user_id",
			Message: "user_id must be a non-empty string",
		})
	}

	if _, exists := data["login_at"]; !exists {
		errs = append(errs, EventValidationError{
			Field:   "data.login_at",
			Message: "login_at is required for user.logged_in event",
		})
	}

	return errs
}

// validateTokenRefreshedData validates data for token.refreshed event
func (v *EventValidator) validateTokenRefreshedData(data map[string]interface{}) []EventValidationError {
	var errs []EventValidationError

	// Required: user_id, token_issued_at
	if _, exists := data["user_id"]; !exists {
		errs = append(errs, EventValidationError{
			Field:   "data.user_id",
			Message: "user_id is required for token.refreshed event",
		})
	} else if userID, ok := data["user_id"].(string); !ok || userID == "" {
		errs = append(errs, EventValidationError{
			Field:   "data.user_id",
			Message: "user_id must be a non-empty string",
		})
	}

	if _, exists := data["token_issued_at"]; !exists {
		errs = append(errs, EventValidationError{
			Field:   "data.token_issued_at",
			Message: "token_issued_at is required for token.refreshed event",
		})
	}

	return errs
}

// isValidEmail performs basic email validation
func (v *EventValidator) isValidEmail(email string) bool {
	// Simple email validation (checks for @ and basic structure)
	if len(email) < 3 || len(email) > 254 {
		return false
	}

	atCount := 0
	for _, c := range email {
		if c == '@' {
			atCount++
		}
	}

	return atCount == 1 && email[0] != '@' && email[len(email)-1] != '@'
}

// LogValidationErrors logs validation errors with appropriate level
func (v *EventValidator) LogValidationErrors(eventID string, errs []EventValidationError) {
	if len(errs) == 0 {
		return
	}

	fields := make([]zap.Field, len(errs))
	for i, err := range errs {
		fields[i] = zap.String(err.Field, err.Message)
	}

	v.logger.Warn("Event validation errors",
		append([]zap.Field{
			zap.String("event_id", eventID),
			zap.Int("error_count", len(errs)),
		}, fields...)...,
	)
}
