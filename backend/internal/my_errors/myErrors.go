package myerrors

import "errors"

var (
	ErrNotFound      = errors.New("Not found")
	ErrBadRequest    = errors.New("Bad request")
	ErrAlreadyExists = errors.New("already exists")
)
