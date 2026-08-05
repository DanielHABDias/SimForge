package finder

import (
	"path/filepath"
)

// resolveSims4 populates Sims4Root and Sims4Mods for an installation by trying
// the platform-specific candidates and falling back to an interactive prompt.
func resolveSims4(inst *Installation) {
	// Try to auto-detect the game root.
	for _, cand := range sims4Candidates(inst) {
		if directoryExists(cand) {
			inst.Sims4Root = cand
			break
		}
	}

	// Detect the Mods folder.
	if mods := sims4Mods(inst); mods != "" {
		inst.Sims4Mods = mods
	}

	// If the game root is still unknown, ask the user.
	if inst.Sims4Root == "" {
		inst.Sims4Root = promptPath("Could not auto-detect The Sims 4 game folder. Enter its path: ")
		if inst.Sims4Root != "" {
			// Derive the Mods folder from the provided root if possible.
			if mods := filepath.Join(inst.Sims4Root, "Mods"); directoryExists(mods) {
				inst.Sims4Mods = mods
			}
		}
	}
}

// findSims4Mods looks for the standard The Sims 4 Mods folder under the given
// Documents base (native or inside a wine prefix).
func findSims4Mods(documentsBase string) string {
	mods := filepath.Join(documentsBase, "Electronic Arts", "The Sims 4", "Mods")
	if directoryExists(mods) {
		return mods
	}
	// Mods folder may not exist yet; still return the expected path.
	return mods
}
