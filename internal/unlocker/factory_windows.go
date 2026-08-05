//go:build windows

package unlocker

import "github.com/DanielHABDias/SimForge/internal/finder"

func newDllOverrider(inst finder.Installation) DllOverrider {
	return &windowsDllOverrider{inst: inst}
}
