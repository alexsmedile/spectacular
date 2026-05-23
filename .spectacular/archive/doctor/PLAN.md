---
status: verified
shipped_in: v0.3.1
priority: medium
owner: alex
updated: 2026-05-23
summary: "Spectacular doctor — environment/infrastructure self-check (CLI detects, skill repairs)"
related:
  - ../../archive/doc-writer/PLAN.md
  - ../../archive/smart-init/PLAN.md
  - ../../ARCHITECTURE.md
  - ../../../skills/spectacular/references/verification.md
---

# Plan — Doctor

## Goal

Add a `spectacular doctor` command that performs a **self-check of the Spectacular workspace's infrastructure** — skill install, registry parseability, frontmatter validity, snapshot continuity, cross-doc link validity, lifecycle hygiene. CLI detects; the skill handles judgment-requiring fixes.

## Why

Workspaces drift over time. Real drift already observed in this project:
- Missing `PRD@v1.1.md` snapshot (v1.0, v1.2, v1.3 exist; v1.1 gap unexplained)
- Frontmatter schema migrations (PRD v1.x → v2.0 split) leaving older fields stale
- `spectacular init` emits "run diagnostics via `spectacular doctor` once available" for malformed/old-schema cases — doctor is what that message points to

Doctor is also the **skill's own self-check** — when the skill hits an operational failure mid-flow (registry won't parse, kit file malformed), it runs the relevant subset of doctor's checks to surface the substrate issue.

## Scope

**Layer affected:**
- **CLI side** — new `spectacular doctor` subcommand that runs detection scripts and emits a structured report
- **Skill side** — interactive repair flow at `/spectacular doctor --fix` that consumes the CLI's report, proposes fixes per finding, walks user through `y/n/q`

**Conceptual framing:** doctor is an **environment/infrastructure self-check**, NOT a content-quality check. Content quality is what `<doc> review` does. Doctor checks whether the workspace *can do its job*, not whether the writing is good.

**In scope (v1) — CLI detect**
- `spectacular doctor` — runs all checks, prints status report (✅/⚠️/❌), exits with code (0 clean, 1 warnings, 2 errors)
- `spectacular doctor <area>` — scoped checks (`skill`, `workspace`, `frontmatter`, `snapshots`, `links`, `lifecycle`)
- `spectacular doctor --fix` — auto-applies content-free mechanical repairs (gitignore append, missing dirs, broken symlinks, missing always-set stubs). Everything else is reported for skill-side repair.
- `--format text` (default) and `--format json` (for skill consumption)
- Report can be written to `.spectacular/.doctor-report.json` for the skill to read

**In scope (v1) — Skill repair surface**
- `/spectacular doctor` — invokes CLI in detect mode, summarizes report conversationally
- `/spectacular doctor --fix` — reads CLI report; for each finding requiring judgment, proposes a concrete fix with context; walks `y/n/q` per finding; snapshots canonical docs before any edit (per existing `versioning.md` rule)

**In scope (v1) — Skill-invoked subset**
- When skill flows hit infrastructure failures mid-operation (e.g. registry parse error during `status`, malformed kit file during `prd grill`), skill auto-runs the relevant subset of doctor's checks and surfaces results inline without requiring a full sweep

**Check categories (v1)**

| Area | What it verifies |
|---|---|
| `skill` | `.agents/skills/spectacular/` exists; symlink valid; `skills.lock` consistent with installed ref; `SKILL.md` parseable |
| `workspace` | `.spectacular/` directory present; always-set files exist; `config.yaml` parseable |
| `frontmatter` | Every canonical doc has its registry-required fields; schema versions current |
| `snapshots` | Version sequence has no gaps (`v1.0` + `v1.2` without `v1.1`) |
| `links` | `related:` frontmatter entries point at existing files |
| `lifecycle` | `status: active` requests have SESSION.md; `status: review` requests have verification artifact (VERIFY.md OR TASKS § Verification OR PLAN § Validation); `status: verified` requests >7 days old surfaced as archive candidates |
| `kits` | Bundled + project-local kit files parse; `triggers-docs` entries reference known doc IDs |

**Out of scope (v2)**
- Cross-doc semantic validation (PRD goals ↔ PLAN milestones) — belongs in [[doc-writer]] v2
- `prd diff` / `prd merge` subcommands — moved to doc-writer scope (they're PRD operations, not infrastructure checks)
- Vague-word leakage detection — already in `prd review`, no duplication
- Multi-workspace diagnostics
- Performance benchmarks
- Auto-fix without confirmation (no `--yes-to-all`)
- Git history rewriting (snapshot backfill proposes new commits only)
- JSON Schema validation of registry entries (defer)
- `doctor` running automatically on every CLI invocation — explicit trigger only

**Explicit anti-patterns**
- Content-quality checks in doctor — that's `<doc> review`'s job. Doctor checks the substrate, not the writing.
- CLI proposing judgment-requiring fixes — CLI detects mechanical drift, agent handles judgment
- `--fix` bypassing per-finding confirmation — every fix needs explicit `y` (`--yes-to-all` rejected)
- Auto-running doctor on every operation — explicit trigger only
- Duplicating `prd review`'s vague-word scan, slot-presence check, etc. — gate logic lives in the engine

## Verification routing

Applying the 2-of-6 rule from [[verification]]:

1. **User-visible change** — ✓ new CLI command + new skill verb
2. **High reversibility cost** — ⚠️ partial (per-finding confirm, snapshot-before-edit; reversible per change but multi-file scope)
3. **Multi-surface verification** — ✓ CLI detect, CLI mechanical fix, skill repair, skill-invoked subset, multiple areas, json format
4. **Risk surface** — ✓ touches canonical docs during repair (auth/billing parallel: trust + integrity)
5. **External contract change** — ✓ new CLI surface + new skill verb
6. **Rollback** — ✓ snapshot-before-fix is mandatory; revert = `git checkout` the snapshot

→ **Score: 5 of 6** → VERIFY.md scaffolded. Doctor's VERIFY will be more comprehensive than smart-init's (which scored 3).

## Milestones

1. **Detection taxonomy + report format defined** — doc-registry.md gains a `gate-checks:` field per doc OR `references/doctor.md` enumerates checks by area; report format (text + JSON) locked
2. **CLI detect-only mode** — `spectacular doctor [<area>]` walks all checks, emits text report + exit code; `--format json` writes structured findings
3. **CLI mechanical fixes** — `--fix` flag handles the small content-free repair set (gitignore append, missing dirs, dangling symlinks, missing always-set stubs); everything else reported for skill-side repair
4. **Skill repair flow** — `/spectacular doctor --fix` reads CLI report; per-finding agent proposes context-aware fix; user `y/n/q` confirms each; snapshot-before-edit per existing convention
5. **Skill-invoked subset** — wire other skill flows (`status`, `prd grill`, `archive`) to call doctor checks for relevant areas when they hit failures
6. **Smart-init message update** — replace the generic "run diagnostics via `spectacular doctor` once available" with explicit area pointers (`run \`spectacular doctor frontmatter\``)
7. **Tests + VERIFY.md** — `tests/cli/doctor.test.sh` covers detect mode + mechanical fix mode + scoped areas + json output; VERIFY.md mirrors with agent-flow scenarios that can't be automated
8. **Dogfood** — run doctor on this very workspace; surface and (where applicable) fix at least one real drift item (missing `PRD@v1.1.md` snapshot is the obvious candidate)

## Tasks

See `TASKS.md`.

## Dependencies

- **Hard dependency on [[doc-writer]]** ✓ verified — needs the registry to know what "correct" frontmatter looks like for each doc
- **Hard dependency on [[kits-as-plugins]]** ✓ verified — `kits` check needs the kit contract schema
- **Hard dependency on [[smart-init]]** ✓ verified — uses the same `is_empty_file()`, `in_list()`, kit-parser helpers; smart-init's generic diagnostics message becomes specific once doctor ships
- **Touches [[verification]]** — doctor's lifecycle check enforces the verification convention (verified status requires evidence)

## Validation

Per-milestone criteria. Procedural checks live in VERIFY.md.

- M1 — `references/doctor.md` documents every check (area, severity, what it inspects, what it proposes)
- M2 — `spectacular doctor` on a clean workspace exits 0; on this workspace surfaces real findings (missing `PRD@v1.1.md`)
- M3 — `spectacular doctor --fix` on a workspace missing `.gitignore` append-only entry → adds it without prompt (mechanical); on a workspace with frontmatter drift → reports + suggests `/spectacular doctor --fix`
- M4 — `/spectacular doctor --fix` walks each finding requiring judgment; user `y` triggers snapshot + edit; `n` skips; `q` exits
- M5 — `/spectacular status` on a workspace with a malformed registry auto-runs doctor's `kits` + `frontmatter` checks and surfaces results inline
- M6 — smart-init's diagnostic message now reads `run \`spectacular doctor frontmatter\`` instead of the generic placeholder
- M7 — `tests/cli/doctor.test.sh` covers detect + scoped + mechanical-fix + json scenarios with assertions
- M8 — manual dogfood against `/Users/alex/vault/data/skills_db/spectacular/` produces a useful report

## Deliverables

- Updated `cli/spectacular` — new `doctor` subcommand parsing + detection logic + mechanical-fix logic
- New `references/doctor.md` — check definitions, severity levels, repair flow documentation
- New `tests/cli/doctor.test.sh` — automated coverage of CLI behavior
- VERIFY.md at `requests/doctor/VERIFY.md` — manual QA checklist for agent-flow scenarios
- Updated `cli/spectacular` (smart-init's message replaced — moves "v0.3.1 polish" item from smart-init's deferred list)
- SKILL.md routing additions: `spectacular doctor`, `spectacular doctor <area>`, `spectacular doctor --fix`
- Migration note in CHANGELOG.md for v0.4.0 (or v0.3.x if treated as a non-breaking addition)

## CLI vs Skill split (locked)

The CLI detects; the skill repairs (with two exceptions for trivial mechanical fixes).

| Concern | CLI does it | Skill does it |
|---|---|---|
| Stat the filesystem | ✓ | (skill calls into CLI) |
| Parse frontmatter (YAML) | ✓ (awk) | reads CLI's report |
| Detect missing required fields | ✓ | — |
| Detect snapshot version gaps | ✓ | — |
| Detect dangling cross-doc links | ✓ | — |
| Detect lifecycle hygiene issues | ✓ | — |
| Emit text + JSON report | ✓ | (consumes) |
| Append to `.gitignore` (mechanical) | ✓ (`--fix`) | — |
| Create missing always-set dirs (mechanical) | ✓ (`--fix`) | — |
| Repair broken/dangling symlinks (mechanical) | ✓ (`--fix`) | — |
| Re-stub missing always-set file (mechanical) | ✓ (`--fix`) | — |
| Add missing frontmatter field with judgment | — | ✓ (proposes value with context) |
| Reconcile schema drift (v1.x → v2.0) | — | ✓ (proposes field migration) |
| Backfill missing snapshot from git log | — | ✓ (proposes content or punts) |
| Propose link removal vs link rename | — | ✓ (asks user) |
| Surface verified-and-stale requests for archive | (detects) | ✓ (proposes archive flow) |
| Run repair without explicit `y` per finding | ✗ (never) | ✗ (never) |

The CLI is honest about its scope: mechanical detection + a small set of content-free repairs. Anything requiring judgment is the agent's job.
