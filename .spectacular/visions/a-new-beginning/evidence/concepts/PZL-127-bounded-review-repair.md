---
type: concept-piece
id: PZL-127
status: captured
domain: verification-architecture
sources: [source-010, source-012, source-015]
source_authority: unsourced-expanded-synthesis
assessment: strong
evidence_status: partial
disposition: pending
depends_on: [PZL-114, PZL-116]
overlaps_with: [PZL-096, PZL-110]
conflicts_with: []
tags: [review, repair, retry, termination, escalation]
updated: 2026-08-08
---

# Bounded review-repair cycle

## Core message

Route material review failures back with exact evidence, but terminate through an acceptance
rubric, severity policy, attempt budget, unchanged-finding detection, and escalation packet.

## Value

Preserves adversarial improvement without allowing subjective “entire satisfaction” to loop forever.

## Assumptions

- Reviewer findings are deduplicated and tied to a contract or deterministic failure.
- Stop conditions distinguish accepted risk, unresolved defect, and reviewer disagreement.

## Evidence and collisions

Source 010 proposes an unbounded Gauntlet until the Skeptic is satisfied. That conflicts with
PZL-096's evidence-bearing retry budget and makes an agent's confidence the termination authority.
Source 015 adds primary examples of bounded revision/escalation gates and a five-round review
breaker. Their exact thresholds are implementation choices rather than universal constants.

## Trade-offs and recommendation

Repeated review can catch regressions while adding cost and oscillation. Compose existing
hypothesis retries and independent review into one bounded state machine before any automation.

## Decision

Pending.
