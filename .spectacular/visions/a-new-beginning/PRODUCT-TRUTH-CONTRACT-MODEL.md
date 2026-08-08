---
type: refactor-foundation-contract
version: 1.0
status: accepted
decision_session: S03B
accepted_by: owner
accepted_at: 2026-08-08
central_disposition: accept
upstream:
  - PRODUCT-CONSTITUTION.md@1.0
  - TRUTH-PROVENANCE-FLOOR.md@1.0
  - SUCCESS-EVIDENCE-CONSTITUTION.md@1.0
next_session: S04
---

# Product Truth and Contract Model

This is the accepted S03B model for distinguishing intended behavior,
implementation facts, runtime observations, evidence, decisions, history,
named gaps, and derived projections. It defines the Minimum Capability
Contract model and the present relationship ceiling. It does not choose work
units, execution authority, evidence sufficiency, storage, public commands,
migration, or implementation.

## 1. Claim-scoped authority

No source is globally authoritative. Every consequential claim declares the
question it answers, its scope, and the authority permitted to answer it.

| Claim class | Authority for that question |
|---|---|
| Accepted behavior | Latest accepted, in-scope Capability Contract or owner decision |
| Implementation fact | Cited code, schema, revision, and fresh scoped checks when applicable |
| Runtime observation | Attributable observation scoped to environment and time |
| Evidence result | Attributable receipt or pointer describing what a named method checked or observed |
| Decision record | Attributable record of what was accepted or rejected, by whom, and within what scope |
| Historical record | Attributable account of an earlier event, belief, or rationale |
| Recommendation, inference, or assumption | Explicitly non-authoritative input until accepted by the appropriate authority |

One authority cannot silently answer another authority domain's question. Code
does not silently amend accepted behavior, and an accepted contract does not
override facts about what code implements or an environment exhibited.

A mismatch is exposed without selecting a convenient winner:

> Accepted behavior requires X. Revision Y implements Z. Environment E
> exhibited W. Conformance check C reported Q.

Material conflict stops consequential decisions, acceptance, lifecycle
promotion, and external effects until the responsible authority resolves it.
Demonstrably unaffected inspection, evidence collection, option preparation,
and recovery planning may continue.

## 2. Capability Contract and MCC

- A **Capability** is an end-to-end ability of the product to produce a defined,
  externally observable outcome for an actor under declared conditions.
- A **Capability Contract** is one versioned, owner-accepted document describing
  the intended behavior of one capability.
- A **Capability Contract Proposal** is that document before acceptance and is
  not authoritative accepted behavior.
- The **Minimum Capability Contract model (MCC)** is the minimum information
  every Capability Contract contains.

The document is the Capability Contract; the MCC is its minimum schema. A
Capability Contract may add advanced sections while remaining MCC-compliant.

## 3. MCC purpose and limits

The MCC lets a cold human or model determine:

- the intended outcome;
- when the behavior applies;
- what behavior is required or prohibited;
- what happens when success is impossible;
- which persistent information and continuity conditions matter;
- how conformance can be checked;
- who accepted the claims and how current they are;
- where related authoritative material can be found.

The MCC does not describe the complete implementation, prove current
conformance, report production behavior, assign work or execution authority,
define evidence sufficiency or closure, replace other authority classes, or
contain task plans and milestones by default.

## 4. Required MCC envelope

Every Capability Contract contains:

1. **Document purpose** — why the document exists and which capability it governs.
2. **Intended outcome** — the externally observable result the capability is meant to provide.
3. **Applies when** — actors and conditions under which required behavior applies.
4. **Does not apply when** — conditions outside the contract's coverage.
5. **Does not provide** — nearby responsibilities or outcomes readers might otherwise assume.
6. **Required behavior** — normative clauses using:
   - **MUST:** required for conformance;
   - **MUST NOT:** prohibited;
   - **MAY:** permitted but not required;
   - **SHOULD:** expected unless a documented reason justifies deviation.
7. **Operating cases** — starting conditions, trigger, required response,
   observable result, and relevant failure handling.
8. **Persistent information and continuity rules** — information surviving the
   current interaction, applicable invariants, or an explicit statement that
   none is contract-relevant.
9. **Conformance checks** — checks or observations mapped to required behavior
   or invariants. They do not establish evidence sufficiency, authorize work,
   accept a result, or close work.
10. **Authority, provenance, and freshness** — acceptance identity, scope,
    version, time, lineage, and the basis for revalidation.
11. **Related authoritative material** — pointers to relevant decisions,
    policies, implementation facts, observations, evidence, history, and
    technical documents. Links locate other authorities; they do not copy or
    transfer authority.

“Required behavior” replaces “guarantee,” “operating case” replaces “behavior
path,” “persistent information and continuity rules” replaces generic “state,”
and “conformance check” replaces “verification obligation.”

## 5. Minimum plus advanced

Spectacular provides one reusable Capability Contract template. The minimum
sections always apply. Advanced sections are added only when earned, such as
for complex persistence, security or privacy, compatibility, performance,
external interfaces, operational constraints, or unusually complex failure
behavior. Advanced sections do not create new contract types.

## 6. Document types are not contract types

Different documents may answer different questions:

- Capability Contract — accepted intended behavior;
- decision record — what was accepted or rejected, by whom, in what scope, and why;
- policy — an accepted cross-cutting constraint;
- evidence or observation record — what a method or observer established;
- historical record — what occurred or was believed in the past;
- handoff record — a downstream candidate for bounded instructions, authority,
  context, receiver, and return conditions for one dispatch.

These do not become Capability Contract subtypes merely because they are
linked. S04 owns a handoff's work-unit relationship, S05 owns execution
authority, and S08 owns storage and layout. S03B preserves the requirement but
does not choose its final name, folder, schema, or lifecycle.

Existing phase-policy machinery is neither expanded nor removed here because
its usage and comparative value remain unknown.

## 7. Capability and architecture remain distinct

Capabilities own accepted observable behavior. Component, interface,
persistence, and policy material describe technical responsibilities,
boundaries, continuity, or cross-cutting constraints supporting that behavior.

The present model introduces no separate component, interface, state, or
policy contract registries. Technical information stays in advanced Capability
Contract sections or ordinary linked documents.

Structural separation is a later reversible Type-2 proposal only when observed
evidence shows one or more of these conditions:

- copied claims drift across Capability Contracts;
- a shared interface changes independently and breaks several capabilities;
- a cold reader cannot recover a material constraint;
- unrelated Capability Contracts require coordinated edits;
- the information needs a genuinely separate acceptance or freshness cycle.

Promotion is evaluated, not automatic.

## 8. Lean relationship metadata

Frontmatter is the lean filtering and relationship surface. It carries compact
metadata needed to select and interpret documents, such as identity, document
class, disposition, version, authority, scope, freshness, supersession,
provenance, and relationships. The body carries complete human-readable
claims; frontmatter does not duplicate them.

A frontmatter field may be authoritative metadata only when its document
contract explicitly declares that role. Exact names, schema, validation, and
storage belong to S08 and S09.

## 9. Present graph ceiling: small Level 1

The present ceiling is a deliberately small Level 1:

- Capability Contracts use the minimum-plus-advanced model;
- other questions use ordinary document types;
- frontmatter exposes lean relationship metadata;
- retrieval follows pointers to authoritative material;
- generated indexes and decision maps remain projections.

Level 1 introduces no graph database, graph compiler, GraphRAG, relationship
registry, component/interface/state contract taxonomy, or executable graph
transaction model.

Level 2 requires demonstrated recurring independent contract boundaries and a
measured failure of Level 1. Level 3 additionally requires a concrete problem
that ordinary documents, metadata, links, and deterministic checks cannot
solve, plus proportional evidence of owner benefit and acceptable maintenance,
comprehension, recovery, compatibility, context, and attention costs.

## 10. Decision history map

A decision history map is a non-authoritative projection showing how accepted,
rejected, derived, and superseded decisions contributed to later decisions.
“Map” is preferred to “tree” because one decision may depend on several
predecessors.

Decision records remain authoritative. The map may display derivation,
dependency, rejection, supersession, conflict, and missing links, but it may
not accept, reject, supersede, renew, or reconcile decisions.

## 11. Named gaps and bounded continuation

Missing information is explicitly classified as:

- **Unknown:** information not currently available.
- **Assumption:** a temporary stated belief used for bounded work.
- **Question:** an answer required from an authorized person.
- **Discovery task:** bounded investigation needed to obtain evidence.

S04 decides whether any becomes a separate durable work record.

Consequential work stops when the answer could materially affect intended
outcome, required behavior, authority, safety or rights, compatibility,
conformance, acceptance, or an external effect. Unaffected inspection, option
preparation, evidence collection, and recovery planning may continue. Missing
contract information is never filled through silent inference.

## 12. Projection and drill-down contract

Every projection must identify its source documents and generation time,
expose relevant freshness, named gaps, and conflicts, drill down to the
authoritative record, and remain visibly non-authoritative.

A projection may not accept or reject a claim, supersede authority, renew
freshness, hide disagreement, or perform reconciliation. Live generation
improves currency but does not create authority.

## 13. Reconciliation boundary

Reconciliation means resolving a material disagreement through the authority
responsible for the affected claim and updating the appropriate authoritative
material.

S03B defines only the vocabulary and authority boundary. S05 and S06 retain
evidence sufficiency, execution permissions, owner and provider gates, closure,
reconciliation mechanics, and lifecycle promotion.

## Exit condition

Every consequential claim has one declared authority for its question and
scope. Every derived projection can drill down to authoritative material while
exposing freshness, named gaps, and conflict.
