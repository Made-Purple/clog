package command

import (
	"bufio"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/made-purple/clog/internal/fragment"
	"github.com/made-purple/clog/internal/skill"
)

// mute redirects stdout and stderr to /dev/null for the duration of a test so
// command output doesn't clutter the test log.
func mute(t *testing.T) {
	t.Helper()
	devnull, err := os.OpenFile(os.DevNull, os.O_WRONLY, 0)
	if err != nil {
		t.Fatalf("open %s: %v", os.DevNull, err)
	}
	outOrig, errOrig := os.Stdout, os.Stderr
	os.Stdout, os.Stderr = devnull, devnull
	t.Cleanup(func() {
		os.Stdout, os.Stderr = outOrig, errOrig
		devnull.Close()
	})
}

// withStdin feeds the given input to os.Stdin for the duration of a test, so
// interactive prompts read scripted answers.
func withStdin(t *testing.T, input string) {
	t.Helper()
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("pipe: %v", err)
	}
	orig := os.Stdin
	os.Stdin = r
	t.Cleanup(func() { os.Stdin = orig; r.Close() })
	go func() {
		io.WriteString(w, input)
		w.Close()
	}()
}

func runGit(t *testing.T, args ...string) string {
	t.Helper()
	out, err := exec.Command("git", args...).CombinedOutput()
	if err != nil {
		t.Fatalf("git %s: %v\n%s", strings.Join(args, " "), err, out)
	}
	return strings.TrimSpace(string(out))
}

func initRepo(t *testing.T) {
	t.Helper()
	runGit(t, "init", "-b", "master")
	runGit(t, "config", "user.email", "test@example.com")
	runGit(t, "config", "user.name", "Test User")
	runGit(t, "config", "commit.gpgsign", "false")
}

func TestScopeLabel(t *testing.T) {
	if got := scopeLabel(skill.Global); got != "global" {
		t.Errorf("scopeLabel(Global) = %q, want global", got)
	}
	if got := scopeLabel(skill.Project); got != "project" {
		t.Errorf("scopeLabel(Project) = %q, want project", got)
	}
}

func TestDedupeAgents(t *testing.T) {
	in := []skill.Agent{skill.Claude, skill.Claude, skill.Codex, skill.Codex, skill.Claude}
	out := dedupeAgents(in)
	if len(out) != 2 || out[0].Key != "claude" || out[1].Key != "codex" {
		t.Errorf("dedupeAgents = %+v, want [claude codex]", out)
	}
}

func TestDiffStaging(t *testing.T) {
	current := map[string][]string{
		"added": {"a", "b", "c"},
		"fixed": {"f1"},
	}
	base := map[string][]string{
		"added": {"a"},
	}
	got := diffStaging(current, base)
	if !reflect.DeepEqual(got["added"], []string{"b", "c"}) {
		t.Errorf("added diff = %v, want [b c]", got["added"])
	}
	if !reflect.DeepEqual(got["fixed"], []string{"f1"}) {
		t.Errorf("fixed diff = %v, want [f1]", got["fixed"])
	}
}

func TestPromptYesNo(t *testing.T) {
	cases := []struct {
		in   string
		want bool
	}{
		{"\n", true}, // default is yes
		{"y\n", true},
		{"Y\n", true},
		{"yes\n", true},
		{"YES\n", true},
		{"n\n", false},
		{"no\n", false},
		{"maybe\n", false},
	}
	for _, c := range cases {
		got, err := promptYesNo(bufio.NewReader(strings.NewReader(c.in)), "Proceed?")
		if err != nil {
			t.Fatalf("promptYesNo(%q): %v", c.in, err)
		}
		if got != c.want {
			t.Errorf("promptYesNo(%q) = %v, want %v", c.in, got, c.want)
		}
	}
}

func TestPromptYesNoReadError(t *testing.T) {
	if _, err := promptYesNo(bufio.NewReader(strings.NewReader("")), "Proceed?"); err == nil {
		t.Errorf("promptYesNo with empty input: expected an error")
	}
}

func TestInitCommand(t *testing.T) {
	t.Chdir(t.TempDir())
	mute(t)

	if err := initCmd.RunE(initCmd, nil); err != nil {
		t.Fatalf("init: %v", err)
	}
	for _, p := range []string{
		"changelog.d",
		filepath.Join("changelog.d", fragment.SampleFilename),
		"CHANGELOG.md",
	} {
		if _, err := os.Stat(p); err != nil {
			t.Errorf("expected %s to exist: %v", p, err)
		}
	}

	// A second run hits the "CHANGELOG.md already exists" branch.
	if err := initCmd.RunE(initCmd, nil); err != nil {
		t.Fatalf("init (rerun): %v", err)
	}
}

func TestValidateCommand(t *testing.T) {
	t.Chdir(t.TempDir())
	mute(t)

	// No changelog.d directory: handled gracefully, no error.
	if err := validateCmd.RunE(validateCmd, nil); err != nil {
		t.Fatalf("validate (no dir): %v", err)
	}

	if err := os.MkdirAll("changelog.d", 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join("changelog.d", "feature.yaml"), fragment.Template(), 0644); err != nil {
		t.Fatal(err)
	}
	if err := validateCmd.RunE(validateCmd, nil); err != nil {
		t.Fatalf("validate (valid fragment): %v", err)
	}

	// An unknown category fails validation.
	if err := os.WriteFile(filepath.Join("changelog.d", "bad.yaml"), []byte("bogus:\n  - item\n"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := validateCmd.RunE(validateCmd, nil); err == nil {
		t.Errorf("validate (invalid fragment): expected an error")
	}
}

func TestPreviewCommand(t *testing.T) {
	t.Chdir(t.TempDir())
	mute(t)

	// No changelog.d directory.
	if err := previewCmd.RunE(previewCmd, nil); err != nil {
		t.Fatalf("preview (no dir): %v", err)
	}

	if err := os.MkdirAll("changelog.d", 0755); err != nil {
		t.Fatal(err)
	}

	// Only the sample is present -> ReadAll skips it -> no fragments.
	if err := os.WriteFile(filepath.Join("changelog.d", fragment.SampleFilename), fragment.Template(), 0644); err != nil {
		t.Fatal(err)
	}
	if err := previewCmd.RunE(previewCmd, nil); err != nil {
		t.Fatalf("preview (sample only): %v", err)
	}

	// An all-empty fragment yields no entries to preview.
	if err := os.WriteFile(filepath.Join("changelog.d", "empty.yaml"), fragment.Template(), 0644); err != nil {
		t.Fatal(err)
	}
	if err := previewCmd.RunE(previewCmd, nil); err != nil {
		t.Fatalf("preview (empty fragment): %v", err)
	}

	// A fragment with a real entry exercises the render path.
	if err := os.WriteFile(filepath.Join("changelog.d", "feature.yaml"), []byte("added:\n  - \"Did a thing\"\n"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := previewCmd.RunE(previewCmd, nil); err != nil {
		t.Fatalf("preview (with entry): %v", err)
	}
}

func TestNewCommand(t *testing.T) {
	t.Chdir(t.TempDir())
	mute(t)
	initRepo(t)
	runGit(t, "checkout", "-b", "feature/foo")
	runGit(t, "commit", "--allow-empty", "-m", "init")

	if err := newCmd.RunE(newCmd, nil); err != nil {
		t.Fatalf("new: %v", err)
	}
	want := filepath.Join("changelog.d", "feature-foo.yaml")
	if _, err := os.Stat(want); err != nil {
		t.Errorf("expected %s to exist: %v", want, err)
	}

	// A second run hits the "already exists" branch.
	if err := newCmd.RunE(newCmd, nil); err != nil {
		t.Fatalf("new (rerun): %v", err)
	}
}

func TestReleaseHappyPath(t *testing.T) {
	t.Chdir(t.TempDir())
	mute(t)
	initRepo(t)

	if err := os.WriteFile(changelogPath, []byte(defaultChangelog), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(fragmentDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(fragmentDir, "feature.yaml"), []byte("added:\n  - \"New feature\"\n"), 0644); err != nil {
		t.Fatal(err)
	}
	runGit(t, "add", ".")
	runGit(t, "commit", "-m", "init")

	// version, metadata (skip), proceed, auto-commit, tag.
	withStdin(t, "1.0.0\n\ny\ny\ny\n")

	if err := runRelease(releaseCmd, nil); err != nil {
		t.Fatalf("release: %v", err)
	}

	out, err := os.ReadFile(changelogPath)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(out), "1.0.0") {
		t.Errorf("CHANGELOG.md missing version:\n%s", out)
	}
	if !strings.Contains(string(out), "New feature") {
		t.Errorf("CHANGELOG.md missing entry:\n%s", out)
	}
	if _, err := os.Stat(filepath.Join(fragmentDir, "feature.yaml")); !os.IsNotExist(err) {
		t.Errorf("fragment should have been removed")
	}
	if msg := runGit(t, "log", "-1", "--pretty=%s"); msg != "Release v1.0.0" {
		t.Errorf("commit subject = %q, want Release v1.0.0", msg)
	}
	if tags := runGit(t, "tag", "--list"); !strings.Contains(tags, "v1.0.0") {
		t.Errorf("tags = %q, want to contain v1.0.0", tags)
	}
}

func TestReleaseCancelled(t *testing.T) {
	t.Chdir(t.TempDir())
	mute(t)
	initRepo(t)

	if err := os.WriteFile(changelogPath, []byte(defaultChangelog), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(fragmentDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(fragmentDir, "feature.yaml"), []byte("added:\n  - \"New feature\"\n"), 0644); err != nil {
		t.Fatal(err)
	}
	runGit(t, "add", ".")
	runGit(t, "commit", "-m", "init")

	// version, metadata, then decline at the proceed prompt.
	withStdin(t, "1.0.0\n\nn\n")

	if err := runRelease(releaseCmd, nil); err != nil {
		t.Fatalf("release (cancelled): %v", err)
	}

	// Nothing was changed: the fragment is still there and the changelog has no entry.
	if _, err := os.Stat(filepath.Join(fragmentDir, "feature.yaml")); err != nil {
		t.Errorf("fragment should remain after cancel: %v", err)
	}
	out, _ := os.ReadFile(changelogPath)
	if strings.Contains(string(out), "New feature") {
		t.Errorf("CHANGELOG.md should be unchanged after cancel")
	}
}

func TestSkillInstallAndUninstallViaFlags(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	mute(t)

	mustSetFlag(t, skillInstallCmd.Flags().Set("claude", "true"))
	mustSetFlag(t, skillInstallCmd.Flags().Set("global", "true"))
	t.Cleanup(func() {
		skillInstallCmd.Flags().Set("claude", "false")
		skillInstallCmd.Flags().Set("global", "false")
	})

	if err := runSkillInstall(skillInstallCmd, nil); err != nil {
		t.Fatalf("skill install: %v", err)
	}
	if !skill.Claude.Installed(skill.Global) {
		t.Errorf("Claude skill not installed after install --claude --global")
	}
	if !installedAnywhere(skill.Claude) {
		t.Errorf("installedAnywhere(Claude) = false after install")
	}

	// A second install is a no-op (already up to date) and must still succeed.
	if err := runSkillInstall(skillInstallCmd, nil); err != nil {
		t.Fatalf("skill install (rerun): %v", err)
	}

	mustSetFlag(t, skillUninstallCmd.Flags().Set("claude", "true"))
	mustSetFlag(t, skillUninstallCmd.Flags().Set("global", "true"))
	t.Cleanup(func() {
		skillUninstallCmd.Flags().Set("claude", "false")
		skillUninstallCmd.Flags().Set("global", "false")
	})

	if err := runSkillUninstall(skillUninstallCmd, nil); err != nil {
		t.Fatalf("skill uninstall: %v", err)
	}
	if skill.Claude.Installed(skill.Global) {
		t.Errorf("Claude skill still installed after uninstall --claude --global")
	}
}

func TestPromptsRequireTerminal(t *testing.T) {
	// Point stdin at a pipe so isInteractive() reports false deterministically,
	// regardless of how the test process was launched.
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	orig := os.Stdin
	os.Stdin = r
	t.Cleanup(func() { os.Stdin = orig; r.Close(); w.Close() })

	if isInteractive() {
		t.Fatalf("isInteractive() = true with a pipe stdin; test environment unexpected")
	}

	reader := bufio.NewReader(os.Stdin)
	if _, err := promptScope(reader, "Where?"); err == nil {
		t.Errorf("promptScope without a terminal: expected an error")
	}
	if _, err := promptAgents(reader, "Which?", "detected", skill.Agent.Detected); err == nil {
		t.Errorf("promptAgents without a terminal: expected an error")
	}
}

func mustSetFlag(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("setting flag: %v", err)
	}
}

func TestMigrateCommand(t *testing.T) {
	t.Chdir(t.TempDir())
	mute(t)
	initRepo(t)

	const header = "# Change Log\nAll notable changes.\n\n"
	const footer = "# Notes\n[Added] for new features.\n"

	// Base branch: a changelog with no staging entries.
	if err := os.WriteFile(changelogPath, []byte(header+footer), 0644); err != nil {
		t.Fatal(err)
	}
	runGit(t, "add", ".")
	runGit(t, "commit", "-m", "base")

	// Feature branch: add a staging entry.
	runGit(t, "checkout", "-b", "feature/foo")
	staged := header + "## [staging]\n### Added\n- New thing on this branch\n\n" + footer
	if err := os.WriteFile(changelogPath, []byte(staged), 0644); err != nil {
		t.Fatal(err)
	}
	runGit(t, "add", ".")
	runGit(t, "commit", "-m", "add staging entry")

	if err := migrateCmd.Flags().Set("base", "master"); err != nil {
		t.Fatal(err)
	}
	// confirm migrate, then decline auto-commit.
	withStdin(t, "y\nn\n")

	if err := runMigrate(migrateCmd, nil); err != nil {
		t.Fatalf("migrate: %v", err)
	}

	fragPath := filepath.Join(fragmentDir, "feature-foo.yaml")
	data, err := os.ReadFile(fragPath)
	if err != nil {
		t.Fatalf("expected fragment %s: %v", fragPath, err)
	}
	if !strings.Contains(string(data), "New thing on this branch") {
		t.Errorf("fragment missing migrated entry:\n%s", data)
	}

	// The staging entry should have been removed from CHANGELOG.md.
	out, _ := os.ReadFile(changelogPath)
	if strings.Contains(string(out), "New thing on this branch") {
		t.Errorf("staging entry should be removed from CHANGELOG.md:\n%s", out)
	}
}

// setReleaseFixtures writes a changelog and one valid fragment into the current
// directory and returns the fragment path.
func setReleaseFixtures(t *testing.T) string {
	t.Helper()
	if err := os.WriteFile(changelogPath, []byte(defaultChangelog), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(fragmentDir, 0755); err != nil {
		t.Fatal(err)
	}
	frag := filepath.Join(fragmentDir, "feature.yaml")
	if err := os.WriteFile(frag, []byte("added:\n  - \"New feature\"\n"), 0644); err != nil {
		t.Fatal(err)
	}
	return frag
}

func TestReleaseDirtyTreeAndVersionFiles(t *testing.T) {
	t.Chdir(t.TempDir())
	mute(t)
	initRepo(t)
	setReleaseFixtures(t)
	if err := os.WriteFile("VERSION", []byte("0.0.1\n"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile("package.json", []byte(`{"version": "0.0.1"}`), 0644); err != nil {
		t.Fatal(err)
	}
	// Nothing committed, so the working tree is dirty.

	// continue=y, version, metadata, proceed=y, VERSION=y, package.json=y, auto-commit=n
	withStdin(t, "y\n1.0.0\n\ny\ny\ny\nn\n")
	if err := runRelease(releaseCmd, nil); err != nil {
		t.Fatalf("release: %v", err)
	}

	if got, _ := os.ReadFile("VERSION"); string(got) != "1.0.0\n" {
		t.Errorf("VERSION = %q, want 1.0.0", got)
	}
	if got, _ := os.ReadFile("package.json"); !strings.Contains(string(got), "1.0.0") {
		t.Errorf("package.json not updated: %s", got)
	}
	// auto-commit declined: the fragment is deleted from disk, no commit made.
	if _, err := os.Stat(filepath.Join(fragmentDir, "feature.yaml")); !os.IsNotExist(err) {
		t.Errorf("fragment should be deleted from disk")
	}
}

func TestReleaseDirtyTreeAbort(t *testing.T) {
	t.Chdir(t.TempDir())
	mute(t)
	initRepo(t)
	frag := setReleaseFixtures(t)

	withStdin(t, "n\n") // decline to continue with a dirty tree
	if err := runRelease(releaseCmd, nil); err != nil {
		t.Fatalf("release (abort on dirty): %v", err)
	}
	if _, err := os.Stat(frag); err != nil {
		t.Errorf("fragment should be untouched after abort: %v", err)
	}
}

func TestReleaseEmptyVersion(t *testing.T) {
	t.Chdir(t.TempDir())
	mute(t)
	initRepo(t)
	setReleaseFixtures(t)
	runGit(t, "add", ".")
	runGit(t, "commit", "-m", "init") // clean tree, no dirty prompt

	withStdin(t, "\n") // empty version
	if err := runRelease(releaseCmd, nil); err == nil {
		t.Errorf("release with empty version: expected an error")
	}
}

func TestReleaseNoFragments(t *testing.T) {
	t.Chdir(t.TempDir())
	mute(t)
	initRepo(t)
	if err := os.WriteFile(changelogPath, []byte(defaultChangelog), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(fragmentDir, 0755); err != nil {
		t.Fatal(err)
	}
	// Only the sample is present; ReadAll skips it.
	if err := os.WriteFile(filepath.Join(fragmentDir, fragment.SampleFilename), fragment.Template(), 0644); err != nil {
		t.Fatal(err)
	}
	runGit(t, "add", ".")
	runGit(t, "commit", "-m", "init")

	if err := runRelease(releaseCmd, nil); err != nil {
		t.Fatalf("release (no fragments): %v", err)
	}
}

func TestReleaseAllEmptyFragments(t *testing.T) {
	t.Chdir(t.TempDir())
	mute(t)
	initRepo(t)
	if err := os.WriteFile(changelogPath, []byte(defaultChangelog), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(fragmentDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(fragmentDir, "empty.yaml"), fragment.Template(), 0644); err != nil {
		t.Fatal(err)
	}
	runGit(t, "add", ".")
	runGit(t, "commit", "-m", "init")

	if err := runRelease(releaseCmd, nil); err != nil {
		t.Fatalf("release (all empty): %v", err)
	}
}

func TestReleaseNoChangelogDir(t *testing.T) {
	t.Chdir(t.TempDir())
	mute(t)
	initRepo(t)
	if err := os.WriteFile(changelogPath, []byte(defaultChangelog), 0644); err != nil {
		t.Fatal(err)
	}
	runGit(t, "add", ".")
	runGit(t, "commit", "-m", "init")

	if err := runRelease(releaseCmd, nil); err == nil {
		t.Errorf("release without a changelog.d directory: expected an error")
	}
}

func TestReleaseValidationError(t *testing.T) {
	t.Chdir(t.TempDir())
	mute(t)
	initRepo(t)
	if err := os.WriteFile(changelogPath, []byte(defaultChangelog), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(fragmentDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(fragmentDir, "bad.yaml"), []byte("bogus:\n  - item\n"), 0644); err != nil {
		t.Fatal(err)
	}
	runGit(t, "add", ".")
	runGit(t, "commit", "-m", "init")

	if err := runRelease(releaseCmd, nil); err == nil {
		t.Errorf("release with an unknown category: expected an error")
	}
}

func TestReleaseMissingChangelog(t *testing.T) {
	t.Chdir(t.TempDir())
	mute(t)
	initRepo(t)
	runGit(t, "commit", "--allow-empty", "-m", "init") // clean tree, no CHANGELOG.md

	if err := runRelease(releaseCmd, nil); err == nil {
		t.Errorf("release without CHANGELOG.md: expected an error")
	}
}

func TestMigrateNoChangelog(t *testing.T) {
	t.Chdir(t.TempDir())
	mute(t)
	if err := runMigrate(migrateCmd, nil); err == nil {
		t.Errorf("migrate without CHANGELOG.md: expected an error")
	}
}

func TestMigrateNoStaging(t *testing.T) {
	t.Chdir(t.TempDir())
	mute(t)
	if err := os.WriteFile(changelogPath, []byte("# Change Log\n\n# Notes\n[Added] for features.\n"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := runMigrate(migrateCmd, nil); err != nil {
		t.Fatalf("migrate (no staging): %v", err)
	}
}

func TestMigrateNoNewEntries(t *testing.T) {
	t.Chdir(t.TempDir())
	mute(t)
	initRepo(t)

	const header = "# Change Log\n\n"
	const footer = "# Notes\n[Added] for features.\n"
	staged := header + "## [staging]\n### Added\n- Shared entry\n\n" + footer
	if err := os.WriteFile(changelogPath, []byte(staged), 0644); err != nil {
		t.Fatal(err)
	}
	runGit(t, "add", ".")
	runGit(t, "commit", "-m", "base with staging")
	// Feature branch carries the same staging entry — nothing new.
	runGit(t, "checkout", "-b", "feature/foo")
	mustSetFlag(t, migrateCmd.Flags().Set("base", "master"))

	if err := runMigrate(migrateCmd, nil); err != nil {
		t.Fatalf("migrate (no new entries): %v", err)
	}
}

func TestMigrateExistingFragmentAndCommit(t *testing.T) {
	t.Chdir(t.TempDir())
	mute(t)
	initRepo(t)

	const header = "# Change Log\n\n"
	const footer = "# Notes\n[Added] for features.\n"
	if err := os.WriteFile(changelogPath, []byte(header+footer), 0644); err != nil {
		t.Fatal(err)
	}
	runGit(t, "add", ".")
	runGit(t, "commit", "-m", "base")

	runGit(t, "checkout", "-b", "feature/foo")
	staged := header + "## [staging]\n### Added\n- Fresh entry\n\n" + footer
	if err := os.WriteFile(changelogPath, []byte(staged), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(fragmentDir, 0755); err != nil {
		t.Fatal(err)
	}
	// A fragment already exists for this branch; migrate must merge into it.
	if err := os.WriteFile(filepath.Join(fragmentDir, "feature-foo.yaml"), []byte("fixed:\n  - \"Existing fix\"\n"), 0644); err != nil {
		t.Fatal(err)
	}
	runGit(t, "add", ".")
	runGit(t, "commit", "-m", "feature work")
	mustSetFlag(t, migrateCmd.Flags().Set("base", "master"))

	withStdin(t, "y\ny\n") // confirm migrate, then auto-commit
	if err := runMigrate(migrateCmd, nil); err != nil {
		t.Fatalf("migrate: %v", err)
	}

	data, _ := os.ReadFile(filepath.Join(fragmentDir, "feature-foo.yaml"))
	if !strings.Contains(string(data), "Fresh entry") || !strings.Contains(string(data), "Existing fix") {
		t.Errorf("merged fragment missing entries:\n%s", data)
	}
	if msg := runGit(t, "log", "-1", "--pretty=%s"); msg != "Migrated changelog entries to changelog fragments" {
		t.Errorf("commit subject = %q", msg)
	}
}

func TestNewCommandWithEditor(t *testing.T) {
	t.Chdir(t.TempDir())
	mute(t)
	initRepo(t)
	runGit(t, "checkout", "-b", "feature/edit")
	runGit(t, "commit", "--allow-empty", "-m", "init")

	// `true` is a real program that ignores its args and exits 0.
	t.Setenv("EDITOR", "true")
	mustSetFlag(t, newCmd.Flags().Set("edit", "true"))
	t.Cleanup(func() { newCmd.Flags().Set("edit", "false") })

	// First run creates the fragment, then "opens" the editor.
	if err := newCmd.RunE(newCmd, nil); err != nil {
		t.Fatalf("new --edit (create): %v", err)
	}
	want := filepath.Join("changelog.d", "feature-edit.yaml")
	if _, err := os.Stat(want); err != nil {
		t.Errorf("expected %s: %v", want, err)
	}
	// Second run: the fragment already exists, --edit opens it again.
	if err := newCmd.RunE(newCmd, nil); err != nil {
		t.Fatalf("new --edit (existing): %v", err)
	}
}

func TestOpenEditorNoEditor(t *testing.T) {
	// With no editor configured and the well-known fallback absent, openEditor
	// reports an error rather than guessing.
	t.Setenv("EDITOR", "")
	t.Setenv("VISUAL", "")
	if _, err := os.Stat("/usr/bin/editor"); err == nil {
		t.Skip("/usr/bin/editor exists; cannot exercise the no-editor path here")
	}
	if err := openEditor("anything.yaml"); err == nil {
		t.Errorf("openEditor with no editor configured: expected an error")
	}
}

func TestSkillInstallBothAgentsBothScopes(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	t.Chdir(t.TempDir())
	mute(t)

	for _, f := range []string{"claude", "codex", "global", "project"} {
		mustSetFlag(t, skillInstallCmd.Flags().Set(f, "true"))
	}
	t.Cleanup(func() {
		for _, f := range []string{"claude", "codex", "global", "project"} {
			skillInstallCmd.Flags().Set(f, "false")
		}
	})

	if err := runSkillInstall(skillInstallCmd, nil); err != nil {
		t.Fatalf("skill install (both/both): %v", err)
	}
	for _, a := range []skill.Agent{skill.Claude, skill.Codex} {
		for _, s := range []skill.Scope{skill.Global, skill.Project} {
			if !a.Installed(s) {
				t.Errorf("%s not installed at %s", a.Display, scopeLabel(s))
			}
		}
	}
}

func TestSkillUninstallNotInstalled(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	mute(t)

	mustSetFlag(t, skillUninstallCmd.Flags().Set("codex", "true"))
	mustSetFlag(t, skillUninstallCmd.Flags().Set("global", "true"))
	t.Cleanup(func() {
		skillUninstallCmd.Flags().Set("codex", "false")
		skillUninstallCmd.Flags().Set("global", "false")
	})

	// Removing a skill that was never installed is a no-op, not an error.
	if err := runSkillUninstall(skillUninstallCmd, nil); err != nil {
		t.Fatalf("skill uninstall (absent): %v", err)
	}
}

func TestSkillUninstallModifiedKept(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	mute(t)

	// Install, then hand-edit so the skill diverges from the embedded copy.
	if _, err := skill.Claude.Install(skill.Global); err != nil {
		t.Fatal(err)
	}
	path, _ := skill.Claude.TargetPath(skill.Global)
	if err := os.WriteFile(path, []byte("hand-edited\n"), 0644); err != nil {
		t.Fatal(err)
	}

	mustSetFlag(t, skillUninstallCmd.Flags().Set("claude", "true"))
	mustSetFlag(t, skillUninstallCmd.Flags().Set("global", "true"))
	t.Cleanup(func() {
		skillUninstallCmd.Flags().Set("claude", "false")
		skillUninstallCmd.Flags().Set("global", "false")
	})

	// Without --force, a modified skill is kept and the command reports an error.
	if err := runSkillUninstall(skillUninstallCmd, nil); err == nil {
		t.Errorf("uninstall of a modified skill without --force: expected an error")
	}
	if _, err := os.Stat(path); err != nil {
		t.Errorf("modified skill should be kept without --force: %v", err)
	}
}

// fakeInteractive makes isInteractive() report true for the duration of a test
// so the prompt parsing logic can be exercised with a scripted reader.
func fakeInteractive(t *testing.T) {
	t.Helper()
	orig := isInteractive
	isInteractive = func() bool { return true }
	t.Cleanup(func() { isInteractive = orig })
}

func TestPromptScopeChoices(t *testing.T) {
	fakeInteractive(t)
	mute(t)

	cases := []struct {
		in   string
		want skill.Scope
	}{
		{"\n", skill.Global}, // default
		{"1\n", skill.Global},
		{"g\n", skill.Global},
		{"global\n", skill.Global},
		{"2\n", skill.Project},
		{"p\n", skill.Project},
		{"project\n", skill.Project},
	}
	for _, c := range cases {
		got, err := promptScope(bufio.NewReader(strings.NewReader(c.in)), "Where?")
		if err != nil {
			t.Fatalf("promptScope(%q): %v", c.in, err)
		}
		if got != c.want {
			t.Errorf("promptScope(%q) = %v, want %v", c.in, got, c.want)
		}
	}

	if _, err := promptScope(bufio.NewReader(strings.NewReader("nonsense\n")), "Where?"); err == nil {
		t.Errorf("promptScope(invalid): expected an error")
	}
}

func TestPromptScopeRequiresTerminal(t *testing.T) {
	mute(t)
	// isInteractive is not faked here, and the test's stdin is not a TTY.
	orig := isInteractive
	isInteractive = func() bool { return false }
	t.Cleanup(func() { isInteractive = orig })

	if _, err := promptScope(bufio.NewReader(strings.NewReader("1\n")), "Where?"); err == nil {
		t.Errorf("promptScope without a terminal: expected an error")
	}
}

func TestPromptAgentsSelection(t *testing.T) {
	fakeInteractive(t)
	mute(t)

	none := func(skill.Agent) bool { return false }
	keys := func(agents []skill.Agent) []string {
		out := make([]string, len(agents))
		for i, a := range agents {
			out[i] = a.Key
		}
		return out
	}

	cases := []struct {
		in   string
		want []string
	}{
		{"1\n", []string{"claude"}},
		{"2\n", []string{"codex"}},
		{"1,2\n", []string{"claude", "codex"}},
		{"claude\n", []string{"claude"}},
		{"codex,claude\n", []string{"codex", "claude"}},
		{"Claude\n", []string{"claude"}}, // matched by display name, case-insensitive
		{"1,1\n", []string{"claude"}},    // duplicates are removed
		{"claude,\n", []string{"claude"}}, // trailing empty token is ignored
	}
	for _, c := range cases {
		got, err := promptAgents(bufio.NewReader(strings.NewReader(c.in)), "Which?", "detected", none)
		if err != nil {
			t.Fatalf("promptAgents(%q): %v", c.in, err)
		}
		if !reflect.DeepEqual(keys(got), c.want) {
			t.Errorf("promptAgents(%q) = %v, want %v", c.in, keys(got), c.want)
		}
	}

	if _, err := promptAgents(bufio.NewReader(strings.NewReader("99\n")), "Which?", "detected", none); err == nil {
		t.Errorf("promptAgents(invalid): expected an error")
	}
}

func TestPromptAgentsDefaults(t *testing.T) {
	fakeInteractive(t)
	mute(t)

	// Pressing Enter selects the agents for which isDefault returns true.
	onlyClaude := func(a skill.Agent) bool { return a.Key == "claude" }
	got, err := promptAgents(bufio.NewReader(strings.NewReader("\n")), "Which?", "detected", onlyClaude)
	if err != nil {
		t.Fatalf("promptAgents (default): %v", err)
	}
	if len(got) != 1 || got[0].Key != "claude" {
		t.Errorf("promptAgents default = %+v, want [claude]", got)
	}
}

func TestExecuteHelp(t *testing.T) {
	mute(t)
	rootCmd.SetArgs([]string{"--help"})
	t.Cleanup(func() { rootCmd.SetArgs(nil) })
	if err := Execute(); err != nil {
		t.Fatalf("Execute --help: %v", err)
	}
}
