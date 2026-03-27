package command

import (
	"github.com/spf13/cobra"
)

// Version is set at build time via ldflags.
var Version = "dev"

var rootCmd = &cobra.Command{
	Use:     "clog",
	Short:   "Changelog fragment manager",
	Long:    "clog manages changelog fragments in YAML format and merges them into CHANGELOG.md at release time.",
	Version: Version,
}

func Execute() error {
	return rootCmd.Execute()
}
