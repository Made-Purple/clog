package command

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/made-purple/clog/internal/color"
	"github.com/made-purple/clog/internal/skill"
	"github.com/spf13/cobra"
)

func init() {
	skillInstallCmd.Flags().Bool("claude", false, "Install the skill for Claude")
	skillInstallCmd.Flags().Bool("codex", false, "Install the skill for Codex")
	skillInstallCmd.Flags().Bool("global", false, "Install into your home config dir (~/.claude, ~/.codex)")
	skillInstallCmd.Flags().Bool("project", false, "Install into this project (.claude, .codex)")
	skillCmd.AddCommand(skillInstallCmd)
	rootCmd.AddCommand(skillCmd)
}

var skillCmd = &cobra.Command{
	Use:   "skill",
	Short: "Manage the clog assistant skill (Claude / Codex)",
	Long:  "Install the clog skill so Claude Code or Codex can manage changelog fragments for you.",
}

var skillInstallCmd = &cobra.Command{
	Use:   "install",
	Short: "Install the clog skill for Claude and/or Codex",
	Long: `Install the clog skill into the skill directory of one or more AI coding
assistants. The same SKILL.md works in both Claude Code and Codex.

Run with no flags to choose interactively. Use flags to skip the prompts:

  clog skill install --claude --global      # Claude, ~/.claude/skills/clog
  clog skill install --codex --project      # Codex, .codex/skills/clog
  clog skill install --claude --codex --global

Existing skills are compared and only rewritten when the content differs.`,
	RunE: runSkillInstall,
}

func runSkillInstall(cmd *cobra.Command, args []string) error {
	claudeFlag, _ := cmd.Flags().GetBool("claude")
	codexFlag, _ := cmd.Flags().GetBool("codex")
	globalFlag, _ := cmd.Flags().GetBool("global")
	projectFlag, _ := cmd.Flags().GetBool("project")

	reader := bufio.NewReader(os.Stdin)

	// Resolve which agents to install for: flags win, otherwise prompt.
	var agents []skill.Agent
	if claudeFlag || codexFlag {
		if claudeFlag {
			agents = append(agents, skill.Claude)
		}
		if codexFlag {
			agents = append(agents, skill.Codex)
		}
	} else {
		selected, err := promptAgents(reader)
		if err != nil {
			return err
		}
		agents = selected
	}
	if len(agents) == 0 {
		fmt.Println("No assistants selected. Nothing to do.")
		return nil
	}

	// Resolve scope(s): flags win (and may combine), otherwise prompt for one.
	var scopes []skill.Scope
	if globalFlag || projectFlag {
		if globalFlag {
			scopes = append(scopes, skill.Global)
		}
		if projectFlag {
			scopes = append(scopes, skill.Project)
		}
	} else {
		scope, err := promptScope(reader)
		if err != nil {
			return err
		}
		scopes = []skill.Scope{scope}
	}

	// Install each agent at each scope.
	for _, s := range scopes {
		for _, a := range agents {
			res, err := a.Install(s)
			if err != nil {
				return err
			}
			if res.Updated {
				color.Success("%s (%s): installed %s", a.Display, scopeLabel(s), res.Path)
			} else {
				fmt.Printf("%s %s (%s): already up to date %s\n",
					color.Dim("•"), a.Display, scopeLabel(s), color.Dim(res.Path))
			}
		}
	}
	return nil
}

// promptAgents asks which assistants to install for. Detected assistants are
// pre-selected as the default (chosen by pressing Enter).
func promptAgents(reader *bufio.Reader) ([]skill.Agent, error) {
	if !isInteractive() {
		return nil, fmt.Errorf("no terminal available for prompts; pass --claude and/or --codex")
	}

	fmt.Println("Install the clog skill for which assistant(s)?")
	for i, a := range skill.Agents {
		mark := " "
		hint := ""
		if a.Detected() {
			mark = "x"
			hint = color.Dim(" (detected)")
		}
		fmt.Printf("  %d) [%s] %s%s\n", i+1, mark, a.Display, hint)
	}
	color.Prompt("Enter numbers (comma-separated), or press Enter for detected:")

	line, err := reader.ReadString('\n')
	if err != nil && !(errors.Is(err, io.EOF) && line != "") {
		if errors.Is(err, io.EOF) {
			return nil, fmt.Errorf("no terminal available for prompts; pass --claude and/or --codex")
		}
		return nil, err
	}
	line = strings.TrimSpace(line)

	if line == "" {
		var detected []skill.Agent
		for _, a := range skill.Agents {
			if a.Detected() {
				detected = append(detected, a)
			}
		}
		return detected, nil
	}

	var chosen []skill.Agent
	for _, tok := range strings.Split(line, ",") {
		tok = strings.TrimSpace(strings.ToLower(tok))
		if tok == "" {
			continue
		}
		matched := false
		for i, a := range skill.Agents {
			if tok == fmt.Sprint(i+1) || tok == a.Key || tok == strings.ToLower(a.Display) {
				chosen = append(chosen, a)
				matched = true
				break
			}
		}
		if !matched {
			return nil, fmt.Errorf("invalid selection: %q", tok)
		}
	}
	return dedupeAgents(chosen), nil
}

// promptScope asks where to install: home config (global) or this project.
func promptScope(reader *bufio.Reader) (skill.Scope, error) {
	if !isInteractive() {
		return 0, fmt.Errorf("no terminal available for prompts; pass --global or --project")
	}

	fmt.Println("Where should the skill be installed?")
	fmt.Printf("  1) global   %s\n", color.Dim("your home config (~/.claude, ~/.codex)"))
	fmt.Printf("  2) project  %s\n", color.Dim("this repository (.claude, .codex)"))
	color.Prompt("Enter 1 or 2 (default 1):")

	line, err := reader.ReadString('\n')
	if err != nil && !(errors.Is(err, io.EOF) && line != "") {
		if errors.Is(err, io.EOF) {
			return 0, fmt.Errorf("no terminal available for prompts; pass --global or --project")
		}
		return 0, err
	}
	switch strings.TrimSpace(strings.ToLower(line)) {
	case "", "1", "g", "global":
		return skill.Global, nil
	case "2", "p", "project":
		return skill.Project, nil
	default:
		return 0, fmt.Errorf("invalid choice: %q", strings.TrimSpace(line))
	}
}

func scopeLabel(s skill.Scope) string {
	if s == skill.Project {
		return "project"
	}
	return "global"
}

func dedupeAgents(in []skill.Agent) []skill.Agent {
	seen := make(map[string]bool, len(in))
	var out []skill.Agent
	for _, a := range in {
		if !seen[a.Key] {
			seen[a.Key] = true
			out = append(out, a)
		}
	}
	return out
}

// isInteractive reports whether stdin is a terminal we can prompt on.
func isInteractive() bool {
	info, err := os.Stdin.Stat()
	if err != nil {
		return false
	}
	return info.Mode()&os.ModeCharDevice != 0
}
