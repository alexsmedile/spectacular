---
type: Anchor
id: 01a030b4-6159-7a6a-b77b-e2466a25469b
human_ref: ONTOLOGY
title: Domain ontology and ubiquitous language
direction: Give owners and agents one readable model of Spectacular's concepts, relationships,
  rules, and actions.
created: '2026-08-24T01:05:49+02:00'
version: '0.4'
updated: '2026-10-07T20:07:46Z'
---

# Domain vocabulary and relationships

This Anchor defines project terms. [Domain overview](atlas/domain-overview.md)
is its visual companion; [governed execution](atlas/governed-execution-model.md)
holds the detailed mechanical object model. Soft and governed meanings are
explicit; a type alone does not select governance.

## Glossary index

| Term | Definition | Context |
| --- | --- | --- |
| A La Carte | Modular adoption of individual Spectacular surfaces without adopting the full mission lifecycle. | Operational mode |
| Anatomy | Optional explanation of artifact parts and relationships; a reasoning guide, not a lifecycle. | Drafting method |
| Anchor | Durable project truth about a named concern. | Project definition |
| Architectural HUD | The visual and topological map (Atlas + Campaign DAG) giving operators an instant mental model of system boundaries. | Planning and navigation |
| Assessment | Qualitative evaluation of architectural posture, system maturity, or technical debt; captured in Reviews or Retrospectives without rigid binary gates. | Proof and continuity |
| Atlas | Maintained objects, relationships, and lasting explanations; diagrams are optional. | Durable context |
| Atlas Coverage | Verification that generated code implements all state transitions and entity boundaries declared in Atlas maps. | Proof and continuity |
| Audit | Inspection of a stated scope and revision; independent findings may live in audits/. | Inspection |
| Blast Radius | The surface area of files, modules, and dependencies impacted by an execution turn; bounded by authorized paths. | Operational safety |
| Capability Contract | A modular specification of an observable capability. | Project definition |
| Context Amnesia | The loss of settled architectural nuances or review findings across fresh model windows; solved by durable Git Markdown records. | Operational continuity |
| Contract | Accepted agreement about a capability and its constraints. | Project definition |
| Decision | Choice and rationale; soft by default, mechanical validation recommended. | Durable context |
| Decision Compliance | Verification that an implementation strictly adheres to locked architectural rulings (`D<N>`). | Proof and continuity |
| Durable Context | Project knowledge that remains useful across chats, agents, and working phases. | Continuity |
| Evidence | Attributable observation supporting a claim. | Proof and continuity |
| Gap | A stated unresolved limitation or dependency. | Proof and continuity |
| Generation Velocity | The rapid rate at which frontier models produce code, requiring structural containment and bounded charters. | Governed execution |
| Handoff | A bounded transfer of work context between operators. | Proof and continuity |
| High-Stakes Code | Consequential changes (auth, crypto, payments, zero-downtime cutovers) requiring strict review and attributable evidence. | Operational mode |
| Maps | Explanatory visual projections (`.spectacular/atlas/`) that navigate domain models without enforcing schema authority. | Planning and navigation |
| Mission | A frozen execution envelope with authority and proof boundaries. | Governed execution |
| Objective | An outcome-sized claim within a Mission. | Governed execution |
| Owner | Person accountable for consequential direction and acceptance. | Governed execution |
| Plan | Intended outcome, approach, progress, and result; durable project copy in plans/. | Ordinary work |
| Proof Validity | Verification that test receipts and validation runs are authentic, reproducible, and exit with code 0. | Proof and continuity |
| Proposal | Mutable exploration that may inform later accepted work. | Planning |
| Requirement | Independently useful need, outcome, constraint, and observable acceptance; optional requirements/, existing specs/ homes remain valid. | Durable context |
| Retrospective | Freeform milestone post-mortem and reflection (`.spectacular/retrospectives/`) capturing lessons learned and forward recommendations without execution ceremony. | Planning and continuity |
| Review | Governed evaluation of claims; normally a Mission-owned reviews/ record. | Governed proof |
| Routine Fast Code | Authorized everyday work checked against the actual change; no mandatory commit or Mission. | Ordinary work |
| Run | A mutable attempt to advance one Objective or Mission. | Governed execution |
| Schema | Field and constraint definition; a frontmatter schema claim requires command-owned mechanical validation. | Structure and validation |
| Skeleton | Optional starting outline, evolved into the actual artifact without a redundant companion file. | Drafting method |
| Spec | Reusable description of intended behavior or interfaces; a soft file in optional specs/. | Durable context |

## Bounded contexts

- Durable context: Anchors, raw, Atlas, Decisions, plans, and optional requirements/specs/audits.
- Drafting methods: optional skeletons and anatomies, distinct from enforced schemas.
- Governed execution: Contracts, Missions, Objectives, Runs, authority, and proof.
- Product documentation: root docs/ explains shipped behavior to human readers.

## Objects

### Entity: Anchor
Named current project truth. Soft Anchors need basic metadata and no UUID.
This project preserves its historical governed Anchor identities. Only create
an additional Anchor when its concern has independent use.

### Entity: Decision
Soft choices retain rationale and supersession. Existing governed Decisions keep
UUIDs, numbered references, and immutable history; validation is recommended for
new consequential choices but does not activate a Mission.

### Entity: Contract
Mechanically validated agreement with versioned capabilities. Amend gaps/editorial
fields through the supported command; agreed behavior requires a version change.

### Entity: Plan
A stable work document in plans/, reused through planning, implementation, and
results. Native plan mode controls actions; Spectacular provides durable storage.

### Entity: Requirement
A linked need with scope and observable acceptance. Its home is optional
requirements/ or an existing specs/ file; retain one owner. A short PRD overview
links needs rather than duplicating them. See [requirements](../skills/spectacular/references/requirements.md).

### Entity: Proposal
A formally tracked open question with existing governed metadata and retirement
rules. Ordinary exploration can stay in raw or a plan.

## Relationships and actions

Plans reference Requirements, Anchors, Atlas objects, Decisions, and specs. Specs
describe behavior satisfying needs; plans sequence delivery and verification.
Prototypes explore uncertain needs or approaches; they do not certify acceptance.
[Skeletons and anatomies](../skills/spectacular/references/drafting-methods.md) aid
drafting only when useful. Optional Missions bind Contracts and own proof. See the
[mechanical object model](atlas/governed-execution-model.md) for detailed relations.

Relationships are directed labelled edges; use owns, contains, references,
governed_by, reads_from, writes_to, emits, and transitions_to where meaningful.
Use 1, 0..1, 1..*, or 0..* only when cardinality matters.

## Invariants and policies

- Ordinary work reads relevant linked context and does not require a lifecycle.
- Folder agreements define soft metadata conventions; schema claims mean mechanical validation.
- Requirements and drafting methods add no activation, lifecycle, or validation gate.
- Owner instructions and accepted constraints apply in both routes.
- Audit findings, a governed Review, and an owner acceptance are distinct facts.
- Guardrails are project context, not a claim that preflight mechanically certifies their prose.
- Ontology-impact declarations remain guidance, not a typed Mission field.

Historical Decisions D25–D30 remain unchanged. Current ordinary-work policy is
stated in [AGENTS.md](../AGENTS.md) and [PROJECT.md](PROJECT.md).
