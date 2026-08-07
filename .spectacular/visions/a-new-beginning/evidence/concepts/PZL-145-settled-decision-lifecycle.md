---
type: concept-piece
id: PZL-145
status: captured
domain: decision-governance
sources: [source-011, source-014]
source_authority: owner-authored-proposal-corpus
assessment: mixed
evidence_status: partial
disposition: pending
depends_on: [PZL-086]
overlaps_with: [PZL-076, PZL-100]
conflicts_with: []
tags: [adr, lifecycle, proposed, rejected, deprecated, question]
updated: 2026-08-07
---

# Settled-decision lifecycle

## Core message

A decision record should represent an actual resolution; unresolved choices stay questions or
proposals, while rejected and deprecated outcomes remain durable history.

## Value

Keeps the decision corpus trustworthy and prevents proposals from masquerading as authority.

## Assumptions

The system can distinguish proposal capture from authorization without losing discussion context.

## Evidence and collisions

Issue #39 asks whether ADR states should include proposed, rejected, and deprecated. The open choice
is whether “accepted” is explicit or simply the normal settled state.

## Trade-offs and recommendation

Keep the lifecycle small: unresolved as QUE/proposal; settled current; rejected; deprecated or
superseded. Avoid a workflow engine for decisions.

## Decision

Pending.
