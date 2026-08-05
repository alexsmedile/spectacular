---
status: active
updated: 2026-08-05
related:
  - PLAN.md
source_type: issue
source_ref: "alexsmedile/spectacular#12"
---

# Tasks — soft-db-projection-inventory

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

### M1 — Align with merged #11 retrieval baseline
- [x] Record #11's merged `status --brief --json` orientation contract and exclude a competing workspace briefing from #12.
- [x] Confirm the first #12 slice is collection-only named details plus omission tests; filters, caches, and YAML-schema hardening remain deferred.

### M2 — Projection contract and regression fixtures
- [x] Implement the named entity-detail projections from `PROJECTION-INVENTORY.md` for the highest-value missing views, preserving literal `--full` evidence expansion.
- [x] Add per-entity omission fixtures and deterministic human/JSON output tests for requests, decisions, memories, questions, research, spikes, ideas, audits, and fixes.

### M3 — Compatibility, guidance, and measured follow-up
- [x] Update CLI-first guidance with the explicit direct-named-file escape hatch, only after commands and tests exist.
- [x] Re-run the agreed retrieval measurements and decide, from evidence, whether compiled workspace briefing, constrained filters, index checks, or caching earns a separate request.
## v2 (deferred)

- [~] General equality filters or an indexed-frontmatter subset contract (separate approved scope only).
- [~] Derived-index rebuild/doctor repair and cache design (separate approved scope only).
