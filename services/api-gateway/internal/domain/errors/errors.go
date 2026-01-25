package errors

import "errors"

var (
	// Common errors
	ErrInvalidInput         = errors.New("invalid input")
	ErrMissingRequiredField = errors.New("missing required field")
	ErrNotFound             = errors.New("resource not found")
	ErrUnauthorized         = errors.New("unauthorized access")
	ErrForbidden            = errors.New("forbidden")
	ErrInternalServer       = errors.New("internal server error")
	ErrBadGateway           = errors.New("bad gateway")
	ErrServiceUnavailable   = errors.New("service unavailable")

	// Validation errors
	ErrInvalidUUID       = errors.New("invalid UUID format")
	ErrInvalidEmail      = errors.New("invalid email format")
	ErrInvalidPagination = errors.New("invalid pagination parameters")
	ErrStringTooShort    = errors.New("string too short")
	ErrStringTooLong     = errors.New("string too long")

	// Company errors
	ErrCompanyNotFound = errors.New("company not found")
	ErrCompanyExists   = errors.New("company already exists")

	// Invoice errors
	ErrInvoiceNotFound = errors.New("invoice not found")
	ErrInvoiceExists   = errors.New("invoice already exists")

	// Catalog errors
	ErrCatalogNotFound = errors.New("catalog item not found")
	ErrCatalogExists   = errors.New("catalog item already exists")

	// Bank account errors
	ErrBankAccountNotFound = errors.New("bank account not found")
	ErrBankAccountExists   = errors.New("bank account already exists")
)

// APIError структура для HTTP API ошибок
type APIError struct {
	Code       string `json:"code"`
	Message    string `json:"message"`
	Details    string `json:"details,omitempty"`
	StatusCode int    `json:"-"`
}

func (e *APIError) Error() string {
	if e.Details != "" {
		return e.Message + ": " + e.Details
	}
	return e.Message
}

// NewAPIError создает новую API ошибку
func NewAPIError(code, message string, statusCode int) *APIError {
	return &APIError{
		Code:       code,
		Message:    message,
		StatusCode: statusCode,
	}
}

// NewAPIErrorWithDetails создает новую API ошибку с деталями
func NewAPIErrorWithDetails(code, message, details string, statusCode int) *APIError {
	return &APIError{
		Code:       code,
		Message:    message,
		Details:    details,
		StatusCode: statusCode,
	}
}

// Предопределенные API ошибки
var (
	ErrAPIInvalidInput = &APIError{
		Code:       "INVALID_INPUT",
		Message:    "Invalid input data",
		StatusCode: 400,
	}

	ErrAPINotFound = &APIError{
		Code:       "NOT_FOUND",
		Message:    "Resource not found",
		StatusCode: 404,
	}

	ErrAPIUnauthorized = &APIError{
		Code:       "UNAUTHORIZED",
		Message:    "Unauthorized access",
		StatusCode: 401,
	}

	ErrAPIForbidden = &APIError{
		Code:       "FORBIDDEN",
		Message:    "Access forbidden",
		StatusCode: 403,
	}

	ErrAPIInternalServer = &APIError{
		Code:       "INTERNAL_SERVER_ERROR",
		Message:    "Internal server error",
		StatusCode: 500,
	}

	ErrAPIBadGateway = &APIError{
		Code:       "BAD_GATEWAY",
		Message:    "Bad gateway",
		StatusCode: 502,
	}

	ErrAPIServiceUnavailable = &APIError{
		Code:       "SERVICE_UNAVAILABLE",
		Message:    "Service temporarily unavailable",
		StatusCode: 503,
	}
)

// Helper functions для создания типовых ошибок

// NewValidationError создает ошибку валидации
func NewValidationError(details string) *APIError {
	return &APIError{
		Code:       "VALIDATION_ERROR",
		Message:    "Validation failed",
		Details:    details,
		StatusCode: 400,
	}
}

// NewNotFoundError создает ошибку "не найдено"
func NewNotFoundError(message string) *APIError {
	return &APIError{
		Code:       "NOT_FOUND",
		Message:    message,
		StatusCode: 404,
	}
}

// NewUnauthorizedError создает ошибку "не авторизован"
func NewUnauthorizedError(message string) *APIError {
	return &APIError{
		Code:       "UNAUTHORIZED",
		Message:    message,
		StatusCode: 401,
	}
}

// NewForbiddenError создает ошибку "доступ запрещен"
func NewForbiddenError(message string) *APIError {
	return &APIError{
		Code:       "FORBIDDEN",
		Message:    message,
		StatusCode: 403,
	}
}

// NewInternalServerError создает ошибку "внутренняя ошибка сервера"
func NewInternalServerError() *APIError {
	return &APIError{
		Code:       "INTERNAL_SERVER_ERROR",
		Message:    "Internal server error",
		StatusCode: 500,
	}
}

// NewBadGatewayError создает ошибку "bad gateway"
func NewBadGatewayError(details string) *APIError {
	return &APIError{
		Code:       "BAD_GATEWAY",
		Message:    "Bad gateway",
		Details:    details,
		StatusCode: 502,
	}
}

// NewServiceUnavailableError создает ошибку "сервис недоступен"
func NewServiceUnavailableError(details string) *APIError {
	return &APIError{
		Code:       "SERVICE_UNAVAILABLE",
		Message:    "Service temporarily unavailable",
		Details:    details,
		StatusCode: 503,
	}
}
