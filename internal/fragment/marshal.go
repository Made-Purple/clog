package fragment

import (
	"strings"

	"gopkg.in/yaml.v3"
)

// MarshalEntries produces YAML bytes from a category->entries map.
// Categories with no entries get a single empty string placeholder.
// Categories follow the fixed CategoryOrder.
func MarshalEntries(entries map[string][]string) []byte {
	var b strings.Builder
	for _, cat := range CategoryOrder {
		b.WriteString(cat + ":\n")
		items := entries[cat]
		if len(items) == 0 {
			b.WriteString("  - \"\"\n")
		} else {
			for _, item := range items {
				itemBytes, _ := yaml.Marshal(item)
				b.WriteString("  - " + strings.TrimSpace(string(itemBytes)) + "\n")
			}
		}
	}
	return []byte(b.String())
}
