package domain

import "errors"

var (
	ErrDocumentNotFound       = errors.New("document not found")
	ErrInvalidDocumentID      = errors.New("invalid document ID")
	ErrInvalidOrganizationID  = errors.New("invalid organization ID")
	ErrInvalidStatus          = errors.New("invalid status")
	ErrUnauthorized           = errors.New("unauthorized")
	ErrInvalidInput           = errors.New("invalid input")
	ErrVersionConflict        = errors.New("version conflict: document has been modified by another request")
	ErrConcurrentModification = errors.New("concurrent modification detected")
)
