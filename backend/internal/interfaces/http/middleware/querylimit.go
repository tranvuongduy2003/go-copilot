package middleware

import (
	"fmt"
	"net/http"
)

const (
	DefaultMaxQueryParameterLength = 256
	DefaultMaxQueryParameterCount  = 50
	DefaultMaxTotalQueryLength     = 2048
)

type QueryLimitConfig struct {
	MaxParameterLength int
	MaxParameterCount  int
	MaxTotalLength     int
	ExcludedPaths      []string
	ErrorHandler       func(http.ResponseWriter, *http.Request, error)
}

func DefaultQueryLimitConfig() QueryLimitConfig {
	return QueryLimitConfig{
		MaxParameterLength: DefaultMaxQueryParameterLength,
		MaxParameterCount:  DefaultMaxQueryParameterCount,
		MaxTotalLength:     DefaultMaxTotalQueryLength,
		ExcludedPaths:      []string{},
		ErrorHandler:       defaultQueryLimitErrorHandler,
	}
}

func defaultQueryLimitErrorHandler(writer http.ResponseWriter, request *http.Request, err error) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusBadRequest)
	fmt.Fprintf(writer, `{"error":{"code":"QUERY_PARAMETER_LIMIT_EXCEEDED","message":"%s"}}`, err.Error())
}

func QueryLimit(config QueryLimitConfig) func(http.Handler) http.Handler {
	if config.MaxParameterLength <= 0 {
		config.MaxParameterLength = DefaultMaxQueryParameterLength
	}

	if config.MaxParameterCount <= 0 {
		config.MaxParameterCount = DefaultMaxQueryParameterCount
	}

	if config.MaxTotalLength <= 0 {
		config.MaxTotalLength = DefaultMaxTotalQueryLength
	}

	if config.ErrorHandler == nil {
		config.ErrorHandler = defaultQueryLimitErrorHandler
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			if isPathExcluded(request.URL.Path, config.ExcludedPaths) {
				next.ServeHTTP(writer, request)
				return
			}

			rawQuery := request.URL.RawQuery
			if len(rawQuery) > config.MaxTotalLength {
				err := fmt.Errorf("query string too long: %d characters exceeds limit of %d",
					len(rawQuery), config.MaxTotalLength)
				config.ErrorHandler(writer, request, err)
				return
			}

			queryParams := request.URL.Query()
			parameterCount := 0
			for key, values := range queryParams {
				parameterCount += len(values)

				if len(key) > config.MaxParameterLength {
					err := fmt.Errorf("query parameter name '%s' too long: %d characters exceeds limit of %d",
						truncateString(key, 20), len(key), config.MaxParameterLength)
					config.ErrorHandler(writer, request, err)
					return
				}

				for _, value := range values {
					if len(value) > config.MaxParameterLength {
						err := fmt.Errorf("query parameter '%s' value too long: %d characters exceeds limit of %d",
							key, len(value), config.MaxParameterLength)
						config.ErrorHandler(writer, request, err)
						return
					}
				}
			}

			if parameterCount > config.MaxParameterCount {
				err := fmt.Errorf("too many query parameters: %d exceeds limit of %d",
					parameterCount, config.MaxParameterCount)
				config.ErrorHandler(writer, request, err)
				return
			}

			next.ServeHTTP(writer, request)
		})
	}
}

func QueryLimitDefault(next http.Handler) http.Handler {
	config := DefaultQueryLimitConfig()
	return QueryLimit(config)(next)
}

func isPathExcluded(requestPath string, excludedPaths []string) bool {
	for _, path := range excludedPaths {
		if matchPath(requestPath, path) {
			return true
		}
	}
	return false
}

func truncateString(stringValue string, maxLength int) string {
	if len(stringValue) <= maxLength {
		return stringValue
	}
	return stringValue[:maxLength] + "..."
}
