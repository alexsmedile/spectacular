---
status: verified
updated: 2026-08-04
related:
  - PLAN.md
source_spec: SPC-004
source_type: spec
source_ref: SPC-004
---

# Tasks — idea-destination-routing

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

### M1 — Implement SPC-004
- [x] Implement explicit `request`, `roadmap`, and `shared` IDEA destinations.
  - [x] Preserve legacy request promotion and undo semantics.
  - [x] Require a stable shared reference and keep the path local-only.
  - [x] Require an explicit Icebox choice for roadmap transfer.
- [x] Add focused CLI tests for every destination and compatibility guardrail.

### M2 — Align routing contract
- [x] Update the idea, soft-DB, GitHub bridge, roadmap, and request references.
- [x] Update the skill, template, and architecture description without changing
      unrelated GitHub bridge or signal-capture behavior.

### M3 — Verify boundaries
- [x] Run focused and regression CLI tests plus syntax/version checks.
- [x] Record verification evidence and assess documentation impact.
## v2 (deferred)

- [~] <Deferred task>
- [~] <Deferred task>
