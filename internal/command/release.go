package command

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/made-purple/clog/internal/changelog"
	"github.com/made-purple/clog/internal/color"
	"github.com/made-purple/clog/internal/fragment"
	"github.com/made-purple/clog/internal/gitutil"
	"github.com/made-purple/clog/internal/merge"
	"github.com/spf13/cobra"
)

const fragmentDir = "changelog.d"
const changelogPath = "CHANGELOG.md"

func init() {
	rootCmd.AddCommand(releaseCmd)
}

var releaseCmd = &cobra.Command{
	Use:   "release",
	Short: "Merge changelog fragments into CHANGELOG.md and create a release entry",
	RunE:  runRelease,
}

func runRelease(cmd *cobra.Command, args []string) error {
	reader := bufio.NewReader(os.Stdin)

	// 0. Check for dirty working tree
	staged, unstaged, untracked, err := gitutil.WorkingTreeStatus()
	if err == nil && (staged != "" || unstaged != "" || untracked != "") {
		color.Warn("Your working tree has uncommitted changes:")
		if staged != "" {
			fmt.Printf("\n%s\n%s", color.Yellow("Staged:"), staged)
		}
		if unstaged != "" {
			fmt.Printf("\n%s\n%s", color.Yellow("Unstaged:"), unstaged)
		}
		if untracked != "" {
			fmt.Printf("\n%s\n%s", color.Yellow("Untracked:"), untracked)
		}
		color.Prompt("\nContinue with release anyway? (y/n):")
		answer, err := reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("reading answer: %w", err)
		}
		if strings.ToLower(strings.TrimSpace(answer)) != "y" {
			fmt.Println("Release cancelled. Clean your working tree and try again.")
			return nil
		}
	}

	// 1. Read and parse CHANGELOG.md
	cl, err := changelog.Read(changelogPath)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("CHANGELOG.md not found. Run `clog init` first")
		}
		return err
	}

	lastVersion := changelog.LastVersion(cl)

	// 2. Read all fragments
	if _, err := os.Stat(fragmentDir); os.IsNotExist(err) {
		return fmt.Errorf("no %s/ directory found. Run `clog init` first", fragmentDir)
	}

	fragments, readErrs := fragment.ReadAll(fragmentDir)
	for _, err := range readErrs {
		fmt.Fprintf(os.Stderr, "%s %s\n", color.Red("Error:"), err)
	}
	if len(readErrs) > 0 {
		return fmt.Errorf("failed to read some fragment files, aborting release")
	}

	if len(fragments) == 0 {
		fmt.Println("No fragment files found. Nothing to release.")
		return nil
	}

	// Validate all fragments
	hasValidationErrors := false
	for _, f := range fragments {
		errs := fragment.Validate(f)
		for _, err := range errs {
			fmt.Fprintf(os.Stderr, "%s %s\n", color.Red("Error:"), err)
			hasValidationErrors = true
		}
	}
	if hasValidationErrors {
		return fmt.Errorf("fragment validation failed, aborting release")
	}

	// 3. Merge fragments
	merged := merge.Merge(fragments)
	if len(merged) == 0 {
		fmt.Println("No entries found in fragments (all empty). Nothing to release.")
		return nil
	}

	// 4. Prompt for version
	if lastVersion != "" {
		fmt.Printf("%s %s\n", color.Dim("Last version:"), lastVersion)
	} else {
		fmt.Println(color.Dim("No previous version found."))
	}
	color.Prompt("Enter new version:")
	version, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("reading version: %w", err)
	}
	version = strings.TrimSpace(version)
	if version == "" {
		return fmt.Errorf("version cannot be empty")
	}

	// 5. Prompt for optional metadata
	color.Prompt("Enter optional metadata (e.g. (98%)(Dev)), or press Enter to skip:")
	metadata, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("reading metadata: %w", err)
	}
	metadata = strings.TrimSpace(metadata)

	// 6. Render the entry
	date := time.Now().Format("2006-01-02")
	entry := merge.Render(version, date, metadata, merged)

	// 7. Preview and confirm
	fmt.Println()
	fmt.Println(color.Dim("─── Preview ───"))
	fmt.Print(entry)
	fmt.Println(color.Dim("─── End Preview ───"))
	color.Prompt("\nProceed with release? (y/n):")
	confirm, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("reading confirmation: %w", err)
	}
	if strings.ToLower(strings.TrimSpace(confirm)) != "y" {
		fmt.Println("Release cancelled.")
		return nil
	}

	// 8. Insert into CHANGELOG.md
	result := changelog.Insert(cl, entry)
	if err := os.WriteFile(changelogPath, []byte(result), 0644); err != nil {
		return fmt.Errorf("writing CHANGELOG.md: %w", err)
	}
	color.Success("Updated CHANGELOG.md")

	// 9. Delete fragment files
	dirEntries, err := os.ReadDir(fragmentDir)
	if err != nil {
		return fmt.Errorf("reading %s: %w", fragmentDir, err)
	}
	for _, e := range dirEntries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), ".yaml") && e.Name() != fragment.SampleFilename {
			path := filepath.Join(fragmentDir, e.Name())
			if err := os.Remove(path); err != nil {
				fmt.Fprintf(os.Stderr, "%s could not remove %s: %s\n", color.Yellow("!"), path, err)
			}
		}
	}
	color.Success("Removed changelog fragments")

	// 10. Ask about auto-commit
	color.Prompt("Auto-commit release? (y/n):")
	commitAnswer, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("reading commit answer: %w", err)
	}
	if strings.ToLower(strings.TrimSpace(commitAnswer)) == "y" {
		if err := gitutil.CommitRelease(version, fragmentDir, changelogPath); err != nil {
			return fmt.Errorf("commit failed: %w", err)
		}
		color.Success("Committed: Release v%s", version)

		// 11. Ask about tagging
		color.Prompt("Tag release? (y/n):")
		tagAnswer, err := reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("reading tag answer: %w", err)
		}
		if strings.ToLower(strings.TrimSpace(tagAnswer)) == "y" {
			if err := gitutil.TagRelease(version); err != nil {
				return fmt.Errorf("tagging failed: %w", err)
			}
			color.Success("Tagged: v%s", version)
		}
	}

	fmt.Println()
	fmt.Println(color.BoldGreen("Release complete!"))
	return nil
}
