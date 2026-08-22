---
type: Decision
id: 01a029c7-e30c-740d-8fdf-91891d8342f5
title: Enforce automatic branch and worktree guardrails at Mission activation
created_by: Alex
created: "2026-08-22T16:02:40Z"
updated: "2026-08-22T16:02:40Z"
actor: Alex
actor_role: owner
ref: D15-branch-guardrail-at-activation
question: How should Spectacular guard against accidental in-place execution on the default branch?
disposition: enforce-feature-branch-or-worktree
rationale: >-
    Executing changes directly on the `main` or default branch risks corrupting the clean trunk,
    complicating rollbacks, and causing branch lock collisions during concurrent development.
    `spectacular mission start` mechanically allocates a dedicated branch or worktree (`feat/M<N>-<slug>`),
    refusing execution directly on `main` unless an explicit `--allow-main` flag is provided.
alternatives:
    - unblocked execution directly on main
    - warning messages that do not enforce branch isolation
authority_basis: Owner explicitly approved Option A (Mechanical Branch/Worktree Guardrail) in the design interview.
---

# Enforce automatic branch and worktree guardrails at Mission activation

## Decision
- `spectacular mission start` allocates a clean branch/worktree (`feat/M<N>-<slug>`).
- Direct activation on `main` is mechanically refused without `--allow-main`.
