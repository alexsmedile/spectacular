---
type: concept-piece
id: PZL-122
status: captured
domain: semantic-contract
sources: [source-009, source-010]
source_authority: unsourced-expanded-synthesis
assessment: strong
evidence_status: supported
disposition: pending
depends_on: [PZL-107, PZL-110]
overlaps_with: [PZL-083]
conflicts_with: []
tags: [ontology, inference, validation, shacl, schema]
updated: 2026-08-07
---

# Separate semantics, inference, and validation

## Core message

Model meaning, derive implications, and enforce closed constraints as distinct functions;
choose the lightest representation and validator that satisfies each one.

## Value

Prevents an ontology from being mistaken for complete runtime policy enforcement.

## Assumptions

- Important invariants have explicit ownership and failure behavior.
- Technology is chosen after the required reasoning and validation semantics are known.

## Evidence and collisions

RDF represents graph data and OWL supplies formal ontology semantics; SHACL explicitly
validates a data graph against shapes. Missing axioms or different open/closed-world
assumptions can still cause missed or surprising results.

## Trade-offs and recommendation

Formal semantic tooling can improve interoperability while adding a specialized stack.
Start with Markdown frontmatter schemas and deterministic checks; test RDF/OWL/SHACL only
when a concrete relationship or validation problem exceeds them.

## Decision

Pending.
