package common

import "fmt"

type BusinessError struct {
	HTTPCode int
	Message  string
}

func (e *BusinessError) Error() string {
	return fmt.Sprintf("[%d]: %s", e.HTTPCode, e.Message)
}

type ErrServerError struct {
	Err error // underlying DB error
}

func (e *ErrServerError) Error() string {
	return fmt.Sprintf("server error, error: %v", e.Err)
}

func (e *ErrServerError) Unwrap() error {
	return e.Err
}

func (e *ErrServerError) Is(
	target error,
) bool { // NOTE: Sentinel struct pattern
	_, ok := target.(*ErrServerError)
	return ok
}
