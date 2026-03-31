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
	"github.com/koubae/game-hangar/pkg/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAuthController_RegisterByUsername(t *testing.T) {
	t.Parallel()

	_, handlerPtr, mocker := testunit.NewTestRouterAndContainer(t)
	handler := *handlerPtr

	mocker.MockGetDefaultUsernameProvider()
	mocker.MockGetCredentialByProvider(testutil.ProviderUsernameID, testutil.UsernameTest01, nil, errs.ResourceNotFound)
	mocker.MockCreateAccountCredential(testutil.CredIDTest01, nil)
	mocker.MockCreateAccount(testutil.UsernameTest01, nil, testutil.AccountIDTest01Str, nil)

	payload := fmt.Sprintf(
		`{
		"source": "global",	
		"username": "%s",
		"password": "StrongPassword123!"
	}`, testutil.UsernameTest01,
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
		AccountID:    testutil.AccountIDTest01Str,
		Username:     testutil.UsernameTest01,
		LoggedCredID: testutil.CredIDTest01,
	}
	assert.Equal(t, expected, response)
	assert.Equal(t, http.StatusCreated, rr.Code)
}

func TestAuthController_RegisterByUsername_ErrOnInValidPassword(t *testing.T) {
	t.Parallel()

	_, handlerPtr, mocker := testunit.NewTestRouterAndContainer(t)
	handler := *handlerPtr

	mocker.MockGetDefaultUsernameProvider()
	mocker.MockGetCredentialByProvider(testutil.ProviderUsernameID, testutil.UsernameTest01, nil, errs.ResourceNotFound)
	mocker.MockCreateAccountCredential(testutil.CredIDTest01, nil)
	mocker.MockCreateAccount(testutil.UsernameTest01, nil, testutil.AccountIDTest01Str, nil)

	payload := fmt.Sprintf(
		`{
		"source": "global",	
		"username": "%s",
		"password": "pass-not-strong"
	}`, testutil.UsernameTest01,
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
