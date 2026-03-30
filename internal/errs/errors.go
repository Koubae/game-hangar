package errs

import (
	"errors"
	"fmt"
	"strings"
)

var (
	Unmapped = errors.New("unmapped error")

	ResourceNotFound  = errors.New("resource not found")
	ResourceDuplicate = errors.New("resource already exists")
)

// AppError Application error that wraps any other errors.
// The Application should be aware and directly handle this error when returning anything as reply to
// the outside world.
// Acts as "Proxy" or "filter" for any other kind of errors.
type AppError struct {
	// Err original wrapped error. It could have multiple errors wrapped inside.
	Err error
	// Optional message to be returned to the outside world.
	Msg string
}

func (e *AppError) Error() string {
	if strings.TrimSpace(e.Msg) == "" {
		return e.Err.Error()
	}
	return fmt.Sprintf("%s, error: %s", e.Msg, e.Err.Error())
}

func (e *AppError) Unwrap() error {
	return e.Err
}

func (e *AppError) IsUnmapped() bool {
	return errors.Is(e.Err, Unmapped)
}

// AsAppError converts a given error to an AppError.
// Wraps Unmapped errors with an "unknown error" message and Unmapped tag.
// NOTE:
//   - Why Join the original error with Unmapped?
//   - So that, at a different app layer, we can check for a smaller set of errors
//     and everything else "defaults" to Unmapped.
func AsAppError(err error) *AppError {
	if appErr, ok := errors.AsType[*AppError](err); ok {
		return appErr
	}
	return &AppError{Err: errors.Join(Unmapped, err), Msg: "unknown error"}
}
