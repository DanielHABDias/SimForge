//go:build linux

package finder

import (
	"os"
	"path/filepath"
	"testing"
)

// makeFlatpakSteamLayout reproduces the realistic flatpak Steam Proton layout
// for The Sims 4, matching the structure found on real Linux systems:
//
//	<home>/.var/app/com.valvesoftware.Steam/.local/share/Steam/steamapps/
//	    common/The Sims 4/                                  (game root -> DLCs)
//	    compatdata/1222670/pfx/                             (the prefix)
//	    compatdata/1222670/pfx/drive_c/users/steamuser/Documents/
//	        Electronic Arts/The Sims 4/Mods/                (Mods folder)
func makeFlatpakSteamLayout(t *testing.T) (libRoot, prefixPath, gameRoot, modsPath string) {
	t.Helper()
	libRoot = t.TempDir()

	gameRoot = filepath.Join(libRoot, "steamapps", "common", "The Sims 4")
	prefixPath = filepath.Join(libRoot, "steamapps", "compatdata", "1222670", "pfx")
	modsPath = filepath.Join(prefixPath, "drive_c", "users", "steamuser", "Documents", "Electronic Arts", "The Sims 4", "Mods")

	for _, dir := range []string{gameRoot, modsPath} {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			t.Fatal(err)
		}
	}
	return libRoot, prefixPath, gameRoot, modsPath
}

// TestResolveSims4FlatpakSteamLayout verifies both the Sims 4 game root (where
// DLCs go) and the Mods folder are detected from a flatpak Steam Proton layout.
func TestResolveSims4FlatpakSteamLayout(t *testing.T) {
	_, prefixPath, gameRoot, modsPath := makeFlatpakSteamLayout(t)

	inst := &Installation{
		AppID:      "1222670", // The Sims 4
		PrefixPath: prefixPath,
		User:       "steamuser",
	}

	oldPrompt := promptPath
	promptPath = func(p string) string { return "" }
	defer func() { promptPath = oldPrompt }()

	resolveSims4(inst)

	// DLC location: the game root.
	if inst.Sims4Root == "" {
		t.Error("Sims4Root not detected from flatpak Steam layout")
	} else if inst.Sims4Root != gameRoot {
		t.Errorf("Sims4Root = %q, want %q", inst.Sims4Root, gameRoot)
	} else if !directoryExists(inst.Sims4Root) {
		t.Errorf("Sims4Root points to non-existing dir: %s", inst.Sims4Root)
	}

	// Mods folder location.
	if inst.Sims4Mods == "" {
		t.Error("Sims4Mods not detected from flatpak Steam layout")
	} else if inst.Sims4Mods != modsPath {
		t.Errorf("Sims4Mods = %q, want %q", inst.Sims4Mods, modsPath)
	} else if !directoryExists(inst.Sims4Mods) {
		t.Errorf("Sims4Mods points to non-existing dir: %s", inst.Sims4Mods)
	}
}

// TestSims4CandidatesSteamCommon verifies sims4Candidates returns the game root
// located under steamapps/common for a Steam Proton prefix.
func TestSims4CandidatesSteamCommon(t *testing.T) {
	_, prefixPath, gameRoot, _ := makeFlatpakSteamLayout(t)

	inst := &Installation{
		AppID:      "1222670",
		PrefixPath: prefixPath,
	}

	cands := sims4Candidates(inst)
	found := false
	for _, c := range cands {
		if c == gameRoot {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("sims4Candidates = %v, expected to include %q", cands, gameRoot)
	}
}

// TestSims4ModsSteamUser verifies sims4Mods resolves the Mods folder under the
// document folder of the prefix's steamuser.
func TestSims4ModsSteamUser(t *testing.T) {
	_, prefixPath, _, modsPath := makeFlatpakSteamLayout(t)

	inst := &Installation{
		PrefixPath: prefixPath,
		User:       "steamuser",
	}

	got := sims4Mods(inst)
	if got != modsPath {
		t.Errorf("sims4Mods = %q, want %q", got, modsPath)
	}
}
