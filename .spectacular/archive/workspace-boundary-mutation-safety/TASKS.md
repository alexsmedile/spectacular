---
status: verified
updated: 2026-08-06
related:
  - PLAN.md
source_spec: SPC-006
source_type: spec
source_ref: SPC-006
---

# Tasks — workspace-boundary-mutation-safety

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

### M1 — Implement SPC-006
- [x] Define one shared schema-relation helper for `older`, `equal`, and `newer`
- [x] Every command that can write `.spectacular/` must refuse before its first
- [x] `spectacular status --against-latest` must distinguish older from newer:
- [x] Preserve existing older-schema migration behavior and equal-schema behavior.
- [x] Replace guidance that says `.spectacular.local/` takes precedence with an
- [x] Add filename-only detection for tracked `.spectacular.local/` paths before a
- [x] A detected tracked local path fails closed: stop the operation, report only
- [x] Retain lazy local creation and the existing root `.gitignore` protection. Do
- [x] Add synthetic fixtures covering older, equal, and newer schema relations for
- [x] Add synthetic tracked-local-path fixtures that prove filename-only output and
- [x] Align CLI help/reference guidance with the executable boundary and DEC-022.
- [x] Keep `CURRENT_SCHEMA` at `2.0`; do not add a `2.0 -> 3.0` registry entry,
- [x] Do not create, activate, or implement a migration request from this SPC.
- [x] Do not inspect, copy, archive, convert, back up, delete, or expose contents
- [x] Do not change GitHub configuration or broaden SPC-003's protected-security
## v2 (deferred)

- [~] <Deferred task>
- [~] <Deferred task>
