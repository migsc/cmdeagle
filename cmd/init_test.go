package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/charmbracelet/log"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

var tempDir string
var err error

func TestInitCmd(t *testing.T) {
	// Set test environment variable to skip interactive prompts
	os.Setenv("GO_TEST", "1")
	defer os.Unsetenv("GO_TEST")

	// Create a temporary directory for testing
	tempDir, err = afero.TempDir(afero.NewOsFs(), "", "cmdeagle-tests")

	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Store original working directory
	originalWd, err := os.Getwd()
	log.Info("Original working directory:", "path", originalWd)
	if err != nil {
		t.Fatalf("Failed to get current working directory: %v", err)
	}

	// Change to temp directory for tests
	err = os.Chdir(tempDir)
	log.Info("Changed to temp directory:", "path", tempDir)
	if err != nil {
		t.Fatalf("Failed to change to temp directory: %v", err)
	}
	defer os.Chdir(originalWd)

	tests := []struct {
		name         string
		args         []string
		expectedName string
		expectError  bool
		setupFunc    func() error
		validateFunc func(t *testing.T, yamlContent []byte)
	}{
		{
			name:         "with provided project name",
			args:         []string{"init", "test-project"},
			expectedName: "test-project",
			validateFunc: func(t *testing.T, yamlContent []byte) {
				assert.Contains(t, string(yamlContent), "name: \"test-project\"")
			},
		},
		{
			name:         "using directory name",
			args:         []string{"init"},
			expectedName: "using_directory_name",
			validateFunc: func(t *testing.T, yamlContent []byte) {
				assert.Contains(t, string(yamlContent), "name: \"using_directory_name\"")
			},
		},
		{
			name:        "error when .cmd.yaml already exists",
			args:        []string{"init", "test-project"},
			expectError: true,
			setupFunc: func() error {
				return os.WriteFile(".cmd.yaml", []byte("existing content"), 0644)
			},
		},
		{
			name:        "error with invalid directory permissions",
			args:        []string{"init", "test-project"},
			expectError: true,
			setupFunc: func() error {
				// First create a subdirectory with no write permissions
				err := os.Mkdir("readonly", 0555)
				if err != nil {
					return err
				}
				return os.Chdir("readonly")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a new subdirectory for each test
			testSubDir := filepath.Join(tempDir, strings.ReplaceAll(tt.name, " ", "_"))
			log.Info("Creating test subdirectory:", "path", testSubDir)
			err := os.MkdirAll(testSubDir, 0755)
			assert.NoError(t, err)

			// Change to test directory
			err = os.Chdir(testSubDir)
			assert.NoError(t, err)

			// Ensure we clean up after the test
			defer func() {
				// Change back to temp directory
				err = os.Chdir(tempDir)
				assert.NoError(t, err)

				// If this was the permissions test, we need to fix permissions to clean up
				if tt.name == "error with invalid directory permissions" {
					// Change back to parent dir to modify readonly dir
					readonlyDir := filepath.Join(testSubDir, "readonly")
					err = os.Chmod(readonlyDir, 0755)
					assert.NoError(t, err)
				}

				// Remove the test subdirectory
				err = os.RemoveAll(testSubDir)
				assert.NoError(t, err)
			}()

			// Run setup if provided
			if tt.setupFunc != nil {
				err := tt.setupFunc()
				if err != nil {
					t.Fatalf("Setup failed: %v", err)
				}
			}

			// Execute the init command
			log.Info("Using args:", "args", tt.args)
			rootCmd.SetArgs(tt.args)
			err = rootCmd.Execute()

			if tt.expectError {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)

			// Verify the .cmd.yaml file was created
			yamlContent, err := os.ReadFile(".cmd.yaml")
			assert.NoError(t, err)
			assert.NotEmpty(t, yamlContent)

			// Run custom validation if provided
			if tt.validateFunc != nil {
				tt.validateFunc(t, yamlContent)
			}
		})
	}
}
