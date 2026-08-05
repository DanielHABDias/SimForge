package config

import (
	"os"
	"path/filepath"
	"testing"
)

// TestLoadCreatesDefault verifies that Load creates a config file with the
// default paths when none exists.
func TestLoadCreatesDefault(t *testing.T) {
	dir := t.TempDir()
	// Point configPath at our temp dir by changing working directory.
	oldWd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(oldWd)
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load returned error: %v", err)
	}

	if cfg.DLCPath != DefaultDLCPath {
		t.Errorf("DLCPath = %q, want %q", cfg.DLCPath, DefaultDLCPath)
	}
	if cfg.ModsPath != DefaultModsPath {
		t.Errorf("ModsPath = %q, want %q", cfg.ModsPath, DefaultModsPath)
	}

	// The config file should now exist.
	if _, err := os.Stat(FileName); err != nil {
		t.Errorf("config file was not created: %v", err)
	}
}

// TestLoadAppliesDefaults verifies that empty values in an existing config
// are replaced by defaults.
func TestLoadAppliesDefaults(t *testing.T) {
	dir := t.TempDir()
	oldWd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(oldWd)
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}

	// Write a config with one empty field.
	writeTestConfig(t, `{"dlc_path": "", "mods_path": "custom/mods"}`)

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load returned error: %v", err)
	}

	if cfg.DLCPath != DefaultDLCPath {
		t.Errorf("DLCPath = %q, want default %q", cfg.DLCPath, DefaultDLCPath)
	}
	if cfg.ModsPath != "custom/mods" {
		t.Errorf("ModsPath = %q, want %q", cfg.ModsPath, "custom/mods")
	}
}

// TestSavePersists verifies Save writes the values so a subsequent Load reads
// them back.
func TestSavePersists(t *testing.T) {
	dir := t.TempDir()
	oldWd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(oldWd)
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}

	cfg := &Config{DLCPath: "my/dlcs", ModsPath: "my/mods"}
	if err := cfg.Save(); err != nil {
		t.Fatalf("Save returned error: %v", err)
	}

	loaded, err := Load()
	if err != nil {
		t.Fatalf("Load returned error: %v", err)
	}
	if loaded.DLCPath != "my/dlcs" {
		t.Errorf("DLCPath = %q, want %q", loaded.DLCPath, "my/dlcs")
	}
	if loaded.ModsPath != "my/mods" {
		t.Errorf("ModsPath = %q, want %q", loaded.ModsPath, "my/mods")
	}
}

// TestInvalidConfig verifies an invalid config returns an error.
func TestInvalidConfig(t *testing.T) {
	dir := t.TempDir()
	oldWd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(oldWd)
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}

	writeTestConfig(t, `{invalid json`)

	if _, err := Load(); err == nil {
		t.Error("Load should return error for invalid config")
	}
}

func writeTestConfig(t *testing.T, content string) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(".", FileName), []byte(content), 0o644); err != nil {
		t.Fatalf("could not write test config: %v", err)
	}
}
