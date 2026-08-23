---
type: Decision
id: 01a030b6-a873-735f-897d-cfd017c99ae8
title: Use Atlas domain maps as non-governing visual projections
created_by: Alex
created: "2026-08-23T22:21:12Z"
updated: "2026-08-23T22:21:12Z"
actor: Alex
actor_role: owner
alternatives:
    - make the domain graph the canonical ontology
    - add a graph database or formal semantic-web store
    - retain only prose and tables in VOCABULARY.md
disposition: accepted
question: How should Spectacular make a project's ontology easy to scan without creating a second source of authority?
rationale: A graph makes the whole domain legible, but maps are easier to simplify, split, and redraw than detailed definitions. Keeping the Vocabulary canonical and Atlas maps non-governing preserves one source of semantic truth while giving owners and agents a digestible overview.
ref: D26-accepted
scope:
    - v2
---
# Use Atlas domain maps as non-governing visual projections

`atlas/domain-overview.md` is the whole-product visual companion to the
canonical Vocabulary. It renders bounded contexts as clusters and uses the
approved node and edge legend. A project adds Atlas slices only when a bounded
context or value slice needs more legible detail than the global overview can
hold.

The Atlas explains and navigates. It grants no authority, does not create
Mission drift, and carries no schema claim. When its content differs from
`VOCABULARY.md`, the Vocabulary definition is authoritative and the map must be
updated in the same accepted change.

Implementation mappings remain in the Vocabulary or supporting architecture
documents. They may appear in a slice Atlas when useful, but they must not
obscure the primary domain model with every API, table, UI surface, or test.
