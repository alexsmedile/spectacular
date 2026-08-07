---
type: concept-piece
id: PZL-167
status: captured
domain: decision-method
sources: [source-014]
source_authority: user-provided-unsourced-proposal
assessment: strong
evidence_status: partial
disposition: pending
depends_on: [PZL-115, PZL-155]
overlaps_with: [PZL-031, PZL-048, PZL-160]
conflicts_with: []
tags: [reversibility, blast-radius, decision-effort, one-way-door, evidence]
updated: 2026-08-07
---

# Reversibility-calibrated decision investment

## Core message

Scale decision effort and evidence fidelity to reversal cost, blast radius, migration burden,
external commitments, and consequence rather than applying one ceremony to every choice.

## Value

Prevents analysis paralysis on cheap choices while protecting architecture, security, data, public
contracts, and other decisions that become expensive to unwind.

## Assumptions

- Reversibility can be estimated conservatively before commitment and reassessed as adoption grows.
- Fast decisions still respect separate authority, safety, privacy, and compatibility gates.

## Evidence and collisions

Risk-calibrated review and uncertainty-shaped artifacts already use consequence and reversibility.
Source 014 makes reversal cost the explicit first classification but presents a misleading binary:
an API format or library can change class once consumers, data, or operations depend on it.

## Trade-offs and recommendation

Adopt a small spectrum rather than permanent Type 1/Type 2 labels. Record impact, reversibility,
uncertainty, and evidence threshold; cap effort for low-risk choices and escalate when facts change.

## Decision

Pending.
