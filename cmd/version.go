package cmd

import (
	_ "embed"
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
)

//go:embed package.json
var packageJSON []byte

type PackageJSON struct {
	Version string `json:"version"`
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of cmdeagle",
	Run: func(cmd *cobra.Command, args []string) {
		var pkg PackageJSON
		if err := json.Unmarshal(packageJSON, &pkg); err != nil {
			fmt.Println("Error parsing package.json:", err)
			return
		}

		fmt.Println(pkg.Version)
	},
}
