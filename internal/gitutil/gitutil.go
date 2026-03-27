package gitutil

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
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

// CommitRelease stages the changelog and removed fragments, then commits.
func CommitRelease(version string, fragmentDir string, changelogPath string) error {
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
		if err := run("git", "rm", f); err != nil {
			// If not tracked by git, just stage the removal
			if err2 := run("git", "add", f); err2 != nil {
				return fmt.Errorf("removing fragment %s: %w", f, err)
			}
		}
	}

	msg := fmt.Sprintf("Release v%s", version)
	if err := run("git", "commit", "-m", msg); err != nil {
		return fmt.Errorf("committing release: %w", err)
	}

	return nil
}

// TagRelease creates a git tag for the release version.
func TagRelease(version string) error {
	tag := fmt.Sprintf("v%s", version)
	if err := run("git", "tag", tag); err != nil {
		return fmt.Errorf("creating tag %s: %w", tag, err)
	}
	return nil
}

func run(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdout = nil
	cmd.Stderr = nil
	return cmd.Run()
}
