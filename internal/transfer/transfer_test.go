package transfer

import (
	"archive/zip"
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"testing"
)

// writeFile creates a file at path with the given content, creating parent dirs.
func writeFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

// createZip creates a zip archive at zipPath containing the given entries
// (name -> content).
func createZip(t *testing.T, zipPath string, entries map[string]string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(zipPath), 0o755); err != nil {
		t.Fatal(err)
	}
	f, err := os.Create(zipPath)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	w := zip.NewWriter(f)
	for name, content := range entries {
		zw, err := w.Create(name)
		if err != nil {
			t.Fatal(err)
		}
		if _, err := zw.Write([]byte(content)); err != nil {
			t.Fatal(err)
		}
	}
	if err := w.Close(); err != nil {
		t.Fatal(err)
	}
}

func fileExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}

func dirExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}

// TestCountFiles verifies CountFiles counts regular files recursively.
func TestCountFiles(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, filepath.Join(dir, "a.txt"), "a")
	writeFile(t, filepath.Join(dir, "sub", "b.txt"), "b")
	writeFile(t, filepath.Join(dir, "sub", "deep", "c.bin"), "c")

	count, err := CountFiles(dir)
	if err != nil {
		t.Fatalf("CountFiles error: %v", err)
	}
	if count != 3 {
		t.Errorf("CountFiles = %d, want 3", count)
	}
}

// TestCopyDLCsToRoot verifies DLC files are copied into the game root.
func TestCopyDLCsToRoot(t *testing.T) {
	src := filepath.Join(t.TempDir(), "dlcs_source")
	dst := filepath.Join(t.TempDir(), "game_root")

	writeFile(t, filepath.Join(src, "EP01", "data.package"), "dlc-data")
	writeFile(t, filepath.Join(src, "EP02", "data.package"), "dlc-data2")

	res, err := CopyDLCs(src, dst)
	if err != nil {
		t.Fatalf("CopyDLCs error: %v", err)
	}

	if !fileExists(filepath.Join(dst, "EP01", "data.package")) {
		t.Error("EP01/data.package not copied to game root")
	}
	if !fileExists(filepath.Join(dst, "EP02", "data.package")) {
		t.Error("EP02/data.package not copied to game root")
	}
	if res.CopiedFiles != 2 {
		t.Errorf("CopiedFiles = %d, want 2", res.CopiedFiles)
	}
}

// TestCopyModsToModsFolder verifies mods are copied into the Mods folder.
func TestCopyModsToModsFolder(t *testing.T) {
	src := filepath.Join(t.TempDir(), "mods_source")
	dst := filepath.Join(t.TempDir(), "game", "Mods")

	writeFile(t, filepath.Join(src, "CC.package"), "cc")
	writeFile(t, filepath.Join(src, "script.mod"), "script")

	res, err := CopyMods(src, dst)
	if err != nil {
		t.Fatalf("CopyMods error: %v", err)
	}

	if !fileExists(filepath.Join(dst, "CC.package")) {
		t.Error("CC.package not copied to Mods folder")
	}
	if !fileExists(filepath.Join(dst, "script.mod")) {
		t.Error("script.mod not copied to Mods folder")
	}
	if res.CopiedFiles != 2 {
		t.Errorf("CopiedFiles = %d, want 2", res.CopiedFiles)
	}
}

// TestCopyDLCsExtractsZip verifies a zip in the source is extracted into the
// game root.
func TestCopyDLCsExtractsZip(t *testing.T) {
	src := filepath.Join(t.TempDir(), "dlcs_source")
	dst := filepath.Join(t.TempDir(), "game_root")

	createZip(t, filepath.Join(src, "EP05.zip"), map[string]string{
		"EP05/data.package": "ep05-data",
	})

	res, err := CopyDLCs(src, dst)
	if err != nil {
		t.Fatalf("CopyDLCs error: %v", err)
	}

	if !fileExists(filepath.Join(dst, "EP05", "data.package")) {
		t.Error("EP05/data.package not extracted to game root")
	}
	if res.ExtractedZips < 1 {
		t.Errorf("ExtractedZips = %d, want >= 1", res.ExtractedZips)
	}
}

// TestCopyModsExtractsZip verifies a zip in the source is extracted into the
// Mods folder.
func TestCopyModsExtractsZip(t *testing.T) {
	src := filepath.Join(t.TempDir(), "mods_source")
	dst := filepath.Join(t.TempDir(), "game", "Mods")

	createZip(t, filepath.Join(src, "modpack.zip"), map[string]string{
		"mods/CC.package": "cc",
		"mods/script.mod": "script",
	})

	res, err := CopyMods(src, dst)
	if err != nil {
		t.Fatalf("CopyMods error: %v", err)
	}

	if !fileExists(filepath.Join(dst, "mods", "CC.package")) {
		t.Error("CC.package not extracted to Mods folder")
	}
	if !fileExists(filepath.Join(dst, "mods", "script.mod")) {
		t.Error("script.mod not extracted to Mods folder")
	}
	if res.ExtractedZips < 1 {
		t.Errorf("ExtractedZips = %d, want >= 1", res.ExtractedZips)
	}
}

// TestCopyExtractsNestedZip verifies a zip containing another zip is fully
// extracted (nested archive handled recursively).
func TestCopyExtractsNestedZip(t *testing.T) {
	src := filepath.Join(t.TempDir(), "dlcs_source")
	dst := filepath.Join(t.TempDir(), "game_root")

	// Create an inner zip file on disk.
	innerPath := filepath.Join(t.TempDir(), "inner.zip")
	createZip(t, innerPath, map[string]string{
		"GP08/data.package": "gp08-data",
	})
	innerData, err := os.ReadFile(innerPath)
	if err != nil {
		t.Fatal(err)
	}

	// Outer zip contains the inner zip.
	outerPath := filepath.Join(src, "outer.zip")
	if err := os.MkdirAll(filepath.Dir(outerPath), 0o755); err != nil {
		t.Fatal(err)
	}
	fo, err := os.Create(outerPath)
	if err != nil {
		t.Fatal(err)
	}
	zw := zip.NewWriter(fo)
	entry, err := zw.Create("nested/inner.zip")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := entry.Write(innerData); err != nil {
		t.Fatal(err)
	}
	if err := zw.Close(); err != nil {
		t.Fatal(err)
	}
	fo.Close()

	res, err := CopyDLCs(src, dst)
	if err != nil {
		t.Fatalf("CopyDLCs error: %v", err)
	}

	if !fileExists(filepath.Join(dst, "nested", "GP08", "data.package")) {
		t.Error("nested GP08/data.package not extracted")
	}
	if res.ExtractedZips < 2 {
		t.Errorf("ExtractedZips = %d, want >= 2", res.ExtractedZips)
	}
}

// TestCopyNonExistentSource verifies an error is returned for a missing source.
func TestCopyNonExistentSource(t *testing.T) {
	src := filepath.Join(t.TempDir(), "does_not_exist")
	dst := filepath.Join(t.TempDir(), "game_root")

	if _, err := CopyDLCs(src, dst); err == nil {
		t.Error("expected error for missing source")
	}
}

// TestZipSlipGuard verifies a malicious zip entry cannot escape the destination.
func TestZipSlipGuard(t *testing.T) {
	dst := t.TempDir()
	malicious := filepath.Join("..", "evil.txt")

	if _, err := zipSlipGuard(dst, malicious); err == nil {
		t.Error("expected zipSlipGuard to reject path traversal entry")
	}
}

// TestExtractPreservesContent verifies extracted file content matches.
func TestExtractPreservesContent(t *testing.T) {
	src := filepath.Join(t.TempDir(), "dlcs_source")
	dst := filepath.Join(t.TempDir(), "game_root")

	const want = "important-dlc-content"
	createZip(t, filepath.Join(src, "EP01.zip"), map[string]string{
		"EP01/data.package": want,
	})

	if _, err := CopyDLCs(src, dst); err != nil {
		t.Fatalf("CopyDLCs error: %v", err)
	}

	got, err := os.ReadFile(filepath.Join(dst, "EP01", "data.package"))
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != want {
		t.Errorf("content = %q, want %q", string(got), want)
	}
}

// helper to compute a file hash for content verification.
func hashFile(t *testing.T, path string) string {
	t.Helper()
	f, err := os.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		t.Fatal(err)
	}
	return fmt.Sprintf("%x", h.Sum(nil))
}
