---
type: Mission
id: 01a02bdf-17a0-7514-a626-51ac67c5be6a
title: Enforce disjoint write reservations and Handoff validation
status: active
created: "2026-08-22T23:49:28Z"
updated: "2026-08-22T23:52:49Z"
activation:
    at: "2026-08-22T23:49:28Z"
    by: Alex
    fingerprint: sha256:6cbe3976064e575995e14fbd58fc15527d967b5c7ae8f7bc6af5a265970cdd6c
authority:
    operator:
        - inspect
        - edit-in-scope
        - choose-reversible-implementation
        - run-checks
        - generate-derived-files
        - bounded-repair
        - commit-local
    requires_owner:
        - activate-mission
        - change-outcome-or-completion
        - expand-scope
        - push
        - merge
        - release
        - irreversible-change
        - destructive-data
baseline:
    branch: codex/m19-disjoint-write-reservations
    commit: 78dddac4c2fe67f5cb8cb5168ad2cd25fa278856
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
contract:
    fingerprint: sha256:62cc4645130cdade0ae4e5a25b32ae00e8ed3e4df74860c70fae1b91c7968339
    ref: Contract:01a00aae-8921-7b27-96a9-1a4c175e7dc6
dependencies:
    - M18 completed with Objective-scoped Run lifecycle
gaps: []
objectives:
    - claims:
        - disjoint-write-reservations
      id: 01a02bdf-17a0-7b0a-bc29-3f0949c43547
      outcome: Validate and enforce disjoint writable path reservations across Handoffs.
      ref: O1
      status: implemented
    - claims:
        - dependency-locked-runs
      id: 01a02bdf-17a0-7891-bf4d-ae7b82a67bab
      outcome: Enforce upstream dependency locks on Objective Run starts.
      ref: O2
      status: implemented
    - claims:
        - passive-git-state-inspection
      id: 01a02bdf-17a0-73f0-aa04-60750538f0b8
      outcome: Verify passive Git sanity without destructive side effects.
      ref: O3
      status: implemented
outcome: Orchestrators can safely parallelize Objectives because Handoffs enforce exact disjoint writable path reservations and upstream dependency locks.
owner: Alex
ref: M19
repair_budget: 2
review: independent
reviews:
    - file: reviews/RV1-independent-review-of-m19-disjoint-write-reservations-and-dependency-locks.md
      id: 01a02be1-1b40-7796-9699-ecff63f5892b
      ref: RV1
      verdict: pass
run:
    current_objective: O1
    id: 01a02bdf-17a0-7515-a902-bc090cf28f8c
    operator: Alex
    ref: R1
    repairs: 0
    started_at: "2026-08-22T23:49:28Z"
    status: active
scope:
    mechanical:
        - internal/missionbundle
        - internal/command
        - internal/runtime
        - test/evals/spectacular
    semantic:
        - Handoff write reservations
        - path disjointness validation
        - upstream run dependency locking
        - passive Git sanity check
start_key: sha256:42f3d54be7b8e98e6824d446825161bce8b9df3fd3f91b832329a7eb8fb6266a
stops:
    - overlapping-write-reservations
    - blocked-upstream-dependency
    - git-mutation-by-spectacular
    - path-escape
    - data-loss
validation:
    mode: cli
    schema: mission.v2
---
# Mission: Enforce Disjoint Write Reservations and Handoff Validation

## Purpose & Scope
Guarantees concurrent subagents never edit overlapping files. Slices M19 to 3 core invariants: path disjointness checking, upstream dependency locking, and passive Git safety.

## Key Deliverables & Actions
1. **Disjoint Write Path Math (`internal/missionbundle/handoff.go`)**: Validate `writes: [...]` paths; refuse intersecting write perimeters across active Handoffs.
2. **Upstream Dependency Lock (`internal/missionbundle/service.go`)**: In `StartRun`, verify all upstream dependencies in `depends_on` are clean and unblocked.
3. **Passive Git Check (`internal/runtime/`)**: Inspect `.git` state for active rebase/merge conflicts without modifying or stashing files.
