---
type: MissionPlan
title: Enforce disjoint write reservations and Handoff validation
owner: Alex
contract:
  ref: Contract:01a00aae-8921-7b27-96a9-1a4c175e7dc6
outcome: Orchestrators can safely parallelize Objectives because Handoffs enforce exact disjoint writable path reservations and upstream dependency locks.
review: independent
completion:
  - claim: disjoint-write-reservations
    pass_boundary: Handoff write scopes accept repo-relative exact files/folders, forbid globs and parent traversal (../), and refuse dispatch when overlapping with any active Handoff.
    proof_requirement: Table-driven matrix tests assert that overlapping paths, nested folders, and path escapes trigger typed RefusalInvalidScope errors before write.
  - claim: dependency-locked-runs
    pass_boundary: run start on an Objective refuses if an upstream dependent Objective has an active Run in blocked or stopped state without an owner Decision unblock.
    proof_requirement: Test fixture proves run start refuses on blocked/stopped upstream dependencies and passes when completed or unblocked by Decision D<N>.
  - claim: passive-git-state-inspection
    pass_boundary: Spectacular verifies clean Git branch state and rejects execution during active rebase/merge conflicts without mutating Git or deleting worktrees.
    proof_requirement: Unit fixtures with clean, dirty, and conflicting Git states assert typed refusal without executing destructive Git side effects.
objectives:
  - outcome: Validate and enforce disjoint writable path reservations across Handoffs.
    claims: [disjoint-write-reservations]
  - outcome: Enforce upstream dependency locks on Objective Run starts.
    claims: [dependency-locked-runs]
  - outcome: Verify passive Git sanity without destructive side effects.
    claims: [passive-git-state-inspection]
authority:
  operator: [inspect, edit-in-scope, choose-reversible-implementation, run-checks, generate-derived-files, bounded-repair, commit-local]
  requires_owner: [activate-mission, change-outcome-or-completion, expand-scope, push, merge, release, irreversible-change, destructive-data]
scope:
  mechanical: [internal/missionbundle, internal/command, internal/runtime, test/evals/spectacular]
  semantic: [Handoff write reservations, path disjointness validation, upstream run dependency locking, passive Git sanity check]
repair_budget: 2
dependencies: [M18 completed with Objective-scoped Run lifecycle]
gaps: []
stops: [overlapping-write-reservations, blocked-upstream-dependency, git-mutation-by-spectacular, path-escape, data-loss]
---

# Mission: Enforce Disjoint Write Reservations and Handoff Validation

## Purpose & Scope
Guarantees concurrent subagents never edit overlapping files. Slices M19 to 3 core invariants: path disjointness checking, upstream dependency locking, and passive Git safety.

## Key Deliverables & Actions
1. **Disjoint Write Path Math (`internal/missionbundle/handoff.go`)**: Validate `writes: [...]` paths; refuse intersecting write perimeters across active Handoffs.
2. **Upstream Dependency Lock (`internal/missionbundle/service.go`)**: In `StartRun`, verify all upstream dependencies in `depends_on` are clean and unblocked.
3. **Passive Git Check (`internal/runtime/`)**: Inspect `.git` state for active rebase/merge conflicts without modifying or stashing files.
