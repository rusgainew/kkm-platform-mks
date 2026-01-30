// Файл foreign-company-server/internal/domain/errors.go содержит реализацию пакета domain.
package domain

import "errors"

var (
	// ErrForeignCompanyNotFound возвращается когда компания не найдена
	ErrForeignCompanyNotFound = errors.New("foreign company not found")

	// ErrForeignCompanyAlreadyExists возвращается когда компания с таким PIN уже существует
	ErrForeignCompanyAlreadyExists = errors.New("foreign company with this PIN already exists")

	// ErrInvalidPIN возвращается при некорректном формате PIN
	ErrInvalidPIN = errors.New("invalid PIN format: must be 5-50 characters")

	// ErrInvalidFullName возвращается при некорректном названии
	ErrInvalidFullName = errors.New("invalid full name: must be 2-500 characters")

	// ErrInvalidCountryCode возвращается при некорректном коде страны
	ErrInvalidCountryCode = errors.New("invalid country code: must be 2-letter ISO 3166-1 alpha-2 code")

	// ErrInvalidCreatedBy возвращается при отсутствии создателя
	ErrInvalidCreatedBy = errors.New("created_by user ID is required")

	// ErrCannotDeleteActiveCompany возвращается при попытке удалить активную компанию
	ErrCannotDeleteActiveCompany = errors.New("cannot delete active foreign company")

	// ErrDatabaseConnection возвращается при проблемах с БД
	ErrDatabaseConnection = errors.New("database connection error")

	// ErrInvalidID возвращается при некорректном ID
	ErrInvalidID = errors.New("invalid foreign company ID")
)
