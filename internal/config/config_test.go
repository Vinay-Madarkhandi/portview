package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadDefaultWhenConfigMissing(t *testing.T) {
	t.Setenv(envConfigPath, filepath.Join(t.TempDir(), "missing.json"))

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if cfg.RefreshIntervalSeconds != 3 {
		t.Fatalf("refresh interval = %d, want 3", cfg.RefreshIntervalSeconds)
	}
	if cfg.Theme != "purple" {
		t.Fatalf("theme = %q, want purple", cfg.Theme)
	}
}

func TestLoadConfig(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.json")
	t.Setenv(envConfigPath, path)

	data := []byte(`{
		"refresh_interval_seconds": 7,
		"theme": "green",
		"accent_color": "#00ff00"
	}`)
	if err := os.WriteFile(path, data, 0o600); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if cfg.RefreshIntervalSeconds != 7 {
		t.Fatalf("refresh interval = %d, want 7", cfg.RefreshIntervalSeconds)
	}
	if cfg.Theme != "green" {
		t.Fatalf("theme = %q, want green", cfg.Theme)
	}
	if cfg.AccentColor != "#00ff00" {
		t.Fatalf("accent = %q, want #00ff00", cfg.AccentColor)
	}
}

func TestLoadConfigAppliesDefaults(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.json")
	t.Setenv(envConfigPath, path)

	if err := os.WriteFile(path, []byte(`{}`), 0o600); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if cfg.RefreshIntervalSeconds != 3 {
		t.Fatalf("refresh interval = %d, want 3", cfg.RefreshIntervalSeconds)
	}
	if cfg.Theme != "purple" {
		t.Fatalf("theme = %q, want purple", cfg.Theme)
	}
}

func TestLoadInvalidConfig(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.json")
	t.Setenv(envConfigPath, path)

	if err := os.WriteFile(path, []byte(`{`), 0o600); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	if _, err := Load(); err == nil {
		t.Fatal("expected parse error")
	}
}
