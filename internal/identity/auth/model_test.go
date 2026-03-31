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
