// Package config loads PortView runtime configuration.
package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	envConfigPath = "PORTVIEW_CONFIG"
)

type Config struct {
	RefreshIntervalSeconds int    `json:"refresh_interval_seconds"`
	Theme                  string `json:"theme"`
	AccentColor            string `json:"accent_color"`
	HeaderForegroundColor  string `json:"header_foreground_color"`
	MutedColor             string `json:"muted_color"`
	StatusColor            string `json:"status_color"`
	ErrorColor             string `json:"error_color"`
}

func Default() Config {
	return Config{
		RefreshIntervalSeconds: 3,
		Theme:                  "purple",
	}
}

func Load() (Config, error) {
	path, err := Path()
	if err != nil {
		return Default(), err
	}

	cfg := Default()
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return cfg, nil
		}
		return cfg, fmt.Errorf("read config: %w", err)
	}

	if err := json.Unmarshal(data, &cfg); err != nil {
		return Default(), fmt.Errorf("parse config: %w", err)
	}

	cfg.ApplyDefaults()
	return cfg, nil
}

func Path() (string, error) {
	if path := strings.TrimSpace(os.Getenv(envConfigPath)); path != "" {
		return path, nil
	}

	dir, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("resolve config directory: %w", err)
	}
	return filepath.Join(dir, "portview", "config.json"), nil
}

func (c *Config) ApplyDefaults() {
	defaults := Default()
	if c.RefreshIntervalSeconds <= 0 {
		c.RefreshIntervalSeconds = defaults.RefreshIntervalSeconds
	}
	if strings.TrimSpace(c.Theme) == "" {
		c.Theme = defaults.Theme
	}
}
