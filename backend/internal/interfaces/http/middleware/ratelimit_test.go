package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRateLimiter_Allow(t *testing.T) {
	tests := []struct {
		name           string
		config         RateLimiterConfig
		requestCount   int
		expectedAllow  int
		expectedDenied int
	}{
		{
			name: "allows requests within burst limit",
			config: RateLimiterConfig{
				RequestsPerSecond: 10,
				BurstSize:         5,
				CleanupInterval:   time.Minute,
			},
			requestCount:   5,
			expectedAllow:  5,
			expectedDenied: 0,
		},
		{
			name: "denies requests exceeding burst limit",
			config: RateLimiterConfig{
				RequestsPerSecond: 10,
				BurstSize:         3,
				CleanupInterval:   time.Minute,
			},
			requestCount:   5,
			expectedAllow:  3,
			expectedDenied: 2,
		},
		{
			name: "login rate limiter config",
			config: LoginRateLimiterConfig(),
			requestCount:   7,
			expectedAllow:  5,
			expectedDenied: 2,
		},
		{
			name: "register rate limiter config",
			config: RegisterRateLimiterConfig(),
			requestCount:   5,
			expectedAllow:  3,
			expectedDenied: 2,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			limiter := NewRateLimiter(testCase.config)
			defer limiter.Stop()

			allowed := 0
			denied := 0
			key := "test-client"

			for i := 0; i < testCase.requestCount; i++ {
				if limiter.Allow(key) {
					allowed++
				} else {
					denied++
				}
			}

			assert.Equal(t, testCase.expectedAllow, allowed)
			assert.Equal(t, testCase.expectedDenied, denied)
		})
	}
}

func TestRateLimiter_TokenRefill(t *testing.T) {
	config := RateLimiterConfig{
		RequestsPerSecond: 100,
		BurstSize:         1,
		CleanupInterval:   time.Minute,
	}
	limiter := NewRateLimiter(config)
	defer limiter.Stop()

	key := "test-client"

	assert.True(t, limiter.Allow(key))
	assert.False(t, limiter.Allow(key))

	time.Sleep(20 * time.Millisecond)

	assert.True(t, limiter.Allow(key))
}

func TestRateLimiter_RetryAfter(t *testing.T) {
	config := RateLimiterConfig{
		RequestsPerSecond: 10,
		BurstSize:         1,
		CleanupInterval:   time.Minute,
	}
	limiter := NewRateLimiter(config)
	defer limiter.Stop()

	key := "test-client"

	assert.True(t, limiter.Allow(key))
	assert.False(t, limiter.Allow(key))

	retryAfter := limiter.RetryAfter(key)
	assert.Greater(t, retryAfter, time.Duration(0))
	assert.LessOrEqual(t, retryAfter, 200*time.Millisecond)
}

func TestRateLimiter_DifferentKeys(t *testing.T) {
	config := RateLimiterConfig{
		RequestsPerSecond: 10,
		BurstSize:         2,
		CleanupInterval:   time.Minute,
	}
	limiter := NewRateLimiter(config)
	defer limiter.Stop()

	assert.True(t, limiter.Allow("client1"))
	assert.True(t, limiter.Allow("client1"))
	assert.False(t, limiter.Allow("client1"))

	assert.True(t, limiter.Allow("client2"))
	assert.True(t, limiter.Allow("client2"))
	assert.False(t, limiter.Allow("client2"))
}

func TestRateLimitMiddleware(t *testing.T) {
	config := RateLimiterConfig{
		RequestsPerSecond: 10,
		BurstSize:         2,
		CleanupInterval:   time.Minute,
	}
	limiter := NewRateLimiter(config)
	defer limiter.Stop()

	nextHandler := http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(http.StatusOK)
	})

	handler := RateLimit(limiter)(nextHandler)

	request := httptest.NewRequest(http.MethodGet, "/test", nil)
	request.RemoteAddr = "192.168.1.1:12345"

	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, request)
	assert.Equal(t, http.StatusOK, recorder.Code)

	recorder = httptest.NewRecorder()
	handler.ServeHTTP(recorder, request)
	assert.Equal(t, http.StatusOK, recorder.Code)

	recorder = httptest.NewRecorder()
	handler.ServeHTTP(recorder, request)
	assert.Equal(t, http.StatusTooManyRequests, recorder.Code)
	assert.NotEmpty(t, recorder.Header().Get("Retry-After"))
	assert.Equal(t, "0", recorder.Header().Get("X-RateLimit-Remaining"))
}

func TestRateLimitMiddleware_ExtractsIP(t *testing.T) {
	config := RateLimiterConfig{
		RequestsPerSecond: 10,
		BurstSize:         1,
		CleanupInterval:   time.Minute,
	}
	limiter := NewRateLimiter(config)
	defer limiter.Stop()

	nextHandler := http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(http.StatusOK)
	})

	handler := RateLimit(limiter)(nextHandler)

	tests := []struct {
		name          string
		setupRequest  func(*http.Request)
		expectedAllow bool
	}{
		{
			name: "uses X-Forwarded-For header",
			setupRequest: func(request *http.Request) {
				request.Header.Set("X-Forwarded-For", "10.0.0.1")
				request.RemoteAddr = "192.168.1.1:12345"
			},
			expectedAllow: true,
		},
		{
			name: "uses X-Real-IP header when no X-Forwarded-For",
			setupRequest: func(request *http.Request) {
				request.Header.Set("X-Real-IP", "10.0.0.2")
				request.RemoteAddr = "192.168.1.1:12345"
			},
			expectedAllow: true,
		},
		{
			name: "uses RemoteAddr when no headers",
			setupRequest: func(request *http.Request) {
				request.RemoteAddr = "10.0.0.3:12345"
			},
			expectedAllow: true,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, "/test", nil)
			testCase.setupRequest(request)

			recorder := httptest.NewRecorder()
			handler.ServeHTTP(recorder, request)

			if testCase.expectedAllow {
				assert.Equal(t, http.StatusOK, recorder.Code)
			} else {
				assert.Equal(t, http.StatusTooManyRequests, recorder.Code)
			}
		})
	}
}

func TestDefaultRateLimiterConfig(t *testing.T) {
	config := DefaultRateLimiterConfig()

	assert.Equal(t, 10, config.RequestsPerSecond)
	assert.Equal(t, 20, config.BurstSize)
	assert.Equal(t, time.Minute, config.CleanupInterval)
}

func TestLoginRateLimiterConfig(t *testing.T) {
	config := LoginRateLimiterConfig()

	assert.Equal(t, 1, config.RequestsPerSecond)
	assert.Equal(t, 5, config.BurstSize)
}

func TestRegisterRateLimiterConfig(t *testing.T) {
	config := RegisterRateLimiterConfig()

	assert.Equal(t, 1, config.RequestsPerSecond)
	assert.Equal(t, 3, config.BurstSize)
}

func TestPasswordResetRateLimiterConfig(t *testing.T) {
	config := PasswordResetRateLimiterConfig()

	assert.Equal(t, 1, config.RequestsPerSecond)
	assert.Equal(t, 3, config.BurstSize)
}

func TestTokenRefreshRateLimiterConfig(t *testing.T) {
	config := TokenRefreshRateLimiterConfig()

	assert.Equal(t, 1, config.RequestsPerSecond)
	assert.Equal(t, 30, config.BurstSize)
}
