package http

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rusgainew/kkm-project-mks/api-gateway/internal/application/services"
	"github.com/rusgainew/kkm-project-mks/api-gateway/internal/domain/errors"
	"github.com/rusgainew/kkm-project-mks/api-gateway/internal/domain/models"
	"github.com/rusgainew/kkm-project-mks/lib/conversion"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// UserHandler обработчик HTTP запросов для пользователей
type UserHandler struct {
	service *services.UserService
	logger  *zap.Logger
}

// NewUserHandler создает новый UserHandler
func NewUserHandler(service *services.UserService, logger *zap.Logger) *UserHandler {
	return &UserHandler{
		service: service,
		logger:  logger,
	}
}

// RegisterRequest запрос на регистрацию
type RegisterRequest struct {
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=8"`
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
}

// LoginRequest запрос на вход
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// AssignRoleRequest запрос на назначение роли пользователю
type AssignRoleRequest struct {
	Role string `json:"role" binding:"required,oneof=user manager"`
}

// RefreshTokenRequest запрос на обновление токенов
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// LogoutRequest запрос на выход
type LogoutRequest struct {
	RevokeAll bool `json:"revoke_all_sessions"`
}

// UpdateProfileRequest запрос на обновление профиля
type UpdateProfileRequest struct {
	UserID    string `json:"user_id" binding:"required,uuid4"`
	FirstName string `json:"first_name" binding:"omitempty,min=1,max=100"`
	LastName  string `json:"last_name" binding:"omitempty,min=1,max=100"`
}

// ChangePasswordRequest запрос на смену пароля
type ChangePasswordRequest struct {
	UserID      string `json:"user_id" binding:"required,uuid4"`
	OldPassword string `json:"old_password" binding:"required,min=8"`
	NewPassword string `json:"new_password" binding:"required,min=8"`
}

// ForgotPasswordRequest запрос на восстановление пароля
type ForgotPasswordRequest struct {
	Email string `json:"email" binding:"required,email"`
}

// ResetPasswordRequest запрос на сброс пароля
type ResetPasswordRequest struct {
	ResetToken  string `json:"reset_token" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=8"`
}

// Register регистрирует нового пользователя
//
//	@Summary		User registration
//	@Description	Register a new user account. Returns access and refresh tokens.
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Param			request	body		RegisterRequest		true	"Registration data with email, password (min 8 chars), first_name and last_name"
//	@Success		201		{object}	models.AuthResponse	"User successfully registered"
//	@Failure		400		{object}	models.APIResponse	"Invalid input - validation errors"
//	@Failure		500		{object}	models.APIResponse	"Internal server error"
//	@Router			/users/register [post]
func (h *UserHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("Invalid register request", zap.Error(err))
		c.JSON(http.StatusBadRequest, errors.NewValidationError(err.Error()))
		return
	}

	authResp, err := h.service.Register(c.Request.Context(), req.Email, req.Password, req.FirstName, req.LastName)
	if err != nil {
		h.logger.Error("Failed to register user", zap.Error(err))
		statusCode, apiErr := errors.MapGRPCErrorToHTTP(err)
		c.JSON(statusCode, apiErr)
		return
	}

	c.JSON(http.StatusCreated, authResp)
}

// Login выполняет вход пользователя
//
//	@Summary		User login
//	@Description	User authentication with email and password. Returns JWT access and refresh tokens.
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Param			request	body		LoginRequest		true	"Login credentials (email and password)"
//	@Success		200		{object}	models.AuthResponse	"Successfully authenticated"
//	@Failure		400		{object}	models.APIResponse	"Invalid input - validation errors"
//	@Failure		401		{object}	models.APIResponse	"Unauthorized - invalid credentials"
//	@Failure		500		{object}	models.APIResponse	"Internal server error"
//	@Router			/users/login [post]
func (h *UserHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("Invalid login request", zap.Error(err))
		c.JSON(http.StatusBadRequest, errors.NewValidationError(err.Error()))
		return
	}

	authResp, err := h.service.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		h.logger.Error("Failed to login", zap.Error(err))
		statusCode, apiErr := errors.MapGRPCErrorToHTTP(err)
		c.JSON(statusCode, apiErr)
		return
	}

	c.JSON(http.StatusOK, authResp)
}

// RefreshToken обновляет access/refresh токены
//
//	@Summary	Refresh tokens
//	@Description	Refresh access and refresh tokens using refresh_token
//	@Tags		Users
//	@Accept		json
//	@Produce	json
//	@Param		request	body	RefreshTokenRequest	true	"Refresh token payload"
//	@Success	200	{object}	models.TokenResponse	"New tokens"
//	@Failure	400	{object}	models.APIResponse	"Invalid input"
//	@Failure	401	{object}	models.APIResponse	"Unauthorized"
//	@Failure	500	{object}	models.APIResponse	"Internal server error"
//	@Router		/users/refresh [post]
func (h *UserHandler) RefreshToken(c *gin.Context) {
	var req RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("Invalid refresh token request", zap.Error(err))
		c.JSON(http.StatusBadRequest, errors.NewValidationError(err.Error()))
		return
	}

	resp, err := h.service.RefreshToken(c.Request.Context(), req.RefreshToken)
	if err != nil {
		st, ok := status.FromError(err)
		if ok {
			switch st.Code() {
			case codes.InvalidArgument:
				c.JSON(http.StatusBadRequest, errors.NewValidationError(st.Message()))
				return
			case codes.Unauthenticated:
				c.JSON(http.StatusUnauthorized, errors.NewUnauthorizedError("Невалидный refresh token"))
				return
			case codes.PermissionDenied:
				c.JSON(http.StatusForbidden, errors.NewForbiddenError("Refresh token заблокирован"))
				return
			}
		}

		h.logger.Error("Failed to refresh token", zap.Error(err))
		c.JSON(http.StatusInternalServerError, errors.NewInternalServerError())
		return
	}

	c.JSON(http.StatusOK, resp)
}

// GetMe получает информацию о текущем пользователе
//
//	@Summary		Get current user
//	@Description	Get information about the currently authenticated user from JWT token
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	models.APIResponse{data=models.User}	"Current user information"
//	@Failure		401	{object}	models.APIResponse						"Unauthorized - invalid or missing token"
//	@Failure		500	{object}	models.APIResponse						"Internal server error"
//	@Router			/users/me [get]
//	@Security		BearerAuth
func (h *UserHandler) GetMe(c *gin.Context) {
	// Получаем user_id из контекста (установлен auth middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		h.logger.Warn("user_id not found in context")
		c.JSON(http.StatusUnauthorized, errors.NewUnauthorizedError("Требуется аутентификация"))
		return
	}

	user, err := h.service.GetUser(c.Request.Context(), userID.(string))
	if err != nil {
		h.logger.Error("Failed to get user", zap.Error(err))
		c.JSON(http.StatusInternalServerError, errors.NewInternalServerError())
		return
	}

	c.JSON(http.StatusOK, user)
}

// GetUser получает информацию о пользователе по ID
//
//	@Summary		Get user by ID
//	@Description	Get detailed information about a specific user by their ID
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		string				true	"User ID (UUID)"
//	@Success		200	{object}	models.User			"User information"
//	@Failure		400	{object}	models.APIResponse	"Invalid user ID"
//	@Failure		401	{object}	models.APIResponse	"Unauthorized"
//	@Failure		404	{object}	models.APIResponse	"User not found"
//	@Failure		500	{object}	models.APIResponse	"Internal server error"
//	@Router			/users/{id} [get]
func (h *UserHandler) GetUser(c *gin.Context) {
	userID := c.Param("id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, errors.NewValidationError("user_id is required"))
		return
	}

	user, err := h.service.GetUser(c.Request.Context(), userID)
	if err != nil {
		h.logger.Error("Failed to get user", zap.Error(err))
		c.JSON(http.StatusNotFound, errors.NewNotFoundError("Пользователь не найден"))
		return
	}

	c.JSON(http.StatusOK, user)
}

// Logout инвалидирует текущий токен
//
//	@Summary	Logout user
//	@Description	Invalidate access token (and optionally all sessions)
//	@Tags		Users
//	@Accept		json
//	@Produce	json
//	@Security	BearerAuth
//	@Param		request	body	LogoutRequest	true	"Revoke options"
//	@Success	200	{object}	models.APIResponse	"Logged out"
//	@Failure	401	{object}	models.APIResponse	"Unauthorized"
//	@Failure	500	{object}	models.APIResponse	"Internal server error"
//	@Router		/users/logout [post]
func (h *UserHandler) Logout(c *gin.Context) {
	var req LogoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("Invalid logout request", zap.Error(err))
		c.JSON(http.StatusBadRequest, errors.NewValidationError(err.Error()))
		return
	}

	accessToken, exists := c.Get("user_token")
	if !exists {
		h.logger.Warn("user_token not found in context")
		c.JSON(http.StatusUnauthorized, errors.NewUnauthorizedError("Требуется аутентификация"))
		return
	}

	if err := h.service.Logout(c.Request.Context(), accessToken.(string), req.RevokeAll); err != nil {
		st, ok := status.FromError(err)
		if ok && st.Code() == codes.Unauthenticated {
			c.JSON(http.StatusUnauthorized, errors.NewUnauthorizedError("Токен невалиден"))
			return
		}

		h.logger.Error("Failed to logout", zap.Error(err))
		c.JSON(http.StatusInternalServerError, errors.NewInternalServerError())
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{Success: true})
}

// AssignRole назначает роль пользователю (только admin)
//
//	@Summary		Assign role to user
//	@Description	Assign role (user/manager) to a target user. Admin only.
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id		path	string	true	"Target user ID (UUID)"
//	@Param			request	body	AssignRoleRequest	true	"Role payload (user or manager)"
//	@Success		200	{object}	models.User	"Updated user with new role"
//	@Failure		400	{object}	models.APIResponse	"Invalid input"
//	@Failure		401	{object}	models.APIResponse	"Unauthorized"
//	@Failure		403	{object}	models.APIResponse	"Forbidden"
//	@Failure		404	{object}	models.APIResponse	"User not found"
//	@Failure		500	{object}	models.APIResponse	"Internal server error"
//	@Router			/users/{id}/role [put]
func (h *UserHandler) AssignRole(c *gin.Context) {
	var req AssignRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("Invalid assign role request", zap.Error(err))
		c.JSON(http.StatusBadRequest, errors.NewValidationError(err.Error()))
		return
	}

	userID := c.Param("id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, errors.NewValidationError("user_id is required"))
		return
	}

	accessToken, exists := c.Get("user_token")
	if !exists {
		h.logger.Warn("user_token not found in context")
		c.JSON(http.StatusUnauthorized, errors.NewUnauthorizedError("Требуется аутентификация"))
		return
	}

	updatedUser, err := h.service.AssignRole(c.Request.Context(), accessToken.(string), userID, req.Role)
	if err != nil {
		st, ok := status.FromError(err)
		if ok {
			switch st.Code() {
			case codes.InvalidArgument:
				c.JSON(http.StatusBadRequest, errors.NewValidationError(st.Message()))
				return
			case codes.PermissionDenied:
				c.JSON(http.StatusForbidden, errors.NewForbiddenError("Недостаточно прав для назначения роли"))
				return
			case codes.NotFound:
				c.JSON(http.StatusNotFound, errors.NewNotFoundError("Пользователь не найден"))
				return
			case codes.Unauthenticated:
				c.JSON(http.StatusUnauthorized, errors.NewUnauthorizedError("Требуется аутентификация"))
				return
			}
		}

		h.logger.Error("Failed to assign role", zap.Error(err))
		c.JSON(http.StatusInternalServerError, errors.NewInternalServerError())
		return
	}

	c.JSON(http.StatusOK, updatedUser)
}

// UpdateProfile обновляет имя/фамилию пользователя
//
//	@Summary	Update user profile
//	@Description	Update first_name and/or last_name
//	@Tags		Users
//	@Accept		json
//	@Produce	json
//	@Security	BearerAuth
//	@Param		request	body	UpdateProfileRequest	true	"Profile fields"
//	@Success	200	{object}	models.User	"Updated user"
//	@Failure	400	{object}	models.APIResponse	"Invalid input"
//	@Failure	401	{object}	models.APIResponse	"Unauthorized"
//	@Failure	403	{object}	models.APIResponse	"Forbidden"
//	@Failure	404	{object}	models.APIResponse	"User not found"
//	@Failure	500	{object}	models.APIResponse	"Internal server error"
//	@Router		/users/profile [put]
func (h *UserHandler) UpdateProfile(c *gin.Context) {
	var req UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("Invalid update profile request", zap.Error(err))
		c.JSON(http.StatusBadRequest, errors.NewValidationError(err.Error()))
		return
	}

	accessToken, exists := c.Get("user_token")
	if !exists {
		h.logger.Warn("user_token not found in context")
		c.JSON(http.StatusUnauthorized, errors.NewUnauthorizedError("Требуется аутентификация"))
		return
	}

	user, err := h.service.UpdateProfile(c.Request.Context(), accessToken.(string), req.UserID, req.FirstName, req.LastName)
	if err != nil {
		st, ok := status.FromError(err)
		if ok {
			switch st.Code() {
			case codes.InvalidArgument:
				c.JSON(http.StatusBadRequest, errors.NewValidationError(st.Message()))
				return
			case codes.PermissionDenied:
				c.JSON(http.StatusForbidden, errors.NewForbiddenError("Недостаточно прав"))
				return
			case codes.NotFound:
				c.JSON(http.StatusNotFound, errors.NewNotFoundError("Пользователь не найден"))
				return
			case codes.Unauthenticated:
				c.JSON(http.StatusUnauthorized, errors.NewUnauthorizedError("Требуется аутентификация"))
				return
			}
		}

		h.logger.Error("Failed to update profile", zap.Error(err))
		c.JSON(http.StatusInternalServerError, errors.NewInternalServerError())
		return
	}

	c.JSON(http.StatusOK, user)
}

// ChangePassword меняет пароль пользователя
//
//	@Summary	Change password
//	@Description	Change password for current user
//	@Tags		Users
//	@Accept		json
//	@Produce	json
//	@Security	BearerAuth
//	@Param		request	body	ChangePasswordRequest	true	"Passwords payload"
//	@Success	200	{object}	models.APIResponse	"Password changed"
//	@Failure	400	{object}	models.APIResponse	"Invalid input"
//	@Failure	401	{object}	models.APIResponse	"Unauthorized"
//	@Failure	403	{object}	models.APIResponse	"Forbidden"
//	@Failure	500	{object}	models.APIResponse	"Internal server error"
//	@Router		/users/password [put]
func (h *UserHandler) ChangePassword(c *gin.Context) {
	var req ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("Invalid change password request", zap.Error(err))
		c.JSON(http.StatusBadRequest, errors.NewValidationError(err.Error()))
		return
	}

	accessToken, exists := c.Get("user_token")
	if !exists {
		h.logger.Warn("user_token not found in context")
		c.JSON(http.StatusUnauthorized, errors.NewUnauthorizedError("Требуется аутентификация"))
		return
	}

	if err := h.service.ChangePassword(c.Request.Context(), accessToken.(string), req.UserID, req.OldPassword, req.NewPassword); err != nil {
		st, ok := status.FromError(err)
		if ok {
			switch st.Code() {
			case codes.InvalidArgument:
				c.JSON(http.StatusBadRequest, errors.NewValidationError(st.Message()))
				return
			case codes.PermissionDenied:
				c.JSON(http.StatusForbidden, errors.NewForbiddenError("Недостаточно прав"))
				return
			case codes.Unauthenticated:
				c.JSON(http.StatusUnauthorized, errors.NewUnauthorizedError("Неверные учетные данные"))
				return
			}
		}

		h.logger.Error("Failed to change password", zap.Error(err))
		c.JSON(http.StatusInternalServerError, errors.NewInternalServerError())
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{Success: true})
}

// ListUsers возвращает пользователей (admin)
//
//	@Summary	List users
//	@Description	List users with pagination (admin only)
//	@Tags		Users
//	@Accept		json
//	@Produce	json
//	@Security	BearerAuth
//	@Param		page	query	int	false	"Page (0-based)"
//	@Param		size	query	int	false	"Page size (1-100)"
//	@Param		status	query	string	false	"Status filter"
//	@Param		role	query	string	false	"Role filter"
//	@Success	200	{object}	models.ListUsersResponse	"Users list"
//	@Failure	400	{object}	models.APIResponse	"Invalid input"
//	@Failure	401	{object}	models.APIResponse	"Unauthorized"
//	@Failure	403	{object}	models.APIResponse	"Forbidden"
//	@Failure	500	{object}	models.APIResponse	"Internal server error"
//	@Router		/users [get]
func (h *UserHandler) ListUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "0"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))
	statusFilter := c.Query("status")
	role := c.Query("role")

	resp, err := h.service.ListUsers(c.Request.Context(), conversion.SafeIntToInt32WithDefault(page, 0), conversion.SafeIntToInt32WithDefault(size, 20), statusFilter, role)
	if err != nil {
		st, ok := status.FromError(err)
		if ok {
			switch st.Code() {
			case codes.InvalidArgument:
				c.JSON(http.StatusBadRequest, errors.NewValidationError(st.Message()))
				return
			case codes.PermissionDenied:
				c.JSON(http.StatusForbidden, errors.NewForbiddenError("Недостаточно прав"))
				return
			case codes.Unauthenticated:
				c.JSON(http.StatusUnauthorized, errors.NewUnauthorizedError("Требуется аутентификация"))
				return
			}
		}

		h.logger.Error("Failed to list users", zap.Error(err))
		c.JSON(http.StatusInternalServerError, errors.NewInternalServerError())
		return
	}

	c.JSON(http.StatusOK, resp)
}

// ForgotPassword инициирует восстановление пароля
//
//	@Summary	Forgot password
//	@Description	Send reset instructions to email
//	@Tags		Users
//	@Accept		json
//	@Produce	json
//	@Param		request	body	ForgotPasswordRequest	true	"Email"
//	@Success	200	{object}	models.APIResponse	"Email sent"
//	@Failure	400	{object}	models.APIResponse	"Invalid input"
//	@Failure	500	{object}	models.APIResponse	"Internal server error"
//	@Router		/users/forgot-password [post]
func (h *UserHandler) ForgotPassword(c *gin.Context) {
	var req ForgotPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("Invalid forgot password request", zap.Error(err))
		c.JSON(http.StatusBadRequest, errors.NewValidationError(err.Error()))
		return
	}

	if err := h.service.ForgotPassword(c.Request.Context(), req.Email); err != nil {
		h.logger.Error("Failed to process forgot password", zap.Error(err))
		c.JSON(http.StatusInternalServerError, errors.NewInternalServerError())
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{Success: true})
}

// ResetPassword сбрасывает пароль по reset token
//
//	@Summary	Reset password
//	@Description	Reset password using reset token
//	@Tags		Users
//	@Accept		json
//	@Produce	json
//	@Param		request	body	ResetPasswordRequest	true	"Reset payload"
//	@Success	200	{object}	models.AuthResponse	"New tokens and user"
//	@Failure	400	{object}	models.APIResponse	"Invalid input"
//	@Failure	401	{object}	models.APIResponse	"Token expired"
//	@Failure	404	{object}	models.APIResponse	"User not found"
//	@Failure	500	{object}	models.APIResponse	"Internal server error"
//	@Router		/users/reset-password [post]
func (h *UserHandler) ResetPassword(c *gin.Context) {
	var req ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("Invalid reset password request", zap.Error(err))
		c.JSON(http.StatusBadRequest, errors.NewValidationError(err.Error()))
		return
	}

	resp, err := h.service.ResetPassword(c.Request.Context(), req.ResetToken, req.NewPassword)
	if err != nil {
		st, ok := status.FromError(err)
		if ok {
			switch st.Code() {
			case codes.InvalidArgument:
				c.JSON(http.StatusBadRequest, errors.NewValidationError(st.Message()))
				return
			case codes.Unauthenticated:
				c.JSON(http.StatusUnauthorized, errors.NewUnauthorizedError("Reset token истек"))
				return
			case codes.NotFound:
				c.JSON(http.StatusNotFound, errors.NewNotFoundError("Пользователь не найден"))
				return
			}
		}

		h.logger.Error("Failed to reset password", zap.Error(err))
		c.JSON(http.StatusInternalServerError, errors.NewInternalServerError())
		return
	}

	c.JSON(http.StatusOK, resp)
}
