package changelog

import (
	"strings"
	"testing"
)

const sampleChangelog = `# Change Log
All notable changes to this project will be documented in this file.
This project adheres to [Semantic Versioning](http://semver.org/).

## [3.81.0] - 2026-03-26 (86.8%)(Dev)
### Deployment
- Make sure to add new environment variable
### Added
- Added setting the licence key

## [3.80.0] - 2026-03-16 (86.8%)(Dev)(Prod)
### Added
- Added sonnet 4.6 model to scribe

# Notes
[Deployment] Notes for deployment
[Added] for new features.
`

func TestParseContent(t *testing.T) {
	cl := ParseContent(sampleChangelog)

	// Header should contain title and description
	if !strings.Contains(cl.Header, "# Change Log") {
		t.Error("Header should contain title")
	}
	if !strings.Contains(cl.Header, "Semantic Versioning") {
		t.Error("Header should contain semver reference")
	}

	// Header should NOT contain version entries
	if strings.Contains(cl.Header, "## [3.81.0]") {
		t.Error("Header should not contain version entries")
	}

	// Entries should contain version sections
	if !strings.Contains(cl.Entries, "## [3.81.0]") {
		t.Error("Entries should contain version 3.81.0")
	}
	if !strings.Contains(cl.Entries, "## [3.80.0]") {
		t.Error("Entries should contain version 3.80.0")
	}

	// Footer should contain Notes
	if !strings.Contains(cl.Footer, "# Notes") {
		t.Error("Footer should contain # Notes")
	}
	if !strings.Contains(cl.Footer, "[Added] for new features") {
		t.Error("Footer should contain notes content")
	}
}

func TestParseContentNoVersions(t *testing.T) {
	content := `# Change Log
All notable changes to this project will be documented in this file.

# Notes
[Added] for new features.
`
	cl := ParseContent(content)

	if cl.Entries != "" {
		t.Errorf("expected empty entries, got %q", cl.Entries)
	}
	if !strings.Contains(cl.Footer, "# Notes") {
		t.Error("Footer should contain # Notes")
	}
}

func TestParseContentNoFooter(t *testing.T) {
	content := `# Change Log

## [1.0.0] - 2026-01-01
### Added
- Initial release
`
	cl := ParseContent(content)

	if cl.Footer != "" {
		t.Errorf("expected empty footer, got %q", cl.Footer)
	}
	if !strings.Contains(cl.Entries, "## [1.0.0]") {
		t.Error("Entries should contain version 1.0.0")
	}
}

func TestLastVersion(t *testing.T) {
	tests := []struct {
		name    string
		entries string
		want    string
	}{
		{
			name:    "standard version",
			entries: "## [3.81.0] - 2026-03-26\n### Added\n- stuff\n",
			want:    "3.81.0",
		},
		{
			name:    "version with metadata",
			entries: "## [3.81.0] - 2026-03-26 (86.8%)(Dev)\n### Added\n- stuff\n",
			want:    "3.81.0",
		},
		{
			name:    "no version",
			entries: "",
			want:    "",
		},
		{
			name:    "multiple versions returns first",
			entries: "## [2.0.0] - 2026-02-01\n## [1.0.0] - 2026-01-01\n",
			want:    "2.0.0",
		},
		{
			name:    "staging is not a version",
			entries: "## [staging]\n## [1.0.0] - 2026-01-01\n",
			want:    "staging",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cl := &Changelog{Entries: tt.entries}
			got := LastVersion(cl)
			if got != tt.want {
				t.Errorf("LastVersion() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestInsert(t *testing.T) {
	cl := ParseContent(sampleChangelog)
	newEntry := "## [3.82.0] - 2026-03-27\n### Added\n- New feature\n"

	result := Insert(cl, newEntry)

	// New entry should appear before 3.81.0
	newIdx := strings.Index(result, "## [3.82.0]")
	oldIdx := strings.Index(result, "## [3.81.0]")
	if newIdx == -1 || oldIdx == -1 {
		t.Fatal("missing version entries in result")
	}
	if newIdx >= oldIdx {
		t.Error("new entry should appear before existing entries")
	}

	// Header should still be present
	if !strings.Contains(result, "# Change Log") {
		t.Error("Header missing from result")
	}

	// Footer should still be present
	if !strings.Contains(result, "# Notes") {
		t.Error("Footer missing from result")
	}
}

func TestInsertNoExistingEntries(t *testing.T) {
	content := `# Change Log
All notable changes to this project will be documented in this file.

# Notes
[Added] for new features.
`
	cl := ParseContent(content)
	newEntry := "## [1.0.0] - 2026-03-27\n### Added\n- Initial release\n"

	result := Insert(cl, newEntry)

	if !strings.Contains(result, "## [1.0.0]") {
		t.Error("new entry missing from result")
	}
	if !strings.Contains(result, "# Notes") {
		t.Error("Footer missing from result")
	}

	// Version should come after header, before footer
	vIdx := strings.Index(result, "## [1.0.0]")
	fIdx := strings.Index(result, "# Notes")
	if vIdx >= fIdx {
		t.Error("new entry should appear before footer")
	}
}
