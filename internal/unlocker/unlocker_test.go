package unlocker

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/DanielHABDias/SimForge/internal/finder"
)

// fakeOverrider is a no-op DllOverrider used in tests so we don't touch a real
// Wine prefix / user.reg.
type fakeOverrider struct{}

func (fakeOverrider) Apply() error  { return nil }
func (fakeOverrider) Remove() error { return nil }

// makeAssets builds a fake unlocker assets tree in a temp dir and returns it.
func makeAssets(t *testing.T) *Assets {
	t.Helper()
	root := t.TempDir()

	// config.ini
	writeTestFile(t, filepath.Join(root, "config.ini"), "[config]\ndefaultDisabled=0\n")

	// ea_app/version.dll and origin/version.dll
	writeTestFile(t, filepath.Join(root, "ea_app", "version.dll"), "ea-dll")
	writeTestFile(t, filepath.Join(root, "origin", "version.dll"), "origin-dll")

	// Some game configs.
	writeTestFile(t, filepath.Join(root, "g_The Sims 4.ini"), "[config]\n")
	writeTestFile(t, filepath.Join(root, "g_Need For Speed Heat.ini"), "[config]\n")

	return &Assets{
		RootDir:   root,
		DllEA:     filepath.Join(root, "ea_app", "version.dll"),
		DllOrigin: filepath.Join(root, "origin", "version.dll"),
		Config:    filepath.Join(root, "config.ini"),
	}
}

// makeInstallation builds a fake finder.Installation pointing into a temp dir
// that mimics a Wine/Proton prefix.
func makeInstallation(t *testing.T) finder.Installation {
	t.Helper()
	prefix := t.TempDir()
	appData := filepath.Join(prefix, "drive_c", "users", "testuser", "AppData", "Roaming", "anadius", "EA DLC Unlocker v2")
	localAppData := filepath.Join(prefix, "drive_c", "users", "testuser", "AppData", "Local", "anadius", "EA DLC Unlocker v2")

	return finder.Installation{
		Platform:   "linux",
		Launcher:   "Wine",
		Name:       "Test Wine",
		PrefixPath: prefix,
		User:       "testuser",
		Client:     "ea_app",
		ClientPath: filepath.Join(prefix, "drive_c", "Program Files", "Electronic Arts", "EA Desktop", "EA Desktop"),
		DllPath:    filepath.Join(prefix, "drive_c", "Program Files", "Electronic Arts", "EA Desktop", "EA Desktop", "version.dll"),
		DllPath2:   filepath.Join(prefix, "drive_c", "Program Files", "Electronic Arts", "EA Desktop", "StagedEADesktop", "EA Desktop", "version.dll"),
		ConfigPath: appData,
		LogsPath:   localAppData,
	}
}

// newTestManager creates a Manager wired with the fake assets, a fake
// installation and a fake DllOverrider.
func newTestManager(t *testing.T) (*Manager, finder.Installation) {
	t.Helper()
	assets := makeAssets(t)
	inst := makeInstallation(t)
	return &Manager{Installation: inst, Assets: assets, DllOverride: fakeOverrider{}}, inst
}

func writeTestFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

// TestLocateAssets verifies LocateAssets finds the unlocker assets when the
// working directory contains them.
func TestLocateAssets(t *testing.T) {
	dir := t.TempDir()
	oldWd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(oldWd)
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}

	// Populate the default search candidate ".".
	writeTestFile(t, filepath.Join(dir, "config.ini"), "cfg")
	writeTestFile(t, filepath.Join(dir, "ea_app", "version.dll"), "dll")
	writeTestFile(t, filepath.Join(dir, "origin", "version.dll"), "dll2")

	assets, err := LocateAssets()
	if err != nil {
		t.Fatalf("LocateAssets error: %v", err)
	}
	if assets == nil {
		t.Fatal("LocateAssets returned nil assets")
	}
	if !fileExists(assets.Config) || !fileExists(assets.DllEA) || !fileExists(assets.DllOrigin) {
		t.Error("assets paths do not all exist")
	}
}

// TestLocateAssetsNotFound verifies LocateAssets returns ErrAssetsNotFound when
// no assets are present.
func TestLocateAssetsNotFound(t *testing.T) {
	dir := t.TempDir()
	oldWd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(oldWd)
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}

	if _, err := LocateAssets(); err != ErrAssetsNotFound {
		t.Errorf("expected ErrAssetsNotFound, got %v", err)
	}
}

// TestInstall verifies Install copies the DLL and config, and reports installed.
func TestInstall(t *testing.T) {
	mgr, inst := newTestManager(t)

	if mgr.IsInstalled() {
		t.Fatal("should not be installed before Install")
	}

	if err := mgr.Install(); err != nil {
		t.Fatalf("Install error: %v", err)
	}

	if !mgr.IsInstalled() {
		t.Error("IsInstalled should be true after install")
	}

	// DLL in primary and staged locations.
	if !fileExists(inst.DllPath) {
		t.Error("primary version.dll not installed")
	}
	if !fileExists(inst.DllPath2) {
		t.Error("staged version.dll not installed")
	}

	// config.ini in the config folder.
	if !fileExists(filepath.Join(inst.ConfigPath, "config.ini")) {
		t.Error("config.ini not installed in config folder")
	}
}

// TestSourceDll verifies the correct source DLL is returned per client type.
func TestSourceDll(t *testing.T) {
	mgr, _ := newTestManager(t)
	mgr.Installation.Client = "ea_app"

	dll, err := mgr.SourceDll()
	if err != nil {
		t.Fatalf("SourceDll error: %v", err)
	}
	if dll != mgr.Assets.DllEA {
		t.Errorf("SourceDll = %q, want %q", dll, mgr.Assets.DllEA)
	}

	mgr.Installation.Client = "origin"
	dll, err = mgr.SourceDll()
	if err != nil {
		t.Fatalf("SourceDll error: %v", err)
	}
	if dll != mgr.Assets.DllOrigin {
		t.Errorf("SourceDll = %q, want %q", dll, mgr.Assets.DllOrigin)
	}
}

// TestGameConfigNames verifies all game configs are listed.
func TestGameConfigNames(t *testing.T) {
	mgr, _ := newTestManager(t)

	names, err := mgr.GameConfigNames()
	if err != nil {
		t.Fatalf("GameConfigNames error: %v", err)
	}
	if len(names) != 2 {
		t.Fatalf("GameConfigNames len = %d, want 2", len(names))
	}

	// Should include both names (order not guaranteed).
	got := make(map[string]bool)
	for _, n := range names {
		got[n] = true
	}
	if !got["The Sims 4"] || !got["Need For Speed Heat"] {
		t.Errorf("GameConfigNames = %v, want The Sims 4 and Need For Speed Heat", names)
	}
}

// TestAddGameConfig verifies a game config is copied and the etag removed.
func TestAddGameConfig(t *testing.T) {
	mgr, inst := newTestManager(t)

	// Simulate a log file that should be removed.
	writeTestFile(t, filepath.Join(inst.LogsPath, "The Sims 4.etag"), "etag")

	if err := mgr.AddGameConfig("The Sims 4"); err != nil {
		t.Fatalf("AddGameConfig error: %v", err)
	}

	if !fileExists(filepath.Join(inst.ConfigPath, "The Sims 4.ini")) {
		t.Error("The Sims 4.ini not installed")
	}
	if fileExists(filepath.Join(inst.LogsPath, "The Sims 4.etag")) {
		t.Error("The Sims 4.etag should have been removed")
	}
}

// TestInstalledGameConfigs verifies installed configs are listed correctly.
func TestInstalledGameConfigs(t *testing.T) {
	mgr, inst := newTestManager(t)

	// Nothing installed yet.
	if got := mgr.InstalledGameConfigs(); len(got) != 0 {
		t.Errorf("InstalledGameConfigs = %v, want empty", got)
	}

	// Install two configs.
	writeTestFile(t, filepath.Join(inst.ConfigPath, "g_The Sims 4.ini"), "x")
	writeTestFile(t, filepath.Join(inst.ConfigPath, "g_It Takes Two.ini"), "y")

	got := mgr.InstalledGameConfigs()
	if len(got) != 2 {
		t.Fatalf("InstalledGameConfigs len = %d, want 2", len(got))
	}
	gotMap := make(map[string]bool)
	for _, g := range got {
		gotMap[g] = true
	}
	if !gotMap["The Sims 4"] || !gotMap["It Takes Two"] {
		t.Errorf("InstalledGameConfigs = %v", got)
	}
}

// TestUninstall verifies Uninstall removes DLLs, config and logs.
func TestUninstall(t *testing.T) {
	mgr, inst := newTestManager(t)

	if err := mgr.Install(); err != nil {
		t.Fatalf("Install error: %v", err)
	}

	if err := mgr.Uninstall(); err != nil {
		t.Fatalf("Uninstall error: %v", err)
	}

	if mgr.IsInstalled() {
		t.Error("IsInstalled should be false after uninstall")
	}
	if fileExists(inst.DllPath) {
		t.Error("primary version.dll still present after uninstall")
	}
	if fileExists(inst.DllPath2) {
		t.Error("staged version.dll still present after uninstall")
	}
	if fileExists(filepath.Join(inst.ConfigPath, "config.ini")) {
		t.Error("config.ini still present after uninstall")
	}
}
