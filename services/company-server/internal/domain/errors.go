// Файл company-server/internal/domain/errors.go содержит реализацию пакета domain.
package domain

import "errors"

var (
	// Organization errors
	ErrOrganizationNotFound     = errors.New("organization not found")
	ErrOrganizationExists       = errors.New("organization already exists")
	ErrInvalidOrganizationID    = errors.New("invalid organization ID")
	ErrInvalidOrganizationName  = errors.New("invalid organization name")
	ErrOrganizationNameRequired = errors.New("organization name is required")

	// Employee errors
	ErrEmployeeNotFound  = errors.New("employee not found")
	ErrEmployeeExists    = errors.New("employee already exists in organization")
	ErrInvalidEmployeeID = errors.New("invalid employee ID")
	ErrInvalidUserID     = errors.New("invalid user ID")
	ErrInvalidRole       = errors.New("invalid employee role")

	// Permission errors
	ErrUnauthorized     = errors.New("unauthorized access")
	ErrPermissionDenied = errors.New("permission denied")

	// Validation errors
	ErrInvalidInput         = errors.New("invalid input")
	ErrMissingRequiredField = errors.New("missing required field")
)
