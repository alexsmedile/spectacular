---
status: verified
updated: 2026-08-01
related:
  - PLAN.md
---

# Tasks — wayfinding-sequencer

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

### M1 — Dependency and readiness engine
- [x] Implement graph loading from canonical cross-references
- [x] Derive frontier/fog/deferred/resolved views and detect invalid graphs
- [x] Rank spike/research ahead of lower-uncertainty work at equal priority
- [x] → check: DAG and ranking fixtures pass

### M2 — Wayfinding commands and briefing
- [x] Add `wayfind status|next|resolve|defer` surfaces
- [x] Surface concise user-input blockers at session startup
- [x] Preserve deferred open loops without presenting them as current blockers
- [x] → check: CLI and session briefing fixtures pass

### M3 — Metaphoric intent layer
- [x] Route park, icebox, find-your-way, and act-on-goal language
- [x] Gate act-on-goal on a current spec and normal request lifecycle
- [x] Park execution discoveries without expanding current milestone scope
- [x] → check: phrase-routing and boundary tests pass

### M4 — Cross-layer coherence
- [x] Analyze PRD/spec/plan/roadmap links for missing declared dependencies
- [x] Warn on discovery/execution target-version inversions
- [x] Add doctor output with proposed remediation, never silent reslotting
- [x] → check: cross-layer fixtures produce deterministic warnings

## v2 (deferred)

- [~] Mermaid/ASCII exploration graph rendering — nice-to-have after the resolver proves useful
- [~] Parallel research fan-out — enable only after three independent frontier nodes create real orchestration pain
