package errs

import (
	"github.com/koubae/game-hangar/pkg/errspkg"
)

var (
	ProviderNotFound = &errspkg.AppError{
		Err:         errspkg.ClientErr,
		Msg:         "provider not found",
		DefaultCode: 404,
	}
	ProviderDisabled = &errspkg.AppError{
		Err:         errspkg.ClientErr,
		Msg:         "provider is disabled",
		DefaultCode: 403,
	}

	AccountCredVerifiedAtRequired = &errspkg.AppError{
		Err:         errspkg.ClientErr,
		Msg:         "verified_at is required when verified is true",
		DefaultCode: 400,
	}
	AccountCredVerifiedNilWhenIsFalse = &errspkg.AppError{
		Err:         errspkg.ClientErr,
		Msg:         "verified_at must be nil when verified is false",
		DefaultCode: 400,
	}
	AccountCredCredentialRequired = &errspkg.AppError{
		Err:         errspkg.ClientErr,
		Msg:         "credential is required",
		DefaultCode: 400,
	}
	AccountCredCredentialTooShort = &errspkg.AppError{
		Err:         errspkg.ClientErr,
		Msg:         "credential is too short",
		DefaultCode: 400,
	}
	AccountCredCredentialTooLong = &errspkg.AppError{
		Err:         errspkg.ClientErr,
		Msg:         "credential is too long",
		DefaultCode: 400,
	}
	AccountCredCredentialInvalid = &errspkg.AppError{
		Err:         errspkg.ClientErr,
		Msg:         "credential contains invalid characters",
		DefaultCode: 400,
	}
	AccountCredCredentialReserved = &errspkg.AppError{
		Err:         errspkg.ClientErr,
		Msg:         "credential contains reserved characters",
		DefaultCode: 400,
	}

	AccountCredCreateIncorrectProviderType = &errspkg.AppError{
		Err:         errspkg.ClientErr,
		Msg:         "incorrect provider type",
		DefaultCode: 400,
	}
	AccountCredDuplicate = &errspkg.AppError{
		Err:         errspkg.ClientErr,
		Msg:         "credential already exists",
		DefaultCode: 409,
	}

	UsernameRequired = &errspkg.AppError{
		Err:         errspkg.ClientErr,
		Msg:         "username is required",
		DefaultCode: 400,
	}
	InvalidEmailFormat = &errspkg.AppError{
		Err:         errspkg.ClientErr,
		Msg:         "invalid email format",
		DefaultCode: 400,
	}
)
