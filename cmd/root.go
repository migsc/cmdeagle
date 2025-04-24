/*
Copyright © 2024 Miguel Chateloin
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "cmdeagle",
	Version: "dev",
	Short:   "Build powerful CLI tools with simple YAML configuration",
	Long: `Create professional command-line tools without the complexity.

cmdeagle helps you build CLI applications that:
- Are easy to configure using simple YAML files - no complex boilerplate needed
- Work with your favorite scripting language - optimized for JavaScript, 
  Python, and shell scripts (compiled languages coming soon)
- Package your code and assets into a single executable 
  (note: runtimes like Node.js/Python need separate installation)
- Run consistently across platforms when properly containerized
- Scale from simple scripts to complex tools

Get started with 'cmdeagle init myapp' to create your first CLI application.`,
	Example: `  # Initialize a new CLI project
  cmdeagle init mycli

  # Build your CLI
  cmdeagle build

  # Build for a specific platform
  cmdeagle build --os linux --arch arm64`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

// SetVersion allows setting the version at runtime, typically from main
func SetVersion(v string) {
	rootCmd.Version = v
}

func init() {
	// Add version flag
	rootCmd.SetVersionTemplate("cmdeagle version {{.Version}}\n")

	// Add app metadata to help template
	cobra.AddTemplateFunc("AppMetadata", func() string {
		return fmt.Sprintf("Author: %s\nLicense: %s\n", "Miguel Chateloin <miguel@chateloin.com>", "MIT")
	})

	rootCmd.SetHelpTemplate(`{{with .Long}}{{. | trimTrailingWhitespaces}}{{end}}

Usage:{{if .Runnable}}
  {{.UseLine}}{{end}}{{if .HasAvailableSubCommands}}
  {{.CommandPath}} [command]{{end}}{{if gt (len .Aliases) 0}}

Aliases:
  {{.NameAndAliases}}{{end}}{{if .HasExample}}

Examples:
{{.Example}}{{end}}{{if .HasAvailableSubCommands}}

Available Commands:{{range .Commands}}{{if (or .IsAvailableCommand (eq .Name "help"))}}
  {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableLocalFlags}}

Flags:
{{.LocalFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasAvailableInheritedFlags}}

Global Flags:
{{.InheritedFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasHelpSubCommands}}

Additional help topics:{{range .Commands}}{{if .IsAdditionalHelpTopicCommand}}
  {{rpad .CommandPath .CommandPathPadding}} {{.Short}}{{end}}{{end}}{{end}}

{{AppMetadata}}
Use "{{.CommandPath}} [command] --help" for more information about a command.
`)
}
