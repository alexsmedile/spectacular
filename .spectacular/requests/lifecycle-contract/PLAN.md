---
status: verified
priority: high
owner: alex
updated: 2026-08-01
build: b35
docs_impact: required
summary: "Unify Spectacular entity lifecycles, evidence gates, AFK authorization, archival behavior, and documentation impact under one enforceable contract"
related:
  - PRD.md
docs_impact_evidence: Updated docs/commands.md, docs/workflow.md, docs/scaffold.md, canonical architecture/spec index, and skill references
---

# Plan — lifecycle-contract

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

Make every Spectacular entity lifecycle explicit, evidence-backed, compatible with existing workspaces, and enforced from one canonical contract so the workspace reduces uncertainty without pretending Markdown is fresher than code.

## Constraints

<!-- What's fixed before you start? Inherited from PRD/STACK/PRINCIPLES + request-specific limits. -->

- Actual code is authoritative; an `implemented` specification is a historical claim anchored to a commit/build, not continuously synchronized documentation.
- Existing files remain readable. Legacy status and path migrations are preview-first, archive-first, and require explicit user confirmation.
- Request lifecycle semantics stay `planned → active → review → verified → archived`; holds remain orthogonal.
- Root anchors remain current-by-location and snapshot-before-edit; they do not gain a status machine.
- Bash 3.2 compatibility, deterministic Markdown storage, canonical IDs, dry-run defaults, and existing HITL boundaries remain intact.
- `PRT` and `TSK` remain reserved only; their storage/lifecycles are explicitly deferred in `TODO.md`.

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

Request, idea, session, and debug lifecycles are mostly explicit, but specification status is contradictory across canonical docs, references, doctor, CLI, tests, and public docs. Research/spikes overstate outcomes as `verified`; decisions accept arbitrary autonomous evidence; fixes may be created unverified; memories cannot be corrected; AFK authorization is Git-only; docs impact is advisory; branch cleanup relies on expiring reflogs; plural collection paths drift.

### What changes

Add one lifecycle contract and align the CLI, doctor, templates, migration behavior, skill references, architecture, system spec, and public docs. Introduce the approved specification lifecycle, evidence-anchored implementation, terminal spec archival, completed/result discovery semantics, decision/memory supersession, verified-only fixes, durable AFK runs, durable Git cleanup refs, documentation-impact gates, unique per-collection slugs, and conservative legacy mappings.

### What stays the same

Code remains the real source of behavior; specs remain lightweight execution context. Requests still own implementation progress and `source_spec`; release state remains separate; root docs use snapshots; archives preserve history; remote deletion, merge, sensitive access, breaking changes, and product/business judgment remain human-gated.

## Decisions

<!--
  Design calls made inside this request. Format: chose X over Y — because Z.
  Rejected alternatives stay listed; deleting them re-litigates them later.
  Project-wide calls go to DECISIONS.md via `spectacular decide` instead
  (see decisions-rules.md routing table). Empty is fine — no decisions yet.
-->

- Chose sibling entry states `draft` (collaborative) and `unconfirmed` (AFK) because provenance changes the review gate but neither is authorized execution.
- Chose `approved → implemented` over `current` because Spectacular cannot promise a Markdown spec remains synchronized with code.
- Chose request-owned `source_spec` and derived inverse links over storing `execution_request` on specs because one spec may produce multiple requests and duplicated links drift.
- Chose `completed` plus `result: supported | refuted | inconclusive` for research/spikes because finishing an investigation is distinct from proving its hypothesis.
- Chose durable supersession/retraction over edits or deletion for decisions and memories.
- Chose `candidate` before roadmap `planned` because likely-near work is not yet committed scope.
- Chose one canonical lifecycle contract over repeated enum copies because drift has already made structurally clean workspaces semantically contradictory.

## Milestones

<!-- Ordered, demoable checkpoints. Outcomes, not tasks. -->
<!-- 3-7 milestones for a typical request. Each is something someone can see working. -->

- M1 — One canonical lifecycle contract defines every entity, transition, evidence gate, relationship owner, and archive behavior.
- M2 — Specifications and legacy migrations follow `draft|unconfirmed → approved → implemented → superseded|deprecated → archived` without false freshness claims or duplicated provenance.
- M3 — Discovery, decisions, questions, memories, and fixes preserve uncertainty and history with enforceable evidence and supersession rules.
- M4 — Goal-scoped AFK authorization, durable branch recovery, and documentation-impact checkpoints prevent silent scope or data loss.
- M5 — Doctor, templates, tests, architecture, skill routing, and public docs enforce and explain the same contract.

## Tasks

<!-- Pointer. The executable checklist lives in TASKS.md, grouped by milestone. -->

See `TASKS.md`.

## Dependencies

<!-- Other requests, skills, blocking decisions. Use [[request-slug]] notation. -->

- Builds directly on verified [[wayfinding-contract]], [[wayfinding-sequencer]], and [[afk-git-hygiene]].
- Uses the existing request lifecycle, snapshot, migration, policy, archive, and doctor substrates without replacing them.

## Validation

<!--
  How each milestone is verified. Per-milestone checks.
  Each check states its AUTHORITY: a run: command, an assertable property,
  a judgable artifact, or a human-observable behavior (see verify.md kinds).
  A check with no authority can't fail. Aspiration verbs (improve, enhance,
  optimize, handle gracefully) are not checks.
-->

- M1 — assert: the canonical contract contains a closed state/transition table for every lifecycle-bearing collection and all other references link to it.
- M2 — run: focused lifecycle CLI tests prove approved-only action, evidence-anchored implementation, atomic supersession, safe archival, duplicate refusal, and conservative migration.
- M3 — run: focused collection tests prove evidence-required discovery outcomes, sourced DEC capture, memory correction states, question provenance, and verified-only fixes.
- M4 — run: AFK fixtures prove run gating and durable branch restoration refs; request fixtures prove documentation-impact closure gates.
- M5 — run: `bash -n cli/spectacular cli/install.sh scripts/hooks/pre-commit`, `scripts/hooks/pre-commit --check`, focused tests, full `bash tests/run.sh`, and `./cli/spectacular doctor` all pass.

## Deliverables

<!-- Artifacts that ship out of this request. Concrete files, docs, behaviors. -->

- Canonical lifecycle contract and aligned entity rules.
- Bash 3.2-compatible CLI verbs, migrations, doctor findings, and templates.
- Focused regression tests covering all new gates and compatibility paths.
- Updated architecture, system spec, workflow, commands, scaffold, TODO, and skill routing documentation.
