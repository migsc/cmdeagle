/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "cmdeagle",
	Version: "0.11.17",
	Short:   "Build powerful CLI tools with simple YAML configuration",
	Long: `Create professional command-line tools without the complexity.

cmdeagle helps you build CLI applications that:
- Are easy to configure using simple YAML files - no complex boilerplate needed
- Work with your favorite scripting language - currently optimized for JavaScript, Python, 
  and shell scripts (support for compiled languages coming soon)
- Package your code and assets into a single executable (note: language runtimes like Node.js 
  or Python must be installed separately - consider Docker/Podman for full environment portability)
- Run consistently across platforms when properly containerized
- Scale from simple scripts to complex tools

Get started with 'cmdeagle init myapp' to create your first CLI application.`,
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

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.cmdeagle.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
