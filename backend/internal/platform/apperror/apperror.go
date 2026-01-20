package apperror

import (
	"errors"
	"fmt"
)

type ErrorCode string

const (
	ErrCodeNotFound         ErrorCode = "NOT_FOUND"
	ErrCodeInvalidInput     ErrorCode = "INVALID_INPUT"
	ErrCodeUnauthorized     ErrorCode = "UNAUTHORIZED"
	ErrCodeForbidden        ErrorCode = "FORBIDDEN"
	ErrCodeConflict         ErrorCode = "CONFLICT"
	ErrCodeInternalError    ErrorCode = "INTERNAL_ERROR"
	ErrCodeValidationFailed ErrorCode = "VALIDATION_FAILED"
	ErrCodeDatabaseError    ErrorCode = "DATABASE_ERROR"
	ErrCodeExternalService  ErrorCode = "EXTERNAL_SERVICE_ERROR"
	ErrCodeBadRequest       ErrorCode = "BAD_REQUEST"
	ErrCodeTooManyRequests  ErrorCode = "TOO_MANY_REQUESTS"
)

func (e ErrorCode) String() string { return string(e) }

type AppError struct {
	Code    ErrorCode `json:"code"`
	Message string    `json:"message"`
	Details any       `json:"details,omitempty"`
	Cause   error     `json:"-"`
}

func (e *AppError) Error() string {
	if e.Details != nil {
		return fmt.Sprintf("%s: %s (%v)", e.Code, e.Message, e.Details)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

func (e *AppError) Unwrap() error { return e.Cause }

func IsAppError(err error) (*AppError, bool) {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr, true
	}
	return nil, false
}

func New(code ErrorCode, message string) *AppError {
	return &AppError{Code: code, Message: message}
}

func NewWithCause(code ErrorCode, message string, cause error) *AppError {
	return &AppError{Code: code, Message: message, Cause: cause}
}

func NewWithDetails(code ErrorCode, message string, details any) *AppError {
	return &AppError{Code: code, Message: message, Details: details}
}

// Wrap returns an AppError if err is not nil. If err is already an AppError, it is returned as-is.
func Wrap(err error, code ErrorCode, message string) error {
	if err == nil {
		return nil
	}
	if app, ok := IsAppError(err); ok {
		return app
	}
	return NewWithCause(code, message, err)
}

// Some common predefined errors (optional conveniences)
var (
	ErrNotFound      = New(ErrCodeNotFound, "resource not found")
	ErrUnauthorized  = New(ErrCodeUnauthorized, "authentication required")
	ErrForbidden     = New(ErrCodeForbidden, "access forbidden")
	ErrInternalError = New(ErrCodeInternalError, "internal server error")
	ErrInvalidInput  = New(ErrCodeInvalidInput, "invalid input provided")
	ErrConflict      = New(ErrCodeConflict, "resource conflict")
)

// ValidationErrors represents a collection of validation errors keyed by field.
type ValidationErrors struct {
	*AppError
	Fields map[string]string `json:"fields"`
}

func NewValidationErrors(fields map[string]string) *ValidationErrors {
	return &ValidationErrors{
		AppError: &AppError{Code: ErrCodeValidationFailed, Message: "validation failed", Details: fields},
		Fields:   fields,
	}
}

func IsValidationError(err error) (*ValidationErrors, bool) {
	var ve *ValidationErrors
	if errors.As(err, &ve) {
		return ve, true
	}
	return nil, false
}
