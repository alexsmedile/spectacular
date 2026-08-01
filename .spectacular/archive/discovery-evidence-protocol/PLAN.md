---
status: archived
priority: medium
owner: alex
updated: 2026-08-02
build: b37
docs_impact: required
summary: "Define progressive routing for research, spikes, prototypes, tracer bullets, artifacts, and technical debt"
related:
  - PRD.md
docs_impact_evidence: Updated README.md, docs/workflow.md, docs/scaffold.md, architecture, system spec, and routed skill references
archived: 2026-08-02
---

# Plan — discovery-evidence-protocol

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

Define a progressive discovery protocol that routes genuine uncertainty through research, spikes, prototypes, or tracer execution without manufacturing nodes, folders, or ceremony when the implementation path is already clear.

## Constraints

<!-- What's fixed before you start? Inherited from PRD/STACK/PRINCIPLES + request-specific limits. -->

- Research gathers evidence but never chooses among product or architectural alternatives by itself.
- Existing canonical identities stay stable: `RES-NNN` for research and `SPK-NNN` for technical spikes; conversational speech stays compact (`R1`, `SPK1`).
- Prototype and tracer-bullet concepts must reuse existing lifecycle surfaces rather than create overlapping databases prematurely.
- Discovery work is demand-driven: create a node only for a named uncertainty that cannot be answered more cheaply from existing code, tests, documentation, or direct user clarification.
- Production code may come only from an approved specification and normal request lifecycle; disposable spike/prototype code never silently becomes implementation.
- Technical debt remains execution/backlog context, not a new soft database, unless later evidence proves a dedicated lifecycle is necessary.

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

Research (`RES`) and spike (`SPK`) nodes already carry evidence and `supported|refuted|inconclusive` results. `PRT` is reserved but undefined, findings are reserved, roadmap guidance mentions prototypes, and request artifacts can store prototypes or research. These pieces do not yet form one concise routing protocol, leaving room for duplicate records and unnecessary discovery ceremony.

### What changes

Add one canonical progressive-discovery protocol that distinguishes read-only research, technical spikes, UX/workflow prototypes, production tracer bullets, and durable knowledge artifacts. Define the cheapest-answer-first gate, artifact promotion rules, safe branch/code destinations, and technical-debt routing. Align the skill router, collection rules, architecture, system spec, and user workflow docs.

### What stays the same

No new CLI mutator, resolver, lifecycle, collection, or `DEB`/`ART` identity ships. `PRT` stays reserved; UX prototypes remain attached artifacts until a standalone workflow earns its cost. `FND` stays reserved and completed research/spike records remain sufficient unless a reusable cross-cutting finding workflow ships later.

## Decisions

<!--
  Design calls made inside this request. Format: chose X over Y — because Z.
  Rejected alternatives stay listed; deleting them re-litigates them later.
  Project-wide calls go to DECISIONS.md via `spectacular decide` instead
  (see decisions-rules.md routing table). Empty is fine — no decisions yet.
-->

- Chose `RES` over `RCH`, `RSC`, `SER`, or `SRC` because `RES` is already shipped, clearly expands to research, and lets users speak the friendlier `R1` without migration.
- Chose `SPK` as the technical-experiment node and kept `PRT` reserved because prototypes validate human experience and can live as request/vision artifacts without another project-wide lifecycle.
- Chose `execution_mode: tracer` on an approved `SPC` over a tracer node type because tracer code is production implementation, not disposable discovery evidence.
- Chose “artifact” as an umbrella term over an `ART` ID or catch-all folder because location and owning record carry the useful semantics.
- Chose existing requests, roadmap candidates, ideas, and decisions for technical debt over `DEB-NNN` because Spectacular is not a ticket tracker and debt needs ownership/priority, not another passive ledger.

## Milestones

<!-- Ordered, demoable checkpoints. Outcomes, not tasks. -->
<!-- 3-7 milestones for a typical request. Each is something someone can see working. -->

- M1 — One decision table routes uncertainty to the cheapest sufficient mechanism and explicitly says when to skip discovery.
- M2 — Spike, prototype, tracer-bullet, artifact, and technical-debt boundaries map to existing IDs, owners, code destinations, and durability rules without overlap.
- M3 — Canonical and user-facing documentation presents the same progressive workflow and passes structural checks.

## Tasks

<!-- Pointer. The executable checklist lives in TASKS.md, grouped by milestone. -->

See `TASKS.md`.

## Dependencies

<!-- Other requests, skills, blocking decisions. Use [[request-slug]] notation. -->

- Builds on verified [[lifecycle-contract]], [[wayfinding-contract]], [[wayfinding-sequencer]], and [[afk-git-hygiene]] behavior.
- Complements [[advanced-engineering-collections]] without activating its reserved finding or benchmark workflows.

## Validation

<!--
  How each milestone is verified. Per-milestone checks.
  Each check states its AUTHORITY: a run: command, an assertable property,
  a judgable artifact, or a human-observable behavior (see verify.md kinds).
  A check with no authority can't fail. Aspiration verbs (improve, enhance,
  optimize, handle gracefully) are not checks.
-->

- M1 — judge: every route starts with an explicit uncertainty and the protocol selects no node when code/tests/docs or a direct question can answer it.
- M2 — assert: `RES`, `SPK`, reserved `PRT`, `SPC` tracer mode, generic artifacts, and debt routing have exactly one owner and no conflicting lifecycle claim.
- M3 — run: `bash -n cli/spectacular cli/install.sh scripts/hooks/pre-commit`, `scripts/hooks/pre-commit --check`, targeted doctor checks, `bash tests/run.sh`, and `git diff --check` all pass.

## Deliverables

<!-- Artifacts that ship out of this request. Concrete files, docs, behaviors. -->

- Canonical `discovery-protocol.md` skill reference and router entry.
- Aligned research, spike, lifecycle, soft-DB, ID, architecture, system-spec, workflow, and scaffold guidance.
- Explicit no-new-folder routing for technical debt and no-new-entity routing for generic artifacts.
