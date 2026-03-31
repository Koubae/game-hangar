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
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAuthController_RegisterByUsername(t *testing.T) {
	t.Parallel()

	_, handlerPtr, mocker := testunit.NewTestRouterAndContainer(t)
	handler := *handlerPtr

	mocker.MockGetDefaultUsernameProvider()
	mocker.MockGetCredentialByProvider(testunit.ProviderUsernameID, testunit.UsernameTest01, nil, errs.ResourceNotFound)
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
	mocker.MockGetCredentialByProvider(testunit.ProviderUsernameID, testunit.UsernameTest01, nil, errs.ResourceNotFound)
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
	mocker.MockGetCredentialByProvider(testunit.ProviderUsernameID, username, nil, errs.ResourceNotFound)
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

	_, handlerPtr, mocker := testunit.NewTestRouterAndContainer(t)
	handler := *handlerPtr

	mocker.MockGetProvider(
		"global",
		"username",
		&auth.Provider{ID: 1, Source: "global", Type: "username", Disabled: false},
		nil,
	)
	mocker.MockGetCredentialByProvider(
		testunit.ProviderUsernameID,
		testunit.UsernameTest01,
		&auth.AccountCredential{
			ID:         1,
			Credential: testunit.UsernameTest01,
			AccountID:  testunit.AccountIDTest01,
			ProviderID: 1,
			Secret:     testunit.StrongPasswordHash,
		},
		nil,
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

	var response auth.DTOAccessToken
	err = json.Unmarshal([]byte(rr.Body.String()), &response)
	require.NoError(t, err)

	assert.IsType(t, auth.DTOAccessToken{}, response)
	assert.NotEmpty(t, response.AccessToken)
	assert.NotEmpty(t, response.ExpiresIn)
	assert.Equal(t, http.StatusOK, rr.Code)
}
