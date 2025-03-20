package web

import "errors"

type shutdownError struct {
	Message string
}

func NewShutdownError(message string) error {
	return &shutdownError{Message: message}
}

// Error is the implementation of the error interface
func (se *shutdownError) Error() string {
	return se.Message
}

// IsShutdown tests if the error contains a shutdown error
func IsShutdown(err error) bool {
	var se *shutdownError

	return errors.As(err, &se)
}
