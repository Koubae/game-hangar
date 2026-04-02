package account_test

import (
	"context"
	"testing"

	"github.com/koubae/game-hangar/internal/identity/account"
	"github.com/koubae/game-hangar/internal/testunit"
	"github.com/koubae/game-hangar/pkg/errspkg"
	"github.com/stretchr/testify/assert"
)

func TestManagementService_GetAccount(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		accountID   string
		expected    *account.Account
		expectedErr error
	}{
		"get-account": {
			accountID:   testunit.AccountIDTest01Str,
			expected:    testunit.AccountTest01,
			expectedErr: nil,
		},
		"record-is-not-found": {
			accountID:   testunit.AccountIDTest01Str,
			expected:    nil,
			expectedErr: errspkg.ResourceNotFound,
		},
	}

	ctx := context.Background()
	for id, tt := range tests {
		t.Run(
			id, func(t *testing.T) {
				t.Parallel()

				container, _, mocker := testunit.NewTestRouterAndContainer(t)
				mocker.MockGetAccount(tt.accountID, tt.expected, tt.expectedErr)

				service := container.AccountManagementService(nil)
				_account, err := service.GetAccount(ctx, tt.accountID)

				if tt.expectedErr != nil {
					assert.Error(t, err)
					assert.ErrorAs(t, err, &tt.expectedErr)
					assert.Nil(t, _account)
				} else {
					assert.NoError(t, err)
					assert.Equal(t, tt.expected, _account)
				}
			},
		)
	}

}
