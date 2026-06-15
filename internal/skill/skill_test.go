package skill

import (
	"os"
	"path/filepath"
	"testing"
)

func TestTargetPath(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	global, err := Claude.TargetPath(Global)
	if err != nil {
		t.Fatalf("global TargetPath: %v", err)
	}
	want := filepath.Join(home, ".claude", "skills", "clog", "SKILL.md")
	if global != want {
		t.Errorf("global path = %q, want %q", global, want)
	}

	project, err := Codex.TargetPath(Project)
	if err != nil {
		t.Fatalf("project TargetPath: %v", err)
	}
	if project != filepath.Join(".codex", "skills", "clog", "SKILL.md") {
		t.Errorf("project path = %q, want relative .codex/skills/clog/SKILL.md", project)
	}
}

func TestInstallWritesAndIsIdempotent(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	// First install writes the embedded content and reports an update.
	res, err := Claude.Install(Global)
	if err != nil {
		t.Fatalf("install: %v", err)
	}
	if !res.Updated {
		t.Errorf("first install: Updated = false, want true")
	}
	got, err := os.ReadFile(res.Path)
	if err != nil {
		t.Fatalf("reading installed skill: %v", err)
	}
	if string(got) != string(content) {
		t.Errorf("installed content does not match embedded SKILL.md")
	}

	// Second install with identical content is a no-op.
	res2, err := Claude.Install(Global)
	if err != nil {
		t.Fatalf("re-install: %v", err)
	}
	if res2.Updated {
		t.Errorf("second install: Updated = true, want false (already up to date)")
	}

	// Divergent content is overwritten and reported as updated.
	if err := os.WriteFile(res.Path, []byte("tampered\n"), 0644); err != nil {
		t.Fatalf("tampering: %v", err)
	}
	res3, err := Claude.Install(Global)
	if err != nil {
		t.Fatalf("re-install after tamper: %v", err)
	}
	if !res3.Updated {
		t.Errorf("install after tamper: Updated = false, want true")
	}
}

func TestUninstall(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	// Removing when nothing is installed reports neither existed nor removed.
	res, err := Claude.Uninstall(Global, false)
	if err != nil {
		t.Fatalf("uninstall (absent): %v", err)
	}
	if res.Existed || res.Removed {
		t.Errorf("absent uninstall = %+v, want Existed=false Removed=false", res)
	}

	// Install, then uninstall: the file and its skills/clog dir are gone, but
	// skills/ remains.
	if _, err := Claude.Install(Global); err != nil {
		t.Fatalf("install: %v", err)
	}
	res, err = Claude.Uninstall(Global, false)
	if err != nil {
		t.Fatalf("uninstall: %v", err)
	}
	if !res.Removed {
		t.Errorf("uninstall: Removed = false, want true")
	}
	if _, err := os.Stat(res.Path); !os.IsNotExist(err) {
		t.Errorf("SKILL.md still present after uninstall")
	}
	clogDir := filepath.Join(home, ".claude", "skills", "clog")
	if _, err := os.Stat(clogDir); !os.IsNotExist(err) {
		t.Errorf("empty skills/clog dir was not pruned")
	}
	skillsDir := filepath.Join(home, ".claude", "skills")
	if _, err := os.Stat(skillsDir); err != nil {
		t.Errorf("skills/ dir should be left in place, got: %v", err)
	}
}

func TestUninstallCustomizedRequiresForce(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	res, err := Codex.Install(Project)
	if err != nil {
		t.Fatalf("install: %v", err)
	}
	if err := os.WriteFile(res.Path, []byte("hand-edited\n"), 0644); err != nil {
		t.Fatalf("tamper: %v", err)
	}

	// Without force, a modified skill is kept.
	got, err := Codex.Uninstall(Project, false)
	if err != nil {
		t.Fatalf("uninstall (no force): %v", err)
	}
	if got.Removed || !got.Customized || !got.Existed {
		t.Errorf("modified uninstall = %+v, want Existed=true Customized=true Removed=false", got)
	}
	if _, err := os.Stat(res.Path); err != nil {
		t.Errorf("modified skill should be kept without --force, got: %v", err)
	}

	// With force, it is removed.
	got, err = Codex.Uninstall(Project, true)
	if err != nil {
		t.Fatalf("uninstall (force): %v", err)
	}
	if !got.Removed {
		t.Errorf("force uninstall: Removed = false, want true")
	}
}

func TestInstalled(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	if Claude.Installed(Global) {
		t.Errorf("Installed = true before install")
	}
	if _, err := Claude.Install(Global); err != nil {
		t.Fatal(err)
	}
	if !Claude.Installed(Global) {
		t.Errorf("Installed = false after install")
	}
}

func TestDetected(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	if Claude.Detected() {
		t.Errorf("Claude.Detected() = true before ~/.claude exists")
	}
	if err := os.MkdirAll(filepath.Join(home, ".claude"), 0755); err != nil {
		t.Fatal(err)
	}
	if !Claude.Detected() {
		t.Errorf("Claude.Detected() = false after ~/.claude created")
	}
}

func TestAgentByKey(t *testing.T) {
	if a, ok := AgentByKey("codex"); !ok || a.Key != "codex" {
		t.Errorf("AgentByKey(codex) = %+v, %v", a, ok)
	}
	if _, ok := AgentByKey("nope"); ok {
		t.Errorf("AgentByKey(nope) ok = true, want false")
	}
}
