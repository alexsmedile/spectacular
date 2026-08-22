---
type: Decision
id: 01a02a64-3e8a-7450-a10d-52047fffc54b
title: Inspect native Git state at Mission start without owning Git effects
created_by: Alex
created: "2026-08-22T16:54:04Z"
updated: "2026-08-22T16:54:04Z"
actor: Alex
actor_role: owner
ref: D18-inspect-native-git-at-mission-start
question: Which Git responsibilities belong to mission start and which remain native Git operations?
disposition: inspect-and-refuse-with-native-git-effects
rationale: >-
    Spectacular needs an honest activation baseline but should not become a wrapper around Git.
    Branches and worktrees are prepared with native Git before activation; mission start can then
    deterministically inspect and refuse unsafe state without stashing, switching, allocating, or
    repairing repository state behind the owner's back.
alternatives:
    - automatic branch and worktree allocation inside mission start
    - accepting dirty or interrupted Git state as an activation baseline
authority_basis: Owner confirmed that Git branching and worktree creation remain separate native commands and mission start only inspects a clean prepared branch/worktree.
authorized_effects:
    - contract.version-bump
conditions:
    - clean-working-tree
    - non-default-branch
    - no-interrupted-git-operation
scope:
    - v2
targets:
    - Proposal:01a029be-b7d3-703c-a7ee-50c6b8bae3a2
supersedes: D15-branch-guardrail-at-activation
---

# Inspect native Git state at Mission start without owning Git effects

## Decision

- Native Git creates and switches the Mission branch/worktree before `mission start`.
- `mission start` refuses a dirty tree, a default branch, or an interrupted merge, rebase, cherry-pick, revert, or bisect.
- `mission start` never creates, switches, stashes, or repairs Git state.
- Concurrent Missions use separate integration branches/worktrees. Objective branches/worktrees start from their Mission branch and merge back into it.
- The Spectacular CLI records and verifies Git isolation but performs no branch, worktree, merge, push, or deletion effect.
