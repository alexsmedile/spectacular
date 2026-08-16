---
type: Mission
id: 01a009ff-ce94-724e-a6f8-66783f1a4003
ref: M5
title: Align the Spectacular Skill with compact Missions
status: completed
owner: Alex
created: "2026-08-16T09:55:54Z"
updated: "2026-08-16T10:52:24Z"

contract:
  ref: Contract:019fe381-5d61-7223-b362-03a5f99a7b10
  fingerprint: sha256:aa2f59e740e9526bacef1dd9999127861836460e5f2f96b5fe05bc86a458ee1a

baseline:
  commit: 67315afc50a197d245f14c90da7b45dad1506973
  branch: codex/lean-launch-context

outcome: The Spectacular Skill operates the compact Mission model immediately, assigns meaning to LLM judgment and repeatable invariants to mechanics, and leaves a frozen Contract and separate Mission plan for the CLI implementation.
review: independent

completion:
  - claim: skill-model
    pass_boundary: The core Skill and routed references consistently teach optional Proposal exploration, compact MISSION.md execution, earned expansion, ADR-like Decisions, ordinary specification edits, and one completion gate.
    proof_requirement: Focused content checks find the new vocabulary in its owning references and find no mandatory preparation receipt, lifecycle Decision, Mission index, or reconciliation path in active guidance.
  - claim: judgment-mechanics-boundary
    pass_boundary: Guidance assigns contextual meaning, prose, decomposition, and problem-solving to the LLM while reserving schema invariants, exact bindings, atomic transitions, and refusals for supported mechanics.
    proof_requirement: The core Skill gives concrete routing rules for supported tooling and manual-bootstrap work without claiming that either surface replaces owner judgment.
  - claim: progressive-context
    pass_boundary: Routine orientation and execution load MISSION.md, the current Objective, and exact sources before optional detail; obsolete CLI catalogs and unchanged project history are not routine context.
    proof_requirement: Skill structure and reference checks confirm one routed file per job, compact continuity, earned file promotion, and no duplicate policy blocks.
  - claim: cli-work-defined
    pass_boundary: A focused mechanical CLI Contract and planned M6 freeze the accepted decoder, command, validation, review, refusal, atomicity, and stress-test properties without starting CLI implementation inside M5.
    proof_requirement: The Contract and M6 are readable from the filesystem, bind exact identity, cover every accepted property, and keep M6 inactive until M5 completion.

objectives:
  - ref: O1
    id: 01a009ff-ce94-7249-accc-9c2a089d3080
    outcome: Rewrite the core Skill around the compact Mission schema and the judgment/mechanics boundary.
    status: implemented
    claims: [skill-model, judgment-mechanics-boundary]
  - ref: O2
    id: 01a009ff-ce94-7585-ae70-34adfff1c44b
    outcome: Rewrite routed preparation, execution, runtime, review, completion, and audit guidance with progressive disclosure.
    status: implemented
    after: [O1]
    claims: [progressive-context]
  - ref: O3
    id: 01a009ff-ce94-7c74-a7d8-0f384b429fbc
    outcome: Freeze the mechanical CLI Contract and a separate planned Mission with adversarial validation properties.
    status: implemented
    after: [O2]
    claims: [cli-work-defined]

run:
  ref: R1
  id: 01a009ff-ce94-7323-a942-754ef422e264
  status: completed
  operator: Codex primary session
  started_at: "2026-08-16T10:31:30Z"
  current_objective: O3
  repairs: 1

reviews:
  - ref: RV1
    id: 01a00a33-87b8-7e0e-9100-450696ad1e80
    file: reviews/RV1-independent-review.md
    verdict: pass

completion_record:
  by: Alex
  at: "2026-08-16T10:52:24Z"
  authorization: Owner instructed M5 and M6 to continue through completion if no unresolved discrepancy remained.
  reviewed_commit: 4074708c26c1158f4eb778b55c86aabe80979e76
  review: RV1
  limitations:
    - Legacy CLI acceptance validation remains out of service until M6 replaces it.

activation:
  by: Alex
  at: "2026-08-16T10:31:30Z"
  fingerprint: sha256:3314b10debed8b94a6482dac8109685b039426fd8757a410c39f06dad892569f

validation:
  schema: mission.v2
  mode: manual-bootstrap

authority:
  operator: [inspect, edit-in-scope, choose-reversible-implementation, run-checks, generate-derived-files, bounded-repair, commit-local]
  requires_owner: [change-outcome-or-completion, expand-scope, push, merge, release, irreversible-change, destructive-data, secret-change]

scope:
  mechanical: [skills/spectacular/, cmd/assemble-release/main_test.go, internal/humanlayout/index.go, internal/humanlayout/layout.go, internal/humanlayout/layout_test.go, .spectacular/PROJECT.md, .spectacular/index.md, .spectacular/contracts/, .spectacular/missions/index.md, .spectacular/missions/M5-implement-compact-expandable-missions/, .spectacular/missions/M6-implement-compact-mission-cli/]
  semantic:
    - Spectacular Skill behavior for compact Mission preparation, execution, review, completion, and recovery.
    - Division of work between LLM judgment and deterministic mechanics.
    - Mechanical CLI requirements and stress properties, but not CLI implementation.

repair_budget: 2
dependencies: []
gaps: []

stops:
  - Expanding M5 into CLI implementation instead of leaving that work to M6.
  - Moving contextual product judgment into rigid mechanics or exact repeatable invariants into prose alone.
  - Reintroducing mandatory Proposal, preparation receipt, lifecycle Decision, Mission index, or reconciliation ceremony.
---
# Mission

## Origin and amendment

The owner split the work after the compact schema checkpoint: update the Skill
first, then implement the mechanical CLI in a separate Mission. This
owner-confirmed amendment replaces M5's original CLI implementation scope while
preserving its identities and compact file structure.

The owner had also accepted indexes as rebuildable navigation caches rather
than mandatory per-Mission files. The first independent review exposed stale
workspace caches and an old generator that still created Mission-local
`index.md` files. This bounded repair adds collection-cache integration to M5
without making any cache authoritative or entering M6's command work.

## LLM and mechanical work

Use the LLM directly for fuzzy intent, adaptive questions, criteria wording,
semantic scope, coherent Objectives, prose, contextual problem-solving, and
drafting canonical Markdown. These jobs benefit from broad understanding and
usually cost less than encoding every edit through a rigid command payload.

Use scripts or CLI for YAML/schema validity, UUIDv7 and reference integrity,
fingerprints, baseline checks, dependency graphs, authority vocabulary, safe
paths, atomic multi-file transitions, retries/concurrency, state transitions,
compact projections, and exact refusals. These jobs are repeated, objectively
checkable, and expensive to perform inconsistently.

The plan supplies meaning; supported mechanics supply invariants. Neither
surface decides owner intent or claims that structural validity proves quality.

## Execution plan

### O1 — Core Skill

- Replace the old Proposal/receipt/Decision/reconciliation lifecycle with the
  approved working model and compact Mission schema.
- Teach one preflight, progressive context, earned expansion, owner/operator
  authority, and the manual-bootstrap fallback.
- State the LLM/mechanics decision rule in concrete terms.

### O2 — Routed guidance

- Make Proposal optional exploration and adaptive grilling conditional.
- Freeze Mission activation directly with an owner-attributed boundary
  fingerprint.
- Promote Objectives for delegation and Runs for real execution boundaries.
- Record earned Evidence/review only when completion requires it.
- Close through one Mission owner gate with product specifications in the same
  worktree.

### O3 — CLI Contract and next Mission

- Write one focused Contract for the typed Mission-bundle decoder, minimal
  command surface, schema-owned validation, independent-review record, precise
  refusals, atomic transitions, and stress properties.
- Create M6 as a planned, inactive Mission bound to that Contract and dependent
  on M5 completion.
- Keep current CLI mutation and validation out of service for compact Missions.

## Bootstrap conditions

Current CLI behavior implements the superseded lifecycle and is not proof for
M5. Manual checks own YAML syntax, UUIDv7 uniqueness, Contract/baseline and
activation fingerprints, completion coverage, Objective dependency order,
current Run, authority/scope, and one-file placement.

## Independent review

After O3, create one earned review file under `reviews/`. A fresh reviewer checks
all four claims against this exact tree, reports claim verdicts, findings, and
limitations, and does not rewrite criteria or repeat implementation work. M5
then returns one owner completion gate; M6 remains planned until that gate closes.
