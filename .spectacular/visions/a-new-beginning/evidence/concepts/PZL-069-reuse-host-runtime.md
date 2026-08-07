---
type: concept-piece
id: PZL-069
status: captured
domain: runtime-boundary
sources: [source-006, source-007]
source_authority: proposal
assessment: strong
evidence_status: supported
disposition: pending
depends_on: [PZL-062]
overlaps_with: [PZL-051, PZL-054, PZL-071, PZL-091, PZL-098]
conflicts_with: []
tags: [runtime, coding-agent, scheduler, platform, boundary]
updated: 2026-08-07
---

# Reuse the host coding runtime

## Core message

Do not build a model host, queue, daemon, or agent platform; express bounded work
through the coding-agent runtime the user already chose.

## Value

Keeps Spectacular focused on durable context, authority, evidence, and reconciliation.

## Assumptions

- Host runtimes expose enough checkpoint, resume, Git, and stop behavior.
- The convention remains portable across runtime vendors.

## Evidence and collisions

This aligns with the PRD's orchestration-platform non-goal, but a shell command
that starts and manages a mission run may itself become a runtime abstraction.
Source 007 proposes several typed runs and agents while still rejecting generalized
scheduling, making the platform boundary more important and less settled.

## Trade-offs and recommendation

Low platform burden versus runtime-specific adapters. Make the durable run contract
canonical and keep launch mechanics thin or skill-owned.

## Decision

Pending.
