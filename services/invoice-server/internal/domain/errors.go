// Файл invoice-server/internal/domain/errors.go содержит реализацию пакета domain.
package domain

import "errors"

var (
	// Invoice errors
	ErrInvoiceNotFound       = errors.New("invoice not found")
	ErrInvoiceExists         = errors.New("invoice already exists")
	ErrInvalidInvoiceID      = errors.New("invalid invoice ID")
	ErrInvalidInvoiceNumber  = errors.New("invalid invoice number")
	ErrInvalidInvoiceStatus  = errors.New("invalid invoice status")
	ErrInvoiceNumberRequired = errors.New("invoice number is required")

	// Invoice state transition errors
	ErrCannotSign      = errors.New("cannot sign invoice in current state")
	ErrCannotAccept    = errors.New("cannot accept invoice in current state")
	ErrCannotReject    = errors.New("cannot reject invoice in current state")
	ErrCannotRevoke    = errors.New("cannot revoke invoice in current state")
	ErrAlreadySigned   = errors.New("invoice already signed")
	ErrAlreadyAccepted = errors.New("invoice already accepted")
	ErrAlreadyRejected = errors.New("invoice already rejected")
	ErrAlreadyRevoked  = errors.New("invoice already revoked")

	// Party errors
	ErrInvalidParty        = errors.New("invalid party data")
	ErrLegalPersonRequired = errors.New("legal person is required")
	ErrContractorRequired  = errors.New("contractor is required")

	// Invoice detail errors
	ErrDetailNotFound    = errors.New("invoice detail not found")
	ErrInvalidDetailData = errors.New("invalid detail data")
	ErrInvalidQuantity   = errors.New("invalid quantity")
	ErrInvalidPrice      = errors.New("invalid price")
	ErrInvalidAmount     = errors.New("invalid amount")
	ErrEmptyDetails      = errors.New("invoice must have at least one detail")

	// Financial errors
	ErrInvalidTotalAmount = errors.New("invalid total amount")
	ErrAmountMismatch     = errors.New("total amount does not match details sum")

	// Validation errors
	ErrInvalidInput         = errors.New("invalid input")
	ErrMissingRequiredField = errors.New("missing required field")
	ErrInvalidDateFormat    = errors.New("invalid date format")

	// Permission errors
	ErrUnauthorized     = errors.New("unauthorized access")
	ErrPermissionDenied = errors.New("permission denied")
)
