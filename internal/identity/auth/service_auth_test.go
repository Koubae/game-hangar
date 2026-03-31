package auth_test

import (
	"errors"
	"testing"

	"github.com/koubae/game-hangar/internal/identity/auth"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewPasswordValidationRules_MinMaxLengthValidations(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		minLength   int
		maxLength   int
		expectedMin int
		expectedMax int
	}{
		"min-max-inbounds": {
			minLength:   8,
			maxLength:   16,
			expectedMin: 8,
			expectedMax: 16,
		},

		"min-below-lower-bounds": {
			minLength:   auth.MinPasswordLengthLowerBound - 1,
			maxLength:   16,
			expectedMin: auth.MinPasswordLengthLowerBound,
			expectedMax: 16,
		},
		"min-above-upper-bounds": {
			minLength:   auth.MaxPasswordLengthUpperBound + 1,
			maxLength:   16,
			expectedMin: auth.MaxPasswordLengthUpperBound,
			expectedMax: auth.MaxPasswordLengthUpperBound,
		},

		"max-below-lower-bounds": {
			minLength:   8,
			maxLength:   auth.MinPasswordLengthLowerBound - 1,
			expectedMin: 8,
			expectedMax: 8,
		},
		"max-above-upper-bounds": {
			minLength:   8,
			maxLength:   auth.MaxPasswordLengthUpperBound + 1,
			expectedMin: 8,
			expectedMax: auth.MaxPasswordLengthUpperBound,
		},
	}
	for id, tt := range tests {
		t.Run(
			id, func(t *testing.T) {
				t.Parallel()

				rules := auth.NewPasswordValidationRules(tt.minLength, tt.maxLength, true, true, true, true)

				assert.Equal(t, tt.expectedMin, rules.MinLength)
				assert.Equal(t, tt.expectedMax, rules.MaxLength)
			},
		)
	}
}

func TestSecretsService_ValidatePassword(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		password    string
		rules       *auth.PasswordValidationRules
		expectedErr *auth.PasswordValidationErrList
	}{
		"valid-password": {
			password: "PassUnit123!",
			rules: auth.NewPasswordValidationRules(
				8,
				12,
				true,
				true,
				true,
				true,
			),
			expectedErr: nil,
		},
		"valid-password-dev-rules": {
			password: "test",
			rules: auth.NewPasswordValidationRules(
				4,
				72,
				false,
				false,
				false,
				false,
			),
			expectedErr: nil,
		},
		"err-too-short": {
			password: "PaU1!",
			rules: auth.NewPasswordValidationRules(
				8,
				12,
				true,
				true,
				true,
				true,
			),
			expectedErr: &auth.PasswordValidationErrList{MinLength: "password below minimum length 8\n"},
		},
		"err-too-long": {
			password: "PassUnit123!PassUnit123!PassUnit123!PassUnit123!PassUnit123!PassUnit123!",
			rules: auth.NewPasswordValidationRules(
				8,
				12,
				true,
				true,
				true,
				true,
			),
			expectedErr: &auth.PasswordValidationErrList{MaxLength: "password above maximum length 12\n"},
		},
		"err-missing-upper": {
			password: "passunit123!",
			rules: auth.NewPasswordValidationRules(
				8,
				12,
				true,
				true,
				true,
				true,
			),
			expectedErr: &auth.PasswordValidationErrList{Uppercase: "at least one uppercase letter is required\n"},
		},
		"err-missing-lower": {
			password: "PASSUNIT123!",
			rules: auth.NewPasswordValidationRules(
				8,
				12,
				true,
				true,
				true,
				true,
			),
			expectedErr: &auth.PasswordValidationErrList{Lowercase: "at least one lowercase letter is required\n"},
		},
		"err-missing-digit": {
			password: "PassUnitOne!",
			rules: auth.NewPasswordValidationRules(
				8,
				12,
				true,
				true,
				true,
				true,
			),
			expectedErr: &auth.PasswordValidationErrList{Digits: "at least one digit is required\n"},
		},
		"err-missing-special": {
			password: "PassUnit123",
			rules: auth.NewPasswordValidationRules(
				8,
				12,
				true,
				true,
				true,
				true,
			),
			expectedErr: &auth.PasswordValidationErrList{Special: "at least one special character is required\n"},
		},
		"err-multiple": {
			password: "pass",
			rules: auth.NewPasswordValidationRules(
				8,
				12,
				true,
				true,
				true,
				true,
			),
			expectedErr: &auth.PasswordValidationErrList{
				MinLength: "password below minimum length 8\n",
				Uppercase: "at least one uppercase letter is required\n",
				Digits:    "at least one digit is required\n",
				Special:   "at least one special character is required\n",
			},
		},
	}

	service := auth.NewSecretsService()
	for id, tt := range tests {
		t.Run(
			id, func(t *testing.T) {
				t.Parallel()

				err := service.ValidatePassword(tt.password, *tt.rules)
				if tt.expectedErr == nil {
					assert.NoError(t, err)
					return
				}

				errValidation, ok := errors.AsType[*auth.PasswordValidationErrList](err)
				require.True(t, ok)

				assert.Error(t, errValidation)
				assert.Equal(t, tt.expectedErr.Error(), errValidation.Error())

				assert.Equal(t, tt.expectedErr.MinLength, errValidation.MinLength)
				assert.Equal(t, tt.expectedErr.MaxLength, errValidation.MaxLength)
				assert.Equal(t, tt.expectedErr.Uppercase, errValidation.Uppercase)
				assert.Equal(t, tt.expectedErr.Lowercase, errValidation.Lowercase)
				assert.Equal(t, tt.expectedErr.Digits, errValidation.Digits)
				assert.Equal(t, tt.expectedErr.Special, errValidation.Special)

			},
		)
	}
}
