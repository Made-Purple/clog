package command

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/made-purple/clog/internal/fragment"
	"github.com/made-purple/clog/internal/gitutil"
	"github.com/spf13/cobra"
)

func init() {
	newCmd.Flags().Bool("edit", false, "Open the fragment file in $EDITOR after creation")
	rootCmd.AddCommand(newCmd)
}

var newCmd = &cobra.Command{
	Use:   "new",
	Short: "Create a new changelog fragment for the current branch",
	RunE: func(cmd *cobra.Command, args []string) error {
		branch, err := gitutil.BranchName()
		if err != nil {
			return err
		}

		filename := gitutil.SanitizeBranchName(branch) + ".yaml"
		path := filepath.Join("changelog.d", filename)

		edit, _ := cmd.Flags().GetBool("edit")

		if _, err := os.Stat(path); err == nil {
			fmt.Printf("Fragment already exists: %s\n", path)
			if edit {
				return openEditor(path)
			}
			return nil
		}

		if err := os.MkdirAll("changelog.d", 0755); err != nil {
			return fmt.Errorf("creating changelog.d: %w", err)
		}

		if err := os.WriteFile(path, fragment.Template(), 0644); err != nil {
			return fmt.Errorf("writing fragment: %w", err)
		}

		fmt.Printf("Created %s\n", path)

		if edit {
			return openEditor(path)
		}

		return nil
	},
}

func openEditor(path string) error {
	editor := os.Getenv("EDITOR")
	if editor == "" {
		return fmt.Errorf("EDITOR environment variable not set")
	}

	cmd := exec.Command(editor, path)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
