package ports

import (
	"fmt"
	"regexp"
	"strings"
	"time"
)

const (
	// MaxFilterLength is the maximum length for filter string fields
	MaxFilterLength = 255
	// MaxSearchTextLength is the maximum length for search text
	MaxSearchTextLength = 500
	// MaxPageSize is the maximum number of items per page
	MaxPageSize = 100
	// MaxInvoiceAmount is the maximum invoice amount for validation
	MaxInvoiceAmount = 999999999999.99
)

var (
	// dateRegex validates YYYY-MM-DD format
	dateRegex = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)
)

// ValidateInvoiceFilter validates filter parameters
func ValidateInvoiceFilter(filter *InvoiceFilter) error {
	if filter == nil {
		return nil
	}

	if len(filter.InvoiceNumber) > MaxFilterLength {
		return fmt.Errorf("invoice_number filter too long: %d characters (max %d)", len(filter.InvoiceNumber), MaxFilterLength)
	}

	if len(filter.SearchText) > MaxSearchTextLength {
		return fmt.Errorf("search_text filter too long: %d characters (max %d)", len(filter.SearchText), MaxSearchTextLength)
	}

	// Validate date format
	if filter.DateFrom != "" {
		if !dateRegex.MatchString(filter.DateFrom) {
			return fmt.Errorf("invalid date_from format: %s (expected YYYY-MM-DD)", filter.DateFrom)
		}
		if _, err := time.Parse("2006-01-02", filter.DateFrom); err != nil {
			return fmt.Errorf("invalid date_from value: %w", err)
		}
	}

	if filter.DateTo != "" {
		if !dateRegex.MatchString(filter.DateTo) {
			return fmt.Errorf("invalid date_to format: %s (expected YYYY-MM-DD)", filter.DateTo)
		}
		if _, err := time.Parse("2006-01-02", filter.DateTo); err != nil {
			return fmt.Errorf("invalid date_to value: %w", err)
		}
	}

	// Validate date range
	if filter.DateFrom != "" && filter.DateTo != "" {
		dateFrom, _ := time.Parse("2006-01-02", filter.DateFrom)
		dateTo, _ := time.Parse("2006-01-02", filter.DateTo)
		if dateFrom.After(dateTo) {
			return fmt.Errorf("date_from cannot be after date_to")
		}
	}

	// Validate amounts
	if filter.MinAmount < 0 {
		return fmt.Errorf("min_amount cannot be negative: %.2f", filter.MinAmount)
	}

	if filter.MaxAmount < 0 {
		return fmt.Errorf("max_amount cannot be negative: %.2f", filter.MaxAmount)
	}

	if filter.MinAmount > MaxInvoiceAmount {
		return fmt.Errorf("min_amount too large: %.2f (max %.2f)", filter.MinAmount, MaxInvoiceAmount)
	}

	if filter.MaxAmount > MaxInvoiceAmount {
		return fmt.Errorf("max_amount too large: %.2f (max %.2f)", filter.MaxAmount, MaxInvoiceAmount)
	}

	if filter.MinAmount > 0 && filter.MaxAmount > 0 && filter.MinAmount > filter.MaxAmount {
		return fmt.Errorf("min_amount cannot be greater than max_amount")
	}

	return nil
}

// ValidateInvoiceSort validates sort parameters
func ValidateInvoiceSort(sort *InvoiceSort) error {
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
		"invoice_date":   true,
		"total_amount":   true,
		"invoice_number": true,
		"created_date":   true,
		"created_at":     true,
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
