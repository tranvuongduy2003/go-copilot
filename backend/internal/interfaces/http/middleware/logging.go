package middleware

import (
	"net/http"
	"time"

	"github.com/tranvuongduy2003/go-copilot/pkg/logger"
)

type responseWriter struct {
	http.ResponseWriter
	statusCode    int
	bytesWritten  int
	headerWritten bool
}

func newResponseWriter(writer http.ResponseWriter) *responseWriter {
	return &responseWriter{
		ResponseWriter: writer,
		statusCode:     http.StatusOK,
	}
}

func (rw *responseWriter) WriteHeader(statusCode int) {
	if !rw.headerWritten {
		rw.statusCode = statusCode
		rw.headerWritten = true
		rw.ResponseWriter.WriteHeader(statusCode)
	}
}

func (rw *responseWriter) Write(bytes []byte) (int, error) {
	if !rw.headerWritten {
		rw.WriteHeader(http.StatusOK)
	}
	bytesWritten, err := rw.ResponseWriter.Write(bytes)
	rw.bytesWritten += bytesWritten
	return bytesWritten, err
}

func (rw *responseWriter) Status() int {
	return rw.statusCode
}

func (rw *responseWriter) BytesWritten() int {
	return rw.bytesWritten
}

func Logging(log logger.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			startTime := time.Now()

			wrappedWriter := newResponseWriter(writer)

			requestID := GetRequestID(request.Context())
			method := request.Method
			path := request.URL.Path
			query := request.URL.RawQuery
			remoteAddr := request.RemoteAddr
			userAgent := request.UserAgent()

			if !isHealthCheckPath(path) {
				log.Info("request started",
					logger.String("request_id", requestID),
					logger.String("method", method),
					logger.String("path", path),
					logger.String("query", query),
					logger.String("remote_addr", remoteAddr),
					logger.String("user_agent", userAgent),
				)
			}

			next.ServeHTTP(wrappedWriter, request)

			duration := time.Since(startTime)
			statusCode := wrappedWriter.Status()
			bytesWritten := wrappedWriter.BytesWritten()

			if !isHealthCheckPath(path) {
				logLevel := determineLogLevel(statusCode)
				logLevel(log, "request completed",
					logger.String("request_id", requestID),
					logger.String("method", method),
					logger.String("path", path),
					logger.Int("status", statusCode),
					logger.Int("bytes", bytesWritten),
					logger.Duration("duration", duration),
				)
			}
		})
	}
}

func isHealthCheckPath(path string) bool {
	return path == "/health/live" || path == "/health/ready" || path == "/health"
}

type logFunc func(log logger.Logger, message string, fields ...logger.Field)

func determineLogLevel(statusCode int) logFunc {
	if statusCode >= 500 {
		return func(log logger.Logger, message string, fields ...logger.Field) {
			log.Error(message, fields...)
		}
	}
	if statusCode >= 400 {
		return func(log logger.Logger, message string, fields ...logger.Field) {
			log.Warn(message, fields...)
		}
	}
	return func(log logger.Logger, message string, fields ...logger.Field) {
		log.Info(message, fields...)
	}
}
