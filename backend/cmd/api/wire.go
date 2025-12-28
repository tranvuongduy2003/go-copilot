//go:build wireinject
// +build wireinject

package main

import (
	"net/http"

	"github.com/google/wire"

	authcommand "github.com/tranvuongduy2003/go-copilot/internal/application/auth/command"
	authquery "github.com/tranvuongduy2003/go-copilot/internal/application/auth/query"
	permissioncommand "github.com/tranvuongduy2003/go-copilot/internal/application/permission/command"
	permissionquery "github.com/tranvuongduy2003/go-copilot/internal/application/permission/query"
	rolecommand "github.com/tranvuongduy2003/go-copilot/internal/application/role/command"
	rolequery "github.com/tranvuongduy2003/go-copilot/internal/application/role/query"
	usercommand "github.com/tranvuongduy2003/go-copilot/internal/application/user/command"
	userquery "github.com/tranvuongduy2003/go-copilot/internal/application/user/query"
	"github.com/tranvuongduy2003/go-copilot/internal/domain/auth"
	"github.com/tranvuongduy2003/go-copilot/internal/domain/permission"
	"github.com/tranvuongduy2003/go-copilot/internal/domain/role"
	"github.com/tranvuongduy2003/go-copilot/internal/domain/shared"
	"github.com/tranvuongduy2003/go-copilot/internal/domain/user"
	"github.com/tranvuongduy2003/go-copilot/internal/infrastructure/audit"
	"github.com/tranvuongduy2003/go-copilot/internal/infrastructure/cache/redis"
	"github.com/tranvuongduy2003/go-copilot/internal/infrastructure/messaging/memory"
	"github.com/tranvuongduy2003/go-copilot/internal/infrastructure/persistence/postgres"
	"github.com/tranvuongduy2003/go-copilot/internal/infrastructure/persistence/repository"
	"github.com/tranvuongduy2003/go-copilot/internal/interfaces/http/handler"
	"github.com/tranvuongduy2003/go-copilot/internal/interfaces/http/middleware"
	"github.com/tranvuongduy2003/go-copilot/internal/interfaces/http/router"
	"github.com/tranvuongduy2003/go-copilot/pkg/config"
	"github.com/tranvuongduy2003/go-copilot/pkg/logger"
	"github.com/tranvuongduy2003/go-copilot/pkg/security"
	"github.com/tranvuongduy2003/go-copilot/pkg/validator"
)

func provideLogger() logger.Logger {
	return logger.L()
}

func provideUserRepository(database *postgres.DB) *repository.UserRepository {
	return repository.NewUserRepository(database.Pool())
}

func providePermissionRepository(database *postgres.DB) *repository.PermissionRepository {
	return repository.NewPermissionRepository(database.Pool())
}

func provideRoleRepository(database *postgres.DB) *repository.RoleRepository {
	return repository.NewRoleRepository(database.Pool())
}

func provideRefreshTokenRepository(database *postgres.DB) *repository.RefreshTokenRepository {
	return repository.NewRefreshTokenRepository(database.Pool())
}

func providePasswordHasher() security.PasswordHasher {
	return security.NewDefaultPasswordHasher()
}

func provideValidator() *validator.Validator {
	return validator.New()
}

func provideTokenGenerator(cfg *config.Config) auth.TokenGenerator {
	return security.NewJWTTokenGenerator(security.JWTConfig{
		SecretKey:       cfg.JWT.Secret,
		AccessTokenTTL:  cfg.JWT.AccessTokenTTL,
		RefreshTokenTTL: cfg.JWT.RefreshTokenTTL,
		Issuer:          cfg.JWT.Issuer,
		Audience:        cfg.JWT.Audience,
	})
}

func provideTokenBlacklist(redisClient *redis.Client) auth.TokenBlacklist {
	return security.NewRedisTokenBlacklist(redisClient.Client())
}

func providePasswordResetTokenStore(redisClient *redis.Client) authcommand.PasswordResetTokenStore {
	return security.NewRedisPasswordResetTokenStore(redisClient.Client())
}

func provideAuthMiddleware(tokenGenerator auth.TokenGenerator, tokenBlacklist auth.TokenBlacklist) *middleware.AuthMiddleware {
	return middleware.NewAuthMiddleware(tokenGenerator, tokenBlacklist)
}

func provideAccountLockout(redisClient *redis.Client) security.AccountLockout {
	return security.NewRedisAccountLockout(redisClient.Client(), security.AccountLockoutConfig{
		MaxAttempts:     5,
		LockoutDuration: 15 * 60 * 1000000000,
		AttemptWindow:   15 * 60 * 1000000000,
	})
}

func provideAuditLogger(database *postgres.DB, log logger.Logger) audit.AuditLogger {
	loggerAudit := audit.NewLoggerAuditLogger(log)
	postgresAudit := audit.NewPostgresAuditLogger(database.Pool())
	return audit.NewCompositeAuditLogger(loggerAudit, postgresAudit)
}

func provideAuthAuditHandler(auditLogger audit.AuditLogger, log logger.Logger) *audit.AuthAuditHandler {
	return audit.NewAuthAuditHandler(auditLogger, log)
}

func provideEventBusWithAudit(log logger.Logger, auditHandler *audit.AuthAuditHandler) *memory.InMemoryEventBus {
	eventBus := memory.NewInMemoryEventBus(log)
	for _, eventType := range auditHandler.SubscribedEventTypes() {
		eventBus.Subscribe(eventType, auditHandler.HandleEvent)
	}
	return eventBus
}

func provideRegisterHandler(
	userRepo user.Repository,
	roleRepo role.Repository,
	permissionRepo permission.Repository,
	refreshTokenRepo auth.RefreshTokenRepository,
	tokenGen auth.TokenGenerator,
	passwordHasher security.PasswordHasher,
	eventBus shared.EventBus,
	cfg *config.Config,
	log logger.Logger,
) *authcommand.RegisterHandler {
	return authcommand.NewRegisterHandler(authcommand.RegisterHandlerParams{
		UserRepository:         userRepo,
		RoleRepository:         roleRepo,
		PermissionRepository:   permissionRepo,
		RefreshTokenRepository: refreshTokenRepo,
		TokenGenerator:         tokenGen,
		PasswordHasher:         passwordHasher,
		EventBus:               eventBus,
		RefreshTokenTTL:        cfg.JWT.RefreshTokenTTL,
		Logger:                 log,
	})
}

func provideLoginHandler(
	userRepo user.Repository,
	roleRepo role.Repository,
	permissionRepo permission.Repository,
	refreshTokenRepo auth.RefreshTokenRepository,
	tokenGen auth.TokenGenerator,
	passwordHasher security.PasswordHasher,
	accountLockout security.AccountLockout,
	eventBus shared.EventBus,
	cfg *config.Config,
	log logger.Logger,
) *authcommand.LoginHandler {
	return authcommand.NewLoginHandler(authcommand.LoginHandlerParams{
		UserRepository:         userRepo,
		RoleRepository:         roleRepo,
		PermissionRepository:   permissionRepo,
		RefreshTokenRepository: refreshTokenRepo,
		TokenGenerator:         tokenGen,
		PasswordHasher:         passwordHasher,
		AccountLockout:         accountLockout,
		EventBus:               eventBus,
		RefreshTokenTTL:        cfg.JWT.RefreshTokenTTL,
		Logger:                 log,
	})
}

func provideRefreshTokenHandler(
	userRepo user.Repository,
	roleRepo role.Repository,
	permissionRepo permission.Repository,
	refreshTokenRepo auth.RefreshTokenRepository,
	tokenGen auth.TokenGenerator,
	tokenBlacklist auth.TokenBlacklist,
	eventBus shared.EventBus,
	cfg *config.Config,
	log logger.Logger,
) *authcommand.RefreshTokenHandler {
	return authcommand.NewRefreshTokenHandler(authcommand.RefreshTokenHandlerParams{
		UserRepository:         userRepo,
		RoleRepository:         roleRepo,
		PermissionRepository:   permissionRepo,
		RefreshTokenRepository: refreshTokenRepo,
		TokenGenerator:         tokenGen,
		TokenBlacklist:         tokenBlacklist,
		EventBus:               eventBus,
		RefreshTokenTTL:        cfg.JWT.RefreshTokenTTL,
		Logger:                 log,
	})
}

func provideLogoutHandler(
	refreshTokenRepo auth.RefreshTokenRepository,
	tokenBlacklist auth.TokenBlacklist,
	eventBus shared.EventBus,
	log logger.Logger,
) *authcommand.LogoutHandler {
	return authcommand.NewLogoutHandler(authcommand.LogoutHandlerParams{
		RefreshTokenRepository: refreshTokenRepo,
		TokenBlacklist:         tokenBlacklist,
		EventBus:               eventBus,
		Logger:                 log,
	})
}

func provideForgotPasswordHandler(
	userRepo user.Repository,
	tokenGen auth.TokenGenerator,
	passwordResetStore authcommand.PasswordResetTokenStore,
	eventBus shared.EventBus,
	log logger.Logger,
) *authcommand.ForgotPasswordHandler {
	return authcommand.NewForgotPasswordHandler(authcommand.ForgotPasswordHandlerParams{
		UserRepository:          userRepo,
		TokenGenerator:          tokenGen,
		PasswordResetTokenStore: passwordResetStore,
		EventBus:                eventBus,
		ResetTokenTTL:           15 * 60 * 1000000000,
		Logger:                  log,
	})
}

func provideResetPasswordHandler(
	userRepo user.Repository,
	tokenGen auth.TokenGenerator,
	passwordResetStore authcommand.PasswordResetTokenStore,
	passwordHasher security.PasswordHasher,
	eventBus shared.EventBus,
	log logger.Logger,
) *authcommand.ResetPasswordHandler {
	return authcommand.NewResetPasswordHandler(authcommand.ResetPasswordHandlerParams{
		UserRepository:          userRepo,
		TokenGenerator:          tokenGen,
		PasswordResetTokenStore: passwordResetStore,
		PasswordHasher:          passwordHasher,
		EventBus:                eventBus,
		Logger:                  log,
	})
}

func provideGetCurrentUserHandler(
	userRepo user.Repository,
	roleRepo role.Repository,
	permissionRepo permission.Repository,
	log logger.Logger,
) *authquery.GetCurrentUserHandler {
	return authquery.NewGetCurrentUserHandler(authquery.GetCurrentUserHandlerParams{
		UserRepository:       userRepo,
		RoleRepository:       roleRepo,
		PermissionRepository: permissionRepo,
		Logger:               log,
	})
}

func provideGetUserSessionsHandler(
	refreshTokenRepo auth.RefreshTokenRepository,
	log logger.Logger,
) *authquery.GetUserSessionsHandler {
	return authquery.NewGetUserSessionsHandler(authquery.GetUserSessionsHandlerParams{
		RefreshTokenRepository: refreshTokenRepo,
		Logger:                 log,
	})
}

func provideHealthHandler(database *postgres.DB, redisClient *redis.Client) *handler.HealthHandler {
	return handler.NewHealthHandler(database, redisClient)
}

func provideMetricsHandler() *handler.MetricsHandler {
	return handler.NewMetricsHandler()
}

func provideDocsHandler(log logger.Logger) *handler.DocsHandler {
	return handler.NewDocsHandler(log)
}

func provideRevokeSessionHandler(
	refreshTokenRepo auth.RefreshTokenRepository,
	eventBus shared.EventBus,
	log logger.Logger,
) *authcommand.RevokeSessionHandler {
	return authcommand.NewRevokeSessionHandler(authcommand.RevokeSessionHandlerParams{
		RefreshTokenRepository: refreshTokenRepo,
		EventBus:               eventBus,
		Logger:                 log,
	})
}

func provideAuthHandler(
	registerHandler *authcommand.RegisterHandler,
	loginHandler *authcommand.LoginHandler,
	refreshTokenHandler *authcommand.RefreshTokenHandler,
	logoutHandler *authcommand.LogoutHandler,
	forgotPasswordHandler *authcommand.ForgotPasswordHandler,
	resetPasswordHandler *authcommand.ResetPasswordHandler,
	revokeSessionHandler *authcommand.RevokeSessionHandler,
	getCurrentUserHandler *authquery.GetCurrentUserHandler,
	getUserSessionsHandler *authquery.GetUserSessionsHandler,
	val *validator.Validator,
	log logger.Logger,
) *handler.AuthHandler {
	return handler.NewAuthHandler(handler.AuthHandlerParams{
		RegisterHandler:        registerHandler,
		LoginHandler:           loginHandler,
		RefreshTokenHandler:    refreshTokenHandler,
		LogoutHandler:          logoutHandler,
		ForgotPasswordHandler:  forgotPasswordHandler,
		ResetPasswordHandler:   resetPasswordHandler,
		RevokeSessionHandler:   revokeSessionHandler,
		GetCurrentUserHandler:  getCurrentUserHandler,
		GetUserSessionsHandler: getUserSessionsHandler,
		Validator:              val,
		Logger:                 log,
	})
}

func provideRouter(
	userHandler *handler.UserHandler,
	authHandler *handler.AuthHandler,
	permissionHandler *handler.PermissionHandler,
	roleHandler *handler.RoleHandler,
	healthHandler *handler.HealthHandler,
	metricsHandler *handler.MetricsHandler,
	docsHandler *handler.DocsHandler,
	authMiddleware *middleware.AuthMiddleware,
	log logger.Logger,
	cfg *config.Config,
) http.Handler {
	return router.NewRouter(router.RouterDependencies{
		UserHandler:       userHandler,
		AuthHandler:       authHandler,
		PermissionHandler: permissionHandler,
		RoleHandler:       roleHandler,
		HealthHandler:     healthHandler,
		MetricsHandler:    metricsHandler,
		DocsHandler:       docsHandler,
		AuthMiddleware:    authMiddleware,
		Logger:            log,
		Config:            cfg,
	})
}

var InfrastructureSet = wire.NewSet(
	provideLogger,
	providePasswordHasher,
	provideValidator,
	provideAuditLogger,
	provideAuthAuditHandler,
	provideEventBusWithAudit,
	provideTokenGenerator,
	provideTokenBlacklist,
	providePasswordResetTokenStore,
	provideAccountLockout,
	provideAuthMiddleware,
	wire.Bind(new(shared.EventBus), new(*memory.InMemoryEventBus)),
)

var RepositorySet = wire.NewSet(
	provideUserRepository,
	providePermissionRepository,
	provideRoleRepository,
	provideRefreshTokenRepository,
	wire.Bind(new(user.Repository), new(*repository.UserRepository)),
	wire.Bind(new(permission.Repository), new(*repository.PermissionRepository)),
	wire.Bind(new(role.Repository), new(*repository.RoleRepository)),
	wire.Bind(new(auth.RefreshTokenRepository), new(*repository.RefreshTokenRepository)),
)

var UserCommandHandlerSet = wire.NewSet(
	usercommand.NewCreateUserHandler,
	usercommand.NewUpdateUserHandler,
	usercommand.NewDeleteUserHandler,
	usercommand.NewChangePasswordHandler,
	usercommand.NewActivateUserHandler,
	usercommand.NewDeactivateUserHandler,
	usercommand.NewBanUserHandler,
	usercommand.NewAssignRoleToUserHandler,
	usercommand.NewRevokeRoleFromUserHandler,
	usercommand.NewSetUserRolesHandler,
)

var AuthCommandHandlerSet = wire.NewSet(
	provideRegisterHandler,
	provideLoginHandler,
	provideRefreshTokenHandler,
	provideLogoutHandler,
	provideForgotPasswordHandler,
	provideResetPasswordHandler,
	provideRevokeSessionHandler,
)

var UserQueryHandlerSet = wire.NewSet(
	userquery.NewGetUserHandler,
	userquery.NewListUsersHandler,
	userquery.NewGetUserRolesHandler,
	userquery.NewGetUserPermissionsHandler,
)

var AuthQueryHandlerSet = wire.NewSet(
	provideGetCurrentUserHandler,
	provideGetUserSessionsHandler,
)

var PermissionCommandHandlerSet = wire.NewSet(
	permissioncommand.NewCreatePermissionHandler,
	permissioncommand.NewUpdatePermissionHandler,
	permissioncommand.NewDeletePermissionHandler,
)

var PermissionQueryHandlerSet = wire.NewSet(
	permissionquery.NewListPermissionsHandler,
	permissionquery.NewGetPermissionHandler,
	permissionquery.NewGetPermissionsForRoleHandler,
)

var RoleCommandHandlerSet = wire.NewSet(
	rolecommand.NewCreateRoleHandler,
	rolecommand.NewUpdateRoleHandler,
	rolecommand.NewDeleteRoleHandler,
	rolecommand.NewAssignPermissionToRoleHandler,
	rolecommand.NewRemovePermissionFromRoleHandler,
	rolecommand.NewSetRolePermissionsHandler,
)

var RoleQueryHandlerSet = wire.NewSet(
	rolequery.NewListRolesHandler,
	rolequery.NewGetRoleHandler,
	rolequery.NewGetUsersWithRoleHandler,
)

var HandlerSet = wire.NewSet(
	wire.Struct(new(handler.UserHandlerParams), "*"),
	handler.NewUserHandler,
	wire.Struct(new(handler.PermissionHandlerParams), "*"),
	handler.NewPermissionHandler,
	wire.Struct(new(handler.RoleHandlerParams), "*"),
	handler.NewRoleHandler,
	provideAuthHandler,
	provideHealthHandler,
	provideMetricsHandler,
	provideDocsHandler,
)

var RouterSet = wire.NewSet(
	provideRouter,
)

func InitializeApplication(
	cfg *config.Config,
	database *postgres.DB,
	redisClient *redis.Client,
) (*Application, error) {
	wire.Build(
		InfrastructureSet,
		RepositorySet,
		UserCommandHandlerSet,
		AuthCommandHandlerSet,
		PermissionCommandHandlerSet,
		RoleCommandHandlerSet,
		UserQueryHandlerSet,
		AuthQueryHandlerSet,
		PermissionQueryHandlerSet,
		RoleQueryHandlerSet,
		HandlerSet,
		RouterSet,
		NewApplication,
	)
	return nil, nil
}
