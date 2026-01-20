package usecase

import "errors"

var (
	ErrUserAlreadyExists   = errors.New("username is already registered in this institution")
	ErrUserNotFound        = errors.New("user not found")
	ErrInstitutionNotFound = errors.New("institution code not found")
	ErrInvalidCredentials  = errors.New("invalid username or password")
	ErrInternalServerError = errors.New("internal server error")
)
