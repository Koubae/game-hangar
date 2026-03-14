package web

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/koubae/game-hangar/pkg/common"
	"github.com/rs/cors"
	"go.uber.org/zap"
)

func testRouter(common.Logger, *common.Config) *http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc(
		"/", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		},
	)

	return new(cors.New(cors.Options{}).Handler(mux))
}

func TestAppInitialization(t *testing.T) {
	app := NewApp("IDENTITY_", testRouter)

	if app.Config == nil {
		t.Fatal("Config should not be nil")
	}
	if app.Logger == nil {
		t.Fatal("Logger should not be nil")
	}
	if app.Server == nil {
		t.Fatal("Server should not be nil")
	}

	// Verify the router is working
	req := httptest.NewRequest("GET", "/", nil)
	rr := httptest.NewRecorder()
	app.Server.Handler().ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status OK, got %v", rr.Code)
	}
}

func TestAppStartStop(t *testing.T) {
	app := NewApp("", testRouter)

	// Use a random port to avoid conflicts
	app.Config.Port = 0
	// We need a way to set Addr if we use the wrapper, or just use the concrete server for this test
	if srv, ok := app.Server.(*httpServerWrapper); ok {
		srv.Server.Addr = ":0"
	}

	ctx, cancel := context.WithCancel(t.Context())

	// Start the app in a goroutine
	done := make(chan struct{})
	go func() {
		app.Start(ctx)
		close(done)
	}()

	// Give it a moment to start
	time.Sleep(100 * time.Millisecond)

	// Trigger shutdown
	cancel()

	// Wait for Start to return
	select {
	case <-done:
		// Success
	case <-time.After(2 * time.Second):
		t.Fatal("App.Start did not return after context cancellation")
	}

	// Test Stop
	if err := app.Stop(); err != nil {
		t.Errorf("app.Stop() failed: %v", err)
	}
}

type MockServer struct {
	ListenAndServeFunc func() error
	ShutdownFunc       func(ctx context.Context) error
	HandlerFunc        func() http.Handler
}

func (m *MockServer) ListenAndServe() error {
	if m.ListenAndServeFunc != nil {
		return m.ListenAndServeFunc()
	}
	return nil
}

func (m *MockServer) Shutdown(ctx context.Context) error {
	if m.ShutdownFunc != nil {
		return m.ShutdownFunc(ctx)
	}
	return nil
}

func (m *MockServer) Handler() http.Handler {
	if m.HandlerFunc != nil {
		return m.HandlerFunc()
	}
	return nil
}

type MockLogger struct {
	FatalCalled bool
	FatalMsg    string
	FatalFields []zap.Field
}

func (m *MockLogger) Debug(msg string, fields ...zap.Field) {}
func (m *MockLogger) Info(msg string, fields ...zap.Field)  {}
func (m *MockLogger) Warn(msg string, fields ...zap.Field)  {}
func (m *MockLogger) Error(msg string, fields ...zap.Field) {}
func (m *MockLogger) Panic(msg string, fields ...zap.Field) {}
func (m *MockLogger) Fatal(msg string, fields ...zap.Field) {
	m.FatalCalled = true
	m.FatalMsg = msg
	m.FatalFields = fields
}

func TestAppStopError(t *testing.T) {
	app := NewApp("", testRouter)
	expectedErr := errors.New("shutdown error")
	app.Server = &MockServer{
		ShutdownFunc: func(ctx context.Context) error {
			return expectedErr
		},
	}

	err := app.Stop()
	if !errors.Is(err, expectedErr) {
		t.Errorf("expected error %v, got %v", expectedErr, err)
	}
}

func TestAppStartError(t *testing.T) {
	app := NewApp("", testRouter)
	expectedErr := errors.New("listen and serve error")

	mockLogger := &MockLogger{}
	app.Logger = mockLogger

	app.Server = &MockServer{
		ListenAndServeFunc: func() error {
			return expectedErr
		},
	}

	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()

	go app.Start(ctx)

	// Wait for Fatal to be called
	start := time.Now()
	for time.Since(start) < 2*time.Second {
		if mockLogger.FatalCalled {
			break
		}
		time.Sleep(50 * time.Millisecond)
	}

	if !mockLogger.FatalCalled {
		t.Fatal("timed out waiting for Fatal call from App.Start")
	}

	if mockLogger.FatalMsg != "Server failed to start" {
		t.Errorf("expected fatal msg %q, got %q", "Server failed to start", mockLogger.FatalMsg)
	}

	found := false
	for _, f := range mockLogger.FatalFields {
		if f.Key == "error" && f.Interface.(error) == expectedErr {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected error field with %v in fatal fields", expectedErr)
	}
}
