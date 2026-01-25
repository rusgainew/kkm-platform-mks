package domain

import "errors"

var (
	// BankAccount errors
	ErrBankAccountNotFound      = errors.New("bank account not found")
	ErrBankAccountAlreadyExists = errors.New("bank account with this number already exists")
	ErrInvalidBankAccount       = errors.New("invalid bank account data")

	// Validation errors
	ErrInvalidAccountNumber = errors.New("invalid account number")
	ErrInvalidIBAN          = errors.New("invalid IBAN")
	ErrInvalidCurrency      = errors.New("invalid currency code")
	ErrInvalidBIC           = errors.New("invalid BIC code")
	ErrInvalidBankID        = errors.New("invalid bank ID")

	// Business logic errors
	ErrCannotDeleteDefaultAccount = errors.New("cannot delete default bank account")
	ErrNoActiveAccounts           = errors.New("organization has no active bank accounts")

	// Permission errors
	ErrUnauthorized            = errors.New("unauthorized access")
	ErrInsufficientPermissions = errors.New("insufficient permissions")

	// Organization errors
	ErrOrganizationNotFound = errors.New("organization not found")
)
