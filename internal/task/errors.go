package task

import "errors"

var (
	ErrTaskNotFound      = errors.New("task not found")
	ErrInvalidTransition = errors.New("invalid status transition")

	errNotImplemented = errors.New("not implemented")
)
