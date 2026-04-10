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

	type catBlock struct {
		headerLine string
		entries    []string
	}
	var blocks []catBlock
	var currentCat string
	var currentHeader string
	var currentEntries []string
	flush := func() {
		if currentHeader == "" {
			return
		}
		blocks = append(blocks, catBlock{headerLine: currentHeader, entries: currentEntries})
		currentHeader = ""
		currentEntries = nil
		currentCat = ""
	}

	for i := startIdx + 1; i < endIdx; i++ {
		line := lines[i]
		if m := categoryLineRe.FindStringSubmatch(line); m != nil {
			flush()
			catDisplay := strings.TrimSpace(m[1])
			currentCat = reverse[strings.ToLower(catDisplay)]
			currentHeader = line
			continue
		}
		if currentCat != "" && strings.HasPrefix(line, "- ") {
			entry := strings.TrimPrefix(line, "- ")
			if removeSets[currentCat] != nil && removeSets[currentCat][entry] {
				continue
			}
			currentEntries = append(currentEntries, line)
			continue
		}
		// Blank lines or other content within the current category — keep only
		// if we're currently inside a category block (they'll be re-emitted).
		if currentCat != "" {
			currentEntries = append(currentEntries, line)
		}
	}
	flush()

	// Trim trailing blank lines from each block and drop blocks with no entries.
	var kept []catBlock
	for _, b := range blocks {
		hasEntry := false
		for _, e := range b.entries {
			if strings.HasPrefix(e, "- ") {
				hasEntry = true
				break
			}
		}
		if !hasEntry {
			continue
		}
		// Drop trailing blank lines.
		for len(b.entries) > 0 && strings.TrimSpace(b.entries[len(b.entries)-1]) == "" {
			b.entries = b.entries[:len(b.entries)-1]
		}
		kept = append(kept, b)
	}

	var newLines []string
	newLines = append(newLines, lines[:startIdx]...)
	if len(kept) > 0 {
		newLines = append(newLines, lines[startIdx]) // "## [staging]" header
		for i, b := range kept {
			if i > 0 {
				newLines = append(newLines, "")
			}
			newLines = append(newLines, b.headerLine)
			newLines = append(newLines, b.entries...)
		}
		newLines = append(newLines, "")
	}
	newLines = append(newLines, lines[endIdx:]...)

	newEntries := strings.TrimSpace(strings.Join(newLines, "\n"))
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
