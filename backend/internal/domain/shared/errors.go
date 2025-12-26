package shared

import (
	"errors"
	"fmt"
)

type ErrorCode string

const (
	ErrCodeNotFound           ErrorCode = "NOT_FOUND"
	ErrCodeValidation         ErrorCode = "VALIDATION_ERROR"
	ErrCodeConflict           ErrorCode = "CONFLICT"
	ErrCodeAuthorization      ErrorCode = "AUTHORIZATION_ERROR"
	ErrCodeBusinessRule       ErrorCode = "BUSINESS_RULE_VIOLATION"
	ErrCodeInternal           ErrorCode = "INTERNAL_ERROR"
	ErrCodeInvalidStatusTrans ErrorCode = "INVALID_STATUS_TRANSITION"
)

type DomainError interface {
	error
	Code() ErrorCode
	Unwrap() error
}

type baseDomainError struct {
	code    ErrorCode
	message string
	cause   error
}

func (e *baseDomainError) Error() string {
	if e.cause != nil {
		return fmt.Sprintf("%s: %v", e.message, e.cause)
	}
	return e.message
}

func (e *baseDomainError) Code() ErrorCode {
	return e.code
}

func (e *baseDomainError) Unwrap() error {
	return e.cause
}

type NotFoundError struct {
	baseDomainError
	EntityType string
	Identifier string
}

func NewNotFoundError(entityType, identifier string) *NotFoundError {
	return &NotFoundError{
		baseDomainError: baseDomainError{
			code:    ErrCodeNotFound,
			message: fmt.Sprintf("%s with identifier '%s' not found", entityType, identifier),
		},
		EntityType: entityType,
		Identifier: identifier,
	}
}

func (e *NotFoundError) Is(target error) bool {
	t, ok := target.(*NotFoundError)
	if !ok {
		return false
	}
	return e.EntityType == t.EntityType || t.EntityType == ""
}

type ValidationError struct {
	baseDomainError
	Field   string
	Message string
}

func NewValidationError(field, message string) *ValidationError {
	return &ValidationError{
		baseDomainError: baseDomainError{
			code:    ErrCodeValidation,
			message: fmt.Sprintf("validation error on field '%s': %s", field, message),
		},
		Field:   field,
		Message: message,
	}
}

func (e *ValidationError) Is(target error) bool {
	t, ok := target.(*ValidationError)
	if !ok {
		return false
	}
	return e.Field == t.Field || t.Field == ""
}

type ConflictError struct {
	baseDomainError
	EntityType string
	Field      string
	Value      string
}

func NewConflictError(entityType, field, value string) *ConflictError {
	return &ConflictError{
		baseDomainError: baseDomainError{
			code:    ErrCodeConflict,
			message: fmt.Sprintf("%s with %s '%s' already exists", entityType, field, value),
		},
		EntityType: entityType,
		Field:      field,
		Value:      value,
	}
}

func (e *ConflictError) Is(target error) bool {
	t, ok := target.(*ConflictError)
	if !ok {
		return false
	}
	return (e.EntityType == t.EntityType || t.EntityType == "") &&
		(e.Field == t.Field || t.Field == "")
}

type AuthorizationError struct {
	baseDomainError
	Action   string
	Resource string
}

func NewAuthorizationError(action, resource string) *AuthorizationError {
	return &AuthorizationError{
		baseDomainError: baseDomainError{
			code:    ErrCodeAuthorization,
			message: fmt.Sprintf("not authorized to %s on %s", action, resource),
		},
		Action:   action,
		Resource: resource,
	}
}

func (e *AuthorizationError) Is(target error) bool {
	_, ok := target.(*AuthorizationError)
	return ok
}

type BusinessRuleViolationError struct {
	baseDomainError
	Rule string
}

func NewBusinessRuleViolationError(rule, message string) *BusinessRuleViolationError {
	return &BusinessRuleViolationError{
		baseDomainError: baseDomainError{
			code:    ErrCodeBusinessRule,
			message: message,
		},
		Rule: rule,
	}
}

func (e *BusinessRuleViolationError) Is(target error) bool {
	t, ok := target.(*BusinessRuleViolationError)
	if !ok {
		return false
	}
	return e.Rule == t.Rule || t.Rule == ""
}

type InvalidStatusTransitionError struct {
	baseDomainError
	CurrentStatus string
	TargetStatus  string
}

func NewInvalidStatusTransitionError(current, target string) *InvalidStatusTransitionError {
	return &InvalidStatusTransitionError{
		baseDomainError: baseDomainError{
			code:    ErrCodeInvalidStatusTrans,
			message: fmt.Sprintf("cannot transition from '%s' to '%s'", current, target),
		},
		CurrentStatus: current,
		TargetStatus:  target,
	}
}

func (e *InvalidStatusTransitionError) Is(target error) bool {
	_, ok := target.(*InvalidStatusTransitionError)
	return ok
}

type InternalError struct {
	baseDomainError
}

func NewInternalError(message string, cause error) *InternalError {
	return &InternalError{
		baseDomainError: baseDomainError{
			code:    ErrCodeInternal,
			message: message,
			cause:   cause,
		},
	}
}

func (e *InternalError) Is(target error) bool {
	_, ok := target.(*InternalError)
	return ok
}

func IsNotFoundError(err error) bool {
	var e *NotFoundError
	return errors.As(err, &e)
}

func IsValidationError(err error) bool {
	var e *ValidationError
	return errors.As(err, &e)
}

func IsConflictError(err error) bool {
	var e *ConflictError
	return errors.As(err, &e)
}

func IsAuthorizationError(err error) bool {
	var e *AuthorizationError
	return errors.As(err, &e)
}

func IsBusinessRuleViolationError(err error) bool {
	var e *BusinessRuleViolationError
	return errors.As(err, &e)
}

func IsInvalidStatusTransitionError(err error) bool {
	var e *InvalidStatusTransitionError
	return errors.As(err, &e)
}

func GetErrorCode(err error) ErrorCode {
	var domainErr DomainError
	if errors.As(err, &domainErr) {
		return domainErr.Code()
	}
	return ErrCodeInternal
}
