---
type: Decision
id: 01a02a64-3e86-7eee-ad75-04fc4beb25f3
title: Scope Runs to Objectives and govern concurrency by reservations
created_by: Alex
created: "2026-08-22T16:54:04Z"
updated: "2026-08-22T16:54:04Z"
actor: Alex
actor_role: owner
ref: D17-objective-scoped-runs-and-concurrency
question: How should Missions, Objectives, Runs, Handoffs, and external agents relate when work is concurrent?
disposition: objective-scoped-runs-and-reservation-gated-concurrency
rationale: >-
    The canonical Mission, Objective, Run, Contract, and Handoff model already expresses the
    required hierarchy; agents are external execution destinations rather than new governance
    objects. Runs are attempts on Objectives, so serial recovery belongs beneath one Objective
    while independent Objectives and Missions may proceed concurrently. Numeric ceilings are a
    poor proxy for safety: explicit dependencies and non-overlapping writable reservations provide
    the real boundary.
alternatives:
    - one live Run for an entire Mission, which prevents safe Objective concurrency
    - a new canonical Agent or Worker record, which duplicates host-runtime state
    - a fixed concurrency ceiling independent of dependencies and writable scope
authority_basis: Owner confirmed Objective-scoped Runs, concurrent Missions and Objectives, no hard concurrency ceiling, frozen Handoff assignments, explicit Run transitions, clustered Evidence, and the 14-to-18 public command surface during the P11 design interview.
authorized_effects:
    - command.surface-growth
    - contract.version-bump
conditions:
    - command-surface-growth-14-to-18
    - run-start-becomes-objective-scoped
    - no-overlapping-active-write-reservations
scope:
    - v2
targets:
    - Proposal:01a029be-b7d3-703c-a7ee-50c6b8bae3a2
supersedes: D12-isolation-and-context-compilation
---

# Scope Runs to Objectives and govern concurrency by reservations

## Decision

- One Orchestrator session is bounded to one Mission; the repository may contain concurrent active Missions in separate integration worktrees.
- A Mission owns Objectives. Each Run belongs to exactly one Objective, and an Objective may retain several serial Run attempts.
- At most one reserving Run exists for an Objective. Runs on distinct eligible Objectives may execute concurrently without a numeric ceiling.
- `active`, `paused`, `blocked`, and `awaiting-review` reserve the Handoff's writable perimeter. `completed` and `stopped` release it.
- A Handoff names its Mission, Objective, Run, selected sources, and machine-readable `writes:` paths. Exact files and trailing-slash directory subtrees are supported; globs and parent traversal are not.
- Concurrent dispatch refuses overlapping active Handoff perimeters, including across Missions. Potential future Mission overlap does not prevent activation.
- An external agent or host task owns no canonical authority or status. Spectacular preserves the Run, frozen Handoff, and returned Evidence.
- Evidence may cover one Run, Objective, Objective cluster, or final Mission gate and links to every Mission, Objective, Run, Handoff, and claim it answers. Changed instructions require a new Handoff.
- A completed Run means execution ended successfully; it does not claim that Evidence or Review is complete.
- `objective finish` requires no reserving Run and at least one completed Run. Its claims may remain unproven until a later clustered or final gate.
- `after:` waits for implemented work. Optional `after_proof:` waits for the Evidence or Review level frozen by the Mission plan.

## Mechanical interface authorization

- Add `spectacular charter`, `spectacular decide`, `spectacular run transition`, and `spectacular evidence record`, growing the generated public surface from 14 to 18 commands.
- Change `spectacular run start` to accept `<mission-ref>/<objective-ref>` and refuse while that Objective has a reserving Run.
- `run transition` records the actor, target state, reason, and next action; it never chooses a recovery from an observed error.
- `evidence record` atomically preserves a complete Evidence package without judging its own sufficiency or requiring one package per Run or Objective.
