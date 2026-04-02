package errspkg

import (
	"errors"
	"fmt"
	"strings"

	"github.com/koubae/game-hangar/pkg/database"
)

var (
	ServerErr = &AppError{
		Err:         errors.New("server error"),
		Msg:         "unexpected error",
		DefaultCode: 500,
	}
	ClientErr = &AppError{
		Err:         errors.New("client error"),
		Msg:         "client error",
		DefaultCode: 400,
	}

	Unmapped = &AppError{
		Err:         ServerErr,
		Msg:         "unmapped error",
		DefaultCode: 500,
	}

	DBError = &AppError{
		Err:         ServerErr,
		Msg:         "database error",
		DefaultCode: 503,
	}

	PayloadValidation = &AppError{
		Err:         ClientErr,
		Msg:         "payload validation error",
		DefaultCode: 400,
	}
	ResourceNotFound = &AppError{
		Err:         ClientErr,
		Msg:         "resource not found",
		DefaultCode: 404,
	}
	ResourceDuplicate = &AppError{
		Err:         ClientErr,
		Msg:         "resource already exists",
		DefaultCode: 409,
	}

	AuthPermissionsScopeEmpty = &AppError{
		Err:         ClientErr,
		Msg:         "no permissions found",
		DefaultCode: 401,
	}
	AuthPermissionsScopeFormat = &AppError{
		Err:         ClientErr,
		Msg:         "permissions scope format error",
		DefaultCode: 401,
	}

	AuthSecretHash = &AppError{
		Err:         ServerErr,
		Msg:         "secret hash error",
		DefaultCode: 500,
	}
	AuthPasswordValidation = &AppError{
		Err:         ClientErr,
		Msg:         "password validation error",
		DefaultCode: 400,
	}
	AuthLoginPasswordMismatch = &AppError{
		Err:         ClientErr,
		Msg:         "password mismatch",
		DefaultCode: 401,
	}
	AuthLoginFailed = &AppError{
		Err:         ClientErr,
		Msg:         "login failed",
		DefaultCode: 401,
	}
	AuthNotLoggedIn = &AppError{
		Err:         ClientErr,
		Msg:         "not logged in",
		DefaultCode: 401,
	}

	PayloadMissingID = &AppError{
		Err:         ClientErr,
		Msg:         "payload missing id",
		DefaultCode: 400,
	}
	InvalidUUID = &AppError{
		Err:         ClientErr,
		Msg:         "invalid uuid",
		DefaultCode: 400,
	}
)

// Validator interface for validating a Request Body Payload
type Validator interface {
	Validate() error
}

// AppError Application error that wraps any other errors.
// The Application should be aware and directly handle this error when returning anything as reply to
// the outside world.
// Acts as "Proxy" or "filter" for any other kind of errors.
type AppError struct {
	// Err original wrapped error. It could have multiple errors wrapped inside.
	Err error
	// Optional message to be returned to the outside world.
	Msg string
	// DefaultCode represents the default HTTP status code associated with the error for external responses.
	DefaultCode int
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

func (e *AppError) GetDefaultCode() int {
	if e.DefaultCode < 200 || e.DefaultCode > 599 {
		if e.IsClientErr() {
			return 400
		}
		return 500
	}
	return e.DefaultCode
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

// Wrap wraps an error with an AppError.
// NOTE: Avoid Wrap the same error multiple times as this will "concatenate" the previous Msg
func Wrap(appErr *AppError, err error) *AppError {
	msg := fmt.Sprintf("%s, error: %s", appErr.Msg, err.Error())
	return &AppError{Err: errors.Join(appErr, err), Msg: msg, DefaultCode: appErr.DefaultCode}
}

func NewWithMsg(appErr *AppError, msg string) *AppError {
	msg = fmt.Sprintf("%s %s", appErr.Msg, msg)
	return &AppError{Err: appErr.Err, Msg: msg, DefaultCode: appErr.DefaultCode}
}

func IsAny(err error, targets ...error) bool {
	for _, target := range targets {
		if errors.Is(err, target) {
			return true
		}
	}
	return false

}

func DBErrToAppErr(err error, resource string) *AppError {
	var mappedErr *AppError
	wrapErr := true
	switch {
	case errors.Is(err, database.ErrNotFound):
		mappedErr = ResourceNotFound
		wrapErr = false
	case errors.Is(err, &database.ErrDuplicate{}):
		mappedErr = ResourceDuplicate
		wrapErr = false
	case errors.Is(err, &database.ErrOpenTransaction{}):
		mappedErr = DBError
	default:
		mappedErr = DBError
	}

	mappedErr = &AppError{
		Err:         mappedErr,
		Msg:         fmt.Sprintf("%s %s", resource, mappedErr.Msg),
		DefaultCode: mappedErr.DefaultCode,
	}
	if wrapErr {
		mappedErr = Wrap(mappedErr, err)
	}
	return mappedErr
}
