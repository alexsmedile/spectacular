---
type: concept-piece
id: PZL-055
status: captured
domain: metadata-authority
sources: [source-005]
source_authority: code-audit-proposal
assessment: promising
evidence_status: partial
disposition: pending
depends_on: [PZL-051, PZL-058]
overlaps_with: [PZL-002, PZL-003, PZL-056]
conflicts_with: []
tags: [registry, commands, aliases, classification, handlers]
updated: 2026-08-07
---

# Command registry authority

## Core message

Define command path, aliases, read/write class, schema gate, help, and handler in
one command registry that drives executable dispatch and its projections.

## Value

Could prevent dispatcher/help/docs/test drift and make the public surface auditable.

## Assumptions

- Bash 3.2 can consume or generate the registry reliably.
- Registry data can own dispatch without becoming parallel metadata.

## Evidence and collisions

Current command docs and dispatch disagree. This registry is distinct from the
document-rules registry in PZL-003, though combining them could create an overly
broad authority.

## Trade-offs and recommendation

One authority and generated checks versus schema machinery and central coupling.
Prototype the smallest representation against real dispatch before adoption.

## Decision

Pending.
