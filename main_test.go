package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRouting(t *testing.T) {
	// We don't need real handler — test just the ServeHTTP routing.
	// Use a simple handler that records it was called.
	apiCalled := false
	srv := &server{
		staticDir: "static",
		apiHandlers: map[string]http.HandlerFunc{
			"/api/system/resource": func(w http.ResponseWriter, r *http.Request) {
				apiCalled = true
				w.WriteHeader(200)
			},
		},
	}

	tests := []struct {
		path       string
		wantStatus int
		wantAPI    bool
	}{
		{"/", 200, false},
		{"/css/style.css", 200, false},
		{"/api/system/resource", 200, true},
		{"/api/nonexistent", 404, false},
	}

	for _, tt := range tests {
		apiCalled = false
		req := httptest.NewRequest("GET", tt.path, nil)
		rec := httptest.NewRecorder()
		srv.ServeHTTP(rec, req)

		if rec.Code != tt.wantStatus {
			t.Errorf("%s: got status %d, want %d", tt.path, rec.Code, tt.wantStatus)
		}
		if tt.wantAPI && !apiCalled {
			t.Errorf("%s: expected API handler to be called", tt.path)
		}
		if !tt.wantAPI && apiCalled {
			t.Errorf("%s: expected static file handler, got API handler", tt.path)
		}
	}
}
