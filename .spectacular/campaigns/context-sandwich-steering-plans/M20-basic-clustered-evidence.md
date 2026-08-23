---
type: MissionPlan
title: Record basic clustered Evidence
owner: Alex
contract:
  ref: Contract:019fe381-5d61-7223-b362-03a5f99a7b10
outcome: Attributable Evidence packages cover Runs, Objectives, or clusters without per-run paperwork.
review: independent
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
objectives:
  - outcome: Implement atomic clustered Evidence recording.
    claims: [atomic-clustered-evidence]
  - outcome: Preserve completed-vs-proved lifecycle boundary.
    claims: [completed-is-not-proved]
  - outcome: Version CLI Contract and expose 18-command surface.
    claims: [authorized-surface-18]
authority:
  operator: [inspect, edit-in-scope, choose-reversible-implementation, run-checks, generate-derived-files, bounded-repair, commit-local]
  requires_owner: [activate-mission, change-outcome-or-completion, expand-scope, push, merge, release, irreversible-change, destructive-data]
scope:
  mechanical: [cmd/spectacular, internal/command, internal/missionbundle, skills/spectacular, .spectacular/contracts]
  semantic: [atomic clustered Evidence recording, completed vs proved lifecycle boundary, 18-command CLI surface]
repair_budget: 2
dependencies: [M19 completed with disjoint write reservations and Handoff validation]
gaps: []
stops: [evidence-self-certification, stale-commit-or-tree, command-count-other-than-18, second-proof-dependency-graph, data-loss]
---

# Mission: Record Basic Clustered Evidence

## Purpose & Scope
Introduces `spectacular evidence record` (Command #18) to attach attributable test outputs, logs, and artifacts to individual Runs or Objective clusters.

## Key Deliverables & Actions
1. **Evidence Record Model (`internal/missionbundle/evidence.go`)**: Parse `EvidenceDraft` frontmatter and write to `.spectacular/missions/<slug>/evidence/` with atomic indexing.
2. **CLI Registration (`internal/command/command.go`)**: Register `evidence record <mission-ref> <evidence.md|->` (18 commands). Bump `CC-missioncli` to v6.
3. **Lifecycle Verification (`test/evals/spectacular/`)**: Test `run completed` $\to$ `evidence recorded` $\to$ `review passed` progression.
