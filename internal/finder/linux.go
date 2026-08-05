//go:build linux

package finder

import (
	"os"
	"path/filepath"
)

func find() ([]Installation, error) {
	var result []Installation

	result = append(result, findWine()...)
	result = append(result, findSteam()...)
	result = append(result, findLutris()...)
	result = append(result, findBottles()...)

	return result, nil
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
