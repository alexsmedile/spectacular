---
status: verified
updated: 2026-08-06
related:
  - PLAN.md
source_spec: SPC-005
source_type: spec
source_ref: SPC-005
---

# Tasks — afk-orchestration-authority

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

### M1 — Implement SPC-005
- [x] Define the optional new-format run fields and CLI event grammar; preserve legacy records.
  - [x] Include authorization identity, session, base, integration branch, goal nodes, and action/gate declarations.
  - [x] Keep all current AFK Git-mutating operations unchanged and outside M1.
- [x] Implement read-only run inspection and doctor validation for authority/event evidence.
  - [x] Detect malformed fields, duplicate/non-monotonic IDs, invalid chronology, and missing required evidence.
- [x] Add focused Bash tests for new-format, legacy, gated, invalid, and no-Git-mutation behavior.
- [x] Run focused and baseline verification; assess documentation impact.
## v2 (deferred)

- [~] <Deferred task>
- [~] <Deferred task>
