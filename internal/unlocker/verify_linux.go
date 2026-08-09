//go:build linux

package unlocker

import "strings"

// verifyOverride checks whether the version DLL override was written to the
// prefix's user.reg file (the mechanism used on Linux/Wine/Proton).
func verifyOverride(m *Manager) []VerificationItem {
	override, ok := m.DllOverride.(*linuxDllOverrider)
	if !ok {
		return nil
	}

	reg, err := override.regPath()
	if err != nil {
		return []VerificationItem{{
			Name:    "version DLL override (user.reg)",
			Path:    reg,
			OK:      false,
			Missing: true,
		}}
	}

	ok = hasDllOverride(reg)
	return []VerificationItem{{
		Name:    "version DLL override (user.reg)",
		Path:    reg,
		OK:      ok,
		Missing: !ok,
	}}
}

// hasDllOverride reports whether a version.dll override entry is present in the
// given user.reg file.
func hasDllOverride(reg string) bool {
	if !fileExists(reg) {
		return false
	}
	lines, err := readLines(reg)
	if err != nil {
		return false
	}
	for _, line := range lines {
		if strings.Contains(line, `"version"="native,builtin"`) {
			return true
		}
	}
	return false
}
