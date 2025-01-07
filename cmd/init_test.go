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
		// {
		// 	name:         "using directory name",
		// 	args:         []string{},
		// 	expectedName: filepath.Base(tempDir),
		// 	validateFunc: func(t *testing.T, yamlContent []byte) {
		// 		assert.Contains(t, string(yamlContent), "name: \""+filepath.Base(tempDir)+"\"")
		// 	},
		// },
		// {
		// 	name:        "error when .cmd.yaml already exists",
		// 	args:        []string{"test-project"},
		// 	expectError: true,
		// 	setupFunc: func() error {
		// 		return os.WriteFile(".cmd.yaml", []byte("existing content"), 0644)
		// 	},
		// },
		// {
		// 	name:        "error with invalid directory permissions",
		// 	args:        []string{"test-project"},
		// 	expectError: true,
		// 	setupFunc: func() error {
		// 		return os.Chmod(tempDir, 0444) // Read-only permissions
		// 	},
		// },
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

			// Change back to temp directory for next test
			err = os.Chdir(tempDir)
			assert.NoError(t, err)
		})
	}
}
