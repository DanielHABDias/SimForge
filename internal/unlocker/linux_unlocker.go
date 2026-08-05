//go:build linux

package unlocker

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// linuxDllOverrider edits the Wine/Proton prefix's user.reg to register
// the "version" DLL override (version="native,builtin").
type linuxDllOverrider struct {
	prefixPath string
}

func (o *linuxDllOverrider) regPath() (string, error) {
	reg := filepath.Join(o.prefixPath, "user.reg")
	if !fileExists(reg) {
		return "", fmt.Errorf("user.reg not found in prefix %s", o.prefixPath)
	}
	return reg, nil
}

func (o *linuxDllOverrider) Apply() error {
	reg, err := o.regPath()
	if err != nil {
		return err
	}

	// Append the override to the end of user.reg.
	entry := fmt.Sprintf(`[Software\\Wine\\DllOverrides] %d`, timestamp())
	content := fmt.Sprintf("\n%s\n\"version\"=\"native,builtin\"\n", entry)

	f, err := os.OpenFile(reg, os.O_APPEND|os.O_WRONLY, 0)
	if err != nil {
		return fmt.Errorf("could not open user.reg for writing: %w", err)
	}
	defer f.Close()

	if _, err := f.WriteString(content); err != nil {
		return fmt.Errorf("could not write to user.reg: %w", err)
	}
	return nil
}

func (o *linuxDllOverrider) Remove() error {
	reg, err := o.regPath()
	if err != nil {
		return err
	}

	lines, err := readLines(reg)
	if err != nil {
		return err
	}

	var out []string
	for _, line := range lines {
		// Remove the version override line and the section header we added.
		if strings.Contains(line, `"version"="native,builtin"`) {
			continue
		}
		if strings.Contains(line, `[Software\\Wine\\DllOverrides]`) {
			continue
		}
		out = append(out, line)
	}

	return writeLines(reg, out)
}

func timestamp() int64 {
	return time.Now().Unix()
}

func readLines(path string) ([]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var lines []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

func writeLines(path string, lines []string) error {
	return os.WriteFile(path, []byte(strings.Join(lines, "\n")+"\n"), 0)
}
