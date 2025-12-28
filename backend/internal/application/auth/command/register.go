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

type RegisterCommand struct {
	Email     string
	Password  string
	FullName  string
	IPAddress net.IP
	UserAgent string
}

type RegisterHandler struct {
	userRepository         user.Repository
	roleRepository         role.Repository
	permissionRepository   permission.Repository
	refreshTokenRepository auth.RefreshTokenRepository
	tokenGenerator         auth.TokenGenerator
	passwordHasher         security.PasswordHasher
	eventBus               shared.EventBus
	refreshTokenTTL        time.Duration
	logger                 logger.Logger
}

type RegisterHandlerParams struct {
	UserRepository         user.Repository
	RoleRepository         role.Repository
	PermissionRepository   permission.Repository
	RefreshTokenRepository auth.RefreshTokenRepository
	TokenGenerator         auth.TokenGenerator
	PasswordHasher         security.PasswordHasher
	EventBus               shared.EventBus
	RefreshTokenTTL        time.Duration
	Logger                 logger.Logger
}

func NewRegisterHandler(params RegisterHandlerParams) *RegisterHandler {
	return &RegisterHandler{
		userRepository:         params.UserRepository,
		roleRepository:         params.RoleRepository,
		permissionRepository:   params.PermissionRepository,
		refreshTokenRepository: params.RefreshTokenRepository,
		tokenGenerator:         params.TokenGenerator,
		passwordHasher:         params.PasswordHasher,
		eventBus:               params.EventBus,
		refreshTokenTTL:        params.RefreshTokenTTL,
		logger:                 params.Logger,
	}
}

func (handler *RegisterHandler) Handle(ctx context.Context, command RegisterCommand) (*authdto.AuthResponseDTO, error) {
	if err := shared.ValidatePassword(command.Password); err != nil {
		return nil, err
	}

	exists, err := handler.userRepository.ExistsByEmail(ctx, command.Email)
	if err != nil {
		return nil, fmt.Errorf("check email exists: %w", err)
	}
	if exists {
		return nil, user.NewEmailAlreadyExistsError(command.Email)
	}

	hashedPassword, err := handler.passwordHasher.Hash(command.Password)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}

	newUser, err := user.NewUser(user.NewUserParams{
		Email:        command.Email,
		PasswordHash: hashedPassword,
		FullName:     command.FullName,
	})
	if err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}

	defaultRole, err := handler.roleRepository.FindDefault(ctx)
	if err != nil {
		handler.logger.Warn("no default role found, proceeding without role assignment",
			logger.Err(err),
		)
	} else if defaultRole != nil {
		if err := newUser.AssignRole(defaultRole.ID()); err != nil {
			handler.logger.Warn("failed to assign default role",
				logger.String("role_id", defaultRole.ID().String()),
				logger.Err(err),
			)
		}
	}

	if err := newUser.Activate(); err != nil {
		handler.logger.Warn("failed to activate user",
			logger.Err(err),
		)
	}

	if err := handler.userRepository.Create(ctx, newUser); err != nil {
		return nil, fmt.Errorf("save user: %w", err)
	}

	roles, permissions := handler.loadUserRolesAndPermissions(ctx, newUser)

	accessToken, err := handler.tokenGenerator.GenerateAccessToken(
		newUser.ID(),
		newUser.Email().String(),
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
		UserID:     newUser.ID(),
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
		events := newUser.DomainEvents()
		events = append(events, auth.NewUserRegisteredEvent(
			newUser.ID(),
			newUser.Email().String(),
			newUser.FullName().String(),
		))

		if err := handler.eventBus.Publish(ctx, events...); err != nil {
			handler.logger.Error("failed to publish domain events",
				logger.String("user_id", newUser.ID().String()),
				logger.Err(err),
			)
		}
		newUser.ClearDomainEvents()
	}

	handler.logger.Info("user registered successfully",
		logger.String("user_id", newUser.ID().String()),
		logger.String("email", newUser.Email().String()),
	)

	return &authdto.AuthResponseDTO{
		User:         userdto.UserFromDomain(newUser),
		AccessToken:  accessToken.Token(),
		RefreshToken: refreshTokenString,
		ExpiresAt:    accessToken.ExpiresAt(),
	}, nil
}

func (handler *RegisterHandler) loadUserRolesAndPermissions(ctx context.Context, domainUser *user.User) ([]string, []string) {
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
