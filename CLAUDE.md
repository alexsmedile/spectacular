# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## What This Is

**Spectacular** is an AI-native operational workspace framework. It ships as three layers:

1. **Convention** — a documented `.spectacular/` directory structure and file contract
2. **Skill** — a Claude Code slash command (`/spectacular`) that operates the workspace
3. **CLI** — a one-time bootstrap tool (`spectacular init`) that scaffolds the directory

This repo *uses* the Spectacular skill to manage itself — the `.spectacular/` directory here is the live workspace for this project.

## Repository Structure

```
spectacular/
├── .spectacular/          # Live workspace (this project uses its own system)
│   ├── config.yaml        # Project config, naming rules, agents.file + tool_overrides
│   ├── AGENTS.md          # Context loading rules and available skills
│   ├── PRD.md             # Product intent
│   ├── STACK.md           # Technology choices and rules
│   ├── DECISIONS.md       # Architectural decisions log
│   └── requests/          # Active and planned work
│       └── cli-bootstrap/ # CLI tool (status: active)
├── cli/                   # CLI implementation
│   ├── spectacular        # Bash binary — spectacular init
│   └── install.sh         # curl-installable installer
└── skills/spectacular/    # The skill itself
    ├── SKILL.md           # Lean orchestrator — triggers + routing table
    └── references/        # Reference docs loaded on demand by the skill
```

## Active Requests

| Slug | Status | Summary |
|---|---|---|
| `cli-bootstrap` | active | `spectacular init` Bash CLI — built, pending GitHub publish for skill install |

## Skill Architecture

The skill (`skills/spectacular/SKILL.md`) is a **lean orchestrator**: it reads triggers, routes to a reference doc, and that doc contains the actual instructions. This mirrors Spectacular's own philosophy — small files, progressive context loading.

Reference docs in `skills/spectacular/references/` are loaded *on demand*:

| File | Loaded when |
|---|---|
| `status.md` | `/spectacular` with no args |
| `new-request.md` | `spectacular new <description>` or `spectacular promote` |
| `active-request.md` | Actively working on a request |
| `lifecycle.md` | Lifecycle transitions |
| `archive.md` | `spectacular archive <slug>` |
| `memory.md` | `spectacular remember this` |
| `versioning.md` | `spectacular snapshot <file>` |
| `init-workflow.md` | `spectacular init` (CLI context) |
| `onboarding.md` | First invocation on existing `.spectacular/` project |
| `scaffold-reference.md` | File template reference — frontmatter stubs for all file types |

## Key Conventions

**Frontmatter is the signal layer** — the skill reads frontmatter, not full file content, during briefings.

**Lifecycle state lives in `PLAN.md` frontmatter only** (`status: planned | active | review | verified`). Never duplicated.

**Canonical documents are never overwritten in place** — always snapshot first (`PRD@v1.0.md` naming). The unversioned filename always points to current.

**Memory (`spectacular remember this`) writes to `.spectacular/memory/`** — git-committed, team-visible. Never to `.claude/` personal memory.

**`.spectacular/`** is fully committed to git. **`.spectacular.local/`** is always gitignored.

## Context Loading Rules (from AGENTS.md)

| Task type | Load |
|---|---|
| Planning / design | PRD.md, DECISIONS.md |
| Skill implementation | STACK.md, PLAN.md, TASKS.md, current/skill/ |
| CLI implementation | STACK.md, PLAN.md, TASKS.md, current/cli/ |
| Review / QA | VERIFY.md, current/<capability>, RISKS.md |

Load only the capability spec relevant to the current task, not all of `current/`. Never read `archive/` during normal operation.

## CLI

`cli/spectacular` — Bash binary, installed to `~/.local/bin/spectacular` via `cli/install.sh`.

```
spectacular init                    # zero prompts, defaults
spectacular init -i                 # interactive mode
spectacular init --name my-app --agents-file CLAUDE.md
spectacular init --global           # install to ~/.agents/ and ~/.claude/
spectacular init --update           # re-download latest skill release
```

Skill is fetched from GitHub (tagged release tarball → `.agents/skills/spectacular/`). `.claude/skills/spectacular/` is a symlink to the `.agents/` target. Version tracked in `.spectacular/skills.lock`.

See `.spectacular/requests/cli-bootstrap/` for full scope and task checklist.
