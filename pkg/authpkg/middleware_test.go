package authpkg

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/koubae/game-hangar/pkg/common"
	"github.com/stretchr/testify/assert"
)

func TestProtected(t *testing.T) {
	tests := map[string]struct {
		permissions  Permissions
		wantStatus   int
		wantNextCall bool
	}{
		"denied_without_permissions": {
			permissions:  nil,
			wantStatus:   http.StatusForbidden,
			wantNextCall: false,
		},
		"allowed_with_permission": {
			permissions: Permissions{
				common.AppID: {
					"account": {
						READ,
					},
				},
			},
			wantStatus:   http.StatusOK,
			wantNextCall: true,
		},
	}

	for name, tc := range tests {
		t.Run(
			name, func(t *testing.T) {
				nextCalled := false

				next := func(w http.ResponseWriter, r *http.Request) {
					nextCalled = true
					w.WriteHeader(http.StatusOK)
				}

				req := httptest.NewRequest(http.MethodGet, "/x", nil)
				rr := httptest.NewRecorder()

				if tc.permissions != nil {
					req = req.WithContext(WithPermissions(req.Context(), tc.permissions))
				}

				handler := Protected("account", READ, next)
				handler(rr, req)

				assert.Equal(t, tc.wantNextCall, nextCalled)
				assert.Equal(t, tc.wantStatus, rr.Code)

				if !tc.wantNextCall {
					assert.NotEmpty(t, rr.Body.String())
				}
			},
		)
	}
}
