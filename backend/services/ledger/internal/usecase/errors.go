package usecase

import "errors"

var (
	ErrNotFound        = errors.New("not found")
	ErrConflict        = errors.New("already exists")
	ErrInvalidArgument = errors.New("invalid argument")
)
