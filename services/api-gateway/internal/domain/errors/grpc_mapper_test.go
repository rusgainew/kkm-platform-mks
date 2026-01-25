package errors_test

import (
	"net/http"
	"testing"

	"github.com/rusgainew/kkm-project-mks/api-gateway/internal/domain/errors"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestMapGRPCErrorToHTTP(t *testing.T) {
	tests := []struct {
		name           string
		grpcErr        error
		expectedStatus int
		expectedCode   string
	}{
		{
			name:           "nil_error",
			grpcErr:        nil,
			expectedStatus: http.StatusOK,
			expectedCode:   "",
		},
		{
			name:           "not_found",
			grpcErr:        status.Error(codes.NotFound, "resource not found"),
			expectedStatus: http.StatusNotFound,
			expectedCode:   "NOT_FOUND",
		},
		{
			name:           "invalid_argument",
			grpcErr:        status.Error(codes.InvalidArgument, "invalid field"),
			expectedStatus: http.StatusBadRequest,
			expectedCode:   "VALIDATION_ERROR",
		},
		{
			name:           "already_exists",
			grpcErr:        status.Error(codes.AlreadyExists, "resource exists"),
			expectedStatus: http.StatusConflict,
			expectedCode:   "ALREADY_EXISTS",
		},
		{
			name:           "permission_denied",
			grpcErr:        status.Error(codes.PermissionDenied, "access denied"),
			expectedStatus: http.StatusForbidden,
			expectedCode:   "FORBIDDEN",
		},
		{
			name:           "unauthenticated",
			grpcErr:        status.Error(codes.Unauthenticated, "not authenticated"),
			expectedStatus: http.StatusUnauthorized,
			expectedCode:   "UNAUTHORIZED",
		},
		{
			name:           "resource_exhausted",
			grpcErr:        status.Error(codes.ResourceExhausted, "rate limited"),
			expectedStatus: http.StatusTooManyRequests,
			expectedCode:   "RATE_LIMITED",
		},
		{
			name:           "unavailable",
			grpcErr:        status.Error(codes.Unavailable, "service down"),
			expectedStatus: http.StatusServiceUnavailable,
			expectedCode:   "SERVICE_UNAVAILABLE",
		},
		{
			name:           "deadline_exceeded",
			grpcErr:        status.Error(codes.DeadlineExceeded, "timeout"),
			expectedStatus: http.StatusGatewayTimeout,
			expectedCode:   "TIMEOUT",
		},
		{
			name:           "internal",
			grpcErr:        status.Error(codes.Internal, "internal error"),
			expectedStatus: http.StatusInternalServerError,
			expectedCode:   "INTERNAL_SERVER_ERROR",
		},
		{
			name:           "unimplemented",
			grpcErr:        status.Error(codes.Unimplemented, "not implemented"),
			expectedStatus: http.StatusNotImplemented,
			expectedCode:   "NOT_IMPLEMENTED",
		},
		{
			name:           "canceled",
			grpcErr:        status.Error(codes.Canceled, "request canceled"),
			expectedStatus: http.StatusRequestTimeout,
			expectedCode:   "REQUEST_CANCELED",
		},
		{
			name:           "aborted",
			grpcErr:        status.Error(codes.Aborted, "operation aborted"),
			expectedStatus: http.StatusConflict,
			expectedCode:   "ABORTED",
		},
		{
			name:           "out_of_range",
			grpcErr:        status.Error(codes.OutOfRange, "out of range"),
			expectedStatus: http.StatusBadRequest,
			expectedCode:   "OUT_OF_RANGE",
		},
		{
			name:           "failed_precondition",
			grpcErr:        status.Error(codes.FailedPrecondition, "precondition failed"),
			expectedStatus: http.StatusPreconditionFailed,
			expectedCode:   "PRECONDITION_FAILED",
		},
		{
			name:           "data_loss",
			grpcErr:        status.Error(codes.DataLoss, "data corrupted"),
			expectedStatus: http.StatusInternalServerError,
			expectedCode:   "DATA_LOSS",
		},
		{
			name:           "unknown",
			grpcErr:        status.Error(codes.Unknown, "unknown error"),
			expectedStatus: http.StatusInternalServerError,
			expectedCode:   "UNKNOWN_ERROR",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			statusCode, apiErr := errors.MapGRPCErrorToHTTP(tt.grpcErr)

			assert.Equal(t, tt.expectedStatus, statusCode)

			if tt.grpcErr != nil {
				assert.NotNil(t, apiErr)
				assert.Equal(t, tt.expectedCode, apiErr.Code)
			}
		})
	}
}

func TestIsNotFoundError(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{
			name:     "nil_error",
			err:      nil,
			expected: false,
		},
		{
			name:     "not_found_error",
			err:      status.Error(codes.NotFound, "not found"),
			expected: true,
		},
		{
			name:     "other_error",
			err:      status.Error(codes.Internal, "internal"),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := errors.IsNotFoundError(tt.err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestIsInvalidArgumentError(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{
			name:     "nil_error",
			err:      nil,
			expected: false,
		},
		{
			name:     "invalid_argument_error",
			err:      status.Error(codes.InvalidArgument, "invalid"),
			expected: true,
		},
		{
			name:     "other_error",
			err:      status.Error(codes.NotFound, "not found"),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := errors.IsInvalidArgumentError(tt.err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestIsAlreadyExistsError(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{
			name:     "nil_error",
			err:      nil,
			expected: false,
		},
		{
			name:     "already_exists_error",
			err:      status.Error(codes.AlreadyExists, "exists"),
			expected: true,
		},
		{
			name:     "other_error",
			err:      status.Error(codes.NotFound, "not found"),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := errors.IsAlreadyExistsError(tt.err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestIsUnavailableError(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{
			name:     "nil_error",
			err:      nil,
			expected: false,
		},
		{
			name:     "unavailable_error",
			err:      status.Error(codes.Unavailable, "unavailable"),
			expected: true,
		},
		{
			name:     "other_error",
			err:      status.Error(codes.NotFound, "not found"),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := errors.IsUnavailableError(tt.err)
			assert.Equal(t, tt.expected, result)
		})
	}
}
