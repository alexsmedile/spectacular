---
type: concept-piece
id: PZL-060
status: captured
domain: migration-strategy
sources: [source-005]
source_authority: code-audit-proposal
assessment: promising
evidence_status: unverified
disposition: pending
depends_on: [PZL-058]
overlaps_with: [PZL-009, PZL-052]
conflicts_with: []
tags: [compatibility, deprecation, release, v2]
updated: 2026-08-07
---

# Compatibility release before v2 removal

## Core message

Warn on deprecated commands for a defined compatibility release, print exact
replacements, then remove the redundant surface in v2.

## Value

Separates migration assistance from permanent alias retention.

## Assumptions

- One minor release reaches affected users and automation.
- Warnings are visible without breaking machine-readable output.

## Evidence and collisions

The source proposes one minor release but supplies no adoption telemetry or
compatibility policy. Different commands may require different windows.

## Trade-offs and recommendation

Cleaner v2 versus migration cost and temporary dual surface. Define compatibility
criteria per command class rather than accepting one universal duration.

## Decision

Pending.
