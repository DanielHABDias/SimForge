package finder

// eaInfo holds the EA app specific details found inside a Wine/Proton prefix.
type eaInfo struct {
	Client     string // "ea_app"
	ClientPath string
	DllPath    string
	DllPath2   string
	ConfigPath string
	LogPath    string
	User       string
}
