package bundle

import (
	"fmt"
	"os"
	"path/filepath"
)

// ExtractBundle extracts all files for a command
func ExtractBundle(manifest *BundleManifest, targetDir string) error {
	for _, file := range manifest.Files {
		targetPath := filepath.Join(targetDir, file.Path)

		if file.IsDir {
			if err := os.MkdirAll(targetPath, file.Mode); err != nil {
				return fmt.Errorf("failed to create directory %s: %w", targetPath, err)
			}
			continue
		}

		if err := os.WriteFile(targetPath, file.Content, file.Mode); err != nil {
			return fmt.Errorf("failed to write file %s: %w", targetPath, err)
		}
	}
	return nil
}

// ExtractCommandFiles extracts only files needed for a specific command
func ExtractCommandFiles(manifest *BundleManifest, commandPath string, targetDir string) error {
	if manifest.CommandPath != commandPath {
		return fmt.Errorf("manifest is for command %s, not %s", manifest.CommandPath, commandPath)
	}
	return ExtractBundle(manifest, targetDir)
}
