---
type: MissionPlan
title: Expose charter and atomic Decision recording
owner: Alex
contract:
  ref: Contract:019fe381-5d61-7223-b362-03a5f99a7b10
outcome: An Orchestrator can invoke the proved charter compiler and atomically record a complete owner Decision while the generated public interface grows only from 14 to 16 commands.
review: independent
completion:
  - claim: charter-public-receipt
    pass_boundary: spectacular charter accepts one Mission/Objective ref and explicit additional source refs, returns the proved three-layer charter plus deterministic attribution and size disposition, and remains read-only.
    proof_requirement: Real-process tests cover help, human and JSON receipts, exact source order, all token boundaries, malformed refs, inactive or ambiguous targets, retry identity, and byte-identical canonical state.
  - claim: atomic-decision-recording
    pass_boundary: spectacular decide accepts a complete Orchestrator-authored Decision from a file or stdin, assigns valid identity and ref, preserves exact owner wording and digest, validates targets and optional supersession, and atomically updates the Decision plus Decision, catalog, and root indexes.
    proof_requirement: Fault injection at every write boundary proves all-or-nothing recovery and retry convergence; tests cover collisions, traversal, target and lineage refusal, wording fidelity, index rebuild equivalence, and no Mission, Objective, or Run mutation.
  - claim: explicit-eligibility-reporting
    pass_boundary: decide reports newly eligible work only where a frozen blocker explicitly names the recorded Decision; all other work is left for Orchestrator judgment and no Run state changes.
    proof_requirement: Positive and negative fixtures prove exact blocker matching, no meaning-based inference, no implicit dispatch, and byte-identical execution state.
  - claim: authorized-surface-16
    pass_boundary: The registry, versioned CC-missioncli Contract, help, schemas, generated interface, installer data, and forbidden-command policy expose exactly 16 commands through only charter and decide.
    proof_requirement: Registry-derived tests fail on any missing, extra, stale, or misclassified command; preflight, quick, acceptance, and final all gates pass before completion.
objectives:
  - outcome: Expose the proved charter engine with compact attributable receipts.
    claims: [charter-public-receipt]
  - outcome: Implement atomic Decision and generated-index recording without runtime mutation.
    claims: [atomic-decision-recording, explicit-eligibility-reporting]
  - outcome: Version and verify the exact 16-command public surface.
    claims: [authorized-surface-16]
authority:
  operator: [inspect, edit-in-scope, choose-reversible-implementation, run-checks, generate-derived-files, bounded-repair, commit-local]
  requires_owner: [activate-mission, change-outcome-or-completion, expand-scope, push, merge, release, irreversible-change]
scope:
  mechanical: [cmd/spectacular, internal/command, internal/runtime, internal/governance, internal/humanlayout, skills/spectacular, install, test, .spectacular/contracts]
  semantic: [public charter retrieval, atomic durable Decision recording, public CLI growth from 14 to 16]
repair_budget: 2
dependencies: [M16 Evidence and independent Review pass the 40-percent and zero-regression gates]
gaps: []
stops: [benchmark-gate-not-passed, command-count-other-than-16, non-atomic-index-write, semantic-eligibility-inference, execution-state-mutation, contract-conflict]
---

# Mission

> **Future Mission sketch.** Preserve as design input. Do not activate, maintain,
> validate, or review as a final plan until its predecessor closes and the
> Orchestrator re-prepares this block from current Evidence.

This Mission exposes only behavior already proved internally and adds Decision
recording as a separate atomic path. It does not change the Run hierarchy,
dispatch agents, or add Evidence mechanics.
