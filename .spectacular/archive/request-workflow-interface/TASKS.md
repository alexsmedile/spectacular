---
status: verified
updated: 2026-08-02
related:
  - PLAN.md
---

# Tasks — request-workflow-interface

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

### M1 — Request workflow contract
- [x] Write the canonical command, phase-ownership, request-view, activation-baseline, alias, and transition contract.
- [x] Create SPC-002 with the confirmed request execution-context requirements and obtain explicit approval before activation.
- [x] Define stable human and JSON output schemas, failure exit codes, and the exact boundaries between stored, derived, and just-in-time context.
- [x] → check: the contract classifies every affected command as read-only, mechanical mutation, or agentic workflow with no overlapping canonical owner.

### M2 — Request context compiler
- [x] Extend `spectacular request <slug>` overview with source, mode, current task, blockers, and stable progress while remaining cheap to read.
- [x] Add `--brief` as an active-request-only implementation prompt assembled from PLAN, TASKS, SESSION, approved spec baseline, and execution boundaries.
- [x] Add milestone selection with `--milestone M2`, `-m 2`, and `-m2` normalization and deterministic handling of missing or completed milestones.
- [x] Make `--full` return the ordered request-owned Markdown bundle and list, but not dump, linked external records or binary artifacts.
- [x] Provide equivalent stable `--json` schemas and focused tests for every view and failure mode.
- [x] → check: focused fixtures prove the agent can resume one active milestone without loading full discovery history or unrelated request state.

### M3 — Approved-spec execution handoff
- [x] Add `spectacular request new <slug> --from SPC-001` with approved-only validation, duplicate detection, and reviewed PLAN/TASKS derivation.
- [x] Permit initial TASKS generation from an approved specification while forbidding silent scope addition, removal, or reordering after activation.
- [x] Record flat activation provenance fields for actor, time, specification version, and Git baseline without copying the specification body.
- [x] Make terminal `spec act` redirect to the agentic flow; make `/spectacular spec act SPC-001` find-or-create one request, run gates, activate, retrieve `--brief`, and initialize native agent planning.
- [x] Refuse activation on unresolved required-user questions, unapproved specs, incomplete PLAN/TASKS, held requests, ambiguous duplicate requests, or declared HITL gates.
- [x] → check: focused fixtures prove an approved spec reaches one baseline-anchored active request while every invalid or ambiguous path stops safely.

### M4 — Lean and logical command grammar
- [x] Add canonical noun-first request forms for new, list, advance, and archive while keeping `new`, `requests`, `advance`, `archive`, and progress views as compatibility aliases.
- [x] Add verb-first `/spectacular grill|refine|review <document> [target]` routing that prints the resolved target before mutation; keep document-first compatibility aliases.
- [x] Make `/spectacular verify` the normal owner of `review → verified` and require explicit recorded override semantics for any exceptional direct transition.
- [x] Consolidate status/help guidance so deterministic CLI reads and agentic interpretation are clearly layered rather than competing commands.
- [x] Keep AFK branch cleanup behavior and docs-impact storage compatible, but describe them as advanced branch hygiene and internal closure assessment rather than everyday workflow steps.
- [x] → check: routing tests prove canonical forms and aliases converge on the same behavior without removing or silently changing existing functionality.

### M5 — Integration, documentation, and dogfood
- [x] Update lifecycle, active-request, build-workflow, new-request, spec-lifecycle, status, archive, doc-writing, and task-tracking references with one consistent workflow.
- [x] Update CLI help, README, PRD interface description, commands, workflow, scaffold, architecture, and system-spec surfaces; keep the everyday path concise and advanced commands progressively disclosed.
- [x] Add doctor checks for source-spec baseline completeness, request view/schema drift, and invalid evidence-free verified transitions where mechanical validation is possible.
- [x] Dogfood `request --brief` on this request, compare its output with the native session plan, and record any missing context without expanding the v1 selector surface automatically.
- [x] Run syntax, version guard, focused suites, full regression suite, doctor, and diff hygiene checks; prepare SPEC-DELTA and verification evidence.
- [x] → check: all automated checks pass and a fresh agent can move from approved SPC to a bounded native implementation plan using the documented happy path.

## v2 (deferred)

- [~] Add arbitrary `--artifact` and repeatable `--section` request selectors — defer until real use shows the compiled brief and direct file reads are insufficient.
- [~] Remove compatibility aliases or redesign all non-request entity commands — defer to a separately approved breaking CLI release.
- [~] Introduce `REQ-NNN` identity or migrate the existing `build: bN` ledger — defer until the identity migration is explicitly designed and confirmed.
- [~] Add `spectacular afk branch ...` nesting — defer until AFK has more than one cleanup category.
