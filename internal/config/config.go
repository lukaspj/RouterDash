package config

import (
	"os"
	"strings"
)

type Config struct {
	RouterOSAddress    string
	RouterOSUser       string
	RouterOSPass       string
	ListenAddr         string
	InsecureSkipVerify bool
}

func Load() *Config {
	cfg := &Config{
		RouterOSAddress:    getEnv("ROUTEROS_ADDRESS", ""),
		RouterOSUser:       getEnv("ROUTEROS_USER", "admin"),
		RouterOSPass:       getEnv("ROUTEROS_PASS", ""),
		ListenAddr:         getEnv("LISTEN_ADDR", ":8080"),
		InsecureSkipVerify: getEnv("ROUTEROS_INSECURE_SKIP_VERIFY", "true") == "true",
	}

	if cfg.RouterOSAddress != "" && !strings.HasPrefix(cfg.RouterOSAddress, "http") {
		cfg.RouterOSAddress = "https://" + cfg.RouterOSAddress
	}

	cfg.RouterOSAddress = strings.TrimRight(cfg.RouterOSAddress, "/")

	return cfg
}

func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}
