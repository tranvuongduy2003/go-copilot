package middleware

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

type requestIDContextKey struct{}

const RequestIDHeader = "X-Request-ID"

func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		requestID := request.Header.Get(RequestIDHeader)
		if requestID == "" {
			requestID = uuid.New().String()
		}

		contextWithRequestID := context.WithValue(request.Context(), requestIDContextKey{}, requestID)
		request = request.WithContext(contextWithRequestID)

		writer.Header().Set(RequestIDHeader, requestID)

		next.ServeHTTP(writer, request)
	})
}

func GetRequestID(context context.Context) string {
	if context == nil {
		return ""
	}
	if requestID, ok := context.Value(requestIDContextKey{}).(string); ok {
		return requestID
	}
	return ""
}
