package tests

import (
	"context"
	"testing"

	"github.com/rusgainew/kkm-project-mks/document-server/internal/infrastructure/middleware"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

// TestTimeoutInterceptor проверяет timeout middleware
func TestTimeoutInterceptor(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	defer logger.Sync()

	ti := middleware.NewTimeoutInterceptor(logger)

	t.Run("Request completes successfully", func(t *testing.T) {
		handler := ti.UnaryServerInterceptor()
		resp, err := handler(
			context.Background(),
			nil,
			&grpc.UnaryServerInfo{FullMethod: "/api.DocumentService/GetDocument"},
			func(ctx context.Context, req interface{}) (interface{}, error) {
				return "success", nil
			},
		)
		assert.NoError(t, err)
		assert.Equal(t, "success", resp)
	})

	t.Run("Sets timeout on context", func(t *testing.T) {
		handler := ti.UnaryServerInterceptor()
		handler(
			context.Background(),
			nil,
			&grpc.UnaryServerInfo{FullMethod: "/api.DocumentService/CreateDocument"},
			func(ctx context.Context, req interface{}) (interface{}, error) {
				_, hasDeadline := ctx.Deadline()
				assert.True(t, hasDeadline, "Context should have deadline")
				return "ok", nil
			},
		)
	})
}
