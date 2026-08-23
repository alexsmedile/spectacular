---
type: MissionPlan
title: Repair Context Sandwich steering invariants
owner: Alex
contract:
  ref: Contract:01a00a20-63dd-7670-97f1-9eb8e12adc3a
outcome: Repair proof and behavior gaps across Handoff reservations, Decision unblocks, worktrees, Evidence targets, and pinned fixtures.
review: independent
completion:
  - claim: m19-reservation-and-git-repair
    pass_boundary: Handoff write reservations walk full supersedes chains; Decision unblocks permit run start on blocked dependencies; git safety resolves linked-worktree gitdirs.
    proof_requirement: Table-driven tests prove worktree gitdir conflict detection, Decision-based run start unblocks, and active handoff reservation calculations.
  - claim: m20-evidence-target-integrity
    pass_boundary: Evidence recording strictly validates that all named Objectives, Runs, and Claims resolve to real identifiers on the target Mission.
    proof_requirement: Refusal tests prove invalid or foreign Objective, Run, or Claim references trigger typed RefusalInvalidReference errors before write.
  - claim: m21-pinned-benchmark-matrix
    pass_boundary: Scope hardening and context savings benchmarks execute against immutable pinned fixtures rather than live repo files with zero false rejections.
    proof_requirement: Pinned fixture tests prove ≥40% context savings and 100% discriminating scope protection across scaffolding, renames, and broad directories.
objectives:
  - outcome: Repair Handoff reservation chains, Decision unblocks, and linked-worktree Git safety.
    claims: [m19-reservation-and-git-repair]
  - outcome: Enforce strict target reference validation on Evidence packages.
    claims: [m20-evidence-target-integrity]
  - outcome: Anchor scope hardening and context benchmarks in pinned immutable fixtures.
    claims: [m21-pinned-benchmark-matrix]
authority:
  operator: [inspect, edit-in-scope, choose-reversible-implementation, run-checks, generate-derived-files, bounded-repair, commit-local]
  requires_owner: [activate-mission, change-outcome-or-completion, expand-scope, push, merge, release, irreversible-change, destructive-data]
scope:
  mechanical: [internal/missionbundle, internal/runtime, internal/command, test/evals/spectacular, test]
  semantic: [linked worktree git inspection, Decision-based dependency unblocking, strict Evidence target validation, pinned benchmark fixtures]
repair_budget: 2
dependencies: [M21 completed with campaign baseline]
gaps: []
stops: [second-dependency-graph, live-file-benchmark-drift, unvalidated-evidence-targets, silent-worktree-git-skip, data-loss]
---

# Mission: Repair Context Sandwich Steering Invariants

## Purpose & Scope
Repairs subtle semantic gaps identified by independent review across M19, M20, and M21 without altering completed historical mission bundles.

## Key Deliverables & Actions
1. **M19 Repair (`internal/missionbundle/`)**:
   - `git_safety.go`: Parse `gitdir:` pointer in `.git` files for linked worktrees.
   - `service.go`: In `StartRun`, allow starting runs on objectives with blocked dependencies if a Decision in the workspace explicitly resolves it.
   - `service.go`: In `recordHandoff`, walk the full `NewestHandoff` chain across active missions in the workspace.
2. **M20 Repair (`internal/missionbundle/evidence.go`)**:
   - Validate every objective, run, and claim in `e` resolves to real elements on `bundle`.
3. **M21 Repair (`test/evals/spectacular/`)**:
   - Pin static charter benchmark fixtures to eliminate drift.
   - Build a comprehensive paired matrix in `scope_hardening_bench_test.go`.
