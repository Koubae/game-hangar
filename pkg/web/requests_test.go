package web

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testPayload struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func TestLoadJsonBodyOrBadRequestResponse(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		body            string
		expectedOk      bool
		expectedPayload *testPayload
	}{
		"valid-json-body": {
			body:            `{"name":"alice","age":20}`,
			expectedOk:      true,
			expectedPayload: &testPayload{Name: "alice", Age: 20},
		},
		"empty-json-body": {
			body:            ``,
			expectedOk:      false,
			expectedPayload: nil,
		},

		"malformed-json-body": {
			body:            `{"name":"Alice","age":`,
			expectedOk:      false,
			expectedPayload: nil,
		},
		"incorrect-json-expected-structure": {
			body:            `{"name":"alice","age":"should-be-a-number"}`,
			expectedOk:      false,
			expectedPayload: nil,
		},
	}

	for id, tt := range tests {
		t.Run(
			id, func(t *testing.T) {
				t.Parallel()

				r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(tt.body))
				w := httptest.NewRecorder()

				payload, ok := LoadJsonBody[testPayload](w, r)

				assert.Equal(t, tt.expectedOk, ok)
				if tt.expectedOk {
					assert.NotNil(t, payload)
					assert.Equal(t, tt.expectedPayload, payload)
				} else {

					var resp ResponseError
					err := json.Unmarshal(w.Body.Bytes(), &resp)
					require.NoError(t, err, "failed to unmarshal response body")

					ct := w.Header().Get("Content-Type")

					assert.Nil(t, payload)
					assert.Equal(t, http.StatusBadRequest, w.Code)
					assert.Equal(t, "application/json", ct)
					assert.Equal(t, 400, resp.Code)
					assert.True(t, strings.HasPrefix(resp.Message, "invalid json:"))
				}

			},
		)
	}
}
