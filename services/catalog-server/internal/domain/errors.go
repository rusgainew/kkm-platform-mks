// Файл catalog-server/internal/domain/errors.go содержит реализацию пакета domain.
package domain

import "errors"

var (
	// Catalog errors
	ErrCatalogItemNotFound      = errors.New("catalog item not found")
	ErrCatalogItemAlreadyExists = errors.New("catalog item with this code already exists")
	ErrInvalidCatalogItem       = errors.New("invalid catalog item data")

	// Validation errors
	ErrInvalidName     = errors.New("invalid item name")
	ErrInvalidCode     = errors.New("invalid item code")
	ErrInvalidPrice    = errors.New("invalid price: must be non-negative")
	ErrInvalidVATRate  = errors.New("invalid VAT rate")
	ErrInvalidUnitType = errors.New("invalid unit type")
	ErrInvalidCategory = errors.New("invalid category")

	// Permission errors
	ErrUnauthorized            = errors.New("unauthorized access")
	ErrInsufficientPermissions = errors.New("insufficient permissions")

	// Organization errors
	ErrOrganizationNotFound = errors.New("organization not found")
)
