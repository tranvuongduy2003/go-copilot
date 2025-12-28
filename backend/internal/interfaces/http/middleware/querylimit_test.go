package middleware

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestQueryLimit(t *testing.T) {
	successHandler := http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(http.StatusOK)
		writer.Write([]byte("success"))
	})

	tests := []struct {
		name           string
		config         QueryLimitConfig
		queryString    string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "allows normal query parameters",
			config:         DefaultQueryLimitConfig(),
			queryString:    "page=1&limit=10",
			expectedStatus: http.StatusOK,
			expectedBody:   "success",
		},
		{
			name: "rejects query string exceeding total length",
			config: QueryLimitConfig{
				MaxParameterLength: 256,
				MaxParameterCount:  50,
				MaxTotalLength:     20,
			},
			queryString:    "page=1&limit=10&search=verylongsearchterm",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "QUERY_PARAMETER_LIMIT_EXCEEDED",
		},
		{
			name: "rejects parameter value exceeding length",
			config: QueryLimitConfig{
				MaxParameterLength: 10,
				MaxParameterCount:  50,
				MaxTotalLength:     2048,
			},
			queryString:    "search=thisisaverylongsearchstring",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "value too long",
		},
		{
			name: "rejects parameter name exceeding length",
			config: QueryLimitConfig{
				MaxParameterLength: 5,
				MaxParameterCount:  50,
				MaxTotalLength:     2048,
			},
			queryString:    "verylongparametername=value",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "name",
		},
		{
			name: "rejects too many parameters",
			config: QueryLimitConfig{
				MaxParameterLength: 256,
				MaxParameterCount:  2,
				MaxTotalLength:     2048,
			},
			queryString:    "a=1&b=2&c=3",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "too many query parameters",
		},
		{
			name: "allows empty query string",
			config: QueryLimitConfig{
				MaxParameterLength: 256,
				MaxParameterCount:  50,
				MaxTotalLength:     2048,
			},
			queryString:    "",
			expectedStatus: http.StatusOK,
			expectedBody:   "success",
		},
		{
			name: "handles multiple values for same parameter",
			config: QueryLimitConfig{
				MaxParameterLength: 256,
				MaxParameterCount:  5,
				MaxTotalLength:     2048,
			},
			queryString:    "status=active&status=inactive&status=pending",
			expectedStatus: http.StatusOK,
			expectedBody:   "success",
		},
		{
			name: "counts multiple values toward parameter count",
			config: QueryLimitConfig{
				MaxParameterLength: 256,
				MaxParameterCount:  2,
				MaxTotalLength:     2048,
			},
			queryString:    "status=active&status=inactive&status=pending",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "too many query parameters",
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			middleware := QueryLimit(testCase.config)
			handler := middleware(successHandler)

			url := "/test"
			if testCase.queryString != "" {
				url += "?" + testCase.queryString
			}

			request := httptest.NewRequest(http.MethodGet, url, nil)
			recorder := httptest.NewRecorder()

			handler.ServeHTTP(recorder, request)

			assert.Equal(t, testCase.expectedStatus, recorder.Code)
			assert.Contains(t, recorder.Body.String(), testCase.expectedBody)
		})
	}
}

func TestQueryLimitDefault(t *testing.T) {
	successHandler := http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(http.StatusOK)
	})

	handler := QueryLimitDefault(successHandler)

	request := httptest.NewRequest(http.MethodGet, "/test?page=1&limit=10", nil)
	recorder := httptest.NewRecorder()

	handler.ServeHTTP(recorder, request)

	assert.Equal(t, http.StatusOK, recorder.Code)
}

func TestQueryLimitExcludedPaths(t *testing.T) {
	successHandler := http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(http.StatusOK)
		writer.Write([]byte("success"))
	})

	config := QueryLimitConfig{
		MaxParameterLength: 10,
		MaxParameterCount:  2,
		MaxTotalLength:     20,
		ExcludedPaths:      []string{"/health*", "/metrics"},
	}

	middleware := QueryLimit(config)
	handler := middleware(successHandler)

	tests := []struct {
		name           string
		path           string
		queryString    string
		expectedStatus int
	}{
		{
			name:           "excluded path with wildcard allows long query",
			path:           "/health/ready",
			queryString:    "verylongparameter=verylongvalue",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "exact excluded path allows long query",
			path:           "/metrics",
			queryString:    "verylongparameter=verylongvalue",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "non-excluded path rejects long query",
			path:           "/api/users",
			queryString:    "verylongparameter=verylongvalue",
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			url := testCase.path + "?" + testCase.queryString
			request := httptest.NewRequest(http.MethodGet, url, nil)
			recorder := httptest.NewRecorder()

			handler.ServeHTTP(recorder, request)

			assert.Equal(t, testCase.expectedStatus, recorder.Code)
		})
	}
}

func TestQueryLimitCustomErrorHandler(t *testing.T) {
	successHandler := http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(http.StatusOK)
	})

	customErrorHandlerCalled := false
	config := QueryLimitConfig{
		MaxParameterLength: 10,
		MaxParameterCount:  50,
		MaxTotalLength:     2048,
		ErrorHandler: func(writer http.ResponseWriter, request *http.Request, err error) {
			customErrorHandlerCalled = true
			writer.WriteHeader(http.StatusUnprocessableEntity)
			writer.Write([]byte("custom error"))
		},
	}

	middleware := QueryLimit(config)
	handler := middleware(successHandler)

	request := httptest.NewRequest(http.MethodGet, "/test?verylongparametername=value", nil)
	recorder := httptest.NewRecorder()

	handler.ServeHTTP(recorder, request)

	assert.True(t, customErrorHandlerCalled)
	assert.Equal(t, http.StatusUnprocessableEntity, recorder.Code)
	assert.Equal(t, "custom error", recorder.Body.String())
}

func TestTruncateString(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		maxLength int
		expected  string
	}{
		{
			name:      "short string unchanged",
			input:     "short",
			maxLength: 10,
			expected:  "short",
		},
		{
			name:      "exact length unchanged",
			input:     "exactly10!",
			maxLength: 10,
			expected:  "exactly10!",
		},
		{
			name:      "long string truncated",
			input:     "this is a very long string",
			maxLength: 10,
			expected:  "this is a ...",
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			result := truncateString(testCase.input, testCase.maxLength)
			assert.Equal(t, testCase.expected, result)
		})
	}
}

func TestQueryLimitWithLargeQueryString(t *testing.T) {
	successHandler := http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(http.StatusOK)
	})

	config := QueryLimitConfig{
		MaxParameterLength: 256,
		MaxParameterCount:  50,
		MaxTotalLength:     100,
	}

	middleware := QueryLimit(config)
	handler := middleware(successHandler)

	longValue := strings.Repeat("a", 200)
	request := httptest.NewRequest(http.MethodGet, "/test?data="+longValue, nil)
	recorder := httptest.NewRecorder()

	handler.ServeHTTP(recorder, request)

	assert.Equal(t, http.StatusBadRequest, recorder.Code)
	assert.Contains(t, recorder.Body.String(), "query string too long")
}
