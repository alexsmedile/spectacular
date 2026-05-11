---
status: planned
priority: medium
owner: alex
updated: 2026-05-11
summary: "CLI tool that bootstraps a new project with the Spectacular scaffold"
related:
  - current/cli
---

# CLI Bootstrap Tool

## Goal

Build `spectacular init` — a one-time CLI tool that scaffolds the `.spectacular/` directory structure on a new project, installs the skill, and writes stub root files.

## Why

Without a CLI, onboarding a new project requires manual directory creation and file copying. The CLI removes that friction and enforces the correct minimal scaffold (only `current/` and `requests/` by default, not all directories).

## Scope

- `spectacular init` command
- Interactive prompts for project name and summary
- Writes `config.yaml`
- Creates `current/` and `requests/` directories
- Creates stub root files (`PRD.md`, `STACK.md`, `DECISIONS.md`, `AGENTS.md`) with frontmatter
- Installs skill: symlinks `~/.claude/skills/spectacular/` into `.claude/skills/spectacular/` if global exists, otherwise copies
- Appends `.spectacular.local/` to `.gitignore` (creates if missing)
- Idempotent: safe to re-run, skips existing files, reports what was skipped

## Out of scope

- `spectacular` subcommands beyond `init` (those are skill-level operations)
- Workspace support (`.spectacular.<name>/`) — v2
- Nested workspace support — v2
- Auto-updating the skill (separate concern)
- Package registry publishing (separate concern)

## Approach

Simple CLI binary. Stack TBD — candidates are Node.js (broad ecosystem, easy publishing to npm) or a shell script (zero dependencies, simple distribution). Decision needed before implementation.

Key constraint: the CLI should have minimal or zero dependencies so it's trivially installable.

## Success criteria

- `spectacular init` runs on a blank directory and produces the correct scaffold
- Re-running on an existing project skips existing files and reports clearly
- Skill is accessible as `/spectacular` after init completes
- `.spectacular.local/` is gitignored
- Output is friendly and tells the user what was created
