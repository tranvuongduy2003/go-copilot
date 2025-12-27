package router

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/tranvuongduy2003/go-copilot/internal/interfaces/http/handler"
	"github.com/tranvuongduy2003/go-copilot/internal/interfaces/http/middleware"
	"github.com/tranvuongduy2003/go-copilot/pkg/config"
	"github.com/tranvuongduy2003/go-copilot/pkg/logger"
)

type RouterDependencies struct {
	UserHandler   *handler.UserHandler
	HealthHandler *handler.HealthHandler
	Logger        logger.Logger
	Config        *config.Config
}

func NewRouter(dependencies RouterDependencies) http.Handler {
	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.Logging(dependencies.Logger))
	router.Use(middleware.Recovery(dependencies.Logger))
	router.Use(middleware.CORS(dependencies.Config.CORS))
	router.Use(middleware.Timeout(30 * time.Second))

	router.Route("/health", func(healthRouter chi.Router) {
		healthRouter.Get("/", dependencies.HealthHandler.Readiness)
		healthRouter.Get("/live", dependencies.HealthHandler.Liveness)
		healthRouter.Get("/ready", dependencies.HealthHandler.Readiness)
	})

	router.Route("/api/v1", func(apiRouter chi.Router) {
		apiRouter.Route("/users", func(userRouter chi.Router) {
			userRouter.Post("/", dependencies.UserHandler.Create)
			userRouter.Get("/", dependencies.UserHandler.List)

			userRouter.Route("/{id}", func(userIDRouter chi.Router) {
				userIDRouter.Get("/", dependencies.UserHandler.Get)
				userIDRouter.Put("/", dependencies.UserHandler.Update)
				userIDRouter.Delete("/", dependencies.UserHandler.Delete)

				userIDRouter.Post("/password", dependencies.UserHandler.ChangePassword)
				userIDRouter.Post("/activate", dependencies.UserHandler.Activate)
				userIDRouter.Post("/deactivate", dependencies.UserHandler.Deactivate)
				userIDRouter.Post("/ban", dependencies.UserHandler.Ban)
			})
		})
	})

	return router
}
