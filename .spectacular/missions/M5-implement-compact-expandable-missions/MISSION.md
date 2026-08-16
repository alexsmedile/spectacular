---
type: Mission
id: 01a009ff-ce94-724e-a6f8-66783f1a4003
ref: M5
title: Implement compact expandable Missions
status: active
owner: Alex
created: "2026-08-16T09:55:54Z"
updated: "2026-08-16T10:15:36Z"

contract:
  ref: Contract:019fe381-5d61-7223-b362-03a5f99a7b10
  fingerprint: sha256:aa2f59e740e9526bacef1dd9999127861836460e5f2f96b5fe05bc86a458ee1a

baseline:
  commit: ddcd153622e5f5bd036ff090fff992f2d0885c64
  branch: codex/lean-launch-context

outcome: Spectacular starts, runs, expands, and completes a Mission through one compact MISSION.md while restoring deterministic CLI support for the new model.
review: independent

completion:
  - claim: mission-model
    pass_boundary: Proposal is optional exploration, Mission is the frozen execution plan, specifications are ordinary Mission work, and Decision is an ADR-like record rather than lifecycle approval.
    proof_requirement: Domain, workspace, Skill, and command tests exercise the same vocabulary and refuse contradictory lifecycle authority.
  - claim: atomic-start
    pass_boundary: One mission start command validates an approved plan, generates stable UUIDv7 identities, creates only MISSION.md with inline Objectives and R1, records activation, and either commits the whole creation or writes nothing.
    proof_requirement: Focused real-process tests cover success, invalid input, retry safety, rollback, compact output, and cold recovery from MISSION.md.
  - claim: selective-expansion
    pass_boundary: Objective promotion and additional Run creation produce dedicated files only when requested, preserve identity, replace inline detail with pointers, and require no Mission index.
    proof_requirement: Tests cover inline state, promotion, second-Run expansion, stable references, and reconstruction through mission show.
  - claim: single-closure
    pass_boundary: One mission complete flow checks frozen criteria, includes required specification edits in ordinary work, and presents one owner gate without a separate Contract reconciliation command.
    proof_requirement: Closure tests cover satisfied, incomplete, stale, and owner-gated outcomes, followed by one compact full verification run.

objectives:
  - ref: O1
    id: 01a009ff-ce94-7249-accc-9c2a089d3080
    outcome: Replace the old noun responsibilities and Mission schema with the compact freeze-and-expand model.
    status: pending
    claims: [mission-model]
  - ref: O2
    id: 01a009ff-ce94-7585-ae70-34adfff1c44b
    outcome: Implement atomic Mission start, Objective promotion, and additional Run creation.
    status: pending
    after: [O1]
    claims: [atomic-start, selective-expansion]
  - ref: O3
    id: 01a009ff-ce94-7c74-a7d8-0f384b429fbc
    outcome: Implement single Mission completion and restore proportionate mechanical validation for the new model.
    status: pending
    after: [O2]
    claims: [single-closure]

run:
  ref: R1
  id: 01a009ff-ce94-7323-a942-754ef422e264
  status: active
  operator: Codex primary session
  started_at: "2026-08-16T09:55:54Z"
  current_objective: O1
  repairs: 0

activation:
  by: Alex
  at: "2026-08-16T09:55:54Z"

validation: manual-bootstrap
mechanical_checks: [yaml-schema, uuidv7-identity, reference-integrity, contract-binding, baseline-binding, completion-claim-coverage, objective-dependency-dag, run-state, authority-vocabulary, file-layout]

authority:
  operator: [inspect, edit-in-scope, choose-reversible-implementation, run-checks, generate-derived-files, bounded-repair, commit-local]
  requires_owner: [change-outcome-or-completion, expand-scope, push, merge, release, irreversible-change, destructive-data, secret-change]

scope:
  mechanical: [cmd/spectacular/, internal/, skills/spectacular/, install/, .spectacular/]
  semantic:
    - Spectacular v2 Mission preparation, execution, expansion, validation, and closure.
    - Proposal, Mission, Contract, Decision, Objective, Run, Evidence, and owner/operator responsibilities.
    - Existing v2 Missions remain readable without migration or a compatibility reader.

repair_budget: 3
dependencies: []
gaps: []

stops:
  - Changing the approved outcome, completion criteria, or observable Mission lifecycle.
  - Weakening atomic creation, stable identity, cold recovery, or owner authority.
  - Requiring migration, a compatibility reader, generic record mutation, or restored ceremony.
---
# Mission

## Origin and product impact

This Mission comes from the owner-confirmed 2026-08-16 design session; no
persistent Proposal file was needed. It implements Spectacular product Contract
version 2. Product specifications are working files in the same change as code,
tests, Skill guidance, and CLI behavior, with no later reconciliation step.

## Execution plan

### O1 — Canonical model

- Make Proposal an optional, mutable exploration artifact that may live in a
  Spectacular file, an issue, or the preparation conversation.
- Freeze the selected outcome, completion criteria, boundaries, inline
  Objectives, and initial Run in `MISSION.md` at start.
- Record activation directly in the Mission; reserve Decision for durable
  ADR-like choices.
- Remove mandatory preparation receipts, lifecycle Decisions, Mission-local
  indexes, and separate Contract reconciliation from the new path.
- Keep existing v2 Missions readable without rewriting them; one current reader
  must understand compact inline records and expanded referenced records.

### O2 — Start and expand

- Add `spectacular mission start` as the atomic typed entry point. A persisted
  Proposal may be supplied, but it is not required and is never created by the
  command.
- Generate UUIDv7 identities, retry identity, baseline, activation, O1…On, and
  R1 mechanically without exposing those details in routine input.
- Initially create only `<mission>/MISSION.md`.
- Add `spectacular objective promote M5/O1`; preserve the Objective UUID and
  replace its inline detail with a file pointer.
- Add `spectacular run start M5 --title <title>`; when R2 is needed, create
  `runs/R1-*.md` and `runs/R2-*.md`, preserve R1 identity, and point the Mission
  at the active Run.

### O3 — Complete and restore validation

- Add `spectacular mission complete M5` as the single closure entry point.
- Check frozen completion criteria, required Evidence, independent review,
  relevant code and specification edits, and freshness before presenting one
  owner gate.
- Replace the old mandatory create/transition/reconcile path rather than
  retaining it as a second public workflow.
- Update the Skill, generated interface, fixtures, and focused tests.
- Run changed-scope tests during each cluster and `bash test/verify.sh all` once
  after integration, with detailed logs only on failure.
- Confirm that the restored CLI can inspect and validate this hand-authored M5
  without modifying it.

## Bootstrap conditions

The current CLI implements the previous Mission model. Until O3 restores the
mechanical checks named in frontmatter, M5 is maintained directly as canonical
Markdown and current CLI validation or mutation output is not completion proof.

During the bootstrap, manual checks own YAML syntax, identity uniqueness,
references, Contract and baseline bindings, completion claim coverage,
Objective dependencies, Run state, authority vocabulary, and file placement.
The bootstrap ends only when the restored CLI validates M5 without rewriting it.

## Independent review

After implementation and the single full verification gate, a fresh reviewer
assesses all four completion claims against their frozen pass boundaries and
proof requirements. The reviewer reports claim results and concrete defects;
it does not redefine success criteria or repeat implementation work.
