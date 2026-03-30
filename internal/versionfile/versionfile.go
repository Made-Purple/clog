package versionfile

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

var packageVersionRe = regexp.MustCompile(`("version"\s*:\s*)"[^"]*"`)

// UpdateVersionFile writes the version (without leading "v") followed by a newline to path.
func UpdateVersionFile(path string, version string) error {
	content := strings.TrimPrefix(version, "v") + "\n"
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		return fmt.Errorf("writing VERSION file: %w", err)
	}
	return nil
}

// UpdatePackageJSON reads the file at path, replaces the first "version" field value
// with the given version (without leading "v"), and writes back preserving all formatting.
func UpdatePackageJSON(path string, version string) error {
	version = strings.TrimPrefix(version, "v")

	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("reading package.json: %w", err)
	}

	loc := packageVersionRe.FindIndex(data)
	if loc == nil {
		return fmt.Errorf("no \"version\" field found in package.json")
	}

	match := data[loc[0]:loc[1]]
	replaced := packageVersionRe.ReplaceAll(match, []byte(fmt.Sprintf(`${1}"%s"`, version)))

	var buf []byte
	buf = append(buf, data[:loc[0]]...)
	buf = append(buf, replaced...)
	buf = append(buf, data[loc[1]:]...)

	if err := os.WriteFile(path, buf, 0644); err != nil {
		return fmt.Errorf("writing package.json: %w", err)
	}
	return nil
}
