---
type: concept-piece
id: PZL-003
status: captured
domain: metadata-authority
sources: [source-001, source-003, source-004, source-005]
source_authority: proposal
assessment: mixed
evidence_status: partial
disposition: pending
depends_on: []
overlaps_with: [PZL-002, PZL-004, PZL-055]
conflicts_with: []
tags: [registry, metadata, rules, duplication]
updated: 2026-08-07
---

# Single dispatch registry

## Core message

Represent document dispatch metadata once in a compact machine-readable registry
instead of opening many thin rules files.

## Value

Could remove stub repetition, enable validation, and make the supported document
surface cheaply inspectable.

## Assumptions

- Most thin rules files contain data rather than meaningful behavior.
- One representation can serve the skill, CLI, doctor, and templates.
- The format is practical under Bash 3.2.

## Evidence and collisions

There are 23 `*-rules.md` files, but rules frontmatter and `doc-index.md` already
provide registry-like authority. Adding `registry.yaml` without retiring those
sources would worsen duplication. YAML parsing may also cost more than the stubs.
Source 005 proposes a separate command registry; merging document dispatch and
command dispatch could centralize unrelated schemas, while keeping both requires
clear ownership boundaries.

## Trade-offs and recommendation

Compact authority and validation versus a central schema bottleneck and parser
complexity. Inventory which rule files are truly data-only, then choose the sole
authority and generated projections. Mixed; format remains open.

## Decision

Pending.
