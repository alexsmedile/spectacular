# AGENTS.md — Contributor & Agent Guide

Onboarding for any agent or human contributor working on the **Spectacular** repo
(the framework itself). This is the neutral, model-agnostic guide; `CLAUDE.md`
carries the same durable knowledge plus Claude-specific project state (active
requests, the skill's reference-doc index). Read this first.

## What This Is

**Spectacular** is an AI-native operational workspace framework. It ships as three layers:

1. **Convention** — a documented `.spectacular/` directory structure and file contract
2. **Skill** — a slash command (`/spectacular`) that operates the workspace
3. **CLI** — a one-time bootstrap tool (`spectacular init`) that scaffolds the directory

This repo *uses* the Spectacular skill to manage itself — the `.spectacular/`
directory here is the live workspace for this project, on the OKF v2.0 layout.

## Repository Structure

```
spectacular/
├── .spectacular/              # Live workspace (this project uses its own system)
│   ├── config.yaml            # Project config, naming rules, agents.file + tool_overrides
│   ├── PRD.md                 # Product intent — what & why & for whom
│   ├── PRINCIPLES.md          # Operating principles + runtime enforcement hooks
│   ├── ARCHITECTURE.md        # .spectacular/ structure, frontmatter, lifecycle, versioning
│   ├── AGENTS.md              # Operating rules for agents working IN .spectacular/
│   ├── STACK.md               # Host project's tech choices
│   ├── POLICY.md              # practice layer / merged policy contract
│   ├── specs/                 # Per-capability specs + index.md system spec
│   ├── roadmaps/              # roadmaps/index.md + shipped v*.md files
│   ├── decisions/             # ADR decision index + D*.md files
│   ├── memories/              # Durable facts memory index + M*.md files
│   ├── sessions/              # work time-log sessions/index.md + S*.md files
│   ├── requests/              # Active and planned work
│   └── archive/               # Completed requests (never deleted)
├── .claude-plugin/            # Claude Code plugin manifest + marketplace.json
├── .codex-plugin/             # Codex plugin manifest
├── cli/spectacular            # Bash binary (targets bash 3.2 — see Conventions)
├── cli/install.sh             # curl-installable installer → ~/.local/bin
├── skills/spectacular/        # The skill itself
│   ├── SKILL.md               # Lean orchestrator — triggers + routing table
│   ├── references/            # Reference docs loaded on demand by the skill
│   ├── templates/             # Canonical templates (PRD kits, base)
│   └── versions/              # Historical SKILL.md snapshots — never overwrite
├── agents/                    # ★ Subagent defs — SOURCE OF TRUTH (edit here, not .claude/agents/)
│                              #   fleet = discover/apply/review × fix/build + specialists:
│                              #   debug-{investigator,fixer,researcher} · repo-explorer · spec-builder
│                              #   · code-reviewer · test-verifier
│                              #   .claude/agents/*.md are relative symlinks → ../../agents/<name>.md
├── packs/                     # App-store convention packs (alex-default)
├── hooks/                     # ⚠ Claude/Codex PLUGIN event handlers (plugin runtime)
├── scripts/hooks/pre-commit   # ⚠ GIT hook — version-consistency guard
├── tests/                     # Bash test suite (cli/ · pipeline/ · agents/)
├── docs/                      # User-facing docs (commands, configuration, install…)
├── CLAUDE.md                  # Claude-specific twin of this guide + project state
└── AGENTS.md                  # This file
```

> **Two `hooks/` directories, same word, different runtimes.** `hooks/` at repo
> root holds **Claude/Codex plugin hooks** (`hooks.json`, `hooks-codex.json`),
> read by the plugin runtime *after* a user installs spectacular as a plugin.
> `scripts/hooks/` holds **git hooks** (`pre-commit`) that run locally during
> development of *this* repo. Don't conflate.

## Skill Architecture

The skill (`skills/spectacular/SKILL.md`) is a **lean orchestrator**: it reads
triggers, routes to a reference doc, and that doc contains the actual
instructions. This mirrors Spectacular's own philosophy — small files,
progressive context loading. Reference docs in `skills/spectacular/references/`
are loaded *on demand*, not up front. (The full per-doc "loaded when" index lives
in `CLAUDE.md` — it changes per release, so it's kept in one place.)

## Dev Commands

No build step. Baseline checks during development:

```bash
bash -n cli/spectacular cli/install.sh scripts/hooks/pre-commit   # syntax
scripts/hooks/pre-commit --check                                   # version guard
./cli/spectacular --help                                           # arg parsing loads
bash tests/run.sh                                                   # full suite
bash tests/cli/doctor.test.sh                                       # one area
```

`bash -n` + the version guard are the required baseline. For a one-off manual
check of a command path, use a throwaway dir — never commit scratch:

```bash
tmpdir="$(mktemp -d)" && cd "$tmpdir" && /path/to/repo/cli/spectacular init --name demo
```

## Conventions

**Coding style**

- **Bash 3.2 target.** macOS ships bash 3.2.57 — **no associative arrays**
  (`declare -A` / `local -A` both fail). Use plain-string membership, not
  assoc-array sets. Load-bearing: assoc-array code passes on Linux and breaks
  for every Mac user.
- Scripts use `set -euo pipefail`, two-space indent, lowercase helper names,
  uppercase constants (`GITHUB_REPO`). Prefer small helpers over inline blocks.
- Request slugs, filenames, generated paths → kebab-case.
- Markdown: concise headings, preserve existing frontmatter, relative paths.
- Add new CLI-behavior scenarios to `tests/cli/*.test.sh` in the existing
  assert style (assert on findings + exit codes).

**How the workspace system works** (durable — the same rules the skill enforces)

- **Frontmatter is the signal layer.** The skill reads frontmatter, not full
  file content, during briefings. To discover the active-request fleet without
  opening files, run **`spectacular status --json`** — it emits one object per
  request (frontmatter + grep-safe body signals: goal, `x/total` progress,
  current milestone). This is the agent opt-in contract; prefer it over
  hand-parsing `requests/*/PLAN.md`.
- **PLAN/TASKS have an enforced structure.** PLAN uses unnumbered canonical
  headings in order (`## Goal / ## Constraints / ## Milestones / ## Tasks /
  ## Dependencies / ## Validation / ## Deliverables`; extra sections allowed
  between). TASKS groups by `### M<N>` with flush-left `- [ ]`/`- [x]`/`- [~]`
  checkboxes. `doctor` errors on drift for active requests (archive skipped);
  `doctor --fix` de-numbers legacy headings.
- **Lifecycle state lives in `PLAN.md` frontmatter** (`status: planned | active
  | review | verified`). `TASKS.md` carries a mirror for skim tooling; PLAN is
  authoritative and `doctor` repairs drift.
- **ReadOnly docs are never overwritten in place** — snapshot first (`PRD@v1.0.md`
  naming); the unversioned filename always points to current. Same rule for
  `skills/spectacular/versions/`: add a new snapshot when releasing a version.
- **Memory (`spectacular remember this`) writes to `.spectacular/memories/`** —
  git-committed, team-visible. Never to `.claude/` personal memory.
- **`.spectacular/` is fully committed to git; `.spectacular.local/` is always
  gitignored.**
- **Convention packs (v0.4.0+)** are opt-in repo-shape opinions declared via
  `config.yaml`'s `convention_pack:` block — four scope locations (project-local
  → user → app-store → bundled) and three modes (`suggest`/`scaffold`/`enforce`).
  Schema: `skills/spectacular/references/packs-contract.md`.
- **Two-layer task tracking:** harness `TaskCreate`/`TaskUpdate` = ephemeral
  session micro-tracker; on-disk `requests/<slug>/TASKS.md` = persistent
  milestone blocks. Anti-pattern: one-for-one duplication. Full convention in
  `.spectacular/AGENTS.md` § Task tracking.

## Before Editing Canonical Docs

When changing Spectacular's canonical docs or skill behavior, read current
project intent first — don't edit from assumption:

- `.spectacular/AGENTS.md` — operating rules + authoritative context-loading table
- `.spectacular/PRD.md` — product intent
- `.spectacular/specs/index.md` — one-page index of what's built now

When working a request in `.spectacular/requests/<slug>/`, load that folder's
`PLAN.md` + `TASKS.md` and only the `specs/<capability>.md` it references — not
the whole `specs/` tree. Keep `.spectacular.local/` personal and uncommitted.

## CLI

`cli/spectacular` — Bash binary, installed to `~/.local/bin/spectacular` via `cli/install.sh`.

```
spectacular init                    # zero prompts, always-set (6 files) + blank kit
spectacular init -i                 # interactive: kit menu + per-doc prompts
spectacular init --kit coding       # adds STACK + ARCHITECTURE (coding kit triggers)
spectacular init --with principles,roadmap   # additive: extras on top of always-set
spectacular init --minimal          # always-set only; kit identity preserved
spectacular init --name my-app --agents-file CLAUDE.md
spectacular init --global           # install to ~/.agents/ and ~/.claude/
spectacular init --update           # re-download latest skill release

spectacular status                  # deterministic active-request fleet table
spectacular status <slug>           # single request card (goal · progress · deps · stale)
spectacular status --json           # machine-readable fleet (agent opt-in contract)

spectacular doctor                  # substrate self-check (all areas)
spectacular doctor <area>           # scoped: skill|workspace|frontmatter|snapshots|links|lifecycle|kits|conventions|specs|docs|personas|memory|sessions|feedback|ideas|debug|policies|vision|decisions|roadmap
spectacular doctor --fix            # apply mechanical fixes

spectacular policy                  # read the merged policy contract (POLICY.md + config overrides)
spectacular policy @<hook>          # one work-phase's policies + linked principle lines
spectacular policy --principle N    # reverse: which policies enforce principle N

spectacular pack list               # show packs across all 4 scopes
spectacular pack install <name>     # copy to ~/.spectacular/packs/<name>/
spectacular pack show <name>        # print scope + frontmatter
spectacular pack remove <name>      # delete user-scope pack
```

Skill is fetched from GitHub (tagged release tarball → `.agents/skills/spectacular/`);
`.claude/skills/spectacular/` is a symlink to that target. Version tracked in
`.spectacular/skills.lock`. Full command reference: `docs/commands.md`. Config
schema: `docs/configuration.md`.

## Commits

Conventional Commit prefixes (`feat:`, `fix:`, `chore:`). Imperative, specific
subjects (`fix: preserve skill install ref fallback`). PRs should state the
behavior change, list verification commands, and link the related
`.spectacular/requests/<slug>/`. **Versioned releases must bump every guarded
version string together** (plugin manifests, README badge, CHANGELOG, tag, skill
frontmatter) — the pre-commit guard blocks on drift.
