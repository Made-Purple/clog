package gitutil

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/made-purple/clog/internal/fragment"
)

// setupRepo creates a fresh git repository in a temp dir, chdirs into it, and
// configures an identity so commits succeed. t.Chdir prevents t.Parallel.
func setupRepo(t *testing.T) {
	t.Helper()
	t.Chdir(t.TempDir())
	mustGit(t, "init", "-b", "master")
	mustGit(t, "config", "user.email", "test@example.com")
	mustGit(t, "config", "user.name", "Test User")
	mustGit(t, "config", "commit.gpgsign", "false")
}

func mustGit(t *testing.T, args ...string) string {
	t.Helper()
	out, err := exec.Command("git", args...).CombinedOutput()
	if err != nil {
		t.Fatalf("git %s: %v\n%s", strings.Join(args, " "), err, out)
	}
	return strings.TrimSpace(string(out))
}

func writeFile(t *testing.T, path, content string) {
	t.Helper()
	if dir := filepath.Dir(path); dir != "." {
		if err := os.MkdirAll(dir, 0755); err != nil {
			t.Fatalf("mkdir %s: %v", dir, err)
		}
	}
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("write %s: %v", path, err)
	}
}

func TestSanitizeBranchName(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"feature/knox-plugin", "feature-knox-plugin"},
		{"bugfix/fix-login-crash", "bugfix-fix-login-crash"},
		{"hotfix/3.81.1", "hotfix-3.81.1"},
		{"main", "main"},
		{"Feature/UPPER-Case", "feature-upper-case"},
		{"a/b/c/d", "a-b-c-d"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := SanitizeBranchName(tt.input)
			if got != tt.want {
				t.Errorf("SanitizeBranchName(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestBranchName(t *testing.T) {
	setupRepo(t)
	writeFile(t, "README.md", "hi\n")
	mustGit(t, "add", ".")
	mustGit(t, "commit", "-m", "init")
	mustGit(t, "checkout", "-b", "feature/foo")

	got, err := BranchName()
	if err != nil {
		t.Fatalf("BranchName: %v", err)
	}
	if got != "feature/foo" {
		t.Errorf("BranchName() = %q, want feature/foo", got)
	}
}

func TestBranchNameOutsideRepo(t *testing.T) {
	t.Chdir(t.TempDir())
	if _, err := BranchName(); err == nil {
		t.Errorf("BranchName() outside a repo: expected an error")
	}
}

func TestFileAtRef(t *testing.T) {
	setupRepo(t)
	writeFile(t, "CHANGELOG.md", "v1 content\n")
	mustGit(t, "add", ".")
	mustGit(t, "commit", "-m", "init")
	head := mustGit(t, "rev-parse", "HEAD")

	got, err := FileAtRef(head, "CHANGELOG.md")
	if err != nil {
		t.Fatalf("FileAtRef: %v", err)
	}
	if got != "v1 content\n" {
		t.Errorf("FileAtRef = %q, want %q", got, "v1 content\n")
	}

	// A file absent at the ref returns ("", nil), not an error.
	missing, err := FileAtRef(head, "does-not-exist.md")
	if err != nil {
		t.Fatalf("FileAtRef(missing): unexpected error %v", err)
	}
	if missing != "" {
		t.Errorf("FileAtRef(missing) = %q, want empty", missing)
	}
}

func TestMergeBase(t *testing.T) {
	setupRepo(t)
	writeFile(t, "a.txt", "1\n")
	mustGit(t, "add", ".")
	mustGit(t, "commit", "-m", "base")
	base := mustGit(t, "rev-parse", "HEAD")

	mustGit(t, "checkout", "-b", "feature")
	writeFile(t, "b.txt", "2\n")
	mustGit(t, "add", ".")
	mustGit(t, "commit", "-m", "feature work")

	got, err := MergeBase("master")
	if err != nil {
		t.Fatalf("MergeBase: %v", err)
	}
	if got != base {
		t.Errorf("MergeBase = %q, want %q", got, base)
	}
}

func TestWorkingTreeStatus(t *testing.T) {
	setupRepo(t)
	writeFile(t, "tracked.txt", "v1\n")
	mustGit(t, "add", ".")
	mustGit(t, "commit", "-m", "init")

	// A clean tree reports nothing.
	staged, unstaged, untracked, err := WorkingTreeStatus()
	if err != nil {
		t.Fatalf("WorkingTreeStatus (clean): %v", err)
	}
	if staged != "" || unstaged != "" || untracked != "" {
		t.Errorf("clean tree = (%q, %q, %q), want all empty", staged, unstaged, untracked)
	}

	writeFile(t, "added.txt", "new\n") // staged
	mustGit(t, "add", "added.txt")
	writeFile(t, "tracked.txt", "v2\n")   // unstaged modification
	writeFile(t, "untracked.txt", "x\n")  // untracked

	staged, unstaged, untracked, err = WorkingTreeStatus()
	if err != nil {
		t.Fatalf("WorkingTreeStatus (dirty): %v", err)
	}
	if !strings.Contains(staged, "added.txt") {
		t.Errorf("staged = %q, want to contain added.txt", staged)
	}
	if !strings.Contains(unstaged, "tracked.txt") {
		t.Errorf("unstaged = %q, want to contain tracked.txt", unstaged)
	}
	if !strings.Contains(untracked, "untracked.txt") {
		t.Errorf("untracked = %q, want to contain untracked.txt", untracked)
	}
}

func TestCommitRelease(t *testing.T) {
	setupRepo(t)
	writeFile(t, "CHANGELOG.md", "# Change Log\n")
	writeFile(t, filepath.Join("changelog.d", "feature-x.yaml"), "added:\n  - thing\n")
	writeFile(t, filepath.Join("changelog.d", fragment.SampleFilename), "added: []\n")
	writeFile(t, "VERSION", "0.9.0\n")
	mustGit(t, "add", ".")
	mustGit(t, "commit", "-m", "init")

	if err := CommitRelease("1.2.3", "changelog.d", "CHANGELOG.md", "VERSION"); err != nil {
		t.Fatalf("CommitRelease: %v", err)
	}

	// The real fragment is removed; the sample is preserved.
	if _, err := os.Stat(filepath.Join("changelog.d", "feature-x.yaml")); !os.IsNotExist(err) {
		t.Errorf("feature-x.yaml should have been removed")
	}
	if _, err := os.Stat(filepath.Join("changelog.d", fragment.SampleFilename)); err != nil {
		t.Errorf("sample fragment should be kept: %v", err)
	}
	if msg := mustGit(t, "log", "-1", "--pretty=%s"); msg != "Release v1.2.3" {
		t.Errorf("commit subject = %q, want %q", msg, "Release v1.2.3")
	}
}

func TestCommitReleaseUntrackedFragment(t *testing.T) {
	setupRepo(t)
	writeFile(t, "CHANGELOG.md", "# Change Log\n")
	mustGit(t, "add", ".")
	mustGit(t, "commit", "-m", "init")

	// A fragment that git never tracked: `git rm -f` fails, so CommitRelease
	// falls back to `git add`.
	writeFile(t, filepath.Join("changelog.d", "feature.yaml"), "added:\n  - thing\n")

	if err := CommitRelease("2.0.0", "changelog.d", "CHANGELOG.md"); err != nil {
		t.Fatalf("CommitRelease (untracked fragment): %v", err)
	}
	if msg := mustGit(t, "log", "-1", "--pretty=%s"); msg != "Release v2.0.0" {
		t.Errorf("commit subject = %q, want Release v2.0.0", msg)
	}
}

func TestCommitMigrate(t *testing.T) {
	setupRepo(t)
	writeFile(t, "CHANGELOG.md", "# Change Log\n")
	mustGit(t, "add", ".")
	mustGit(t, "commit", "-m", "init")

	writeFile(t, "CHANGELOG.md", "# Change Log\nmore\n")
	fragPath := filepath.Join("changelog.d", "feature.yaml")
	writeFile(t, fragPath, "added:\n  - x\n")

	if err := CommitMigrate("CHANGELOG.md", fragPath); err != nil {
		t.Fatalf("CommitMigrate: %v", err)
	}
	if msg := mustGit(t, "log", "-1", "--pretty=%s"); msg != "Migrated changelog entries to changelog fragments" {
		t.Errorf("commit subject = %q", msg)
	}
}

func TestTagRelease(t *testing.T) {
	setupRepo(t)
	writeFile(t, "a.txt", "1\n")
	mustGit(t, "add", ".")
	mustGit(t, "commit", "-m", "init")

	if err := TagRelease("9.9.9"); err != nil {
		t.Fatalf("TagRelease: %v", err)
	}
	if tags := mustGit(t, "tag", "--list"); !strings.Contains(tags, "v9.9.9") {
		t.Errorf("tags = %q, want to contain v9.9.9", tags)
	}
}

func TestRunError(t *testing.T) {
	setupRepo(t)
	if err := run("git", "not-a-real-subcommand"); err == nil {
		t.Errorf("run with an invalid git subcommand: expected an error")
	}
}

func TestMergeBaseError(t *testing.T) {
	setupRepo(t)
	writeFile(t, "a.txt", "1\n")
	mustGit(t, "add", ".")
	mustGit(t, "commit", "-m", "init")
	if _, err := MergeBase("no-such-ref"); err == nil {
		t.Errorf("MergeBase with an unknown ref: expected an error")
	}
}

func TestFileAtRefError(t *testing.T) {
	setupRepo(t)
	writeFile(t, "a.txt", "1\n")
	mustGit(t, "add", ".")
	mustGit(t, "commit", "-m", "init")
	// An invalid ref is a genuine error, distinct from a file missing at a valid ref.
	if _, err := FileAtRef("no-such-ref", "a.txt"); err == nil {
		t.Errorf("FileAtRef with an invalid ref: expected an error")
	}
}

func TestTagReleaseDuplicate(t *testing.T) {
	setupRepo(t)
	writeFile(t, "a.txt", "1\n")
	mustGit(t, "add", ".")
	mustGit(t, "commit", "-m", "init")
	if err := TagRelease("1.0.0"); err != nil {
		t.Fatalf("first TagRelease: %v", err)
	}
	if err := TagRelease("1.0.0"); err == nil {
		t.Errorf("duplicate TagRelease: expected an error")
	}
}
