//go:build windows

package unlocker

import "os/exec"

func openFolder(path string) error {
	return exec.Command("explorer", path).Start()
}
