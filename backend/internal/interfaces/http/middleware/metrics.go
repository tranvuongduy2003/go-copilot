package middleware

import (
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/tranvuongduy2003/go-copilot/pkg/metrics"
)

type metricsResponseWriter struct {
	http.ResponseWriter
	statusCode    int
	bytesWritten  int64
	headerWritten bool
}

func newMetricsResponseWriter(writer http.ResponseWriter) *metricsResponseWriter {
	return &metricsResponseWriter{
		ResponseWriter: writer,
		statusCode:     http.StatusOK,
	}
}

func (writer *metricsResponseWriter) WriteHeader(statusCode int) {
	if !writer.headerWritten {
		writer.statusCode = statusCode
		writer.headerWritten = true
	}
	writer.ResponseWriter.WriteHeader(statusCode)
}

func (writer *metricsResponseWriter) Write(data []byte) (int, error) {
	if !writer.headerWritten {
		writer.headerWritten = true
	}
	bytesWritten, err := writer.ResponseWriter.Write(data)
	writer.bytesWritten += int64(bytesWritten)
	return bytesWritten, err
}

func Metrics(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		start := time.Now()

		metricsWriter := newMetricsResponseWriter(writer)

		next.ServeHTTP(metricsWriter, request)

		duration := time.Since(start).Seconds()

		path := getRoutePath(request)

		metrics.RecordHTTPRequest(
			request.Method,
			path,
			strconv.Itoa(metricsWriter.statusCode),
			duration,
			request.ContentLength,
			metricsWriter.bytesWritten,
		)
	})
}

func getRoutePath(request *http.Request) string {
	routeContext := chi.RouteContext(request.Context())
	if routeContext != nil && routeContext.RoutePattern() != "" {
		return routeContext.RoutePattern()
	}
	return request.URL.Path
}
