package changelog

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/made-purple/clog/internal/fragment"
)

var stagingLineRe = regexp.MustCompile(`(?i)^## \[staging\]`)
var categoryLineRe = regexp.MustCompile(`^### (.+)`)

// categoryReverse builds a reverse lookup from display name (lowercased) to category key.
func categoryReverse() map[string]string {
	m := make(map[string]string, len(fragment.CategoryDisplay))
	for key, display := range fragment.CategoryDisplay {
		m[strings.ToLower(display)] = key
	}
	return m
}

// ExtractStaging parses the ## [staging] section and returns its entries
// as a category->entries map. Returns nil, nil if no staging section exists.
func ExtractStaging(cl *Changelog) (map[string][]string, error) {
	if cl.Entries == "" {
		return nil, nil
	}

	lines := strings.Split(cl.Entries, "\n")
	startIdx := -1
	for i, line := range lines {
		if stagingLineRe.MatchString(line) {
			startIdx = i
			break
		}
	}
	if startIdx == -1 {
		return nil, nil
	}

	endIdx := len(lines)
	for i := startIdx + 1; i < len(lines); i++ {
		if versionLineRe.MatchString(lines[i]) {
			endIdx = i
			break
		}
	}

	reverse := categoryReverse()
	result := make(map[string][]string)
	var currentCat string

	for i := startIdx + 1; i < endIdx; i++ {
		line := lines[i]
		if m := categoryLineRe.FindStringSubmatch(line); m != nil {
			catDisplay := strings.TrimSpace(m[1])
			cat, ok := reverse[strings.ToLower(catDisplay)]
			if !ok {
				return nil, fmt.Errorf("unknown category in staging section: %q", catDisplay)
			}
			currentCat = cat
			continue
		}
		if currentCat != "" && strings.HasPrefix(line, "- ") {
			entry := strings.TrimPrefix(line, "- ")
			if strings.TrimSpace(entry) != "" {
				result[currentCat] = append(result[currentCat], entry)
			}
		}
	}

	if len(result) == 0 {
		return nil, nil
	}

	return result, nil
}

// RemoveStagingEntries removes the given entries (keyed by category) from the
// ## [staging] section, leaving any other staging entries intact. If a category
// becomes empty its header is dropped; if staging becomes empty the whole section
// is removed. Returns the updated full content string.
func RemoveStagingEntries(cl *Changelog, toRemove map[string][]string) string {
	if cl.Entries == "" || len(toRemove) == 0 {
		return buildContent(cl.Header, cl.Entries, cl.Footer)
	}

	lines := strings.Split(cl.Entries, "\n")
	startIdx := -1
	for i, line := range lines {
		if stagingLineRe.MatchString(line) {
			startIdx = i
			break
		}
	}
	if startIdx == -1 {
		return buildContent(cl.Header, cl.Entries, cl.Footer)
	}

	endIdx := len(lines)
	for i := startIdx + 1; i < len(lines); i++ {
		if versionLineRe.MatchString(lines[i]) {
			endIdx = i
			break
		}
	}

	reverse := categoryReverse()
	removeSets := make(map[string]map[string]bool, len(toRemove))
	for cat, entries := range toRemove {
		s := make(map[string]bool, len(entries))
		for _, e := range entries {
			s[e] = true
		}
		removeSets[cat] = s
	}

	// First pass: count surviving entries per category, so we know which
	// category headers (and the whole staging section) to drop.
	survivors := make(map[string]int)
	totalSurvivors := 0
	var cat string
	for i := startIdx + 1; i < endIdx; i++ {
		line := lines[i]
		if m := categoryLineRe.FindStringSubmatch(line); m != nil {
			cat = reverse[strings.ToLower(strings.TrimSpace(m[1]))]
			if _, ok := survivors[cat]; !ok {
				survivors[cat] = 0
			}
			continue
		}
		if cat != "" && strings.HasPrefix(line, "- ") {
			entry := strings.TrimPrefix(line, "- ")
			if removeSets[cat] != nil && removeSets[cat][entry] {
				continue
			}
			survivors[cat]++
			totalSurvivors++
		}
	}

	if totalSurvivors == 0 {
		// Nothing left in staging — drop the entire section.
		return RemoveStaging(cl)
	}

	// Second pass: emit lines, skipping removed entries and any category whose
	// entries are all being removed. This preserves whatever blank-line layout
	// the source already had.
	var out []string
	out = append(out, lines[:startIdx]...)
	out = append(out, lines[startIdx]) // "## [staging]" header
	cat = ""
	skipping := false
	for i := startIdx + 1; i < endIdx; i++ {
		line := lines[i]
		if m := categoryLineRe.FindStringSubmatch(line); m != nil {
			cat = reverse[strings.ToLower(strings.TrimSpace(m[1]))]
			if survivors[cat] == 0 {
				skipping = true
				continue
			}
			skipping = false
			out = append(out, line)
			continue
		}
		if skipping {
			// Drop entries and any blank/other lines belonging to the removed category.
			continue
		}
		if cat != "" && strings.HasPrefix(line, "- ") {
			entry := strings.TrimPrefix(line, "- ")
			if removeSets[cat] != nil && removeSets[cat][entry] {
				continue
			}
		}
		out = append(out, line)
	}
	out = append(out, lines[endIdx:]...)

	newEntries := strings.TrimSpace(strings.Join(out, "\n"))
	if newEntries != "" {
		newEntries += "\n"
	}
	return buildContent(cl.Header, newEntries, cl.Footer)
}

// RemoveStaging removes the ## [staging] section from the changelog
// and returns the updated full content string. Does not modify the Changelog struct.
func RemoveStaging(cl *Changelog) string {
	if cl.Entries == "" {
		return buildContent(cl.Header, "", cl.Footer)
	}

	lines := strings.Split(cl.Entries, "\n")
	startIdx := -1
	for i, line := range lines {
		if stagingLineRe.MatchString(line) {
			startIdx = i
			break
		}
	}
	if startIdx == -1 {
		return buildContent(cl.Header, cl.Entries, cl.Footer)
	}

	endIdx := len(lines)
	for i := startIdx + 1; i < len(lines); i++ {
		if versionLineRe.MatchString(lines[i]) {
			endIdx = i
			break
		}
	}

	var newLines []string
	newLines = append(newLines, lines[:startIdx]...)
	newLines = append(newLines, lines[endIdx:]...)

	newEntries := strings.TrimSpace(strings.Join(newLines, "\n"))
	if newEntries != "" {
		newEntries += "\n"
	}

	return buildContent(cl.Header, newEntries, cl.Footer)
}

func buildContent(header, entries, footer string) string {
	var b strings.Builder

	h := strings.TrimRight(header, "\n")
	b.WriteString(h)
	b.WriteString("\n")

	if entries != "" {
		b.WriteString("\n")
		e := strings.TrimRight(entries, "\n")
		b.WriteString(e)
		b.WriteString("\n")
	}

	if footer != "" {
		b.WriteString("\n")
		b.WriteString(footer)
	}

	return b.String()
}
