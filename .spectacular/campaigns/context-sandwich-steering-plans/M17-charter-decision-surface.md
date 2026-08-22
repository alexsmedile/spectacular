---
type: MissionPlan
title: Expose charter and atomic Decision recording
owner: Alex
contract:
  ref: Contract:019fe381-5d61-7223-b362-03a5f99a7b10
outcome: Unlock instant CLI prompt generation and zero-ceremony atomic Decision recording by exposing spectacular charter and spectacular decide while growing the public surface from 14 to 16 commands.
review: independent
completion:
  - claim: charter-public-receipt
    pass_boundary: spectacular charter accepts a Mission/Objective ref and optional source refs, returning the 3-layer Context Sandwich charter with token count receipt and threshold status as read-only Markdown or JSON.
    proof_requirement: CLI tests cover human/JSON output, token disposition reporting, missing refs, and verify zero canonical files are modified.

  - claim: atomic-decision-recording
    pass_boundary: spectacular decide accepts a decision draft from file or stdin, validates invariants, generates UUIDv7 identity and D<N> ref, writes the file, and updates all indexes atomically with rollback on error.
    proof_requirement: Fault-injection tests at write boundaries prove all-or-nothing rollback; tests verify UUIDv7 generation, D<N> ref allocation, and index synchronization.

  - claim: explicit-eligibility-reporting
    pass_boundary: decide reports newly unblocked work only when an Objective explicitly names the recorded Decision ref as its blocker, with zero meaning-based inference or automatic execution state mutation.
    proof_requirement: Positive and negative test fixtures prove exact blocker matching without mutating Run or Objective states.

  - claim: authorized-surface-16
    pass_boundary: The CLI registry, mechanical interface, generated schemas, and CC-missioncli contract accurately reflect exactly 16 registered commands without extraneous verbs or aliases.
    proof_requirement: Registry-derived tests assert exact 16-command count, and CC-missioncli is version-bumped (v3 -> v4).

objectives:
  - outcome: Expose the read-only charter command with budget receipts (Charter Pillar).
    claims: [charter-public-receipt]
  - outcome: Implement atomic Decision recording and explicit blocker reporting (Decide Pillar).
    claims: [atomic-decision-recording, explicit-eligibility-reporting]
  - outcome: Version and verify the exact 16-command public surface (Surface Growth Pillar).
    claims: [authorized-surface-16]

authority:
  operator: [inspect, edit-in-scope, choose-reversible-implementation, run-checks, generate-derived-files, bounded-repair, commit-local]
  requires_owner: [activate-mission, change-outcome-or-completion, expand-scope, push, merge, release, irreversible-change]
scope:
  mechanical:
    - cmd/spectacular/
    - internal/command/
    - internal/governance/
    - internal/charter/
    - generated/
    - .spectacular/contracts/
    - test/
  semantic:
    - Public spectacular charter command
    - Public spectacular decide atomic transaction command
    - Mechanical interface surface growth from 14 to 16 commands
repair_budget: 2
dependencies: [M16 completed with independent review]
gaps: []
stops:
  - command-count-other-than-16
  - non-atomic-index-write
  - semantic-eligibility-inference
  - execution-state-mutation
  - contract-conflict
---

# Mission: Expose Charter and Atomic Decision Recording

## User Superpower (The Hub)
Eliminates all prompt assembly friction and multi-file governance indexing overhead. Operators can run `spectacular charter` to get instant, razor-sharp worker prompts, and `spectacular decide -` to record durable architectural choices in a single atomic transaction without manually maintaining index tables or UUIDs.

## Technical Pillars (The Spokes)
1. **Charter Pillar (`spectacular charter`)**: Invokes `internal/charter` compiler to emit 3-layer prompt envelopes directly to stdout or structured JSON.
2. **Decide Pillar (`spectacular decide`)**: Validates Decision payload from file/stdin, assigns UUIDv7 & `D<N>`, writes the document, and refreshes all indexes atomically.
3. **Surface Growth Pillar (14 $\to$ 16)**: Bumps `CC-missioncli` (v3 $\to$ v4), updates the command registry, and regenerates interface metadata for exactly 16 commands.

## Key Deliverables & Actions

### 1. Public `charter` Command (`internal/command/`, `cmd/spectacular/`)
- Wire `spectacular charter <mission-ref>/<objective-ref> [extra-sources...] [--json]`.
- Output human-readable Markdown by default, or structured JSON with token count and disposition receipt.

### 2. Public `decide` Command (`internal/command/`, `internal/governance/`)
- Wire `spectacular decide <decision.md|-> [--json]`.
- Parse Decision draft (title, disposition, rationale, scope, targets, supersedes).
- Allocate next canonical `D<N>` reference and UUIDv7 identity.
- Execute atomic transaction: write `.spectacular/decisions/D<N>-slug.md` + update `.spectacular/decisions/index.md` + update `.spectacular/catalog.md` + update `.spectacular/index.md`.
- Check if any Objective in the workspace is explicitly blocked on this ref and report it.

### 3. Surface & Contract Reconcile
- Bump [`CC-missioncli`](.spectacular/contracts/CC-missioncli-spectacular-mechanical-cli.md) from `v3` to `v4` with 16 commands.
- Run `go run ./cmd/generate-interface` to update `generated/mechanical-interface.json` and `catalog.md`.
- Run full verification suite.

