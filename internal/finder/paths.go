package finder

import "os"

// directoryExists reports whether path exists and is a directory.
func directoryExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}
