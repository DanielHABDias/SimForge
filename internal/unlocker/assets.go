package unlocker

import (
	"os"
	"path/filepath"
)

// Assets represents the source files needed for installation.
type Assets struct {
	// RootDir is the directory containing the "ea_app", "origin" sub-folders
	// and game config (g_*.ini) files. Normally the same directory as the
	// original unlocker extraction.
	RootDir string

	DllEA     string // path to ea_app/version.dll
	DllOrigin string // path to origin/version.dll
	Config    string // path to config.ini
}

// SearchCandidates are the directories (relative to the working directory or
// the executable) that may contain the unlocker assets.
var SearchCandidates = []string{
	".",
	"assets",
	"assets/scripts",
	"assets/scripts/EA DLC Unlocker v2",
}

// LocateAssets tries to find the unlocker assets relative to the
// current working directory or the executable path.
func LocateAssets() (*Assets, error) {
	// Try current working directory.
	dir, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	if assets := searchCandidates(dir); assets != nil {
		return assets, nil
	}

	// Try executable directory.
	exec, err := os.Executable()
	if err == nil {
		if assets := searchCandidates(filepath.Dir(exec)); assets != nil {
			return assets, nil
		}
	}

	return nil, ErrAssetsNotFound
}

func searchCandidates(base string) *Assets {
	for _, rel := range SearchCandidates {
		if assets := searchDir(filepath.Join(base, rel)); assets != nil {
			return assets
		}
	}
	return nil
}

func searchDir(dir string) *Assets {
	// Look for the typical unlocker layout: config.ini, ea_app/, origin/
	configPath := filepath.Join(dir, "config.ini")
	eaAppDll := filepath.Join(dir, "ea_app", "version.dll")
	originDll := filepath.Join(dir, "origin", "version.dll")

	if fileExists(configPath) && fileExists(eaAppDll) && fileExists(originDll) {
		return &Assets{
			RootDir:   dir,
			DllEA:     eaAppDll,
			DllOrigin: originDll,
			Config:    configPath,
		}
	}

	return nil
}

// ListGameConfigs returns all g_*.ini files in the assets directory.
func (a *Assets) ListGameConfigs() ([]GameConfig, error) {
	pattern := filepath.Join(a.RootDir, "g_*.ini")
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return nil, err
	}

	var configs []GameConfig
	for _, m := range matches {
		name := filepath.Base(m)
		// Strip "g_" prefix and ".ini" suffix
		if len(name) > 6 && name[:2] == "g_" {
			gameName := name[2 : len(name)-4]
			configs = append(configs, GameConfig{
				Name: gameName,
				Path: m,
			})
		}
	}
	return configs, nil
}

// GameConfig represents a single game configuration file.
type GameConfig struct {
	Name string
	Path string
}

func fileExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}

