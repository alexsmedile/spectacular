---
status: archived
shipped_in: v0.6.1 (Stage 1) + v0.6.2 (Stage 2)
archived: 2026-05-23
priority: high
owner: alex
updated: 2026-05-23
target_version: 0.6.1 (Stage 1) / 0.6.2 (Stage 2)
summary: "Workspace migration framework — track structural schema version, detect drift, migrate between Spectacular versions. Stage 1 ships discoverability + flat-shape support + v0.4→0.5 migration as v0.6.1 hotfix; Stage 2 ships the full registry + skill walk as v0.6.2."
related:
  - ../../archive/spec-rename/PLAN.md
  - ../../archive/doctor/PLAN.md
  - ../../ARCHITECTURE.md
  - ../../../skills/spectacular/references/doctor-substrate.md
provenance:
  - source: SPECTACULAR-FEEDBACK from octopus repo (2026-05-23)
    captures: "User flattened current/specs/ → specs/ manually because no discoverability surface revealed the migration; flat SCHEMA-*.md pattern not handled by capability-only model; PRINCIPLES/ARCH/ROADMAP missing-detection invisible."
---

# Plan — Workspace Migrations

## Goal

Make structural evolution of `.spectacular/` safe, explicit, and **discoverable**. Every time Spectacular changes workspace shape, users with older workspaces should be:

1. **Detected** — `spectacular status --against-latest` gives a one-line "you are behind" verdict; doctor's `specs` area validates shape; future `scaffold` audit covers always-set drift.
2. **Guided** — `spectacular migrate` enumerates pending migrations (description + dry-run + apply); `/spectacular migrate` walks judgment-required ones with snapshot-before-edit.
3. **Tracked** — `workspace_schema:` field in `config.yaml` records the version the workspace was last migrated to.
4. **Verifiable** — each migration is idempotent, dry-run-able; chain validated by doctor.

Replaces the current pattern of ad-hoc fixes (doctor's `specs` area has bespoke `current/ → specs/` logic) with a registered, versioned framework.

## Why (now, with Octopus feedback)

Three real incidents:

- **v0.5.0 `current/ → specs/`** — handled by hand-coded doctor mechanical fix. Worked, but the logic is buried in `check_specs` and **not discoverable**. Octopus user did the rename by hand without ever running doctor.
- **v0.6.0 always-set 5 → 6** (added `specs/`) — `doctor workspace` flags missing dir, `--fix` re-scaffolds. Works but conflates "missing because broken" with "missing because old".
- **Flat schema-doc layout (Octopus)** — `specs/SCHEMA-TASK.md` etc. is a legitimate pattern (frontmatter contracts as primary truth), but `check_specs` assumes capability subfolders. Doctor would nag.

A migration framework + discoverability surface + flat-shape support collapses this to:

- **One verb** (`spectacular migrate`) to apply
- **One verb** (`spectacular status --against-latest`) to detect
- **One field** (`workspace_schema:`) to track
- **Two valid specs layouts** (capability subfolders + flat contract docs) coexist without nagging

## Scope

### Stage 1 — v0.6.1 hotfix (this request, current focus)

Smallest unit that unblocks projects in the Octopus situation.

- **`workspace_schema:` field** in `config.yaml` (default `"0.6"` on init; absent treated as `"0.4"`)
- **`spectacular status --against-latest`** — one-line "you are at schema X, current is Y; run `migrate`"
- **`check_specs` upgrade** — accepts both shapes:
  - Capability subfolders: `specs/<capability>/SPEC.md` (existing behavior)
  - Flat contract docs: top-level `.md` files in `specs/` (new — treated as "contract docs")
  - Mixed: both kinds in same `specs/` tree (validated each per their kind)
- **Doctor scaffold suggestion** — `check_specs` (or a small additional check) emits one info line summarizing missing v0.6+ conventional files (`PRINCIPLES.md`, `ARCHITECTURE.md`, `ROADMAP.md`) if absent, with the exact init command to scaffold them. Single line, not three.
- **v0.4→0.5 backfilled migration** — auto-rename `current/` → `specs/`, preserve flat layout. Mechanical. Idempotent.
- **v0.5→0.6 backfilled migration** — ensure `specs/` exists as always-set (mkdir + .gitkeep). Mechanical. Idempotent.
- **Minimal `spectacular migrate` CLI verb** — `migrate` runs known migrations to bring `workspace_schema` to current; `--dry-run` previews. Both backfilled migrations are mechanical so the CLI does all the work. Judgment-required migrations print a message pointing at `/spectacular migrate` but don't exist yet in Stage 1.

### Stage 2 — v0.6.2 (deferred; tracked here)

- Full migration registry under `skills/spectacular/references/migrations/`
- Migration frontmatter contract (`from`, `to`, `description`, `mechanical`, `reversible`, `affects`)
- `migrations-contract.md` schema doc
- `references/migrate.md` — skill flow for `/spectacular migrate` (judgment migrations + snapshots)
- `--to <version>` / `--from <version>` flags (re-run, partial migrations)
- Per-migration smoke tests under `tests/cli/migrations/`
- Chain validation in doctor (new `migrations` area or extension of existing)
- Auto-prompt in `init --update` after skill upgrade

### Out of scope (both stages)

- Migrations for skill internals (kit format changes, registry format changes)
- Code/data migrations in the host project
- Bidirectional migration (downgrade) for non-reversible changes
- Cross-major-version "skip" migrations (always sequential through chain)
- GUI / interactive UI beyond CLI prompts

## Decisions (locked 2026-05-23 via interview)

- **Specs shape:** both capability subfolders AND flat contract docs at `specs/` root are valid. No mode flag — `check_specs` accepts both and validates each per its kind.
- **PRINCIPLES/ARCH/ROADMAP:** stay opt-in. Doctor emits one info line summarizing missing v0.6+ conventional files (not per-file). No errors, no warnings.
- **Discoverability:** both `spectacular status --against-latest` (quick) and doctor area output (detailed). Same data, two access paths.
- **`workspace_schema:` lives in `config.yaml`** (not a separate lock file, not on SPEC.md frontmatter).
- **`workspace_schema:` defaults to `"0.4"` if absent** — that's the version before the first registered migration.
- **Migrate runner split:** `spectacular migrate` (CLI) does mechanical migrations + emits scriptable errors. `/spectacular migrate` (skill) handles judgment migrations.
- **v0.4→0.5 logic:** auto-rename `current/` → `specs/` unconditionally, preserve whatever layout (flat or subfolders) existed inside. Zero judgment.
- **Doctor area:** extend existing `check_specs` for both shapes. No new `scaffold` area in Stage 1 (only one info line added). Stage 2 may grow this.

## Validation (Stage 1)

- New workspace via `init` writes `workspace_schema: "0.6"`
- Workspace with `workspace_schema: "0.4"` (or absent), `current/` present, `specs/` absent:
  - `spectacular status --against-latest` reports "schema 0.4 detected; current 0.6 — run `spectacular migrate`"
  - `spectacular migrate --dry-run` lists 2 pending migrations with descriptions
  - `spectacular migrate` applies both; `workspace_schema` becomes `"0.6"`; re-running is clean no-op
- Workspace at current schema: `migrate` exits clean ("workspace is up to date")
- Workspace with both `current/` AND `specs/` (broken state): `migrate` refuses, points at doctor for manual resolution (existing behavior in `check_specs`)
- `specs/` with flat `SCHEMA-*.md` files: `doctor specs` reports them as "contract docs" passing validation, not as missing capability folders
- `specs/` with mixed (some subfolders + some flat .md): both kinds validated; no false positives
- `doctor specs` on a v0.6+ project missing PRINCIPLES/ARCH/ROADMAP: one info line listing all three + the `init --with ...` command
- Octopus workspace shape (flat SCHEMA-*.md + missing PRINCIPLES/ARCH/ROADMAP) runs through doctor with zero errors/warnings (only info)

## Stage 1 Milestones

1. **M1 — `workspace_schema:` field** — add to `doc_config()`; `CURRENT_SCHEMA="0.6"` constant; init writes it; init backfills on existing workspace
2. **M2 — `spectacular status --against-latest`** — flag on existing status verb; reads `workspace_schema:`; prints one line
3. **M3 — `check_specs` upgrade** — accept flat contract docs alongside capability subfolders; new pass/info messages per kind
4. **M4 — Scaffold suggestion** — single info line in `check_specs` (or `check_workspace`) for missing v0.6+ conventional files
5. **M5 — `spectacular migrate` CLI (minimal)** — apply known migrations to bring `workspace_schema:` forward; `--dry-run`; mechanical-only in Stage 1
6. **M6 — Backfilled migrations** — v0.4→0.5 (current→specs, preserve layout) + v0.5→0.6 (mkdir specs/.gitkeep)
7. **M7 — Tests + docs + v0.6.1 release** — `tests/cli/migrate.test.sh`; `tests/cli/status.test.sh` extension; CHANGELOG; SPEC.md capability bullet

## Risks

- **Stage 1 ships migrate without the registry pattern from Stage 2** — risk that the CLI's hardcoded migration list grows ugly before Stage 2 lands. Mitigation: only 2 migrations in Stage 1; if Stage 2 slips, can still add a 3rd manually without pain. Hard cap of 4 hardcoded migrations before Stage 2 must ship.
- **Flat contract docs vs capability subfolders ambiguity at `specs/` root** — a single `.md` file could be either "I'm a contract doc" or "I forgot the subfolder for my capability". Mitigation: convention — contract docs are kebab-case files with no `SPEC.md` name; capability subfolders are dirs containing `SPEC.md`. Document explicitly in `doctor-areas.md` and `SPEC.md`'s self-doc.
- **`workspace_schema:` field forgotten on init for existing projects** — `init` doesn't touch existing config.yaml. Mitigation: `doctor --fix` adds the field if missing (defaults to current schema since absence usually means "I'm new" not "I'm on 0.4"). Plus: `spectacular migrate` itself writes the field on first successful run.
- **Status `--against-latest` requires knowing the latest** — CLI hardcodes `CURRENT_SCHEMA`. Stale CLI reports stale "latest". Mitigation: documented; future `init --update` syncs CLI then status answer is correct.

## Open questions (for Stage 2)

- Should migrations be allowed to require user input (interactive)? Probably yes via skill, no via CLI.
- How does this interact with `convention-pack` migrations? Packs may evolve their own schema; for now pack versioning is pack-scoped, not workspace-scoped.
- Cross-host registry: if a third party publishes a Spectacular extension with its own structural rules, can it ship migrations? Stage 2+ concern.
