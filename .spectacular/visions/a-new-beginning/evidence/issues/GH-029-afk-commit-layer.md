---
type: issue-evidence
source: source-011
github_issue: 29
url: https://github.com/alexsmedile/spectacular/issues/29
actionability: architecture-deferred
maps_to: [PZL-044, PZL-048, PZL-054, PZL-137]
retrieved: 2026-08-07
---

# GH-029 — AFK commit-layer gap

## Plain-language focus

AFK defines branches and PRs but leaves commit timing, scope, and session-end handling unspecified.

## Problem and evidence

The workflow's entire commit layer is compressed into “work and verification happen,” so
unattended branch history and goal linkage are incidental.

## Proposed direction

Define logical commit units, goal conformance, uncommitted session handoff, and the same
dry-run/apply/policy posture as other AFK mutations.

## Relationships and collisions

Overlaps closed/related #5 and #25 and should share a branch if retained. Broader sources
question whether AFK or Spectacular should own Git mutations at all.

## Actionability

Defer until execution authority and AFK survival are decided; otherwise document deleted behavior.

## Refactor relevance

Useful contract-gap evidence, not a reason to deepen a subsystem that may be extracted or removed.
