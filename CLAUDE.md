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
│   ├── PRINCIPLES.md          # Operating principles + runtime enforcement hooks
│   ├── ARCHITECTURE.md        # .spectacular/ structure, frontmatter, lifecycle, versioning
│   ├── AGENTS.md              # Onboarding doc for agents working in .spectacular/
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
├── cli/                       # CLI implementation
│   ├── spectacular            # Bash binary
│   └── install.sh             # curl-installable installer
├── skills/spectacular/        # The skill itself
│   ├── SKILL.md               # Lean orchestrator — triggers + routing table
│   ├── references/            # Reference docs loaded on demand by the skill
│   ├── templates/             # Canonical templates (PRD kits, base)
│   └── versions/              # Historical SKILL.md snapshots
├── agents/                    # ★ Subagent defs — SOURCE OF TRUTH (plugin-root convention)
│   │                          #   fleet grid: discover / apply / review × fix / build
│   ├── debug-investigator.md  #   discover · fix — where+why on an open bug (read-only)
│   ├── debug-fixer.md         #   apply · fix — a closed single-site fix (apply-only)
│   ├── debug-researcher.md    #   research — known-external-bug verdict (read-only, web)
│   ├── repo-explorer.md       #   discover · build — map a subsystem before planning (read-only)
│   ├── spec-builder.md        #   apply · build — a closed milestone brief (apply-only)
│   ├── code-reviewer.md       #   review · code — 5-lens findings over a diff (read-only)
│   ├── spec-reviewer.md       #   review · docs — punch list vs a doc's rules-file rubric (read-only)
│   └── test-verifier.md       #   verify — run a check / write a test to spec (apply-only, tests)
│                              #   .claude/agents/*.md are relative symlinks → ../../agents/
├── packs/                     # App-store convention packs (alex-default)
├── hooks/                     # ⚠ Claude Code / Codex PLUGIN event handlers
│                              #   (loaded by the plugin runtime when installed)
├── scripts/                   # Repo-local automation
│   └── hooks/                 # ⚠ GIT hooks (pre-commit, etc.) — managed by git-guard
├── tests/                     # Bash test suite (9 areas)
├── docs/                      # User-facing docs (commands, configuration, install…)
├── CHANGELOG.md
├── README.md
├── CONTRIBUTING.md
├── LICENSE
├── CLAUDE.md                  # This file
└── AGENTS.md                  # Contributor + agent guide (repo map, dev commands, conventions)
```

> **Note on the two `hooks/` directories:** `hooks/` at repo root contains
> **Claude/Codex plugin hooks** (`hooks.json`, `hooks-codex.json`) — these
> are read by the plugin runtime *after* a user installs spectacular as a
> plugin. `scripts/hooks/` contains **git hooks** (e.g. `pre-commit`) that
> run locally during development of *this* repo, managed by the git-guard
> skill. Same word, different runtimes. Don't conflate.

## Active Requests

Run **`spectacular status`** for the live fleet — it renders the active-request
table directly from each `requests/<slug>/PLAN.md` frontmatter plus grep-safe
body signals (Goal line, `x/total` task progress, current milestone). Use
`spectacular status <slug>` for one request's card and `spectacular status --json`
as the machine-readable contract. No table is hand-cached here anymore — that
list drifted, so it now single-sources from the files (status-fleet-view, b23).

Version targets live in the ROADMAP ledger, not here. For full context on any request see `.spectacular/requests/<slug>/PLAN.md`.

**Parked as idea (not active request):** `memory-protocols` — abandoned 2026-05-26 as too broad to ship as a single milestone (6 usecase families × 8 patterns × 6 ideas × 13 open questions). All research consolidated into `.spectacular/ideas/memory-protocols.md` (PLAN + TASKS + IDEAS merged). A narrower spec request will be cut from it.

**Archived (shipped):** see `.spectacular/archive/`.

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
| `spec-sync.md` | Proposing `specs/index.md` + `specs/` updates during archive (was `current-sync.md` pre-v0.5.0) |
| `memory.md` | `spectacular remember this` |
| `versioning.md` | `spectacular snapshot <file>` |
| `init-workflow.md` | `spectacular init` (CLI context) |
| `onboarding.md` | First invocation on existing `.spectacular/` project *with prior work* (references `status.md` for the shared read+briefing flow since v1.21.0) |
| `guided-first-run.md` | First invocation on a *fresh/empty* workspace — ushers new→PRD→first request one step at a time (v1.21.0) |
| `scaffold-reference.md` | File template reference — frontmatter stubs for all file types |
| `doc-index.md` | Human-readable catalog of every doc type. Dispatch lives in each `<doc>-rules.md` frontmatter (since v1.4.0; was `doc-registry.md` pre-v1.4.0) |
| `soft-db-index.md` | Canonical routing index for the 7 soft-DB collections (memory/decisions/sessions/ideas/feedback/audit/fixes) — role, purpose, structure, boundary rules. Loaded when deciding *where* knowledge belongs (v1.25.0) |
| `bug-workflow.md` | Bug handling **runtime core** — check prior fixes first, decide audit-first vs just-fix, log a reusable fix. Ties audit/ + fixes/ into a self-learning loop (v1.25.0). Orchestrator arc + debug-agent fleet routing; Step 2b's two no-fix disposition forks (`folded-into-request`, `wont-fix`) (v1.26.0). Rationale split out to `bug-workflow-doctrine.md` |
| `bug-workflow-doctrine.md` / `build-workflow-doctrine.md` | The *why* behind each workflow's gates (rationale, failure modes, the relation table). Loaded only when a routing call is uncertain or when editing the workflows — never on routine dispatch |
| `debug-trace.md` | In-flight trace schema for a debug job — one folder per job under `.spectacular/debugs/<slug>/`, one JSON artifact per agent turn, all persisted by the orchestrator from returned blocks. Distinct from the audit/fixes ledger. Spines are CLI-validated by `doctor debug` (v1.26.0) |
| `grill.md` / `refine.md` / `review.md` | Generic engine for any registered doc (driven by registry + per-doc rules) |
| `prd-rules.md` | PRD-specific rules: kit selection, slot prompts, vague-word list, gate checks (superseded legacy `prd-grill.md` / `prd-refine.md` / `prd-review.md`) |
| `plan-rules.md` / `tasks-rules.md` | PLAN/TASKS-specific rules consumed by the same engine |
| `pack-rules.md` | Convention-pack-specific grill rules — slot prompts, source-ingestion (`--from`), reserved pack-ids |
| `kits-contract.md` | Kit extension schema: adds-slots, modifies-slots, triggers-docs |
| `packs-contract.md` | Convention-pack schema: 6 rule categories, 4-tier scope precedence, modular-pack v2 sketch |
| `verify.md` | The interactive validation walk — walk-only runtime core since b30 (v1.20.0 had merged three docs; authoring half split back out) |
| `verify-authoring.md` | Authoring-time verification: the 2-of-6 rule (canonical), fold patterns, VERIFY.md shape, promoting checks to `tests/verify/` scripts (b30) |
| `doctor.md` | Substrate self-check spec — areas, severity model, repair flow (CLI mechanical + skill judgment) |

## Key Conventions

**Frontmatter is the signal layer** — the skill reads frontmatter, not full file content, during briefings.

**Lifecycle state lives in `PLAN.md` frontmatter** (`status: planned | active | review | verified`). TASKS.md carries a mirror for skim tooling; PLAN is authoritative and `doctor` repairs drift.

**ReadOnly docs are never overwritten in place** — always snapshot first (`PRD@v1.0.md` naming). The unversioned filename always points to current.

**Memory (`spectacular remember this`) writes to `.spectacular/memories/`** — git-committed, team-visible. Never to `.claude/` personal memory.

**`.spectacular/`** is fully committed to git. **`.spectacular.local/`** is always gitignored.

**Convention packs (v0.4.0+)** are opt-in repo-shape opinions declared via `.spectacular/config.yaml`'s `convention_pack:` block. Packs ship in four scope locations (project-local → user → app-store → bundled, in precedence order). The bundled `minimal` pack enforces only README contract + gitignore baseline; `alex-default` in the app-store is the fully-opinionated reference. Three modes (`suggest` / `scaffold` / `enforce`) control how strictly the pack is applied during init + doctor. Full schema: `skills/spectacular/references/packs-contract.md`.

**Two-layer task tracking:** harness `TaskCreate`/`TaskUpdate` = ephemeral session micro-tracker (drives CLI live progress UI); on-disk `requests/<slug>/TASKS.md` = persistent milestone blocks. Anti-pattern: one-for-one duplication. Full convention in `.spectacular/AGENTS.md` § Task tracking.

## Context Loading Rules (from .spectacular/AGENTS.md)

| Task type | Load |
|---|---|
| Planning / design | PRD.md, PRINCIPLES.md, decisions/index.md |
| Refining intent / PRD work | PRD.md, skill refs prd-rules.md / grill.md / refine.md / review.md |
| Skill implementation | STACK.md, ARCHITECTURE.md, PLAN.md, TASKS.md, specs/index.md, specs/doc-engine.md |
| CLI implementation | STACK.md, ARCHITECTURE.md, PLAN.md, TASKS.md, specs/index.md, relevant specs/<capability>.md |
| Review / QA | VERIFY.md, specs/index.md, relevant specs/<capability>.md, RISKS.md |
| Onboarding cold | PRD.md, specs/index.md, ARCHITECTURE.md, .spectacular/AGENTS.md |

Load only the capability spec relevant to the current task, not all of `specs/`. The top-level `specs/index.md` is cheap and always relevant. Never read `archive/` during normal operation. The authoritative loading table is `.spectacular/AGENTS.md`.

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

Skill is fetched from GitHub (tagged release tarball → `.agents/skills/spectacular/`). `.claude/skills/spectacular/` is a symlink to the `.agents/` target. Version tracked in `.spectacular/skills.lock`.

Full command reference: `docs/commands.md`. Config schema: `docs/configuration.md`.
