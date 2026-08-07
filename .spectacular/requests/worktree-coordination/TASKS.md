---
status: active
updated: 2026-08-07
related:
  - PLAN.md
contract: 019fdd18-96a0-73b5-81b5-38cda98d36f4
---

# Tasks — worktree-coordination

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

### M1 — Evidence model and request metadata
- [x] Map existing Git/AFK/request helpers and add Bash 3.2-safe shared evidence primitives.
- [x] Support optional request `branch:` and `base:` metadata and validate it conservatively.
- [x] → check: focused status/path classification has no writes and reports unknown evidence honestly.

### M2 — Read-only preflight and plan
- [x] Implement deterministic branch, divergence, status-group, request/spec, and path-disposition reports.
- [x] Add one-safe-next-action, explicit blockers, unknown facts, and prospective mutation command output.
- [x] → check: mixed current/unrelated paths separate without mutation.

### M3 — Explicit scoped preservation
- [x] Implement preview-first `workspace preserve` with explicit paths, branch name, staged set, and commit message.
- [x] Apply only after `--apply --yes`; prove undeclared paths remain untouched.
- [x] → check: untracked spec preservation commits only its explicit path and never stashes/resets.

### M4 — Verified cleanup
- [x] Implement read-only branch inventory with local ancestry plus optional forge state.
- [x] Implement guarded deletion requiring fresh shared base, reachability, archive evidence, and explicit confirmation.
- [x] → check: stale-base and open-PR deletion refuse; unmerged/declined work remains preserved until explicit resolution.

### M5 — Enforcement, documentation, and regression
- [x] Gate Spectacular-owned Git mutation/publish/cleanup flows through common preflight without changing provider-neutral core behavior.
- [x] Update skill/help/doctor documentation and add full regression, Bash 3.2, AFK, and lifecycle coverage.
- [x] → check: full suite, syntax, version guard, and doctor checks pass.
## v2 (deferred)

- [~] <Deferred task>
- [~] <Deferred task>
