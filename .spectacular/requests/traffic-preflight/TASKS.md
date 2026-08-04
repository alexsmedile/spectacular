---
status: review
updated: 2026-08-04
related:
  - PLAN.md
source_type: issue
source_ref: "alexsmedile/spectacular#3"
---

# Tasks — traffic-preflight

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

### M1 — Traffic contract
- [x] Document the local durable evidence and read-only/time-bound boundaries.

### M2 — Local preflight
- [x] Implement classification and link visibility without GitHub dependency or automatic mutation.

### M3 — Verification
- [x] Add and run focused state, insufficient-evidence, link, syntax, and regression checks.
## v2 (deferred)

- [~] <Deferred task>
- [~] <Deferred task>
