package changelog

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

var versionLineRe = regexp.MustCompile(`^## \[`)
var versionExtractRe = regexp.MustCompile(`^## \[([^\]]+)\]`)
var footerRe = regexp.MustCompile(`(?m)^# Notes`)

// Changelog represents the parsed sections of a CHANGELOG.md file.
type Changelog struct {
	Header  string // Everything before the first ## [ line
	Entries string // All ## [x.y.z] sections
	Footer  string // From # Notes to EOF (inclusive)
}

// Read parses a CHANGELOG.md file into its sections.
func Read(path string) (*Changelog, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading changelog: %w", err)
	}

	content := string(data)
	return ParseContent(content), nil
}

// ParseContent splits changelog content into Header, Entries, and Footer sections.
func ParseContent(content string) *Changelog {
	cl := &Changelog{}

	// Find the footer (# Notes section)
	footerIdx := footerRe.FindStringIndex(content)
	var body string
	if footerIdx != nil {
		cl.Footer = content[footerIdx[0]:]
		body = content[:footerIdx[0]]
	} else {
		cl.Footer = ""
		body = content
	}

	// Find the first version line in the body
	lines := strings.Split(body, "\n")
	firstVersionLine := -1
	for i, line := range lines {
		if versionLineRe.MatchString(line) {
			firstVersionLine = i
			break
		}
	}

	if firstVersionLine == -1 {
		// No version entries yet
		cl.Header = body
		cl.Entries = ""
	} else {
		cl.Header = strings.Join(lines[:firstVersionLine], "\n")
		cl.Entries = strings.Join(lines[firstVersionLine:], "\n")
	}

	return cl
}

// LastVersion extracts the most recent version string from the changelog entries.
// Returns empty string if no version is found.
func LastVersion(cl *Changelog) string {
	matches := versionExtractRe.FindStringSubmatch(cl.Entries)
	if len(matches) < 2 {
		return ""
	}
	return matches[1]
}

// Insert returns the full file content with newEntry inserted between Header and Entries.
func Insert(cl *Changelog, newEntry string) string {
	var b strings.Builder

	header := strings.TrimRight(cl.Header, "\n")
	b.WriteString(header)
	b.WriteString("\n\n")
	b.WriteString(newEntry)

	if cl.Entries != "" {
		entries := strings.TrimRight(cl.Entries, "\n")
		b.WriteString("\n")
		b.WriteString(entries)
		b.WriteString("\n")
	}

	if cl.Footer != "" {
		b.WriteString("\n")
		b.WriteString(cl.Footer)
	}

	return b.String()
}
