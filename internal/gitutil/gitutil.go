package gitutil

import (
	"bytes"
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/made-purple/clog/internal/fragment"
)

// BranchName returns the current git branch name.
func BranchName() (string, error) {
	out, err := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD").Output()
	if err != nil {
		return "", fmt.Errorf("could not determine branch name (are you in a git repository?): %w", err)
	}
	return strings.TrimSpace(string(out)), nil
}

// SanitizeBranchName converts a branch name to a safe fragment filename (without extension).
// Replaces "/" with "-" and lowercases.
func SanitizeBranchName(branch string) string {
	s := strings.ToLower(branch)
	s = strings.ReplaceAll(s, "/", "-")
	return s
}

// MergeBase returns the merge-base commit between HEAD and the given base ref.
func MergeBase(base string) (string, error) {
	out, err := exec.Command("git", "merge-base", "HEAD", base).Output()
	if err != nil {
		return "", fmt.Errorf("finding merge-base with %s (is it fetched?): %w", base, err)
	}
	return strings.TrimSpace(string(out)), nil
}

// FileAtRef returns the contents of a file at the given git ref.
// Returns ("", nil) if the file did not exist at that ref.
func FileAtRef(ref, path string) (string, error) {
	cmd := exec.Command("git", "show", ref+":"+path)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		msg := stderr.String()
		if strings.Contains(msg, "does not exist") || strings.Contains(msg, "exists on disk, but not in") {
			return "", nil
		}
		return "", fmt.Errorf("reading %s at %s: %w: %s", path, ref, err, strings.TrimSpace(msg))
	}
	return stdout.String(), nil
}

// WorkingTreeStatus returns a summary of uncommitted changes in the working tree.
// Returns empty strings if the tree is clean.
func WorkingTreeStatus() (staged, unstaged, untracked string, err error) {
	out, err := exec.Command("git", "status", "--porcelain").Output()
	if err != nil {
		return "", "", "", fmt.Errorf("checking git status: %w", err)
	}

	var stagedLines, unstagedLines, untrackedLines []string
	for _, line := range strings.Split(strings.TrimSpace(string(out)), "\n") {
		if line == "" {
			continue
		}
		// Porcelain format: XY filename
		// X = index status, Y = working tree status
		x := line[0]
		y := line[1]
		file := strings.TrimSpace(line[3:])

		if x == '?' {
			untrackedLines = append(untrackedLines, file)
		} else {
			if x != ' ' {
				stagedLines = append(stagedLines, file)
			}
			if y != ' ' {
				unstagedLines = append(unstagedLines, file)
			}
		}
	}

	format := func(lines []string) string {
		if len(lines) == 0 {
			return ""
		}
		var b strings.Builder
		for _, l := range lines {
			b.WriteString("  " + l + "\n")
		}
		return b.String()
	}

	return format(stagedLines), format(unstagedLines), format(untrackedLines), nil
}

// CommitRelease stages the changelog, removed fragments, and any extra files, then commits.
func CommitRelease(version string, fragmentDir string, changelogPath string, extraFiles ...string) error {
	// Stage the updated changelog
	if err := run("git", "add", changelogPath); err != nil {
		return fmt.Errorf("staging changelog: %w", err)
	}

	// Find and remove fragment files
	yamls, err := filepath.Glob(filepath.Join(fragmentDir, "*.yaml"))
	if err != nil {
		return fmt.Errorf("globbing fragments: %w", err)
	}
	for _, f := range yamls {
		if filepath.Base(f) == fragment.SampleFilename {
			continue
		}
		if err := run("git", "rm", "-f", f); err != nil {
			// If not tracked by git, just stage the removal
			if err2 := run("git", "add", f); err2 != nil {
				return fmt.Errorf("removing fragment %s: %w", f, err)
			}
		}
	}

	// Stage any extra files (e.g., VERSION, package.json)
	for _, f := range extraFiles {
		if err := run("git", "add", f); err != nil {
			return fmt.Errorf("staging %s: %w", f, err)
		}
	}

	msg := fmt.Sprintf("Release v%s", version)
	if err := run("git", "commit", "-m", msg); err != nil {
		return fmt.Errorf("committing release: %w", err)
	}

	return nil
}

// CommitMigrate stages the changelog and fragment file, then commits with a migration message.
func CommitMigrate(changelogPath string, fragmentPath string) error {
	if err := run("git", "add", changelogPath); err != nil {
		return fmt.Errorf("staging changelog: %w", err)
	}
	if err := run("git", "add", fragmentPath); err != nil {
		return fmt.Errorf("staging fragment: %w", err)
	}
	if err := run("git", "commit", "-m", "Migrated changelog entries to changelog fragments"); err != nil {
		return fmt.Errorf("committing migration: %w", err)
	}
	return nil
}

// TagRelease creates a git tag for the release version.
func TagRelease(version string) error {
	tag := fmt.Sprintf("v%s", version)
	if err := run("git", "tag", "-a", tag, "-m", fmt.Sprintf("Release %s", tag)); err != nil {
		return fmt.Errorf("creating tag %s: %w", tag, err)
	}
	return nil
}

func run(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		if stderr.Len() > 0 {
			return fmt.Errorf("%w: %s", err, strings.TrimSpace(stderr.String()))
		}
		return err
	}
	return nil
}
