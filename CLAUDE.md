# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## What is clog?

A changelog fragment manager written in Go. Developers create YAML fragment files (one per branch) in `changelog.d/`, then at release time `clog release` merges them into `CHANGELOG.md` and optionally commits + tags.

## Build & Test Commands

```bash
make build      # Build binary with version ldflags → ./clog
make test       # Run all tests: go test ./...
make install    # go install ./cmd/clog
make clean      # Remove built binary
```

Run a single test file or function:
```bash
go test ./internal/fragment/ -run TestParse
```

## Architecture

All application code lives under `internal/`:

- **`command/`** — Cobra CLI commands (`init`, `new`, `validate`, `preview`, `release`). The root command is in `root.go`; each subcommand in its own file. Entry point is `cmd/clog/main.go` → `command.Execute()`.
- **`fragment/`** — YAML fragment parsing, validation, and template generation. Fragments are `map[string][]string` keyed by category. Eight categories in fixed order: deployment, added, changed, deprecated, removed, fixed, security, yanked.
- **`changelog/`** — Parses `CHANGELOG.md` into header/entries/footer sections. Handles version extraction and inserting new release entries while preserving the footer (category definitions under `# Notes`).
- **`merge/`** — Combines multiple fragments into a single category→entries map, then renders to markdown. Rendering uses the fixed category display order from `fragment.CategoryOrder`.
- **`gitutil/`** — Git operations via shell exec: branch name retrieval, branch-to-filename sanitization (`feature/foo` → `feature-foo.yaml`), release commits, and tagging.

## Key Design Decisions

- Version is injected at build time via `-ldflags` into `command.Version`.
- Fragment filenames are derived from git branch names (sanitized: lowercased, `/` → `-`).
- The `release` command is interactive — it prompts for version, optional metadata, confirmation, and whether to auto-commit/tag.
- `CHANGELOG.md` footer (everything from `# Notes` onward) is preserved across releases.
- Release metadata format supports parenthetical annotations like `(86.8%)(Dev)(Prod)` appended after the date.
