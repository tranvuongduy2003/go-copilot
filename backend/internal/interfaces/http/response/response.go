package response

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/tranvuongduy2003/go-copilot/internal/domain/shared"
	"github.com/tranvuongduy2003/go-copilot/pkg/validator"
)

type SuccessResponse struct {
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
	Meta    interface{} `json:"meta,omitempty"`
}

type ErrorResponse struct {
	Error   ErrorDetail `json:"error"`
	TraceID string      `json:"trace_id,omitempty"`
}

type ErrorDetail struct {
	Code    string        `json:"code"`
	Message string        `json:"message"`
	Details []FieldError  `json:"details,omitempty"`
}

type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func JSON(writer http.ResponseWriter, statusCode int, data interface{}) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(statusCode)
	if data != nil {
		json.NewEncoder(writer).Encode(data)
	}
}

func Success(writer http.ResponseWriter, data interface{}) {
	JSON(writer, http.StatusOK, SuccessResponse{Data: data})
}

func SuccessWithMessage(writer http.ResponseWriter, data interface{}, message string) {
	JSON(writer, http.StatusOK, SuccessResponse{Data: data, Message: message})
}

func SuccessWithMeta(writer http.ResponseWriter, data interface{}, meta interface{}) {
	JSON(writer, http.StatusOK, SuccessResponse{Data: data, Meta: meta})
}

func Created(writer http.ResponseWriter, data interface{}) {
	JSON(writer, http.StatusCreated, SuccessResponse{Data: data})
}

func CreatedWithLocation(writer http.ResponseWriter, data interface{}, location string) {
	writer.Header().Set("Location", location)
	JSON(writer, http.StatusCreated, SuccessResponse{Data: data})
}

func NoContent(writer http.ResponseWriter) {
	writer.WriteHeader(http.StatusNoContent)
}

func Error(writer http.ResponseWriter, request *http.Request, err error) {
	statusCode, response := mapErrorToResponse(err)
	response.TraceID = GetRequestID(request)
	JSON(writer, statusCode, response)
}

func BadRequest(writer http.ResponseWriter, request *http.Request, message string) {
	response := ErrorResponse{
		Error: ErrorDetail{
			Code:    "BAD_REQUEST",
			Message: message,
		},
		TraceID: GetRequestID(request),
	}
	JSON(writer, http.StatusBadRequest, response)
}

func Unauthorized(writer http.ResponseWriter, request *http.Request, message string) {
	if message == "" {
		message = "unauthorized"
	}
	response := ErrorResponse{
		Error: ErrorDetail{
			Code:    "UNAUTHORIZED",
			Message: message,
		},
		TraceID: GetRequestID(request),
	}
	JSON(writer, http.StatusUnauthorized, response)
}

func Forbidden(writer http.ResponseWriter, request *http.Request, message string) {
	if message == "" {
		message = "forbidden"
	}
	response := ErrorResponse{
		Error: ErrorDetail{
			Code:    "FORBIDDEN",
			Message: message,
		},
		TraceID: GetRequestID(request),
	}
	JSON(writer, http.StatusForbidden, response)
}

func NotFound(writer http.ResponseWriter, request *http.Request, message string) {
	if message == "" {
		message = "resource not found"
	}
	response := ErrorResponse{
		Error: ErrorDetail{
			Code:    "NOT_FOUND",
			Message: message,
		},
		TraceID: GetRequestID(request),
	}
	JSON(writer, http.StatusNotFound, response)
}

func InternalError(writer http.ResponseWriter, request *http.Request) {
	response := ErrorResponse{
		Error: ErrorDetail{
			Code:    "INTERNAL_ERROR",
			Message: "an internal error occurred",
		},
		TraceID: GetRequestID(request),
	}
	JSON(writer, http.StatusInternalServerError, response)
}

func ValidationError(writer http.ResponseWriter, request *http.Request, errs validator.ValidationErrors) {
	details := make([]FieldError, len(errs))
	for i, err := range errs {
		details[i] = FieldError{
			Field:   err.Field,
			Message: err.Message,
		}
	}

	response := ErrorResponse{
		Error: ErrorDetail{
			Code:    "VALIDATION_ERROR",
			Message: "validation failed",
			Details: details,
		},
		TraceID: GetRequestID(request),
	}
	JSON(writer, http.StatusBadRequest, response)
}

func mapErrorToResponse(err error) (int, ErrorResponse) {
	if validationErrs, ok := validator.GetValidationErrors(err); ok {
		details := make([]FieldError, len(validationErrs))
		for i, validationErr := range validationErrs {
			details[i] = FieldError{
				Field:   validationErr.Field,
				Message: validationErr.Message,
			}
		}
		return http.StatusBadRequest, ErrorResponse{
			Error: ErrorDetail{
				Code:    "VALIDATION_ERROR",
				Message: "validation failed",
				Details: details,
			},
		}
	}

	var notFoundErr *shared.NotFoundError
	if errors.As(err, &notFoundErr) {
		return http.StatusNotFound, ErrorResponse{
			Error: ErrorDetail{
				Code:    string(shared.ErrCodeNotFound),
				Message: notFoundErr.Error(),
			},
		}
	}

	var validationErr *shared.ValidationError
	if errors.As(err, &validationErr) {
		return http.StatusBadRequest, ErrorResponse{
			Error: ErrorDetail{
				Code:    string(shared.ErrCodeValidation),
				Message: validationErr.Error(),
				Details: []FieldError{
					{
						Field:   validationErr.Field,
						Message: validationErr.Message,
					},
				},
			},
		}
	}

	var conflictErr *shared.ConflictError
	if errors.As(err, &conflictErr) {
		return http.StatusConflict, ErrorResponse{
			Error: ErrorDetail{
				Code:    string(shared.ErrCodeConflict),
				Message: conflictErr.Error(),
			},
		}
	}

	var authzErr *shared.AuthorizationError
	if errors.As(err, &authzErr) {
		return http.StatusForbidden, ErrorResponse{
			Error: ErrorDetail{
				Code:    string(shared.ErrCodeAuthorization),
				Message: authzErr.Error(),
			},
		}
	}

	var businessRuleErr *shared.BusinessRuleViolationError
	if errors.As(err, &businessRuleErr) {
		return http.StatusUnprocessableEntity, ErrorResponse{
			Error: ErrorDetail{
				Code:    string(shared.ErrCodeBusinessRule),
				Message: businessRuleErr.Error(),
			},
		}
	}

	var statusTransitionErr *shared.InvalidStatusTransitionError
	if errors.As(err, &statusTransitionErr) {
		return http.StatusUnprocessableEntity, ErrorResponse{
			Error: ErrorDetail{
				Code:    string(shared.ErrCodeInvalidStatusTrans),
				Message: statusTransitionErr.Error(),
			},
		}
	}

	return http.StatusInternalServerError, ErrorResponse{
		Error: ErrorDetail{
			Code:    string(shared.ErrCodeInternal),
			Message: "an internal error occurred",
		},
	}
}

func GetRequestID(request *http.Request) string {
	if request == nil {
		return ""
	}
	return request.Header.Get("X-Request-ID")
}
