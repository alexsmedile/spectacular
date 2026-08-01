---
status: verified
updated: 2026-08-01
related:
  - PLAN.md
---

# Tasks — lifecycle-contract

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

### M1 — Canonical lifecycle contract
- [x] Add the single lifecycle contract with entity states, transitions, evidence, ownership, and archive rules
- [x] Route entity references and canonical docs to the contract instead of restating incompatible enums
- [x] Add candidate roadmap state and root-anchor snapshot-only guidance
- [x] → check: contract review finds no undefined or multiply-owned lifecycle transition

### M2 — Specification lifecycle
- [x] Implement draft/unconfirmed creation, approved-only action, and evidence-anchored `spec implement`
- [x] Implement behavior-revision supersession and archive-only terminal movement
- [x] Remove `current` truth claims and derive request backlinks from request-owned `source_spec`
- [x] Make legacy status migration conservative, preview-first, archive-first, and explicitly confirmed
- [x] → check: specification lifecycle and migration fixtures pass

### M3 — Knowledge lifecycle integrity
- [x] Implement completed/result discovery records with mandatory evidence and inconclusive fog behavior
- [x] Add sourced user-derived decisions plus verified→superseded lineage and strict AFK autonomous-decision gates
- [x] Add question resolution provenance and memory active/superseded/retracted states
- [x] Make fixes verified-only and enforce unique slugs within each collection
- [x] → check: collection lifecycle and doctor fixtures pass

### M4 — AFK and documentation gates
- [x] Add durable AFK run authorization states, scope, permissions, and HITL gates
- [x] Replace reflog-only cleanup with durable Spectacular archive refs and restoration evidence
- [x] Add mandatory documentation-impact assessment at major-update, heavy-request verification, and session-end checkpoints
- [x] → check: AFK and documentation-impact gate fixtures pass

### M5 — Alignment and verification
- [x] Align CLI help, templates, doctor, architecture, system spec, skill references, and public docs
- [x] Add `PRT` and `TSK` deferred-design note to TODO.md
- [x] Run Bash syntax, version guard, focused suites, full suite, and doctor
- [x] Review compatibility, migration, archive, and no-data-loss behavior against the confirmed interview
- [x] → check: full verification passes with no unresolved lifecycle contradictions

## v2 (deferred)

- [~] Implement a standalone `prototypes/` collection — `PRT` remains a spike artifact reference until a concrete consumer exists
- [~] Assign canonical IDs to request tasks — `TSK` remains reserved until it can coexist cleanly with TASKS.md and harness tasks
