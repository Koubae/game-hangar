package api_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/koubae/game-hangar/internal/errs"
	"github.com/koubae/game-hangar/internal/identity/account"
	"github.com/koubae/game-hangar/internal/testunit"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAccountManagementController_Me(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		accountID string

		expected    *account.Account
		expectedErr error

		expectedCode        int
		expectedErrResponse string
	}{
		"get-account": {
			accountID: testunit.AccountIDTest01Str,

			expected:    testunit.AccountTest01,
			expectedErr: nil,

			expectedCode:        200,
			expectedErrResponse: "",
		},
		"record-is-not-found": {
			accountID: testunit.AccountIDTest01Str,

			expected:    nil,
			expectedErr: errs.ResourceNotFound,

			expectedCode:        404,
			expectedErrResponse: `{"code":404,"message":"resource not found"}`,
		},
	}
	for id, tt := range tests {
		t.Run(
			id, func(t *testing.T) {
				t.Parallel()

				_, handlerPtr, mocker := testunit.NewTestRouterAndContainer(t)
				handler := *handlerPtr

				mocker.MockGetAccount(tt.accountID, tt.expected, tt.expectedErr)

				req, err := http.NewRequest("GET", "/api/v1/account/me", nil)
				require.NoError(t, err)
				mocker.GenAccessTokenAndSetInReq(t, req, tt.accountID, testunit.UsernameTest01)

				rr := httptest.NewRecorder()
				handler.ServeHTTP(rr, req)

				if tt.expectedErrResponse == "" {
					var response account.DTOAccount
					err = json.Unmarshal([]byte(rr.Body.String()), &response)
					require.NoError(t, err)

					assert.IsType(t, account.DTOAccount{}, response)
					assert.Equal(t, tt.expected.ID, response.ID)
					assert.Equal(t, tt.expectedCode, rr.Code)
				} else {
					response := rr.Body.String()
					assert.Equal(t, tt.expectedErrResponse, strings.TrimSpace(response))
					assert.Equal(t, tt.expectedCode, rr.Code)
				}

			},
		)
	}
}
