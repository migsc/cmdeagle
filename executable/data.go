package executable

import (
	"os"
	"path/filepath"
	"runtime"

	"github.com/charmbracelet/log"
)

func GetAppDataDir(appName string) string {
	dataDir := getOSDataDir()
	appDataDir := filepath.Join(dataDir, appName)

	return appDataDir
}

// getOSDataDir returns the appropriate user-specific data directory for cmdeagle
func getOSDataDir() string {
	// Get user home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		// Fallback to temp directory if we can't get home
		log.Warn("Could not determine user home directory, falling back to temp dir", "error", err)
		return filepath.Join(os.TempDir(), "cmdeagle")
	}

	switch runtime.GOOS {
	case "windows":
		// Windows: %LocalAppData%\cmdeagle
		if localAppData := os.Getenv("LocalAppData"); localAppData != "" {
			return filepath.Join(localAppData, "cmdeagle")
		}
		return filepath.Join(homeDir, "AppData", "Local", "cmdeagle")

	case "darwin":
		// macOS: ~/Library/Application Support/cmdeagle
		return filepath.Join(homeDir, "Library", "Application Support", "cmdeagle")

	default:
		// Linux/Unix: ~/.local/share/cmdeagle
		if xdgData := os.Getenv("XDG_DATA_HOME"); xdgData != "" {
			return filepath.Join(xdgData, "cmdeagle")
		}
		return filepath.Join(homeDir, ".local", "share", "cmdeagle")
	}
}
