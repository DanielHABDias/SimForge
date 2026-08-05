// Package transfer copies downloaded DLCs and mods into the game folders,
// extracting any compressed files found along the way.
package transfer

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// maxDepth limits how deeply we recurse when extracting nested archives.
const maxDepth = 5

// Result summarises a transfer operation.
type Result struct {
	Source        string
	Destination   string
	CopiedFiles   int
	ExtractedZips int
}

// CountFiles returns the number of regular files (recursively) under path.
func CountFiles(path string) (int, error) {
	var count int
	err := filepath.Walk(path, func(p string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			count++
		}
		return nil
	})
	return count, err
}

// CopyDLCs copies all files from src into dst (the game root), extracting any
// archives found. Returns a summary of what was done.
func CopyDLCs(src, dst string) (*Result, error) {
	return copyTree(src, dst)
}

// CopyMods copies all files from src into dst (the game Mods folder), extracting
// any archives found. Returns a summary of what was done.
func CopyMods(src, dst string) (*Result, error) {
	return copyTree(src, dst)
}

// copyTree walks src and copies each file into dst. Archives are extracted
// directly into dst (recursively), so their contents land where they belong.
func copyTree(src, dst string) (*Result, error) {
	info, err := os.Stat(src)
	if err != nil {
		return nil, fmt.Errorf("source path not available: %w", err)
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("source path is not a directory: %s", src)
	}

	res := &Result{
		Source:      src,
		Destination: dst,
	}

	if err := os.MkdirAll(dst, 0o755); err != nil {
		return nil, fmt.Errorf("could not create destination folder: %w", err)
	}

	err = filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		// Archives are extracted directly into the destination root.
		if isArchive(path) {
			if err := extractNested(path, dst, 1, res); err != nil {
				return err
			}
			res.ExtractedZips++
			return nil
		}

		rel, rerr := filepath.Rel(src, path)
		if rerr != nil {
			return rerr
		}
		target := filepath.Join(dst, rel)
		if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
			return err
		}
		if err := copyFile(path, target); err != nil {
			return err
		}
		res.CopiedFiles++
		return nil
	})

	return res, err
}

// extractNested extracts an archive into dst, then recursively extracts any
// archives produced by that extraction (e.g. a zip containing other zips).
func extractNested(archivePath, dst string, depth int, res *Result) error {
	if depth > maxDepth {
		return fmt.Errorf("archive nesting too deep at %s", archivePath)
	}

	if err := extractArchive(archivePath, dst); err != nil {
		return fmt.Errorf("could not extract %s: %w", archivePath, err)
	}

	// Recurse into any archives found inside the destination produced by this
	// extraction. Nested archive files are extracted into their parent folder
	// so their relative structure is preserved, and removed after successful
	// extraction to keep the destination clean.
	return filepath.Walk(dst, func(p string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() || !isArchive(p) {
			return nil
		}
		// Never try to re-extract the archive we are currently processing.
		if isSamePath(p, archivePath) {
			return nil
		}
		// Extract into the archive's parent directory so nested content keeps
		// its relative location (e.g. nested/inner.zip -> nested/GP08/...).
		parent := filepath.Dir(p)
		if err := extractNested(p, parent, depth+1, res); err != nil {
			return err
		}
		res.ExtractedZips++
		_ = os.Remove(p)
		return nil
	})
}

// extractArchive dispatches to the right extractor based on the file extension.
func extractArchive(path, dst string) error {
	switch {
	case strings.HasSuffix(strings.ToLower(path), ".zip"):
		return unzip(path, dst)
	case strings.HasSuffix(strings.ToLower(path), ".tar.gz"),
		strings.HasSuffix(strings.ToLower(path), ".tgz"):
		return untarGz(path, dst)
	case strings.HasSuffix(strings.ToLower(path), ".tar"):
		return untar(path, dst)
	case strings.HasSuffix(strings.ToLower(path), ".gz"):
		return gunzip(path, dst)
	default:
		return fmt.Errorf("unsupported archive format: %s", filepath.Ext(path))
	}
}

// isArchive reports whether the file has an archive extension that we can
// extract with the standard library. (zip, tar, tar.gz, tgz, gz)
func isArchive(path string) bool {
	lower := strings.ToLower(path)
	for _, ext := range []string{".zip", ".tar", ".tar.gz", ".tgz", ".gz"} {
		if strings.HasSuffix(lower, ext) {
			return true
		}
	}
	return false
}

// isSamePath reports whether two paths refer to the same file.
func isSamePath(a, b string) bool {
	ca, _ := filepath.Abs(a)
	cb, _ := filepath.Abs(b)
	return filepath.Clean(ca) == filepath.Clean(cb)
}

// copyFile copies a single regular file from src to dst.
func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	if _, err := io.Copy(out, in); err != nil {
		return err
	}
	return out.Close()
}
