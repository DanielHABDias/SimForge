//go:build linux

package finder

import (
	"os"
	"path/filepath"
)

func find() ([]Installation, error) {
	result := findWine()
	result = append(result, findSteam()...)
	result = append(result, findLutris()...)
	result = append(result, findBottles()...)

	return dedupeInstalls(result), nil
}

// dedupeInstalls removes installations that share the same prefix path.
// The flatpak Steam setup exposes the same prefix through several Steam paths
// (e.g. ~/.var/app/com.valvesoftware.Steam/.local/share/Steam and
// ~/.var/app/com.valvesoftware.Steam/.steam/steam), producing duplicates.
func dedupeInstalls(installs []Installation) []Installation {
	seen := make(map[string]bool)
	result := make([]Installation, 0, len(installs))
	for _, inst := range installs {
		key := inst.PrefixPath
		if key == "" {
			// Fall back to the client path for Windows-like entries.
			key = inst.ClientPath
		}
		if key == "" || seen[key] {
			continue
		}
		seen[key] = true
		result = append(result, inst)
	}
	return result
}

// checkPrefix verifies if the given prefix contains an EA app or Origin
// installation and returns a completed Installation.
func checkPrefix(path string, launcher string, name string) *Installation {
	var ea *eaInfo

	if e := findEA(path, "ea_app"); e != nil {
		ea = e
	} else if e := findEA(path, "origin"); e != nil {
		ea = e
	}

	if ea == nil {
		return nil
	}

	inst := &Installation{
		Platform:   "linux",
		Launcher:   launcher,
		Name:       name,
		PrefixPath: path,
	}
	applyEAInfo(inst, ea)
	resolveSims4(inst)
	return inst
}

func findWine() []Installation {
	home, _ := os.UserHomeDir()
	prefix := filepath.Join(home, ".wine")
	if item := checkPrefix(prefix, "Wine", "Default Wine"); item != nil {
		return []Installation{*item}
	}
	return nil
}

func findLutris() []Installation {
	home, _ := os.UserHomeDir()
	base := filepath.Join(home, "Games")

	var result []Installation
	files, _ := os.ReadDir(base)
	for _, f := range files {
		if !f.IsDir() {
			continue
		}
		p := filepath.Join(base, f.Name())
		if item := checkPrefix(p, "Lutris", f.Name()); item != nil {
			result = append(result, *item)
		}
	}
	return result
}

func findBottles() []Installation {
	home, _ := os.UserHomeDir()
	paths := []string{
		".local/share/bottles/bottles",
		".var/app/com.usebottles.bottles/data/bottles/bottles",
	}

	var result []Installation
	for _, base := range paths {
		base = filepath.Join(home, base)
		files, _ := os.ReadDir(base)
		for _, f := range files {
			if !f.IsDir() {
				continue
			}
			p := filepath.Join(base, f.Name())
			if item := checkPrefix(p, "Bottles", f.Name()); item != nil {
				result = append(result, *item)
			}
		}
	}
	return result
}
