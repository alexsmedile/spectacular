---
type: concept-piece
id: PZL-004
status: captured
domain: reference-retrieval
sources: [source-001, source-003, source-004]
source_authority: proposal
assessment: promising
evidence_status: partial
disposition: pending
depends_on: [PZL-003]
overlaps_with: [PZL-001, PZL-005]
conflicts_with: []
tags: [references, filesystem, retrieval, migration]
updated: 2026-08-07
---

# Tiered reference layout

## Core message

Use filesystem paths such as `workflows/`, `engines/`, `rules/`, and `contracts/`
as routing metadata instead of keeping every reference in one flat directory.

## Value

Improves human navigation, signals loading intent before a file opens, and makes
the runtime-reachable reference surface easier to audit.

## Assumptions

- The proposed categories are stable and mutually understandable.
- Path depth helps retrieval more than it complicates links and packaging.

## Evidence and collisions

The current directory contains 77 flat Markdown files. Moving them affects wiki
links, tests, plugin packaging, snapshots, and external references. Some files may
span categories, and a forced taxonomy can hide rather than remove complexity.

## Trade-offs and recommendation

Clear navigation versus migration cost and category debates. Create a complete
current-file-to-target map and measure actual reachable files before moving
anything. Promising; exact tiers remain provisional.

## Decision

Pending.
