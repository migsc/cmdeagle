package executable

import (
	"bytes"
	"embed"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/migsc/cmdeagle/file"

	log "github.com/charmbracelet/log"
)

//go:embed *
var PackageFS embed.FS

var SrcDirPath string

var packageSrcDirPath string

func GetPackageSrcDirPath() string {
	log.Debug("Getting package src dir path:")
	if packageSrcDirPath != "" {
		log.Debug("Package src dir path already set:", "path", packageSrcDirPath)
		return packageSrcDirPath
	}

	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		log.Fatal("failed to get current file path")
	}

	log.Debug("Package src dir path not set, setting it now:", "path", filepath.Dir(filename))

	packageSrcDirPath = filepath.Dir(filename)

	return packageSrcDirPath
}

// BuildEnv holds build environment configuration
type BuildEnv struct {
	CGOEnabled string
	GOOS       string
	GOARCH     string
}

var DefaultBuildEnv = BuildEnv{
	CGOEnabled: "0",
	GOOS:       "linux",
	GOARCH:     "amd64",
}

func init() {
	// _, filename, _, ok := runtime.Caller(0)
	// if !ok {
	// 	log.Fatalf("Failed to get the source directory for the executable.")
	// 	os.Exit(1)
	// }

	// SrcDirPath = filepath.Join(filepath.Dir(filename))

	// fmt.Printf("Using source directory: %s\n", SrcDirPath)
}

func GetDestDir() (string, error) {
	// Determine default binary location based on OS
	var binDir string
	switch runtime.GOOS {
	case "windows":
		// Use %LOCALAPPDATA%\Programs
		localAppData := os.Getenv("LOCALAPPDATA")
		if localAppData == "" {
			return "", fmt.Errorf("LOCALAPPDATA environment variable not set")
		}
		binDir = filepath.Join(localAppData, "Programs")
	default: // linux, darwin, etc.
		// Try /usr/local/bin if we have write permissions, otherwise use ~/.local/bin
		if _, err := os.Stat("/usr/local/bin"); err == nil {
			if isWriteable("/usr/local/bin") {
				binDir = "/usr/local/bin"
				break
			}
		}
		// Fallback to user's local bin directory
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("could not determine user home directory: %w", err)
		}
		binDir = filepath.Join(homeDir, ".local", "bin")
	}

	// Ensure the directory exists
	if err := os.MkdirAll(binDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create binary directory: %w", err)
	}

	log.Debug("Using binary directory:", "path", binDir)
	return binDir, nil
}

// Helper function to check if a directory is writeable
func isWriteable(path string) bool {
	// Try to create a temporary file
	tmpFile := filepath.Join(path, ".write_test")
	err := os.WriteFile(tmpFile, []byte{}, 0644)
	if err != nil {
		return false
	}
	os.Remove(tmpFile)
	return true
}

func BuildBinary(mainFilePath string, targetBinaryPath string, moduleName string) error {
	expandedTargetPath, err := file.ExpandPath(targetBinaryPath)
	if err != nil {
		return fmt.Errorf("failed to expand bin-path: %w", err)
	}

	log.Info("Building CLI",
		"from", mainFilePath,
		"to", expandedTargetPath,
	)

	// Initialize the module
	initCmd := exec.Command("go", "mod", "init", moduleName)
	initCmd.Dir = filepath.Dir(mainFilePath)
	initCmd.Stdout = os.Stdout
	initCmd.Stderr = os.Stderr

	if err := initCmd.Run(); err != nil {
		return fmt.Errorf("error initializing module: %w", err)
	}

	// Run go mod tidy
	tidyCmd := exec.Command("go", "mod", "tidy")
	tidyCmd.Dir = filepath.Dir(mainFilePath)
	tidyCmd.Stdout = os.Stdout
	tidyCmd.Stderr = os.Stderr

	if err := tidyCmd.Run(); err != nil {
		return fmt.Errorf("error running go mod tidy: %w", err)
	}

	// Build the binary - explicitly specify the main package
	buildCmd := exec.Command("go", "build", "-o", expandedTargetPath, ".")
	buildCmd.Dir = filepath.Dir(mainFilePath)

	// Set up environment variables for the build
	// Always keep CGO_ENABLED=0 as it's crucial for the functionality
	env := append(os.Environ(),
		"CGO_ENABLED=0",
		fmt.Sprintf("GOOS=%s", DefaultBuildEnv.GOOS),
		fmt.Sprintf("GOARCH=%s", DefaultBuildEnv.GOARCH),
	)

	log.Debug("Building with environment",
		"CGO_ENABLED", "0",
		"GOOS", DefaultBuildEnv.GOOS,
		"GOARCH", DefaultBuildEnv.GOARCH,
	)

	buildCmd.Env = env
	buildCmd.Stdout = os.Stdout
	buildCmd.Stderr = os.Stderr

	if err := buildCmd.Run(); err != nil {
		return fmt.Errorf("error building binary: %w", err)
	}

	return nil
}

func CheckDependencyExists(name string) (string, bool) {
	path, err := exec.LookPath(name)
	if err != nil {
		return "", false
	}
	return path, true
}

// RunCommand executes a command with the given arguments and returns its output
func RunCommand(binaryPath string, args []string) (string, error) {
	cmd := exec.Command(binaryPath, args...)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("failed to run command %s %v: %v", binaryPath, args, err)
	}
	return strings.TrimSpace(out.String()), nil
}
