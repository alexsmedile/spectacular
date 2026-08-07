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
- [ ] Map existing Git/AFK/request helpers and add Bash 3.2-safe shared evidence primitives.
- [ ] Support optional request `branch:` and `base:` metadata and validate it conservatively.
- [ ] → check: focused status/path classification has no writes and reports unknown evidence honestly.

### M2 — Read-only preflight and plan
- [ ] Implement deterministic branch, divergence, status-group, request/spec, and path-disposition reports.
- [ ] Add one-safe-next-action, explicit blockers, unknown facts, and prospective mutation command output.
- [ ] → check: mixed current/unrelated paths separate without mutation.

### M3 — Explicit scoped preservation
- [x] Implement preview-first `workspace preserve` with explicit paths, branch name, staged set, and commit message.
- [x] Apply only after `--apply --yes`; prove undeclared paths remain untouched.
- [ ] → check: untracked spec preservation commits only its explicit path and never stashes/resets.

### M4 — Verified cleanup
- [x] Implement read-only branch inventory with local ancestry plus optional forge state.
- [ ] Implement guarded deletion requiring fresh shared base, reachability or durable declined evidence, and explicit confirmation.
- [ ] → check: stale-base and open-PR deletion refuse; declined committed branch remains preview-only until confirmed.

### M5 — Enforcement, documentation, and regression
- [ ] Gate Spectacular-owned Git mutation/publish/cleanup flows through common preflight without changing provider-neutral core behavior.
- [ ] Update skill/help/doctor documentation and add full regression, Bash 3.2, AFK, and lifecycle coverage.
- [ ] → check: full suite, syntax, version guard, and doctor checks pass.
## v2 (deferred)

- [~] <Deferred task>
- [~] <Deferred task>
