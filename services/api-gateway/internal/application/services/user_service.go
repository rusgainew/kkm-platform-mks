package services

import (
	"context"
	"fmt"
	"strconv"

	"github.com/rusgainew/kkm-project-mks/api-gateway/internal/domain/models"
	"github.com/rusgainew/kkm-project-mks/api-gateway/internal/infrastructure/client"
	"github.com/rusgainew/kkm-project-mks/api-gateway/internal/infrastructure/observability"
	authpb "github.com/rusgainew/kkm-project-mks/proto-lib/auth"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

// formatTimestamp converts int64 timestamp to string
func formatTimestamp(ts int64) string {
	return strconv.FormatInt(ts, 10)
}

// UserService сервис для работы с пользователями и аутентификацией
type UserService struct {
	connManager *client.ConnectionManager
	serviceAddr string
	metrics     *observability.Metrics
	tracer      *observability.Tracer
	logger      *zap.Logger
}

// NewUserService создает новый UserService
func NewUserService(
	connManager *client.ConnectionManager,
	serviceAddr string,
	metrics *observability.Metrics,
	tracer *observability.Tracer,
	logger *zap.Logger,
) *UserService {
	return &UserService{
		connManager: connManager,
		serviceAddr: serviceAddr,
		metrics:     metrics,
		tracer:      tracer,
		logger:      logger,
	}
}

// getClient возвращает gRPC клиент для user service
func (s *UserService) getClient(ctx context.Context) (authpb.AuthServiceClient, error) {
	conn, err := s.connManager.GetConnection(ctx, s.serviceAddr)
	if err != nil {
		return nil, fmt.Errorf("failed to get connection: %w", err)
	}
	return authpb.NewAuthServiceClient(conn), nil
}

// Register регистрирует нового пользователя
func (s *UserService) Register(ctx context.Context, email, password, firstName, lastName string) (*models.AuthResponse, error) {
	ctx, endSpan := s.tracer.StartSpan(ctx, "UserService.Register")
	defer endSpan()

	client, err := s.getClient(ctx)
	if err != nil {
		s.logger.Error("Failed to get user client", zap.Error(err))
		return nil, err
	}

	req := &authpb.RegisterRequest{
		Email:     email,
		Password:  password,
		FirstName: firstName,
		LastName:  lastName,
	}

	resp, err := client.Register(ctx, req)
	if err != nil {
		s.logger.Error("Failed to register user", zap.Error(err))
		s.metrics.IncrementErrorCount("user_service", "register")
		return nil, fmt.Errorf("failed to register user: %w", err)
	}

	s.logger.Info("User registered successfully", zap.String("user_id", resp.User.Id))

	return &models.AuthResponse{
		AccessToken:  resp.Token.AccessToken,
		RefreshToken: resp.Token.RefreshToken,
		ExpiresIn:    resp.Token.ExpiresIn,
		RefreshExp:   resp.Token.RefreshExpiresIn,
		TokenType:    resp.Token.TokenType,
		Timestamp:    resp.Timestamp,
		User: models.User{
			UserID:    resp.User.Id,
			Email:     resp.User.Email,
			FirstName: resp.User.FirstName,
			LastName:  resp.User.LastName,
			Role:      resp.User.Role,
			CreatedAt: formatTimestamp(resp.User.CreatedAt),
			UpdatedAt: formatTimestamp(resp.User.UpdatedAt),
			Status:    resp.User.Status,
		},
	}, nil
}

// Login выполняет вход пользователя
func (s *UserService) Login(ctx context.Context, email, password string) (*models.AuthResponse, error) {
	ctx, endSpan := s.tracer.StartSpan(ctx, "UserService.Login")
	defer endSpan()

	client, err := s.getClient(ctx)
	if err != nil {
		s.logger.Error("Failed to get user client", zap.Error(err))
		return nil, err
	}

	req := &authpb.LoginRequest{
		Email:    email,
		Password: password,
	}

	resp, err := client.Login(ctx, req)
	if err != nil {
		s.logger.Error("Failed to login user", zap.Error(err))
		s.metrics.IncrementErrorCount("user_service", "login")
		return nil, fmt.Errorf("failed to login: %w", err)
	}

	s.logger.Info("User logged in successfully", zap.String("user_id", resp.User.Id))

	return &models.AuthResponse{
		AccessToken:  resp.Token.AccessToken,
		RefreshToken: resp.Token.RefreshToken,
		ExpiresIn:    resp.Token.ExpiresIn,
		RefreshExp:   resp.Token.RefreshExpiresIn,
		TokenType:    resp.Token.TokenType,
		Timestamp:    resp.Timestamp,
		User: models.User{
			UserID:    resp.User.Id,
			Email:     resp.User.Email,
			FirstName: resp.User.FirstName,
			LastName:  resp.User.LastName,
			Role:      resp.User.Role,
			CreatedAt: formatTimestamp(resp.User.CreatedAt),
			UpdatedAt: formatTimestamp(resp.User.UpdatedAt),
			Status:    resp.User.Status,
		},
	}, nil
}

// ValidateToken проверяет токен и возвращает информацию о пользователе
func (s *UserService) ValidateToken(ctx context.Context, token string) (*models.User, error) {
	ctx, endSpan := s.tracer.StartSpan(ctx, "UserService.ValidateToken")
	defer endSpan()

	client, err := s.getClient(ctx)
	if err != nil {
		s.logger.Error("Failed to get user client", zap.Error(err))
		return nil, err
	}

	req := &authpb.ValidateTokenRequest{
		AccessToken: token,
	}

	resp, err := client.ValidateToken(ctx, req)
	if err != nil {
		s.logger.Error("Failed to validate token", zap.Error(err))
		s.metrics.IncrementErrorCount("user_service", "validate_token")
		return nil, fmt.Errorf("failed to validate token: %w", err)
	}

	return &models.User{
		UserID:    resp.Id,
		Email:     resp.Email,
		FirstName: resp.FirstName,
		LastName:  resp.LastName,
		Role:      resp.Role,
		CreatedAt: formatTimestamp(resp.CreatedAt),
		UpdatedAt: formatTimestamp(resp.UpdatedAt),
		Status:    resp.Status,
	}, nil
}

// RefreshToken обновляет токены по refresh token
func (s *UserService) RefreshToken(ctx context.Context, refreshToken string) (*models.TokenResponse, error) {
	ctx, endSpan := s.tracer.StartSpan(ctx, "UserService.RefreshToken")
	defer endSpan()

	client, err := s.getClient(ctx)
	if err != nil {
		s.logger.Error("Failed to get user client", zap.Error(err))
		return nil, err
	}

	req := &authpb.RefreshTokenRequest{RefreshToken: refreshToken}
	resp, err := client.RefreshToken(ctx, req)
	if err != nil {
		s.logger.Error("Failed to refresh token", zap.Error(err))
		s.metrics.IncrementErrorCount("user_service", "refresh_token")
		return nil, fmt.Errorf("failed to refresh token: %w", err)
	}

	return &models.TokenResponse{
		AccessToken:      resp.AccessToken,
		RefreshToken:     resp.RefreshToken,
		ExpiresIn:        resp.ExpiresIn,
		RefreshExpiresIn: resp.RefreshExpiresIn,
		TokenType:        resp.TokenType,
	}, nil
}

// Logout инвалидирует токен
func (s *UserService) Logout(ctx context.Context, accessToken string, revokeAll bool) error {
	ctx, endSpan := s.tracer.StartSpan(ctx, "UserService.Logout")
	defer endSpan()

	client, err := s.getClient(ctx)
	if err != nil {
		s.logger.Error("Failed to get user client", zap.Error(err))
		return err
	}

	req := &authpb.LogoutRequest{AccessToken: accessToken, RevokeAllSessions: revokeAll}
	_, err = client.Logout(ctx, req)
	if err != nil {
		s.logger.Error("Failed to logout", zap.Error(err))
		s.metrics.IncrementErrorCount("user_service", "logout")
		return fmt.Errorf("failed to logout: %w", err)
	}

	return nil
}

// UpdateProfile обновляет имя/фамилию пользователя
func (s *UserService) UpdateProfile(ctx context.Context, accessToken, userID, firstName, lastName string) (*models.User, error) {
	ctx, endSpan := s.tracer.StartSpan(ctx, "UserService.UpdateProfile")
	defer endSpan()

	client, err := s.getClient(ctx)
	if err != nil {
		s.logger.Error("Failed to get user client", zap.Error(err))
		return nil, err
	}

	req := &authpb.UpdateProfileRequest{UserId: userID}
	if firstName != "" {
		req.FirstName = wrapperspb.String(firstName)
	}
	if lastName != "" {
		req.LastName = wrapperspb.String(lastName)
	}
	resp, err := client.UpdateProfile(ctx, req)
	if err != nil {
		s.logger.Error("Failed to update profile", zap.Error(err))
		s.metrics.IncrementErrorCount("user_service", "update_profile")
		return nil, fmt.Errorf("failed to update profile: %w", err)
	}

	return &models.User{
		UserID:    resp.Id,
		Email:     resp.Email,
		FirstName: resp.FirstName,
		LastName:  resp.LastName,
		Role:      resp.Role,
		CreatedAt: formatTimestamp(resp.CreatedAt),
		UpdatedAt: formatTimestamp(resp.UpdatedAt),
		Status:    resp.Status,
	}, nil
}

// ChangePassword меняет пароль пользователя
func (s *UserService) ChangePassword(ctx context.Context, accessToken, userID, oldPassword, newPassword string) error {
	ctx, endSpan := s.tracer.StartSpan(ctx, "UserService.ChangePassword")
	defer endSpan()

	client, err := s.getClient(ctx)
	if err != nil {
		s.logger.Error("Failed to get user client", zap.Error(err))
		return err
	}

	req := &authpb.ChangePasswordRequest{UserId: userID, OldPassword: oldPassword, NewPassword: newPassword}
	_, err = client.ChangePassword(ctx, req)
	if err != nil {
		s.logger.Error("Failed to change password", zap.Error(err))
		s.metrics.IncrementErrorCount("user_service", "change_password")
		return fmt.Errorf("failed to change password: %w", err)
	}

	return nil
}

// ListUsers возвращает пользователей с пагинацией (admin only)
func (s *UserService) ListUsers(ctx context.Context, page, size int32, status, role string) (*models.ListUsersResponse, error) {
	ctx, endSpan := s.tracer.StartSpan(ctx, "UserService.ListUsers")
	defer endSpan()

	client, err := s.getClient(ctx)
	if err != nil {
		s.logger.Error("Failed to get user client", zap.Error(err))
		return nil, err
	}

	req := &authpb.ListUsersRequest{Page: page, Size: size}
	if status != "" {
		req.Status = wrapperspb.String(status)
	}
	if role != "" {
		req.Role = wrapperspb.String(role)
	}

	resp, err := client.ListUsers(ctx, req)
	if err != nil {
		s.logger.Error("Failed to list users", zap.Error(err))
		s.metrics.IncrementErrorCount("user_service", "list_users")
		return nil, fmt.Errorf("failed to list users: %w", err)
	}

	users := make([]models.User, 0, len(resp.Users))
	for _, u := range resp.Users {
		users = append(users, models.User{
			UserID:    u.Id,
			Email:     u.Email,
			FirstName: u.FirstName,
			LastName:  u.LastName,
			Role:      u.Role,
			CreatedAt: formatTimestamp(u.CreatedAt),
			UpdatedAt: formatTimestamp(u.UpdatedAt),
			Status:    u.Status,
			IsActive:  u.Status == "active",
		})
	}

	return &models.ListUsersResponse{
		Users: users,
		Page: models.PageInfo{
			Page:       resp.PageInfo.Page,
			Size:       resp.PageInfo.PerPage,
			TotalCount: int32(resp.PageInfo.Total),
		},
	}, nil
}

// ForgotPassword инициирует восстановление пароля
func (s *UserService) ForgotPassword(ctx context.Context, email string) error {
	ctx, endSpan := s.tracer.StartSpan(ctx, "UserService.ForgotPassword")
	defer endSpan()

	client, err := s.getClient(ctx)
	if err != nil {
		s.logger.Error("Failed to get user client", zap.Error(err))
		return err
	}

	req := &authpb.ForgotPasswordRequest{Email: email}
	_, err = client.ForgotPassword(ctx, req)
	if err != nil {
		s.logger.Error("Failed to send forgot password", zap.Error(err))
		s.metrics.IncrementErrorCount("user_service", "forgot_password")
		return fmt.Errorf("failed to send forgot password: %w", err)
	}

	return nil
}

// ResetPassword сбрасывает пароль по reset token
func (s *UserService) ResetPassword(ctx context.Context, resetToken, newPassword string) (*models.AuthResponse, error) {
	ctx, endSpan := s.tracer.StartSpan(ctx, "UserService.ResetPassword")
	defer endSpan()

	client, err := s.getClient(ctx)
	if err != nil {
		s.logger.Error("Failed to get user client", zap.Error(err))
		return nil, err
	}

	req := &authpb.ResetPasswordRequest{ResetToken: resetToken, NewPassword: newPassword}
	resp, err := client.ResetPassword(ctx, req)
	if err != nil {
		s.logger.Error("Failed to reset password", zap.Error(err))
		s.metrics.IncrementErrorCount("user_service", "reset_password")
		return nil, fmt.Errorf("failed to reset password: %w", err)
	}

	return &models.AuthResponse{
		AccessToken:  resp.Token.AccessToken,
		RefreshToken: resp.Token.RefreshToken,
		ExpiresIn:    resp.Token.ExpiresIn,
		RefreshExp:   resp.Token.RefreshExpiresIn,
		TokenType:    resp.Token.TokenType,
		Timestamp:    resp.Timestamp,
		User: models.User{
			UserID:    resp.User.Id,
			Email:     resp.User.Email,
			FirstName: resp.User.FirstName,
			LastName:  resp.User.LastName,
			Role:      resp.User.Role,
			CreatedAt: formatTimestamp(resp.User.CreatedAt),
			UpdatedAt: formatTimestamp(resp.User.UpdatedAt),
			Status:    resp.User.Status,
		},
	}, nil
}

// GetUser получает информацию о пользователе по ID
func (s *UserService) GetUser(ctx context.Context, userID string) (*models.User, error) {
	ctx, endSpan := s.tracer.StartSpan(ctx, "UserService.GetUser")
	defer endSpan()

	client, err := s.getClient(ctx)
	if err != nil {
		s.logger.Error("Failed to get user client", zap.Error(err))
		return nil, err
	}

	req := &authpb.GetUserRequest{
		UserId: userID,
	}

	resp, err := client.GetUser(ctx, req)
	if err != nil {
		s.logger.Error("Failed to get user", zap.Error(err))
		s.metrics.IncrementErrorCount("user_service", "get_user")
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &models.User{
		UserID:    resp.Id,
		Email:     resp.Email,
		FirstName: resp.FirstName,
		LastName:  resp.LastName,
		Role:      resp.Role,
		CreatedAt: formatTimestamp(resp.CreatedAt),
		UpdatedAt: formatTimestamp(resp.UpdatedAt),
		Status:    resp.Status,
	}, nil
}

// AssignRole вызывает user-server для назначения роли (admin only)
func (s *UserService) AssignRole(ctx context.Context, accessToken, targetUserID, role string) (*models.User, error) {
	ctx, endSpan := s.tracer.StartSpan(ctx, "UserService.AssignRole")
	defer endSpan()

	client, err := s.getClient(ctx)
	if err != nil {
		s.logger.Error("Failed to get user client", zap.Error(err))
		return nil, err
	}

	req := &authpb.AssignRoleRequest{
		AccessToken:  accessToken,
		TargetUserId: targetUserID,
		Role:         role,
	}

	resp, err := client.AssignRole(ctx, req)
	if err != nil {
		s.logger.Error("Failed to assign role", zap.Error(err))
		s.metrics.IncrementErrorCount("user_service", "assign_role")
		return nil, fmt.Errorf("failed to assign role: %w", err)
	}

	return &models.User{
		UserID:    resp.Id,
		Email:     resp.Email,
		FirstName: resp.FirstName,
		LastName:  resp.LastName,
		Role:      resp.Role,
		CreatedAt: formatTimestamp(resp.CreatedAt),
		UpdatedAt: formatTimestamp(resp.UpdatedAt),
		Status:    resp.Status,
	}, nil
}
