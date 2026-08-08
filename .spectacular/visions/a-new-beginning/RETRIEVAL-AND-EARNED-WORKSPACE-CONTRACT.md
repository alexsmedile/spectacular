---
type: refactor-foundation-contract
version: 1.0
status: accepted
decision_session: S08
accepted_by: owner
accepted_at: 2026-08-08
central_disposition: accept
upstream:
  - PRODUCT-CONSTITUTION.md@1.0
  - TRUTH-PROVENANCE-FLOOR.md@1.0
  - SUCCESS-EVIDENCE-CONSTITUTION.md@1.0
  - PRODUCT-TRUTH-CONTRACT-MODEL.md@1.0
  - WORK-UNIT-LIFECYCLE-CONTRACT.md@1.0
  - EXECUTION-AUTHORITY-CONTRACT.md@1.0
  - EVIDENCE-CLOSURE-CONTINUITY-CONTRACT.md@1.0
  - RESPONSIBILITY-PLACEMENT-CONTRACT.md@1.0
next_session: S09
---

# Retrieval and Earned Workspace Contract

This accepted S08 contract defines semantic context classes, authority-safe
projections, earned growth, and retrieval tiers. It does not choose public
terms, paths, commands, schemas, physical representation, caching, or
implementation.

## Scope-relative cold-start context

“Always loaded” is relative to the entered operational scope, never a single
universal preload. Every operation loads a very small authoritative set of
universal operating invariants. Project entry then loads concise authoritative
project anchors and a compact project overview. Mission selection or resumption
loads bounded authoritative Mission anchors.

Current-job context is selected from the relevant contracts, repository facts,
tests, tools, constraints, and evidence requirements. Claim evidence and human
rationale are loaded only for validation, repair, conflict resolution,
historical explanation, or audit. Retrieval stops at the shallowest layer that
supports an authority-correct, high-quality next action and provides a
deterministic route to deeper canonical context when insufficient.

| Context class | Authority | Trigger | Purpose |
|---|---|---|---|
| Universal invariants | Authoritative operating constraints | Every operation | Permitted action, stops, authority limits, forbidden inference |
| Project anchors | Authoritative project context | Project entry | Identity, accepted direction, boundaries, critical constraints, current-truth location |
| Project overview | Non-authoritative projection | Orientation | Current state, relationships, Gaps/conflicts, drill-down |
| Mission anchors | Authoritative bounded work context | Mission selection/resume | Outcome, Objective, scope, authority, dependencies, evidence, safe continuation |
| Current-job context | Mixed claim-scoped sources | Planning/execution need | Implementation, contracts, tests, tools, constraints, checks |
| Claim evidence | Attributable evidence | Validate/assess/repair/audit | Support, freshness, method, limits, contrary evidence |
| Human rationale | Attributed history/explanation | Ambiguity/audit/consequence review | Alternatives, rejected choices, design history |

## Authority and projection discipline

Canonical anchors, contracts, records, decisions, receipts, and owner
dispositions retain their own claims. Deterministic registries own only
operating metadata: supported semantic classes, validation, locations,
operation characteristics, and mechanical handling.

Overviews, summaries, indexes, status views, diagrams, relationship maps,
compiled context, and generated documentation are non-authoritative
projections. They expose source scope, generation/freshness basis, omissions,
Gaps, conflicts, unknowns, and direct drill-down. A stored projection is a
dated snapshot. It may propose navigation or a safe next action only while
exposing the underlying canonical state and governing rule.

Projections never own accepted intent, current Capability truth, Mission
lifecycle, authorization, evidence sufficiency, owner disposition,
reconciliation, or resolution of unknowns. Repair targets the canonical source
through its authorized procedure or repairs the projection mechanism, then
regenerates the projection. A projection is never hand-edited to override
canonical material.

## Earned workspace growth

The initial semantic floor is workspace identity/format information; project
anchors for intent, boundaries, critical constraints, and current truth or
explicit unknown; project-specific operating invariants that materially
constrain work; deterministic metadata to locate/validate those anchors; and a
generated orientation view from canonical sources.

Mission, Run, evidence, Gap, decision, archive, historical, and optional
specialist material is not eagerly created. A record or collection is created
deterministically and idempotently on its first qualifying write. A separate
material class earns creation only when it answers a distinct durable owner
question, has independent authority/lifecycle/validation, must cross Missions
without duplicate truth, or embedded treatment demonstrably causes conflict or
failed recovery. An explicit owner request for a demonstrated project need also
earns it. Kits, categories, raw counts, theoretical availability, and aesthetic
completeness do not.

Every material class communicates an absence state: present/current, supported
but unused, required but missing, present but stale, conflicting, or unknown.

## Retrieval tiers and deferrals

| Semantic tier | Owns | Load rule |
|---|---|---|
| Workflow guidance | Judgment sequence, gates, stops, escalation, return | Selected work kind or phase |
| Deterministic operating rules | Mechanical inputs/effects, schema, validation, failure | After semantic intent selects operation |
| Canonical contracts | Direction, behavior, boundaries, authority, Mission/current truth | Affected claim and scope |
| Projections | No canonical claim; orientation and drill-down | Orientation/navigation only |
| Claim evidence | Attributable observation/proof | Planned, executed, assessed, repaired, audited claim |
| Human rationale | Attributed historical explanation | Ambiguity, maintenance, conflict, audit |

Deterministic capabilities are authoritative for mechanical facts that code can
validate or explain. Model instructions retain irreducible judgment. Generated
syntax documentation and tests derive from or validate against one owner of an
operating fact; deterministic mechanisms cannot redefine project authority,
evidence sufficiency, or lifecycle acceptance.

Semantic/vector/graph-database, hosted, and provider-dependent retrieval are
deferred from core. A deterministic relationship map from explicit canonical
links is permitted as a non-authoritative projection. An advisory adapter earns
an experiment only after deterministic retrieval is behaviorally validated and
representative scenarios show repeated material misses. It returns candidate
sources with provenance/freshness, never authority, and core inspection, repair,
resume, and audit never depend on it.

## Validation rule

Validate the contract through cold new-Mission, interrupted-Mission, and
material-claim-audit scenarios. Diagnostics may include loaded/total context,
calls, full/repeated reads, and time to first justified action, but must pair
with correct authority and orientation, conflict/unknown exposure, safe cold
recovery, owner comprehension, source drill-down/repair, and high-quality task
output. A smaller surface is a regression if it hides constraints, increases
unsafe inference, weakens recovery, degrades output, or creates a provider
requirement.
