---
status: planned
updated: 2026-07-07
related:
  - PLAN.md
---

# Tasks — cli-path-abstraction

## v1

### M1 — Centralized variable declaration & audit pass
- [ ] Define global path variables at the top of `cli/spectacular` (e.g. `DIR_ROADMAPPED`, `DIR_DECISIONS`, `DIR_MEMORIES`, `DIR_SESSIONS`, `DIR_AUDITS`, `DIR_FIXES`, `DIR_IDEAS`, `DIR_DEBUGS`)
- [ ] Audit every function in `cli/spectacular` to find all hardcoded occurrences of these paths
- [ ] → check: `bash -n cli/spectacular` passes cleanly

### M2 — Migrate CLI commands to path variables
- [ ] Replace hardcoded directory strings with global variables inside CLI functions (e.g. `check_decisions()`, `cmd_roadmap()`, etc.)
- [ ] Ensure all fallback checks (e.g. `.spectacular/ROADMAP.md`) still work properly using variables
- [ ] → check: `bash -n cli/spectacular` passes cleanly and CLI loads

### M3 — Full verification and regression check
- [ ] Run the E2E test suite to verify no regressions in functionality
- [ ] Run pre-commit hooks checks
- [ ] → check: `bash tests/run.sh` passes 100% green

## v2 (deferred)

- [ ] Support custom project structure overrides in config.yaml mapping to custom paths
