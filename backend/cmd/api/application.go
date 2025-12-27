package main

import (
	"context"
	"net/http"
	"time"

	"github.com/tranvuongduy2003/go-copilot/internal/infrastructure/cache/redis"
	"github.com/tranvuongduy2003/go-copilot/internal/infrastructure/messaging/memory"
	"github.com/tranvuongduy2003/go-copilot/internal/infrastructure/persistence/postgres"
	"github.com/tranvuongduy2003/go-copilot/internal/interfaces/http/server"
	"github.com/tranvuongduy2003/go-copilot/pkg/config"
	"github.com/tranvuongduy2003/go-copilot/pkg/logger"
)

type Application struct {
	Config      *config.Config
	Database    *postgres.DB
	RedisClient *redis.Client
	EventBus    *memory.InMemoryEventBus
	Router      http.Handler
	Server      *server.Server
}

func NewApplication(
	cfg *config.Config,
	database *postgres.DB,
	redisClient *redis.Client,
	eventBus *memory.InMemoryEventBus,
	routerHandler http.Handler,
) *Application {
	httpServer := server.New(routerHandler, cfg.Server, logger.L())

	return &Application{
		Config:      cfg,
		Database:    database,
		RedisClient: redisClient,
		EventBus:    eventBus,
		Router:      routerHandler,
		Server:      httpServer,
	}
}

func (app *Application) Start() error {
	return app.Server.Start()
}

func (app *Application) Shutdown(ctx context.Context) error {
	return app.Server.Shutdown(ctx)
}

func (app *Application) Close() {
	if app.EventBus != nil {
		shutdownContext, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		app.EventBus.Shutdown(shutdownContext)
	}

	if app.RedisClient != nil {
		app.RedisClient.Close()
	}

	if app.Database != nil {
		app.Database.Close()
	}
}

func (app *Application) Address() string {
	return app.Server.Address()
}
