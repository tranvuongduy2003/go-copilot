package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	HTTPRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "path", "status"},
	)

	HTTPRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP request duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path"},
	)

	HTTPRequestSize = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_size_bytes",
			Help:    "HTTP request size in bytes",
			Buckets: prometheus.ExponentialBuckets(100, 10, 8),
		},
		[]string{"method", "path"},
	)

	HTTPResponseSize = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_response_size_bytes",
			Help:    "HTTP response size in bytes",
			Buckets: prometheus.ExponentialBuckets(100, 10, 8),
		},
		[]string{"method", "path"},
	)

	DatabaseQueryDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "database_query_duration_seconds",
			Help:    "Database query duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"operation", "table"},
	)

	DatabaseQueryErrors = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "database_query_errors_total",
			Help: "Total number of database query errors",
		},
		[]string{"operation", "table", "error_type"},
	)

	DatabaseConnectionsActive = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "database_connections_active",
			Help: "Number of active database connections",
		},
	)

	DatabaseConnectionsIdle = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "database_connections_idle",
			Help: "Number of idle database connections",
		},
	)

	CacheHits = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "cache_hits_total",
			Help: "Total number of cache hits",
		},
		[]string{"cache_type"},
	)

	CacheMisses = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "cache_misses_total",
			Help: "Total number of cache misses",
		},
		[]string{"cache_type"},
	)

	UsersRegistered = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "users_registered_total",
			Help: "Total number of registered users",
		},
	)

	UsersActive = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "users_active",
			Help: "Number of active users",
		},
	)

	EventsPublished = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "domain_events_published_total",
			Help: "Total number of domain events published",
		},
		[]string{"event_type"},
	)

	EventsHandled = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "domain_events_handled_total",
			Help: "Total number of domain events handled",
		},
		[]string{"event_type", "handler"},
	)

	EventHandlerErrors = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "domain_event_handler_errors_total",
			Help: "Total number of domain event handler errors",
		},
		[]string{"event_type", "handler"},
	)

	CacheOperationDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "cache_operation_duration_seconds",
			Help:    "Cache operation duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"operation", "success"},
	)
)

func RecordHTTPRequest(method, path, status string, durationSeconds float64, requestSize, responseSize int64) {
	HTTPRequestsTotal.WithLabelValues(method, path, status).Inc()
	HTTPRequestDuration.WithLabelValues(method, path).Observe(durationSeconds)
	HTTPRequestSize.WithLabelValues(method, path).Observe(float64(requestSize))
	HTTPResponseSize.WithLabelValues(method, path).Observe(float64(responseSize))
}

func RecordDatabaseQuery(operation, table string, durationSeconds float64, err error) {
	DatabaseQueryDuration.WithLabelValues(operation, table).Observe(durationSeconds)
	if err != nil {
		DatabaseQueryErrors.WithLabelValues(operation, table, "error").Inc()
	}
}

func UpdateDatabaseConnections(active, idle int) {
	DatabaseConnectionsActive.Set(float64(active))
	DatabaseConnectionsIdle.Set(float64(idle))
}

func RecordCacheHit(cacheType string) {
	CacheHits.WithLabelValues(cacheType).Inc()
}

func RecordCacheMiss(cacheType string) {
	CacheMisses.WithLabelValues(cacheType).Inc()
}

func RecordUserRegistration() {
	UsersRegistered.Inc()
}

func SetActiveUsers(count int) {
	UsersActive.Set(float64(count))
}

func RecordEventPublished(eventType string) {
	EventsPublished.WithLabelValues(eventType).Inc()
}

func RecordEventHandled(eventType, handler string) {
	EventsHandled.WithLabelValues(eventType, handler).Inc()
}

func RecordEventHandlerError(eventType, handler string) {
	EventHandlerErrors.WithLabelValues(eventType, handler).Inc()
}

func RecordCacheOperation(operation string, success bool, durationSeconds float64) {
	successStr := "false"
	if success {
		successStr = "true"
	}
	CacheOperationDuration.WithLabelValues(operation, successStr).Observe(durationSeconds)

	if success {
		CacheHits.WithLabelValues(operation).Inc()
	} else {
		CacheMisses.WithLabelValues(operation).Inc()
	}
}
