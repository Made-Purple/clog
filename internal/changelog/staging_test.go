package changelog

import (
	"strings"
	"testing"
)

func TestExtractStaging(t *testing.T) {
	tests := []struct {
		name       string
		entries    string
		wantNil    bool
		wantErr    bool
		wantCats   []string
		wantCounts map[string]int
	}{
		{
			name:    "no entries",
			entries: "",
			wantNil: true,
		},
		{
			name:    "no staging section",
			entries: "## [1.0.0] - 2026-01-01\n### Added\n- Feature\n",
			wantNil: true,
		},
		{
			name:       "single category",
			entries:    "## [staging]\n### Changed\n- Updated the UI\n\n## [1.0.0] - 2026-01-01\n### Added\n- Feature\n",
			wantCats:   []string{"changed"},
			wantCounts: map[string]int{"changed": 1},
		},
		{
			name:       "multiple categories",
			entries:    "## [staging]\n### Added\n- New thing\n### Fixed\n- Bug fix\n- Another fix\n\n## [1.0.0] - 2026-01-01\n",
			wantCats:   []string{"added", "fixed"},
			wantCounts: map[string]int{"added": 1, "fixed": 2},
		},
		{
			name:       "staging is only entry",
			entries:    "## [staging]\n### Changed\n- First change\n",
			wantCats:   []string{"changed"},
			wantCounts: map[string]int{"changed": 1},
		},
		{
			name:       "case insensitive header",
			entries:    "## [Staging]\n### Added\n- Item\n",
			wantCats:   []string{"added"},
			wantCounts: map[string]int{"added": 1},
		},
		{
			name:    "unknown category",
			entries: "## [staging]\n### Bogus\n- Item\n",
			wantErr: true,
		},
		{
			name:    "empty staging section",
			entries: "## [staging]\n\n## [1.0.0] - 2026-01-01\n",
			wantNil: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cl := &Changelog{Entries: tt.entries}
			result, err := ExtractStaging(cl)

			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if tt.wantNil {
				if result != nil {
					t.Fatalf("expected nil, got %v", result)
				}
				return
			}
			if result == nil {
				t.Fatal("expected non-nil result")
			}

			for _, cat := range tt.wantCats {
				if _, ok := result[cat]; !ok {
					t.Errorf("missing category %q", cat)
				}
			}
			for cat, count := range tt.wantCounts {
				if len(result[cat]) != count {
					t.Errorf("category %q: got %d entries, want %d", cat, len(result[cat]), count)
				}
			}
		})
	}
}

func TestRemoveStaging(t *testing.T) {
	tests := []struct {
		name    string
		header  string
		entries string
		footer  string
		check   func(t *testing.T, result string)
	}{
		{
			name:    "remove staging from top",
			header:  "# Changelog\n\n",
			entries: "## [staging]\n### Changed\n- Item\n\n## [1.0.0] - 2026-01-01\n### Added\n- Feature\n",
			footer:  "",
			check: func(t *testing.T, result string) {
				if strings.Contains(result, "[staging]") {
					t.Error("staging section not removed")
				}
				if !strings.Contains(result, "## [1.0.0]") {
					t.Error("existing entry removed")
				}
			},
		},
		{
			name:    "remove only staging entry",
			header:  "# Changelog\n\n",
			entries: "## [staging]\n### Changed\n- Item\n",
			footer:  "",
			check: func(t *testing.T, result string) {
				if strings.Contains(result, "[staging]") {
					t.Error("staging section not removed")
				}
				if !strings.Contains(result, "# Changelog") {
					t.Error("header removed")
				}
			},
		},
		{
			name:    "no staging section",
			header:  "# Changelog\n\n",
			entries: "## [1.0.0] - 2026-01-01\n### Added\n- Feature\n",
			footer:  "",
			check: func(t *testing.T, result string) {
				if !strings.Contains(result, "## [1.0.0]") {
					t.Error("entry removed")
				}
			},
		},
		{
			name:    "preserves footer",
			header:  "# Changelog\n\n",
			entries: "## [staging]\n### Added\n- Item\n\n## [1.0.0] - 2026-01-01\n### Added\n- Feature\n",
			footer:  "# Notes\nSome notes\n",
			check: func(t *testing.T, result string) {
				if !strings.Contains(result, "# Notes") {
					t.Error("footer removed")
				}
				if strings.Contains(result, "[staging]") {
					t.Error("staging section not removed")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cl := &Changelog{
				Header:  tt.header,
				Entries: tt.entries,
				Footer:  tt.footer,
			}
			result := RemoveStaging(cl)
			tt.check(t, result)
		})
	}
}

func TestRemoveStagingEntries(t *testing.T) {
	tests := []struct {
		name     string
		header   string
		entries  string
		footer   string
		toRemove map[string][]string
		check    func(t *testing.T, result string)
	}{
		{
			name:    "partial removal within a category keeps header and other entries",
			header:  "# Changelog\n\n",
			entries: "## [staging]\n### Changed\n- Keep me\n- Remove me\n- Also keep\n\n## [1.0.0] - 2026-01-01\n### Added\n- Feature\n",
			toRemove: map[string][]string{
				"changed": {"Remove me"},
			},
			check: func(t *testing.T, result string) {
				if !strings.Contains(result, "[staging]") {
					t.Error("staging section dropped when it should remain")
				}
				if !strings.Contains(result, "### Changed") {
					t.Error("category header dropped when entries remain")
				}
				if !strings.Contains(result, "- Keep me") || !strings.Contains(result, "- Also keep") {
					t.Error("non-removed entries were dropped")
				}
				if strings.Contains(result, "- Remove me") {
					t.Error("targeted entry was not removed")
				}
				if !strings.Contains(result, "## [1.0.0]") {
					t.Error("existing version entry was removed")
				}
			},
		},
		{
			name:    "removing all entries in a category drops its header",
			header:  "# Changelog\n\n",
			entries: "## [staging]\n### Changed\n- Only one\n### Fixed\n- Keep me\n\n## [1.0.0] - 2026-01-01\n",
			toRemove: map[string][]string{
				"changed": {"Only one"},
			},
			check: func(t *testing.T, result string) {
				if !strings.Contains(result, "[staging]") {
					t.Error("staging section dropped prematurely")
				}
				if strings.Contains(result, "### Changed") {
					t.Error("empty category header was not dropped")
				}
				if !strings.Contains(result, "### Fixed") || !strings.Contains(result, "- Keep me") {
					t.Error("other category was incorrectly affected")
				}
			},
		},
		{
			name:    "removing every entry drops the whole staging section",
			header:  "# Changelog\n\n",
			entries: "## [staging]\n### Added\n- A\n### Fixed\n- B\n\n## [1.0.0] - 2026-01-01\n### Added\n- Feature\n",
			toRemove: map[string][]string{
				"added": {"A"},
				"fixed": {"B"},
			},
			check: func(t *testing.T, result string) {
				if strings.Contains(result, "[staging]") {
					t.Error("empty staging section was not removed")
				}
				if !strings.Contains(result, "## [1.0.0]") {
					t.Error("existing version entry was removed")
				}
			},
		},
		{
			name:    "empty toRemove is a no-op",
			header:  "# Changelog\n\n",
			entries: "## [staging]\n### Added\n- A\n\n## [1.0.0] - 2026-01-01\n",
			toRemove: map[string][]string{},
			check: func(t *testing.T, result string) {
				if !strings.Contains(result, "- A") {
					t.Error("no-op dropped an entry")
				}
				if !strings.Contains(result, "[staging]") {
					t.Error("no-op dropped the staging section")
				}
			},
		},
		{
			name:    "entry not present in staging is ignored",
			header:  "# Changelog\n\n",
			entries: "## [staging]\n### Added\n- Real entry\n",
			toRemove: map[string][]string{
				"added": {"Ghost entry"},
			},
			check: func(t *testing.T, result string) {
				if !strings.Contains(result, "- Real entry") {
					t.Error("real entry was incorrectly removed")
				}
			},
		},
		{
			name:    "preserves footer",
			header:  "# Changelog\n\n",
			entries: "## [staging]\n### Added\n- A\n",
			footer:  "# Notes\nSome notes\n",
			toRemove: map[string][]string{
				"added": {"A"},
			},
			check: func(t *testing.T, result string) {
				if !strings.Contains(result, "# Notes") {
					t.Error("footer was removed")
				}
				if strings.Contains(result, "[staging]") {
					t.Error("empty staging section was not removed")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cl := &Changelog{
				Header:  tt.header,
				Entries: tt.entries,
				Footer:  tt.footer,
			}
			result := RemoveStagingEntries(cl, tt.toRemove)
			tt.check(t, result)
		})
	}
}

func TestRemoveStagingEntriesPreservesLayout(t *testing.T) {
	t.Run("no blank lines between sections stays compact", func(t *testing.T) {
		cl := &Changelog{
			Header: "# Changelog\n\n",
			Entries: "## [staging]\n" +
				"### Deployment\n" +
				"- Deploy note\n" +
				"### Added\n" +
				"- Added thing\n" +
				"### Changed\n" +
				"- Changed thing\n" +
				"### Security\n" +
				"- Security thing\n",
		}
		result := RemoveStagingEntries(cl, map[string][]string{
			"changed": {"Changed thing"},
		})
		want := "# Changelog\n\n" +
			"## [staging]\n" +
			"### Deployment\n" +
			"- Deploy note\n" +
			"### Added\n" +
			"- Added thing\n" +
			"### Security\n" +
			"- Security thing\n"
		if result != want {
			t.Errorf("layout not preserved.\nwant:\n%q\ngot:\n%q", want, result)
		}
	})

	t.Run("blank lines between sections are preserved", func(t *testing.T) {
		cl := &Changelog{
			Header: "# Changelog\n\n",
			Entries: "## [staging]\n" +
				"### Added\n" +
				"- A\n" +
				"\n" +
				"### Changed\n" +
				"- B\n" +
				"\n" +
				"### Fixed\n" +
				"- C\n",
		}
		result := RemoveStagingEntries(cl, map[string][]string{
			"changed": {"B"},
		})
		if !strings.Contains(result, "- A\n\n### Fixed") {
			t.Errorf("expected blank line preserved between Added and Fixed, got:\n%s", result)
		}
		if strings.Contains(result, "### Changed") || strings.Contains(result, "- B") {
			t.Errorf("Changed section should be removed, got:\n%s", result)
		}
	})

	t.Run("partial removal in middle category leaves header", func(t *testing.T) {
		cl := &Changelog{
			Header: "# Changelog\n\n",
			Entries: "## [staging]\n" +
				"### Added\n" +
				"- A\n" +
				"### Changed\n" +
				"- Keep\n" +
				"- Drop\n" +
				"### Fixed\n" +
				"- F\n",
		}
		result := RemoveStagingEntries(cl, map[string][]string{
			"changed": {"Drop"},
		})
		want := "# Changelog\n\n" +
			"## [staging]\n" +
			"### Added\n" +
			"- A\n" +
			"### Changed\n" +
			"- Keep\n" +
			"### Fixed\n" +
			"- F\n"
		if result != want {
			t.Errorf("partial removal corrupted layout.\nwant:\n%q\ngot:\n%q", want, result)
		}
	})
}
