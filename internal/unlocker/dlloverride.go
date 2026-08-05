package unlocker

// DllOverrider manages the "version" library override registration.
// On Windows this is done via a scheduled task and machine.ini.
// On Linux this is done by editing the prefix's user.reg.
type DllOverrider interface {
	// Apply registers the version=native,builtin override.
	Apply() error
	// Remove removes the override.
	Remove() error
}
