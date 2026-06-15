// Package skill installs the clog assistant skill into the per-user or
// per-project skill directories of supported AI coding assistants.
//
// The SKILL.md format is shared across assistants (Claude Code, Codex), so the
// same embedded content is written to a parallel directory tree:
//
//	~/.claude/skills/clog/SKILL.md   (Claude, global)
//	.claude/skills/clog/SKILL.md     (Claude, project)
//	~/.codex/skills/clog/SKILL.md    (Codex, global)
//	.codex/skills/clog/SKILL.md      (Codex, project)
package skill

import (
	"bytes"
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
)

// content is the canonical SKILL.md, compiled into the binary so the installer
// has it without a separate download.
//
//go:embed clog/SKILL.md
var content []byte

// Name is the skill's folder name under each assistant's skills/ directory.
const Name = "clog"

// Agent identifies a supported AI coding assistant.
type Agent struct {
	Key     string // stable identifier, e.g. "claude"
	Display string // human-facing name, e.g. "Claude"
	Dir     string // config dir name, e.g. ".claude"
}

// Supported agents, in display order.
var (
	Claude = Agent{Key: "claude", Display: "Claude", Dir: ".claude"}
	Codex  = Agent{Key: "codex", Display: "Codex", Dir: ".codex"}
	Agents = []Agent{Claude, Codex}
)

// Scope is where a skill is installed.
type Scope int

const (
	// Global installs into the user's home config dir (e.g. ~/.claude).
	Global Scope = iota
	// Project installs into the current repository (e.g. .claude).
	Project
)

// HomeConfigDir returns the agent's per-user config directory, e.g. ~/.claude.
func (a Agent) HomeConfigDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, a.Dir), nil
}

// Detected reports whether the agent's per-user config directory exists,
// indicating the user has that assistant set up.
func (a Agent) Detected() bool {
	dir, err := a.HomeConfigDir()
	if err != nil {
		return false
	}
	info, err := os.Stat(dir)
	return err == nil && info.IsDir()
}

// TargetPath returns the SKILL.md path for the agent at the given scope.
// Project paths are relative to the current working directory.
func (a Agent) TargetPath(scope Scope) (string, error) {
	switch scope {
	case Global:
		cfg, err := a.HomeConfigDir()
		if err != nil {
			return "", err
		}
		return filepath.Join(cfg, "skills", Name, "SKILL.md"), nil
	case Project:
		return filepath.Join(a.Dir, "skills", Name, "SKILL.md"), nil
	default:
		return "", fmt.Errorf("unknown scope")
	}
}

// Result reports the outcome of installing a skill.
type Result struct {
	Path    string
	Updated bool // true if written/changed; false if already up to date
}

// Install writes the embedded SKILL.md for the agent at the given scope,
// creating parent directories as needed. If the target already holds identical
// content, it is left untouched and Updated is false.
func (a Agent) Install(scope Scope) (Result, error) {
	path, err := a.TargetPath(scope)
	if err != nil {
		return Result{}, err
	}
	if existing, err := os.ReadFile(path); err == nil && bytes.Equal(existing, content) {
		return Result{Path: path, Updated: false}, nil
	}
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return Result{}, fmt.Errorf("creating %s: %w", filepath.Dir(path), err)
	}
	if err := os.WriteFile(path, content, 0644); err != nil {
		return Result{}, fmt.Errorf("writing %s: %w", path, err)
	}
	return Result{Path: path, Updated: true}, nil
}

// Installed reports whether a clog skill file exists for the agent at the scope.
func (a Agent) Installed(scope Scope) bool {
	path, err := a.TargetPath(scope)
	if err != nil {
		return false
	}
	_, err = os.Stat(path)
	return err == nil
}

// UninstallResult reports the outcome of removing a skill.
type UninstallResult struct {
	Path       string
	Existed    bool // a skill file was present before
	Removed    bool // the file was deleted
	Customized bool // existing content differed from the embedded SKILL.md
}

// Uninstall removes the clog skill for the agent at the given scope. If the
// installed SKILL.md differs from the embedded version (i.e. it was hand-edited),
// it is left in place unless force is true — the returned result reports
// Existed=true, Customized=true, Removed=false so the caller can warn. After a
// successful removal, an empty skills/clog directory is pruned (skills/ is left
// alone).
func (a Agent) Uninstall(scope Scope, force bool) (UninstallResult, error) {
	path, err := a.TargetPath(scope)
	if err != nil {
		return UninstallResult{}, err
	}
	existing, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return UninstallResult{Path: path}, nil
		}
		return UninstallResult{}, fmt.Errorf("reading %s: %w", path, err)
	}
	customized := !bytes.Equal(existing, content)
	if customized && !force {
		return UninstallResult{Path: path, Existed: true, Customized: true}, nil
	}
	if err := os.Remove(path); err != nil {
		return UninstallResult{}, fmt.Errorf("removing %s: %w", path, err)
	}
	// Prune the now-empty skills/clog directory; os.Remove is a no-op if it
	// still holds other files (e.g. references/), and we never touch skills/.
	os.Remove(filepath.Dir(path))
	return UninstallResult{Path: path, Existed: true, Removed: true, Customized: customized}, nil
}

// AgentByKey returns the supported agent with the given key.
func AgentByKey(key string) (Agent, bool) {
	for _, a := range Agents {
		if a.Key == key {
			return a, true
		}
	}
	return Agent{}, false
}
