---
type: concept-piece
id: PZL-143
status: captured
domain: implementation-quality
sources: [source-011]
source_authority: owner-authored-proposal-corpus
assessment: strong
evidence_status: partial
disposition: pending
depends_on: [PZL-111, PZL-125]
overlaps_with: []
conflicts_with: [PZL-010]
tags: [principles, rubric, quality, anti-slop, context-budget]
updated: 2026-08-07
---

# Repository-specific quality rubric

## Core message

Encode the few design constraints that protect a repository from its observed failure modes, and
evaluate changes against them at relevant gates.

## Value

Makes anti-slop expectations concrete without asking every agent to reread generic engineering lore.

## Assumptions

Rules are tied to repository architecture, evidence, or recurring defects and have a clear scope.

## Evidence and collisions

Issue #37 and Sources 008–010 support classic modularity. An always-loaded generic glossary would
collide with the cold-start budget and mostly duplicate model knowledge.

## Trade-offs and recommendation

Start with a short rubric covering deep boundaries, compatibility, scope, and verification. Add a
rule only when it changes a decision or blocks a demonstrated failure.

## Decision

Pending.
