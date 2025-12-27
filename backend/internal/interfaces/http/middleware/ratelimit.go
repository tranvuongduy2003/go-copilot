package middleware

import (
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/tranvuongduy2003/go-copilot/internal/interfaces/http/response"
)

type RateLimiterConfig struct {
	RequestsPerSecond int
	BurstSize         int
	CleanupInterval   time.Duration
}

func DefaultRateLimiterConfig() RateLimiterConfig {
	return RateLimiterConfig{
		RequestsPerSecond: 10,
		BurstSize:         20,
		CleanupInterval:   time.Minute,
	}
}

type tokenBucket struct {
	tokens     float64
	lastRefill time.Time
	mutex      sync.Mutex
}

type RateLimiter struct {
	buckets           map[string]*tokenBucket
	mutex             sync.RWMutex
	requestsPerSecond float64
	burstSize         float64
	cleanupInterval   time.Duration
	stopCleanup       chan struct{}
}

func NewRateLimiter(config RateLimiterConfig) *RateLimiter {
	limiter := &RateLimiter{
		buckets:           make(map[string]*tokenBucket),
		requestsPerSecond: float64(config.RequestsPerSecond),
		burstSize:         float64(config.BurstSize),
		cleanupInterval:   config.CleanupInterval,
		stopCleanup:       make(chan struct{}),
	}

	go limiter.cleanupRoutine()

	return limiter
}

func (limiter *RateLimiter) Allow(key string) bool {
	limiter.mutex.Lock()
	bucket, exists := limiter.buckets[key]
	if !exists {
		bucket = &tokenBucket{
			tokens:     limiter.burstSize,
			lastRefill: time.Now(),
		}
		limiter.buckets[key] = bucket
	}
	limiter.mutex.Unlock()

	bucket.mutex.Lock()
	defer bucket.mutex.Unlock()

	now := time.Now()
	elapsed := now.Sub(bucket.lastRefill).Seconds()
	bucket.tokens += elapsed * limiter.requestsPerSecond
	if bucket.tokens > limiter.burstSize {
		bucket.tokens = limiter.burstSize
	}
	bucket.lastRefill = now

	if bucket.tokens >= 1 {
		bucket.tokens--
		return true
	}

	return false
}

func (limiter *RateLimiter) RetryAfter(key string) time.Duration {
	limiter.mutex.RLock()
	bucket, exists := limiter.buckets[key]
	limiter.mutex.RUnlock()

	if !exists {
		return 0
	}

	bucket.mutex.Lock()
	defer bucket.mutex.Unlock()

	tokensNeeded := 1 - bucket.tokens
	if tokensNeeded <= 0 {
		return 0
	}

	return time.Duration(tokensNeeded/limiter.requestsPerSecond*1000) * time.Millisecond
}

func (limiter *RateLimiter) cleanupRoutine() {
	ticker := time.NewTicker(limiter.cleanupInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			limiter.cleanup()
		case <-limiter.stopCleanup:
			return
		}
	}
}

func (limiter *RateLimiter) cleanup() {
	limiter.mutex.Lock()
	defer limiter.mutex.Unlock()

	threshold := time.Now().Add(-limiter.cleanupInterval)
	for key, bucket := range limiter.buckets {
		bucket.mutex.Lock()
		if bucket.lastRefill.Before(threshold) {
			delete(limiter.buckets, key)
		}
		bucket.mutex.Unlock()
	}
}

func (limiter *RateLimiter) Stop() {
	close(limiter.stopCleanup)
}

func RateLimit(limiter *RateLimiter) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			key := extractClientIP(request)

			if !limiter.Allow(key) {
				retryAfter := limiter.RetryAfter(key)
				writer.Header().Set("Retry-After", strconv.Itoa(int(retryAfter.Seconds())+1))
				writer.Header().Set("X-RateLimit-Limit", strconv.FormatFloat(limiter.requestsPerSecond, 'f', 0, 64))
				writer.Header().Set("X-RateLimit-Remaining", "0")

				response.JSON(writer, http.StatusTooManyRequests, response.ErrorResponse{
					Error: response.ErrorDetail{
						Code:    "RATE_LIMIT_EXCEEDED",
						Message: "too many requests, please try again later",
					},
					TraceID: GetRequestID(request.Context()),
				})
				return
			}

			next.ServeHTTP(writer, request)
		})
	}
}

func extractClientIP(request *http.Request) string {
	if forwardedFor := request.Header.Get("X-Forwarded-For"); forwardedFor != "" {
		return forwardedFor
	}

	if realIP := request.Header.Get("X-Real-IP"); realIP != "" {
		return realIP
	}

	return request.RemoteAddr
}
