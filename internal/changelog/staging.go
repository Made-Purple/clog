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
