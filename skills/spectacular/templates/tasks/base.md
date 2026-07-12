---
status: planned
updated: <DATE>
related:
  - PLAN.md
---

# Tasks — <Request Title>

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

### M1 — <Milestone name>
- [ ] <Task>
- [ ] <Task>
- [ ] <Task>
- [ ] → check: <how M1 proves itself — see PLAN §6 Validation>

### M2 — <Milestone name>
- [ ] <Task>
- [ ] <Task>
- [ ] → check: <how M2 proves itself>

### M3 — <Milestone name>
- [ ] <Task>
- [ ] → check: <how M3 proves itself>

## v2 (deferred)

- [~] <Deferred task>
- [~] <Deferred task>
