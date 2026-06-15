package changelog

import (
	"os"
	"path/filepath"
	"testing"
)

func TestReadParsesFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "CHANGELOG.md")
	content := "# Change Log\n\n## [1.0.0] - 2026-01-01\n### Added\n- Feature\n\n# Notes\n[Added] for new features.\n"
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("write: %v", err)
	}

	cl, err := Read(path)
	if err != nil {
		t.Fatalf("Read: %v", err)
	}
	if got := LastVersion(cl); got != "1.0.0" {
		t.Errorf("LastVersion = %q, want 1.0.0", got)
	}
	if cl.Footer == "" {
		t.Errorf("Footer was not parsed")
	}
}

func TestReadMissingFile(t *testing.T) {
	_, err := Read(filepath.Join(t.TempDir(), "nope.md"))
	if err == nil {
		t.Errorf("Read of a missing file: expected an error")
	}
}
