package finder

import (
	"os"
	"path/filepath"
	"testing"
)

// makeSteamLayout creates a realistic Steam prefix + game library layout
// and returns the library root.
//
// Layout:
//
//	<lib>/steamapps/common/The Sims 4/
//	<lib>/steamapps/compatdata/1222670/pfx/                 (the prefix)
//	<lib>/steamapps/compatdata/1222670/pfx/drive_c/users/steamuser/Documents/Electronic Arts/The Sims 4/Mods/
func makeSteamLayout(t *testing.T) (libRoot, prefixPath string) {
	t.Helper()
	libRoot = t.TempDir()

	gameRoot := filepath.Join(libRoot, "steamapps", "common", "The Sims 4")
	prefix := filepath.Join(libRoot, "steamapps", "compatdata", "1222670", "pfx")
	mods := filepath.Join(prefix, "drive_c", "users", "steamuser", "Documents", "Electronic Arts", "The Sims 4", "Mods")

	for _, dir := range []string{gameRoot, mods} {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			t.Fatal(err)
		}
	}
	return libRoot, prefix
}

// TestResolveSims4SteamLayout verifies resolveSims4 finds both the game root
// and the Mods folder from a realistic Steam Proton layout.
func TestResolveSims4SteamLayout(t *testing.T) {
	_, prefix := makeSteamLayout(t)

	inst := &Installation{
		AppID:      "1222670", // The Sims 4
		PrefixPath: prefix,
		User:       "steamuser",
	}

	// Temporarily disable the interactive prompt so we don't hang on failure.
	oldPrompt := promptPath
	promptPath = func(p string) string { return "" }
	defer func() { promptPath = oldPrompt }()

	resolveSims4(inst)

	if inst.Sims4Root == "" {
		t.Error("Sims4Root not detected from Steam layout")
	} else if !dirExists(inst.Sims4Root) {
		t.Errorf("Sims4Root points to non-existing dir: %s", inst.Sims4Root)
	}

	if inst.Sims4Mods == "" {
		t.Error("Sims4Mods not detected from Steam layout")
	} else if !dirExists(inst.Sims4Mods) {
		t.Errorf("Sims4Mods points to non-existing dir: %s", inst.Sims4Mods)
	}
}

// TestResolveSims4PromptFallback verifies resolveSims4 prompts when the game
// root cannot be auto-detected, and that a user-provided root is used.
func TestResolveSims4PromptFallback(t *testing.T) {
	libRoot, prefix := makeSteamLayout(t)
	// Remove the game root so detection fails.
	if err := os.RemoveAll(filepath.Join(libRoot, "steamapps", "common")); err != nil {
		t.Fatal(err)
	}

	inst := &Installation{
		AppID:      "1222670",
		PrefixPath: prefix,
		User:       "steamuser",
	}

	// Provide a fake game root via the prompt.
	oldPrompt := promptPath
	promptPath = func(p string) string { return filepath.Join(libRoot, "steamapps", "common", "The Sims 4") }
	defer func() { promptPath = oldPrompt }()

	resolveSims4(inst)

	if inst.Sims4Root == "" {
		t.Error("Sims4Root should use the prompted path")
	}
}

// TestFindSims4Mods verifies findSims4Mods returns the expected Mods path.
func TestFindSims4Mods(t *testing.T) {
	documents := t.TempDir()
	mods := filepath.Join(documents, "Electronic Arts", "The Sims 4", "Mods")
	if err := os.MkdirAll(mods, 0o755); err != nil {
		t.Fatal(err)
	}

	got := findSims4Mods(documents)
	if got != mods {
		t.Errorf("findSims4Mods = %q, want %q", got, mods)
	}
}

// TestFindSims4ModsNonExistent verifies findSims4Mods still returns the expected
// path even when the Mods folder does not exist yet.
func TestFindSims4ModsNonExistent(t *testing.T) {
	documents := t.TempDir()
	want := filepath.Join(documents, "Electronic Arts", "The Sims 4", "Mods")

	got := findSims4Mods(documents)
	if got != want {
		t.Errorf("findSims4Mods = %q, want %q", got, want)
	}
}

// TestLookForTheSims4 verifies lookForTheSims4 finds the game in a common folder.
func TestLookForTheSims4(t *testing.T) {
	common := t.TempDir()
	gameRoot := filepath.Join(common, "The Sims 4")
	if err := os.MkdirAll(gameRoot, 0o755); err != nil {
		t.Fatal(err)
	}

	got := lookForTheSims4(common)
	if got != gameRoot {
		t.Errorf("lookForTheSims4 = %q, want %q", got, gameRoot)
	}
}

func dirExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}
