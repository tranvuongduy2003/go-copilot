package middleware

import (
	"net/http"
	"runtime/debug"

	"github.com/tranvuongduy2003/go-copilot/internal/interfaces/http/response"
	"github.com/tranvuongduy2003/go-copilot/pkg/logger"
)

func Recovery(log logger.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			defer func() {
				if recovered := recover(); recovered != nil {
					requestID := GetRequestID(request.Context())
					stackTrace := string(debug.Stack())

					log.Error("panic recovered",
						logger.String("request_id", requestID),
						logger.String("method", request.Method),
						logger.String("path", request.URL.Path),
						logger.Any("panic", recovered),
						logger.String("stack_trace", stackTrace),
					)

					response.InternalError(writer, request)
				}
			}()

			next.ServeHTTP(writer, request)
		})
	}
}
