---
type: MissionPlan
title: Record basic clustered Evidence
owner: Alex
contract:
  ref: Contract:019fe381-5d61-7223-b362-03a5f99a7b10
outcome: One attributable Evidence package can efficiently cover a Run, Objective, Objective cluster, or final Mission gate without creating a second dependency system or requiring per-Run paperwork.
review: independent
completion:
  - claim: atomic-clustered-evidence
    pass_boundary: evidence record atomically accepts a complete package naming every covered Mission, Objective, Run, Handoff, and claim, and supports Run, Objective, cluster, and final-Mission coverage without requiring one package per Run.
    proof_requirement: Real-process and fault-injection tests cover every coverage shape, partial and duplicate coverage, malformed and foreign refs, retry convergence, every write boundary, and no unrelated Mission, Objective, or Run mutation.
  - claim: attributable-evidence-integrity
    pass_boundary: Every package records actor, executor, method, observed_at, commit and tree, declared checks and results, contrary evidence, and limitations; freshness and reference integrity are validated while Evidence never certifies its own broader sufficiency.
    proof_requirement: Table-driven tests cover every required field, stale commit/tree, missing and failed checks, incomplete claim coverage, contradictory observations, absent limitations, activation-fingerprint mismatch at Review, and exact refusal-before-write behavior.
  - claim: completed-is-not-proved
    pass_boundary: A completed Run releases its reservation and reports execution success only; objective finish requires no reserving Run and at least one completed Run while leaving unproved claims visible. Proof-sensitive downstream dispatch remains an explicit condition in its frozen Handoff and an Orchestrator decision, not a second dependency graph.
    proof_requirement: Lifecycle fixtures show clustered and final proof, visible pending claims, manual refusal of a proof-sensitive dispatch before accepted Evidence or Review, later dispatch after acceptance, and no after_proof schema or automatic proof scheduler.
  - claim: authorized-surface-18
    pass_boundary: CC-missioncli, the registry, help, schemas, generated interface, installer data, and forbidden-command policy expose exactly 18 commands through only evidence record; CC-projsurf and its concurrent-run-timelines Gap remain byte-identical.
    proof_requirement: Registry-derived tests fail on any missing, extra, stale, or misclassified command; git diff proves CC-projsurf and completed Mission bindings unchanged; full Mission verification passes.
objectives:
  - outcome: Implement atomic clustered Evidence with complete coverage and attribution.
    claims: [atomic-clustered-evidence, attributable-evidence-integrity]
  - outcome: Preserve the distinction between completed execution and accepted proof without adding another dependency graph.
    claims: [completed-is-not-proved]
  - outcome: Version only the CLI Contract and align the exact 18-command public surface.
    claims: [authorized-surface-18]
authority:
  operator: [inspect, edit-in-scope, choose-reversible-implementation, run-checks, generate-derived-files, bounded-repair, commit-local]
  requires_owner: [activate-mission, change-outcome-or-completion, expand-scope, push, merge, release, irreversible-change]
scope:
  mechanical: [cmd/spectacular, internal/command, internal/missionbundle, internal/runtime, skills/spectacular, install, test, .spectacular/contracts]
  semantic: [atomic clustered Evidence coverage and attribution, evidence freshness and integrity, visible pending proof, manual Orchestrator proof gates, public CLI growth from 17 to 18]
repair_budget: 2
dependencies: [M19 completed with frozen Handoff reservations and versioned CC-v2prod]
gaps: []
stops: [evidence-self-certification, stale-commit-or-tree, incomplete-claim-coverage, automatic-proof-scheduler, second-proof-dependency-graph, command-count-other-than-18, contract-conflict]
---

# Mission

> **Future Mission sketch.** Preserve as design input. Do not activate, maintain,
> validate, or review as a final plan until its predecessor closes and the
> Orchestrator re-prepares this block from current Evidence.

Keep Evidence useful and cheap. Small Runs may continue toward clustered proof.
When downstream work must wait, its frozen Handoff names the required Evidence or
Review and the Orchestrator dispatches only after acceptance. Do not add
`after_proof:`, concurrent timeline projections, or automatic proof propagation.
