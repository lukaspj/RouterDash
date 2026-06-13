package main

import (
	"flag"
	"log"
	"net/http"
	"strings"

	"routerdash/internal/config"
	"routerdash/internal/handler"
	"routerdash/internal/routeros"
)

func main() {
	cfg := config.Load()

	flag.StringVar(&cfg.RouterOSAddress, "address", cfg.RouterOSAddress, "RouterOS address (host:port)")
	flag.StringVar(&cfg.RouterOSUser, "user", cfg.RouterOSUser, "RouterOS username")
	flag.StringVar(&cfg.RouterOSPass, "pass", cfg.RouterOSPass, "RouterOS password")
	flag.StringVar(&cfg.ListenAddr, "listen", cfg.ListenAddr, "Listen address")
	flag.BoolVar(&cfg.InsecureSkipVerify, "insecure", cfg.InsecureSkipVerify, "Skip TLS certificate verification")
	flag.Parse()

	if cfg.RouterOSAddress == "" {
		log.Fatal("ROUTEROS_ADDRESS is required (set via env var or -address flag)")
	}

	client := routeros.NewClient(cfg.RouterOSAddress, cfg.RouterOSUser, cfg.RouterOSPass, cfg.InsecureSkipVerify)
	h := handler.New(client)

	srv := &server{
		handler:    h,
		staticDir:  "static",
		apiHandlers: map[string]http.HandlerFunc{
			"/api/system/resource":       h.HandleSystemResource,
			"/api/system/identity":       h.HandleSystemIdentity,
			"/api/interface":             h.HandleInterfaces,
			"/api/ip/address":            h.HandleIPAddresses,
			"/api/ip/route":              h.HandleIPRoutes,
			"/api/ip/dhcp-server/lease":  h.HandleDHCPLeases,
			"/api/ip/firewall/filter":    h.HandleFirewallFilter,
			"/api/ip/firewall/nat":       h.HandleFirewallNat,
			"/api/ip/firewall/mangle":    h.HandleFirewallMangle,
			"/api/ip/firewall/raw":       h.HandleFirewallRaw,
			"/api/ip/firewall/connection": h.HandleFirewallConnections,
			"/api/ip/firewall/connection/tracking": h.HandleConnectionTracking,
			"/api/system/resource/cpu":     h.HandleSystemResourceCPU,
			"/api/interface/wireless":    h.HandleWirelessInterfaces,
			"/api/interface/bridge/port": h.HandleBridgePorts,
			"/api/interface/ethernet":   h.HandleEthernetInterfaces,
			"/api/system/health":        h.HandleSystemHealth,
			"/api/system/routerboard":   h.HandleSystemRouterboard,
		},
	}

	log.Printf("routerdash starting on %s", cfg.ListenAddr)
	log.Fatal(http.ListenAndServe(cfg.ListenAddr, srv))
}

type server struct {
	handler     *handler.Handler
	staticDir   string
	apiHandlers map[string]http.HandlerFunc
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	if strings.HasPrefix(r.URL.Path, "/api/") {
		if fn, ok := s.apiHandlers[r.URL.Path]; ok {
			fn(w, r)
			return
		}
		http.NotFound(w, r)
		return
	}

	http.FileServer(http.Dir(s.staticDir)).ServeHTTP(w, r)
}
