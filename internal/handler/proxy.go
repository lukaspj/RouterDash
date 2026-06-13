package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"routerdash/internal/routeros"
)

type Handler struct {
	client *routeros.Client
}

func New(client *routeros.Client) *Handler {
	return &Handler{client: client}
}

func (h *Handler) HandleSystemResource(w http.ResponseWriter, r *http.Request) {
	h.proxy(w, "system/resource")
}

func (h *Handler) HandleSystemIdentity(w http.ResponseWriter, r *http.Request) {
	h.proxy(w, "system/identity")
}

func (h *Handler) HandleInterfaces(w http.ResponseWriter, r *http.Request) {
	h.proxy(w, "interface")
}

func (h *Handler) HandleIPAddresses(w http.ResponseWriter, r *http.Request) {
	h.proxy(w, "ip/address")
}

func (h *Handler) HandleIPRoutes(w http.ResponseWriter, r *http.Request) {
	h.proxy(w, "ip/route")
}

func (h *Handler) HandleDHCPLeases(w http.ResponseWriter, r *http.Request) {
	h.proxy(w, "ip/dhcp-server/lease")
}

func (h *Handler) HandleFirewallFilter(w http.ResponseWriter, r *http.Request) {
	h.proxy(w, "ip/firewall/filter")
}

func (h *Handler) HandleFirewallNat(w http.ResponseWriter, r *http.Request) {
	h.proxy(w, "ip/firewall/nat")
}

func (h *Handler) HandleFirewallMangle(w http.ResponseWriter, r *http.Request) {
	h.proxy(w, "ip/firewall/mangle")
}

func (h *Handler) HandleFirewallRaw(w http.ResponseWriter, r *http.Request) {
	h.proxy(w, "ip/firewall/raw")
}

func (h *Handler) HandleWirelessInterfaces(w http.ResponseWriter, r *http.Request) {
	h.proxy(w, "interface/wireless")
}

func (h *Handler) HandleBridgePorts(w http.ResponseWriter, r *http.Request) {
	h.proxy(w, "interface/bridge/port")
}

func (h *Handler) HandleEthernetInterfaces(w http.ResponseWriter, r *http.Request) {
	h.proxy(w, "interface/ethernet")
}

func (h *Handler) HandleSystemHealth(w http.ResponseWriter, r *http.Request) {
	h.proxy(w, "system/health")
}

func (h *Handler) HandleSystemRouterboard(w http.ResponseWriter, r *http.Request) {
	h.proxy(w, "system/routerboard")
}

func (h *Handler) proxy(w http.ResponseWriter, restPath string) {
	data, err := h.client.Get(restPath)
	if err != nil {
		log.Printf("proxy error [%s]: %v", restPath, err)
		if strings.Contains(err.Error(), "no such command") || strings.Contains(err.Error(), "no such item") {
			writeJSON(w, http.StatusOK, []byte("[]"))
			return
		}
		writeError(w, http.StatusBadGateway, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, data)
}

func writeError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

func writeJSON(w http.ResponseWriter, status int, data []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(data)
}
