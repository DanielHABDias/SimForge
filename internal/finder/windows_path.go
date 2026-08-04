//go:build windows

package finder

import (
	"os"
	"path/filepath"
)


func buildWindowsPaths(
	clientName string,
	clientPath string,
) Installation {


	appData := os.Getenv("APPDATA")
	localAppData := os.Getenv("LOCALAPPDATA")


	commonDir := filepath.Join(
		"anadius",
		"EA DLC Unlocker v2",
	)


	staged := filepath.Join(
		filepath.Dir(clientPath),
		"StagedEADesktop",
		"EA Desktop",
	)


	return Installation{

		ClientName: clientName,

		ClientPath: clientPath,


		DllPath: filepath.Join(
			clientPath,
			"version.dll",
		),


		StagedPath: filepath.Join(
			staged,
			"version.dll",
		),


		ConfigPath: filepath.Join(
			appData,
			commonDir,
		),


		LogPath: filepath.Join(
			localAppData,
			commonDir,
		),


		UserProfile: os.Getenv("USERNAME"),
	}
}