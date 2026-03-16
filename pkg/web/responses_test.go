package web

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/koubae/game-hangar/pkg/common"
)

func TestWriteJSONErrorResponse(t *testing.T) {
	tests := []struct {
		name     string
		code     int
		message  string
		expected ResponseError
	}{
		{
			name:    "BadRequest error",
			code:    http.StatusBadRequest,
			message: "invalid input",
			expected: ResponseError{
				Code:    http.StatusBadRequest,
				Message: "invalid input",
			},
		},
		{
			name:    "NotFound error",
			code:    http.StatusNotFound,
			message: "resource not found",
			expected: ResponseError{
				Code:    http.StatusNotFound,
				Message: "resource not found",
			},
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				w := httptest.NewRecorder()
				WriteJSONErrorResponse(w, tt.code, tt.message)

				res := w.Result()
				defer res.Body.Close()

				if res.StatusCode != tt.code {
					t.Errorf("expected status code %d, got %d", tt.code, res.StatusCode)
				}

				if contentType := res.Header.Get("Content-Type"); contentType != "application/json" {
					t.Errorf("expected Content-Type application/json, got %s", contentType)
				}

				var actual ResponseError
				if err := json.NewDecoder(res.Body).Decode(&actual); err != nil {
					t.Fatalf("failed to decode response body: %v", err)
				}

				if actual != tt.expected {
					t.Errorf("expected %+v, got %+v", tt.expected, actual)
				}
			},
		)
	}
}

func TestWriteBusinessErrorResponse(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected ResponseError
	}{
		{
			name: "BusinessError match",
			err: &common.BusinessError{
				HTTPCode: http.StatusForbidden,
				Message:  "access denied",
			},
			expected: ResponseError{
				Code:    http.StatusForbidden,
				Message: "access denied",
			},
		},
		{
			name: "Generic error",
			err:  errors.New("some internal error"),
			expected: ResponseError{
				Code:    http.StatusInternalServerError,
				Message: "unexpected error",
			},
		},
		{
			name: "Nil error (generic path)",
			err:  nil,
			expected: ResponseError{
				Code:    http.StatusInternalServerError,
				Message: "unexpected error",
			},
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				w := httptest.NewRecorder()
				WriteBusinessErrorResponse(w, tt.err)

				res := w.Result()
				defer res.Body.Close()

				if res.StatusCode != tt.expected.Code {
					t.Errorf("expected status code %d, got %d", tt.expected.Code, res.StatusCode)
				}

				if contentType := res.Header.Get("Content-Type"); contentType != "application/json" {
					t.Errorf("expected Content-Type application/json, got %s", contentType)
				}

				var actual ResponseError
				if err := json.NewDecoder(res.Body).Decode(&actual); err != nil {
					t.Fatalf("failed to decode response body: %v", err)
				}

				if actual != tt.expected {
					t.Errorf("expected %+v, got %+v", tt.expected, actual)
				}
			},
		)
	}
}
