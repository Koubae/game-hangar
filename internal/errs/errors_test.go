package errs_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/koubae/game-hangar/internal/errs"
	"github.com/koubae/game-hangar/pkg/database"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAppError_AsAppError_WhenUnmappedWrapsOriginalError(t *testing.T) {
	t.Parallel()

	errCustom := errors.New("custom error")

	appError := errs.AsAppError(errCustom)

	assert.IsType(t, &errs.AppError{}, appError)
	assert.Equal(t, "unknown error", appError.Msg)
	assert.True(t, appError.IsUnmapped())

	assert.True(t, errors.Is(appError, errs.Unmapped))
	assert.True(t, errors.Is(appError, errCustom))
}

func TestAppError_MultipleWrappedErrorsContainChainedErrorContext(t *testing.T) {
	t.Parallel()

	serviceErr := errors.New("service-a-error")
	appError := &errs.AppError{
		Msg: "service-a-error",
		Err: errors.Join(
			serviceErr, &errs.AppError{
				Msg: "repository-error",
				Err: &errs.AppError{
					Msg: "record not found",
					Err: database.ErrNotFound,
				},
			},
		),
	}

	assert.Equal(t, "service-a-error", appError.Msg)
	assert.False(t, appError.IsUnmapped())

	assert.True(t, errors.Is(appError, serviceErr))
	assert.True(t, errors.Is(appError, database.ErrNotFound))
}

func TestAppError_ErrorWrapChaining(t *testing.T) {
	t.Parallel()

	errLevel1 := errors.New("level-1")
	errLevel2 := fmt.Errorf("level-2: %w", errLevel1)
	errLevel3 := fmt.Errorf("level-3: %w", errLevel2)
	errLevel4 := fmt.Errorf("level-4: %w", errLevel3)

	tests := map[string]struct {
		err        error
		switchCase int
		expected   []error
	}{
		"unmapped": {
			err:        errors.New("unknown-error"),
			switchCase: 0,
			expected:   []error{errs.Unmapped},
		},
		"lvl-1": {
			err:        errLevel1,
			switchCase: 1,
			expected:   []error{errLevel1},
		},
		"lvl-2": {
			err:        errLevel2,
			switchCase: 2,
			expected:   []error{errLevel1, errLevel2},
		},
	}

	for id, tt := range tests {
		t.Run(
			id, func(t *testing.T) {
				t.Parallel()

				err := errs.AsAppError(tt.err)

				switch tt.switchCase {
				case 1:
					require.ErrorAs(t, err, &errLevel1)
				case 2:
					require.ErrorAs(t, err, &errLevel2)
				case 3:
					require.ErrorAs(t, err, &errLevel3)
				case 4:
					require.ErrorAs(t, err, &errLevel4)
				default:
					require.ErrorAs(t, err, &errs.Unmapped)
				}

				for _, errExpected := range tt.expected {
					assert.True(t, errors.Is(err, errExpected))
				}

				caseSwitchHit := -1
				switch { // switch for exceptions should be built "upside down", from more specific to less specific
				case errors.Is(err, errLevel4):
					caseSwitchHit = 4
				case errors.Is(err, errLevel3):
					caseSwitchHit = 3
				case errors.Is(err, errLevel2):
					caseSwitchHit = 2
				case errors.Is(err, errLevel1):
					caseSwitchHit = 1
				default:
					caseSwitchHit = 0
				}

				assert.Equal(t, tt.switchCase, caseSwitchHit)
			},
		)
	}
}

func TestAppError_DBErrToAppErr(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		err      error
		expected error
	}{
		"unknown-db-err": {
			err:      errors.New("unknown-error"),
			expected: errs.DBError,
		},
		"open-transaction-error": {
			err:      &database.ErrOpenTransaction{Err: errors.New("open-transaction-error")},
			expected: errs.DBError,
		},
		"duplicate-error": {
			err:      &database.ErrDuplicate{Err: errors.New("duplicate-error")},
			expected: errs.ResourceDuplicate,
		},
		"not-found-error": {
			err:      database.ErrNotFound,
			expected: errs.ResourceNotFound,
		},
	}

	for id, tt := range tests {
		t.Run(
			id, func(t *testing.T) {
				t.Parallel()

				err := errs.DBErrToAppErr(tt.err)

				assert.IsType(t, &errs.AppError{}, err)
				assert.False(t, err.IsUnmapped())
				assert.ErrorAs(t, err, &tt.expected)
			},
		)
	}
}

func TestAppError_IsServerErr_IsClientErr(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		err         error
		isServerErr bool
		isClientErr bool
	}{
		"is-server-err-1": {
			err:         errs.ServerErr,
			isServerErr: true,
			isClientErr: false,
		},
		"is-server-err-2": {
			err:         errs.Unmapped,
			isServerErr: true,
			isClientErr: false,
		},
	}

	for id, tt := range tests {
		t.Run(
			id, func(t *testing.T) {
				t.Parallel()

				err := errs.AsAppError(tt.err)

				isServerErr := err.IsServerErr()
				isClientErr := err.IsClientErr()

				assert.IsType(t, &errs.AppError{}, err)
				assert.Equal(t, tt.isServerErr, isServerErr)
				assert.Equal(t, tt.isClientErr, isClientErr)
			},
		)
	}
}
