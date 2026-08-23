---
type: Mission
id: 01a02cfe-c428-76b7-8e06-f485642347cb
title: Record basic clustered Evidence
status: active
created: "2026-08-23T05:01:35Z"
updated: "2026-08-23T05:01:35Z"
activation:
    at: "2026-08-23T05:01:35Z"
    by: Alex
    fingerprint: sha256:5a9c6de23542574338debd7069bff58fd997d2ff7cf885189177f9c244fae775
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
    branch: codex/m20-basic-clustered-evidence
    commit: 088c828acffc19b958c2c2687d95927eca2800c0
completion:
    - claim: atomic-clustered-evidence
      pass_boundary: evidence record atomically writes an Evidence package covering a Mission, Objective, Run, or cluster with commit/tree coordinates, checks, and limitations.
      proof_requirement: Table-driven tests prove atomic write, retry convergence, and rollback on invalid refs.
    - claim: completed-is-not-proved
      pass_boundary: run transition to completed records execution success; evidence record provides proof, and review record delivers formal verdicts.
      proof_requirement: Lifecycle fixtures verify completed runs with visible pending claims and explicit evidence gating.
    - claim: authorized-surface-18
      pass_boundary: CC-missioncli, catalogs, and CLI registry expose exactly 18 commands through evidence record with rollback safety.
      proof_requirement: Registry assertions verify exactly 18 commands; git diff proves CC-projsurf remains unchanged.
contract:
    fingerprint: sha256:aa2f59e740e9526bacef1dd9999127861836460e5f2f96b5fe05bc86a458ee1a
    ref: Contract:019fe381-5d61-7223-b362-03a5f99a7b10
dependencies:
    - M19 completed with disjoint write reservations and Handoff validation
gaps: []
objectives:
    - claims:
        - atomic-clustered-evidence
      id: 01a02cfe-c428-7603-a188-62c3b564f2b8
      outcome: Implement atomic clustered Evidence recording.
      ref: O1
      status: pending
    - claims:
        - completed-is-not-proved
      id: 01a02cfe-c428-7eb1-880a-d4a069ad9792
      outcome: Preserve completed-vs-proved lifecycle boundary.
      ref: O2
      status: pending
    - claims:
        - authorized-surface-18
      id: 01a02cfe-c428-741a-8632-d6e7451ea82e
      outcome: Version CLI Contract and expose 18-command surface.
      ref: O3
      status: pending
outcome: Attributable Evidence packages cover Runs, Objectives, or clusters without per-run paperwork.
owner: Alex
ref: M20
repair_budget: 2
review: independent
run:
    current_objective: O1
    id: 01a02cfe-c428-7777-9213-a172688bfb87
    operator: Alex
    ref: R1
    repairs: 0
    started_at: "2026-08-23T05:01:35Z"
    status: active
scope:
    mechanical:
        - cmd/spectacular
        - internal/command
        - internal/missionbundle
        - skills/spectacular
        - .spectacular/contracts
    semantic:
        - atomic clustered Evidence recording
        - completed vs proved lifecycle boundary
        - 18-command CLI surface
start_key: sha256:f99c8245795fe7c5427469e797e8f3e27985f80fa5686359f824c7df594ecf4b
stops:
    - evidence-self-certification
    - stale-commit-or-tree
    - command-count-other-than-18
    - second-proof-dependency-graph
    - data-loss
validation:
    mode: cli
    schema: mission.v2
---
# Mission: Record Basic Clustered Evidence

## Purpose & Scope
Introduces `spectacular evidence record` (Command #18) to attach attributable test outputs, logs, and artifacts to individual Runs or Objective clusters.

## Key Deliverables & Actions
1. **Evidence Record Model (`internal/missionbundle/evidence.go`)**: Parse `EvidenceDraft` frontmatter and write to `.spectacular/missions/<slug>/evidence/` with atomic indexing.
2. **CLI Registration (`internal/command/command.go`)**: Register `evidence record <mission-ref> <evidence.md|->` (18 commands). Bump `CC-missioncli` to v6.
3. **Lifecycle Verification (`test/evals/spectacular/`)**: Test `run completed` $\to$ `evidence recorded` $\to$ `review passed` progression.
