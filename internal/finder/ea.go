package finder

import (
	"os"
	"path/filepath"
)

// findWineUser returns the first non-default user directory inside a prefix.
func findWineUser(prefixPath string) string {
	usersPath := filepath.Join(prefixPath, "drive_c", "users")

	users, err := os.ReadDir(usersPath)
	if err != nil {
		return "user"
	}

	for _, user := range users {
		if !user.IsDir() {
			continue
		}
		name := user.Name()
		// ignore default Wine users
		if name != "Public" && name != "Default" && name != "defaultuser0" {
			return name
		}
	}
	return "user"
}

// findVersionedDir looks for EA Desktop installations that live under a
// versioned sub-directory (e.g. ".../EA Desktop/13.735.2.6250/EA Desktop").
func findVersionedDir(base string) (string, bool) {
	entries, err := os.ReadDir(base)
	if err != nil {
		return "", false
	}

	// The versioned folder is usually a numeric directory containing "EA Desktop".
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		inner := filepath.Join(base, e.Name(), "EA Desktop")
		if directoryExists(inner) {
			return inner, true
		}
	}

	// Fallback: the base itself may be the client folder.
	if directoryExists(filepath.Join(base, "EA Desktop")) {
		return filepath.Join(base, "EA Desktop"), true
	}
	return "", false
}

// findEA locates the EA app (and Origin) inside a Wine/Proton prefix.
// clientType selects which client to look for: "ea_app" or "origin".
func findEA(prefixPath string, clientType string) *eaInfo {
	var possibleClients []string

	switch clientType {
	case "origin":
		possibleClients = []string{
			filepath.Join(prefixPath, "drive_c", "Program Files (x86)", "Origin"),
			filepath.Join(prefixPath, "drive_c", "Program Files", "Origin"),
		}
	default: // ea_app
		base := filepath.Join(prefixPath, "drive_c", "Program Files", "Electronic Arts", "EA Desktop")
		bases := []string{
			base,
			filepath.Join(prefixPath, "drive_c", "Program Files (x86)", "Electronic Arts", "EA Desktop"),
		}
		for _, b := range bases {
			if client, ok := findVersionedDir(b); ok {
				possibleClients = append(possibleClients, client)
			}
		}
	}

	var clientPath string
	for _, p := range possibleClients {
		if directoryExists(p) {
			clientPath = p
			break
		}
	}

	if clientPath == "" {
		return nil
	}

	user := findWineUser(prefixPath)

	configPath := filepath.Join(prefixPath, "drive_c", "users", user, "AppData", "Roaming", "anadius", "EA DLC Unlocker v2")
	logPath := filepath.Join(prefixPath, "drive_c", "users", user, "AppData", "Local", "anadius", "EA DLC Unlocker v2")

	info := &eaInfo{
		Client:     clientType,
		ClientPath: clientPath,
		DllPath:    filepath.Join(clientPath, "version.dll"),
		ConfigPath: configPath,
		LogPath:    logPath,
		User:       user,
	}

	// EA app also installs into a staged directory that gets refreshed on update.
	if clientType == "ea_app" {
		stagedParent := filepath.Join(filepath.Dir(clientPath), "..", "StagedEADesktop", "EA Desktop")
		stagedParent = filepath.Clean(stagedParent)
		if directoryExists(stagedParent) {
			info.DllPath2 = filepath.Join(stagedParent, "version.dll")
		}
		// The staged folder may not exist yet; still record the path so install can create it.
		if info.DllPath2 == "" {
			info.DllPath2 = filepath.Join(stagedParent, "version.dll")
		}
	}

	return info
}

// hasEAOrOrigin reports whether the prefix contains either EA app or Origin.
func hasEAOrOrigin(prefixPath string) bool {
	base := filepath.Join(prefixPath, "drive_c", "Program Files", "Electronic Arts", "EA Desktop")
	if _, ok := findVersionedDir(base); ok {
		return true
	}
	if directoryExists(filepath.Join(prefixPath, "drive_c", "Program Files (x86)", "Origin")) {
		return true
	}
	if directoryExists(filepath.Join(prefixPath, "drive_c", "Program Files", "Origin")) {
		return true
	}
	return false
}
