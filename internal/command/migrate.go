package command

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/made-purple/clog/internal/changelog"
	"github.com/made-purple/clog/internal/color"
	"github.com/made-purple/clog/internal/fragment"
	"github.com/made-purple/clog/internal/gitutil"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(migrateCmd)
}

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Migrate entries from the [staging] section of CHANGELOG.md into a fragment file",
	RunE:  runMigrate,
}

func runMigrate(cmd *cobra.Command, args []string) error {
	reader := bufio.NewReader(os.Stdin)

	// 1. Read CHANGELOG.md
	cl, err := changelog.Read(changelogPath)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("CHANGELOG.md not found. Run `clog init` first")
		}
		return err
	}

	// 2. Extract staging entries
	staging, err := changelog.ExtractStaging(cl)
	if err != nil {
		return err
	}
	if staging == nil {
		fmt.Println("No [staging] section found in CHANGELOG.md.")
		return nil
	}

	// 3. Determine target fragment file
	branch, err := gitutil.BranchName()
	if err != nil {
		return err
	}
	filename := gitutil.SanitizeBranchName(branch) + ".yaml"
	path := filepath.Join(fragmentDir, filename)

	// 4. Show preview
	fmt.Println()
	fmt.Println(color.Dim("─── Entries to migrate ───"))
	for _, cat := range fragment.CategoryOrder {
		entries, ok := staging[cat]
		if !ok || len(entries) == 0 {
			continue
		}
		fmt.Printf("### %s\n", fragment.CategoryDisplay[cat])
		for _, e := range entries {
			fmt.Printf("  - %s\n", e)
		}
	}
	fmt.Println(color.Dim("──────────────────────────"))
	fmt.Printf("\nTarget: %s\n", color.Cyan(path))

	// 5. Check if fragment file already exists
	existingEntries := make(map[string][]string)
	if _, err := os.Stat(path); err == nil {
		color.Warn("Fragment file already exists: %s", path)
		fmt.Println("  Existing entries will be merged with migrated entries.")

		data, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("reading existing fragment: %w", err)
		}
		existing, err := fragment.Parse(filename, data)
		if err != nil {
			return fmt.Errorf("parsing existing fragment: %w", err)
		}
		existingEntries = fragment.NonEmptyEntries(existing)
	}

	// 6. Ask for confirmation
	fmt.Println()
	color.Prompt("Migrate these entries? (y/n):")
	confirm, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("reading confirmation: %w", err)
	}
	if strings.ToLower(strings.TrimSpace(confirm)) != "y" {
		fmt.Println("Migration cancelled.")
		return nil
	}

	// 7. Merge entries
	merged := make(map[string][]string)
	for cat, entries := range existingEntries {
		merged[cat] = append(merged[cat], entries...)
	}
	for cat, entries := range staging {
		existing := make(map[string]bool)
		for _, e := range merged[cat] {
			existing[e] = true
		}
		for _, e := range entries {
			if !existing[e] {
				merged[cat] = append(merged[cat], e)
			}
		}
	}

	// 8. Write fragment file
	if err := os.MkdirAll(fragmentDir, 0755); err != nil {
		return fmt.Errorf("creating %s: %w", fragmentDir, err)
	}
	if err := os.WriteFile(path, fragment.MarshalEntries(merged), 0644); err != nil {
		return fmt.Errorf("writing fragment: %w", err)
	}
	color.Success("Written %s", path)

	// 9. Remove staging section from CHANGELOG.md
	updated := changelog.RemoveStaging(cl)
	if err := os.WriteFile(changelogPath, []byte(updated), 0644); err != nil {
		return fmt.Errorf("writing CHANGELOG.md: %w", err)
	}
	color.Success("Removed [staging] section from CHANGELOG.md")

	// 10. Ask about auto-commit
	fmt.Println()
	color.Prompt("Auto-commit changes? (y/n):")
	commitAnswer, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("reading commit answer: %w", err)
	}
	if strings.ToLower(strings.TrimSpace(commitAnswer)) == "y" {
		if err := gitutil.CommitMigrate(changelogPath, path); err != nil {
			return fmt.Errorf("commit failed: %w", err)
		}
		color.Success("Committed: Migrated changelog entries to changelog fragments")
	}

	return nil
}
