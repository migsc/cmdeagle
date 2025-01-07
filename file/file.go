package file

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

//go:embed *
var PackageFS embed.FS

func ResolvePath(path string) (string, error) {
	expandedPath, err := ExpandPath(path)
	if err != nil {
		return "", err
	}

	absPath, err := filepath.Abs(expandedPath)
	if err != nil {
		return "", err
	}

	return absPath, nil
}

// CheckDirExists checks if a directory exists at the given path
func CheckDirExists(path string) bool {
	dirStats, err := os.Stat(path)
	if err != nil {
		return false
	}
	return dirStats.IsDir()
}

func CheckFileExists(path string) bool {
	fileStats, err := os.Stat(path)
	if err != nil {
		return false
	}
	return !fileStats.IsDir()
}

// CreateDir creates a directory at the given path
func CreateDir(path string) error {
	err := os.MkdirAll(path, 0755)
	if err != nil {
		fmt.Printf("Error creating directory %s: %v\n", path, err)
		return err
	}

	fmt.Printf("Created directory %s\n", path)

	return nil
}

func DeleteFile(path string) error {
	err := os.Remove(path)
	if err != nil {
		fmt.Printf("Error deleting file %s: %v\n", path, err)
		return err
	}

	return nil
}

// CleanDir removes all contents of a directory while preserving the directory itself
func CleanDir(path string) error {
	// First check if directory exists
	if !CheckDirExists(path) {
		return fmt.Errorf("directory %s does not exist", path)
	}

	// Read directory contents
	entries, err := os.ReadDir(path)
	if err != nil {
		return fmt.Errorf("failed to read directory %s: %w", path, err)
	}

	// Remove each entry
	for _, entry := range entries {
		entryPath := filepath.Join(path, entry.Name())
		if err := os.RemoveAll(entryPath); err != nil {
			return fmt.Errorf("failed to remove %s: %w", entryPath, err)
		}
	}

	return nil
}

func SetupEmptyDir(path string) error {
	if CheckFileExists(path) {
		DeleteFile(path)
	}

	if CheckDirExists(path) {
		return CleanDir(path)
	} else {
		return CreateDir(path)
	}
}

// FindFileEndsWithPattern looks for a file matching the pattern in the given directory
func FindFileEndsWithPattern(dir string, pattern string) (string, error) {
	files, err := os.ReadDir(dir)
	if err != nil {
		return "", err
	}

	for _, file := range files {
		if !file.IsDir() {
			name := file.Name()
			if filepath.Ext(name) == ".yaml" || filepath.Ext(name) == ".yml" {
				return name, nil
			}
		}
	}

	return "", fmt.Errorf("no file matching the pattern %s found in directory %s", pattern, dir)
}

func ExpandPath(path string) (string, error) {
	if path == "" {
		return "", fmt.Errorf("empty path")
	}

	// Handle home directory expansion
	if strings.HasPrefix(path, "~/") {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		path = filepath.Join(home, path[2:])
	}

	// Convert to absolute path
	absPath, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}

	// Clean the path
	return filepath.Clean(absPath), nil
}
