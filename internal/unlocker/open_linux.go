//go:build linux

package unlocker

import (
	"os/exec"
	"runtime"
)

func openFolder(path string) error {
	// Determine the right opener based on the desktop environment.
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "linux":
		cmd = exec.Command("xdg-open", path)
	default:
		cmd = exec.Command("xdg-open", path)
	}
	return cmd.Start()
}
