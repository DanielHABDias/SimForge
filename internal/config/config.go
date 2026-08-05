// Package config handles the SimForge configuration file.
// It stores the source paths for downloaded DLCs and mods.
package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

// Default paths used when the config does not specify a value.
const (
	DefaultDLCPath  = "assets/dlcs"
	DefaultModsPath = "assets/mods"
)

// Config represents the SimForge configuration.
type Config struct {
	// DLCPath is the folder where downloaded DLC files are stored.
	DLCPath string `json:"dlc_path"`
	// ModsPath is the folder where downloaded mod files are stored.
	ModsPath string `json:"mods_path"`
}

// FileName is the name of the configuration file.
const FileName = "config.json"

// Load reads the configuration from disk. If the file does not exist, it is
// created with the default values. Values are always applied with defaults
// when empty.
func Load() (*Config, error) {
	path := configPath()
	cfg := &Config{
		DLCPath:  DefaultDLCPath,
		ModsPath: DefaultModsPath,
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			// Create the default config on first run.
			if serr := cfg.Save(); serr != nil {
				return nil, fmt.Errorf("could not create default config: %w", serr)
			}
			return cfg, nil
		}
		return nil, fmt.Errorf("could not read config: %w", err)
	}

	if len(data) > 0 {
		var loaded Config
		if uerr := json.Unmarshal(data, &loaded); uerr != nil {
			return nil, fmt.Errorf("invalid config file %s: %w", path, uerr)
		}
		cfg = &loaded
	}

	// Apply defaults for any empty field.
	if cfg.DLCPath == "" {
		cfg.DLCPath = DefaultDLCPath
	}
	if cfg.ModsPath == "" {
		cfg.ModsPath = DefaultModsPath
	}

	return cfg, nil
}

// Save writes the configuration to disk.
func (c *Config) Save() error {
	path := configPath()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("could not create config folder: %w", err)
	}

	// Normalise empty values back to defaults before saving.
	if c.DLCPath == "" {
		c.DLCPath = DefaultDLCPath
	}
	if c.ModsPath == "" {
		c.ModsPath = DefaultModsPath
	}

	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return fmt.Errorf("could not encode config: %w", err)
	}

	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("could not write config: %w", err)
	}
	return nil
}

// Path returns the absolute location of the config file for display.
func Path() string {
	return configPath()
}

// configPath returns the config file path, preferring the current working
// directory and falling back to the executable directory.
func configPath() string {
	if dir, err := os.Getwd(); err == nil {
		return filepath.Join(dir, FileName)
	}
	if exec, err := os.Executable(); err == nil {
		return filepath.Join(filepath.Dir(exec), FileName)
	}
	return FileName
}
