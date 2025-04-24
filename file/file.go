package file

import (
	"embed"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/afero"
)

//go:embed *
var PackageFS embed.FS

// MimeTypes maps file extensions to their corresponding MIME types
var MimeTypes = map[string]string{
	".jpg":  "image/jpeg",
	".jpeg": "image/jpeg",
	".png":  "image/png",
	".gif":  "image/gif",
	".pdf":  "application/pdf",
	".txt":  "text/plain",
	".html": "text/html",
	".htm":  "text/html",
	".json": "application/json",
	".xml":  "application/xml",
	".zip":  "application/zip",
	".doc":  "application/msword",
	".docx": "application/vnd.openxmlformats-officedocument.wordprocessingml.document",
	".xls":  "application/vnd.ms-excel",
	".xlsx": "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
	".csv":  "text/csv",
}

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

// ValidateFileType checks if a file matches the expected type (MIME type or extension)
func ValidateFileType(fs afero.Fs, filePath string, expectedType string) error {
	// Normalize expected type
	if !strings.HasPrefix(expectedType, ".") && !strings.Contains(expectedType, "/") {
		expectedType = "." + expectedType
	}

	// Open and read file for MIME type detection
	file, err := fs.Open(filePath)
	if err != nil {
		return fmt.Errorf("Cannot open file: %v", err)
	}
	defer file.Close()

	// Read first 512 bytes for MIME type detection
	buffer := make([]byte, 512)
	_, err = file.Read(buffer)
	if err != nil && err != io.EOF {
		return fmt.Errorf("Cannot read file: %v", err)
	}

	detectedType := http.DetectContentType(buffer)

	// Handle MIME type constraints
	if strings.Contains(expectedType, "/") {
		// Special handling for binary files
		if detectedType == "application/octet-stream" {
			// For binary files, trust the file extension more than the detected type
			if ext := strings.ToLower(filepath.Ext(filePath)); ext != "" {
				if _, ok := MimeTypes[ext]; ok {
					return nil // Accept if extension matches expected type
				}
			}
		}

		if !strings.HasPrefix(strings.ToLower(detectedType), strings.ToLower(expectedType)) {
			return fmt.Errorf("File is not of type %v (detected: %v)", expectedType, detectedType)
		}
		return nil
	}

	// Handle extension constraints
	ext := strings.ToLower(filepath.Ext(filePath))
	if ext == "" {
		return fmt.Errorf("File has no extension")
	}

	// Check if extension matches expected MIME type
	if expectedMime, ok := MimeTypes[strings.ToLower(expectedType)]; ok {
		if strings.HasPrefix(strings.ToLower(detectedType), strings.ToLower(expectedMime)) {
			return nil
		}
		return fmt.Errorf("File extension %v doesn't match content type (detected: %v, expected: %v)",
			ext, detectedType, expectedMime)
	}

	return fmt.Errorf("Unknown file type: %v", expectedType)
}
