package merge

import (
	"fmt"
	"strings"

	"github.com/made-purple/clog/internal/fragment"
)

// Merge combines multiple fragments into a single map of category -> entries.
// Empty entries are filtered out. Categories are not deduplicated.
func Merge(fragments []*fragment.Fragment) map[string][]string {
	result := make(map[string][]string)
	for _, f := range fragments {
		nonEmpty := fragment.NonEmptyEntries(f)
		for cat, entries := range nonEmpty {
			result[cat] = append(result[cat], entries...)
		}
	}
	return result
}

// Render produces the markdown block for a version entry.
// metadata is optional free-text appended after the date (e.g. "(98%)(Dev)").
func Render(version, date, metadata string, merged map[string][]string) string {
	var b strings.Builder

	// Version header line
	b.WriteString(fmt.Sprintf("## [%s] - %s", version, date))
	if metadata != "" {
		b.WriteString(" ")
		b.WriteString(metadata)
	}
	b.WriteString("\n")

	// Categories in fixed order
	for _, cat := range fragment.CategoryOrder {
		entries, ok := merged[cat]
		if !ok || len(entries) == 0 {
			continue
		}
		b.WriteString(fmt.Sprintf("### %s\n", fragment.CategoryDisplay[cat]))
		for _, entry := range entries {
			b.WriteString(fmt.Sprintf("- %s\n", entry))
		}
	}

	return b.String()
}
