---
type: concept-piece
id: PZL-059
status: captured
domain: implementation-architecture
sources: [source-005]
source_authority: code-audit-proposal
assessment: promising
evidence_status: unverified
disposition: pending
depends_on: [PZL-035, PZL-042, PZL-058]
overlaps_with: [PZL-042]
conflicts_with: []
tags: [bash, modules, bundle, distribution, build]
updated: 2026-08-07
---

# Modular Bash source, single artifact

## Core message

Develop the CLI as core, command, and adapter modules, then assemble one Bash 3.2
zero-dependency executable for release.

## Value

Could improve ownership and reviewability without changing installation simplicity.

## Assumptions

- Assembly is deterministic and preserves source locations useful for debugging.
- Module boundaries reduce coupling in shell code.

## Evidence and collisions

The 16,507-line monolith motivates the proposal, but no prototype measures test,
startup, release, or contributor impact. The repo currently has no build step.

## Trade-offs and recommendation

Maintainable sources versus generator/release complexity and bundled-debug friction.
Prototype only after the v2 surface settles; compare against extraction or wrapping.

## Decision

Pending.
