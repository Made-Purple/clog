package fragment

import (
	"os"
	"path/filepath"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestTemplate(t *testing.T) {
	tmpl := Template()

	// Should be valid YAML
	var parsed map[string][]string
	if err := yaml.Unmarshal(tmpl, &parsed); err != nil {
		t.Fatalf("Template() produced invalid YAML: %v", err)
	}

	// Should contain all categories
	for _, cat := range CategoryOrder {
		entries, ok := parsed[cat]
		if !ok {
			t.Errorf("Template() missing category %q", cat)
			continue
		}
		if len(entries) != 1 || entries[0] != "" {
			t.Errorf("Template() category %q = %v, want [\"\"]", cat, entries)
		}
	}

	// Should not contain extra categories
	if len(parsed) != len(CategoryOrder) {
		t.Errorf("Template() has %d categories, want %d", len(parsed), len(CategoryOrder))
	}
}

func TestParse(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
		check   func(t *testing.T, f *Fragment)
	}{
		{
			name: "valid fragment",
			input: `added:
  - "Added a new feature"
fixed:
  - "Fixed a bug"
`,
			check: func(t *testing.T, f *Fragment) {
				if len(f.Entries["added"]) != 1 || f.Entries["added"][0] != "Added a new feature" {
					t.Errorf("added = %v", f.Entries["added"])
				}
				if len(f.Entries["fixed"]) != 1 || f.Entries["fixed"][0] != "Fixed a bug" {
					t.Errorf("fixed = %v", f.Entries["fixed"])
				}
			},
		},
		{
			name: "uppercase keys are lowercased",
			input: `Added:
  - "feature"
`,
			check: func(t *testing.T, f *Fragment) {
				if _, ok := f.Entries["added"]; !ok {
					t.Error("expected lowercase key 'added'")
				}
			},
		},
		{
			name:    "invalid YAML",
			input:   `not: valid: yaml: [`,
			wantErr: true,
		},
		{
			name:  "empty file",
			input: "",
			check: func(t *testing.T, f *Fragment) {
				if f.Entries != nil && len(f.Entries) != 0 {
					t.Errorf("expected nil or empty entries, got %v", f.Entries)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f, err := Parse("test.yaml", []byte(tt.input))
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if tt.check != nil {
				tt.check(t, f)
			}
		})
	}
}

func TestValidate(t *testing.T) {
	tests := []struct {
		name      string
		entries   map[string][]string
		wantCount int
	}{
		{
			name:      "all valid categories",
			entries:   map[string][]string{"added": {"x"}, "fixed": {"y"}},
			wantCount: 0,
		},
		{
			name:      "unknown category",
			entries:   map[string][]string{"added": {"x"}, "bogus": {"y"}},
			wantCount: 1,
		},
		{
			name:      "multiple unknown",
			entries:   map[string][]string{"foo": {"x"}, "bar": {"y"}},
			wantCount: 2,
		},
		{
			name:      "empty entries",
			entries:   map[string][]string{},
			wantCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &Fragment{Filename: "test.yaml", Entries: tt.entries}
			errs := Validate(f)
			if len(errs) != tt.wantCount {
				t.Errorf("got %d errors, want %d: %v", len(errs), tt.wantCount, errs)
			}
		})
	}
}

func TestNonEmptyEntries(t *testing.T) {
	tests := []struct {
		name    string
		entries map[string][]string
		want    map[string][]string
	}{
		{
			name:    "filters empty strings",
			entries: map[string][]string{"added": {"feature", "", "  "}, "fixed": {""}},
			want:    map[string][]string{"added": {"feature"}},
		},
		{
			name:    "keeps all non-empty",
			entries: map[string][]string{"added": {"one", "two"}},
			want:    map[string][]string{"added": {"one", "two"}},
		},
		{
			name:    "all empty returns empty map",
			entries: map[string][]string{"added": {""}, "fixed": {"", "  "}},
			want:    map[string][]string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &Fragment{Filename: "test.yaml", Entries: tt.entries}
			got := NonEmptyEntries(f)

			if len(got) != len(tt.want) {
				t.Fatalf("got %d categories, want %d", len(got), len(tt.want))
			}
			for cat, wantEntries := range tt.want {
				gotEntries, ok := got[cat]
				if !ok {
					t.Errorf("missing category %q", cat)
					continue
				}
				if len(gotEntries) != len(wantEntries) {
					t.Errorf("category %q: got %d entries, want %d", cat, len(gotEntries), len(wantEntries))
				}
			}
		})
	}
}

func TestReadAll(t *testing.T) {
	dir := t.TempDir()

	// Write a valid fragment
	validYAML := []byte("added:\n  - \"A feature\"\n")
	if err := os.WriteFile(filepath.Join(dir, "valid.yaml"), validYAML, 0644); err != nil {
		t.Fatal(err)
	}

	// Write an invalid fragment
	invalidYAML := []byte("not: valid: yaml: [")
	if err := os.WriteFile(filepath.Join(dir, "invalid.yaml"), invalidYAML, 0644); err != nil {
		t.Fatal(err)
	}

	// Write a non-yaml file (should be ignored)
	if err := os.WriteFile(filepath.Join(dir, "readme.txt"), []byte("ignore me"), 0644); err != nil {
		t.Fatal(err)
	}

	fragments, errs := ReadAll(dir)

	if len(fragments) != 1 {
		t.Errorf("got %d fragments, want 1", len(fragments))
	}
	if len(errs) != 1 {
		t.Errorf("got %d errors, want 1", len(errs))
	}
	if len(fragments) > 0 && fragments[0].Filename != "valid.yaml" {
		t.Errorf("fragment filename = %q, want %q", fragments[0].Filename, "valid.yaml")
	}
}

func TestReadAllMissingDir(t *testing.T) {
	fragments, errs := ReadAll("/nonexistent/path")
	if len(fragments) != 0 {
		t.Errorf("expected no fragments, got %d", len(fragments))
	}
	if len(errs) != 1 {
		t.Errorf("expected 1 error, got %d", len(errs))
	}
}
