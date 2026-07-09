---
status: verified
updated: 2026-07-09
related:
  - PLAN.md
---

# Tasks — okf-restructure

## v1

### M1 — Architecture and Schema Definition
- [~] Define "project-wide anchors" and the OKF layout in `.spectacular/ARCHITECTURE.md` — OKF layout (plural dirs, `index.md` relocation, `D<N>`/`M<N>` prefixes, ID-namespace convention) fully documented; the concept is present (PRD.md called an "anchor doc"; namespace table marked "project-wide") but the exact term "project-wide anchors" was not coined as a defined heading. Cosmetic, deferred.
- [x] Update `skills/spectacular/references/doc-index.md` to register `index.md` file paths and the `specs/` consolidation
- [x] Update workspace base templates (`skills/spectacular/templates/`) to reflect the new layout
- [x] → check: ARCHITECTURE.md and doc-index.md describe the OKF structure

### M2 — CLI Restructuring
- [x] Update path constants, array mappings, and file creation helpers in `cli/spectacular`
- [x] Update `top_usage()` CLI help outputs to reflect plural folder names and indexing
- [x] Update `spectacular doctor` validation areas to check for plural folders and `index.md` targets
- [x] Update `doctor` checks to validate sequential prefix names (`D<N>-<slug>` and `M<N>-<slug>`)
- [x] → check: bash -n cli/spectacular exits 0 and paths command returns new configurations

### M3 — Skill Restructuring
- [x] Update Skill orchestrator triggers in `skills/spectacular/SKILL.md`
- [x] Replace all index references in `skills/spectacular/references/` (root index files → `<dir>/index.md`)
- [x] Update `grill`, `refine`, and `review` logic to load rules and rules overrides from the new paths
- [x] → check: grep search finds no root-level index files referenced in skill files

### M4 — Workspace Schema Migration
- [x] Write the `skills/spectacular/references/migrations/v06-to-v20.md` migration plan/frontmatter
- [x] Add the `v0.6-to-v2.0` mechanical migration functions in `cli/spectacular`
- [x] Implement folder renaming, index moving, and links rewriting inside the migration runner
- [x] → check: running migrate on a mock workspace upgrades it successfully to v2.0

### M5 — Verification & E2E Validation
- [x] Run the complete test suite (`bash tests/run.sh`) — 15/15 areas, 0 failed
- [x] Verify that running `spectacular doctor` on the migrated repository reports zero errors — 0 errors (1 pre-existing unrelated roadmap-drift warning)
- [x] Confirm git hooks are completely green — pre-commit version-guard passed on commit 68e439d
- [x] → check: tests/run.sh exits 0 and pre-commit check passes

## v2 (deferred)
- [~] Implement compose subcommands for multi-part specifications — out of scope for v1; tracked separately under convention-pack-modules
- [~] Add automated migration options for custom external convention packs — out of scope for v1; tracked separately under convention-pack-modules
