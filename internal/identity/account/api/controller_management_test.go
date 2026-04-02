package api_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/koubae/game-hangar/internal/identity/account"
	"github.com/koubae/game-hangar/internal/testunit"
	"github.com/koubae/game-hangar/pkg/errspkg"
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
			expectedErr: errspkg.ResourceNotFound,

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

func TestAccountManagementController_GetAccount(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		scope     string
		accountID string

		expected    *account.Account
		expectedErr error

		expectedCode        int
		expectedErrResponse string
	}{
		"get-account": {
			scope:     "testing:account:read",
			accountID: testunit.AccountIDTest01Str,

			expected:    testunit.AccountTest01,
			expectedErr: nil,

			expectedCode:        200,
			expectedErrResponse: "",
		},
		"record-is-not-found": {
			scope:     "testing:account:read",
			accountID: testunit.AccountIDTest01Str,

			expected:    nil,
			expectedErr: errspkg.ResourceNotFound,

			expectedCode:        404,
			expectedErrResponse: `{"code":404,"message":"resource not found"}`,
		},
		"err-forbidden-scope-missing": {
			scope:     "",
			accountID: testunit.AccountIDTest01Str,

			expected:    nil,
			expectedErr: errspkg.ResourceNotFound,

			expectedCode:        403,
			expectedErrResponse: `{"code":403,"message":"user does not have permission to read account"}`,
		},
		"err-forbidden-scope-missing-permissions": {
			scope:     "testing:auth:read",
			accountID: testunit.AccountIDTest01Str,

			expected:    nil,
			expectedErr: errspkg.ResourceNotFound,

			expectedCode:        403,
			expectedErrResponse: `{"code":403,"message":"user does not have permission to read account"}`,
		},
	}
	for id, tt := range tests {
		t.Run(
			id, func(t *testing.T) {
				t.Parallel()

				_, handlerPtr, mocker := testunit.NewTestRouterAndContainer(t)
				handler := *handlerPtr

				mocker.MockGetAccount(tt.accountID, tt.expected, tt.expectedErr)

				req, err := http.NewRequest("GET", fmt.Sprintf("/api/v1/backoffice/account/%s", tt.accountID), nil)
				require.NoError(t, err)
				mocker.GenAdminAccessTokenAndSetInReq(
					t,
					req,
					testunit.AdminAccountIDTest01Str,
					testunit.AdminUsernameTest01,
					tt.scope,
				)

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
