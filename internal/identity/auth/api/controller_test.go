package api_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/koubae/game-hangar/internal/errs"
	"github.com/koubae/game-hangar/internal/identity/auth"
	"github.com/koubae/game-hangar/internal/testunit"
	"github.com/koubae/game-hangar/pkg/authpkg"
	"github.com/koubae/game-hangar/pkg/errspkg"
	"github.com/koubae/game-hangar/pkg/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAuthController_RegisterByUsername(t *testing.T) {
	t.Parallel()

	_, handlerPtr, mocker := testunit.NewTestRouterAndContainer(t)
	handler := *handlerPtr

	mocker.MockGetDefaultUsernameProvider()
	mocker.MockGetCredentialByProvider(
		testunit.ProviderUsernameID, testunit.UsernameTest01, nil,
		errspkg.ResourceNotFound,
	)
	mocker.MockCreateAccountCredential(testunit.CredIDTest01, nil)
	mocker.MockCreateAccount(testunit.UsernameTest01, nil, testunit.AccountIDTest01Str, nil)

	payload := fmt.Sprintf(
		`{
		"source": "global",	
		"username": "%s",
		"password": "%s"
	}`, testunit.UsernameTest01,
		testunit.StrongPassword,
	)

	req, err := http.NewRequest("POST", "/api/v1/auth/register/username", strings.NewReader(payload))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	var response auth.DTOAccountLoggedIn
	err = json.Unmarshal([]byte(rr.Body.String()), &response)
	require.NoError(t, err)

	expected := auth.DTOAccountLoggedIn{
		AccountID:    testunit.AccountIDTest01Str,
		Username:     testunit.UsernameTest01,
		LoggedCredID: testunit.CredIDTest01,
	}
	assert.Equal(t, expected, response)
	assert.Equal(t, http.StatusCreated, rr.Code)
}

func TestAuthController_RegisterByUsername_ErrOnInValidPassword(t *testing.T) {
	t.Parallel()

	_, handlerPtr, mocker := testunit.NewTestRouterAndContainer(t)
	handler := *handlerPtr

	mocker.MockGetDefaultUsernameProvider()
	mocker.MockGetCredentialByProvider(
		testunit.ProviderUsernameID, testunit.UsernameTest01, nil,
		errspkg.ResourceNotFound,
	)
	mocker.MockCreateAccountCredential(testunit.CredIDTest01, nil)
	mocker.MockCreateAccount(testunit.UsernameTest01, nil, testunit.AccountIDTest01Str, nil)

	payload := fmt.Sprintf(
		`{
		"source": "global",	
		"username": "%s",
		"password": "pass-not-strong"
	}`, testunit.UsernameTest01,
	)

	req, err := http.NewRequest("POST", "/api/v1/auth/register/username", strings.NewReader(payload))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	response := rr.Body.String()
	expected := `{"code":400,"message":"password validation error, error: at least one uppercase letter is required\nat least one digit is required\n"}`
	assert.Equal(t, expected, strings.TrimSpace(response))
	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestAuthController_RegisterByUsername_ErrOnInValidUsername(t *testing.T) {
	t.Parallel()

	_, handlerPtr, mocker := testunit.NewTestRouterAndContainer(t)
	handler := *handlerPtr

	username := "!invalid-username"

	mocker.MockGetDefaultUsernameProvider()
	mocker.MockGetCredentialByProvider(testunit.ProviderUsernameID, username, nil, errspkg.ResourceNotFound)
	mocker.MockCreateAccountCredential(testunit.CredIDTest01, nil)
	mocker.MockCreateAccount(username, nil, testunit.AccountIDTest01Str, nil)

	payload := fmt.Sprintf(
		`{
		"source": "global",	
		"username": "%s",
		"password": "StrongPassword123!"
	}`, username,
	)

	req, err := http.NewRequest("POST", "/api/v1/auth/register/username", strings.NewReader(payload))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	response := rr.Body.String()
	expected := `{"code":400,"message":"could not create account: credential contains invalid characters"}`
	assert.Equal(t, expected, strings.TrimSpace(response))
	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestAuthController_LoginByUsername(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		providerReturned *auth.Provider
		providerGetErr   error

		credentialReturned *auth.AccountCredential
		credentialGetErr   error

		expectedCode        int
		expectedErrResponse string
	}{
		"login-success": {
			providerReturned: &auth.Provider{ID: 1, Source: "global", Type: "username", Disabled: false},
			providerGetErr:   nil,

			credentialReturned: &auth.AccountCredential{
				ID:         1,
				Credential: testunit.UsernameTest01,
				AccountID:  testunit.AccountIDTest01,
				ProviderID: 1,
				Secret:     testunit.StrongPasswordHash,
			},

			credentialGetErr: nil,
			expectedCode:     http.StatusOK,

			expectedErrResponse: "",
		},
		"err-on-provider-not-found": {
			providerReturned: nil,
			providerGetErr:   errs.ProviderNotFound,

			credentialReturned: nil,
			credentialGetErr:   nil,

			expectedCode:        http.StatusUnauthorized,
			expectedErrResponse: `{"code":401,"message":"login failed, error: provider not found"}`,
		},
		"err-on-provider-disabled": {
			providerReturned: &auth.Provider{Disabled: true, ID: 1, Source: "global", Type: "username"},
			providerGetErr:   nil,

			credentialReturned: nil,
			credentialGetErr:   nil,

			expectedCode:        http.StatusUnauthorized,
			expectedErrResponse: `{"code":401,"message":"login failed, error: provider is disabled"}`,
		},
		"err-on-credential-not-found": {
			providerReturned: &auth.Provider{ID: 1, Source: "global", Type: "username", Disabled: false},
			providerGetErr:   nil,

			credentialReturned: nil,
			credentialGetErr:   errs.ProviderNotFound,

			expectedCode:        http.StatusUnauthorized,
			expectedErrResponse: `{"code":401,"message":"credential login failed, error: provider not found"}`,
		},
	}
	for id, tt := range tests {
		t.Run(
			id, func(t *testing.T) {
				t.Parallel()

				_, handlerPtr, mocker := testunit.NewTestRouterAndContainer(t)
				handler := *handlerPtr

				mocker.MockGetProvider("global", "username", tt.providerReturned, tt.providerGetErr)
				mocker.MockGetCredentialByProvider(
					testunit.ProviderUsernameID,
					testunit.UsernameTest01,
					tt.credentialReturned,
					tt.credentialGetErr,
				)

				payload := fmt.Sprintf(
					`{
					"source": "global",	
					"username": "%s",
					"password": "%s"
				}`, testunit.UsernameTest01,
					testunit.StrongPassword,
				)

				req, err := http.NewRequest("POST", "/api/v1/auth/login/username", strings.NewReader(payload))
				require.NoError(t, err)
				req.Header.Set("Content-Type", "application/json")

				rr := httptest.NewRecorder()
				handler.ServeHTTP(rr, req)

				if tt.expectedErrResponse == "" {
					var response auth.DTOAccessToken
					err = json.Unmarshal([]byte(rr.Body.String()), &response)
					require.NoError(t, err)

					assert.IsType(t, auth.DTOAccessToken{}, response)
					assert.NotEmpty(t, response.AccessToken)
					assert.NotEmpty(t, response.ExpiresIn)
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

func TestAuthController_LoginAdminByUsername(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		scope            string
		providerReturned *auth.Provider
		providerGetErr   error

		credentialReturned *auth.AccountCredential
		credentialGetErr   error

		accountPermissionsReturned []*auth.Permission

		expectedScope       string
		expectedCode        int
		expectedErrResponse string
	}{
		"login-success": {
			scope: "identity:account_credential:*|storage:config:read,write,delete|identity:*:read",

			providerReturned: testunit.ProviderUsername,
			providerGetErr:   nil,

			credentialReturned: testunit.AccountCredentialTest01,
			credentialGetErr:   nil,

			accountPermissionsReturned: testunit.AdminAccountPermissions,

			expectedScope:       "identity:auth:read|identity:account:read|storage:config:read",
			expectedCode:        http.StatusOK,
			expectedErrResponse: "",
		},
		"login-return-all-permissions": {
			scope: "*:*:*",

			providerReturned: testunit.ProviderUsername,
			providerGetErr:   nil,

			credentialReturned: testunit.AccountCredentialTest01,
			credentialGetErr:   nil,

			accountPermissionsReturned: testunit.AdminAccountPermissions,

			expectedScope:       "identity:auth:read,write,delete|identity:account:read,write,delete|storage:config:read|storage:setting:read,write|leaderboard:leaderboard:read|chat:*:*",
			expectedCode:        http.StatusOK,
			expectedErrResponse: "",
		},
		"login-return-all-of-action": {
			scope: "*:*:read",

			providerReturned: testunit.ProviderUsername,
			providerGetErr:   nil,

			credentialReturned: testunit.AccountCredentialTest01,
			credentialGetErr:   nil,

			accountPermissionsReturned: testunit.AdminAccountPermissions,

			expectedScope:       "identity:auth:read|identity:account:read|storage:config:read|storage:setting:read|leaderboard:leaderboard:read|chat:*:read",
			expectedCode:        http.StatusOK,
			expectedErrResponse: "",
		},
		"login-return-all-of-resource": {
			scope: "*:auth:*",

			providerReturned: testunit.ProviderUsername,
			providerGetErr:   nil,

			credentialReturned: testunit.AccountCredentialTest01,
			credentialGetErr:   nil,

			accountPermissionsReturned: testunit.AdminAccountPermissions,

			expectedScope:       "identity:auth:read,write,delete|chat:auth:*",
			expectedCode:        http.StatusOK,
			expectedErrResponse: "",
		},
		"login-return-all-of-service": {
			scope: "identity:*:*",

			providerReturned: testunit.ProviderUsername,
			providerGetErr:   nil,

			credentialReturned: testunit.AccountCredentialTest01,
			credentialGetErr:   nil,

			accountPermissionsReturned: testunit.AdminAccountPermissions,

			expectedScope:       "identity:auth:read,write,delete|identity:account:read,write,delete",
			expectedCode:        http.StatusOK,
			expectedErrResponse: "",
		},
		"err-requested-scope-has-no-permission": {
			scope: "matchmaking:*:*|campaigns:*:*|leaderboard:leaderboard:write|storage:config:delete",

			providerReturned: testunit.ProviderUsername,
			providerGetErr:   nil,

			credentialReturned: testunit.AccountCredentialTest01,
			credentialGetErr:   nil,

			accountPermissionsReturned: testunit.AdminAccountPermissions,

			expectedScope:       "",
			expectedCode:        http.StatusUnauthorized,
			expectedErrResponse: `{"code":401,"message":"login failed, error: no permissions found"}`,
		},

		"err-validation-scope-not-provided": {
			scope: "",

			providerReturned: testunit.ProviderUsername,
			providerGetErr:   nil,

			credentialReturned: testunit.AccountCredentialTest01,
			credentialGetErr:   nil,

			accountPermissionsReturned: testunit.AdminAccountPermissions,

			expectedScope:       "",
			expectedCode:        http.StatusBadRequest,
			expectedErrResponse: `{"code":400,"message":"payload validation error, error: field 'scope' is required"}`,
		},

		"err-on-provider-not-found": {
			scope: "identity:account_credential:*|storage:config:read,write,delete|identity:*:read",

			providerReturned: nil,
			providerGetErr:   errs.ProviderNotFound,

			credentialReturned: nil,
			credentialGetErr:   nil,

			accountPermissionsReturned: testunit.AdminAccountPermissions,

			expectedScope:       "",
			expectedCode:        http.StatusUnauthorized,
			expectedErrResponse: `{"code":401,"message":"login failed, error: provider not found"}`,
		},
		"err-on-provider-disabled": {
			scope: "identity:account_credential:*|storage:config:read,write,delete|identity:*:read",

			providerReturned: &auth.Provider{Disabled: true, ID: 1, Source: "global", Type: "username"},
			providerGetErr:   nil,

			credentialReturned: nil,
			credentialGetErr:   nil,

			accountPermissionsReturned: testunit.AdminAccountPermissions,

			expectedScope:       "",
			expectedCode:        http.StatusUnauthorized,
			expectedErrResponse: `{"code":401,"message":"login failed, error: provider is disabled"}`,
		},
		"err-on-credential-not-found": {
			scope: "identity:account_credential:*|storage:config:read,write,delete|identity:*:read",

			providerReturned: &auth.Provider{ID: 1, Source: "global", Type: "username", Disabled: false},
			providerGetErr:   nil,

			credentialReturned: nil,
			credentialGetErr:   errs.ProviderNotFound,

			accountPermissionsReturned: testunit.AdminAccountPermissions,

			expectedScope:       "",
			expectedCode:        http.StatusUnauthorized,
			expectedErrResponse: `{"code":401,"message":"credential login failed, error: provider not found"}`,
		},
		"err-on-permissions-not-found": {
			scope: "identity:account_credential:*|storage:config:read,write,delete|identity:*:read",

			providerReturned: testunit.ProviderUsername,
			providerGetErr:   nil,

			credentialReturned: testunit.AccountCredentialTest01,
			credentialGetErr:   nil,

			accountPermissionsReturned: []*auth.Permission{},

			expectedScope:       "",
			expectedCode:        http.StatusUnauthorized,
			expectedErrResponse: `{"code":401,"message":"login failed, error: no permissions found"}`,
		},
		"err-on-scope-requested-incorrect-format": {
			scope: "identity:account_credential|storage:|",

			providerReturned: testunit.ProviderUsername,
			providerGetErr:   nil,

			credentialReturned: testunit.AccountCredentialTest01,
			credentialGetErr:   nil,

			accountPermissionsReturned: testunit.AdminAccountPermissions,

			expectedScope:       "",
			expectedCode:        http.StatusUnauthorized,
			expectedErrResponse: `{"code":401,"message":"login failed, error: permissions scope format error"}`,
		},
	}
	for id, tt := range tests {
		t.Run(
			id, func(t *testing.T) {
				t.Parallel()

				_, handlerPtr, mocker := testunit.NewTestRouterAndContainer(t)
				handler := *handlerPtr

				mocker.MockGetProvider("global", "username", tt.providerReturned, tt.providerGetErr)
				mocker.MockGetCredentialByProvider(
					testunit.ProviderUsernameID,
					testunit.UsernameTest01,
					tt.credentialReturned,
					tt.credentialGetErr,
				)
				mocker.MockGetAdminAccountPermissions(testunit.AccountIDTest01Str, tt.accountPermissionsReturned)

				payload := fmt.Sprintf(
					`{
					"source": "global",	
					"username": "%s",
					"password": "%s",
					"scope": "%s"
				}`, testunit.UsernameTest01,
					testunit.StrongPassword,
					tt.scope,
				)

				req, err := http.NewRequest(
					"POST",
					"/api/v1/backoffice/auth/login/username",
					strings.NewReader(payload),
				)
				require.NoError(t, err)
				req.Header.Set("Content-Type", "application/json")

				rr := httptest.NewRecorder()
				handler.ServeHTTP(rr, req)

				if tt.expectedErrResponse == "" {
					var response auth.DTOAdminAccessToken
					err = json.Unmarshal([]byte(rr.Body.String()), &response)
					require.NoError(t, err)

					expectedPermission, _ := authpkg.ParsePermissions(tt.expectedScope)
					accessToken, err := testutil.ExtractAdminAccessToken(response.AccessToken)
					require.NoError(t, err)

					assert.IsType(t, auth.DTOAdminAccessToken{}, response)
					assert.NotEmpty(t, response.AccessToken)
					assert.NotEmpty(t, response.ExpiresIn)
					assert.Equal(t, tt.expectedCode, rr.Code)
					assert.Equal(t, expectedPermission, response.Permissions)
					assert.Equal(t, expectedPermission, accessToken.Permissions)
				} else {
					response := rr.Body.String()
					assert.Equal(t, tt.expectedErrResponse, strings.TrimSpace(response))
					assert.Equal(t, tt.expectedCode, rr.Code)
				}

			},
		)
	}
}
