---
type: concept-piece
id: PZL-144
status: captured
domain: metadata-contract
sources: [source-011]
source_authority: owner-authored-proposal-corpus
assessment: promising
evidence_status: partial
disposition: pending
depends_on: [PZL-005, PZL-123]
overlaps_with: [PZL-056]
conflicts_with: []
tags: [dublin-core, metadata, pageworks, adapter, provenance]
updated: 2026-08-07
---

# Standard metadata adapter mapping

## Core message

Map Spectacular's stable document metadata to external standards at an adapter boundary instead of
renaming the internal lifecycle schema around a general-purpose vocabulary.

## Value

Supports interoperability and companion skills while preserving domain-specific internal meaning.

## Assumptions

A real exporter, indexer, or documentation consumer needs the mapping.

## Evidence and collisions

Issue #38 proposes Dublin Core. This is especially relevant to pageworks or a shared doc-engine
handoff, but no current consumer proves that every field belongs in core frontmatter.

## Trade-offs and recommendation

Inventory existing fields and define a one-way mapping table. Implement only alongside a concrete
consumer; do not broaden the default schema speculatively.

## Decision

Pending.
