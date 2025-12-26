package main

import (
	"os"

	"github.com/tranvuongduy2003/go-copilot/pkg/config"
	"github.com/tranvuongduy2003/go-copilot/pkg/logger"
)

var (
	Version   = "dev"
	BuildTime = "unknown"
	GoVersion = "unknown"
)

func main() {
	cfg, err := config.LoadFromEnv()
	if err != nil {
		logger.Fatal("failed to load config", logger.Err(err))
	}

	if err := logger.Init(cfg.Log, cfg.IsProduction()); err != nil {
		logger.Fatal("failed to initialize logger", logger.Err(err))
	}
	defer logger.Sync()

	logger.Info("starting server",
		logger.String("app", cfg.App.Name),
		logger.String("version", Version),
		logger.String("build_time", BuildTime),
		logger.String("go_version", GoVersion),
		logger.String("env", cfg.App.Env),
	)

	logger.Info("TODO: implement server startup")
	os.Exit(0)
}
