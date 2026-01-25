package grpc

import (
	"fmt"

	"github.com/rusgainew/kkm-project-mks/document-server/internal/application/document"
	"github.com/rusgainew/kkm-project-mks/document-server/internal/domain"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// MapDomainErrorToGRPC преобразует доменные ошибки в gRPC статусы
func MapDomainErrorToGRPC(err error) error {
	if err == nil {
		return nil
	}

	switch err {
	case domain.ErrInvalidInput:
		return status.Error(codes.InvalidArgument, "invalid input parameters")
	case domain.ErrInvalidDocumentID:
		return status.Error(codes.InvalidArgument, "invalid document id")
	case domain.ErrInvalidOrganizationID:
		return status.Error(codes.InvalidArgument, "invalid organization id")
	case domain.ErrVersionConflict:
		return status.Error(codes.FailedPrecondition, "version conflict: document has been modified by another request")
	case domain.ErrConcurrentModification:
		return status.Error(codes.FailedPrecondition, "concurrent modification detected")
	default:
		// Если это не известная доменная ошибка, вернуть как Internal
		return status.Error(codes.Internal, "internal server error")
	}
}

// MapValidationErrorsToGRPC преобразует ошибки валидации в gRPC ответ
func MapValidationErrorsToGRPC(validationErrors []document.ValidationError) error {
	if len(validationErrors) == 0 {
		return nil
	}

	errorMsg := fmt.Sprintf("validation failed: %d error(s)\n", len(validationErrors))
	for i, ve := range validationErrors {
		errorMsg += fmt.Sprintf("%d. %s: %s\n", i+1, ve.Field, ve.Message)
	}

	return status.Error(codes.InvalidArgument, errorMsg)
}

// IsNotFoundError проверяет, является ли ошибка ошибкой "не найдено"
func IsNotFoundError(err error) bool {
	if err == nil {
		return false
	}
	return err.Error() == "document not found"
}

// MapRepositoryErrorToGRPC преобразует ошибки репозитория в gRPC статусы
func MapRepositoryErrorToGRPC(err error) error {
	if err == nil {
		return nil
	}

	errMsg := err.Error()

	// Проверяем на ошибки "не найдено"
	if errMsg == "document not found" {
		return status.Error(codes.NotFound, "document not found")
	}

	// По умолчанию - внутренняя ошибка сервера
	return status.Error(codes.Internal, "failed to process request")
}
