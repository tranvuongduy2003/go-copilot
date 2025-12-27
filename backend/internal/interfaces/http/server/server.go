package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/tranvuongduy2003/go-copilot/pkg/config"
	"github.com/tranvuongduy2003/go-copilot/pkg/logger"
)

type Server struct {
	httpServer *http.Server
	logger     logger.Logger
}

func New(handler http.Handler, serverConfig config.ServerConfig, logger logger.Logger) *Server {
	return &Server{
		httpServer: &http.Server{
			Addr:         serverConfig.Address(),
			Handler:      handler,
			ReadTimeout:  serverConfig.ReadTimeout,
			WriteTimeout: serverConfig.WriteTimeout,
			IdleTimeout:  serverConfig.IdleTimeout,
		},
		logger: logger,
	}
}

func (server *Server) Start() error {
	server.logger.Info("starting HTTP server",
		logger.String("address", server.httpServer.Addr),
	)

	if err := server.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("HTTP server error: %w", err)
	}

	return nil
}

func (server *Server) Shutdown(ctx context.Context) error {
	server.logger.Info("shutting down HTTP server")

	shutdownContext, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	if err := server.httpServer.Shutdown(shutdownContext); err != nil {
		return fmt.Errorf("HTTP server shutdown error: %w", err)
	}

	server.logger.Info("HTTP server shutdown complete")
	return nil
}

func (server *Server) Address() string {
	return server.httpServer.Addr
}
