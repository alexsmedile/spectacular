---
status: verified
updated: 2026-08-04
related:
  - PLAN.md
source_spec: SPC-004
source_type: spec
source_ref: SPC-004
---

# Tasks — vision-workflow

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

### M1 — Vision contract and templates
- [x] Define the pre-request Vision lifecycle, spine, fragments, evidence boundary, and legacy compatibility.
  - [x] Replace the mandatory story/UI/architecture seed rule with proportional fragment selection.
  - [x] Add top-level spine and generic fragment templates while retaining legacy templates.

### M2 — CLI workflow
- [x] Implement Vision collection paths and `imagine`, list/show/add/react/propose/approve/derive behavior.
  - [x] Keep derivation agentic and approval explicitly human-authorized.
  - [x] Regenerate manifests mechanically without rewriting judgment state.

### M3 — Workflow composition
- [x] Align Idea, discovery, Spike, feedback, roadmap, SPC, request, and routing references around approved Vision.
  - [x] Remove the `experiment` feedback-loop alias.
  - [x] Prefer `direction-validation` while accepting legacy roadmap `prototype` values.

### M4 — Doctor and tests
- [x] Add compatibility-aware doctor validation and focused Vision workflow tests.
  - [x] Cover new and legacy locations, lifecycle/reaction gates, manifest fix, and derive refusal.
  - [x] Run Bash syntax, focused tests, full suite, and version guard.

### M5 — Truth, docs impact, and verification
- [x] Synchronize canonical architecture/system truth, assess public docs, and record verification evidence.
  - [x] Snapshot canonical docs before updates.
  - [x] Produce a Pageworks handoff instead of directly authoring substantial public docs.

## v2 (deferred)

- [~] Visual browser-based Vision gallery or hosted prototype renderer.
- [~] Automatic migration of legacy request-owned Vision workspaces.
