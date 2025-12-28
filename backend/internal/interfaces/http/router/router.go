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
	UserHandler       *handler.UserHandler
	AuthHandler       *handler.AuthHandler
	PermissionHandler *handler.PermissionHandler
	RoleHandler       *handler.RoleHandler
	HealthHandler     *handler.HealthHandler
	MetricsHandler    *handler.MetricsHandler
	DocsHandler       *handler.DocsHandler
	AuthMiddleware    *middleware.AuthMiddleware
	Logger            logger.Logger
	Config            *config.Config
}

func NewRouter(dependencies RouterDependencies) http.Handler {
	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.Logging(dependencies.Logger))
	router.Use(middleware.Recovery(dependencies.Logger))
	router.Use(middleware.CORS(dependencies.Config.CORS))
	router.Use(middleware.BodyLimitDefault)
	router.Use(middleware.QueryLimitDefault)
	router.Use(middleware.Timeout(30 * time.Second))
	router.Use(middleware.SecureHeaders)
	router.Use(middleware.Metrics)
	router.Use(middleware.GzipCompress)

	router.Route("/health", func(healthRouter chi.Router) {
		healthRouter.Get("/", dependencies.HealthHandler.Readiness)
		healthRouter.Get("/live", dependencies.HealthHandler.Liveness)
		healthRouter.Get("/ready", dependencies.HealthHandler.Readiness)
	})

	if dependencies.MetricsHandler != nil {
		router.Handle("/metrics", dependencies.MetricsHandler.Metrics())
	}

	if dependencies.DocsHandler != nil {
		router.Route("/docs", func(docsRouter chi.Router) {
			docsRouter.Get("/", dependencies.DocsHandler.SwaggerUI)
			docsRouter.Get("/openapi.yaml", dependencies.DocsHandler.OpenAPISpec)
		})
	}

	loginRateLimiter := middleware.NewRateLimiter(middleware.LoginRateLimiterConfig())
	registerRateLimiter := middleware.NewRateLimiter(middleware.RegisterRateLimiterConfig())
	passwordResetRateLimiter := middleware.NewRateLimiter(middleware.PasswordResetRateLimiterConfig())
	tokenRefreshRateLimiter := middleware.NewRateLimiter(middleware.TokenRefreshRateLimiterConfig())

	router.Route("/api/v1", func(apiRouter chi.Router) {
		apiRouter.Route("/auth", func(authRouter chi.Router) {
			authRouter.With(middleware.RateLimit(registerRateLimiter)).Post("/register", dependencies.AuthHandler.Register)
			authRouter.With(middleware.RateLimit(loginRateLimiter)).Post("/login", dependencies.AuthHandler.Login)
			authRouter.With(middleware.RateLimit(tokenRefreshRateLimiter)).Post("/refresh", dependencies.AuthHandler.RefreshToken)
			authRouter.With(middleware.RateLimit(passwordResetRateLimiter)).Post("/forgot-password", dependencies.AuthHandler.ForgotPassword)
			authRouter.With(middleware.RateLimit(passwordResetRateLimiter)).Post("/reset-password", dependencies.AuthHandler.ResetPassword)

			authRouter.Group(func(protectedAuthRouter chi.Router) {
				protectedAuthRouter.Use(dependencies.AuthMiddleware.RequireAuth)
				protectedAuthRouter.Post("/logout", dependencies.AuthHandler.Logout)
				protectedAuthRouter.Post("/logout-all", dependencies.AuthHandler.LogoutAll)
				protectedAuthRouter.Get("/me", dependencies.AuthHandler.GetCurrentUser)
				protectedAuthRouter.Get("/sessions", dependencies.AuthHandler.GetSessions)
				protectedAuthRouter.Delete("/sessions/{id}", dependencies.AuthHandler.RevokeSession)
			})
		})

		apiRouter.Route("/users", func(userRouter chi.Router) {
			userRouter.Use(dependencies.AuthMiddleware.RequireAuth)

			userRouter.With(middleware.RequirePermission("users:create")).Post("/", dependencies.UserHandler.Create)
			userRouter.With(middleware.RequirePermission("users:list")).Get("/", dependencies.UserHandler.List)

			userRouter.Route("/{id}", func(userIDRouter chi.Router) {
				userIDRouter.With(middleware.RequirePermission("users:read")).Get("/", dependencies.UserHandler.Get)
				userIDRouter.With(middleware.RequirePermission("users:update")).Put("/", dependencies.UserHandler.Update)
				userIDRouter.With(middleware.RequirePermission("users:delete")).Delete("/", dependencies.UserHandler.Delete)

				userIDRouter.With(middleware.ResourceOwner("id")).Post("/password", dependencies.UserHandler.ChangePassword)
				userIDRouter.With(middleware.RequirePermission("users:manage")).Post("/activate", dependencies.UserHandler.Activate)
				userIDRouter.With(middleware.RequirePermission("users:manage")).Post("/deactivate", dependencies.UserHandler.Deactivate)
				userIDRouter.With(middleware.RequirePermission("users:manage")).Post("/ban", dependencies.UserHandler.Ban)

				userIDRouter.With(middleware.RequireAnyPermission("users:read", "roles:assign")).Get("/roles", dependencies.UserHandler.GetRoles)
				userIDRouter.With(middleware.RequirePermission("roles:assign")).Put("/roles", dependencies.UserHandler.SetRoles)
				userIDRouter.With(middleware.RequirePermission("roles:assign")).Post("/roles/{roleId}", dependencies.UserHandler.AssignRole)
				userIDRouter.With(middleware.RequirePermission("roles:assign")).Delete("/roles/{roleId}", dependencies.UserHandler.RevokeRole)
				userIDRouter.With(middleware.RequireAnyPermission("users:read", "permissions:read")).Get("/permissions", dependencies.UserHandler.GetPermissions)
			})
		})

		apiRouter.Route("/permissions", func(permissionRouter chi.Router) {
			permissionRouter.Use(dependencies.AuthMiddleware.RequireAuth)

			permissionRouter.With(middleware.RequirePermission("permissions:list")).Get("/", dependencies.PermissionHandler.List)
			permissionRouter.With(middleware.RequirePermission("permissions:create")).Post("/", dependencies.PermissionHandler.Create)

			permissionRouter.Route("/{id}", func(permissionIDRouter chi.Router) {
				permissionIDRouter.With(middleware.RequirePermission("permissions:read")).Get("/", dependencies.PermissionHandler.Get)
				permissionIDRouter.With(middleware.RequirePermission("permissions:update")).Put("/", dependencies.PermissionHandler.Update)
				permissionIDRouter.With(middleware.RequirePermission("permissions:delete")).Delete("/", dependencies.PermissionHandler.Delete)
			})
		})

		apiRouter.Route("/roles", func(roleRouter chi.Router) {
			roleRouter.Use(dependencies.AuthMiddleware.RequireAuth)

			roleRouter.With(middleware.RequirePermission("roles:list")).Get("/", dependencies.RoleHandler.List)
			roleRouter.With(middleware.RequirePermission("roles:create")).Post("/", dependencies.RoleHandler.Create)

			roleRouter.Route("/{id}", func(roleIDRouter chi.Router) {
				roleIDRouter.With(middleware.RequirePermission("roles:read")).Get("/", dependencies.RoleHandler.Get)
				roleIDRouter.With(middleware.RequirePermission("roles:update")).Put("/", dependencies.RoleHandler.Update)
				roleIDRouter.With(middleware.RequirePermission("roles:delete")).Delete("/", dependencies.RoleHandler.Delete)

				roleIDRouter.With(middleware.RequirePermission("roles:update")).Put("/permissions", dependencies.RoleHandler.SetPermissions)
				roleIDRouter.With(middleware.RequirePermission("roles:update")).Post("/permissions/{permissionId}", dependencies.RoleHandler.AddPermission)
				roleIDRouter.With(middleware.RequirePermission("roles:update")).Delete("/permissions/{permissionId}", dependencies.RoleHandler.RemovePermission)
				roleIDRouter.With(middleware.RequirePermission("roles:read")).Get("/users", dependencies.RoleHandler.GetUsersWithRole)
			})
		})
	})

	return router
}
