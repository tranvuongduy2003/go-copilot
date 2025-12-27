//go:build wireinject
// +build wireinject

package main

import (
	"net/http"

	"github.com/google/wire"

	"github.com/tranvuongduy2003/go-copilot/internal/application/command"
	"github.com/tranvuongduy2003/go-copilot/internal/application/query"
	"github.com/tranvuongduy2003/go-copilot/internal/domain/shared"
	"github.com/tranvuongduy2003/go-copilot/internal/domain/user"
	"github.com/tranvuongduy2003/go-copilot/internal/infrastructure/cache/redis"
	"github.com/tranvuongduy2003/go-copilot/internal/infrastructure/messaging/memory"
	"github.com/tranvuongduy2003/go-copilot/internal/infrastructure/persistence/postgres"
	"github.com/tranvuongduy2003/go-copilot/internal/infrastructure/persistence/repository"
	"github.com/tranvuongduy2003/go-copilot/internal/interfaces/http/handler"
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

func providePasswordHasher() security.PasswordHasher {
	return security.NewDefaultPasswordHasher()
}

func provideValidator() *validator.Validator {
	return validator.New()
}

func provideEventBus(log logger.Logger) *memory.InMemoryEventBus {
	return memory.NewInMemoryEventBus(log)
}

func provideHealthHandler(database *postgres.DB, redisClient *redis.Client) *handler.HealthHandler {
	return handler.NewHealthHandler(database, redisClient)
}

func provideRouter(
	userHandler *handler.UserHandler,
	healthHandler *handler.HealthHandler,
	log logger.Logger,
	cfg *config.Config,
) http.Handler {
	return router.NewRouter(router.RouterDependencies{
		UserHandler:   userHandler,
		HealthHandler: healthHandler,
		Logger:        log,
		Config:        cfg,
	})
}

var InfrastructureSet = wire.NewSet(
	provideLogger,
	providePasswordHasher,
	provideValidator,
	provideEventBus,
	wire.Bind(new(shared.EventBus), new(*memory.InMemoryEventBus)),
)

var RepositorySet = wire.NewSet(
	provideUserRepository,
	wire.Bind(new(user.Repository), new(*repository.UserRepository)),
)

var CommandHandlerSet = wire.NewSet(
	command.NewCreateUserHandler,
	command.NewUpdateUserHandler,
	command.NewDeleteUserHandler,
	command.NewChangePasswordHandler,
	command.NewActivateUserHandler,
	command.NewDeactivateUserHandler,
	command.NewBanUserHandler,
)

var QueryHandlerSet = wire.NewSet(
	query.NewGetUserHandler,
	query.NewListUsersHandler,
)

var HandlerSet = wire.NewSet(
	handler.NewUserHandler,
	provideHealthHandler,
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
		CommandHandlerSet,
		QueryHandlerSet,
		HandlerSet,
		RouterSet,
		NewApplication,
	)
	return nil, nil
}
