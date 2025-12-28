package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/tranvuongduy2003/go-copilot/internal/infrastructure/cache/redis"
	"github.com/tranvuongduy2003/go-copilot/internal/infrastructure/persistence/migrations"
	"github.com/tranvuongduy2003/go-copilot/internal/infrastructure/persistence/postgres"
	"github.com/tranvuongduy2003/go-copilot/pkg/config"
	"github.com/tranvuongduy2003/go-copilot/pkg/logger"
)

var (
	Version   = "dev"
	BuildTime = "unknown"
	GoVersion = "unknown"
)

func main() {
	exitCode := run()
	os.Exit(exitCode)
}

func run() int {
	configuration, err := config.LoadFromEnv()
	if err != nil {
		logger.Fatal("failed to load configuration", logger.Err(err))
		return 1
	}

	if err := logger.Init(configuration.Log, configuration.IsProduction()); err != nil {
		logger.Fatal("failed to initialize logger", logger.Err(err))
		return 1
	}
	defer logger.Sync()

	logger.Info("starting application",
		logger.String("app", configuration.App.Name),
		logger.String("version", Version),
		logger.String("build_time", BuildTime),
		logger.String("go_version", GoVersion),
		logger.String("environment", configuration.App.Env),
	)

	applicationContext, cancelApplication := context.WithCancel(context.Background())
	defer cancelApplication()

	database := postgres.NewDB(&configuration.Database, logger.L())
	if err := database.Connect(applicationContext); err != nil {
		logger.Fatal("failed to connect to database", logger.Err(err))
		return 1
	}
	logger.Info("database connection established")

	if configuration.Database.AutoMigrate {
		logger.Info("running database migrations",
			logger.String("path", configuration.Database.MigrationsPath),
		)
		migrationRunner := migrations.NewRunner(
			configuration.Database.MigrationsPath,
			configuration.Database.DSN(),
		)
		if err := migrationRunner.Up(); err != nil {
			logger.Error("failed to run database migrations", logger.Err(err))
			return 1
		}
		version, dirty, err := migrationRunner.Version()
		if err != nil {
			logger.Warn("failed to get migration version", logger.Err(err))
		} else {
			logger.Info("database migrations completed",
				logger.Int64("version", int64(version)),
				logger.Bool("dirty", dirty),
			)
		}
	}

	var redisClient *redis.Client
	if configuration.Redis.Host != "" {
		redisClient = redis.NewClient(configuration.Redis, logger.L())
		if err := redisClient.Connect(applicationContext); err != nil {
			logger.Warn("failed to connect to redis, continuing without redis",
				logger.Err(err),
			)
			redisClient = nil
		} else {
			logger.Info("redis connection established")
		}
	}

	application, err := InitializeApplication(configuration, database, redisClient)
	if err != nil {
		logger.Fatal("failed to initialize application", logger.Err(err))
		return 1
	}
	defer application.Close()

	serverErrorChannel := make(chan error, 1)
	go func() {
		if err := application.Start(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErrorChannel <- err
		}
		close(serverErrorChannel)
	}()

	logger.Info("server started",
		logger.String("address", application.Address()),
	)

	shutdownSignal := make(chan os.Signal, 1)
	signal.Notify(shutdownSignal, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-serverErrorChannel:
		logger.Error("server error", logger.Err(err))
		return 1
	case sig := <-shutdownSignal:
		logger.Info("shutdown signal received", logger.String("signal", sig.String()))
	}

	logger.Info("initiating graceful shutdown")

	shutdownTimeout := 30 * time.Second
	shutdownContext, cancelShutdown := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancelShutdown()

	if err := application.Shutdown(shutdownContext); err != nil {
		logger.Error("server shutdown error", logger.Err(err))
		return 1
	}

	logger.Info("application shutdown complete")
	return 0
}
