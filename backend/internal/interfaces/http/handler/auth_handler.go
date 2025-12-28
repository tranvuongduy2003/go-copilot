package handler

import (
	"encoding/json"
	"net"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	authcommand "github.com/tranvuongduy2003/go-copilot/internal/application/auth/command"
	authquery "github.com/tranvuongduy2003/go-copilot/internal/application/auth/query"
	"github.com/tranvuongduy2003/go-copilot/internal/interfaces/http/dto"
	"github.com/tranvuongduy2003/go-copilot/internal/interfaces/http/middleware"
	"github.com/tranvuongduy2003/go-copilot/internal/interfaces/http/response"
	"github.com/tranvuongduy2003/go-copilot/pkg/logger"
	"github.com/tranvuongduy2003/go-copilot/pkg/validator"
)

type AuthHandler struct {
	registerHandler        *authcommand.RegisterHandler
	loginHandler           *authcommand.LoginHandler
	refreshTokenHandler    *authcommand.RefreshTokenHandler
	logoutHandler          *authcommand.LogoutHandler
	forgotPasswordHandler  *authcommand.ForgotPasswordHandler
	resetPasswordHandler   *authcommand.ResetPasswordHandler
	revokeSessionHandler   *authcommand.RevokeSessionHandler
	getCurrentUserHandler  *authquery.GetCurrentUserHandler
	getUserSessionsHandler *authquery.GetUserSessionsHandler
	validator              *validator.Validator
	logger                 logger.Logger
}

type AuthHandlerParams struct {
	RegisterHandler        *authcommand.RegisterHandler
	LoginHandler           *authcommand.LoginHandler
	RefreshTokenHandler    *authcommand.RefreshTokenHandler
	LogoutHandler          *authcommand.LogoutHandler
	ForgotPasswordHandler  *authcommand.ForgotPasswordHandler
	ResetPasswordHandler   *authcommand.ResetPasswordHandler
	RevokeSessionHandler   *authcommand.RevokeSessionHandler
	GetCurrentUserHandler  *authquery.GetCurrentUserHandler
	GetUserSessionsHandler *authquery.GetUserSessionsHandler
	Validator              *validator.Validator
	Logger                 logger.Logger
}

func NewAuthHandler(params AuthHandlerParams) *AuthHandler {
	return &AuthHandler{
		registerHandler:        params.RegisterHandler,
		loginHandler:           params.LoginHandler,
		refreshTokenHandler:    params.RefreshTokenHandler,
		logoutHandler:          params.LogoutHandler,
		forgotPasswordHandler:  params.ForgotPasswordHandler,
		resetPasswordHandler:   params.ResetPasswordHandler,
		revokeSessionHandler:   params.RevokeSessionHandler,
		getCurrentUserHandler:  params.GetCurrentUserHandler,
		getUserSessionsHandler: params.GetUserSessionsHandler,
		validator:              params.Validator,
		logger:                 params.Logger,
	}
}

func (handler *AuthHandler) Register(writer http.ResponseWriter, request *http.Request) {
	var requestBody dto.RegisterRequest
	if err := json.NewDecoder(request.Body).Decode(&requestBody); err != nil {
		response.BadRequest(writer, request, "invalid request body")
		return
	}

	if err := handler.validator.Validate(requestBody); err != nil {
		if validationErrors, ok := validator.GetValidationErrors(err); ok {
			response.ValidationError(writer, request, validationErrors)
			return
		}
		response.BadRequest(writer, request, err.Error())
		return
	}

	cmd := authcommand.RegisterCommand{
		Email:     requestBody.Email,
		Password:  requestBody.Password,
		FullName:  requestBody.FullName,
		IPAddress: handler.getClientIP(request),
		UserAgent: request.UserAgent(),
	}

	result, err := handler.registerHandler.Handle(request.Context(), cmd)
	if err != nil {
		response.Error(writer, request, err)
		return
	}

	response.Created(writer, dto.AuthResponse{
		User: dto.UserResponse{
			ID:        result.User.ID,
			Email:     result.User.Email,
			FullName:  result.User.FullName,
			Status:    result.User.Status,
			CreatedAt: result.User.CreatedAt,
			UpdatedAt: result.User.UpdatedAt,
			DeletedAt: result.User.DeletedAt,
		},
		AccessToken:  result.AccessToken,
		RefreshToken: result.RefreshToken,
		ExpiresAt:    result.ExpiresAt,
	})
}

func (handler *AuthHandler) Login(writer http.ResponseWriter, request *http.Request) {
	var requestBody dto.LoginRequest
	if err := json.NewDecoder(request.Body).Decode(&requestBody); err != nil {
		response.BadRequest(writer, request, "invalid request body")
		return
	}

	if err := handler.validator.Validate(requestBody); err != nil {
		if validationErrors, ok := validator.GetValidationErrors(err); ok {
			response.ValidationError(writer, request, validationErrors)
			return
		}
		response.BadRequest(writer, request, err.Error())
		return
	}

	cmd := authcommand.LoginCommand{
		Email:     requestBody.Email,
		Password:  requestBody.Password,
		IPAddress: handler.getClientIP(request),
		UserAgent: request.UserAgent(),
	}

	result, err := handler.loginHandler.Handle(request.Context(), cmd)
	if err != nil {
		response.Error(writer, request, err)
		return
	}

	response.Success(writer, dto.AuthResponse{
		User: dto.UserResponse{
			ID:        result.User.ID,
			Email:     result.User.Email,
			FullName:  result.User.FullName,
			Status:    result.User.Status,
			CreatedAt: result.User.CreatedAt,
			UpdatedAt: result.User.UpdatedAt,
			DeletedAt: result.User.DeletedAt,
		},
		AccessToken:  result.AccessToken,
		RefreshToken: result.RefreshToken,
		ExpiresAt:    result.ExpiresAt,
	})
}

func (handler *AuthHandler) RefreshToken(writer http.ResponseWriter, request *http.Request) {
	var requestBody dto.RefreshTokenRequest
	if err := json.NewDecoder(request.Body).Decode(&requestBody); err != nil {
		response.BadRequest(writer, request, "invalid request body")
		return
	}

	if err := handler.validator.Validate(requestBody); err != nil {
		if validationErrors, ok := validator.GetValidationErrors(err); ok {
			response.ValidationError(writer, request, validationErrors)
			return
		}
		response.BadRequest(writer, request, err.Error())
		return
	}

	cmd := authcommand.RefreshTokenCommand{
		RefreshToken: requestBody.RefreshToken,
		IPAddress:    handler.getClientIP(request),
		UserAgent:    request.UserAgent(),
	}

	result, err := handler.refreshTokenHandler.Handle(request.Context(), cmd)
	if err != nil {
		response.Error(writer, request, err)
		return
	}

	response.Success(writer, dto.AuthResponse{
		User: dto.UserResponse{
			ID:        result.User.ID,
			Email:     result.User.Email,
			FullName:  result.User.FullName,
			Status:    result.User.Status,
			CreatedAt: result.User.CreatedAt,
			UpdatedAt: result.User.UpdatedAt,
			DeletedAt: result.User.DeletedAt,
		},
		AccessToken:  result.AccessToken,
		RefreshToken: result.RefreshToken,
		ExpiresAt:    result.ExpiresAt,
	})
}

func (handler *AuthHandler) Logout(writer http.ResponseWriter, request *http.Request) {
	authContext, ok := middleware.GetAuthContext(request.Context())
	if !ok {
		response.Unauthorized(writer, request, "unauthorized")
		return
	}

	var requestBody dto.LogoutRequest
	if request.ContentLength > 0 {
		if err := json.NewDecoder(request.Body).Decode(&requestBody); err != nil {
			response.BadRequest(writer, request, "invalid request body")
			return
		}
	}

	cmd := authcommand.LogoutCommand{
		UserID:    authContext.UserID,
		TokenID:   authContext.TokenID,
		ExpiresAt: authContext.ExpiresAt.Unix(),
		LogoutAll: requestBody.LogoutAll,
	}

	if err := handler.logoutHandler.Handle(request.Context(), cmd); err != nil {
		response.Error(writer, request, err)
		return
	}

	response.SuccessWithMessage(writer, nil, "logged out successfully")
}

func (handler *AuthHandler) LogoutAll(writer http.ResponseWriter, request *http.Request) {
	authContext, ok := middleware.GetAuthContext(request.Context())
	if !ok {
		response.Unauthorized(writer, request, "unauthorized")
		return
	}

	cmd := authcommand.LogoutCommand{
		UserID:    authContext.UserID,
		TokenID:   authContext.TokenID,
		ExpiresAt: authContext.ExpiresAt.Unix(),
		LogoutAll: true,
	}

	if err := handler.logoutHandler.Handle(request.Context(), cmd); err != nil {
		response.Error(writer, request, err)
		return
	}

	response.SuccessWithMessage(writer, nil, "logged out from all devices successfully")
}

func (handler *AuthHandler) ForgotPassword(writer http.ResponseWriter, request *http.Request) {
	var requestBody dto.ForgotPasswordRequest
	if err := json.NewDecoder(request.Body).Decode(&requestBody); err != nil {
		response.BadRequest(writer, request, "invalid request body")
		return
	}

	if err := handler.validator.Validate(requestBody); err != nil {
		if validationErrors, ok := validator.GetValidationErrors(err); ok {
			response.ValidationError(writer, request, validationErrors)
			return
		}
		response.BadRequest(writer, request, err.Error())
		return
	}

	cmd := authcommand.ForgotPasswordCommand{
		Email: requestBody.Email,
	}

	_, err := handler.forgotPasswordHandler.Handle(request.Context(), cmd)
	if err != nil {
		response.Error(writer, request, err)
		return
	}

	response.SuccessWithMessage(writer, nil, "if the email exists, a password reset link has been sent")
}

func (handler *AuthHandler) ResetPassword(writer http.ResponseWriter, request *http.Request) {
	var requestBody dto.ResetPasswordRequest
	if err := json.NewDecoder(request.Body).Decode(&requestBody); err != nil {
		response.BadRequest(writer, request, "invalid request body")
		return
	}

	if err := handler.validator.Validate(requestBody); err != nil {
		if validationErrors, ok := validator.GetValidationErrors(err); ok {
			response.ValidationError(writer, request, validationErrors)
			return
		}
		response.BadRequest(writer, request, err.Error())
		return
	}

	cmd := authcommand.ResetPasswordCommand{
		ResetToken:  requestBody.ResetToken,
		NewPassword: requestBody.NewPassword,
	}

	if err := handler.resetPasswordHandler.Handle(request.Context(), cmd); err != nil {
		response.Error(writer, request, err)
		return
	}

	response.SuccessWithMessage(writer, nil, "password reset successfully")
}

func (handler *AuthHandler) GetCurrentUser(writer http.ResponseWriter, request *http.Request) {
	authContext, ok := middleware.GetAuthContext(request.Context())
	if !ok {
		response.Unauthorized(writer, request, "unauthorized")
		return
	}

	userQuery := authquery.GetCurrentUserQuery{UserID: authContext.UserID}
	result, err := handler.getCurrentUserHandler.Handle(request.Context(), userQuery)
	if err != nil {
		response.Error(writer, request, err)
		return
	}

	response.Success(writer, dto.AuthUserResponse{
		ID:          result.ID,
		Email:       result.Email,
		FullName:    result.FullName,
		Status:      result.Status,
		Roles:       result.Roles,
		Permissions: result.Permissions,
	})
}

func (handler *AuthHandler) GetSessions(writer http.ResponseWriter, request *http.Request) {
	authContext, ok := middleware.GetAuthContext(request.Context())
	if !ok {
		response.Unauthorized(writer, request, "unauthorized")
		return
	}

	currentTokenID, _ := uuid.Parse(authContext.TokenID)

	sessionQuery := authquery.GetUserSessionsQuery{
		UserID:         authContext.UserID,
		CurrentTokenID: currentTokenID,
	}

	sessions, err := handler.getUserSessionsHandler.Handle(request.Context(), sessionQuery)
	if err != nil {
		response.Error(writer, request, err)
		return
	}

	sessionResponses := make([]dto.SessionResponse, len(sessions))
	for i, session := range sessions {
		sessionResponses[i] = dto.SessionResponse{
			ID: session.ID,
			DeviceInfo: dto.DeviceResponse{
				UserAgent: session.DeviceInfo.UserAgent,
				Platform:  session.DeviceInfo.Platform,
				Browser:   session.DeviceInfo.Browser,
			},
			IPAddress:  session.IPAddress,
			CreatedAt:  session.CreatedAt,
			LastUsedAt: session.LastUsedAt,
			IsCurrent:  session.IsCurrent,
		}
	}

	response.Success(writer, sessionResponses)
}

func (handler *AuthHandler) RevokeSession(writer http.ResponseWriter, request *http.Request) {
	authContext, ok := middleware.GetAuthContext(request.Context())
	if !ok {
		response.Unauthorized(writer, request, "unauthorized")
		return
	}

	sessionIDStr := chi.URLParam(request, "id")
	if sessionIDStr == "" {
		response.BadRequest(writer, request, "session id is required")
		return
	}

	sessionID, err := uuid.Parse(sessionIDStr)
	if err != nil {
		response.BadRequest(writer, request, "invalid session id")
		return
	}

	cmd := authcommand.RevokeSessionCommand{
		UserID:    authContext.UserID,
		SessionID: sessionID,
	}

	if err := handler.revokeSessionHandler.Handle(request.Context(), cmd); err != nil {
		response.Error(writer, request, err)
		return
	}

	response.NoContent(writer)
}

func (handler *AuthHandler) getClientIP(request *http.Request) net.IP {
	xff := request.Header.Get("X-Forwarded-For")
	if xff != "" {
		parts := strings.Split(xff, ",")
		if len(parts) > 0 {
			ip := net.ParseIP(strings.TrimSpace(parts[0]))
			if ip != nil {
				return ip
			}
		}
	}

	xri := request.Header.Get("X-Real-IP")
	if xri != "" {
		ip := net.ParseIP(xri)
		if ip != nil {
			return ip
		}
	}

	host, _, err := net.SplitHostPort(request.RemoteAddr)
	if err != nil {
		return net.ParseIP(request.RemoteAddr)
	}

	return net.ParseIP(host)
}
