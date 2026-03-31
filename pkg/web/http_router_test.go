package web

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/koubae/game-hangar/pkg/common"
	"go.uber.org/zap"
)

type MockContainer struct {
	logger common.Logger
}

func (m *MockContainer) Logger() common.Logger {
	return m.logger
}

func (m *MockContainer) Shutdown() error {
	return nil
}

func TestRouterEndpoints(t *testing.T) {
	logger := &routerMockLogger{}
	config := &common.Config{
		AppName:    "TestApp",
		AppVersion: "1.0.0",
		CORSConfig: &common.CORSConfig{},
	}
	routerRegister := func(mux *http.ServeMux) {}

	handlerPtr := Router(&MockContainer{logger: logger}, config, routerRegister)
	handler := *handlerPtr

	tests := []struct {
		name           string
		path           string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "Root index",
			path:           "/",
			expectedStatus: http.StatusOK,
			expectedBody:   "Welcome to TestApp v1.0.0+",
		},
		{
			name:           "Healthz endpoint",
			path:           "/healthz",
			expectedStatus: http.StatusOK,
			expectedBody:   "OK",
		},
		{
			name:           "Readyz endpoint",
			path:           "/readyz",
			expectedStatus: http.StatusOK,
			expectedBody:   "OK",
		},
		{
			name:           "Ping endpoint",
			path:           "/ping",
			expectedStatus: http.StatusOK,
			expectedBody:   "pong",
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				req, err := http.NewRequest("GET", tt.path, nil)
				if err != nil {
					t.Fatal(err)
				}
				rr := httptest.NewRecorder()
				handler.ServeHTTP(rr, req)

				if rr.Code != tt.expectedStatus {
					t.Errorf(
						"handler returned wrong status code for %s: got %v want %v",
						tt.path, rr.Code, tt.expectedStatus,
					)
				}

				if tt.expectedBody != "" && rr.Body.String() != tt.expectedBody {
					t.Errorf(
						"handler returned unexpected body for %s: got %v want %v",
						tt.path, rr.Body.String(), tt.expectedBody,
					)
				}
			},
		)
	}
}

type routerMockLogger struct {
	common.Logger
}

func (m *routerMockLogger) Debug(msg string, fields ...zap.Field)  {}
func (m *routerMockLogger) Info(msg string, fields ...zap.Field)   {}
func (m *routerMockLogger) Warn(msg string, fields ...zap.Field)   {}
func (m *routerMockLogger) Error(msg string, fields ...zap.Field)  {}
func (m *routerMockLogger) Panic(msg string, fields ...zap.Field)  {}
func (m *routerMockLogger) DPanic(msg string, fields ...zap.Field) {}
func (m *routerMockLogger) Fatal(msg string, fields ...zap.Field)  {}
