package command

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "clog",
	Short: "Changelog fragment manager",
	Long:  "clog manages changelog fragments in YAML format and merges them into CHANGELOG.md at release time.",
}

func Execute() error {
	return rootCmd.Execute()
}
