package app

import (
	"net/http"
	"net/http/httptest"
	"testing"
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
	app.Server.Handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status OK, got %v", rr.Code)
	}
}
