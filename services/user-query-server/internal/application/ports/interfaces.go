package ports

import (
	"context"

	pb "github.com/rusgainew/kkm-project-mks/proto-lib/api"
)

// UserRepository интерфейс для репозитория пользователей
type UserRepository interface {
	GetUser(ctx context.Context, userID string) (*pb.UserReadModel, error)
	ListUsers(ctx context.Context, offset, limit int32, status, role string) ([]*pb.UserReadModel, int64, error)
	SearchUsers(ctx context.Context, query string, offset, limit int32, status string) ([]*pb.UserReadModel, int64, error)
	UpsertUser(ctx context.Context, user *pb.UserReadModel) error
}

// UserCache интерфейс для кеша пользователей
type UserCache interface {
	GetUser(ctx context.Context, userID string) (*pb.UserReadModel, error)
	SetUser(ctx context.Context, user *pb.UserReadModel) error
	InvalidateUser(ctx context.Context, userID string) error
}
