package finder

// Installation represents a detected EA app / Origin installation
// within a particular prefix (Windows) or Wine/Proton prefix (Linux).
type Installation struct {
	Platform string // "windows" or "linux"
	Launcher string // "EA App", "Origin", "Steam Proton", "Wine", "Lutris", "Bottles"
	Name     string // display name of the prefix / client
	AppID    string // Steam app ID (empty when not a Steam prefix)

	PrefixPath string // Wine/Proton prefix root (drive_c parent). Empty on Windows.
	User       string // the OS/wine user owning the prefix

	Client string // "ea_app" or "origin"

	ClientPath string // folder containing the EA client executable
	DllPath    string // primary destination for version.dll
	DllPath2   string // staged destination for version.dll (EA app only)

	ConfigPath string // Roaming\anadius\EA DLC Unlocker v2
	LogsPath   string // Local\anadius\EA DLC Unlocker v2

	// Sims4Root is the root folder of the The Sims 4 game installation,
	// where DLC files are placed.
	Sims4Root string
	// Sims4Mods is the game's Mods folder (Documents\Electronic Arts\The Sims 4\Mods).
	Sims4Mods string
}

// Find returns all detected EA/Origin installations for the current OS.
func Find() ([]Installation, error) {
	return find()
}
