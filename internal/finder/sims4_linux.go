//go:build linux

package finder

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// sims4Candidates returns candidate game root folders for The Sims 4 on Linux.
func sims4Candidates(inst *Installation) []string {
	var candidates []string

	// Steam Proton: the game lives in the Steam library, not the prefix.
	if inst.AppID == "1222670" && inst.PrefixPath != "" {
		// Walk up to find the steamapps/common folder.
		prefix := inst.PrefixPath
		// The prefix is usually <library>/steamapps/compatdata/1222670/pfx
		compatData := filepath.Dir(filepath.Dir(prefix))
		if strings.HasSuffix(filepath.ToSlash(compatData), "compatdata") {
			library := filepath.Dir(filepath.Dir(compatData)) // steamapps
			common := filepath.Join(filepath.Dir(library), "common")
			if root := lookForTheSims4(common); root != "" {
				candidates = append(candidates, root)
			}
		}
	}

	// Wine generic installs.
	if inst.PrefixPath != "" {
		driveC := filepath.Join(inst.PrefixPath, "drive_c")
		for _, base := range []string{
			filepath.Join(driveC, "Program Files (x86)"),
			filepath.Join(driveC, "Program Files"),
		} {
			for _, sub := range []string{
				filepath.Join(base, "EA Games", "The Sims 4"),
				filepath.Join(base, "Origin Games", "The Sims 4"),
				filepath.Join(base, "Steam", "steamapps", "common", "The Sims 4"),
			} {
				if directoryExists(sub) {
					candidates = append(candidates, sub)
				}
			}
		}
	}

	return candidates
}

// sims4Mods returns the The Sims 4 Mods folder on Linux.
func sims4Mods(inst *Installation) string {
	if inst.PrefixPath == "" {
		return ""
	}
	// Standard wine Documents folder.
	documents := filepath.Join(inst.PrefixPath, "drive_c", "users", inst.User, "Documents")
	return findSims4Mods(documents)
}

// lookForTheSims4 finds "The Sims 4" inside a Steam common folder.
func lookForTheSims4(common string) string {
	p := filepath.Join(common, "The Sims 4")
	if directoryExists(p) {
		return p
	}
	return ""
}

// promptPath asks the user to type a directory path.
func promptPath(prompt string) string {
	fmt.Print(prompt)
	reader := bufio.NewReader(os.Stdin)
	line, _ := reader.ReadString('\n')
	return strings.TrimSpace(line)
}
