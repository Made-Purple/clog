package command

import (
	"fmt"
	"os"

	"github.com/made-purple/clog/internal/fragment"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(validateCmd)
}

var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate changelog fragment files",
	RunE: func(cmd *cobra.Command, args []string) error {
		if _, err := os.Stat("changelog.d"); os.IsNotExist(err) {
			fmt.Println("No changelog.d/ directory found. Run `clog init` first.")
			return nil
		}

		fragments, readErrs := fragment.ReadAll("changelog.d")

		hasErrors := len(readErrs) > 0
		for _, err := range readErrs {
			fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		}

		if len(fragments) == 0 && !hasErrors {
			fmt.Println("No fragments found.")
			return nil
		}

		for _, f := range fragments {
			errs := fragment.Validate(f)
			for _, err := range errs {
				fmt.Fprintf(os.Stderr, "Error: %s\n", err)
				hasErrors = true
			}
		}

		if hasErrors {
			return fmt.Errorf("validation failed")
		}

		fmt.Printf("All %d fragment(s) are valid.\n", len(fragments))
		return nil
	},
}
