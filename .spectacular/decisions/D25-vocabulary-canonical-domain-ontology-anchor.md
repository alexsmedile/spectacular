---
type: Decision
id: 01a030b6-6ed9-7d84-a50d-7a4274b20fcf
title: Make VOCABULARY.md the canonical domain ontology Anchor
created_by: Alex
created: "2026-08-23T22:20:57Z"
updated: "2026-08-23T22:20:57Z"
actor: Alex
actor_role: owner
alternatives:
    - introduce a new Ontology record type and CRUD command
    - make Atlas files authoritative
    - keep only a loose glossary with no defined model shape
disposition: accepted
question: How should Spectacular preserve a project's detailed ontology without introducing a new record type or CLI surface?
rationale: 'A project needs durable, shared definitions of its concepts, relationships, actions, and rules. The existing earned VOCABULARY.md Anchor is the smallest coherent home: it is canonical Markdown, works with current Anchor discovery, and can remain useful without a graph database or a prose validator.'
ref: D25-vocabulary-canonical-domain-ontology-anchor
scope:
    - v2
---
# Make VOCABULARY.md the canonical domain ontology Anchor

`VOCABULARY.md` is an earned Anchor titled **Domain ontology and ubiquitous
language**. It begins with an alphabetical glossary index, then holds detailed
bounded contexts; entities and value objects; directed relationships,
cardinality, and ownership; actions and events; invariants and policies;
implementation mappings; and semantic gaps/change history.

The map legend is fixed: Bounded Context, Actor, Entity, Value Object, Action,
Event, Policy/Invariant, and External System are node types. Relationships are
edges. Default labels are `owns`, `contains`, `belongs_to`, `has`, `references`,
`requests`, `performs`, `emits`, `transitions_to`, `governed_by`, `reads_from`,
and `writes_to`. Cardinality follows standard UML/ER multiplicity: `1`, `0..1`,
`1..*`, and `0..*`.

A documented relationship is not automatically enforced. Its implementation
mapping must name the schema, code invariant, API contract, or test that does
enforce it.
