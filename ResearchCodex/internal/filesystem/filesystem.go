package filesystem

import (
	"errors"
	"os"
	"path/filepath"
)

// EnsureDir makes sure the provided directory exists.
func EnsureDir(path string) error {
	return os.MkdirAll(path, 0o755)
}

// WriteFile writes content with 0644 permissions, creating parent directories.
func WriteFile(path string, data []byte) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}

// WriteFileIfAbsent writes content only if the file does not already exist.
func WriteFileIfAbsent(path string, data []byte) error {
	if _, err := os.Stat(path); err == nil {
		return nil
	} else if !errors.Is(err, os.ErrNotExist) {
		return err
	}
	return WriteFile(path, data)
}
