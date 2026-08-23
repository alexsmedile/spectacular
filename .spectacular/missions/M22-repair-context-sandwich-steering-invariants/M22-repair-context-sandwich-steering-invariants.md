---
type: Mission
id: 01a02e74-f908-7f50-b05b-7127633d6595
title: Repair Context Sandwich steering invariants
status: active
created: "2026-08-23T11:50:27Z"
updated: "2026-08-23T11:53:18Z"
activation:
    at: "2026-08-23T11:50:27Z"
    by: Alex
    fingerprint: sha256:8e66eff6a64465a5bc61dfcd9bd9d540f65b030fba5446cf41e4c5ce7f92d2ea
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
    branch: codex/m22-repair-context-sandwich-invariants
    commit: 8fd81afc1239fbdba1a8c2253b434ae6ffd5d79d
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
contract:
    fingerprint: sha256:cb380518f027c697e5d2a22f5a4c6ca5f2cab8e996db5718b3ef8b4cbf72c1d4
    ref: Contract:01a00a20-63dd-7670-97f1-9eb8e12adc3a
dependencies:
    - M21 completed with campaign baseline
gaps: []
objectives:
    - claims:
        - m19-reservation-and-git-repair
      id: 01a02e74-f908-7b42-b5b3-420a425829b9
      outcome: Repair Handoff reservation chains, Decision unblocks, and linked-worktree Git safety.
      ref: O1
      status: implemented
    - claims:
        - m20-evidence-target-integrity
      id: 01a02e74-f908-7558-bedb-f46320f9812e
      outcome: Enforce strict target reference validation on Evidence packages.
      ref: O2
      status: implemented
    - claims:
        - m21-pinned-benchmark-matrix
      id: 01a02e74-f908-76c5-b0df-caf6e4a94fe6
      outcome: Anchor scope hardening and context benchmarks in pinned immutable fixtures.
      ref: O3
      status: implemented
outcome: Repair proof and behavior gaps across Handoff reservations, Decision unblocks, worktrees, Evidence targets, and pinned fixtures.
owner: Alex
ref: M22
repair_budget: 2
review: independent
run:
    current_objective: O1
    id: 01a02e74-f908-7977-ba97-7ebf55ebb89d
    operator: Alex
    ref: R1
    repairs: 0
    started_at: "2026-08-23T11:50:27Z"
    status: active
scope:
    mechanical:
        - internal/missionbundle
        - internal/runtime
        - internal/command
        - test/evals/spectacular
        - test
    semantic:
        - linked worktree git inspection
        - Decision-based dependency unblocking
        - strict Evidence target validation
        - pinned benchmark fixtures
start_key: sha256:0a2b918ccbf0fd4d841bac6ea84e9733514a862bf3391cac05af0375e7d38a96
stops:
    - second-dependency-graph
    - live-file-benchmark-drift
    - unvalidated-evidence-targets
    - silent-worktree-git-skip
    - data-loss
validation:
    mode: cli
    schema: mission.v2
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
