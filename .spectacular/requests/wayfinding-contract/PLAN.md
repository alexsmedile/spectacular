---
status: verified
priority: high
owner: alex
updated: 2026-08-01
build: b32
source_spec: SPC-001
summary: "Add canonical IDs, decisions/questions/ideas databases, research and spike node contracts, and confirmed/unconfirmed specification state"
related:
  - PRD.md
blocks:
  - wayfinding-sequencer
  - afk-git-hygiene
---

# Plan — wayfinding-contract

<!--
  Canonical 7-slot PLAN template for a single request.
  Lives at: .spectacular/requests/<slug>/PLAN.md

  Rules:
  - PLAN is per-request. PRD is project-wide. Never put a PRD inside requests/.
  - This file's frontmatter `status:` is the single source of lifecycle state for the request.
  - The 7 required sections must appear IN ORDER, unnumbered:
      ## Goal, ## Constraints, ## Milestones, ## Tasks, ## Dependencies, ## Validation, ## Deliverables
    Extra sections (## Understanding, ## Decisions, request-specific) may appear
    BETWEEN them; doctor enforces the required set's presence + order, not a closed list.
  - All 7 required sections must be filled before this PLAN is considered usable.
  - Replace every <placeholder> with concrete content.
-->

## Goal

<!-- One sentence. What does this request change? -->
<!-- Compress the request's intent. Aligns with PRD's Vision or Goals — this is a slice, not a restatement. -->

Extend Spectacular's shared Markdown workspace so humans and agents can carry ideas, open questions, research, spikes, decisions, and confirmed or unconfirmed specifications across sessions with stable canonical identities.

## Constraints

<!-- What's fixed before you start? Inherited from PRD/STACK/PRINCIPLES + request-specific limits. -->

- Markdown remains the human-readable source of truth; CLI and skill behavior are accelerators, not prerequisites.
- Canonical cross-references use zero-padded IDs; user-facing output prefers compact aliases such as `D1` and `Q1`.
- Canonical prefixes are `DEC`, `QUE`, `IDEA`, `RES`, `SPK`, `PRT`, `SPC`, and reserved `TSK`; legacy `D<N>` and unnumbered ideas must remain readable during migration.
- `IDEA-001` is canonical; `IDE-001` is an accepted legacy/input alias only.
- Specs move through `unconfirmed` to `current`; confirmed specifications are the only specifications eligible to generate implementation requests.
- Discovery/execution release labels use SemVer prerelease forms such as `v1.0.0-discovery` and `v1.0.0-execution`.
- Product/business questions require human resolution; evidence-backed technical decisions may be drafted automatically but remain traceable to sources or prototype evidence.

<!-- ## Understanding and ## Decisions below are OPTIONAL extra sections,
     allowed between Constraints and Milestones. -->

## Understanding

<!--
  OPTIONAL authoring slot, but REQUIRED before `planned → active` by the
  `understand-before-change` policy (@Implementation). Fill it here for a
  typical request; escalate to a dedicated requests/<slug>/UNDERSTANDING.md
  (same three subheads) for large ones — the policy is satisfied by EITHER.
  Not one of the 7 required authoring slots; it gates implementation, not planning.
-->

### How it works now

Decisions use `D<N>` files, ideas use unnumbered kebab-case files, specifications use capability slugs, and there is no dedicated questions, research, or spike substrate. Draft/current spec state and canonical cross-store aliases are not represented consistently.

### What changes

Introduce a canonical ID resolver and schemas/templates for decisions, questions, ideas, research, spikes, prototype artifacts, and specifications. Add additive compatibility reads, deterministic allocation, explicit spec confirmation state, and migration tooling that archives or snapshots before rewrites.

### What stays the same

PLAN/TASKS remain the execution unit, the request lifecycle remains pure, the roadmap ledger remains the request-to-release authority, and all destructive or breaking operations retain human gates.

## Decisions

<!--
  Design calls made inside this request. Format: chose X over Y — because Z.
  Rejected alternatives stay listed; deleting them re-litigates them later.
  Project-wide calls go to DECISIONS.md via `spectacular decide` instead
  (see decisions-rules.md routing table). Empty is fine — no decisions yet.
-->

- Chose three user-facing knowledge domains (`decisions/`, `questions/`, `ideas/`) over one generic exploration collection because users need distinct mental models for settled choices, active blockers, and parked inspiration.
- Chose canonical IDs with compact aliases over filename-only identity because dependency graphs and conversational references must survive cross-store navigation.
- Chose `IDEA` over `IDE` as the canonical prefix because the confirmed reservation says `IDE → IDEA`; `IDE` remains an input alias.
- Chose additive compatibility and explicit migration over immediate destructive renames because existing workspaces and links must remain usable.
- Chose `unconfirmed | current | deprecated` specification state over treating draft specs as requests because ideas/specification/execution are distinct phases.

## Milestones

<!-- Ordered, demoable checkpoints. Outcomes, not tasks. -->
<!-- 3-7 milestones for a typical request. Each is something someone can see working. -->

- M1 — Canonical identity and alias contract is documented, deterministic, and backwards-compatible.
- M2 — Questions, ideas, research, spikes, and decisions can be created and read as typed Markdown records; prototype artifacts are linked from spikes with reserved `PRT` IDs.
- M3 — Specifications visibly move from unconfirmed to current and current specs can seed requests with provenance.
- M4 — Existing workspace records migrate safely with snapshots/archives and all canonical cross-references remain resolvable.

## Tasks

<!-- Pointer. The executable checklist lives in TASKS.md, grouped by milestone. -->

See `TASKS.md`.

## Dependencies

<!-- Other requests, skills, blocking decisions. Use [[request-slug]] notation. -->

- Existing soft-DB, doc-engine, snapshot, and migration substrates.
- Enables [[wayfinding-sequencer]] and [[afk-git-hygiene]].

## Validation

<!--
  How each milestone is verified. Per-milestone checks.
  Each check states its AUTHORITY: a run: command, an assertable property,
  a judgable artifact, or a human-observable behavior (see verify.md kinds).
  A check with no authority can't fail. Aspiration verbs (improve, enhance,
  optimize, handle gracefully) are not checks.
-->

- M1 — run: focused canonical-ID tests resolve canonical, prefixed alias, and context-aware numeric input deterministically; ambiguous naked numbers refuse non-interactive execution.
- M2 — run: focused CLI tests create/list/show/resolve each supported record with canonical IDs and valid frontmatter.
- M3 — assert: `unconfirmed` specs cannot produce an execution request; confirmation produces `current` state and preserves provenance.
- M4 — run: migration fixtures preserve legacy records, links, and content; `spectacular doctor` reports no ID or reference drift.

## Deliverables

<!-- Artifacts that ship out of this request. Concrete files, docs, behaviors. -->

- Canonical-ID and entity contract documentation.
- CLI resolver/allocation helpers and collection verbs.
- Questions/research/spike/prototype/spec templates and rules files.
- Backwards-compatible migration plus doctor checks.
- Focused Bash 3.2-compatible tests and updated system spec.
