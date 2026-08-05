//go:build linux

package unlocker

import "github.com/DanielHABDias/SimForge/internal/finder"

func newDllOverrider(inst finder.Installation) DllOverrider {
	return &linuxDllOverrider{prefixPath: inst.PrefixPath}
}
