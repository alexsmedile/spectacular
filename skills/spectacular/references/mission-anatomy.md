# Mission anatomy

Use this when: Agent authoring, splitting, or auditing detailed Spectacular record structures and fields.

The full field lists. `SKILL.md` carries the summary; this file carries the
inventory. Load it when you are writing or auditing a Mission record, not on every
session.

## Shape

A real Mission file (`M7-render-derived-state.md`), trimmed to one of each thing:

```yaml
---
type: Mission
id: 01a00af1-38c0-7268-9529-5856afc7b2f2   # UUIDv7, generated, never hand-written
ref: M7                                     # human navigation; sub-records are M7/O1, M7/R1
title: Render derived state and validate the Proposal record
status: active
activation:
    at: "2026-08-16T14:19:42Z"
    by: Alex
    fingerprint: sha256:ef7e695f...        # over the frozen envelope only
authority:
    operator: [inspect, edit-in-scope, run-checks, bounded-repair, commit-local]
    requires_owner: [activate-mission, change-outcome-or-completion, expand-scope, push]
baseline:
    branch: m7-derived-state
    commit: 127dac140467a462c3810c85c9ca325c18278a14
contract:
    ref: Contract:01a00aae-8921-7b27-96a9-1a4c175e7dc6
    fingerprint: sha256:1ffd39b4...
completion:
    - claim: drift-flags
      pass_boundary: ...      # what must be observably true
      proof_requirement: ...  # what would demonstrate it
objectives:
    - id: 01a00af1-38c0-72bf-a3b2-0c8f1595b50d
      claims: [drift-flags]   # every Objective serves at least one claim
      outcome: ...
dependencies:
    - M6 completed with independent review and owner acceptance.
gaps: []
---

# Origin, rationale, detailed plans, review instructions
```

Live examples are in `.spectacular/missions/*/*.md`. Read one before
authoring a new Mission by hand.

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

When you split, keep the same UUID and ref, and leave one pointer where the inline
detail was. The root Mission record stays the index. Do not add a Mission-local `index.md`.

## Checkpoints

Plan optional checkpoints in the Run body when a Run needs a named progress,
verification, or resume gate. A checkpoint itself does not grant authority or
require human review. When it produces a decision, observation, verdict, or
handoff, create the corresponding Decision, Evidence, Review/Assessment, or
Handoff record and link it from the Run-body note. See
[execute.md](execute.md) for the routing table and template.

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
- `VOCABULARY.md`: Canonical domain ontology and ubiquitous language. Its glossary index is alphabetical; for the detailed section skeleton, see the body shape in [genesis-examples.md](genesis-examples.md).
  * *Threshold*: If <= 3-4 simple entities with no ambiguous terms or shared rules, keep them inline in `PROJECT.md`.
  * *Earned triggers*: (1) Synonym collision / naming ambiguity (e.g. `User` vs `Account`, `Job` vs `Task`); (2) Non-trivial state machine invariants (e.g. `DRAFT` -> `ACTIVE` -> `REVIEW`); (3) Relationships, permissions, or actions that span several concepts; (4) Bespoke non-standard concepts (e.g. `Anchor`, `Gap`, `Handoff`); (5) Multi-contract shared models.
  * *Visual companion*: `atlas/domain-overview.md` is a non-governing projection. Relationships are labelled edges; use `1`, `0..1`, `1..*`, and `0..*` only when cardinality matters.
- `SECURITY.md`: Project-specific isolation, multi-tenancy, secrets, or compliance rules (only if non-standard).
- `GUARDRAILS.md`: Custom AI operational rules (only upon explicit owner request; defaults suffice).
- `PRODUCT.md`: Dedicated commercial/marketing models (only if distinct from repository engineering).

### Modular Capability Contracts
Capability Contracts are small, component-level specifications (`CC-<module>.md`). A Mission can bind to a primary contract or coordinate across multiple modular contracts, editing them as ordinary Mission work when observable capabilities change.
