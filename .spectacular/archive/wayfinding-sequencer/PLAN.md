---
status: archived
priority: high
owner: alex
updated: 2026-08-02
build: b33
summary: "Compute fog and frontier, prioritize uncertainty, surface blockers, map Wayfinder language, and validate dependency/version order"
related:
  - PRD.md
depends-on:
  - wayfinding-contract
blocks:
  - afk-git-hygiene
archived: 2026-08-02
---

# Plan — wayfinding-sequencer

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

Turn Spectacular's durable unknowns into an ordered discovery frontier that reduces uncertainty before confirmed specifications enter execution.

## Constraints

<!-- What's fixed before you start? Inherited from PRD/STACK/PRINCIPLES + request-specific limits. -->

- Dependencies override target versions, priorities, and feature preference.
- `frontier` and `fog` are derived views; canonical dependency IDs remain the stored source of truth.
- Spikes and research outrank deterministic work when frontier nodes have otherwise comparable priority.
- Product/business questions always require user input and may be deferred without disappearing.
- Session onboarding surfaces concise blockers and next choices, never dumps the whole database.
- The sequencer may recommend or run authorized discovery; it cannot merge, make breaking changes, or silently resolve human judgment.

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

Request links are advisory, ideas contain flat open-question prose, status/summary do not compute discovery readiness, and agents can only recover unresolved work by rereading bodies.

### What changes

Build a DAG resolver over canonical IDs, derive fog/frontier, rank high-uncertainty nodes first, surface human blockers during session startup, map Wayfinder phrases to safe verbs, and warn when specs/plans/roadmap/PRD imply dependency or version inversions.

### What stays the same

The request lifecycle, roadmap ledger, validation gates, single-orchestrator mutation model, and archive-not-delete convention remain authoritative.

## Decisions

<!--
  Design calls made inside this request. Format: chose X over Y — because Z.
  Rejected alternatives stay listed; deleting them re-litigates them later.
  Project-wide calls go to DECISIONS.md via `spectacular decide` instead
  (see decisions-rules.md routing table). Empty is fine — no decisions yet.
-->

- Chose computed fog/frontier over stored readiness over duplicated state because dependency changes must update readiness without rewriting downstream nodes.
- Chose risk-first ranking over FIFO when priorities tie because cheap early uncertainty reduction prevents late specification rewrites.
- Chose concise session surfacing over listing every question because the purpose is fast unblocking, not context flooding.
- Chose warnings and proposals over automatic roadmap rewrites because target-version changes can encode product intent.

## Milestones

<!-- Ordered, demoable checkpoints. Outcomes, not tasks. -->
<!-- 3-7 milestones for a typical request. Each is something someone can see working. -->

- M1 — The DAG resolver deterministically classifies frontier, fog, deferred, and resolved nodes and rejects invalid graphs.
- M2 — `spectacular wayfind` identifies the highest-value discovery action and session startup surfaces user-required blockers.
- M3 — Metaphoric commands route correctly from ideas through discovery, spec confirmation, request execution, and parking/defer flows.
- M4 — Cross-layer analysis warns about missing dependencies and discovery/execution target-version inversions.

## Tasks

<!-- Pointer. The executable checklist lives in TASKS.md, grouped by milestone. -->

See `TASKS.md`.

## Dependencies

<!-- Other requests, skills, blocking decisions. Use [[request-slug]] notation. -->

- Hard dependency on [[wayfinding-contract]].
- Existing roadmap ledger, request-link, session, policy, and build-workflow substrates.

## Validation

<!--
  How each milestone is verified. Per-milestone checks.
  Each check states its AUTHORITY: a run: command, an assertable property,
  a judgable artifact, or a human-observable behavior (see verify.md kinds).
  A check with no authority can't fail. Aspiration verbs (improve, enhance,
  optimize, handle gracefully) are not checks.
-->

- M1 — run: DAG fixtures cover chains, diamonds, cycles, dangling IDs, deferred nodes, and resolved prerequisites.
- M2 — run: focused CLI fixtures prove session startup shows the highest-priority `requires_user_input` blocker and deferred loops are omitted.
- M3 — run: phrase-routing tests map park/icebox/find-your-way/act-on-goal without bypassing confirmation or lifecycle gates.
- M4 — run: dependency/version fixtures warn when an earlier execution target relies on a later discovery target and do not mutate the roadmap.

## Deliverables

<!-- Artifacts that ship out of this request. Concrete files, docs, behaviors. -->

- DAG resolver and uncertainty ranking helpers.
- `spectacular wayfind status|next|resolve|defer` CLI/skill flow.
- Session-start blocker briefing and natural-language routing.
- Cross-layer dependency/version doctor checks.
- Focused tests, capability spec, and updated workflow references.
