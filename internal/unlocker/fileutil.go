package unlocker

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func createDir(path string) error {
	return os.MkdirAll(path, 0o755)
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("could not open source: %w", err)
	}
	defer in.Close()

	if err := createDir(filepath.Dir(dst)); err != nil {
		return err
	}

	out, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("could not create destination: %w", err)
	}
	defer out.Close()

	if _, err := io.Copy(out, in); err != nil {
		return fmt.Errorf("could not copy file: %w", err)
	}
	return out.Close()
}

func removeFile(path string) error {
	if !fileExists(path) {
		return nil
	}
	return os.Remove(path)
}

func removeDirRecursive(path string) error {
	if !directoryExists(path) {
		return nil
	}
	return os.RemoveAll(path)
}

func removeDirIfEmpty(path string) {
	if entries, err := os.ReadDir(path); err == nil && len(entries) == 0 {
		_ = os.Remove(path)
	}
}
