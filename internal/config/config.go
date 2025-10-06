package config

import (
	"os"
)

const (
	listenAddressEnv     = "LISTEN_ADDRESS"
	defaultListenAddress = ":8080"
)

// Config captures runtime configuration knobs for the application.
type Config struct {
	ListenAddress string
}

// Load resolves configuration from environment variables, falling back to defaults.
func Load() Config {
	listenAddr := lookupEnvDefault(listenAddressEnv, defaultListenAddress)

	return Config{
		ListenAddress: listenAddr,
	}
}

func lookupEnvDefault(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok && value != "" {
		return value
	}
	return fallback
}
