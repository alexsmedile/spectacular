---
type: concept-piece
id: PZL-048
status: captured
domain: git-safety
sources: [source-005]
source_authority: code-audit-proposal
assessment: strong
evidence_status: supported
disposition: pending
depends_on: []
overlaps_with: [PZL-031, PZL-044, PZL-054]
conflicts_with: []
tags: [bug, git, remote, deletion, consent, stabilization]
updated: 2026-08-07
---

# Explicit remote-deletion consent

## Core message

Local cleanup must not imply remote branch deletion; require a separate explicit
remote-deletion authorization and reconcile every canonical contract around it.

## Value

Restores a narrow, legible safety boundary for an irreversible provider mutation.

## Assumptions

- Local and remote cleanup are materially different authorities.
- Preserving remote branches by default is the intended product contract.

## Evidence and collisions

Workspace and AFK cleanup currently delete a matching origin branch under
`--apply --yes`. Specs and approved authority records forbid this, while one
runtime reference and tests now expect it. This is verified canonical drift.

## Trade-offs and recommendation

An extra flag or native Git step costs convenience but prevents consent expansion.
Stabilize before broad cleanup refactoring; the exact command owner remains open.

## Decision

Pending.
