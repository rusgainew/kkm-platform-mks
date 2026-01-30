package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"go.uber.org/zap"

	pb "github.com/rusgainew/kkm-project-mks/proto-lib/api"
)

// InitDB initializes PostgreSQL database connection
func InitDB(databaseURL string) (*sqlx.DB, error) {
	db, err := sqlx.Connect("postgres", databaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Test the connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Set connection pool settings
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)

	return db, nil
}

type UserQueryRepository struct {
	db     *sqlx.DB
	logger *zap.Logger
}

func NewUserQueryRepository(db *sqlx.DB, logger *zap.Logger) *UserQueryRepository {
	return &UserQueryRepository{
		db:     db,
		logger: logger,
	}
}

// GetUser retrieves a single user by ID
func (r *UserQueryRepository) GetUser(ctx context.Context, userID string) (*pb.UserReadModel, error) {
	query := `
		SELECT id, email, first_name, last_name, phone, status, role, last_login_at, created_at, updated_at
		FROM user_read_model
		WHERE id = $1 AND deleted_at IS NULL
	`

	var user pb.UserReadModel
	var phone, lastLoginAt sql.NullString
	err := r.db.QueryRowContext(ctx, query, userID).Scan(
		&user.Id,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&phone,
		&user.Status,
		&user.Role,
		&lastLoginAt,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if phone.Valid {
		user.Phone = phone.String
	}
	if lastLoginAt.Valid {
		user.LastLoginAt = lastLoginAt.String
	}

	if err != nil {
		r.logger.Error("Failed to get user",
			zap.String("user_id", userID),
			zap.Error(err))
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &user, nil
}

// ListUsers retrieves users with pagination and optional filters
func (r *UserQueryRepository) ListUsers(ctx context.Context, offset, limit int32, status, role string) ([]*pb.UserReadModel, int64, error) {
	// Build dynamic query
	query := `SELECT id, email, first_name, last_name, phone, status, role, last_login_at, created_at, updated_at
	          FROM user_read_model WHERE deleted_at IS NULL`

	args := []interface{}{}
	argCount := 1

	if status != "" {
		query += fmt.Sprintf(" AND status = $%d", argCount)
		args = append(args, status)
		argCount++
	}

	if role != "" {
		query += fmt.Sprintf(" AND role = $%d", argCount)
		args = append(args, role)
		argCount++
	}

	// Count total
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM (%s) t", query)
	var totalCount int64
	if err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&totalCount); err != nil {
		r.logger.Error("Failed to count users", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to count users: %w", err)
	}

	// Add ordering and pagination
	query += fmt.Sprintf(" ORDER BY created_at DESC LIMIT $%d OFFSET $%d", argCount, argCount+1)
	args = append(args, limit, offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		r.logger.Error("Failed to list users", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list users: %w", err)
	}
	defer rows.Close()

	var users []*pb.UserReadModel
	for rows.Next() {
		var user pb.UserReadModel
		var phone, lastLoginAt sql.NullString
		if err := rows.Scan(
			&user.Id,
			&user.Email,
			&user.FirstName,
			&user.LastName,
			&phone,
			&user.Status,
			&user.Role,
			&lastLoginAt,
			&user.CreatedAt,
			&user.UpdatedAt,
		); err != nil {
			r.logger.Error("Failed to scan user", zap.Error(err))
			return nil, 0, fmt.Errorf("failed to scan user: %w", err)
		}
		if phone.Valid {
			user.Phone = phone.String
		}
		if lastLoginAt.Valid {
			user.LastLoginAt = lastLoginAt.String
		}
		users = append(users, &user)
	}

	if err = rows.Err(); err != nil {
		r.logger.Error("Row iteration error", zap.Error(err))
		return nil, 0, fmt.Errorf("row iteration error: %w", err)
	}

	return users, totalCount, nil
}

// SearchUsers performs full-text search on users
func (r *UserQueryRepository) SearchUsers(ctx context.Context, searchQuery string, offset, limit int32, status string) ([]*pb.UserReadModel, int64, error) {
	// Use full-text search on email and name
	query := `
		SELECT id, email, first_name, last_name, phone, status, role, last_login_at, created_at, updated_at
		FROM user_read_model
		WHERE deleted_at IS NULL
		  AND (
			email ILIKE $1
			OR first_name ILIKE $1
			OR last_name ILIKE $1
		  )
	`

	args := []interface{}{fmt.Sprintf("%%%s%%", searchQuery)}

	if status != "" {
		query += " AND status = $2"
		args = append(args, status)
	}

	// Count total
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM (%s) t", query)
	var totalCount int64
	if err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&totalCount); err != nil {
		r.logger.Error("Failed to count users", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to count users: %w", err)
	}

	// Add ordering and pagination
	paramCount := len(args) + 1
	query += fmt.Sprintf(" ORDER BY created_at DESC LIMIT $%d OFFSET $%d", paramCount, paramCount+1)
	args = append(args, limit, offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		r.logger.Error("Failed to search users", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to search users: %w", err)
	}
	defer rows.Close()

	var users []*pb.UserReadModel
	for rows.Next() {
		var user pb.UserReadModel
		var phone, lastLoginAt sql.NullString
		if err := rows.Scan(
			&user.Id,
			&user.Email,
			&user.FirstName,
			&user.LastName,
			&phone,
			&user.Status,
			&user.Role,
			&lastLoginAt,
			&user.CreatedAt,
			&user.UpdatedAt,
		); err != nil {
			r.logger.Error("Failed to scan user", zap.Error(err))
			return nil, 0, fmt.Errorf("failed to scan user: %w", err)
		}
		if phone.Valid {
			user.Phone = phone.String
		}
		if lastLoginAt.Valid {
			user.LastLoginAt = lastLoginAt.String
		}
		users = append(users, &user)
	}

	if err = rows.Err(); err != nil {
		r.logger.Error("Row iteration error", zap.Error(err))
		return nil, 0, fmt.Errorf("row iteration error: %w", err)
	}

	return users, totalCount, nil
}
