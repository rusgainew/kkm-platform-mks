// Файл api-gateway/internal/domain/errors/grpc_mapper.go содержит реализацию пакета errors.
package errors

import (
	"net/http"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// MapGRPCErrorToHTTP преобразует gRPC ошибку в HTTP статус и APIError
func MapGRPCErrorToHTTP(err error) (int, *APIError) {
	if err == nil {
		return http.StatusOK, nil
	}

	st, ok := status.FromError(err)
	if !ok {
		return http.StatusInternalServerError, NewInternalServerError()
	}

	switch st.Code() {
	case codes.OK:
		return http.StatusOK, nil

	case codes.Canceled:
		return http.StatusRequestTimeout, &APIError{
			Code:       "REQUEST_CANCELED",
			Message:    "Request was canceled",
			Details:    st.Message(),
			StatusCode: http.StatusRequestTimeout,
		}

	case codes.Unknown:
		return http.StatusInternalServerError, &APIError{
			Code:       "UNKNOWN_ERROR",
			Message:    "An unknown error occurred",
			StatusCode: http.StatusInternalServerError,
		}

	case codes.InvalidArgument:
		return http.StatusBadRequest, NewValidationError(st.Message())

	case codes.DeadlineExceeded:
		return http.StatusGatewayTimeout, &APIError{
			Code:       "TIMEOUT",
			Message:    "Request timed out",
			Details:    st.Message(),
			StatusCode: http.StatusGatewayTimeout,
		}

	case codes.NotFound:
		return http.StatusNotFound, NewNotFoundError(st.Message())

	case codes.AlreadyExists:
		return http.StatusConflict, &APIError{
			Code:       "ALREADY_EXISTS",
			Message:    "Resource already exists",
			Details:    st.Message(),
			StatusCode: http.StatusConflict,
		}

	case codes.PermissionDenied:
		return http.StatusForbidden, NewForbiddenError(st.Message())

	case codes.ResourceExhausted:
		return http.StatusTooManyRequests, &APIError{
			Code:       "RATE_LIMITED",
			Message:    "Too many requests",
			Details:    st.Message(),
			StatusCode: http.StatusTooManyRequests,
		}

	case codes.FailedPrecondition:
		return http.StatusPreconditionFailed, &APIError{
			Code:       "PRECONDITION_FAILED",
			Message:    "Operation preconditions not met",
			Details:    st.Message(),
			StatusCode: http.StatusPreconditionFailed,
		}

	case codes.Aborted:
		return http.StatusConflict, &APIError{
			Code:       "ABORTED",
			Message:    "Operation was aborted",
			Details:    st.Message(),
			StatusCode: http.StatusConflict,
		}

	case codes.OutOfRange:
		return http.StatusBadRequest, &APIError{
			Code:       "OUT_OF_RANGE",
			Message:    "Value out of range",
			Details:    st.Message(),
			StatusCode: http.StatusBadRequest,
		}

	case codes.Unimplemented:
		return http.StatusNotImplemented, &APIError{
			Code:       "NOT_IMPLEMENTED",
			Message:    "Operation not implemented",
			StatusCode: http.StatusNotImplemented,
		}

	case codes.Internal:
		return http.StatusInternalServerError, NewInternalServerError()

	case codes.Unavailable:
		return http.StatusServiceUnavailable, NewServiceUnavailableError(st.Message())

	case codes.DataLoss:
		return http.StatusInternalServerError, &APIError{
			Code:       "DATA_LOSS",
			Message:    "Unrecoverable data loss or corruption",
			StatusCode: http.StatusInternalServerError,
		}

	case codes.Unauthenticated:
		return http.StatusUnauthorized, NewUnauthorizedError(st.Message())

	default:
		return http.StatusInternalServerError, NewInternalServerError()
	}
}

// IsNotFoundError проверяет является ли ошибка NotFound
func IsNotFoundError(err error) bool {
	if err == nil {
		return false
	}
	st, ok := status.FromError(err)
	if !ok {
		return false
	}
	return st.Code() == codes.NotFound
}

// IsInvalidArgumentError проверяет является ли ошибка InvalidArgument
func IsInvalidArgumentError(err error) bool {
	if err == nil {
		return false
	}
	st, ok := status.FromError(err)
	if !ok {
		return false
	}
	return st.Code() == codes.InvalidArgument
}

// IsAlreadyExistsError проверяет является ли ошибка AlreadyExists
func IsAlreadyExistsError(err error) bool {
	if err == nil {
		return false
	}
	st, ok := status.FromError(err)
	if !ok {
		return false
	}
	return st.Code() == codes.AlreadyExists
}

// IsUnavailableError проверяет является ли ошибка Unavailable
func IsUnavailableError(err error) bool {
	if err == nil {
		return false
	}
	st, ok := status.FromError(err)
	if !ok {
		return false
	}
	return st.Code() == codes.Unavailable
}
