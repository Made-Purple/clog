package command

import (
	"fmt"
	"os"
	"time"

	"github.com/made-purple/clog/internal/fragment"
	"github.com/made-purple/clog/internal/merge"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(previewCmd)
}

var previewCmd = &cobra.Command{
	Use:   "preview",
	Short: "Preview what the next release entry would look like",
	RunE: func(cmd *cobra.Command, args []string) error {
		if _, err := os.Stat("changelog.d"); os.IsNotExist(err) {
			fmt.Println("No changelog.d/ directory found. Nothing to preview.")
			return nil
		}

		fragments, readErrs := fragment.ReadAll("changelog.d")
		for _, err := range readErrs {
			fmt.Fprintf(os.Stderr, "Warning: %s\n", err)
		}

		if len(fragments) == 0 {
			fmt.Println("No fragments found. Nothing to preview.")
			return nil
		}

		merged := merge.Merge(fragments)
		if len(merged) == 0 {
			fmt.Println("No entries to preview (all fragments are empty).")
			return nil
		}

		date := time.Now().Format("2006-01-02")
		output := merge.Render("X.Y.Z", date, "", merged)
		fmt.Print(output)

		return nil
	},
}
