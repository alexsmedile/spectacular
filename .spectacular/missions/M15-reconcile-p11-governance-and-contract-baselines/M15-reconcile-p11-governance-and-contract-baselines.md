---
type: Mission
id: 01a02b83-14f0-7fde-bdb3-1154b0b7c188
title: Reconcile P11 governance and Contract baselines
status: completed
created: "2026-08-22T22:26:03Z"
updated: "2026-08-22T22:26:49Z"
activation:
    at: "2026-08-22T22:26:03Z"
    by: Alex
    fingerprint: sha256:a184517332b911699542256dbc6b378678f58556eeb32cca4e596d540139092c
authority:
    operator:
        - inspect
        - edit-in-scope
        - generate-derived-files
        - run-checks
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
    branch: codex/m15-governance-contract-baseline
    commit: eff9e7c8585db8d9f21631ee7a55c219e17e9d7d
completion:
    - claim: decision-lineage-and-proof-gates-aligned
      pass_boundary: D21 is the active authority superseding D17, restoring D12's >=40% context reduction threshold and zero-regression requirement on M14 benchmarks.
      proof_requirement: Decision lineage audit confirms D12 -> D17 -> D21 consistency without orphaned or conflicting conditions.
    - claim: cli-contract-reconciled-to-14-commands
      pass_boundary: CC-missioncli is version-bumped (v2 -> v3) to accurately list all 14 registered CLI commands and define the typed sources frontmatter schema for Missions and Objectives.
      proof_requirement: Diff check on CC-missioncli matches internal/command.Registry exactly.
    - claim: charter-tokenizer-and-run-lifecycle-frozen
      pass_boundary: The specification freezes spectacular-charter-tokenizer.v1 (o200k_base UTF-8 byte counting) and the 6 Run states (active, paused, blocked, awaiting-review, completed, stopped).
      proof_requirement: Policy review verifies tokenizer and state transition rules are documented with immutable terminal state handling.
    - claim: campaign-contract-transition-map-frozen
      pass_boundary: Campaign roadmap documents the exact Contract version sequence for M16–M21 without modifying past completed mission fingerprints.
      proof_requirement: Campaign document review verifies Contract binding progression across all 7 blocks.
    - claim: progressive-planning-codified
      pass_boundary: Skill guidance (prepare.md) and evals enforce detailing only the active/next block while keeping downstream blocks as high-level sketches.
      proof_requirement: Test fixture confirms Orchestrator prepares single-block previews without expanding or freezing downstream drafts.
completion_record:
    at: "2026-08-22T22:26:49Z"
    authorization: owner supplied --by after schema checks
    by: Alex
    review: RV1
    reviewed_commit: 42840c8ea10a8b8a8a393eb57358277de9e4065e
contract:
    fingerprint: sha256:aa2f59e740e9526bacef1dd9999127861836460e5f2f96b5fe05bc86a458ee1a
    ref: Contract:019fe381-5d61-7223-b362-03a5f99a7b10
dependencies:
    - M14 completed with independent review
gaps: []
objectives:
    - claims:
        - decision-lineage-and-proof-gates-aligned
        - cli-contract-reconciled-to-14-commands
      id: 01a02b83-14f0-71e9-9176-463afe72c74a
      outcome: Reconcile Decision lineage, update CC-missioncli to 14 commands, and freeze sources schema.
      ref: O1
      status: implemented
    - claims:
        - charter-tokenizer-and-run-lifecycle-frozen
        - campaign-contract-transition-map-frozen
      id: 01a02b83-14f0-7da5-ae7e-738dfeb7a113
      outcome: Freeze charter tokenizer, Run lifecycle rules, and campaign contract transition map.
      ref: O2
      status: implemented
    - claims:
        - progressive-planning-codified
      id: 01a02b83-14f0-78bd-9e26-9fca3c8c9cc3
      outcome: Codify progressive proportional planning in skill guidance and evaluation fixtures.
      ref: O3
      status: implemented
outcome: Reconcile the 14-command CLI contract, freeze charter retrieval inputs and tokenizer spec, and establish progressive 1-block-ahead planning before M16 implementation.
owner: Alex
ref: M15
repair_budget: 2
review: independent
reviews:
    - file: reviews/RV1-independent-review-of-m15-governance-and-contract-baseline-reconciliation.md
      id: 01a02b94-bc78-717a-8270-a90c59e1fbf1
      ref: RV1
      verdict: pass
run:
    current_objective: O1
    id: 01a02b83-14f0-7e6d-81bd-e8506cae8d79
    operator: Alex
    ref: R1
    repairs: 0
    started_at: "2026-08-22T22:26:03Z"
    status: completed
scope:
    mechanical:
        - .spectacular/contracts/
        - .spectacular/decisions/
        - .spectacular/proposals/
        - .spectacular/campaigns/
        - skills/spectacular/references/prepare.md
        - test/evals/spectacular/
    semantic:
        - Decision lineage reconciliation (D12/D17/D21)
        - 14-command CLI contract alignment and sources schema
        - Tokenizer and Run state machine specifications
        - Campaign contract transition sequence
        - Progressive planning rules and evals
start_key: sha256:7b9ed8e4410ac1eeb1e7352b4d1bb7557b78aa4070f748cfd20702a7eadc563d
stops:
    - product-code-change
    - public-command-change
    - completed-mission-repoint
    - behavioral-contract-amendment
    - downstream-plan-freeze
validation:
    mode: cli
    schema: mission.v2
---
# Mission: Reconcile P11 Governance and Contract Baselines

## Purpose & Scope
This is a governance-only baseline mission. It makes **zero Go code changes** and introduces **zero new CLI commands**. Its purpose is to lock down the exact contracts, decisions, tokenizer specs, and lifecycle rules before M16 context compiler development begins.

## Key Deliverables & Actions

### 1. Contract Baseline Alignment
- Version-bump [`CC-missioncli`](.spectacular/contracts/CC-missioncli-spectacular-mechanical-cli.md) (`contract_version: "2"` -> `"3"`).
- Update the `command_surface:` section to list all 14 registered commands.
- Define the optional `sources:` frontmatter field for Missions and Objectives.

### 2. Decision Lineage Reconciliation
- Verify that [`D21`](.spectacular/decisions/D21-context-sandwich-execution-gates.md) cleanly supersedes [`D17`](.spectacular/decisions/D17-objective-scoped-runs-and-concurrency.md).
- Ensure the >=40% token reduction gate and M14 benchmark regression gates from [`D12`](.spectacular/decisions/D12-isolation-and-context-compilation.md) are explicitly preserved.

### 3. Freeze Specifications
- **Tokenizer**: Freeze `spectacular-charter-tokenizer.v1` (`o200k_base` byte-exact UTF-8 counting).
- **Run Lifecycle**: Freeze the 6 states (`active`, `paused`, `blocked`, `awaiting-review`, `completed`, `stopped`) and state transition rules.

### 4. Progressive Planning Codification
- Update [`prepare.md`](skills/spectacular/references/prepare.md) with strict rules:
  - Orchestrators plan only 1 block ahead in full detail.
  - Downstream campaign blocks remain fluid summaries.
  - Prohibit dense academic jargon and compound claims in mission plans.
- Add/update evaluation fixtures in `test/evals/spectacular/` to enforce this behavior.
