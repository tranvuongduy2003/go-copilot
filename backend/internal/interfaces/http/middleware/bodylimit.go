package middleware

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

const (
	DefaultBodyLimit = 1 << 20  // 1 MB
	MaxBodyLimit     = 10 << 20 // 10 MB
)

type BodyLimitConfig struct {
	Limit            int64
	ExcludedPaths    []string
	ExcludedMethods  []string
	ErrorHandler     func(http.ResponseWriter, *http.Request, error)
}

func DefaultBodyLimitConfig() BodyLimitConfig {
	return BodyLimitConfig{
		Limit:           DefaultBodyLimit,
		ExcludedPaths:   []string{},
		ExcludedMethods: []string{http.MethodGet, http.MethodHead, http.MethodOptions},
		ErrorHandler:    defaultBodyLimitErrorHandler,
	}
}

func defaultBodyLimitErrorHandler(writer http.ResponseWriter, request *http.Request, err error) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusRequestEntityTooLarge)
	fmt.Fprintf(writer, `{"error":{"code":"REQUEST_TOO_LARGE","message":"%s"}}`, err.Error())
}

func BodyLimit(config BodyLimitConfig) func(http.Handler) http.Handler {
	if config.Limit <= 0 {
		config.Limit = DefaultBodyLimit
	}

	if config.Limit > MaxBodyLimit {
		config.Limit = MaxBodyLimit
	}

	if config.ErrorHandler == nil {
		config.ErrorHandler = defaultBodyLimitErrorHandler
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			if shouldExclude(request, config.ExcludedPaths, config.ExcludedMethods) {
				next.ServeHTTP(writer, request)
				return
			}

			if request.ContentLength > config.Limit {
				err := fmt.Errorf("request body too large: %d bytes exceeds limit of %d bytes",
					request.ContentLength, config.Limit)
				config.ErrorHandler(writer, request, err)
				return
			}

			request.Body = http.MaxBytesReader(writer, request.Body, config.Limit)

			next.ServeHTTP(writer, request)
		})
	}
}

func BodyLimitDefault(next http.Handler) http.Handler {
	config := DefaultBodyLimitConfig()
	return BodyLimit(config)(next)
}

func shouldExclude(request *http.Request, excludedPaths []string, excludedMethods []string) bool {
	for _, method := range excludedMethods {
		if request.Method == method {
			return true
		}
	}

	for _, path := range excludedPaths {
		if matchPath(request.URL.Path, path) {
			return true
		}
	}

	return false
}

func matchPath(requestPath, pattern string) bool {
	if strings.HasSuffix(pattern, "*") {
		prefix := strings.TrimSuffix(pattern, "*")
		return strings.HasPrefix(requestPath, prefix)
	}
	return requestPath == pattern
}

func ParseSize(sizeString string) (int64, error) {
	sizeString = strings.TrimSpace(strings.ToUpper(sizeString))
	if sizeString == "" {
		return 0, fmt.Errorf("empty size string")
	}

	var multiplier int64 = 1
	var numberPart string

	switch {
	case strings.HasSuffix(sizeString, "KB"):
		multiplier = 1 << 10
		numberPart = strings.TrimSuffix(sizeString, "KB")
	case strings.HasSuffix(sizeString, "MB"):
		multiplier = 1 << 20
		numberPart = strings.TrimSuffix(sizeString, "MB")
	case strings.HasSuffix(sizeString, "GB"):
		multiplier = 1 << 30
		numberPart = strings.TrimSuffix(sizeString, "GB")
	case strings.HasSuffix(sizeString, "K"):
		multiplier = 1 << 10
		numberPart = strings.TrimSuffix(sizeString, "K")
	case strings.HasSuffix(sizeString, "M"):
		multiplier = 1 << 20
		numberPart = strings.TrimSuffix(sizeString, "M")
	case strings.HasSuffix(sizeString, "G"):
		multiplier = 1 << 30
		numberPart = strings.TrimSuffix(sizeString, "G")
	case strings.HasSuffix(sizeString, "B"):
		numberPart = strings.TrimSuffix(sizeString, "B")
	default:
		numberPart = sizeString
	}

	number, err := strconv.ParseInt(strings.TrimSpace(numberPart), 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid size format: %s", sizeString)
	}

	return number * multiplier, nil
}
