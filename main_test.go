package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestRouting verifies that the ServeHTTP multiplexer correctly
// dispatches /api/ paths to their registered handlers and everything
// else to the static file server.
func TestRouting(t *testing.T) {
	handlerCalls := make(map[string]int)
	named := func(name string) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			handlerCalls[name]++
			w.WriteHeader(200)
		}
	}

	srv := &server{
		staticDir: "static",
		apiHandlers: map[string]http.HandlerFunc{
			"/api/system/resource":         named("resource"),
			"/api/system/identity":         named("identity"),
			"/api/system/health":           named("health"),
			"/api/system/routerboard":      named("routerboard"),
			"/api/interface":               named("interface"),
			"/api/interface/wireless":      named("wireless"),
			"/api/interface/ethernet":      named("ethernet"),
			"/api/interface/bridge/port":   named("bridge-port"),
			"/api/ip/address":              named("address"),
			"/api/ip/route":                named("route"),
			"/api/ip/dhcp-server/lease":    named("dhcp-lease"),
			"/api/ip/firewall/filter":      named("firewall"),
		},
	}

	tests := []struct {
		path       string
		wantStatus int
		wantAPI    string // empty means static file, non-empty is handler name
	}{
		// Static files
		{"/", 200, ""},
		{"/css/style.css", 200, ""},
		// System endpoints
		{"/api/system/resource", 200, "resource"},
		{"/api/system/identity", 200, "identity"},
		{"/api/system/health", 200, "health"},
		{"/api/system/routerboard", 200, "routerboard"},
		// Interface endpoints
		{"/api/interface", 200, "interface"},
		{"/api/interface/wireless", 200, "wireless"},
		{"/api/interface/ethernet", 200, "ethernet"},
		{"/api/interface/bridge/port", 200, "bridge-port"},
		// IP endpoints
		{"/api/ip/address", 200, "address"},
		{"/api/ip/route", 200, "route"},
		{"/api/ip/dhcp-server/lease", 200, "dhcp-lease"},
		{"/api/ip/firewall/filter", 200, "firewall"},
		// Unknown paths
		{"/api/nonexistent", 404, ""},
		{"/nope/still/not/a/thing", 404, ""},
	}

	for _, tt := range tests {
		req := httptest.NewRequest("GET", tt.path, nil)
		rec := httptest.NewRecorder()
		handlerCalls = make(map[string]int) // reset between tests
		srv.ServeHTTP(rec, req)

		if rec.Code != tt.wantStatus {
			t.Errorf("%s: got status %d, want %d", tt.path, rec.Code, tt.wantStatus)
		}
		if tt.wantAPI != "" {
			if handlerCalls[tt.wantAPI] != 1 {
				t.Errorf("%s: expected handler %q to be called once, saw %d", tt.path, tt.wantAPI, handlerCalls[tt.wantAPI])
			}
		}
	}
}
