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

	skillUninstallCmd.Flags().Bool("claude", false, "Remove the skill from Claude")
	skillUninstallCmd.Flags().Bool("codex", false, "Remove the skill from Codex")
	skillUninstallCmd.Flags().Bool("global", false, "Remove from your home config dir (~/.claude, ~/.codex)")
	skillUninstallCmd.Flags().Bool("project", false, "Remove from this project (.claude, .codex)")
	skillUninstallCmd.Flags().Bool("force", false, "Remove even if the skill was modified from the installed version")

	skillCmd.AddCommand(skillInstallCmd)
	skillCmd.AddCommand(skillUninstallCmd)
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
		selected, err := promptAgents(reader,
			"Install the clog skill for which assistant(s)?", "detected", skill.Agent.Detected)
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
		scope, err := promptScope(reader, "Where should the skill be installed?")
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

var skillUninstallCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "Remove the clog skill from Claude and/or Codex",
	Long: `Remove the clog skill from the skill directory of one or more AI coding
assistants.

Run with no flags to choose interactively. Use flags to skip the prompts:

  clog skill uninstall --claude --global     # Claude, ~/.claude/skills/clog
  clog skill uninstall --codex --project     # Codex, .codex/skills/clog

A skill that has been modified from the installed version is left in place; pass
--force to remove it anyway.`,
	// A "modified, use --force" stop is a runtime condition, not a usage error,
	// so don't dump the usage text on it.
	SilenceUsage: true,
	RunE:         runSkillUninstall,
}

func runSkillUninstall(cmd *cobra.Command, args []string) error {
	claudeFlag, _ := cmd.Flags().GetBool("claude")
	codexFlag, _ := cmd.Flags().GetBool("codex")
	globalFlag, _ := cmd.Flags().GetBool("global")
	projectFlag, _ := cmd.Flags().GetBool("project")
	force, _ := cmd.Flags().GetBool("force")

	reader := bufio.NewReader(os.Stdin)

	// Resolve which agents to remove from: flags win, otherwise prompt.
	var agents []skill.Agent
	if claudeFlag || codexFlag {
		if claudeFlag {
			agents = append(agents, skill.Claude)
		}
		if codexFlag {
			agents = append(agents, skill.Codex)
		}
	} else {
		selected, err := promptAgents(reader,
			"Remove the clog skill from which assistant(s)?", "installed", installedAnywhere)
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
		scope, err := promptScope(reader, "Remove the skill from where?")
		if err != nil {
			return err
		}
		scopes = []skill.Scope{scope}
	}

	// Remove each agent at each scope.
	modified := 0
	for _, s := range scopes {
		for _, a := range agents {
			res, err := a.Uninstall(s, force)
			if err != nil {
				return err
			}
			switch {
			case res.Removed:
				color.Success("%s (%s): removed %s", a.Display, scopeLabel(s), res.Path)
			case res.Existed && res.Customized:
				fmt.Printf("%s %s (%s): kept %s — modified from the installed version\n",
					color.Yellow("!"), a.Display, scopeLabel(s), res.Path)
				modified++
			default:
				fmt.Printf("%s %s (%s): not installed\n", color.Dim("•"), a.Display, scopeLabel(s))
			}
		}
	}
	if modified > 0 {
		return fmt.Errorf("left %d modified skill(s) in place; re-run with --force to remove", modified)
	}
	return nil
}

// installedAnywhere reports whether the agent has a clog skill in either scope.
func installedAnywhere(a skill.Agent) bool {
	return a.Installed(skill.Global) || a.Installed(skill.Project)
}

// promptAgents asks which assistants to act on. Agents for which isDefault
// returns true are pre-selected (chosen by pressing Enter) and flagged with
// hintLabel, e.g. "detected" for install or "installed" for uninstall.
func promptAgents(reader *bufio.Reader, header, hintLabel string, isDefault func(skill.Agent) bool) ([]skill.Agent, error) {
	if !isInteractive() {
		return nil, fmt.Errorf("no terminal available for prompts; pass --claude and/or --codex")
	}

	fmt.Println(header)
	for i, a := range skill.Agents {
		mark := " "
		hint := ""
		if isDefault(a) {
			mark = "x"
			hint = color.Dim(" (" + hintLabel + ")")
		}
		fmt.Printf("  %d) [%s] %s%s\n", i+1, mark, a.Display, hint)
	}
	color.Prompt(fmt.Sprintf("Enter numbers (comma-separated), or press Enter for %s:", hintLabel))

	line, err := reader.ReadString('\n')
	if err != nil && !(errors.Is(err, io.EOF) && line != "") {
		if errors.Is(err, io.EOF) {
			return nil, fmt.Errorf("no terminal available for prompts; pass --claude and/or --codex")
		}
		return nil, err
	}
	line = strings.TrimSpace(line)

	if line == "" {
		var defaults []skill.Agent
		for _, a := range skill.Agents {
			if isDefault(a) {
				defaults = append(defaults, a)
			}
		}
		return defaults, nil
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

// promptScope asks which scope to act on: home config (global) or this project.
func promptScope(reader *bufio.Reader, header string) (skill.Scope, error) {
	if !isInteractive() {
		return 0, fmt.Errorf("no terminal available for prompts; pass --global or --project")
	}

	fmt.Println(header)
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

// isInteractive reports whether stdin is a terminal we can prompt on. It is a
// variable rather than a plain function so tests can simulate an interactive
// terminal without needing a real PTY.
var isInteractive = func() bool {
	info, err := os.Stdin.Stat()
	if err != nil {
		return false
	}
	return info.Mode()&os.ModeCharDevice != 0
}
