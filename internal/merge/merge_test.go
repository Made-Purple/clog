package merge

import (
	"strings"
	"testing"

	"github.com/made-purple/clog/internal/fragment"
)

func TestMergeEmpty(t *testing.T) {
	result := Merge(nil)
	if len(result) != 0 {
		t.Errorf("expected empty map, got %v", result)
	}
}

func TestMergeSingle(t *testing.T) {
	f := &fragment.Fragment{
		Filename: "test.yaml",
		Entries: map[string][]string{
			"added": {"Feature A", ""},
			"fixed": {"Bug fix"},
			"security": {""},
		},
	}

	result := Merge([]*fragment.Fragment{f})

	if len(result["added"]) != 1 || result["added"][0] != "Feature A" {
		t.Errorf("added = %v, want [Feature A]", result["added"])
	}
	if len(result["fixed"]) != 1 {
		t.Errorf("fixed = %v, want [Bug fix]", result["fixed"])
	}
	if _, ok := result["security"]; ok {
		t.Error("security should be omitted (all empty)")
	}
}

func TestMergeMultiple(t *testing.T) {
	f1 := &fragment.Fragment{
		Filename: "a.yaml",
		Entries: map[string][]string{
			"added": {"Feature A"},
			"fixed": {"Bug 1"},
		},
	}
	f2 := &fragment.Fragment{
		Filename: "b.yaml",
		Entries: map[string][]string{
			"added": {"Feature B"},
			"changed": {"Change 1"},
		},
	}

	result := Merge([]*fragment.Fragment{f1, f2})

	if len(result["added"]) != 2 {
		t.Errorf("added has %d entries, want 2", len(result["added"]))
	}
	if len(result["fixed"]) != 1 {
		t.Errorf("fixed has %d entries, want 1", len(result["fixed"]))
	}
	if len(result["changed"]) != 1 {
		t.Errorf("changed has %d entries, want 1", len(result["changed"]))
	}
}

func TestRender(t *testing.T) {
	merged := map[string][]string{
		"added":   {"Feature A", "Feature B"},
		"fixed":   {"Bug fix"},
		"changed": {"Change 1"},
	}

	result := Render("3.82.0", "2026-03-27", "", merged)

	// Check version line
	if !strings.HasPrefix(result, "## [3.82.0] - 2026-03-27\n") {
		t.Errorf("unexpected version line: %s", strings.SplitN(result, "\n", 2)[0])
	}

	// Check category order: Added should come before Changed, Changed before Fixed
	addedIdx := strings.Index(result, "### Added")
	changedIdx := strings.Index(result, "### Changed")
	fixedIdx := strings.Index(result, "### Fixed")

	if addedIdx == -1 || changedIdx == -1 || fixedIdx == -1 {
		t.Fatalf("missing categories in output:\n%s", result)
	}
	if addedIdx >= changedIdx {
		t.Error("Added should come before Changed")
	}
	if changedIdx >= fixedIdx {
		t.Error("Changed should come before Fixed")
	}

	// Check entries are present
	if !strings.Contains(result, "- Feature A") {
		t.Error("missing Feature A entry")
	}
	if !strings.Contains(result, "- Feature B") {
		t.Error("missing Feature B entry")
	}
}

func TestRenderWithMetadata(t *testing.T) {
	merged := map[string][]string{
		"added": {"Feature A"},
	}

	result := Render("3.82.0", "2026-03-27", "(98%)(Dev)", merged)

	if !strings.HasPrefix(result, "## [3.82.0] - 2026-03-27 (98%)(Dev)\n") {
		t.Errorf("unexpected version line: %s", strings.SplitN(result, "\n", 2)[0])
	}
}

func TestRenderEmptyMerge(t *testing.T) {
	result := Render("1.0.0", "2026-01-01", "", map[string][]string{})

	// Should only have the version line
	lines := strings.Split(strings.TrimSpace(result), "\n")
	if len(lines) != 1 {
		t.Errorf("expected 1 line, got %d: %v", len(lines), lines)
	}
}

func TestRenderCategoryOrder(t *testing.T) {
	// Provide categories in reverse order, ensure they render in correct order
	merged := map[string][]string{
		"yanked":     {"Yanked item"},
		"security":   {"Security item"},
		"fixed":      {"Fixed item"},
		"removed":    {"Removed item"},
		"deprecated": {"Deprecated item"},
		"changed":    {"Changed item"},
		"added":      {"Added item"},
		"deployment": {"Deploy note"},
	}

	result := Render("1.0.0", "2026-01-01", "", merged)

	categories := []string{
		"### Deployment",
		"### Added",
		"### Changed",
		"### Deprecated",
		"### Removed",
		"### Fixed",
		"### Security",
		"### YANKED",
	}

	prevIdx := -1
	for _, cat := range categories {
		idx := strings.Index(result, cat)
		if idx == -1 {
			t.Errorf("missing category %s in output", cat)
			continue
		}
		if idx <= prevIdx {
			t.Errorf("category %s is out of order", cat)
		}
		prevIdx = idx
	}
}
