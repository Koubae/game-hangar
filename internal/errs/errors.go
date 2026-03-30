package errs

import (
	"errors"
	"strings"

	"github.com/koubae/game-hangar/pkg/database"
)

var (
	ServerErr = &AppError{
		Err: errors.New("server error"),
		Msg: "unexpected error",
	}
	ClientErr = &AppError{
		Err: errors.New("client error"),
		Msg: "client error",
	}

	Unmapped = &AppError{
		Err: ServerErr,
		Msg: "unmapped error",
	}

	DBError = &AppError{
		Err: ServerErr,
		Msg: "database error",
	}

	ResourceNotFound = &AppError{
		Err: ClientErr,
		Msg: "resource not found",
	}
	ResourceDuplicate = &AppError{
		Err: ClientErr,
		Msg: "resource already exists",
	}

	AuthSecretHash = &AppError{
		Err: ServerErr,
		Msg: "secret hash error",
	}

	ProviderNotFound = &AppError{
		Err: ClientErr,
		Msg: "provider not found",
	}
	ProviderDisabled = &AppError{
		Err: ClientErr,
		Msg: "provider is disabled",
	}

	AccountCredVerifiedAtRequired = &AppError{
		Err: ClientErr,
		Msg: "verified_at is required when verified is true",
	}
	AccountCredVerifiedNilWhenIsFalse = &AppError{
		Err: ClientErr,
		Msg: "verified_at must be nil when verified is false",
	}

	AccountCredCreateIncorrectProviderType = &AppError{
		Err: ClientErr,
		Msg: "incorrect provider type",
	}
	AccountCredDuplicate = &AppError{
		Err: ClientErr,
		Msg: "credential already exists",
	}

	UsernameRequired = &AppError{
		Err: ClientErr,
		Msg: "username is required",
	}
	InvalidEmailFormat = &AppError{
		Err: ClientErr,
		Msg: "invalid email format",
	}
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
	return e.Msg
}

func (e *AppError) Unwrap() error {
	return e.Err
}

func (e *AppError) IsUnmapped() bool {
	return errors.Is(e, Unmapped)
}

func (e *AppError) IsServerErr() bool {
	return errors.Is(e, ServerErr)
}

func (e *AppError) IsClientErr() bool {
	return errors.Is(e, ClientErr)
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

func Wrap(appErr *AppError, err error) *AppError {
	return &AppError{Err: errors.Join(appErr, err), Msg: appErr.Msg}
}

func IsAny(err error, targets ...error) bool {
	for _, target := range targets {
		if errors.Is(err, target) {
			return true
		}
	}
	return false

}

func DBErrToAppErr(err error) *AppError {
	var mappedErr error
	var message string
	switch {
	case errors.Is(err, database.ErrNotFound):
		mappedErr = ResourceNotFound
		message = "resource not found"
	case errors.Is(err, &database.ErrDuplicate{}):
		mappedErr = ResourceDuplicate
		message = "resource already exists"
	case errors.Is(err, &database.ErrOpenTransaction{}):
		mappedErr = DBError
		message = "database error"
	default:
		mappedErr = DBError
		message = "database error"
	}

	return &AppError{Err: mappedErr, Msg: message}

}
