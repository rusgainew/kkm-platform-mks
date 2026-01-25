package handlers

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"

	"github.com/rusgainew/kkm-project-mks/user-query-server/internal/infrastructure/messaging"
)

// UserEventHandler implements messaging.EventHandler for user events
type UserEventHandler struct {
	db     *sqlx.DB
	logger *zap.Logger
}

// NewUserEventHandler creates a new user event handler
func NewUserEventHandler(db *sqlx.DB, logger *zap.Logger) *UserEventHandler {
	return &UserEventHandler{
		db:     db,
		logger: logger,
	}
}

// HandleUserRegistered handles user.registered events
func (h *UserEventHandler) HandleUserRegistered(ctx context.Context, event messaging.UserEvent) error {
	// Extract data from event
	userID, ok := event.Data["user_id"].(string)
	if !ok || userID == "" {
		return fmt.Errorf("invalid user_id in event")
	}

	email, ok := event.Data["email"].(string)
	if !ok || email == "" {
		return fmt.Errorf("invalid email in event")
	}

	firstName, _ := event.Data["first_name"].(string)
	lastName, _ := event.Data["last_name"].(string)

	// Insert into read-model
	query := `
		INSERT INTO user_read_model (id, email, first_name, last_name, status, role, created_at)
		VALUES ($1, $2, $3, $4, 'active', 'user', NOW())
		ON CONFLICT (id) DO UPDATE SET
			email = EXCLUDED.email,
			first_name = EXCLUDED.first_name,
			last_name = EXCLUDED.last_name
	`

	if _, err := h.db.ExecContext(ctx, query, userID, email, firstName, lastName); err != nil {
		h.logger.Error("Failed to insert user",
			zap.String("user_id", userID),
			zap.String("email", email),
			zap.Error(err))
		return fmt.Errorf("failed to insert user: %w", err)
	}

	h.logger.Info("User registered event processed",
		zap.String("user_id", userID),
		zap.String("email", email))

	return nil
}

// HandleUserLoggedIn handles user.logged_in events
func (h *UserEventHandler) HandleUserLoggedIn(ctx context.Context, event messaging.UserEvent) error {
	userID, ok := event.Data["user_id"].(string)
	if !ok || userID == "" {
		return fmt.Errorf("invalid user_id in event")
	}

	// Update last login timestamp
	query := `
		UPDATE user_read_model 
		SET last_login_at = NOW()
		WHERE id = $1
	`

	result, err := h.db.ExecContext(ctx, query, userID)
	if err != nil {
		h.logger.Error("Failed to update last login",
			zap.String("user_id", userID),
			zap.Error(err))
		return fmt.Errorf("failed to update last login: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		h.logger.Warn("User not found for login update",
			zap.String("user_id", userID))
		// Still return nil to not requeue the message
		return nil
	}

	h.logger.Info("User login event processed",
		zap.String("user_id", userID))

	return nil
}

// HandleTokenRefreshed handles token.refreshed events
func (h *UserEventHandler) HandleTokenRefreshed(ctx context.Context, event messaging.UserEvent) error {
	userID, ok := event.Data["user_id"].(string)
	if !ok || userID == "" {
		return fmt.Errorf("invalid user_id in event")
	}

	// This is just for tracking, update status if needed
	query := `
		UPDATE user_read_model 
		SET updated_at = NOW()
		WHERE id = $1
	`

	if _, err := h.db.ExecContext(ctx, query, userID); err != nil {
		h.logger.Error("Failed to update token refresh",
			zap.String("user_id", userID),
			zap.Error(err))
		return fmt.Errorf("failed to update token refresh: %w", err)
	}

	h.logger.Debug("Token refresh event processed",
		zap.String("user_id", userID))

	return nil
}
