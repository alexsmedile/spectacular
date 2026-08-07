---
type: concept-piece
id: PZL-071
status: captured
domain: autonomous-execution
sources: [source-006, source-007]
source_authority: proposal
assessment: disputed
evidence_status: partial
disposition: pending
depends_on: [PZL-067, PZL-069, PZL-070, PZL-072]
overlaps_with: [PZL-044, PZL-075, PZL-099]
conflicts_with: [PZL-051, PZL-098]
tags: [single-agent, run, checkpoint, branch, worktree]
updated: 2026-08-07
---

# Single-agent checkpointed run

## Core message

Begin with one active run, one coding agent, one isolated branch/worktree, narrow
checks, bounded retries, checkpoints, commits, and optional draft PR delivery.

## Value

Tests useful autonomy without queues, parallelism, leases, or background scheduling.

## Assumptions

- One coding agent may safely mutate under the locked mission boundary.
- Checkpoints are durable enough for reliable resume.

## Evidence and collisions

AFK and build workflows already cover parts, but current agent contracts reserve
general mutation and reconciliation for the orchestrator. CLI-owned agentic run
commands also conflict with PZL-051.
Source 007 proposes multiple specialized runs and agents inside one Mission, directly
conflicting with the single-agent MVP boundary.

## Trade-offs and recommendation

Focused autonomy versus a major authority reversal. Prototype as a host-runtime
recipe only after the run envelope and stop semantics are approved.

## Decision

Pending.
