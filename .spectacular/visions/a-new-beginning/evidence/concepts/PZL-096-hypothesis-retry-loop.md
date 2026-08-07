---
type: concept-piece
id: PZL-096
status: captured
domain: diagnostic-loop
sources: [source-007]
source_authority: synthesized-proposal
assessment: strong
evidence_status: supported
disposition: pending
depends_on: [PZL-072]
overlaps_with: [PZL-093, PZL-095]
conflicts_with: []
tags: [retry, hypothesis, evidence, checkpoint, blocker]
updated: 2026-08-07
---

# Hypothesis-driven retry budget

## Core message

Allow another repair attempt only when it adds a new hypothesis or evidence; on
budget exhaustion, preserve failure, commands, hypotheses, attempts, checkpoint, and next investigation.

## Value

Prevents blind guess loops while making failed work resumable and informative.

## Assumptions

- Hypotheses and attempts can be recorded concisely.
- Retry budgets can vary by risk and cost.

## Evidence and collisions

The senior-harness attachment recommends a bounded three-attempt self-correction loop;
Source 007 adds a substantive-progress criterion and richer failure packet.

## Trade-offs and recommendation

Disciplined diagnosis versus trace overhead. Adopt the new-evidence rule and compact
failure handoff; do not hard-code one universal attempt count.

## Decision

Pending.
