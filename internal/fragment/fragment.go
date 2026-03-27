package fragment

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// CategoryOrder defines the fixed display order for changelog categories.
var CategoryOrder = []string{
	"deployment",
	"added",
	"changed",
	"deprecated",
	"removed",
	"fixed",
	"security",
	"yanked",
}

// CategoryDisplay maps lowercase category names to their display form.
var CategoryDisplay = map[string]string{
	"deployment": "Deployment",
	"added":      "Added",
	"changed":    "Changed",
	"deprecated": "Deprecated",
	"removed":    "Removed",
	"fixed":      "Fixed",
	"security":   "Security",
	"yanked":     "YANKED",
}

// allowedCategories is a set for O(1) lookup.
var allowedCategories = func() map[string]bool {
	m := make(map[string]bool, len(CategoryOrder))
	for _, c := range CategoryOrder {
		m[c] = true
	}
	return m
}()

// SampleFilename is the name of the sample fragment file created by init.
// It is skipped during ReadAll so it doesn't appear in releases.
const SampleFilename = "sample.yaml"

// Fragment represents a parsed changelog fragment file.
type Fragment struct {
	Filename string
	Entries  map[string][]string
}

// Template returns the YAML template bytes with all categories pre-populated
// with a single empty string entry.
func Template() []byte {
	var b strings.Builder
	for _, cat := range CategoryOrder {
		b.WriteString(cat)
		b.WriteString(":\n  - \"\"\n")
	}
	return []byte(b.String())
}

// Parse parses YAML data into a Fragment.
func Parse(filename string, data []byte) (*Fragment, error) {
	var raw map[string][]string
	if err := yaml.Unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("%s: %w", filename, err)
	}

	entries := make(map[string][]string, len(raw))
	for k, v := range raw {
		entries[strings.ToLower(k)] = v
	}

	return &Fragment{
		Filename: filename,
		Entries:  entries,
	}, nil
}

// Validate checks that all category keys in a fragment are from the allowed set.
// Returns a slice of errors (empty if valid).
func Validate(f *Fragment) []error {
	var errs []error
	for k := range f.Entries {
		if !allowedCategories[k] {
			errs = append(errs, fmt.Errorf("%s: unknown category %q", f.Filename, k))
		}
	}
	return errs
}

// NonEmptyEntries returns a copy of the fragment's entries with blank/empty
// strings filtered out. Categories where all entries are empty are omitted.
func NonEmptyEntries(f *Fragment) map[string][]string {
	result := make(map[string][]string)
	for cat, entries := range f.Entries {
		var filtered []string
		for _, e := range entries {
			if strings.TrimSpace(e) != "" {
				filtered = append(filtered, e)
			}
		}
		if len(filtered) > 0 {
			result[cat] = filtered
		}
	}
	return result
}

// ReadAll reads and parses all .yaml files from the given directory.
// Returns all successfully parsed fragments and any errors encountered.
func ReadAll(dir string) ([]*Fragment, []error) {
	dirEntries, err := os.ReadDir(dir)
	if err != nil {
		return nil, []error{fmt.Errorf("reading %s: %w", dir, err)}
	}

	var fragments []*Fragment
	var errs []error

	for _, entry := range dirEntries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".yaml") || entry.Name() == SampleFilename {
			continue
		}

		path := filepath.Join(dir, entry.Name())
		data, err := os.ReadFile(path)
		if err != nil {
			errs = append(errs, fmt.Errorf("reading %s: %w", path, err))
			continue
		}

		frag, err := Parse(entry.Name(), data)
		if err != nil {
			errs = append(errs, err)
			continue
		}

		fragments = append(fragments, frag)
	}

	return fragments, errs
}
