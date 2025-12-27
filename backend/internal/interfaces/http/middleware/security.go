package middleware

import (
	"net/http"
)

type SecurityHeadersConfig struct {
	ContentTypeNosniff     bool
	XFrameOptions          string
	XSSProtection          bool
	ContentSecurityPolicy  string
	StrictTransportSecurity string
	ReferrerPolicy         string
	PermissionsPolicy      string
}

func DefaultSecurityHeadersConfig() SecurityHeadersConfig {
	return SecurityHeadersConfig{
		ContentTypeNosniff:      true,
		XFrameOptions:           "DENY",
		XSSProtection:           true,
		ContentSecurityPolicy:   "default-src 'self'",
		StrictTransportSecurity: "max-age=31536000; includeSubDomains",
		ReferrerPolicy:          "strict-origin-when-cross-origin",
		PermissionsPolicy:       "geolocation=(), microphone=(), camera=()",
	}
}

func SecurityHeaders(config SecurityHeadersConfig) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			headers := writer.Header()

			if config.ContentTypeNosniff {
				headers.Set("X-Content-Type-Options", "nosniff")
			}

			if config.XFrameOptions != "" {
				headers.Set("X-Frame-Options", config.XFrameOptions)
			}

			if config.XSSProtection {
				headers.Set("X-XSS-Protection", "1; mode=block")
			}

			if config.ContentSecurityPolicy != "" {
				headers.Set("Content-Security-Policy", config.ContentSecurityPolicy)
			}

			if config.StrictTransportSecurity != "" {
				headers.Set("Strict-Transport-Security", config.StrictTransportSecurity)
			}

			if config.ReferrerPolicy != "" {
				headers.Set("Referrer-Policy", config.ReferrerPolicy)
			}

			if config.PermissionsPolicy != "" {
				headers.Set("Permissions-Policy", config.PermissionsPolicy)
			}

			headers.Set("X-Permitted-Cross-Domain-Policies", "none")

			next.ServeHTTP(writer, request)
		})
	}
}

func SecureHeaders(next http.Handler) http.Handler {
	config := DefaultSecurityHeadersConfig()
	return SecurityHeaders(config)(next)
}
