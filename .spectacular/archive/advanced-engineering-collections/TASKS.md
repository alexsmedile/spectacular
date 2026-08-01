---
status: verified
updated: 2026-08-02
related:
  - PLAN.md
---

# Tasks — advanced-engineering-collections

<!--
  Executable checklist for one request.
  Lives at: .spectacular/requests/<slug>/TASKS.md

  Rules:
  - Group tasks by milestone using `### M<N> — <name>` headings.
  - Flush-left checkboxes are the COUNTED units: `- [ ]` open, `- [x]` done,
    `- [~]` deferred (not-open-not-done; shown separately in progress).
  - Indented `  - [ ]` sub-bullets are allowed as a nested acceptance checklist
    under a task, but are NOT counted — progress counts top-level only, so
    x/total stays comparable across requests.
  - `status:` in frontmatter should match parent PLAN.md.
  - Tasks are owned by the user. Engine never adds/removes/reorders tasks.
-->

## v1

### M1 — Optional init substrate
- [x] Add the five collection IDs to validated `--with` inputs
- [x] Add idempotent Git-preservable directory scaffolders without touching existing entries
- [x] Preserve the bare/minimal init shape
- [x] → check: focused init fixtures cover individual, bundled, and repeat scaffolding

### M2 — Reserved contracts
- [x] Reserve `FND`, `FIX`, `BUG`, `SEC`, and `BMK` with paths and collision-safe aliases
- [x] Document `security/` naming and legacy `F<N>` compatibility
- [x] Keep workflows/lifecycles explicitly unimplemented except the existing verified-fix ledger
- [x] → check: contract review finds no ambiguous active aliases or fabricated behavior

### M3 — Verification and docs
- [x] Update architecture, system spec, init workflow, commands, scaffold guide, and CLI help
- [x] Run syntax, version guard, focused tests, full suite, and diff integrity checks
- [x] Record docs-impact evidence and verification outcome
- [x] → check: all executable gates pass

## v2 (deferred)

- [~] Implement finding, bug, security, and benchmark lifecycle/mutator/doctor behavior after separate design interviews
- [~] Migrate legacy `F<N>` fixes to `FIX-NNN` and activate `f1` for findings only through an explicit preview-first, archive-first migration
