package erros

import "errors"

var (
	ErrUserAlreadyExists = errors.New("user with this email or phone already exists")
	ErrUserNoExists      = errors.New("user with this email or phone no exists")
	ErrInvalidPassword   = errors.New("invalid password")
)
