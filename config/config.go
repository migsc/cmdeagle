package config

import (
	"embed"
	"fmt"
	"io"
	"os"

	"github.com/migsc/cmdeagle/file"
	"github.com/migsc/cmdeagle/types"

	"github.com/charmbracelet/log"
	"gopkg.in/yaml.v3"
)

// TODO: Convert all YAML field names from camelCase to kebab-case to follow Go community conventions
// Examples:
// - allowUnknownFlags -> allow-unknown-flags
// - globalFlags -> global-flags
// - logLevel -> log-level
// This is a breaking change that should be done in a major version update.
//
// TODO: Exception to kebab-case: Use CapitalCase for Cobra-compatible arg validations
// to match Cobra's function names, while keeping custom boolean logic (and/or/not) in kebab-case
// Examples:
// - MinimumNArgs (matches cobra.MinimumNArgs)
// - OnlyValidArgs (matches cobra.OnlyValidArgs)
// - and/or/not (our custom boolean logic)
// This makes it clear which validations map directly to Cobra functions.

//go:embed *
var PackageFS embed.FS

func LoadFromBundle(bundleFS embed.FS) ([]byte, *types.CmdeagleConfig, error) {
	log.Debug("Loading config file from embedded bundle")

	configFile, err := bundleFS.Open("config.cmd.yaml")
	if err != nil {
		return nil, nil, err
	}

	content, err := io.ReadAll(configFile)
	if err != nil {
		return nil, nil, err
	}

	config, err := Parse(content)
	if err != nil {
		return content, nil, err
	}

	return content, config, nil
}

func Load(dirPath string) ([]byte, *types.CmdeagleConfig, error) {
	log.Debug("Checking for config file in", "dir", dirPath)

	configFilePath, err := file.FindFileEndsWithPattern(dirPath, "cmd.yaml")
	if err != nil {
		return nil, nil, err
	}

	log.Debug("Found config file:", "path", configFilePath)

	content, err := os.ReadFile(configFilePath)
	if err != nil {
		return nil, nil, fmt.Errorf("Error reading config file %s: %v", configFilePath, err)
	}

	log.Debug("Loaded config file from:", "path", configFilePath)

	config, err := Parse(content)
	if err != nil {
		return content, nil, err
	}

	return content, config, nil
}

func Parse(content []byte) (*types.CmdeagleConfig, error) {
	log.Debug("Parsing config file content.")

	var config types.CmdeagleConfig
	err := yaml.Unmarshal(content, &config)
	if err != nil {
		return nil, fmt.Errorf("error parsing YAML: %v", err)
	}

	log.Debug("Parsed config file:",
		"name", config.Name,
		"version", config.Version,
	)

	for _, cmd := range config.Commands {
		log.Debug("Found top-level command definition in config file:", "name", cmd.Name)
	}

	return &config, nil
}
