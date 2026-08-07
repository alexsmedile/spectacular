---
type: concept-piece
id: PZL-169
status: captured
domain: design-authority
sources: [source-014]
source_authority: user-provided-unsourced-proposal
assessment: strong
evidence_status: partial
disposition: pending
depends_on: [PZL-064, PZL-085, PZL-086]
overlaps_with: [PZL-066, PZL-083, PZL-087]
conflicts_with: []
tags: [product-contract, architecture-decision, prd, adr, decision-order]
updated: 2026-08-07
---

# Product contract before architecture choice

## Core message

Agree on the outcome, business rules, non-goals, observable behavior, and material constraints before
locking the implementation mechanism that claims to satisfy them.

## Value

Prevents architecture from silently redefining product intent and lets multiple technical options be
judged against the same accepted problem contract.

## Assumptions

- Product and architecture concerns can be linked without forcing every technical detail into a PRD.
- Urgent discovery may precede product lock but cannot masquerade as the accepted solution.

## Evidence and collisions

Existing Mission-delta and authority concepts already separate human outcomes from technical
synthesis. Source 014 frames this as PRD before ADR, but its PRD template risks becoming a mixed
authority containing UX, schemas, dependencies, KPIs, SLOs, and implementation detail.

## Trade-offs and recommendation

Adopt the ordering, not the mega-document. Keep intent compact; link capability contracts, decision
records, schemas, and evidence as independently governed artifacts when complexity earns them.

## Decision

Pending.
