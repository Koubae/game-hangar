package api_test

// import (
// 	"encoding/json"
// 	"fmt"
// 	"net/http"
// 	"net/http/httptest"
// 	"strings"
// 	"testing"
//
// 	"github.com/koubae/game-hangar/internal/identity/auth"
// 	"github.com/koubae/game-hangar/internal/testunit"
// 	"github.com/stretchr/testify/assert"
// 	"github.com/stretchr/testify/require"
// )
//
// func TestAccountManagementController_Me(t *testing.T) {
// 	t.Parallel()
//
// 	tests := map[string]struct{
//
// 	}{
// 		"account-me": {},
// 	}
// 	for id, tt := range tests {
// 		t.Run(
// 			id, func(t *testing.T) {
// 				t.Parallel()
//
// 				_, handlerPtr, mocker := testunit.NewTestRouterAndContainer(t)
// 				handler := *handlerPtr
//
// 				if tt.expectedErrResponse == "" {
// 					var response auth.DTOAccessToken
// 					err = json.Unmarshal([]byte(rr.Body.String()), &response)
// 					require.NoError(t, err)
//
// 					assert.IsType(t, auth.DTOAccessToken{}, response)
// 					assert.NotEmpty(t, response.AccessToken)
// 					assert.NotEmpty(t, response.ExpiresIn)
// 					assert.Equal(t, tt.expectedCode, rr.Code)
// 				} else {
// 					response := rr.Body.String()
// 					assert.Equal(t, tt.expectedErrResponse, strings.TrimSpace(response))
// 					assert.Equal(t, tt.expectedCode, rr.Code)
// 				}
//
// 			},
// 		)
// 	}
// }
