package middleware

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	"go.opentelemetry.io/otel/trace"
)

type TracingConfig struct {
	ServiceName       string
	TracerName        string
	SkipPaths         []string
	RequestIDHeader   string
	PropagateHeaders  bool
}

func DefaultTracingConfig() TracingConfig {
	return TracingConfig{
		ServiceName:      "go-copilot",
		TracerName:       "http-server",
		SkipPaths:        []string{"/health", "/health/live", "/health/ready", "/metrics"},
		RequestIDHeader:  "X-Request-ID",
		PropagateHeaders: true,
	}
}

func Tracing(config TracingConfig) func(http.Handler) http.Handler {
	tracer := otel.Tracer(config.TracerName)
	propagator := otel.GetTextMapPropagator()

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			for _, path := range config.SkipPaths {
				if request.URL.Path == path {
					next.ServeHTTP(writer, request)
					return
				}
			}

			ctx := request.Context()

			if config.PropagateHeaders {
				ctx = propagator.Extract(ctx, propagation.HeaderCarrier(request.Header))
			}

			spanName := buildSpanName(request)

			ctx, span := tracer.Start(ctx, spanName,
				trace.WithSpanKind(trace.SpanKindServer),
				trace.WithAttributes(
					semconv.HTTPRequestMethodKey.String(request.Method),
					semconv.URLPath(request.URL.Path),
					semconv.URLScheme(request.URL.Scheme),
					semconv.ServerAddress(request.Host),
					semconv.UserAgentOriginal(request.UserAgent()),
					semconv.ClientAddress(request.RemoteAddr),
				),
			)
			defer span.End()

			if requestID := request.Header.Get(config.RequestIDHeader); requestID != "" {
				span.SetAttributes(attribute.String("request.id", requestID))
			}

			if query := request.URL.RawQuery; query != "" {
				span.SetAttributes(attribute.String("http.query", query))
			}

			tracingWriter := &tracingResponseWriter{
				ResponseWriter: writer,
				statusCode:     http.StatusOK,
			}

			next.ServeHTTP(tracingWriter, request.WithContext(ctx))

			span.SetAttributes(
				semconv.HTTPResponseStatusCode(tracingWriter.statusCode),
			)

			if tracingWriter.statusCode >= 400 {
				span.SetAttributes(attribute.Bool("error", true))
			}
		})
	}
}

func TracingDefault(next http.Handler) http.Handler {
	config := DefaultTracingConfig()
	return Tracing(config)(next)
}

func buildSpanName(request *http.Request) string {
	routeCtx := chi.RouteContext(request.Context())
	if routeCtx != nil && routeCtx.RoutePattern() != "" {
		return request.Method + " " + routeCtx.RoutePattern()
	}
	return request.Method + " " + request.URL.Path
}

type tracingResponseWriter struct {
	http.ResponseWriter
	statusCode int
	written    bool
}

func (w *tracingResponseWriter) WriteHeader(statusCode int) {
	if !w.written {
		w.statusCode = statusCode
		w.written = true
	}
	w.ResponseWriter.WriteHeader(statusCode)
}

func (w *tracingResponseWriter) Write(data []byte) (int, error) {
	if !w.written {
		w.written = true
	}
	return w.ResponseWriter.Write(data)
}

func (w *tracingResponseWriter) Unwrap() http.ResponseWriter {
	return w.ResponseWriter
}
