package middleware

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/tranvuongduy2003/go-copilot/pkg/config"
)

func CORS(corsConfig config.CORSConfig) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			origin := request.Header.Get("Origin")

			if isOriginAllowed(origin, corsConfig.AllowedOrigins) {
				writer.Header().Set("Access-Control-Allow-Origin", origin)
			}

			writer.Header().Set("Access-Control-Allow-Methods", strings.Join(corsConfig.AllowedMethods, ", "))
			writer.Header().Set("Access-Control-Allow-Headers", strings.Join(corsConfig.AllowedHeaders, ", "))
			writer.Header().Set("Access-Control-Allow-Credentials", "true")
			writer.Header().Set("Access-Control-Max-Age", strconv.Itoa(corsConfig.MaxAge))

			if len(corsConfig.ExposedHeaders) > 0 {
				writer.Header().Set("Access-Control-Expose-Headers", strings.Join(corsConfig.ExposedHeaders, ", "))
			}

			if request.Method == http.MethodOptions {
				writer.WriteHeader(http.StatusNoContent)
				return
			}

			next.ServeHTTP(writer, request)
		})
	}
}

func isOriginAllowed(origin string, allowedOrigins []string) bool {
	if origin == "" {
		return false
	}

	for _, allowed := range allowedOrigins {
		if allowed == "*" {
			return true
		}
		if allowed == origin {
			return true
		}
	}

	return false
}

func DefaultCORSConfig() config.CORSConfig {
	return config.CORSConfig{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{
			"Accept",
			"Authorization",
			"Content-Type",
			"X-Request-ID",
			"X-Requested-With",
		},
		ExposedHeaders: []string{
			"X-Request-ID",
			"X-RateLimit-Limit",
			"X-RateLimit-Remaining",
			"X-RateLimit-Reset",
			"X-Total-Count",
		},
		MaxAge: 86400,
	}
}
