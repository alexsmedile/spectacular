---
type: concept-piece
id: PZL-149
status: captured
domain: capture-boundary
sources: [source-011]
source_authority: owner-authored-proposal-corpus
assessment: strong
evidence_status: partial
disposition: pending
depends_on: [PZL-051, PZL-054]
overlaps_with: [PZL-119, PZL-124]
conflicts_with: []
tags: [capture, proposal, adapters, side-effects, local]
updated: 2026-08-07
---

# Proposal-only capture boundary

## Core message

A generic capture surface may normalize an observation into a local proposal, but external writes
remain explicit provider operations owned by the user or native tool.

## Value

Creates one predictable intake boundary without hiding network mutations behind an agentic verb.

## Assumptions

The local proposal schema can retain provenance and offer clear next destinations.

## Evidence and collisions

Issue #53 proposes a capture surface derived from SPC-007. It converges with mechanical CLI versus
agentic skill ownership and the native-provider mutation boundary.

## Trade-offs and recommendation

Prototype `capture propose` as a side-effect-limited local operation. Do not auto-create GitHub
issues, decisions, or specs; present those as explicit promotions.

## Decision

Pending.
