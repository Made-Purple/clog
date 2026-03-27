package gitutil

import "testing"

func TestSanitizeBranchName(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"feature/knox-plugin", "feature-knox-plugin"},
		{"bugfix/fix-login-crash", "bugfix-fix-login-crash"},
		{"hotfix/3.81.1", "hotfix-3.81.1"},
		{"main", "main"},
		{"Feature/UPPER-Case", "feature-upper-case"},
		{"a/b/c/d", "a-b-c-d"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := SanitizeBranchName(tt.input)
			if got != tt.want {
				t.Errorf("SanitizeBranchName(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}
