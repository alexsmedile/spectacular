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
├── .spectacular/              # Live workspace (this project uses its own system)
│   ├── config.yaml            # Project config, naming rules, agents.file + tool_overrides
│   ├── PRD.md                 # Product intent — what & why & for whom
│   ├── SPEC.md                # System spec — index of what's built right now (v0.5.0+)
│   ├── PRINCIPLES.md          # Operating principles + runtime enforcement hooks
│   ├── ARCHITECTURE.md        # .spectacular/ structure, frontmatter, lifecycle, versioning
│   ├── ROADMAP.md             # Time-ordered "what's next"
│   ├── AGENTS.md              # Onboarding doc for agents working in .spectacular/
│   ├── STACK.md               # Host project's tech choices
│   ├── DECISIONS.md           # ADR-style decision log
│   ├── specs/                 # Per-capability specs (optional; SPEC.md is the index)
│   └── requests/              # Active and planned work
│       ├── spec-rename/             # SPEC.md + specs/ rename (status: active, v0.5.0)
│       ├── public-docs-foundation/  # docs/ first-class surface (status: planned)
│       └── public-docs-advanced/    # v2 docs surface (status: planned)
├── cli/                       # CLI implementation
│   ├── spectacular            # Bash binary — spectacular init
│   └── install.sh             # curl-installable installer
└── skills/spectacular/        # The skill itself
    ├── SKILL.md               # Lean orchestrator — triggers + routing table
    ├── references/            # Reference docs loaded on demand by the skill
    └── templates/             # Canonical templates (PRD kits, base)
```

## Active Requests

| Slug | Status | Summary |
|---|---|---|
| `spec-rename` | verified | `SPEC.md` + `specs/` replaced legacy `current/` (shipped v0.5.0) |
| `public-docs-foundation` | review | First-class `docs/` surface — `docs.yaml` + doc verbs + doctor area (shipping v0.6.0) |
| `public-docs-advanced` | planned | v2 docs — renderer adapters + versioning + spec→doc sync (gated on real demand) |
| `convention-pack-schema` | verified | Pack schema + bundled `minimal` pack + app-store folder (shipped v0.4.0) |
| `convention-pack-fabricator` | review | `pack-overrides.md` grill + `alex-default` dogfood — live grill scenarios pending |
| `convention-pack-application` | review | CLI `pack` subcommand + `convention_pack:` config + init/doctor wiring — live three-mode scenarios pending |
| `convention-pack-modules` | planned | v2 modular packs — stays planned until composition pain surfaces from v1 use |
| `doctor` | review | Substrate self-check (shipped v0.3.1) — interactive skill-side scenarios pending |
| `cli-bootstrap` | planned | Parked — kept for v0.2.x maintenance fixes |

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
| `spec-sync.md` | Proposing `SPEC.md` + `specs/` updates during archive (was `current-sync.md` pre-v0.5.0) |
| `memory.md` | `spectacular remember this` |
| `versioning.md` | `spectacular snapshot <file>` |
| `init-workflow.md` | `spectacular init` (CLI context) |
| `onboarding.md` | First invocation on existing `.spectacular/` project |
| `scaffold-reference.md` | File template reference — frontmatter stubs for all file types |
| `doc-registry.md` | Doc-type registry: template + slots + mode + location + overrides for every registered doc |
| `grill.md` / `refine.md` / `review.md` | Generic engine for any registered doc (driven by registry + per-doc overrides) |
| `prd-overrides.md` | PRD-specific rules: kit selection, slot prompts, vague-word list, gate checks |
| `plan-overrides.md` / `tasks-overrides.md` | PLAN/TASKS-specific rules consumed by the same engine |
| `pack-overrides.md` | Convention-pack-specific grill rules — slot prompts, source-ingestion (`--from`), reserved pack-ids |
| `kits-contract.md` | Kit extension schema: adds-slots, modifies-slots, triggers-docs |
| `packs-contract.md` | Convention-pack schema: 6 rule categories, 4-tier scope precedence, modular-pack v2 sketch |
| `verification.md` | 2-of-6 rule for when VERIFY.md is required vs folded into PLAN/TASKS |
| `doctor.md` | Substrate self-check spec — areas, severity model, repair flow (CLI mechanical + skill judgment) |
| `prd-grill.md` / `prd-refine.md` / `prd-review.md` | Legacy — superseded by generic engine + `prd-overrides.md` |

## Key Conventions

**Frontmatter is the signal layer** — the skill reads frontmatter, not full file content, during briefings.

**Lifecycle state lives in `PLAN.md` frontmatter only** (`status: planned | active | review | verified`). Never duplicated.

**Canonical documents are never overwritten in place** — always snapshot first (`PRD@v1.0.md` naming). The unversioned filename always points to current.

**Memory (`spectacular remember this`) writes to `.spectacular/memory/`** — git-committed, team-visible. Never to `.claude/` personal memory.

**`.spectacular/`** is fully committed to git. **`.spectacular.local/`** is always gitignored.

**Convention packs (v0.4.0+)** are opt-in repo-shape opinions declared via `.spectacular/config.yaml`'s `convention_pack:` block. Packs ship in four scope locations (project-local → user → app-store → bundled, in precedence order). The bundled `minimal` pack enforces only README contract + gitignore baseline; `alex-default` in the app-store is the fully-opinionated reference. Three modes (`suggest` / `scaffold` / `enforce`) control how strictly the pack is applied during init + doctor. Full schema: `skills/spectacular/references/packs-contract.md`.

**Two-layer task tracking:** harness `TaskCreate`/`TaskUpdate` = ephemeral session micro-tracker (drives CLI live progress UI); on-disk `requests/<slug>/TASKS.md` = persistent milestone blocks. Anti-pattern: one-for-one duplication. Full convention in `.spectacular/AGENTS.md` § Task tracking.

## Context Loading Rules (from .spectacular/AGENTS.md)

| Task type | Load |
|---|---|
| Planning / design | PRD.md, PRINCIPLES.md, DECISIONS.md |
| Refining intent / PRD work | PRD.md, skill refs prd-grill.md / prd-refine.md / prd-review.md |
| Skill implementation | STACK.md, ARCHITECTURE.md, PLAN.md, TASKS.md, SPEC.md, specs/skill/SPEC.md |
| CLI implementation | STACK.md, ARCHITECTURE.md, PLAN.md, TASKS.md, SPEC.md, specs/cli/SPEC.md |
| Review / QA | VERIFY.md, SPEC.md, specs/<capability>/SPEC.md, RISKS.md |
| Onboarding cold | PRD.md, SPEC.md, ARCHITECTURE.md, .spectacular/AGENTS.md |

Load only the capability spec relevant to the current task, not all of `specs/`. The top-level `SPEC.md` is cheap and always relevant. Never read `archive/` during normal operation. The authoritative loading table is `.spectacular/AGENTS.md`.

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

spectacular doctor                  # substrate self-check (8 areas)
spectacular doctor <area>           # scoped: skill|workspace|frontmatter|snapshots|links|lifecycle|kits|conventions
spectacular doctor --fix            # apply mechanical fixes

spectacular pack list               # show packs across all 4 scopes
spectacular pack install <name>     # copy to ~/.spectacular/packs/<name>/
spectacular pack show <name>        # print scope + frontmatter
spectacular pack remove <name>      # delete user-scope pack
```

Skill is fetched from GitHub (tagged release tarball → `.agents/skills/spectacular/`). `.claude/skills/spectacular/` is a symlink to the `.agents/` target. Version tracked in `.spectacular/skills.lock`.

Full command reference: `docs/commands.md`. Config schema: `docs/configuration.md`.
