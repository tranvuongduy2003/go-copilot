package handler

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type MetricsHandler struct{}

func NewMetricsHandler() *MetricsHandler {
	return &MetricsHandler{}
}

func (handler *MetricsHandler) Metrics() http.Handler {
	return promhttp.Handler()
}
