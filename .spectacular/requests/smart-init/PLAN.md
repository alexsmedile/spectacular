---
status: verified
priority: high
owner: alex
updated: 2026-05-22
summary: "Smart spectacular init — minimal always-set by default, kit/flag-driven doc scaffolding, interactive doc-by-doc prompt, non-overwrite pre-flight"
related:
  - ../doc-writer/PLAN.md
  - ../kits-as-plugins/PLAN.md
  - ../doctor/PLAN.md
  - ../cli-bootstrap/PLAN.md
  - ../../ARCHITECTURE.md
  - ../../../skills/spectacular/references/verification.md
---

# Plan — Smart Init

## Goal

Upgrade the existing `spectacular init` CLI to scaffold **only what the project needs** instead of all 7 root docs every time. Default to a 5-file always-set; expand via kit declarations (always + suggested) or explicit flags. Interactive mode asks per suggested doc. Never overwrite existing files.

## Why

Today's `spectacular init` (v0.2.0) scaffolds all 7 root docs (PRD, PRINCIPLES, ARCHITECTURE, ROADMAP, STACK, DECISIONS, AGENTS) for every project. Reality:
- A **skill project** rarely needs STACK or ARCHITECTURE — `SKILL.md` is its architecture
- A **content project** rarely needs STACK or ARCHITECTURE at all
- A **research project** needs PRD + ROADMAP + maybe DECISIONS
- A **CLI tool** needs STACK + ARCHITECTURE

Scaffolding everything unconditionally creates **stub fatigue** — empty files that nobody fills in, that the skill nonetheless reads on every briefing, that dilute the signal. Better: scaffold the minimum, let the kit declare extras, let interactive mode ask.

## Scope

**Layer affected:** CLI only. The skill's runtime triggers/routing/state-awareness are unchanged. Smart-init upgrades the bootstrap entry point.

**In scope (v1)**
- Define the always-set: `PRD.md`, `requests/`, `current/`, `config.yaml`, `AGENTS.md`
- Non-interactive default: scaffold only the always-set + `blank` kit (no inference)
- Flag plumbing: `--kit <name>`, `--with <doc1,doc2>`, `--minimal` (always-set only, ignore kit defaults)
- Interactive mode (`-i`): ask kit → scaffold kit's always-docs → ask y/n per suggested doc
- Pre-flight non-overwrite: if a target file exists, skip it, report "already present, leaving alone"
- Defer drift/repair to [[doctor]] — smart-init emits a **generic "run diagnostics" message** when it detects malformed/old-schema files; the message stays generic until doctor ships, then references the doctor verbs explicitly
- Update `init-workflow.md` reference doc to reflect new behavior
- VERIFY.md (see § Verification routing)
- Tests at repo root `tests/` (see § Tests)

**Out of scope (v2)**
- Auto-detect project type from repo signals (`package.json`, `SKILL.md`, `.claude-plugin/`) — explicitly deferred per user decision
- Retrofit command (`spectacular scaffold <type>`) for existing repos
- Kit marketplace / fetching kits from URLs
- Multi-kit application at init time (single-kit per [[kits-as-plugins]] v1)

**Explicit anti-patterns**
- Scaffolding every root doc by default — replaced by minimal + kit-driven
- Silent overwrite of existing files — always pre-flight, always skip
- Coupling kit selection to PRD-craft flow only — kits also drive init
- Project-type inference in v1 — explicitly punted to v2

## Constraints

- Backwards compatible: existing workspaces must continue to work; this changes init behavior, not runtime behavior
- Idempotent: re-running `spectacular init` on an existing workspace must be safe (skip everything that exists)
- Bash-only CLI — no new language dependencies
- Must work with `--global` flag (install to `~/.agents/` and `~/.claude/`)
- `cli-bootstrap` remains open in parallel for v0.2.x maintenance fixes — smart-init does not absorb it

## Verification routing

This request hits the **2-of-6 rule** in [[verification]] with 3 axes:
1. User-visible change — CLI behavior changes for every existing user ✓
2. Multi-surface verification — fresh dir, existing workspace, interactive, all flag combos ✓
3. External contract change — new flags become part of the user-facing CLI contract ✓

→ **VERIFY.md is justified.** Smart-init becomes the first VERIFY.md-bearing request in this project.

VERIFY.md lives at `requests/smart-init/VERIFY.md` and carries the manual QA / edge cases / regression / rollback checklists. PLAN § Validation lists per-milestone criteria; VERIFY.md carries the procedural checks. Per the convention, **both are load-bearing** — `verified` status requires every `- [x]` in VERIFY.md.

## Milestones

1. **Always-set defined** — document the 5-file minimum; rationale in `init-workflow.md`
2. **Flag interface** — `--kit`, `--with`, `--minimal` flags wired in `cli/spectacular`
3. **Pre-flight check** — non-overwrite logic with clear stdout reporting (see Pre-flight behavior table below)
4. **Interactive mode** — `-i` walks kit selection + per-suggested-doc y/n
5. **Kit consumption** — CLI reads `triggers-docs` from selected kit's frontmatter
6. **Tests + VERIFY.md** — `tests/cli/` Bash test harness + VERIFY.md manual QA checklist
7. **Dogfood** — fresh project with `init`, `init -i`, `init --kit coding`; each produces correct doc set; re-running is no-op

## Pre-flight behavior (re-init semantics)

`spectacular init` is **always idempotent + non-destructive**. Re-running on an initialized workspace is safe by design.

| State | Behavior | Stdout |
|---|---|---|
| File doesn't exist | Create with stub | `✓ created PRD.md` |
| File exists, empty (0 bytes or only whitespace) | Fill with stub | `✓ filled empty PRD.md` |
| File exists, has content | **Skip** — never overwrite | `⊘ PRD.md already present, leaving alone` |
| File exists with malformed frontmatter | Skip; emit generic diagnostics message | `⊘ PRD.md present (issues detected — run diagnostics via \`spectacular doctor\` once available)` |
| File exists with old schema version | Skip; emit generic diagnostics message | `⊘ PRD.md present (schema mismatch — run diagnostics via \`spectacular doctor\` once available)` |
| Directory doesn't exist | Create silently | (no output) |
| Directory exists | No-op silently | (no output) |
| `.gitignore` entry missing | Append entry only | `✓ added .spectacular.local/ to .gitignore` |
| `.gitignore` entry present | No-op | (no output) |

**Decided exclusions:**
- **No `--force` flag** — violates the "never destructive" principle. To regenerate a stub, the user deletes the file manually first, then re-inits.
- **No schema migration in init** — that's [[doctor]]'s job. Init only creates; doctor diagnoses + opt-in repairs. Until doctor ships, init's diagnostics messages stay generic.
- **No prompt on existing files** — silent skip + report. Asking would slow down the common case (re-running init after adding a kit).
- **No project-type inference** — bare init scaffolds the blank kit unconditionally. Inference deferred to v2.

**Adding a kit later is safe:** `spectacular init --kit coding` on an existing workspace only adds the kit's missing always-docs; it never touches anything that already exists. This is the canonical way to "upgrade" a project from blank kit to coding kit.

## Tests

Bash test harness at repo root: `tests/cli/init.test.sh`. Future non-CLI tests can sit beside it (`tests/skill/`, `tests/registry/`).

Test runner:
- `tests/run.sh` — runs all `tests/**/*.test.sh`, reports pass/fail
- Each test creates an isolated `/tmp/spectacular-test-<n>/` directory, runs CLI commands, asserts file presence/content, cleans up

Coverage maps 1:1 with VERIFY.md scenarios (6 core) — see § Verification routing.

## Tasks

See `TASKS.md`.

## Dependencies

- **Hard dependency on [[doc-writer]]** ✓ verified — registry exists, init consumes it for doc IDs
- **Hard dependency on [[kits-as-plugins]]** ✓ verified — kit `triggers-docs.always/suggested` parseable
- **Soft dependency on [[doctor]]** — diagnostics messages stay generic until doctor ships; **no blocker**. Replace generic message with `spectacular doctor frontmatter` / `spectacular doctor schema` references when doctor lands.
- **Touches [[cli-bootstrap]]** — `cli-bootstrap` stays open in parallel for v0.2.x maintenance; smart-init modifies the same `cli/spectacular` binary but ships under v0.3.0

## Validation

Per-milestone criteria. Procedural checks live in VERIFY.md.

- M1 — `init-workflow.md` documents the 5-file always-set with rationale
- M2 — `--kit`, `--with`, `--minimal` flags parse correctly; `--help` updated; invalid combos error cleanly
- M3 — pre-flight matches the 9-state table; idempotent re-run exits 0; `.gitignore` append-only
- M4 — `-i` flow walks kit menu (built from `templates/prd/kits/*.md` frontmatter) + per-suggested-doc y/n prompts
- M5 — CLI parses `triggers-docs.always/suggested` from selected kit's frontmatter; resolves to scaffold list; unknown doc IDs error
- M6 — `tests/cli/init.test.sh` covers all 6 VERIFY scenarios with pass/fail assertions; VERIFY.md scaffolded with 6 checklist items
- M7 — manual dogfood on fresh `/tmp/` workspace matches expected file sets per flag combination

## Deliverables

- Updated `cli/spectacular` with new flag parsing + pre-flight + kit consumption
- Updated `cli/install.sh` if flag changes affect install messaging
- Updated `references/init-workflow.md` reflecting always-set + kit-driven model
- Updated `.spectacular/ARCHITECTURE.md § Init flow` section
- New `tests/cli/init.test.sh` + `tests/run.sh` harness
- VERIFY.md at `requests/smart-init/VERIFY.md`
- Migration note in CHANGELOG.md for v0.3.0 (init behavior change)
