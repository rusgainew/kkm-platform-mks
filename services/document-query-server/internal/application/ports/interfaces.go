// Файл document-query-server/internal/application/ports/interfaces.go содержит реализацию пакета ports.
package ports

import (
	"context"

	pb "github.com/rusgainew/kkm-project-mks/proto-lib/api"
)

// DocumentRepository интерфейс для репозитория документов
type DocumentRepository interface {
	GetDocument(ctx context.Context, documentID string) (*pb.DocumentReadModel, error)
	ListDocuments(ctx context.Context, offset, limit int32, status, docType, companyID, approvalStatus string) ([]*pb.DocumentReadModel, int64, error)
	SearchDocuments(ctx context.Context, query string, offset, limit int32, docType, companyID, approvalStatus string) ([]*pb.DocumentReadModel, int64, error)
	GetPendingApprovalDocuments(ctx context.Context, companyID string, offset, limit int32) ([]*pb.DocumentReadModel, int64, error)
}

// DocumentCache интерфейс для кеша документов
type DocumentCache interface {
	GetDocument(ctx context.Context, documentID string) (*pb.DocumentReadModel, error)
	SetDocument(ctx context.Context, doc *pb.DocumentReadModel) error
	InvalidateDocument(ctx context.Context, documentID string) error
}
