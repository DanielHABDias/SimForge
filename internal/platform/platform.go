package platform

import (
	"os"
	"runtime"
)

type OS string

const (
	Windows OS = "windows"
	Linux   OS = "linux"
	MacOS   OS = "darwin"
)

func Current() OS {
	return OS(runtime.GOOS)
}

func IsWindows() bool {
	return Current() == Windows
}

func IsLinux() bool {
	return Current() == Linux
}

func IsMacOS() bool {
	return Current() == MacOS
}

func HomeDir() (string, error) {
	return os.UserHomeDir()
}

func Arch() string {
	return runtime.GOARCH
}
