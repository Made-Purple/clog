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
