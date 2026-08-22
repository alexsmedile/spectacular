---
type: Decision
id: 01a02afd-9607-7161-8cf9-4e5d7d13bb1f
title: Preserve proof gates while staging the context-sandwich command surface
created_by: Alex
created: "2026-08-22T19:40:57Z"
updated: "2026-08-22T19:40:57Z"
actor: Alex
actor_role: owner
ref: D21-context-sandwich-execution-gates
question: How should the Objective-scoped Run model retain D12's proof conditions and stage the authorized public command growth?
disposition: preserve-proof-gates-and-stage-four-authorized-commands
rationale: >-
    D17 correctly moved Runs beneath Objectives and replaced a numeric concurrency ceiling with
    dependency and writable-reservation checks, but its whole-Decision supersession omitted D12's
    minimum context-reduction and zero-regression conditions. The four authorized commands also
    cross different proof boundaries, so exposing them in smaller reviewed stages reduces the
    amount of behavior an owner must accept at once without changing the final 18-command surface.
alternatives:
    - retain D17 as written and rely on Proposal prose for the omitted proof gates
    - expose all four commands in one Mission after one combined review
    - restore a fixed concurrency ceiling as a substitute for dependency and write-scope checks
authority_basis: Owner approved the Objective-scoped Run hierarchy, reservation-gated concurrency, clustered Evidence, the final 18-command surface, the 1200/1400/1440 charter boundaries, and a finer Mission split so each update can be evaluated independently.
authorized_effects:
    - command.surface-growth
    - contract.version-bump
conditions:
    - context-ingestion-reduction-at-least-40-percent
    - zero-regression-in-safety-success-recovery-routing-and-decision-fidelity
    - command-surface-growth-staged-14-to-16-to-17-to-18
    - no-overlapping-active-write-reservations
scope:
    - v2
targets:
    - Proposal:01a029be-b7d3-703c-a7ee-50c6b8bae3a2
supersedes: D17-objective-scoped-runs-and-concurrency
---

# Preserve proof gates while staging the context-sandwich command surface

## Canonical hierarchy and concurrency

- One Orchestrator session is bounded to one Mission. Multiple Missions may be active in separate integration worktrees.
- A Mission owns Objectives. Every Run belongs to exactly one Objective, and an Objective may retain multiple serial Run attempts.
- At most one reserving Run exists for an Objective. Distinct eligible Objectives may run concurrently without a numeric ceiling.
- Eligibility follows explicit dependencies and disjoint writable reservations, including reservations held by other active Missions.
- `active`, `paused`, `blocked`, and `awaiting-review` retain the frozen Handoff's reservation. `completed` and `stopped` release it.
- External agents are execution destinations, not canonical governance objects. Spectacular preserves the assignment and returned Evidence.

## Handoff, proof, and history

- A frozen Handoff names Mission, Objective, Run, selected sources, and machine-readable `writes:` entries. Exact files and trailing-slash directory subtrees are valid; globs and parent traversal are not.
- Changed instructions require a new Handoff. Concurrent dispatch refuses overlapping reserving perimeters.
- Evidence may cover a Run, Objective, Objective cluster, or final Mission gate and names every claim and governed record it answers.
- A completed Run reports successful execution, not accepted proof. `objective finish` requires no reserving Run and at least one completed Run; proof may remain visibly pending for a clustered gate.
- `after:` waits for implementation. When downstream work must wait for proof, the Orchestrator records that condition in the frozen Handoff and dispatches only after the named Evidence or Review is accepted; P11 adds no second proof-dependency graph.
- Completed Missions remain historical plans bound to the Contract version they activated against. New semantics apply only through new Contract versions and new Missions.

## Proof conditions

- Before public charter mechanics ship, the paired M14 benchmark must demonstrate at least 40 percent lower total context ingestion than the full-scan baseline.
- That reduction cannot trade away task success, safety, recovery, routing, or decision fidelity. Any regression or contaminated or missing telemetry fails or returns inconclusive; it never passes by absence.
- Measurements report governance-envelope, named-source, and repair-diagnostic tokens separately.
- Evidence-level proof requires declared claim coverage, a fresh commit/tree binding, and every declared check passing.
- Review-level proof additionally validates the activation fingerprint, commit/tree freshness, claim coverage, verdict, and reviewer independence.

## Staged mechanical interface

- Stage one grows the generated surface from exactly 14 to exactly 16 commands through only `spectacular charter` and `spectacular decide`, after charter benchmark proof.
- Stage two grows 16 to 17 through `spectacular run transition` and changes `run start` to accept `<mission-ref>/<objective-ref>`, after Objective-scoped lifecycle and historical-decoder proof.
- Stage three grows 17 to 18 through `spectacular evidence record`, after clustered-Evidence integrity and atomicity validation.
- `run transition` records actor, target state, reason, and next action; it never chooses recovery from an observed error.
- `evidence record` atomically preserves a complete Evidence package without requiring one package per Run or Objective and without judging its own sufficiency.
