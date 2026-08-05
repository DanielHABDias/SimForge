//go:build linux

package finder

import (
	"bufio"
	"os"
	"path/filepath"
	"regexp"
)

func findSteam() []Installation {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil
	}

	steamPaths := []string{
		filepath.Join(home, ".steam", "steam"),
		filepath.Join(home, ".local", "share", "Steam"),
		filepath.Join(home, "snap", "steam", "common", ".local", "share", "Steam"),
		filepath.Join(home, "steam"),
		filepath.Join(home, ".var", "app", "com.valvesoftware.Steam", ".steam", "steam"),
		filepath.Join(home, ".var", "app", "com.valvesoftware.Steam", ".local", "share", "Steam"),
	}

	var installations []Installation

	for _, steamPath := range steamPaths {
		steamApps := filepath.Join(steamPath, "steamapps")
		if !directoryExists(steamApps) {
			continue
		}

		for _, library := range findSteamLibraries(steamApps) {
			installations = append(installations, findProtonPrefixes(library)...)
		}
	}

	return installations
}

func findSteamLibraries(steamApps string) []string {
	libraries := []string{steamApps}

	file := filepath.Join(steamApps, "libraryfolders.vdf")
	f, err := os.Open(file)
	if err != nil {
		return libraries
	}
	defer f.Close()

	regex := regexp.MustCompile(`"path"\s+"([^"]+)"`)
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		match := regex.FindStringSubmatch(scanner.Text())
		if len(match) > 1 {
			librarySteamApps := filepath.Join(match[1], "steamapps")
			if directoryExists(librarySteamApps) {
				libraries = append(libraries, librarySteamApps)
			}
		}
	}
	return libraries
}

func findProtonPrefixes(steamApps string) []Installation {
	compatData := filepath.Join(steamApps, "compatdata")
	if !directoryExists(compatData) {
		return nil
	}

	files, err := os.ReadDir(compatData)
	if err != nil {
		return nil
	}

	var installations []Installation
	for _, file := range files {
		if !file.IsDir() {
			continue
		}

		appID := file.Name()
		prefix := filepath.Join(compatData, appID, "pfx")
		if !directoryExists(prefix) {
			continue
		}

		installation := Installation{
			Platform:   "linux",
			Launcher:   "Steam Proton",
			Name:       GetGameName(appID),
			AppID:      appID,
			PrefixPath: prefix,
			User:       "steamuser",
		}

		// Prefer EA app, fall back to Origin.
		if ea := findEA(prefix, "ea_app"); ea != nil {
			applyEAInfo(&installation, ea)
		} else if ea := findEA(prefix, "origin"); ea != nil {
			applyEAInfo(&installation, ea)
		}

		// Only include prefixes that actually contain EA app/Origin.
		if installation.ClientPath != "" {
			installations = append(installations, installation)
		}
	}

	return installations
}

func applyEAInfo(installation *Installation, ea *eaInfo) {
	installation.Client = ea.Client
	installation.ClientPath = ea.ClientPath
	installation.DllPath = ea.DllPath
	installation.DllPath2 = ea.DllPath2
	installation.ConfigPath = ea.ConfigPath
	installation.LogsPath = ea.LogPath
	installation.User = ea.User
}
