//go:build windows

package finder

import (
	"os"
	"path/filepath"

	"golang.org/x/sys/windows/registry"
)

func find() ([]Installation, error) {
	clients := []struct {
		Client string
		Name   string
		Key    string
	}{
		{Client: "ea_app", Name: "EA App", Key: `SOFTWARE\Electronic Arts\EA Desktop`},
		{Client: "origin", Name: "Origin", Key: `SOFTWARE\WOW6432Node\Origin`},
		{Client: "origin", Name: "Origin", Key: `SOFTWARE\Origin`},
	}

	var result []Installation
	for _, client := range clients {
		path, err := getRegistryPath(client.Key)
		if err != nil {
			continue
		}
		result = append(result, createWindowsInstallation(client.Client, client.Name, path))
	}
	return result, nil
}

func getRegistryPath(keyPath string) (string, error) {
	key, err := registry.OpenKey(registry.LOCAL_MACHINE, keyPath, registry.QUERY_VALUE)
	if err != nil {
		return "", err
	}
	defer key.Close()

	value, _, err := key.GetStringValue("ClientPath")
	return value, err
}

func createWindowsInstallation(client, name, clientPath string) Installation {
	appData := os.Getenv("APPDATA")
	local := os.Getenv("LOCALAPPDATA")
	common := filepath.Join("anadius", "EA DLC Unlocker v2")

	inst := Installation{
		Platform:   "windows",
		Launcher:   name,
		Name:       name,
		Client:     client,
		ClientPath: clientPath,
		DllPath:    filepath.Join(clientPath, "version.dll"),
		ConfigPath: filepath.Join(appData, common),
		LogsPath:   filepath.Join(local, common),
		User:       os.Getenv("USERNAME"),
	}

	// EA app also installs into a staged directory.
	if client == "ea_app" {
		staged := filepath.Join(filepath.Dir(clientPath), "StagedEADesktop", "EA Desktop")
		inst.DllPath2 = filepath.Join(staged, "version.dll")
	}

	resolveSims4(&inst)
	return inst
}
