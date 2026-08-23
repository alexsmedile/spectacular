---
type: Proposal
id: 01a030b4-6159-7d23-a146-61693eac8556
ref: P13
title: Ontology anchors and visual domain maps
status: draft
created_by: Alex
created: "2026-08-24T00:15:00Z"
updated: "2026-08-24T00:15:00Z"
scope:
    - v2
target_contract: Contract:019fe381-5d61-7223-b362-03a5f99a7b10
atlas: ../atlas/domain-overview.md
---

# Ontology anchors and visual domain maps

Spectacular should help a project formulate and maintain a shared domain model,
so an agent grounds its work in explicit concepts, relationships, rules, and
actions rather than inferring them from isolated screens, tables, or prompts.

## Proposed shape

`.spectacular/VOCABULARY.md` becomes an earned project Anchor titled **Domain
ontology and ubiquitous language**. It is the detailed, canonical source for:

1. an alphabetical glossary index;
2. bounded contexts;
3. objects: entities, value objects, identity, attributes, and lifecycle;
4. directed relationships, cardinality, and ownership;
5. actions and events;
6. invariants, policies, and permission or HITL rules;
7. implementation mappings; and
8. semantic gaps and change history.

`.spectacular/atlas/domain-overview.md` is the linked, non-governing visual
projection. It presents the whole domain as a compact graph and may be expanded
by bounded-context or value-slice Atlases when one map no longer remains clear.
The Vocabulary is canonical if a map and definition drift.

## Shared notation

The domain map uses bounded contexts as clusters; Actors, Entities, Value
Objects, Actions, Events, Policies/Invariants, and External Systems as node
types. Relationships are labelled edges, not nodes. The default edge labels
are `owns`, `contains`, `belongs_to`, `has`, `references`, `requests`,
`performs`, `emits`, `transitions_to`, `governed_by`, `reads_from`, and
`writes_to`. Cardinality uses established UML/ER multiplicity: `1`, `0..1`,
`1..*`, and `0..*`.

## Workflow impact

At Genesis, formulate domain nouns before implementation verbs whenever the
project has non-trivial concepts, state, relationships, or policies. A material
Mission later cites the affected Vocabulary concepts in its existing `sources:`
or body, or explicitly records that it has no ontology impact. Changes to public
behaviour, data, permissions, workflow, or lifecycle update the Vocabulary and
the relevant Atlas in the same scope.

## First implementation

Dogfood the convention in Spectacular itself; add the Anchor, its domain overview
Atlas, templates, and guidance. This owner-approved documentation/workflow patch
does not activate a Mission. It records the direction through Decisions and is
subject to adversarial review at the end of this session.

## Non-goals

- No new public CLI command, record type, generic entity CRUD system, graph
  database, RDF/OWL format, or compatibility layer.
- No automatic graph generation or hard prose validator.
- No claim that a documented cardinality is enforced unless its implementation
  mapping identifies a concrete schema, code invariant, or test.
