//go:build windows

package unlocker

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// verifyOverride checks whether the EA App override was applied on Windows:
// the scheduled "copy_dlc_unlocker" task exists and machine.ini contains
// machine.bgsstandaloneenabled=0.
func verifyOverride(m *Manager) []VerificationItem {
	if m.Installation.Client != "ea_app" {
		// Origin needs no override registration.
		return []VerificationItem{{
			Name: "version DLL override (Windows)",
			Path: "(not required for Origin)",
			OK:   true,
		}}
	}

	var items []VerificationItem

	// Scheduled task.
	taskOK := scheduledTaskExists("copy_dlc_unlocker")
	items = append(items, VerificationItem{
		Name:    "Scheduled task (copy_dlc_unlocker)",
		Path:    "schtasks /Query /TN copy_dlc_unlocker",
		OK:      taskOK,
		Missing: !taskOK,
	})

	// machine.ini setting.
	programData := os.Getenv("ProgramData")
	if programData == "" {
		programData = `C:\ProgramData`
	}
	ini := filepath.Join(programData, "EA Desktop", "machine.ini")
	iniOK := machineIniHasSetting(ini)
	items = append(items, VerificationItem{
		Name:    "machine.ini setting",
		Path:    ini,
		OK:      iniOK,
		Missing: !iniOK,
	})

	return items
}

func scheduledTaskExists(name string) bool {
	err := exec.Command("schtasks", "/Query", "/TN", name).Run()
	return err == nil
}

func machineIniHasSetting(ini string) bool {
	data, err := os.ReadFile(ini)
	if err != nil {
		return false
	}
	return strings.Contains(string(data), "machine.bgsstandaloneenabled=0")
}
