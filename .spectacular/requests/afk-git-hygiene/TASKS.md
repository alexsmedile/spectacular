---
status: verified
updated: 2026-08-01
related:
  - PLAN.md
---

# Tasks — afk-git-hygiene

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

### M1 — AFK policy and opt-in
- [x] Define AFK authorization, primary-branch, dirty-tree, and breaking-change gates
- [x] Add project configuration and read-only policy/status output
- [x] Reconcile Spectacular names with host-required branch prefixes
- [x] → check: default invocation mutates nothing and explains the gate

### M2 — Isolated branch proposals
- [x] Generate names for draft specs, spikes, forks, and confirmed execution
- [x] Add clean-tree and current-branch preflight
- [x] Record branch-to-node/spec/request provenance
- [x] → check: branch fixtures enforce naming and safety rules

### M3 — Archive-first hygiene
- [x] Capture spike outcome, evidence, and disposition before cleanup
- [x] Produce dry-run cleanup plans for merged or abandoned playground branches
- [x] Require confirmation before local/remote deletion
- [x] → check: cleanup fixtures preserve recoverable records and default to no deletion

### M4 — Human-review handoff
- [x] Gate handoff on confirmed spec, verification evidence, and passing tests
- [x] Open `[SpecTACular] Executed: <version> - <name>` PRs when authorized
- [x] Stop before merge and breaking API/schema changes
- [x] → check: PR handoff fixtures pass and HITL gates remain enforced

## v2 (deferred)

- [~] Automatic remote branch deletion — keep human-confirmed until repeated cleanup volume proves a safer policy
- [~] Worktree-based concurrent writers — defer until actual collision/latency evidence overrides the single-mutator model
