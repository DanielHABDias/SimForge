//go:build windows

package finder

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// sims4Candidates returns candidate game root folders for The Sims 4 on Windows.
func sims4Candidates(inst *Installation) []string {
	var candidates []string

	drives := []string{"C:", "D:", "E:"}
	for _, drive := range drives {
		for _, base := range []string{
			filepath.Join(drive+"\\", "Program Files (x86)"),
			filepath.Join(drive+"\\", "Program Files"),
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

// sims4Mods returns the The Sims 4 Mods folder on Windows.
func sims4Mods(inst *Installation) string {
	documents := os.Getenv("USERPROFILE")
	if documents == "" {
		return ""
	}
	return findSims4Mods(filepath.Join(documents, "Documents"))
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
func promptPathImpl(prompt string) string {
	fmt.Print(prompt)
	reader := bufio.NewReader(os.Stdin)
	line, _ := reader.ReadString('\n')
	return strings.TrimSpace(line)
}
