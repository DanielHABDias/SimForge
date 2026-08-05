//go:build windows

package finder

import (
	"os"
	"path/filepath"
	"testing"
)

// TestSims4CandidatesWindowsProgramFiles verifies sims4Candidates checks the
// standard Windows install locations for the game root (where DLCs go).
func TestSims4CandidatesWindowsProgramFiles(t *testing.T) {
	// The candidates are static paths; we just verify the expected entries are
	// present so the game root detection covers the common install locations.
	inst := &Installation{}

	cands := sims4Candidates(inst)
	if len(cands) == 0 {
		t.Fatal("sims4Candidates returned no candidates")
	}

	want := []string{
		`C:\Program Files (x86)\EA Games\The Sims 4`,
		`C:\Program Files (x86)\Origin Games\The Sims 4`,
		`C:\Program Files (x86)\Steam\steamapps\common\The Sims 4`,
		`C:\Program Files\EA Games\The Sims 4`,
		`C:\Program Files\Origin Games\The Sims 4`,
		`C:\Program Files\Steam\steamapps\common\The Sims 4`,
	}

	gotSet := make(map[string]bool)
	for _, c := range cands {
		gotSet[c] = true
	}
	for _, w := range want {
		if !gotSet[w] {
			t.Errorf("sims4Candidates missing %q (got %v)", w, cands)
		}
	}
}

// TestSims4ModsWindowsDocuments verifies the Mods folder is resolved under the
// user's Documents\Electronic Arts\The Sims 4 folder on Windows.
func TestSims4ModsWindowsDocuments(t *testing.T) {
	// We can't easily override USERPROFILE in-process, so verify the shared
	// findSims4Mods helper produces the expected Mods path from any Documents
	// base (this is what sims4Mods uses on Windows).
	documents := filepath.Join("C:", "Users", "TestUser", "Documents")
	mods := filepath.Join(documents, "Electronic Arts", "The Sims 4", "Mods")

	got := findSims4Mods(documents)
	if got != mods {
		t.Errorf("findSims4Mods = %q, want %q", got, mods)
	}
}

// TestSims4ModsNonExistentWindows verifies the Mods path is still returned even
// when the folder does not exist yet (it will be created by the transfer step).
func TestSims4ModsNonExistentWindows(t *testing.T) {
	documents := filepath.Join("C:", "Users", "TestUser", "Documents")
	want := filepath.Join(documents, "Electronic Arts", "The Sims 4", "Mods")

	got := findSims4Mods(documents)
	if got != want {
		t.Errorf("findSims4Mods = %q, want %q", got, want)
	}
}

// TestSims4ModsWindowsUsesUserProfile verifies sims4Mods uses USERPROFILE to
// build the Documents path.
func TestSims4ModsWindowsUsesUserProfile(t *testing.T) {
	old := os.Getenv("USERPROFILE")
	defer os.Setenv("USERPROFILE", old)

	os.Setenv("USERPROFILE", `C:\Users\TestUser`)
	inst := &Installation{}

	got := sims4Mods(inst)
	want := filepath.Join(`C:\Users\TestUser`, "Documents", "Electronic Arts", "The Sims 4", "Mods")
	if got != want {
		t.Errorf("sims4Mods = %q, want %q", got, want)
	}
}
