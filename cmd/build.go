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
)

const (
	defaultPermissions = 0644
	configFileName     = "config.cmd.yaml"
	mainFileName       = "main.go"
	binDirName         = "bin"
)

func init() {
	log.SetLevel(log.InfoLevel)
	log.SetFormatter(log.TextFormatter)

	rootCmd.AddCommand(buildCmd)

	buildCmd.Flags().StringVar(&executable.DefaultBuildEnv.GOOS, "os", "linux", "Target operating system")
	buildCmd.Flags().StringVar(&executable.DefaultBuildEnv.GOARCH, "arch", "amd64", "Target architecture")

	// Add verbose flag and connect it to log level
	buildCmd.Flags().Bool("verbose", false, "Enable verbose logging")
	buildCmd.PreRun = func(cmd *cobra.Command, args []string) {
		verbose, _ := cmd.Flags().GetBool("verbose")
		if verbose {
			log.SetLevel(log.DebugLevel)
		}
	}
}

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

	outFile := filepath.Join(bundleStagingDirPath, "config.cmd.yaml")
	err = os.WriteFile(outFile, configFileContent, 0644) // TODO do we need to preserve original permissions?
	if err != nil {
		return fmt.Errorf("error writing %s: %w", outFile, err)
	}

	// Create a visitor to handle the build-specific command processing in a recursive manner.
	cmdVisitor := BuildCommandVisitor{
		config:  cmdConfig,
		baseDir: bundleStagingDirPath,
	}

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

	// There's no need to setup an empty directory for the binary to be built in because we expect it to be
	// shared with other binaries already on the system.
	var binDirPath string
	binDirPath, err = executable.GetDestDir(cmdConfig)
	if err != nil {
		return err
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
	config  *types.CmdeagleConfig
	baseDir string
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

	if commandDef.Build == "" {
		return nil
	}

	log.Info("Running build script",
		"command", commandDef.Name,
		"script", commandDef.Build,
	)

	// TODO: Should probably allow the user to specify the shell to run the build script in.
	cmd := exec.Command("sh", "-c", commandDef.Build)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return err
	}

	if err := bundle.CopyIncludedFiles(v.config, commandDef, path, bundleStagingDirPath); err != nil {
		return err
	}

	return nil
}