package domain

import "errors"

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrInvalidCredential  = errors.New("invalid credentials")
	ErrUnauthorized       = errors.New("unauthorized")
)
