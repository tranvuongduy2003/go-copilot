package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/tranvuongduy2003/go-copilot/internal/infrastructure/cache/redis"
	"github.com/tranvuongduy2003/go-copilot/internal/infrastructure/persistence/postgres"
	"github.com/tranvuongduy2003/go-copilot/internal/interfaces/http/response"
)

type Pinger interface {
	Ping(ctx context.Context) error
}

type HealthHandler struct {
	database    Pinger
	redisClient Pinger
}

func NewHealthHandler(database *postgres.DB, redisClient *redis.Client) *HealthHandler {
	var dbPinger Pinger
	var redisPinger Pinger

	if database != nil {
		dbPinger = database
	}
	if redisClient != nil {
		redisPinger = redisClient
	}

	return &HealthHandler{
		database:    dbPinger,
		redisClient: redisPinger,
	}
}

func NewHealthHandlerWithPinger(database Pinger, redisClient Pinger) *HealthHandler {
	return &HealthHandler{
		database:    database,
		redisClient: redisClient,
	}
}

type LivenessResponse struct {
	Status string `json:"status"`
}

type ReadinessResponse struct {
	Status string                   `json:"status"`
	Checks map[string]CheckResponse `json:"checks"`
}

type CheckResponse struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
}

func (handler *HealthHandler) Liveness(writer http.ResponseWriter, request *http.Request) {
	response.Success(writer, LivenessResponse{
		Status: "ok",
	})
}

func (handler *HealthHandler) Readiness(writer http.ResponseWriter, request *http.Request) {
	checks := make(map[string]CheckResponse)
	allHealthy := true

	databaseCheck := handler.checkDatabase(request.Context())
	checks["database"] = databaseCheck
	if databaseCheck.Status != "ok" && databaseCheck.Status != "unknown" {
		allHealthy = false
	}

	redisCheck := handler.checkRedis(request.Context())
	checks["redis"] = redisCheck
	if redisCheck.Status != "ok" && redisCheck.Status != "unknown" {
		allHealthy = false
	}

	status := "ok"
	if !allHealthy {
		status = "degraded"
	}

	readinessResponse := ReadinessResponse{
		Status: status,
		Checks: checks,
	}

	if allHealthy {
		response.Success(writer, readinessResponse)
	} else {
		response.JSON(writer, http.StatusServiceUnavailable, response.SuccessResponse{
			Data: readinessResponse,
		})
	}
}

func (handler *HealthHandler) checkDatabase(ctx context.Context) CheckResponse {
	if handler.database == nil {
		return CheckResponse{
			Status:  "unknown",
			Message: "database not configured",
		}
	}

	checkContext, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := handler.database.Ping(checkContext); err != nil {
		return CheckResponse{
			Status:  "unhealthy",
			Message: "database connection failed",
		}
	}

	return CheckResponse{
		Status: "ok",
	}
}

func (handler *HealthHandler) checkRedis(ctx context.Context) CheckResponse {
	if handler.redisClient == nil {
		return CheckResponse{
			Status:  "unknown",
			Message: "redis not configured",
		}
	}

	checkContext, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := handler.redisClient.Ping(checkContext); err != nil {
		return CheckResponse{
			Status:  "unhealthy",
			Message: "redis connection failed",
		}
	}

	return CheckResponse{
		Status: "ok",
	}
}
