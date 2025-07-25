package errs

import (
	"errors"
)

type AppError struct {
	Code int
	Msg  string
	Err  error
}

func (e *AppError) Error() string {
	return e.Msg
}

func (e *AppError) Unwrap() error {
	return e.Err
}

func New(code int, msg string, err error) *AppError {
	return &AppError{code, msg, err}
}

var (
	ErrUserNotFound      = errors.New("user is not found")
	ErrTaskNotFound      = errors.New("task in not found")
	ErrUnexpected        = errors.New("unexpected error")
	ErrInvalidUserId     = errors.New("invalid user id")
	ErrInvalidTaskId     = errors.New("invalid task id")
	ErrIncorrectPassword = errors.New("incorrect password")
	ErrUsernameExists    = errors.New("username is already exists")
)
