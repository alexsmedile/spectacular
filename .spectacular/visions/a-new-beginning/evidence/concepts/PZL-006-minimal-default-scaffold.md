---
type: concept-piece
id: PZL-006
status: captured
domain: initialization
sources: [source-001, source-003, source-006]
source_authority: proposal
assessment: mixed
evidence_status: partial
disposition: pending
depends_on: []
overlaps_with: [PZL-007, PZL-008, PZL-009, PZL-063, PZL-077]
conflicts_with: [PZL-077]
tags: [init, scaffold, onboarding, substrate]
updated: 2026-08-07
---

# Minimal default scaffold

## Core message

Make the default initialized workspace contain only the substrate that pays for
itself immediately, proposed as PRD, config, and an empty requests directory.

## Value

Reduces minute-one burden and makes the directory tree reflect actual project use
instead of Spectacular's maximum capability surface.

## Assumptions

- Missing AGENTS, POLICY, and specs do not weaken safe cold-agent operation.
- Current lifecycle commands can work when directories and indexes are absent.
- Config can be machine-owned or substantially simplified.

## Evidence and collisions

Current blank init creates PRD, specs/index, config, AGENTS, POLICY, requests, and
specs. The proposed floor conflicts with spec-first execution, user-configurable
policy/tool overrides, and the current cold-agent instruction source. The “three
files” description is actually two files plus one directory. Source 006 proposes
a different minimum—PROJECT, SYSTEM, capabilities, missions, and archive—showing
that “minimal” depends on the protected product loop rather than file count alone.

## Trade-offs and recommendation

Honest onboarding versus hidden implicit behavior and weaker discoverability.
Define the minimum operational invariants before choosing artifacts. Mixed; the
direction is valuable but the exact floor is not ready.

## Decision

Pending.
