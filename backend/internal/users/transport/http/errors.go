package http

import (
	"errors"
	"mutabaahapi/internal/platform/apperror"
	"mutabaahapi/internal/users/usecase"
)

func mapDomainError(err error) error {
	switch {
	case errors.Is(err, usecase.ErrUserAlreadyExists):
		return apperror.New(apperror.ErrCodeConflict, err.Error())
	case errors.Is(err, usecase.ErrInstitutionNotFound):
		return apperror.New(apperror.ErrCodeNotFound, err.Error())
	case errors.Is(err, usecase.ErrInvalidCredentials), errors.Is(err, usecase.ErrUserNotFound):
		return apperror.New(apperror.ErrCodeUnauthorized, "invalid username, password, or institution code")
	default:
		return err
	}
}
