package handlers

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"

	pb "github.com/rusgainew/kkm-project-mks/proto-lib/api"
	"github.com/rusgainew/kkm-project-mks/user-query-server/internal/application/ports"
	"github.com/rusgainew/kkm-project-mks/user-query-server/internal/infrastructure/messaging"
)

// UserEventHandler implements messaging.EventHandler for user events
type UserEventHandler struct {
	repo   ports.UserRepository
	logger *zap.Logger
}

// NewUserEventHandler creates a new user event handler
func NewUserEventHandler(repo ports.UserRepository, logger *zap.Logger) *UserEventHandler {
	return &UserEventHandler{
		repo:   repo,
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
	phone, _ := event.Data["phone"].(string)
	role, _ := event.Data["role"].(string)
	if role == "" {
		role = "user"
	}

	// Create user read model
	user := &pb.UserReadModel{
		Id:        userID,
		Email:     email,
		FirstName: firstName,
		LastName:  lastName,
		Phone:     phone,
		Status:    "active",
		Role:      role,
		CreatedAt: time.Now().Format(time.RFC3339),
		UpdatedAt: time.Now().Format(time.RFC3339),
	}

	// Upsert into in-memory repository
	if err := h.repo.UpsertUser(ctx, user); err != nil {
		h.logger.Error("Failed to upsert user",
			zap.String("user_id", userID),
			zap.String("email", email),
			zap.Error(err))
		return fmt.Errorf("failed to upsert user: %w", err)
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

	// Get existing user from repository
	user, err := h.repo.GetUser(ctx, userID)
	if err != nil {
		h.logger.Error("Failed to get user",
			zap.String("user_id", userID),
			zap.Error(err))
		return fmt.Errorf("failed to get user: %w", err)
	}

	if user == nil {
		h.logger.Warn("User not found for login update",
			zap.String("user_id", userID))
		// Still return nil to not requeue the message
		return nil
	}

	// Update last login timestamp
	user.LastLoginAt = time.Now().Format(time.RFC3339)
	user.UpdatedAt = time.Now().Format(time.RFC3339)

	// Update in repository
	if err := h.repo.UpsertUser(ctx, user); err != nil {
		h.logger.Error("Failed to update last login",
			zap.String("user_id", userID),
			zap.Error(err))
		return fmt.Errorf("failed to update last login: %w", err)
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

	// Get existing user from repository
	user, err := h.repo.GetUser(ctx, userID)
	if err != nil {
		h.logger.Error("Failed to get user",
			zap.String("user_id", userID),
			zap.Error(err))
		return fmt.Errorf("failed to get user: %w", err)
	}

	if user == nil {
		h.logger.Warn("User not found for token refresh update",
			zap.String("user_id", userID))
		return nil
	}

	// Update timestamp
	user.UpdatedAt = time.Now().Format(time.RFC3339)

	// Update in repository
	if err := h.repo.UpsertUser(ctx, user); err != nil {
		h.logger.Error("Failed to update token refresh",
			zap.String("user_id", userID),
			zap.Error(err))
		return fmt.Errorf("failed to update token refresh: %w", err)
	}

	h.logger.Debug("Token refresh event processed",
		zap.String("user_id", userID))

	return nil
}
