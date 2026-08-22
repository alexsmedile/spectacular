---
type: MissionPlan
title: Introduce the Objective-scoped Run lifecycle
owner: Alex
contract:
  ref: Contract:019fe381-5d61-7223-b362-03a5f99a7b10
outcome: New Missions retain serial Run attempts beneath each Objective and transition them explicitly while completed historical Missions remain readable under their frozen bindings.
review: independent
completion:
  - claim: objective-scoped-runs
    pass_boundary: Every new Run belongs to exactly one Objective, an Objective retains serial attempts, and run start accepts mission-ref/objective-ref while refusing a second reserving Run for that Objective.
    proof_requirement: Model, schema, and real-process tests cover identity, promotion, retry, objective ownership, reserving conflicts, multiple completed attempts, inline and promoted parity, and atomic refusal-before-write.
  - claim: attributable-run-transitions
    pass_boundary: A new Run starts active; active may move to paused, blocked, awaiting-review, completed, or stopped; paused may move to active, blocked, or stopped; blocked may move to active or stopped; awaiting-review may move to active, completed, or stopped; completed and stopped are terminal. Every transition records actor, reason, and next action without selecting recovery from an error.
    proof_requirement: Table-driven legal and illegal transition tests, fault injection, retry convergence, and concurrent mutation tests assert reservation retention and release, complete attribution, and byte-identical refusal.
  - claim: historical-decoder-integrity
    pass_boundary: Every completed v2 Mission remains byte-identical, inspectable, and valid under its recorded Contract fingerprint while new semantics apply only to new Missions bound to the versioned product Contract.
    proof_requirement: Golden fixtures include all completed Mission shapes and old single-Run records; show and check decode them without rewrite, and git diff proves no completed Mission file changed.
  - claim: authorized-surface-17
    pass_boundary: The generated surface grows from exactly 16 to exactly 17 through only run transition, while run start changes signature and all interfaces and Contracts agree.
    proof_requirement: Registry and generated-interface tests assert the exact set, signature, effect class, refusals, and forbidden commands; full Mission verification passes.
objectives:
  - outcome: Version the CLI and project-surface Run model and implement Objective ownership with historical decoding.
    claims: [objective-scoped-runs, historical-decoder-integrity]
  - outcome: Implement explicit transitions and reservation state semantics.
    claims: [attributable-run-transitions]
  - outcome: Align the exact 17-command mechanical interface.
    claims: [authorized-surface-17]
authority:
  operator: [inspect, edit-in-scope, choose-reversible-implementation, run-checks, generate-derived-files, bounded-repair, commit-local]
  requires_owner: [activate-mission, change-outcome-or-completion, expand-scope, push, merge, release, irreversible-change]
scope:
  mechanical: [cmd/spectacular, internal/command, internal/missionbundle, internal/runtime, skills/spectacular, install, test, .spectacular/contracts]
  semantic: [Objective-scoped Run identity and lifecycle, explicit transition attribution, CLI and project-surface Contract versioning, historical compatibility, public CLI growth from 16 to 17]
repair_budget: 2
dependencies: [M17 completed with the exact 16-command surface]
gaps: []
stops: [completed-mission-rewrite, ambiguous-run-ownership, implicit-error-recovery, command-count-other-than-17, transition-atomicity-failure]
---

# Mission

> **Future Mission sketch.** Preserve as design input. Do not activate, maintain,
> validate, or review as a final plan until its predecessor closes and the
> Orchestrator re-prepares this block from current Evidence.

This Mission changes lifecycle shape but does not dispatch external agents, enforce
cross-Mission write perimeters, or record clustered Evidence.
