// Package config loads the host/daemon configuration from a JSON file in the
// user's config directory. Nothing is trusted by default: with no config the
// host exposes no filesystem roots.
package config

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

// Config is the on-disk configuration.
type Config struct {
	// AllowedRoots are the only directories fs.* methods may touch.
	AllowedRoots []string `json:"allowedRoots"`
	// Sidecar is the path to the Rust sidecar binary (optional).
	Sidecar string `json:"sidecar"`
	// Daemon settings (toolboxd).
	Daemon Daemon `json:"daemon"`
}

// Daemon holds the local HTTP daemon settings.
type Daemon struct {
	Addr           string   `json:"addr"`           // default 127.0.0.1:48731
	Token          string   `json:"token"`          // required: sent as x-api-key
	AllowedOrigins []string `json:"allowedOrigins"` // e.g. chrome-extension://<id>
}

// DefaultAddr is the loopback address used when none is configured.
const DefaultAddr = "127.0.0.1:48731"

// Dir returns the per-user config directory for the toolbox.
func Dir() (string, error) {
	base, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(base, "toolbox"), nil
}

// DefaultPath is <config dir>/config.json.
func DefaultPath() (string, error) {
	d, err := Dir()
	if err != nil {
		return "", err
	}
	return filepath.Join(d, "config.json"), nil
}

// Load reads path. A missing file yields an empty (deny-all) config.
func Load(path string) (Config, error) {
	var c Config
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		c.Daemon.Addr = DefaultAddr
		return c, nil
	}
	if err != nil {
		return c, err
	}
	if err := json.Unmarshal(data, &c); err != nil {
		return c, err
	}
	if c.Daemon.Addr == "" {
		c.Daemon.Addr = DefaultAddr
	}
	return c, nil
}
