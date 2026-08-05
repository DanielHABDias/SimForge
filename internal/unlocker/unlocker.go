// Package unlocker implements the operations of the EA DLC Unlocker
// (install, uninstall, add/update game config) in a cross-platform way.
package unlocker

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/DanielHABDias/SimForge/internal/finder"
)

// ErrAssetsNotFound is returned when the unlocker source files could not be found.
var ErrAssetsNotFound = errors.New("could not locate the unlocker assets (config.ini, ea_app/, origin/)")

// Manager performs unlocker operations against a single installation.
type Manager struct {
	Installation finder.Installation
	Assets       *Assets
	// DllOverride is platform-specific (Linux writes to user.reg).
	DllOverride DllOverrider
}

// New creates a Manager for the given installation.
func New(inst finder.Installation) (*Manager, error) {
	assets, err := LocateAssets()
	if err != nil {
		return nil, err
	}
	return &Manager{
		Installation: inst,
		Assets:       assets,
		DllOverride:  newDllOverrider(inst),
	}, nil
}

// SourceDll returns the path of the version.dll to copy for the client.
func (m *Manager) SourceDll() (string, error) {
	if m.Installation.Client == "origin" {
		if !fileExists(m.Assets.DllOrigin) {
			return "", fmt.Errorf("%s missing, you didn't extract all files", m.Assets.DllOrigin)
		}
		return m.Assets.DllOrigin, nil
	}
	if !fileExists(m.Assets.DllEA) {
		return "", fmt.Errorf("%s missing, you didn't extract all files", m.Assets.DllEA)
	}
	return m.Assets.DllEA, nil
}

// IsInstalled reports whether the unlocker is present in the client.
func (m *Manager) IsInstalled() bool {
	return fileExists(m.Installation.DllPath) &&
		fileExists(filepath.Join(m.Installation.ConfigPath, "config.ini"))
}

// Install installs the unlocker: copies the DLL to the client (and staged
// directory for EA app) and copies the main config.
func (m *Manager) Install() error {
	srcDll, err := m.SourceDll()
	if err != nil {
		return err
	}
	if !fileExists(m.Assets.Config) {
		return fmt.Errorf("%s missing, you didn't extract all files", m.Assets.Config)
	}

	if err := createDir(m.Installation.ConfigPath); err != nil {
		return fmt.Errorf("could not create the configs folder: %w", err)
	}

	// Copy main config.
	if err := copyFile(m.Assets.Config, filepath.Join(m.Installation.ConfigPath, "config.ini")); err != nil {
		return fmt.Errorf("could not copy the main config: %w", err)
	}

	// Register the DLL override (platform-specific).
	if m.DllOverride != nil {
		if err := m.DllOverride.Apply(); err != nil {
			return err
		}
	}

	// Copy the DLL to the client folder.
	if err := copyFile(srcDll, m.Installation.DllPath); err != nil {
		return fmt.Errorf("could not install the Unlocker: %w", err)
	}

	// Copy to the staged directory (EA app only).
	if m.Installation.DllPath2 != "" {
		_ = createDir(filepath.Dir(m.Installation.DllPath2))
		if err := copyFile(srcDll, m.Installation.DllPath2); err != nil {
			return fmt.Errorf("could not install the Unlocker to staged folder: %w", err)
		}
	}

	return nil
}

// Uninstall removes the unlocker: DLLs, config folder, logs folder, and the
// DLL override registration.
func (m *Manager) Uninstall() error {
	// Register removal of the DLL override (platform-specific).
	if m.DllOverride != nil {
		if err := m.DllOverride.Remove(); err != nil {
			return err
		}
	}

	// Remove DLLs.
	if err := removeFile(m.Installation.DllPath); err != nil {
		return err
	}
	if m.Installation.DllPath2 != "" {
		_ = removeFile(m.Installation.DllPath2)
	}

	// Remove config folder.
	if err := removeDirRecursive(m.Installation.ConfigPath); err != nil {
		return err
	}
	removeDirIfEmpty(filepath.Dir(m.Installation.ConfigPath))

	// Remove logs folder.
	if err := removeDirRecursive(m.Installation.LogsPath); err != nil {
		return err
	}
	removeDirIfEmpty(filepath.Dir(m.Installation.LogsPath))

	return nil
}

// AddGameConfig copies a game config (g_<name>.ini) into the config folder
// and deletes the related .etag file.
func (m *Manager) AddGameConfig(name string) error {
	config, err := m.Assets.ListGameConfigs()
	if err != nil {
		return err
	}

	for _, c := range config {
		if c.Name == name {
			if err := createDir(m.Installation.ConfigPath); err != nil {
				return fmt.Errorf("could not create the configs folder: %w", err)
			}
			if err := copyFile(c.Path, filepath.Join(m.Installation.ConfigPath, c.Name+".ini")); err != nil {
				return fmt.Errorf("could not copy the game config: %w", err)
			}
			// Remove etag to force refresh.
			_ = removeFile(filepath.Join(m.Installation.LogsPath, name+".etag"))
			return nil
		}
	}
	return fmt.Errorf("game config %q not found", name)
}

// InstalledGameConfigs returns the names of game configs currently installed.
func (m *Manager) InstalledGameConfigs() []string {
	matches, _ := filepath.Glob(filepath.Join(m.Installation.ConfigPath, "g_*.ini"))
	var names []string
	for _, m := range matches {
		base := filepath.Base(m)
		if len(base) > 6 && base[:2] == "g_" {
			names = append(names, base[2:len(base)-4])
		}
	}
	return names
}

// OpenConfigsFolder opens the config folder in the file manager.
func (m *Manager) OpenConfigsFolder() error {
	if !directoryExists(m.Installation.ConfigPath) {
		return errors.New("configs folder not found. Install the Unlocker first")
	}
	return openFolder(m.Installation.ConfigPath)
}

// OpenLogsFolder opens the logs folder in the file manager.
func (m *Manager) OpenLogsFolder() error {
	if !directoryExists(m.Installation.LogsPath) {
		return errors.New("logs folder not found. Install the Unlocker and run EA app/Origin first")
	}
	return openFolder(m.Installation.LogsPath)
}

// GameConfigNames lists the available game configs from the assets.
func (m *Manager) GameConfigNames() ([]string, error) {
	configs, err := m.Assets.ListGameConfigs()
	if err != nil {
		return nil, err
	}
	names := make([]string, 0, len(configs))
	for _, c := range configs {
		names = append(names, c.Name)
	}
	return names, nil
}

// helper definitions

func directoryExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}
