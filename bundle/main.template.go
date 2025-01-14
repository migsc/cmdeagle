package bundle

import (
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/migsc/cmdeagle/args"
	"github.com/migsc/cmdeagle/config"
	"github.com/migsc/cmdeagle/executable"
	"github.com/migsc/cmdeagle/flags"
	"github.com/migsc/cmdeagle/types"

	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
)

//go:embed *
var bundleFS embed.FS

// var config schema.CmdeagleConfig
var rootCmd *cobra.Command

var cobraCommands = make(map[string]*cobra.Command)

var LOG_LEVEL = log.InfoLevel

func init() {
	// DEBUG_MODE will be replaced during build
	log.SetLevel(LOG_LEVEL)
	log.SetFormatter(log.TextFormatter)
	cobra.EnableTraverseRunHooks = true
}

func main_template() {
	if err := execute(); err != nil {
		log.Fatal("execution failed", "error", err)
		os.Exit(1)
	}
}

func execute() error {
	log.Info("Executing...")

	var err error
	var cmdConfig *types.CmdeagleConfig

	_, cmdConfig, err = config.LoadFromBundle(bundleFS)
	if err != nil {
		return fmt.Errorf("Failed to load configuration from embedded bundle: %w", err)
	}

	rootCommandDef := &types.CommandDefinition{
		Name:        cmdConfig.Name,
		Description: cmdConfig.Description,
		Aliases:     []string{},
		Args:        cmdConfig.Args,
		Flags:       cmdConfig.Flags,
		Requires:    cmdConfig.Requires,
		Includes:    cmdConfig.Includes,
		Build:       cmdConfig.Build,
		Start:       cmdConfig.Start,
	}

	// Set up the root command
	rootCmd, err = setupCommand(cmdConfig, rootCommandDef, nil, []string{})

	if err != nil {
		return fmt.Errorf("failed to setup root command: %w", err)
	}

	rootCmd.Version = cmdConfig.Version

	// Set up all other subcommands
	visitor := &RunnerCommandVisitor{config: cmdConfig}
	if err := config.WalkCommands(&cmdConfig.Commands, nil, visitor, []string{}); err != nil {
		return fmt.Errorf("failed to process commands: %w", err)
	}

	log.Info("Inspecting embedded bundle filesystem...")

	// Walk through the embedded filesystem and print the directory tree
	// err = fs.WalkDir(bundleFS, ".", func(path string, d fs.DirEntry, err error) error {
	// 	if err != nil {
	// 		return err
	// 	}

	// 	indent := strings.Repeat("  ", len(strings.Split(path, "/")))
	// 	if d.IsDir() {
	// 		fmt.Println(indent + "📁 " + d.Name())
	// 	} else {
	// 		fmt.Println(indent + "📄 " + d.Name())
	// 	}

	// 	return nil
	// })

	// if err != nil {
	// 	return fmt.Errorf("failed to walk embedded filesystem: %w", err)
	// }

	log.Debug("Setting up data directory")
	if err := setupDataDirectory(bundleFS, cmdConfig.Name); err != nil {
		log.Fatalf("Failed to setup data directory: %v", err)
	}

	// log.Debug("Command tree structure:")
	// var printCommandTree func(cmd *cobra.Command, level int)
	// printCommandTree = func(cmd *cobra.Command, level int) {
	// 	indent := strings.Repeat("  ", level)
	// 	// fmt.Printf("%s%s\n", indent, cmd.Name())
	// 	for _, subCmd := range cmd.Commands() {
	// 		printCommandTree(subCmd, level+1)
	// 	}
	// }
	// printCommandTree(rootCmd, 0)

	log.Debug("Done")

	return rootCmd.Execute()
}

// RunnerCommandVisitor handles the command processing during runtime
type RunnerCommandVisitor struct {
	config *types.CmdeagleConfig
}

func (visitor *RunnerCommandVisitor) Visit(commandDef *types.CommandDefinition, parent *types.CommandDefinition, path []string) error {
	log.Debug("Visiting command", "name", commandDef.Name, "path", path)
	_, err := setupSubCommand(visitor.config, commandDef, parent, path)
	return err
}

func getCommandPath(parts ...string) string {
	log.Debug("Getting command path", "parts", parts)
	return strings.Join(parts, ":")
}

func setupSubCommand(cmdConfig *types.CmdeagleConfig, commandDef *types.CommandDefinition, parent *types.CommandDefinition, path []string) (*cobra.Command, error) {
	log.Debug("Setting up subcommand", "name", commandDef.Name, "path", path)
	cobraCmd, err := setupCommand(cmdConfig, commandDef, parent, path)
	if err != nil {
		return nil, err
	}

	// Add to parent
	if parent == nil {
		log.Debug("Adding command to root", "name", commandDef.Name, "path", path)
		rootCmd.AddCommand(cobraCmd)
	} else {
		log.Debug("Adding command to parent", "name", commandDef.Name, "path", path)
		cobraCmd.InheritedFlags()
		parentPath := getCommandPath(path[:len(path)-1]...)
		if parentCmd, exists := cobraCommands[parentPath]; exists {
			parentCmd.AddCommand(cobraCmd)
		}
	}

	// Store command in map using namespace path
	cmdPath := getCommandPath(path...)
	cobraCommands[cmdPath] = cobraCmd

	return cobraCmd, nil
}

func setupCommand(cmdConfig *types.CmdeagleConfig, commandDef *types.CommandDefinition, parent *types.CommandDefinition, path []string) (*cobra.Command, error) {
	appDataDirPath := executable.GetAppDataDir(cmdConfig.Name)
	commandPath := filepath.Join(appDataDirPath, filepath.Join(path...))

	// TODO: This is how we want to organize our logic now. we set up with cobra's lifecycle hooks and then we
	// validate the args and flags and then we run the command. We need to have...
	// 1. A way to define the command args and flags, etc
	// 2. A way to validate the args and flags, and probably we may put the "requires" configuration in here?
	// 2. A way to interpolate the args and flags into the start script
	// 3. A way to run the command

	cobraCmd := &cobra.Command{
		Use:          commandDef.Name,
		Short:        commandDef.Description,
		SilenceUsage: true, // Silence usage on validation errors
		RunE: func(cmd *cobra.Command, args []string) error {
			// Show help if it's a leaf command with no args
			if len(args) == 0 && len(cmd.Commands()) == 0 {
				return cmd.Help()
			}

			// If there's no start script, just show help
			if commandDef.Start == "" || parent == nil {
				return cmd.Help()
			}

			// Continue with normal command execution
			return nil
		},
	}

	if len(commandDef.Aliases) > 0 {
		cobraCmd.Aliases = commandDef.Aliases
	}

	// Create flag store
	flagStore := flags.CreateFlagsStore(cobraCmd, commandDef)
	log.Debug("Created flagStore", "path", commandPath, "flagStore", flagStore)

	// 1. Global setup for the entire top-level command.
	cobraCmd.PersistentPreRunE = func(cobraCmd *cobra.Command, args []string) error {
		log.Debug("PersistentPreRunE / Triggering hook", "path", commandPath)
		log.Debug("Requires", "path", commandPath, "requires", commandDef.Requires)
		if commandDef.Requires != nil {
			for name, versionDeclared := range commandDef.Requires {
				path, exists := executable.CheckDependencyExists(name)

				if !exists {
					return fmt.Errorf("dependency %s not found. You need to install it to run this command.", name)
				}

				if versionDeclared == "*" {
					continue
				}

				versionFound, err := executable.GetVersion(path, []string{})
				if err != nil {
					return fmt.Errorf("failed to get version for dependency %s: %w", name, err)
				}

				log.Debug("Dependency version", "name", name, "versionFound", versionFound, "versionDeclared", versionDeclared)

				matchesRequirement, reason := executable.CheckVersionCompatibility(versionFound, versionDeclared)
				if !matchesRequirement {
					return fmt.Errorf("dependency %s version %s does not meet the required version %s: %s", name, versionFound, versionDeclared, reason)
				}

			}
		}
		log.Debug("Triggering hook `PersistentPreRunE`", "path", commandPath)
		return nil
	}

	// 2. Parse, validate and load arguments and flags..
	var argStore *args.ArgsStateStore
	var paramsStore *config.ParamsStateStore
	// var flagStore *state.FlagsStateStore

	cobraCmd.Args = func(cobraCommand *cobra.Command, arguments []string) error {
		log.Debug("Triggering hook `Args`", "path", commandPath)
		argStore = args.CreateArgsStore(cobraCmd, &commandDef.Args, arguments)
		log.Debug("Created argsStore", "path", commandPath, "argsStore", argStore)

		paramsStore = config.CreateParamsStore(argStore, flagStore)
		log.Debug("Created paramsStore", "path", commandPath, "paramsStore", paramsStore)

		err := args.ValidateArgs(cobraCmd, &commandDef.Args, argStore)
		if err != nil {
			return err
		}

		err = flags.ValidateFlags(cobraCmd, commandDef.Flags, flagStore)
		if err != nil {
			return err
		}

		return nil
	}

	// 3. Interpolate args and flags into the start script
	cobraCmd.PreRunE = func(cobraCmd *cobra.Command, args []string) error {
		log.Debug("Triggering hook `PreRunE`", "path", commandPath)
		return nil
	}

	// 4. Run the start script
	cobraCmd.RunE = func(cobraCmd *cobra.Command, args []string) error {
		log.Debug("Run / Triggering hook", "path", commandPath)

		if commandDef.Start == "" {
			log.Debug("No start script defined for command", "path", commandPath)
			return nil
		}

		// Interpolate both args and flags into the start script
		script := commandDef.Start

		// Interpolate args and flags
		script = argStore.Interpolate(script)
		script = flagStore.Interpolate(script)
		script = paramsStore.Interpolate(script)

		// TODO: Useful?
		// cmdVisitor.paramStore.Set("LOG_LEVEL", "debug")
		paramsStore.Set("cli.name", cmdConfig.Name)
		// cmdVisitor.paramStore.Set("CMD_VERSION", cmdConfig.Version)
		// cmdVisitor.paramStore.Set("CMD_DESCRIPTION", cmdConfig.Description)
		// cmdVisitor.paramStore.Set("CMD_AUTHOR", cmdConfig.Author)
		// cmdVisitor.paramStore.Set("CMD_LICENSE", cmdConfig.License)
		binDirPath, err := executable.GetDestDir()
		if err != nil {
			return fmt.Errorf("failed to get binary directory: %w", err)
		}
		paramsStore.Set("cli.bin_dir", binDirPath)
		paramsStore.Set("cli.data_dir", appDataDirPath)

		// 4. Run the command
		log.Debug("Run / Running start script for", "path", commandPath, "script", script)

		execCmd := exec.Command("sh", "-c", script)

		// Copy the current environment and add new variables iteratively
		envVars := append(paramsStore.GetEnvVariables(), append(argStore.GetEnvVariables(), flagStore.GetEnvVariables()...)...)
		execCmd.Env = os.Environ() // Start with the current environment
		for _, env := range envVars {
			log.Debug("Run / Setting environment variable", "path", commandPath, "env", env.Name+"="+env.Value)
			execCmd.Env = append(execCmd.Env, env.Name+"="+env.Value)
		}

		// We don't want to run the command in the command's directory, we want to run it in the root command's directory
		// This is because the command may need to access files in the root command's directory
		// We don't have a way to opt out of this yet, but we can change this later
		// execCmd.Dir = commandPath
		// execCmd.Dir = filepath.Join(runnerSrcDistDir, cmdConfig.Name)
		execCmd.Dir = appDataDirPath
		execCmd.Stdout = os.Stdout
		execCmd.Stderr = os.Stderr

		log.Debug("#########OUTPUT############")
		if err := execCmd.Run(); err != nil {
			return err
		}
		log.Debug("###########################")

		return nil
	}

	cobraCmd.PostRunE = func(cobraCmd *cobra.Command, args []string) error {
		log.Debug("PostRunE / Triggering hook", "path", commandPath)
		return nil
	}

	// 1. Global teardown for the entire top-level command.
	// We should eventaully do things here like
	// - Evaluate global settings
	// - Handle JIT bundle extraction
	// - Handle dependencies specified via `requires` settings
	cobraCmd.PersistentPostRunE = func(cobraCmd *cobra.Command, args []string) error {
		log.Debug("PersistentPostRunE / Triggering hook", "path", commandPath)
		return nil
	}

	// Allow help to bypass SilenceUsage
	// cobraCmd.SetHelpFunc(func(cmd *cobra.Command, args []string) {
	// 	cmd.SilenceUsage = false
	// 	cmd.Parent().HelpFunc()(cmd, args)
	// })

	return cobraCmd, nil
}

func setupDataDirectory(embeddedFS embed.FS, appName string) error {
	// Get the app-specific data directory
	appDataDir := executable.GetAppDataDir(appName)

	// Check if manifest already exists
	if _, err := os.Stat(filepath.Join(appDataDir, ".manifest.json")); err == nil {
		log.Debug("Manifest already exists, skipping bundle extraction")
		return nil
	}

	// Create the directory if it doesn't exist
	if err := os.MkdirAll(appDataDir, 0755); err != nil {
		return fmt.Errorf("failed to create data directory: %w", err)
	}

	// Initialize manifest to track all files
	manifest := struct {
		Files []string          `json:"files"`
		Time  time.Time         `json:"created_at"`
		Meta  map[string]string `json:"meta"`
	}{
		Files: []string{},
		Time:  time.Now(),
		Meta:  map[string]string{"app": appName},
	}

	// Walk through embedded filesystem and copy files
	err := fs.WalkDir(embeddedFS, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Skip root directory
		if path == "." {
			return nil
		}

		targetPath := filepath.Join(appDataDir, path)

		// Create directories
		if d.IsDir() {
			if err := os.MkdirAll(targetPath, 0755); err != nil {
				return fmt.Errorf("failed to create directory %s: %w", targetPath, err)
			}
			// // Print directory creation (preserving existing logging)
			// indent := strings.Repeat("  ", len(strings.Split(path, "/")))
			// fmt.Println(indent + "📁 " + d.Name())
			return nil
		}

		// Read and write files
		data, err := embeddedFS.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to read embedded file %s: %w", path, err)
		}

		if err := os.WriteFile(targetPath, data, 0644); err != nil {
			return fmt.Errorf("failed to write file %s: %w", targetPath, err)
		}

		// Add to manifest
		manifest.Files = append(manifest.Files, path)

		// // Print file creation (preserving existing logging)
		// indent := strings.Repeat("  ", len(strings.Split(path, "/")))
		// fmt.Println(indent + "📄 " + d.Name())

		return nil
	})

	if err != nil {
		return fmt.Errorf("failed to copy embedded files: %w", err)
	}

	// Write manifest file
	manifestData, err := json.MarshalIndent(manifest, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to create manifest: %w", err)
	}

	manifestPath := filepath.Join(appDataDir, ".manifest.json")
	if err := os.WriteFile(manifestPath, manifestData, 0644); err != nil {
		return fmt.Errorf("failed to write manifest: %w", err)
	}

	return nil
}
