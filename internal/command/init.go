package command

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/made-purple/clog/internal/fragment"
	"github.com/spf13/cobra"
)

const defaultChangelog = `# Change Log
All notable changes to this project will be documented in this file.
This project adheres to [Semantic Versioning](http://semver.org/).

# Notes
[Deployment] Notes for deployment
[Added] for new features.
[Changed] for changes in existing functionality.
[Deprecated] for once-stable features removed in upcoming releases.
[Removed] for deprecated features removed in this release.
[Fixed] for any bug fixes.
[Security] to invite users to upgrade in case of vulnerabilities.
[YANKED] Note the emphasis, used for Hotfixes
`

func init() {
	rootCmd.AddCommand(initCmd)
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize changelog.d directory and CHANGELOG.md",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := os.MkdirAll("changelog.d", 0755); err != nil {
			return fmt.Errorf("creating changelog.d: %w", err)
		}
		fmt.Println("Created changelog.d/")

		samplePath := filepath.Join("changelog.d", fragment.SampleFilename)
		if err := os.WriteFile(samplePath, fragment.Template(), 0644); err != nil {
			return fmt.Errorf("creating %s: %w", samplePath, err)
		}
		fmt.Printf("Created %s\n", samplePath)

		if _, err := os.Stat("CHANGELOG.md"); os.IsNotExist(err) {
			if err := os.WriteFile("CHANGELOG.md", []byte(defaultChangelog), 0644); err != nil {
				return fmt.Errorf("creating CHANGELOG.md: %w", err)
			}
			fmt.Println("Created CHANGELOG.md")
		} else {
			fmt.Println("CHANGELOG.md already exists")
		}

		return nil
	},
}
