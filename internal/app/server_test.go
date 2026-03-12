package app

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestAppInitialization(t *testing.T) {
	app := NewApp()

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
	app := NewApp()

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
	SetHandlerFunc     func(h http.Handler)
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

func (m *MockServer) SetHandler(h http.Handler) {
	if m.SetHandlerFunc != nil {
		m.SetHandlerFunc(h)
	}
}

func TestAppStopError(t *testing.T) {
	app := NewApp()
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
