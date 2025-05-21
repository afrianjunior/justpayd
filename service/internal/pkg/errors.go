package pkg

import "errors"

// Common errors
var (
	ErrNotFound = errors.New("resource not found")
)

type ValidationError struct {
	Message string
}

func (e ValidationError) Error() string {
	return e.Message
}

func NewValidationError(message string) ValidationError {
	return ValidationError{Message: message}
}

type UnauthorizedError struct {
	Message string
}

func (e UnauthorizedError) Error() string {
	return e.Message
}

func NewUnauthorizedError(message string) UnauthorizedError {
	return UnauthorizedError{Message: message}
}
