//go:build windows

package unlocker

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/DanielHABDias/SimForge/internal/finder"
)

// windowsDllOverrider handles the EA app behaviour on Windows:
//   - creates a scheduled task that copies version.dll into the staged folder
//   - writes machine.bgsstandaloneenabled=0 into machine.ini
//
// For Origin there is nothing to register.
type windowsDllOverrider struct {
	inst finder.Installation
}

func (o *windowsDllOverrider) Apply() error {
	if o.inst.Client != "ea_app" {
		return nil
	}
	if err := o.createScheduledTask(); err != nil {
		return err
	}
	return o.writeMachineIni()
}

func (o *windowsDllOverrider) Remove() error {
	if o.inst.Client != "ea_app" {
		return nil
	}
	return o.deleteScheduledTask()
}

// createScheduledTask registers a one-time (recreated on EA app update)
// task that copies version.dll into the staged folder.
func (o *windowsDllOverrider) createScheduledTask() error {
	// Paths use forward slashes to satisfy schtasks quoting.
	staged := filepath.ToSlash(o.inst.DllPath2)
	dest := o.inst.DllPath

	tr := fmt.Sprintf("xcopy.exe /Y '%s' '%s'", dest, staged)

	cmd := exec.Command("schtasks",
		"/Create", "/F", "/RL", "HIGHEST",
		"/SC", "ONCE", "/ST", "00:00", "/SD", "01/01/2000",
		"/TN", "copy_dlc_unlocker",
		"/TR", tr,
	)
	if err := cmd.Run(); err != nil {
		// Try alternate date format.
		cmd = exec.Command("schtasks",
			"/Create", "/F", "/RL", "HIGHEST",
			"/SC", "ONCE", "/ST", "00:00", "/SD", "2000/01/01",
			"/TN", "copy_dlc_unlocker",
			"/TR", tr,
		)
		if err2 := cmd.Run(); err2 != nil {
			return fmt.Errorf("could not create the copy scheduled task: %w", err2)
		}
	}
	return nil
}

func (o *windowsDllOverrider) deleteScheduledTask() error {
	cmd := exec.Command("schtasks", "/Delete", "/TN", "copy_dlc_unlocker", "/F")
	if err := cmd.Run(); err != nil {
		// Task may not exist — not fatal.
		return nil
	}
	return nil
}

// writeMachineIni appends machine.bgsstandaloneenabled=0 to the EA Desktop
// machine.ini so the staged copy is actually used.
func (o *windowsDllOverrider) writeMachineIni() error {
	programData := os.Getenv("ProgramData")
	if programData == "" {
		programData = `C:\ProgramData`
	}
	ini := filepath.Join(programData, "EA Desktop", "machine.ini")

	f, err := os.OpenFile(ini, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0)
	if err != nil {
		return fmt.Errorf("could not write machine.ini: %w", err)
	}
	defer f.Close()

	_, err = f.WriteString("machine.bgsstandaloneenabled=0\n")
	if err != nil {
		return fmt.Errorf("could not write machine.ini: %w", err)
	}
	return nil
}
