---
type: Decision
id: 01a02a64-3e8b-7cc5-950b-d99773bdf0e5
title: Propose verified Git cleanup and confirm it at Mission closure
created_by: Alex
created: "2026-08-22T16:54:04Z"
updated: "2026-08-22T16:54:04Z"
actor: Alex
actor_role: owner
ref: D19-confirm-cleanup-at-mission-closure
question: When and under whose authority should Objective and Mission Git state be cleaned up?
disposition: owner-confirmed-cleanup-after-proof
rationale: >-
    Removing a clean worktree after its work is merged is recoverable because the branch and commits
    remain, while deleting the remaining branches is a broader lifecycle cleanup. Both actions must
    follow verified integration, and the final cleanup belongs to Mission closure. The Orchestrator
    should propose exact native Git operations and the owner should confirm them rather than allowing
    automatic deletion.
alternatives:
    - automatic worktree deletion immediately after Evidence merge
    - retaining every worktree and branch indefinitely
authority_basis: Owner required merged and tested state before cleanup, explicit confirmation of the Orchestrator's proposal, and final cleanup as a Mission-closure responsibility.
authorized_effects:
    - contract.version-bump
conditions:
    - exact-cleanup-targets-proposed
    - merged-and-tested
    - clean-and-verifiable-state
    - owner-confirmation
scope:
    - v2
targets:
    - Proposal:01a029be-b7d3-703c-a7ee-50c6b8bae3a2
supersedes: D16-auto-prune-merged-worktrees
---

# Propose verified Git cleanup and confirm it at Mission closure

## Decision

- After a Run or Objective is merged and tested, the Orchestrator may propose removing its clean worktree with native Git.
- The owner confirms every cleanup operation before execution.
- Objective branches remain available until the Mission's final integration and closure boundary.
- Final removal of remaining worktrees and fully merged temporary branches is a Mission-closure responsibility.
- Dirty, unmerged, conflicted, or unverifiable state is preserved and never force-removed.
