/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(initCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// initCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// initCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	// rootCmd.Args
}

//go:embed template.cmd.yaml
var sampleYAMLConfig []byte

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a cmdeagle project in a directory.",
	// Long:  ``, // TODO: Add long description
	// Args: cobra.PositionalArgs(cobra.ExactArgs(1)),

	Run: func(cmd *cobra.Command, args []string) {
		// Determine name of the project
		log.Info("Initializing cmdeagle project")

		var cliName string
		if len(args) > 0 {
			cliName = args[0]
		} else {
			// Get current directory name as fallback
			dir, err := os.Getwd()
			if err != nil {
				fmt.Printf("Error getting current directory: %v\n", err)
				return
			}

			var folderName = filepath.Base(dir)

			cliName = folderName

			form := huh.NewForm(
				// Gather some final details about the order.
				huh.NewGroup(
					huh.NewInput().
						Title("What’s your CLI's name?.").
						Description("This will be the name of the binary and the directory.").
						// TODO implement rename command then change to: Description("This will be the name of the binary and the directory. You can change this later with `cmdeagle rename` command. It must be lowercase.").
						Placeholder("mycli").
						Value(&cliName).
						// Validating fields is easy. The form will mark erroneous fields
						// and display error messages accordingly.
						Validate(func(str string) error {
							if str == "" {
								return fmt.Errorf("cli name cannot be empty")
							}

							if !regexp.MustCompile(`^[a-z0-9\_\-]+$`).MatchString(str) {
								return fmt.Errorf("cli name must be only lowercase letters(a-z), numbers(0-9), hypens( - ), and underscores( _ ) are allowed")
							}

							return nil
						}),
				),
			)

			err = form.Run()

			if err != nil {
				log.Fatal(err)
			}
		}

		// Create the cmd.yaml file
		fmt.Printf("Creating cmd.yaml for %s\n", cliName)
		file, err := os.Create("cmd.yaml")
		if err != nil {
			fmt.Printf("Error creating cmd.yaml: %v\n", err)
			return
		}

		// Interpolate the name into the sample YAML
		var interpolatedYAML string
		interpolatedYAML = strings.Replace(string(sampleYAMLConfig), "{{name}}", cliName, -1)
		interpolatedYAML = strings.Replace(interpolatedYAML, "{{license}}", "Your License", -1)

		// Write the sample YAML to the file
		countBytes, err := file.Write([]byte(interpolatedYAML))
		if err != nil {
			fmt.Printf("Error writing to cmd.yaml: %v\n", err)
			return
		}

		fmt.Printf("Wrote %d bytes to cmd.yaml\n", countBytes)
		defer file.Close()

	},
}
