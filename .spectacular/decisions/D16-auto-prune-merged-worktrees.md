---
type: Decision
id: 01a029c7-e312-7569-9afd-4ed9b25709f9
title: Prune background execution worktrees automatically upon Evidence merge
created_by: Alex
created: "2026-08-22T16:02:40Z"
updated: "2026-08-22T16:02:40Z"
actor: Alex
actor_role: owner
ref: D16-auto-prune-merged-worktrees
question: When should ephemeral background runner worktrees be cleaned up?
disposition: prune-on-evidence-merge
rationale: >-
    Leaving completed worktrees on disk consumes disk space and clutters the filesystem. Once an
    Objective's tests pass, its verifiable Evidence is accepted, and its branch is merged into the
    Mission integration branch, the temporary worktree is automatically pruned (`git worktree remove`).
alternatives:
    - retaining all worktrees until entire Mission completion
    - requiring manual human worktree deletion
authority_basis: Owner explicitly approved Option A (Auto-prune on merge) in the design interview.
---

# Prune background execution worktrees automatically upon Evidence merge

## Decision
- Background worker worktrees (`../<repo>-<purpose>`) are automatically deleted as soon as their Evidence is verified and merged.
- Workspace disk stays clean and free of dead directories.
