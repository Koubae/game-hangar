package auth_test

import (
	"testing"

	"github.com/koubae/game-hangar/internal/errs"
	"github.com/koubae/game-hangar/internal/identity/auth"
	"github.com/koubae/game-hangar/internal/testunit"
	"github.com/koubae/game-hangar/pkg/testutil"
	"github.com/stretchr/testify/assert"
)

func TestNewAccountCredential_Validate(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		model       *auth.NewAccountCredential
		expectedErr error
	}{
		"validation-success": {
			model: &auth.NewAccountCredential{
				Credential: testunit.UsernameTest01,
				AccountID:  testunit.AccountIDTest01,
				ProviderID: testunit.ProviderUsernameID,
				Secret:     "sha255-secret",
				SecretType: "password",
				Verified:   true,
				VerifiedAt: &testutil.Now,
			},
			expectedErr: nil,
		},
		"validation-err-verified-required": {
			model: &auth.NewAccountCredential{
				Credential: testunit.UsernameTest01,
				AccountID:  testunit.AccountIDTest01,
				ProviderID: testunit.ProviderUsernameID,
				Secret:     "sha255-secret",
				SecretType: "password",
				Verified:   true,
				VerifiedAt: nil,
			},
			expectedErr: errs.AccountCredVerifiedAtRequired,
		},
		"validation-err-nil-when-is-not-verified": {
			model: &auth.NewAccountCredential{
				Credential: testunit.UsernameTest01,
				AccountID:  testunit.AccountIDTest01,
				ProviderID: testunit.ProviderUsernameID,
				Secret:     "sha255-secret",
				SecretType: "password",
				Verified:   false,
				VerifiedAt: &testutil.Now,
			},
			expectedErr: errs.AccountCredVerifiedNilWhenIsFalse,
		},
	}
	for id, tt := range tests {
		t.Run(
			id, func(t *testing.T) {
				t.Parallel()

				err := tt.model.Validate()

				if tt.expectedErr != nil {
					assert.ErrorAs(t, err, &tt.expectedErr)
				} else {
					assert.NoError(t, err)
				}

			},
		)
	}
}

func TestNewAccountCredential_ValidateForTypeUsername(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		model       *auth.NewAccountCredential
		expectedErr error
	}{
		"validation-success": {
			model: &auth.NewAccountCredential{
				Credential: testunit.UsernameTest01,
				AccountID:  testunit.AccountIDTest01,
				ProviderID: testunit.ProviderUsernameID,
				Secret:     "sha255-secret",
				SecretType: "password",
				Verified:   true,
				VerifiedAt: &testutil.Now,
			},
			expectedErr: nil,
		},
		"validation-success-after-trim-space": {
			model: &auth.NewAccountCredential{
				Credential: "    unit-test-user-01    ",
				AccountID:  testunit.AccountIDTest01,
				ProviderID: testunit.ProviderUsernameID,
				Secret:     "sha255-secret",
				SecretType: "password",
				Verified:   true,
				VerifiedAt: &testutil.Now,
			},
			expectedErr: nil,
		},
		"validation-err-username-too-short": {
			model: &auth.NewAccountCredential{
				Credential: "abc",
				AccountID:  testunit.AccountIDTest01,
				ProviderID: testunit.ProviderUsernameID,
				Secret:     "sha255-secret",
				SecretType: "password",
				Verified:   true,
				VerifiedAt: &testutil.Now,
			},
			expectedErr: errs.AccountCredCredentialTooShort,
		},
		"validation-err-username-too-long": {
			model: &auth.NewAccountCredential{
				Credential: "username-exceeds-lengthusername-exceeds-lengthusername-exceeds-lengthusername-exceeds-length",
				AccountID:  testunit.AccountIDTest01,
				ProviderID: testunit.ProviderUsernameID,
				Secret:     "sha255-secret",
				SecretType: "password",
				Verified:   true,
				VerifiedAt: &testutil.Now,
			},
			expectedErr: errs.AccountCredCredentialTooLong,
		},

		"validation-err-pattern-invalid": {
			model: &auth.NewAccountCredential{
				Credential: "!" + testunit.UsernameTest01,
				AccountID:  testunit.AccountIDTest01,
				ProviderID: testunit.ProviderUsernameID,
				Secret:     "sha255-secret",
				SecretType: "password",
				Verified:   true,
				VerifiedAt: &testutil.Now,
			},
			expectedErr: errs.AccountCredCredentialInvalid,
		},

		"validation-err-username-reserved": {
			model: &auth.NewAccountCredential{
				Credential: "admin",
				AccountID:  testunit.AccountIDTest01,
				ProviderID: testunit.ProviderUsernameID,
				Secret:     "sha255-secret",
				SecretType: "password",
				Verified:   true,
				VerifiedAt: &testutil.Now,
			},
			expectedErr: errs.AccountCredCredentialReserved,
		},
	}
	for id, tt := range tests {
		t.Run(
			id, func(t *testing.T) {
				t.Parallel()

				err := tt.model.ValidateForTypeUsername()

				if tt.expectedErr != nil {
					assert.ErrorAs(t, err, &tt.expectedErr)
				} else {
					assert.NoError(t, err)
					assert.Equal(t, testunit.UsernameTest01, tt.model.Credential)
				}

			},
		)
	}
}
