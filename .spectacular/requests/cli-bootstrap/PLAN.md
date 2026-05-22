---
status: planned
priority: medium
owner: alex
updated: 2026-05-22
summary: "Bash CLI that bootstraps a new project with the Spectacular scaffold and installs the skill from GitHub (parked: smart-init shipped at v0.3.0; this stays for v0.2.x maintenance)"
related:
  - ../../PRD.md
  - ../../ARCHITECTURE.md
  - ../../archive/canonical-docs-rework/PLAN.md
---

# CLI Bootstrap Tool

## Goal

Build `spectacular init` — a zero-dependency Bash CLI that scaffolds the `.spectacular/` directory structure, installs the skill from GitHub, and writes stub root files. Fast by default, configurable via flags.

## Why

Without a CLI, onboarding a new project requires manual directory creation and file copying. The CLI removes that friction and enforces the correct minimal scaffold.

## Approach

Single Bash script. Distributed via curl-install from the GitHub repo. Installed to `~/.local/bin/spectacular`.

```
cli/spectacular     # the binary (no extension)
cli/install.sh      # curl-installable installer
```

## CLI interface

```
spectacular init                    # fast, zero prompts, all defaults
spectacular init -i                 # interactive — prompts for all foundational settings
spectacular init --name my-app      # override project name
spectacular init --summary "..."    # override project summary
spectacular init --agents-file CLAUDE.md   # override default agents file (default: AGENTS.md)
spectacular init --global           # install skill to ~/.agents/skills/ + ~/.claude/skills/
spectacular init --update           # re-download latest skill release, update skills.lock
spectacular init --help
```

Bare `init` infers project name from the current folder slug. No prompts. No summary. Defaults to `AGENTS.md`.

## Skill installation

Skill source is always pulled from GitHub (tagged release tarball, falls back to `main`).

| Layer | Default (project-local) | `--global` |
|---|---|---|
| Source of truth | `.agents/skills/spectacular/` | `~/.agents/skills/spectacular/` |
| Claude symlink | `.claude/skills/spectacular/` → above | `~/.claude/skills/spectacular/` → above |

Version is recorded in `.spectacular/skills.lock` (CLI-written):

```yaml
spectacular:
  ref: v1.0.0
  sha: abc123def456
  installed: 2026-05-11
  source: https://github.com/alexsmedile/spectacular/archive/refs/tags/v1.0.0.tar.gz
```

## Idempotency

- Existing scaffold files: skip, report skipped
- Existing skill install: skip download, report installed ref from `skills.lock`
- `--update`: re-download latest tag, overwrite skill dirs, update `skills.lock`
- Network failure: fail loudly, exit non-zero — partial scaffold (files) is kept, skill install step fails cleanly

## Success output

```
Spectacular initialized.

Created:
  .spectacular/config.yaml
  .spectacular/PRD.md
  .spectacular/PRINCIPLES.md
  .spectacular/ARCHITECTURE.md
  .spectacular/ROADMAP.md
  .spectacular/STACK.md
  .spectacular/DECISIONS.md
  .spectacular/AGENTS.md
  .spectacular/current/
  .spectacular/requests/

Skill installed:
  .agents/skills/spectacular/     (v1.0.0)
  .claude/skills/spectacular/     → symlink

Gitignore:
  .spectacular.local/ added to .gitignore

Run /spectacular to get started.
```

## Scope

- `spectacular init` with flags above
- `skills.lock` written by CLI
- Multiplatform: `.agents/skills/` + `.claude/skills/` both populated
- **Full canonical doc set on every init** — PRD, PRINCIPLES, ARCHITECTURE, ROADMAP, STACK, DECISIONS, AGENTS. Aligned with the post-rework root layer ([`canonical-docs-rework`](../canonical-docs-rework/PLAN.md)).
- PRD stub uses the 6-slot shape (problem / who / success / non-goals / constraints / milestone)
- AGENTS.md stub uses the new onboarding shape (folder purpose, context loading, don'ts)

## Out of scope

- `spectacular` subcommands beyond `init` (skill-level operations)
- `--force` flag — v2
- Workspace support — v2
- Nested workspace support — v2
- Package registry / homebrew publishing — separate concern

## Success criteria

- Bare `spectacular init` on a blank directory produces correct scaffold with zero prompts
- `-i` mode prompts for name, summary, agents-file, global scope
- Skill installed from GitHub tarball, both `.agents/` and `.claude/` targets populated
- `skills.lock` written with ref + sha
- Re-run skips existing files and reports clearly
- `--update` refreshes skill and updates `skills.lock`
- Network failure exits non-zero with actionable message
- `.spectacular.local/` is gitignored
