# Mission anatomy

Use this when: Agent authoring, splitting, or auditing detailed Spectacular record structures and fields.

The full field lists. `SKILL.md` carries the summary; this file carries the
inventory. Load it when you are writing or auditing a Mission record, not on every
session.

## Shape

Retrieve a governed Mission input shape through `spectacular mission start
--schema`, then round-trip it through the emitting validator. Inspect existing
records as examples of accepted history, not as templates for fabricating an
active Mission. Activation identity, bindings, and fingerprints belong to the CLI.

## Frontmatter

Frozen at activation, covered by the activation fingerprint:

- `id` — UUIDv7, durable identity
- `ref` — human-readable navigation reference
- `owner` — who alone may change the frozen parts
- activation time and activation fingerprint
- exact Contract binding
- exact Git baseline binding
- `outcome` — what this Mission is for
- completion claims — the frozen list to prove
- review level, including whether review must be independent
- authority, and the forbidden-effect ceiling
- mechanical scope and semantic scope
- budgets, including the repair budget
- dependencies, Gaps, stops

Mutable, and deliberately **outside** the fingerprint:

- `status`
- inline Objectives and their progress
- current Run and its state
- repair count
- validation mode

## Markdown body

- origin — where this Mission came from, including any Proposal reference
- rationale — why this approach, and what was rejected
- detailed Objective plans that outgrew the frontmatter
- bootstrap conditions, when the Mission declares `manual-bootstrap`
- examples
- review instructions for the reviewer

## What the tooling owns

The plan never chooses these. The active schema registry does:

- schema and vocabulary validation
- UUIDv7 and ref allocation
- fingerprint computation
- baseline checks
- dependency integrity
- safe path handling
- atomic multi-file transitions
- concurrency and retry behavior
- compact projections
- exact refusal messages

## What the plan owns

The tooling never decides these:

- outcome
- completion criteria
- decomposition into Objectives
- semantic scope
- authority
- dependencies, Gaps, stops
- rationale and prose

## Divide meaning from mechanics

| Use tooling when | Use judgment when |
|---|---|
| failure is expensive | meaning depends on context |
| the rule is exact and repeated | the prose is the value |
| the transition must be atomic | several answers are valid |
| | encoding it mechanically costs more than checking the result |

## Growth

Start with one file: `<mission-dir>/<mission-ref>-<slug>.md` (e.g. `.spectacular/missions/M5-implement-compact-missions/M5-implement-compact-missions.md`).

- Add `objectives/` (`O<N>-<slug>.md`) when an Objective earns its own detail, delegation, owner, or
  independent review.
- Add `runs/` (`<run-dir>/<run-ref>-<slug>.md`) when a Run has a distinct job, operator, baseline, or recovery
  boundary.

These decomposition choices apply to unactivated drafts. Existing governed
Objectives and Runs cannot be moved or replaced with pointers by ordinary edits.
Use only a supported governed transition and validate the result; where none
exists, draft a restructuring under [bootstrap.md](bootstrap.md) with explicit
owner authorization. The root Mission remains its own bundle entry point.

## Checkpoints

Plan optional checkpoints in the Run body when a Run needs a named progress,
verification, or resume gate. A checkpoint itself does not grant authority or
require human review. When it produces a decision, observation, verdict, or
handoff, create the corresponding Decision, Evidence, Review/Assessment, or
Handoff record and link it from the Run-body note. See
[close.md](close.md) for claim-to-evidence accounting.

## Campaign context

When a Mission comes from a Campaign, cite the Campaign file and block in the
Mission body's origin or rationale. This context is non-binding: Campaigns are
mutable roadmap maps under `.spectacular/campaigns/`, while the Mission's frozen
outcome, scope, authority, and completion claims remain authoritative. Do not
add a Campaign binding to Mission frontmatter.

## Anchor anatomy & modular contracts

**The Anchor naming rule**: Bare single-word uppercase names (`<NOUN>.md` e.g. `PROJECT.md`, `STACK.md`, `ARCHITECTURE.md`, `README.md`, `AGENTS.md`) are reserved exclusively for Project Anchors and workspace landmark contracts. All governed records (Missions, Runs, Objectives, Proposals, Reviews, Decisions, Evidence, Gaps) carry their scoped prefix in their filename.

### Core Triad (Required at Kickoff)
- `PROJECT.md`: Direction, immutable boundaries, non-goals, and `current_truth` binding.
- `STACK.md`: Language/runtime versions, allowed dependencies, database engines, baseline verification command.
- `ARCHITECTURE.md`: Layering pattern (e.g. hexagonal/clean), directory layout, dependency directions between Domain, Store, Server, and API.

### On-Demand Anchors (Earned only)
Specialized anchors emerge only when domain or operational complexity exceeds inline thresholds:
- `ROADMAP.md`: Macro-level product evolution and strategic multi-horizon themes. Decomposes into mid-term Campaigns as milestones enter active planning.
- `VOCABULARY.md` (retain an existing `ONTOLOGY.md` without duplicate authority): Canonical domain ontology and ubiquitous language (D25, D29). Its glossary index is alphabetical; for the detailed section skeleton, see the body shape in [genesis-examples.md](genesis-examples.md).
  * *Threshold*: If <= 3-4 simple entities with no ambiguous terms or shared rules, keep them inline in `PROJECT.md`.
  * *Earned triggers*: (1) Synonym collision / naming ambiguity (e.g. `User` vs `Account`, `Job` vs `Task`); (2) Non-trivial state machine invariants (e.g. `DRAFT` -> `ACTIVE` -> `REVIEW`); (3) Relationships, permissions, or actions that span several concepts; (4) Bespoke non-standard concepts (e.g. `Anchor`, `Gap`, `Handoff`); (5) Multi-contract shared models.
  * *Ontology structure*: Must include an explicit **Permitted Actions & Banned Synonyms** table (`Canonical Action` vs `BANNED Synonyms`) and **Permitted Entity States** enumeration to eliminate LLM synonym drift and state-machine fragmentation across fresh context windows.
  * *Writer Authority & Single-Writer Rule*: **Owner / Lead Orchestrator only.** Worker subagents are strictly read-only consumers and NEVER edit the domain Anchor.
  * *When updated (3 Triggers)*: (1) **Genesis Kickoff** from PRD; (2) **Upfront in Planning** when a Mission declares `Ontology impact` (D27) *before* workers write code; (3) **Domain Refactoring** where the domain Anchor is updated first, then code is renamed to match.
  * *Visual companion*: `atlas/domain-overview.md` is a non-governing projection. Relationships are labelled edges; use `1`, `0..1`, `1..*`, and `0..*` only when cardinality matters.
- `SECURITY.md`: Project-specific isolation, multi-tenancy, secrets, or compliance rules (only if non-standard).
- `GUARDRAILS.md`: Custom AI operational rules (only upon explicit owner request; defaults suffice).
- `PRODUCT.md`: Dedicated commercial/marketing models (only if distinct from repository engineering).

### Modular Capability Contracts
Capability Contracts are small, component-level specifications (`CC-<module>.md`).
A Mission can bind to a primary Contract or coordinate across modular Contracts.
Amend bound Contract gaps/editorial fields through `contract amend`; changing
agreed capabilities requires a `contract_version` bump through the supported
governed workflow. Never edit a bound agreement by hand as ordinary Mission work.
