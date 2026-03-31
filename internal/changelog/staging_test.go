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
