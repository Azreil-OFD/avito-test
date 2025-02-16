package errors

import "errors"

var (
	ErrInvalidPassword   = errors.New("invalid password")
	ErrUnauthorized      = errors.New("unauthorized access")
	ErrInternal          = errors.New("internal server error")
	ErrBadRequest        = errors.New("bad request")
	ErrConflict          = errors.New("resource conflict")
	ErrUserNotFound      = errors.New("user not found")
	ErrInsufficientFound = errors.New("insufficient found")
	ErrSelfTransfer      = errors.New("self transfer")
)
