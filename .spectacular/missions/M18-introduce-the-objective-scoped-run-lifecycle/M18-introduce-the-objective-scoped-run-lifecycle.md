---
type: Mission
id: 01a02ba9-cf00-7f30-a440-60e4d88ec2dc
title: Introduce the Objective-scoped Run lifecycle
status: active
created: "2026-08-22T22:52:28Z"
updated: "2026-08-22T22:56:18Z"
activation:
    at: "2026-08-22T22:52:28Z"
    by: Alex
    fingerprint: sha256:e7c06db2f40c66f6fafe983c9433e482424a22cf40b8cc6afe6e0046598d9425
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
baseline:
    branch: codex/m18-objective-run-lifecycle
    commit: a44c5ea424e2e6aec244f54c907f4f064e996f2b
completion:
    - claim: objective-scoped-runs
      pass_boundary: Every new Run belongs to exactly one Objective, an Objective retains serial attempts, and spectacular run start accepts <mission-ref>/<objective-ref> while refusing a second active reserving Run for that Objective.
      proof_requirement: Model, schema, and CLI tests cover objective ownership, serial attempt history, reservation conflict refusal, and parity between inline and promoted representations.
    - claim: attributable-run-transitions
      pass_boundary: spectacular run transition validates legal state transitions (active to paused, blocked, awaiting-review, completed, stopped), recording actor, reason, and optional next action with all-or-nothing rollback on failure.
      proof_requirement: Table-driven transition tests cover legal/illegal state paths, reservation release on terminal states, attribution persistence, and zero mutation on error.
    - claim: historical-decoder-integrity
      pass_boundary: Completed historical Missions (M1-M17) remain byte-identical, valid, and inspectable under their recorded Contract fingerprints without requiring file rewrites.
      proof_requirement: Golden tests assert backward-compatible decoding of legacy single-run Missions across mission show and mission check.
    - claim: authorized-surface-17
      pass_boundary: The CLI registry, mechanical interface, generated schemas, and CC-missioncli contract accurately reflect exactly 17 registered commands through the addition of run transition.
      proof_requirement: Registry-derived tests assert exact 17-command count, and CC-missioncli is version-bumped (v4 -> v5).
contract:
    fingerprint: sha256:aa2f59e740e9526bacef1dd9999127861836460e5f2f96b5fe05bc86a458ee1a
    ref: Contract:019fe381-5d61-7223-b362-03a5f99a7b10
dependencies:
    - M17 completed with the exact 16-command surface
gaps: []
objectives:
    - claims:
        - objective-scoped-runs
        - historical-decoder-integrity
      id: 01a02ba9-cf00-74e1-8fd0-cce49c08e1c9
      outcome: Implement Objective-scoped Run ownership and historical backward compatibility (Objective Scope Pillar).
      ref: O1
      status: implemented
    - claims:
        - attributable-run-transitions
      id: 01a02ba9-cf00-754c-a618-ae38e497fd2f
      outcome: Implement attributable Run transitions and reservation state semantics (Transitions Pillar).
      ref: O2
      status: implemented
    - claims:
        - authorized-surface-17
      id: 01a02ba9-cf00-77ad-abe5-b5b616de1bbd
      outcome: Version and verify the exact 17-command public surface (Surface Growth Pillar).
      ref: O3
      status: implemented
outcome: Unlock resilient multi-attempt execution by scoping Runs to individual Objectives and exposing spectacular run transition while growing the public surface from 16 to 17 commands.
owner: Alex
ref: M18
repair_budget: 2
review: independent
run:
    current_objective: O1
    id: 01a02ba9-cf00-78cd-92b8-013767d2695e
    operator: Alex
    ref: R1
    repairs: 0
    started_at: "2026-08-22T22:52:28Z"
    status: active
scope:
    mechanical:
        - cmd/spectacular/
        - internal/command/
        - internal/missionbundle/
        - generated/
        - .spectacular/contracts/
        - test/
    semantic:
        - Objective-scoped Run identity and lifecycle
        - Attributable run transition state machine
        - Historical decoder compatibility
        - Mechanical interface surface growth from 16 to 17 commands
start_key: sha256:f828830909be43729954640fd8f30fd00d7e3896ecf77c118473242bfa5db4aa
stops:
    - completed-mission-rewrite
    - ambiguous-run-ownership
    - implicit-error-recovery
    - command-count-other-than-17
    - transition-atomicity-failure
validation:
    mode: cli
    schema: mission.v2
---
# Mission: Introduce the Objective-scoped Run Lifecycle

## User Superpower (The Hub)
Enables resilient multi-attempt execution directly at the Objective level. If a subagent hits a blocker or exhausts repairs on an Objective, the operator can transition that run to `paused` or `stopped` with clear attribution, then launch a fresh serial attempt (`R2`, `R3`) under that exact Objective without resetting or corrupting the rest of the Mission.

## Technical Pillars (The Spokes)
1. **Objective Scope Pillar**: Moves the primary Run anchor from the whole Mission to individual Objectives (`M/O/R`), retaining attempt history.
2. **Transitions Pillar (`spectacular run transition`)**: Enforces the formal state machine (`active`, `paused`, `blocked`, `awaiting-review`, `completed`, `stopped`) with mandatory `--by` and `--reason`.
3. **Historical Compatibility Pillar**: Ensures all earlier completed Missions continue to decode through the unified bundle model with zero drift.
4. **Surface Growth Pillar (16 $\to$ 17)**: Bumps `CC-missioncli` (v4 $\to$ v5), updates the command registry, and regenerates interface catalogs.

## Key Deliverables & Actions

### 1. Objective Run Hierarchy & Decoder (`internal/missionbundle/`)
- Update `Objective` model to hold `Runs []Run` and `Run *Run`.
- Update `decode.go` to handle both new Objective-level runs and legacy Mission-level runs seamlessly.
- Update `service.StartRun` to accept `<mission-ref>/<objective-ref>` and enforce single-active reservation per Objective.

### 2. Attributable Run Transition (`internal/missionbundle/run.go`, `internal/command/`)
- Implement `service.TransitionRun(targetRef, toState, actor, reason, nextAction)`.
- Validate transition state machine graph and enforce attribution fields.
- Update Objective readiness dynamically upon Run completion/stop.

### 3. Public `run transition` Command (`internal/command/`, `cmd/spectacular/`)
- Wire `spectacular run transition <mission-ref>/<objective-ref>/<run-ref> --to <state> --by <actor> --reason <text> [--next-action <action>] [--json]`.
- Update `spectacular run start` argument shape to `<mission-ref>/<objective-ref> [--title <title>]`.

### 4. Surface & Contract Reconcile
- Bump [`CC-missioncli`](.spectacular/contracts/CC-missioncli-spectacular-mechanical-cli.md) from `v4` to `v5` with 17 commands.
- Run `go run ./cmd/generate-interface` to update `generated/mechanical-interface.json` and `catalog.md`.
- Run full verification test suite.
