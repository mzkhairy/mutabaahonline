package httperr

import (
	"mutabaahapi/internal/platform/apperror"
	"net/http"
)

// StatusFor maps an error (preferably AppError) to an HTTP status code.
func StatusFor(err error) int {
	if err == nil {
		return http.StatusOK
	}
	if _, ok := apperror.IsValidationError(err); ok {
		return http.StatusBadRequest
	}
	if app, ok := apperror.IsAppError(err); ok {
		switch app.Code {
		case apperror.ErrCodeNotFound:
			return http.StatusNotFound
		case apperror.ErrCodeInvalidInput, apperror.ErrCodeBadRequest, apperror.ErrCodeValidationFailed:
			return http.StatusBadRequest
		case apperror.ErrCodeUnauthorized:
			return http.StatusUnauthorized
		case apperror.ErrCodeForbidden:
			return http.StatusForbidden
		case apperror.ErrCodeConflict:
			return http.StatusConflict
		case apperror.ErrCodeTooManyRequests:
			return http.StatusTooManyRequests
		case apperror.ErrCodeExternalService:
			return http.StatusServiceUnavailable
		case apperror.ErrCodeDatabaseError, apperror.ErrCodeInternalError:
			fallthrough
		default:
			return http.StatusInternalServerError
		}
	}
	// default
	return http.StatusInternalServerError
}

type ErrorBody struct {
	RequestID string `json:"request_id,omitempty"`
	Error     struct {
		Code    string `json:"code"`
		Message string `json:"message"`
		Details any    `json:"details,omitempty"`
	} `json:"error"`
}

func Body(err error, reqID string) ErrorBody {
	eb := ErrorBody{RequestID: reqID}
	if ve, ok := apperror.IsValidationError(err); ok {
		eb.Error.Code = apperror.ErrCodeValidationFailed.String()
		if ve.AppError != nil {
			eb.Error.Message = ve.AppError.Message
		} else {
			eb.Error.Message = "validation failed"
		}
		eb.Error.Details = ve.Fields
		return eb
	}
	if app, ok := apperror.IsAppError(err); ok {
		eb.Error.Code = app.Code.String()
		eb.Error.Message = app.Message
		eb.Error.Details = app.Details
		return eb
	}
	// Fallback
	eb.Error.Code = apperror.ErrCodeInternalError.String()
	eb.Error.Message = "internal server error"
	return eb
}
