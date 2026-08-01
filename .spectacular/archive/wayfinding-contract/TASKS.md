---
status: verified
updated: 2026-08-01
related:
  - PLAN.md
---

# Tasks — wayfinding-contract

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

### M1 — Canonical identity contract
- [x] Specify canonical prefixes, aliases, padding, ambiguity, and compatibility rules
- [x] Implement shared ID allocation and resolution helpers
- [x] Add resolver fixtures for explicit, contextual, ambiguous, and legacy references
- [x] → check: canonical-ID focused tests pass under Bash 3.2 syntax

### M2 — Typed knowledge records
- [x] Add questions, research, and spike templates/rules/CLI verbs
  - [x] Questions template, rules, and CLI verbs
  - [x] Research and spike records; reserve `PRT` for resulting artifacts
- [x] Upgrade decision and idea creation to canonical filenames with compatibility reads
- [x] Add source/evidence and human-input-required fields
- [x] → check: collection CLI and doctor tests pass

### M3 — Specification confirmation
- [x] Define `unconfirmed | current | deprecated` specification state and transitions
- [x] Gate request generation on current specs and preserve source provenance
- [x] Document discovery/execution prerelease labels
- [x] → check: unconfirmed/current transition and request-gate tests pass

### M4 — Safe migration
- [x] Add dry-run migration for legacy decisions, ideas, and spec state
- [x] Rewrite canonical references while preserving archived originals/snapshots
- [x] Extend doctor with ID, prefix, padding, and dangling-reference checks
- [x] → check: migration fixtures and full doctor pass

## v2 (deferred)

- [~] Automatically widen every canonical ID to four digits at 1,000 entries — defer until a scale fixture proves the rewrite cost and UX
- [~] Assign canonical IDs to TASKS.md actions — `TSK` remains reserved until task identity has a concrete consumer
