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
)

type RefreshTokenCommand struct {
	RefreshToken string
	IPAddress    net.IP
	UserAgent    string
}

type RefreshTokenHandler struct {
	userRepository         user.Repository
	roleRepository         role.Repository
	permissionRepository   permission.Repository
	refreshTokenRepository auth.RefreshTokenRepository
	tokenGenerator         auth.TokenGenerator
	tokenBlacklist         auth.TokenBlacklist
	eventBus               shared.EventBus
	refreshTokenTTL        time.Duration
	logger                 logger.Logger
}

type RefreshTokenHandlerParams struct {
	UserRepository         user.Repository
	RoleRepository         role.Repository
	PermissionRepository   permission.Repository
	RefreshTokenRepository auth.RefreshTokenRepository
	TokenGenerator         auth.TokenGenerator
	TokenBlacklist         auth.TokenBlacklist
	EventBus               shared.EventBus
	RefreshTokenTTL        time.Duration
	Logger                 logger.Logger
}

func NewRefreshTokenHandler(params RefreshTokenHandlerParams) *RefreshTokenHandler {
	return &RefreshTokenHandler{
		userRepository:         params.UserRepository,
		roleRepository:         params.RoleRepository,
		permissionRepository:   params.PermissionRepository,
		refreshTokenRepository: params.RefreshTokenRepository,
		tokenGenerator:         params.TokenGenerator,
		tokenBlacklist:         params.TokenBlacklist,
		eventBus:               params.EventBus,
		refreshTokenTTL:        params.RefreshTokenTTL,
		logger:                 params.Logger,
	}
}

func (handler *RefreshTokenHandler) Handle(ctx context.Context, command RefreshTokenCommand) (*authdto.AuthResponseDTO, error) {
	tokenHash := handler.tokenGenerator.HashRefreshToken(command.RefreshToken)

	existingToken, err := handler.refreshTokenRepository.FindByTokenHash(ctx, tokenHash)
	if err != nil {
		return nil, auth.ErrRefreshTokenInvalid
	}

	if !existingToken.IsValid() {
		return nil, auth.ErrRefreshTokenInvalid
	}

	existingUser, err := handler.userRepository.FindByID(ctx, existingToken.UserID())
	if err != nil {
		return nil, auth.ErrRefreshTokenInvalid
	}

	if !existingUser.Status().IsActive() {
		return nil, auth.ErrAccountInactive
	}

	if err := handler.refreshTokenRepository.Revoke(ctx, existingToken.ID()); err != nil {
		return nil, fmt.Errorf("revoke old refresh token: %w", err)
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

	newRefreshTokenString, err := handler.tokenGenerator.GenerateRefreshToken()
	if err != nil {
		return nil, fmt.Errorf("generate refresh token: %w", err)
	}

	newRefreshTokenHash := handler.tokenGenerator.HashRefreshToken(newRefreshTokenString)
	deviceInfo := &auth.DeviceInfo{
		UserAgent: command.UserAgent,
	}

	newRefreshToken, err := auth.NewRefreshToken(auth.NewRefreshTokenParams{
		UserID:     existingUser.ID(),
		TokenHash:  newRefreshTokenHash,
		ExpiresAt:  time.Now().UTC().Add(handler.refreshTokenTTL),
		DeviceInfo: deviceInfo,
		IPAddress:  command.IPAddress,
	})
	if err != nil {
		return nil, fmt.Errorf("create refresh token: %w", err)
	}

	if err := handler.refreshTokenRepository.Create(ctx, newRefreshToken); err != nil {
		return nil, fmt.Errorf("save refresh token: %w", err)
	}

	if handler.eventBus != nil {
		event := auth.NewRefreshTokenRotatedEvent(
			existingUser.ID(),
			existingToken.ID(),
			newRefreshToken.ID(),
		)
		if err := handler.eventBus.Publish(ctx, event); err != nil {
			handler.logger.Error("failed to publish refresh token rotated event",
				logger.String("user_id", existingUser.ID().String()),
				logger.Err(err),
			)
		}
	}

	handler.logger.Info("token refreshed successfully",
		logger.String("user_id", existingUser.ID().String()),
	)

	return &authdto.AuthResponseDTO{
		User:         userdto.UserFromDomain(existingUser),
		AccessToken:  accessToken.Token(),
		RefreshToken: newRefreshTokenString,
		ExpiresAt:    accessToken.ExpiresAt(),
	}, nil
}

func (handler *RefreshTokenHandler) loadUserRolesAndPermissions(ctx context.Context, domainUser *user.User) ([]string, []string) {
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
		for _, permissionID := range roleEntity.PermissionIDs() {
			permissionIDSet[permissionID] = true
		}
	}

	if len(permissionIDSet) == 0 {
		return roleNames, []string{}
	}

	permissionIDs := make([]uuid.UUID, 0, len(permissionIDSet))
	for permissionID := range permissionIDSet {
		permissionIDs = append(permissionIDs, permissionID)
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
	for _, permissionEntity := range permissions {
		permissionCodes = append(permissionCodes, permissionEntity.Code().String())
	}

	return roleNames, permissionCodes
}
