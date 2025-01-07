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

type TemplateVariable = struct {
	Name        string
	Value       *string
	Placeholder string
}

var templateVariablePrefix = "{{"
var templateVariableSuffix = "}}"

var cliName string
var cliDescription string
var cliVersion string
var cliAuthor string
var cliLicense string
var cliLanguages []string

var templateVariables map[string]TemplateVariable = map[string]TemplateVariable{
	"name":        {Name: "name", Value: &cliName, Placeholder: "My CLI"},
	"description": {Name: "description", Value: &cliDescription, Placeholder: "My CLI is a tool to manage my projects."},
	"version":     {Name: "version", Value: &cliVersion, Placeholder: "0.0.1"},
	"author":      {Name: "author", Value: &cliAuthor, Placeholder: "John Doe"},
	"license":     {Name: "license", Value: &cliLicense, Placeholder: "MIT"},
}

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a cmdeagle project in a directory.",
	// Long:  ``, // TODO: Add long description
	// Args: cobra.PositionalArgs(cobra.ExactArgs(1)),

	RunE: func(cmd *cobra.Command, args []string) error {
		// Determine name of the project
		log.Info("Initializing cmdeagle project")

		// Check if cmd.yaml already exists
		if _, err := os.Stat(".cmd.yaml"); err == nil {
			log.Error(".cmd.yaml already exists in current directory")
			return fmt.Errorf(".cmd.yaml already exists in current directory")
		}

		if len(args) > 0 {
			cliName = args[len(args)-1]
		} else {
			// Get current directory name as fallback
			dir, err := os.Getwd()
			if err != nil {
				fmt.Printf("Error getting current directory: %v\n", err)
				return err
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
						Placeholder(templateVariables["name"].Placeholder).
						Value(templateVariables["name"].Value).
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
				huh.NewGroup(
					huh.NewInput().
						Title("Description (optional).").
						Description("This will be the description of the binary and the directory.").
						Placeholder(templateVariables["description"].Placeholder).
						Value(templateVariables["description"].Value),
					huh.NewInput().
						Title("Version (optional).").
						Placeholder(templateVariables["version"].Placeholder).
						Value(templateVariables["version"].Value),
					huh.NewInput().
						Title("Author (optional).").
						Placeholder(templateVariables["author"].Placeholder).
						Value(templateVariables["author"].Value),
					huh.NewInput().
						Title("License (optional).").
						Placeholder(templateVariables["license"].Placeholder).
						Value(templateVariables["license"].Value),
				).Description("Some optional metadata to document your CLI. Displayed by the help command."),
				huh.NewGroup(
					huh.NewMultiSelect[string]().
						Title("Which languages/runtimes do you want to support? (optional).").
						Description("We'll generate sample code showing how to integrate languages/runtimes you choose.").
						Options(
							huh.NewOption("go", "Go"),
							huh.NewOption("python", "Python"),
							huh.NewOption("rust", "Rust"),
							huh.NewOption("javascript", "JavaScript"),
							huh.NewOption("typescript", "TypeScript"),
						).
						Value(&cliLanguages),
				),
			)

			err = form.Run()

			if err != nil {
				log.Fatal(err)
			}
		}

		// Create the .cmd.yaml file
		fmt.Printf("Creating .cmd.yaml for %s\n", cliName)
		_, err := os.Create(".cmd.yaml")
		if err != nil {
			fmt.Printf("Error creating .cmd.yaml: %v\n", err)
			return err
		}

		// Replace template variables
		interpolatedYAML := string(sampleYAMLConfig)
		for name, variable := range templateVariables {
			interpolatedYAML = strings.ReplaceAll(interpolatedYAML, templateVariablePrefix+name+templateVariableSuffix, *variable.Value)
		}

		// Write the file
		err = os.WriteFile(".cmd.yaml", []byte(interpolatedYAML), 0644)
		if err != nil {
			log.Error("Failed to write .cmd.yaml", "error", err)
			return fmt.Errorf("failed to write .cmd.yaml: %w", err)
		}

		log.Info("Successfully initialized cmdeagle project", "name", cliName)
		return nil
	},
}
