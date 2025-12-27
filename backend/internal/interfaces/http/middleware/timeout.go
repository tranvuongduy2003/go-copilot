package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/tranvuongduy2003/go-copilot/internal/interfaces/http/response"
)

func Timeout(timeout time.Duration) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			contextWithTimeout, cancel := context.WithTimeout(request.Context(), timeout)
			defer cancel()

			request = request.WithContext(contextWithTimeout)

			done := make(chan struct{})
			panicChannel := make(chan interface{}, 1)

			go func() {
				defer func() {
					if recovered := recover(); recovered != nil {
						panicChannel <- recovered
					}
				}()
				next.ServeHTTP(writer, request)
				close(done)
			}()

			select {
			case <-done:
				return
			case recovered := <-panicChannel:
				panic(recovered)
			case <-contextWithTimeout.Done():
				if contextWithTimeout.Err() == context.DeadlineExceeded {
					response.JSON(writer, http.StatusGatewayTimeout, response.ErrorResponse{
						Error: response.ErrorDetail{
							Code:    "GATEWAY_TIMEOUT",
							Message: "request timed out",
						},
						TraceID: GetRequestID(request.Context()),
					})
				}
			}
		})
	}
}
