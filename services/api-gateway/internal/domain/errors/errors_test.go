package errors

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewValidationError(t *testing.T) {
	err := NewValidationError("invalid input")

	assert.NotNil(t, err)
	assert.Equal(t, "VALIDATION_ERROR", err.Code)
	assert.Equal(t, "Validation failed", err.Message)
	assert.Equal(t, "invalid input", err.Details)
	assert.Equal(t, http.StatusBadRequest, err.StatusCode)
}

func TestNewNotFoundError(t *testing.T) {
	err := NewNotFoundError("User not found")

	assert.NotNil(t, err)
	assert.Equal(t, "NOT_FOUND", err.Code)
	assert.Equal(t, "User not found", err.Message)
	assert.Equal(t, http.StatusNotFound, err.StatusCode)
}

func TestNewUnauthorizedError(t *testing.T) {
	err := NewUnauthorizedError("Invalid credentials")

	assert.NotNil(t, err)
	assert.Equal(t, "UNAUTHORIZED", err.Code)
	assert.Equal(t, "Invalid credentials", err.Message)
	assert.Equal(t, http.StatusUnauthorized, err.StatusCode)
}

func TestNewInternalServerError(t *testing.T) {
	err := NewInternalServerError()

	assert.NotNil(t, err)
	assert.Equal(t, "INTERNAL_SERVER_ERROR", err.Code)
	assert.Equal(t, "Internal server error", err.Message)
	assert.Equal(t, http.StatusInternalServerError, err.StatusCode)
}

func TestAPIErrorError(t *testing.T) {
	err := NewValidationError("test error")
	result := err.Error()

	assert.NotNil(t, result)
	assert.Contains(t, result, "Validation failed")
	assert.Contains(t, result, "test error")
}

func TestAPIErrorErrorWithoutDetails(t *testing.T) {
	err := NewInternalServerError()
	result := err.Error()

	assert.Equal(t, "Internal server error", result)
}
