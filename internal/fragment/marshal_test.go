package fragment

import (
	"strings"
	"testing"
)

func TestMarshalEntries(t *testing.T) {
	t.Run("empty entries produces template-like output", func(t *testing.T) {
		result := string(MarshalEntries(map[string][]string{}))
		for _, cat := range CategoryOrder {
			if !strings.Contains(result, cat+":") {
				t.Errorf("missing category %q", cat)
			}
		}
		if strings.Count(result, "- \"\"") != len(CategoryOrder) {
			t.Error("expected all categories to have empty placeholder")
		}
	})

	t.Run("populated categories", func(t *testing.T) {
		entries := map[string][]string{
			"changed": {"Updated the UI", "Changed the API"},
			"fixed":   {"Bug fix"},
		}
		result := string(MarshalEntries(entries))

		if !strings.Contains(result, "Updated the UI") {
			t.Error("missing entry 'Updated the UI'")
		}
		if !strings.Contains(result, "Changed the API") {
			t.Error("missing entry 'Changed the API'")
		}
		if !strings.Contains(result, "Bug fix") {
			t.Error("missing entry 'Bug fix'")
		}

		// Empty categories should still have placeholder
		if !strings.Contains(result, "deployment:\n  - \"\"") {
			t.Error("empty category missing placeholder")
		}
	})

	t.Run("roundtrip parse", func(t *testing.T) {
		entries := map[string][]string{
			"added":   {"New feature"},
			"changed": {"Updated something"},
		}
		data := MarshalEntries(entries)

		frag, err := Parse("test.yaml", data)
		if err != nil {
			t.Fatalf("parse error: %v", err)
		}

		nonEmpty := NonEmptyEntries(frag)
		if len(nonEmpty["added"]) != 1 || nonEmpty["added"][0] != "New feature" {
			t.Errorf("roundtrip failed for added: %v", nonEmpty["added"])
		}
		if len(nonEmpty["changed"]) != 1 || nonEmpty["changed"][0] != "Updated something" {
			t.Errorf("roundtrip failed for changed: %v", nonEmpty["changed"])
		}
	})
}
