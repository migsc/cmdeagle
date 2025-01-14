/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bytes"
	"embed"
	_ "embed"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/migsc/cmdeagle/bundle"
	"github.com/migsc/cmdeagle/config"
	"github.com/migsc/cmdeagle/envvar"
	"github.com/migsc/cmdeagle/file"
	"github.com/migsc/cmdeagle/flags"
	"github.com/migsc/cmdeagle/params"

	"github.com/migsc/cmdeagle/args"
	"github.com/migsc/cmdeagle/executable"
	"github.com/migsc/cmdeagle/types"

	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

const (
	defaultPermissions = 0644
	configFileName     = "config.cmd.yaml"
	mainFileName       = "main.go"
	binDirName         = "bin"
)

// buildCmd represents the build command
var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Build your CLI from your cmd.yaml configuration.",
	Long: `Finds the configuration file in the current directory 
	and builds the CLI.`,
	Run: func(cmd *cobra.Command, arguments []string) {
		if err := runBuild(); err != nil {
			fmt.Printf("Build failed: %v\n", err)
			os.Exit(1)
		}
	},
}

var flagSet *pflag.FlagSet

func init() {
	// Add debug flag

	log.SetLevel(log.InfoLevel)
	buildCmd.Flags().Bool("debug", false, "Enable debug logging in both build and generated CLI")

	log.SetFormatter(log.TextFormatter)

	rootCmd.AddCommand(buildCmd)

	buildCmd.Flags().StringVar(&executable.DefaultBuildEnv.GOOS, "os", runtime.GOOS, "Target operating system")
	buildCmd.Flags().StringVar(&executable.DefaultBuildEnv.GOARCH, "arch", runtime.GOARCH, "Target architecture")

	// Add experimental-imports flag
	buildCmd.Flags().Bool("experimental-imports", false, "Enable experimental import resolution")

	// Add verbose flag and connect it to log level
	// buildCmd.Flags().Bool("verbose", false, "Enable verbose logging")
	// buildCmd.PreRun = func(cmd *cobra.Command, args []string) {
	// 	verbose, _ := cmd.Flags().GetBool("verbose")
	// 	if verbose {
	// 		log.SetLevel(log.DebugLevel)
	// 	}
	// }

	flagSet = buildCmd.Flags()
}

// TODO: Hacky way to get the bundle staging directory path to be used by the build command visitor
var bundleStagingDirPath string

type Package struct {
	FS   embed.FS
	Name string
}

var packages = []Package{
	{FS: args.PackageFS, Name: "args"},
	{FS: bundle.PackageFS, Name: "bundle"},
	{FS: config.PackageFS, Name: "config"},
	{FS: envvar.PackageFS, Name: "envvar"},
	{FS: executable.PackageFS, Name: "executable"},
	{FS: file.PackageFS, Name: "file"},
	{FS: flags.PackageFS, Name: "flags"},
	{FS: params.PackageFS, Name: "params"},
	{FS: types.PackageFS, Name: "types"},
}

func runBuild() error {
	var err error

	// Get debug flag value
	// TODO: Move this lower into the appropriate place in the code
	debug, _ := flagSet.GetBool("debug")
	if debug {
		log.SetLevel(log.DebugLevel)
		bundle.MainTemplateReplacements["var LOG_LEVEL = log.InfoLevel"] = []byte("var LOG_LEVEL = log.DebugLevel")
	}

	// TODO we should probably allow the user to specify the working directory
	workingDirPath, err := os.Getwd()
	if err != nil {
		return err
	}
	log.Debug("Using working directory", "path", workingDirPath)

	// Set up a directory for the bundle to be assembled in.

	bundleStagingDirPath, err = bundle.SetupStagingDir()
	if err != nil {
		return err
	}

	// // Get the relative path from main.go to the bundle staging directory
	// mainFileDir := bundle.GetPackageSrcDirPath()
	// relativeBundlePath, err := filepath.Rel(mainFileDir, bundleStagingDirPath)
	// if err != nil {
	// 	return fmt.Errorf("failed to get relative path: %w", err)
	// }

	// Load the config file, parse it, and write it out to the bundle staging directory
	log.Debug("Config file loaded from", "path", workingDirPath)
	configFileContent, cmdConfig, err := config.Load(workingDirPath)
	if err != nil {
		return err
	}

	// Print experimental-imports flag value
	experimentalImports, _ := flagSet.GetBool("experimental-imports")

	if experimentalImports {
		log.Warn("Experimental imports enabled")
		// 1. Determine target OS
		goos := runtime.GOOS
		resolverBinaryName := "resolve-imports-"
		if goos == "windows" {
			resolverBinaryName += "win.exe"
		} else if goos == "darwin" {
			resolverBinaryName += "macos"
		} else {
			resolverBinaryName += "linux"
		}

		// 2. Run the resolver binary
		resolverPath := filepath.Join("bin", resolverBinaryName)
		cmd := exec.Command(resolverPath)
		cmd.Dir = workingDirPath // Run in the same directory as the config file

		// 3. Capture the output
		var stdout bytes.Buffer
		cmd.Stdout = &stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to run import resolver: %w", err)
		}

		// 4. Parse the resolved config
		configFileContent = stdout.Bytes()
		log.Info("Resolved config", "content", string(configFileContent))
		cmdConfig, err = config.Parse(configFileContent)
		if err != nil {
			return fmt.Errorf("failed to parse resolved config: %w", err)
		}
	}

	outFile := filepath.Join(bundleStagingDirPath, "config.cmd.yaml")
	err = os.WriteFile(outFile, configFileContent, 0644) // TODO do we need to preserve original permissions?
	if err != nil {
		return fmt.Errorf("error writing %s: %w", outFile, err)
	}

	// There's no need to setup an empty directory for the binary to be built in because we expect it to be
	// shared with other binaries already on the system.
	var binDirPath string
	binDirPath, err = executable.GetDestDir()
	if err != nil {
		return err
	}

	// Create a visitor to handle the build-specific command processing in a recursive manner.
	cmdVisitor := BuildCommandVisitor{
		config:   cmdConfig,
		baseDir:  bundleStagingDirPath,
		envStore: envvar.CreateEnvStore(),
	}

	// TODO: Useful?
	// cmdVisitor.paramStore.Set("LOG_LEVEL", "debug")
	cmdVisitor.envStore.Set("cli.name", cmdConfig.Name)
	// cmdVisitor.paramStore.Set("CMD_VERSION", cmdConfig.Version)
	// cmdVisitor.paramStore.Set("CMD_DESCRIPTION", cmdConfig.Description)
	// cmdVisitor.paramStore.Set("CMD_AUTHOR", cmdConfig.Author)
	// cmdVisitor.paramStore.Set("CMD_LICENSE", cmdConfig.License)
	cmdVisitor.envStore.Set("cli.bin_dir", binDirPath)
	cmdVisitor.envStore.Set("cli.data_dir", executable.GetAppDataDir(cmdConfig.Name))

	// First we build the root command by creating a command definition from the root configuration.
	rootCommandDef := &types.CommandDefinition{
		Name:        cmdConfig.Name,
		Description: cmdConfig.Description,
		Commands:    cmdConfig.Commands,
		Flags:       cmdConfig.Flags,
		Requires:    cmdConfig.Requires,
		Includes:    cmdConfig.Includes,
		Build:       cmdConfig.Build,
		Start:       cmdConfig.Start,
	}
	err = cmdVisitor.Build(rootCommandDef, nil, []string{})
	if err != nil {
		return fmt.Errorf("Failed to build root command: %w", err)
	}

	// Then we visit each of the subcommands and build them
	if len(cmdConfig.Commands) > 0 {
		if err := config.WalkCommands(&cmdConfig.Commands, nil, &cmdVisitor, []string{}); err != nil {
			return fmt.Errorf("failed to process commands: %w", err)
		}
	}

	// We need to do some code generation to create a main.go file from an existing static template in order to
	// leverage Go's embed features to embded the bundle into the binary.
	templateMainFileContent, err := bundle.GetMainTemplateContent()
	if err != nil {
		return err
	}

	bundle.MainTemplateReplacements["github.com/migsc/cmdeagle/"] = []byte(fmt.Sprintf("%s/", cmdConfig.Name))

	resultingMainFileContent := bundle.InterpolateMainContent(templateMainFileContent)

	resultingMainFilePath := filepath.Join(bundleStagingDirPath, "main.go")
	// resultingMainFilePath := filepath.Join(executable.GetPackageSrcDirPath(), "main.go")

	// Write the resulting main.go file to the bundle staging directory
	err = os.WriteFile(resultingMainFilePath, resultingMainFileContent, 0644) // TODO do we need to preserve original permissions?
	if err != nil {
		return fmt.Errorf("error writing %s: %w", resultingMainFilePath, err)
	}
	log.Debug("Wrote resulting main.go file to", "path", resultingMainFilePath)

	// And now we need to also copy over every package that the executable depends on.
	for _, pkg := range packages {
		fs.WalkDir(pkg.FS, ".", func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}

			// Get source file info
			srcInfo, err := d.Info()
			if err != nil {
				return fmt.Errorf("failed to get source file info: %w", err)
			}

			// Create target path
			targetPath := filepath.Join(bundleStagingDirPath, pkg.Name, path)
			log.Debug("target path", "path", targetPath)

			// If it's a directory, create it
			if d.IsDir() {
				return nil
			}

			// Read source file
			data, err := pkg.FS.ReadFile(path)
			if err != nil {
				return fmt.Errorf("failed to read file %s: %w", path, err)
			}

			// Replace module name in file contents
			data = bytes.ReplaceAll(data, []byte("github.com/migsc/cmdeagle/"), []byte(fmt.Sprintf("%s/", cmdConfig.Name)))

			// Create parent directories if needed
			if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
				return fmt.Errorf("failed to create parent directories for %s: %w", targetPath, err)
			}

			// Write file to target location with modified contents
			if err := os.WriteFile(targetPath, data, srcInfo.Mode()); err != nil {
				return fmt.Errorf("failed to write file %s: %w", targetPath, err)
			}

			return nil
		})
	}

	// Finally we build the binary
	log.Debug("Preparing to build binary with", "binDirPath", binDirPath, "targetBinaryPath")
	targetBinaryPath := filepath.Join(binDirPath, cmdConfig.Name)
	err = executable.BuildBinary(resultingMainFilePath, targetBinaryPath, cmdConfig.Name)
	if err != nil {
		return err
	}

	log.Info("CLI built successfully", "location", binDirPath)

	// Delete the existing bundle data directory
	err = os.RemoveAll(executable.GetAppDataDir(cmdConfig.Name))
	if err != nil {
		return fmt.Errorf("failed to delete existing bundle data directory: %w", err)
	}

	return nil
}

// BuildCommandVisitor handles the build-specific command processing
// during command tree traversal
type BuildCommandVisitor struct {
	config   *types.CmdeagleConfig
	baseDir  string
	envStore *envvar.EnvStateStore
}

func (v *BuildCommandVisitor) Visit(commandDef *types.CommandDefinition, parent *types.CommandDefinition, path []string) error {
	return v.Build(commandDef, parent, path)
}

func (v *BuildCommandVisitor) Build(commandDef *types.CommandDefinition, parent *types.CommandDefinition, path []string) error {
	log.Debug("Building command",
		"name", commandDef.Name,
		"bundle-path", strings.Join(path, "/"),
	)

	commandPath := filepath.Join(v.baseDir, filepath.Join(path...))
	log.Debug("target path", "path", commandPath)

	if commandDef.Build != "" {
		log.Info("Running build script",
			"command", commandDef.Name,
			"script", commandDef.Build,
		)

		// Interpolate params such as environment variables
		script := commandDef.Build
		script = v.envStore.Interpolate(script)

		// TODO: Highly inefficient, but it works for now

		// TODO: Should probably allow the user to specify the shell to run the build script in.
		cmd := exec.Command("sh", "-c", script)

		// Copy the current environment and add new variables iteratively
		envVars := v.envStore.GetEnvVariables()
		cmd.Env = os.Environ() // Start with the current environment
		for _, env := range envVars {
			log.Debug("Run / Setting environment variable", "path", commandPath, "env", env.Name+"="+env.Value)
			cmd.Env = append(cmd.Env, env.Name+"="+env.Value)
		}

		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil {
			return err
		}
	}

	if err := bundle.CopyIncludedFiles(v.config, commandDef, path, bundleStagingDirPath); err != nil {
		return err
	}

	return nil
}
