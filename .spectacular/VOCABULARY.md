---
type: Anchor
id: 01a030b4-6159-7a6a-b77b-e2466a25469b
human_ref: VOCABULARY
title: Domain ontology and ubiquitous language
direction: Give owners and agents one readable model of Spectacular's concepts, relationships, rules, and actions.
---

# Domain ontology and ubiquitous language

This Anchor is the canonical definition of Spectacular's domain. The visual
companion at [atlas/domain-overview.md](atlas/domain-overview.md) is a navigable
projection; this document wins if they differ.

## Glossary index

| Term | Definition | Context |
| --- | --- | --- |
| Anchor | Durable project truth about a named concern. | Project definition |
| A La Carte | Modular adoption of individual Spectacular surfaces without adopting the full mission lifecycle. | Operational mode |
| Architectural HUD | The visual and topological map (Atlas + Campaign DAG) giving operators an instant mental model of system boundaries. | Planning and navigation |
| Assessment | Qualitative evaluation of architectural posture, system maturity, or technical debt; captured in Reviews or Retrospectives without rigid binary gates. | Proof and continuity |
| Atlas | Non-governing visual map of a product or domain slice. | Planning |
| Atlas Coverage | Verification that generated code implements all state transitions and entity boundaries declared in Atlas maps. | Proof and continuity |
| Audit | Independent, read-only forensic inspection performed by an auditor agent or operator to verify contract compliance and claim validity. | Governed execution |
| Blast Radius | The surface area of files, modules, and dependencies impacted by an execution turn; bounded by authorized paths. | Operational safety |
| Capability Contract | A modular specification of an observable capability. | Project definition |
| Context Amnesia | The loss of settled architectural nuances or review findings across fresh model windows; solved by durable Git Markdown records. | Operational continuity |
| Contract | Accepted agreement about a capability and its constraints. | Project definition |
| Decision | An attributable owner choice that resolves a question. | Proof and continuity |
| Decision Compliance | Verification that an implementation strictly adheres to locked architectural rulings (`D<N>`). | Proof and continuity |
| Evidence | Attributable observation supporting a claim. | Proof and continuity |
| Gap | A stated unresolved limitation or dependency. | Proof and continuity |
| Generation Velocity | The rapid rate at which frontier models produce code, requiring structural containment and bounded charters. | Governed execution |
| Handoff | A bounded transfer of work context between operators. | Proof and continuity |
| High-Stakes Code | Consequential changes (auth, crypto, payments, zero-downtime cutovers) requiring strict review and attributable evidence. | Operational mode |
| Maps | Explanatory visual projections (`.spectacular/atlas/`) that navigate domain models without enforcing schema authority. | Planning and navigation |
| Mission | A frozen execution envelope with authority and proof boundaries. | Governed execution |
| Objective | An outcome-sized claim within a Mission. | Governed execution |
| Owner | Person accountable for consequential direction and acceptance. | Governed execution |
| Proof Validity | Verification that test receipts and validation runs are authentic, reproducible, and exit with code 0. | Proof and continuity |
| Proposal | Mutable exploration that may inform later accepted work. | Planning |
| Retrospective | Freeform milestone post-mortem and reflection (`.spectacular/retrospectives/`) capturing lessons learned and forward recommendations without execution ceremony. | Planning and continuity |
| Review | Attributable evaluation of code or claims against an agreed specification or PR (stored in `.spectacular/reviews/`), yielding an explicit verdict (`passed`/`failed`). | Proof and continuity |
| Routine Fast Code | Everyday features, refactors, and fixes where passing tests (`exit 0`) and clean Git commits serve as sufficient proof. | Operational mode |
| Run | A mutable attempt to advance one Objective or Mission. | Governed execution |

## Bounded contexts

| Context | Purpose | Primary objects |
| --- | --- | --- |
| Project definition | States durable direction and product agreements. | Anchor, Capability Contract, Contract |
| Planning | Makes a possible change legible before it gains authority. | Atlas, Proposal |
| Governed execution | Bounds work, responsibility, and recoverable progress. | Mission, Objective, Run, Owner |
| Proof and continuity | Preserves what was observed, reviewed, delegated, or remains unresolved. | Evidence, Review, Handoff, Gap, Decision |

## Notation

The domain model and its visual projections use a fixed notation:

- **Bounded Context**: Rendered as a grouping or Mermaid subgraph; scopes the meaning of terms.
- **Node types**: `Actor`, `Entity`, `Value Object`, `Action`, `Event`, `Policy/Invariant`, `External System`.
- **Relationships**: Directed labelled edges, never nodes. Approved default labels are `owns`, `contains`, `belongs_to`, `has`, `references`, `requests`, `performs`, `emits`, `transitions_to`, `governed_by`, `reads_from`, and `writes_to`.
- **Multiplicity**: Standard UML/ER notation: `1` (exactly one), `0..1` (optional), `1..*` (one or more), `0..*` (zero or more).

## Objects

### Entity: Anchor

Identity: uppercase `<NOUN>.md` filename and UUIDv7. Core Triad (`PROJECT.md`,
`STACK.md`, `ARCHITECTURE.md`) is required at kickoff; specialized Anchors (such
as `VOCABULARY.md`, `SECURITY.md`, `GUARDRAILS.md`) are earned as domain complexity
grows.

### Entity: Capability Contract

Identity: `CC-<name>.md` under `contracts/` with UUIDv7 and `contract_version`.
A modular specification of an observable capability and its constraints.

### Entity: Contract

Identity: `CC-<name>.md` or contract document under `contracts/` with UUIDv7 and
`contract_version`. Lifecycle: immutable version bump on behavioral or schema
change; amended via `contract amend` exclusively for Gap closure or editorial
updates.

### Entity: Decision

Identity: UUIDv7 plus readable `D<n>` reference. Lifecycle: recorded and
immutable; later Decisions may supersede it. Records an attributable owner choice.

### Entity: Evidence

Identity: UUIDv7 plus readable `M<n>/E<n>` reference. Lifecycle: recorded
observation supporting or refuting a frozen completion claim.

### Entity: Gap

Identity: UUIDv7 plus readable `G<n>` reference or Contract gap.
Lifecycle: declared unresolved limitation or dependency; closed by resolution.

### Entity: Handoff

Identity: UUIDv7 plus readable `M<n>/H<n>-<key>` reference. Lifecycle: recorded
bounded transfer of operational context between agents or operators.

### Entity: Mission

Identity: UUIDv7 plus readable `M<n>` reference. Lifecycle: `defined` → `active`
→ `awaiting-assessment` → `resolved` → `completed`. A Mission freezes authority,
scope, claims, and its Contract binding; it owns Objectives and mutable Runs.
Archiving moves the complete bundle to `archive/missions/`.

### Entity: Objective

Identity: UUIDv7 plus readable `M<n>/O<n>` reference. Lifecycle: defined within
the Mission envelope; promoted to a dedicated `objectives/` file when earning
independent review, delegation, or detailed plans.

### Entity: Project

Identity: workspace root and `PROJECT.md` Anchor. A Project owns the
canonical Anchors, Contracts, and history that describe one product.

### Entity: Proposal

Identity: UUIDv7 plus readable `P<n>` reference. Lifecycle: `draft` → `submitted`
→ `accepted`, `rejected`, or `withdrawn`. A Proposal explores direction but
grants no execution authority. Retirement is an archive move to `archive/proposals/`
once resolved.

### Entity: Review

Identity: UUIDv7 plus readable `M<n>/RV<n>` or `A<n>` reference. Lifecycle:
recorded assessment of claims and evidence prior to owner completion.

### Entity: Run

Identity: UUIDv7 plus readable `M<n>/R<n>` reference. Lifecycle: `active` ↔
`paused` / `blocked` / `awaiting-review` → `completed` | `stopped`. A Run bounds
mutable progress and operator execution within frozen Mission scope.

### Policy: Authority

The allowed operator actions, owner-required actions, stop conditions, and
forbidden effects of a Mission or Run. Authority is contextual; it is not a
transfer of owner accountability.

### Value object: Proof boundary

The pass boundary and proof requirement attached to a completion claim. A
passing command, Evidence, Review, and owner completion are distinct facts.

### Policy: Autonomous Execution & Frontier Governance

- **A La Carte**: Independent adoption of isolated Spectacular capabilities (Decisions, Projections, Interview Mode) with zero obligation to adopt the full mission lifecycle.
- **Architectural HUD & Maps**: The non-governing visual projections (`.spectacular/atlas/`) and topological sequence (`campaign check`) that give agents and owners situational awareness.
- **Generation Velocity & Blast Radius**: High-velocity code generation by frontier models is structurally contained using lean charter prompts ($\le 400\text{--}600$ tokens) and OS filesystem watchdogs (`spectacular guard`).
- **Context Amnesia**: Multi-session context resets are bridged by storing immutable Decisions and durable Campaign flight plans directly in Git Markdown.
- **Traceability Dial**: Calibrates proof ceremony: **Routine Fast Code** requires only passing test commands (`exit 0`) and clean Git commits; **High-Stakes Code** demands formal immutable Decisions (`D<N>`), attributable Evidence (`E<N>`), and independent Reviews (`RV<N>`).
- **4-Point Review Rubric (Observe ≠ Act)**: Evaluates **Decision Compliance** (adherence to `D<N>`), **Atlas Coverage** (satisfaction of declared map states), **Blast Radius** (no unauthorized path leaks), and **Proof Validity** (authentic, reproducible test execution).

## Relationships

| Relationship | Meaning | Context |
| --- | --- | --- |
| Project `1` `has` `1..*` Anchor | Every Project has durable named truth; an Anchor belongs to one Project. | `PROJECT`, `STACK`, `ARCHITECTURE`, and earned Anchors |
| Project `1` `has` `0..1` Vocabulary | An earned domain ontology Anchor exists when domain complexity warrants it. | Ubiquitous language |
| Project `1` `has` `0..*` Contract | A Project can govern multiple capabilities. | Capability agreement |
| Mission `0..*` `governed_by` `1` Contract | A Contract may constrain many Missions; a Mission binds one primary Contract. | Frozen baseline |
| Mission `1` `contains` `1..*` Objective | A Mission has one or more outcome-sized claims. | Execution decomposition |
| Objective `1` `has` `0..*` Run | An Objective may have no attempt yet or several attempts. | Mutable progress |
| Mission `0..*` `governed_by` `1` Authority | A Mission or Run executes strictly within an explicit authority envelope. | Operational boundary |
| Mission `1` `has` `0..*` Evidence | Evidence supports claims but does not itself accept them. | Proof |
| Mission `1` `has` `0..*` Review | A Review assesses claims against evidence. | Assessment |
| Mission `1` `has` `0..*` Handoff | Handoffs transfer operational context across runs or operators. | Continuity |
| Mission `1` `has` `0..*` Gap | Gaps record stated unresolved limitations or dependencies. | Governance |
| Owner `1` `performs` `0..*` Decision | Owner decisions are attributable. | Governance |
| Proposal `0..*` `references` `0..1` Atlas | An exploration may attach one overview map; an Atlas may explain many explorations. | Planning |
| Vocabulary `1` `references` `0..*` Atlas | The detailed ontology may have one overview and several slice maps. | Ontology navigation |

## Actions and events

| Action | Preconditions | Effect |
| --- | --- | --- |
| Decide | A consequential question and owner rationale are available. | Records an immutable Decision. |
| Start Mission | Owner-approved Mission plan and valid baseline. | Creates a frozen Mission envelope. |
| Run work | A valid authority envelope exists. | Advances a Run without changing frozen scope. |
| Record evidence | An attributable observation exists. | Preserves a claim-supporting observation. |
| Complete Mission | Frozen criteria and required proof/review are satisfied; owner accepts. | Records completion without rewriting historical bindings. |

Key events: `Decision recorded`, `Mission activated`, `Run transitioned`,
`Evidence recorded`, `Review recorded`, and `Mission completed`.

## Invariants and policies

- Canonical Markdown is readable without a database or projection.
- The owner decides outcome, semantic scope, and consequential acceptance.
- A Proposal or Atlas never authorizes execution.
- A documented relationship is not an enforced invariant unless its implementation
  mapping names a schema, code rule, API contract, or test.
- The Vocabulary is canonical; its Atlas maps are explanatory and non-governing.

### Policy: The Verification Triad (Audit, Assessment, Review)

- **Audit (The Action / Verb)**: An independent, read-only inspection process conducted by an auditor subagent or human reviewer to discover contract drift or verify invariants before completion. The audit itself produces observations.
- **Review (The Artifact / Verdict)**: The canonical, durable evaluation record stored in `.spectacular/reviews/` (either standalone for PRs/branches or mission-scoped). It states an attributable reviewer verdict (`passed`, `failed`), proof basis, and claim checks.
- **Assessment (The Qualitative Scale)**: A spectrum-based evaluation of system health, tech debt, or architectural fit. Unlike reviews, assessments do not enforce binary pass/fail gates and are recorded within Reviews or Retrospectives.

### Policy: Guardrails and Dynamic CLI Injection

- `GUARDRAILS.md` is an authoritative project Anchor containing repository-wide technical constraints, safety rules, and banned patterns.
- Unlike static agent harness prompts (e.g. `.claude/rules/` or `.agents/rules/`), Spectacular guardrails are injected dynamically by the CLI during mission preparation and validated deterministically by `test/verify.sh preflight`.
- Thin harness pointers or relative symlinks to `GUARDRAILS.md` remain an optional convenience for local IDEs, never a mandatory workspace dependency.

## Implementation mappings

| Domain concern | Implementation boundary | Proof |
| --- | --- | --- |
| Typed identity and lifecycle | `internal/domain/`, `internal/missionbundle/` | Go unit and acceptance tests |
| Canonical record discovery | `internal/discovery/`, `.spectacular/` | discovery tests and CLI checks |
| Owner Decisions | `spectacular decide` and `decisions/` | atomic command receipt |
| Ontology Anchor and maps | `skills/spectacular/references/`, `atlas/`, project Anchors | documentation consistency review |

## Semantic gaps and change history

- **Current gap:** ontology-impact declarations are workflow guidance, not a
  typed Mission field or mechanically checked conformance rule.
- **D25-vocabulary-canonical-domain-ontology-anchor:** establishes this Anchor and its notation.
- **D26-atlas-domain-maps-non-governing-visual-projections:** establishes non-governing Atlas domain maps.
- **D27-ontology-impact-explicit-in-planning:** establishes ontology impact in planning guidance without a typed Mission field.
- **D28-dynamic-operating-dial-and-anchors:** establishes the Dynamic Operating Dial, 5 Foundational Anchors, and Tiered Verification Protocol.
