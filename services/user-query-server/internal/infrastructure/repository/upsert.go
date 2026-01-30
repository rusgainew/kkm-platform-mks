// Файл user-query-server/internal/infrastructure/repository/upsert.go содержит реализацию пакета repository.
package repository

import (
	"context"
	"fmt"

	"go.uber.org/zap"

	pb "github.com/rusgainew/kkm-project-mks/proto-lib/api"
)

// UpsertUser inserts or updates a user in the read model
func (r *UserQueryRepository) UpsertUser(ctx context.Context, user *pb.UserReadModel) error {
	query := `
		INSERT INTO user_read_model (
			id, email, first_name, last_name, phone, status, role,
			last_login_at, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		ON CONFLICT (id) DO UPDATE SET
			email = EXCLUDED.email,
			first_name = EXCLUDED.first_name,
			last_name = EXCLUDED.last_name,
			phone = EXCLUDED.phone,
			status = EXCLUDED.status,
			role = EXCLUDED.role,
			last_login_at = EXCLUDED.last_login_at,
			updated_at = EXCLUDED.updated_at
	`

	_, err := r.db.ExecContext(ctx, query,
		user.Id,
		user.Email,
		user.FirstName,
		user.LastName,
		user.Phone,
		user.Status,
		user.Role,
		user.LastLoginAt,
		user.CreatedAt,
		user.UpdatedAt,
	)

	if err != nil {
		r.logger.Error("Failed to upsert user",
			zap.String("user_id", user.Id),
			zap.Error(err))
		return fmt.Errorf("failed to upsert user: %w", err)
	}

	return nil
}
