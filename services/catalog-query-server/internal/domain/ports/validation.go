package ports

import (
	"fmt"
	"strings"
)

const (
	// MaxFilterLength is the maximum length for filter string fields
	MaxFilterLength = 255
	// MaxSearchTextLength is the maximum length for search text
	MaxSearchTextLength = 500
	// MaxPageSize is the maximum number of items per page
	MaxPageSize = 100
)

// ValidateCatalogFilter validates filter parameters
func ValidateCatalogFilter(filter *CatalogFilter) error {
	if filter == nil {
		return nil
	}

	if len(filter.Name) > MaxFilterLength {
		return fmt.Errorf("name filter too long: %d characters (max %d)", len(filter.Name), MaxFilterLength)
	}

	if len(filter.Number) > MaxFilterLength {
		return fmt.Errorf("number filter too long: %d characters (max %d)", len(filter.Number), MaxFilterLength)
	}

	if len(filter.TnvedCode) > MaxFilterLength {
		return fmt.Errorf("tnved_code filter too long: %d characters (max %d)", len(filter.TnvedCode), MaxFilterLength)
	}

	if len(filter.GkedCode) > MaxFilterLength {
		return fmt.Errorf("gked_code filter too long: %d characters (max %d)", len(filter.GkedCode), MaxFilterLength)
	}

	if len(filter.SearchText) > MaxSearchTextLength {
		return fmt.Errorf("search_text filter too long: %d characters (max %d)", len(filter.SearchText), MaxSearchTextLength)
	}

	return nil
}

// ValidateCatalogSort validates sort parameters
func ValidateCatalogSort(sort *CatalogSort) error {
	if sort == nil {
		return nil
	}

	// Validate sort order
	orderUpper := strings.ToUpper(string(sort.Order))
	if orderUpper != "" && orderUpper != "ASC" && orderUpper != "DESC" {
		return fmt.Errorf("invalid sort order: %s (must be ASC or DESC)", sort.Order)
	}

	// Validate sort field
	validFields := map[string]bool{
		"name":       true,
		"number":     true,
		"code":       true,
		"tnved_code": true,
		"tnved":      true,
		"gked_code":  true,
		"gked":       true,
	}

	if sort.Field != "" && !validFields[strings.ToLower(sort.Field)] {
		return fmt.Errorf("invalid sort field: %s", sort.Field)
	}

	return nil
}

// ValidatePagination validates pagination parameters
func ValidatePagination(page, size int32) (int32, int32, error) {
	if page < 1 {
		page = 1
	}

	if size < 1 {
		size = 10
	} else if size > MaxPageSize {
		return 0, 0, fmt.Errorf("page size too large: %d (max %d)", size, MaxPageSize)
	}

	return page, size, nil
}
