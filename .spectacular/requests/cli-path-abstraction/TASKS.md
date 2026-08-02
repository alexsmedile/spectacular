---
status: review
updated: 2026-08-02
related:
  - PLAN.md
---

# Tasks — cli-path-abstraction

## v1

### M1 — Centralized variable declaration & audit pass
- [x] Define global path variables at the top of `cli/spectacular` (e.g. `SPEC_ROADMAPS_DIR`, `SPEC_DECISIONS_DIR`, `SPEC_MEMORIES_DIR`, `SPEC_SESSIONS_DIR`, `SPEC_AUDITS_DIR`, `SPEC_FIXES_DIR`, `SPEC_IDEAS_DIR`, `SPEC_DEBUGS_DIR`)
- [x] Audit every function in `cli/spectacular` to find all hardcoded occurrences of these paths
- [x] → check: `bash -n cli/spectacular` passes cleanly

### M2 — Migrate CLI commands to path variables
- [x] Replace hardcoded directory strings with global variables inside CLI functions (e.g. `check_decisions()`, `cmd_roadmap()`, etc.)
- [x] Ensure all fallback checks (e.g. `.spectacular/ROADMAP.md`) still work properly using variables
- [x] → check: `bash -n cli/spectacular` passes cleanly and CLI loads

### M3 — Full verification and regression check
- [x] Run the E2E test suite to verify no regressions in functionality
- [x] Run pre-commit hooks checks
- [x] → check: `bash tests/run.sh` passes 100% green

## v2 (deferred)

- [~] Support custom project structure overrides in config.yaml mapping to custom paths — deferred to a separately designed v2 request; this change intentionally fixes the canonical layout
