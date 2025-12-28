package authcommand

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/google/uuid"

	authdto "github.com/tranvuongduy2003/go-copilot/internal/application/auth/dto"
	userdto "github.com/tranvuongduy2003/go-copilot/internal/application/user/dto"
	"github.com/tranvuongduy2003/go-copilot/internal/domain/auth"
	"github.com/tranvuongduy2003/go-copilot/internal/domain/permission"
	"github.com/tranvuongduy2003/go-copilot/internal/domain/role"
	"github.com/tranvuongduy2003/go-copilot/internal/domain/shared"
	"github.com/tranvuongduy2003/go-copilot/internal/domain/user"
	"github.com/tranvuongduy2003/go-copilot/pkg/logger"
	"github.com/tranvuongduy2003/go-copilot/pkg/security"
)

type LoginCommand struct {
	Email     string
	Password  string
	IPAddress net.IP
	UserAgent string
}

type LoginHandler struct {
	userRepository         user.Repository
	roleRepository         role.Repository
	permissionRepository   permission.Repository
	refreshTokenRepository auth.RefreshTokenRepository
	tokenGenerator         auth.TokenGenerator
	passwordHasher         security.PasswordHasher
	accountLockout         security.AccountLockout
	eventBus               shared.EventBus
	refreshTokenTTL        time.Duration
	logger                 logger.Logger
}

type LoginHandlerParams struct {
	UserRepository         user.Repository
	RoleRepository         role.Repository
	PermissionRepository   permission.Repository
	RefreshTokenRepository auth.RefreshTokenRepository
	TokenGenerator         auth.TokenGenerator
	PasswordHasher         security.PasswordHasher
	AccountLockout         security.AccountLockout
	EventBus               shared.EventBus
	RefreshTokenTTL        time.Duration
	Logger                 logger.Logger
}

func NewLoginHandler(params LoginHandlerParams) *LoginHandler {
	return &LoginHandler{
		userRepository:         params.UserRepository,
		roleRepository:         params.RoleRepository,
		permissionRepository:   params.PermissionRepository,
		refreshTokenRepository: params.RefreshTokenRepository,
		tokenGenerator:         params.TokenGenerator,
		passwordHasher:         params.PasswordHasher,
		accountLockout:         params.AccountLockout,
		eventBus:               params.EventBus,
		refreshTokenTTL:        params.RefreshTokenTTL,
		logger:                 params.Logger,
	}
}

func (handler *LoginHandler) Handle(ctx context.Context, command LoginCommand) (*authdto.AuthResponseDTO, error) {
	lockoutIdentifier := command.Email

	if handler.accountLockout != nil {
		locked, remainingTime, err := handler.accountLockout.IsLocked(ctx, lockoutIdentifier)
		if err != nil {
			handler.logger.Error("failed to check account lockout status",
				logger.String("email", command.Email),
				logger.Err(err),
			)
		}
		if locked {
			handler.logger.Warn("login attempt on locked account",
				logger.String("email", command.Email),
				logger.String("remaining_lockout", remainingTime.String()),
			)
			return nil, auth.ErrAccountLocked
		}
	}

	existingUser, err := handler.userRepository.FindByEmail(ctx, command.Email)
	if err != nil {
		if handler.accountLockout != nil {
			_, _ = handler.accountLockout.RecordFailedAttempt(ctx, lockoutIdentifier)
		}
		return nil, auth.ErrInvalidCredentials
	}

	if !existingUser.Status().IsActive() {
		return nil, auth.ErrAccountInactive
	}

	valid, err := handler.passwordHasher.Verify(existingUser.PasswordHash().String(), command.Password)
	if err != nil || !valid {
		if handler.accountLockout != nil {
			attemptCount, _ := handler.accountLockout.RecordFailedAttempt(ctx, lockoutIdentifier)
			handler.publishLoginFailedEvent(ctx, existingUser, command.IPAddress.String(), "invalid_password", attemptCount)
		} else {
			handler.publishLoginFailedEvent(ctx, existingUser, command.IPAddress.String(), "invalid_password", 0)
		}
		return nil, auth.ErrInvalidCredentials
	}

	if handler.accountLockout != nil {
		if err := handler.accountLockout.ResetAttempts(ctx, lockoutIdentifier); err != nil {
			handler.logger.Error("failed to reset login attempts",
				logger.String("email", command.Email),
				logger.Err(err),
			)
		}
	}

	roles, permissions := handler.loadUserRolesAndPermissions(ctx, existingUser)

	accessToken, err := handler.tokenGenerator.GenerateAccessToken(
		existingUser.ID(),
		existingUser.Email().String(),
		roles,
		permissions,
	)
	if err != nil {
		return nil, fmt.Errorf("generate access token: %w", err)
	}

	refreshTokenString, err := handler.tokenGenerator.GenerateRefreshToken()
	if err != nil {
		return nil, fmt.Errorf("generate refresh token: %w", err)
	}

	refreshTokenHash := handler.tokenGenerator.HashRefreshToken(refreshTokenString)
	deviceInfo := &auth.DeviceInfo{
		UserAgent: command.UserAgent,
	}

	refreshToken, err := auth.NewRefreshToken(auth.NewRefreshTokenParams{
		UserID:     existingUser.ID(),
		TokenHash:  refreshTokenHash,
		ExpiresAt:  time.Now().UTC().Add(handler.refreshTokenTTL),
		DeviceInfo: deviceInfo,
		IPAddress:  command.IPAddress,
	})
	if err != nil {
		return nil, fmt.Errorf("create refresh token: %w", err)
	}

	if err := handler.refreshTokenRepository.Create(ctx, refreshToken); err != nil {
		return nil, fmt.Errorf("save refresh token: %w", err)
	}

	if handler.eventBus != nil {
		event := auth.NewUserLoggedInEvent(
			existingUser.ID(),
			existingUser.Email().String(),
			command.IPAddress.String(),
			command.UserAgent,
		)
		if err := handler.eventBus.Publish(ctx, event); err != nil {
			handler.logger.Error("failed to publish login event",
				logger.String("user_id", existingUser.ID().String()),
				logger.Err(err),
			)
		}
	}

	handler.logger.Info("user logged in successfully",
		logger.String("user_id", existingUser.ID().String()),
		logger.String("email", existingUser.Email().String()),
	)

	return &authdto.AuthResponseDTO{
		User:         userdto.UserFromDomain(existingUser),
		AccessToken:  accessToken.Token(),
		RefreshToken: refreshTokenString,
		ExpiresAt:    accessToken.ExpiresAt(),
	}, nil
}

func (handler *LoginHandler) loadUserRolesAndPermissions(ctx context.Context, domainUser *user.User) ([]string, []string) {
	roleIDs := domainUser.RoleIDs()
	if len(roleIDs) == 0 {
		return []string{}, []string{}
	}

	roles, err := handler.roleRepository.FindByIDs(ctx, roleIDs)
	if err != nil {
		handler.logger.Error("failed to load user roles",
			logger.String("user_id", domainUser.ID().String()),
			logger.Err(err),
		)
		return []string{}, []string{}
	}

	roleNames := make([]string, 0, len(roles))
	permissionIDSet := make(map[uuid.UUID]bool)

	for _, roleEntity := range roles {
		roleNames = append(roleNames, roleEntity.Name())
		for _, permID := range roleEntity.PermissionIDs() {
			permissionIDSet[permID] = true
		}
	}

	permissionIDs := make([]uuid.UUID, 0, len(permissionIDSet))
	for permID := range permissionIDSet {
		permissionIDs = append(permissionIDs, permID)
	}

	if len(permissionIDs) == 0 {
		return roleNames, []string{}
	}

	permissions, err := handler.permissionRepository.FindByIDs(ctx, permissionIDs)
	if err != nil {
		handler.logger.Error("failed to load permissions",
			logger.String("user_id", domainUser.ID().String()),
			logger.Err(err),
		)
		return roleNames, []string{}
	}

	permissionCodes := make([]string, 0, len(permissions))
	for _, perm := range permissions {
		permissionCodes = append(permissionCodes, perm.CodeString())
	}

	return roleNames, permissionCodes
}

func (handler *LoginHandler) publishLoginFailedEvent(ctx context.Context, domainUser *user.User, ipAddress, reason string, attemptCount int) {
	if handler.eventBus == nil {
		return
	}

	event := auth.NewLoginFailedEvent(
		domainUser.ID(),
		domainUser.Email().String(),
		ipAddress,
		reason,
		attemptCount,
	)

	if err := handler.eventBus.Publish(ctx, event); err != nil {
		handler.logger.Error("failed to publish login failed event",
			logger.String("user_id", domainUser.ID().String()),
			logger.Err(err),
		)
	}
}
