---
type: issue-evidence
source: source-011
github_issue: 30
url: https://github.com/alexsmedile/spectacular/issues/30
actionability: prototype-active
maps_to: [PZL-033, PZL-086, PZL-133, PZL-138]
retrieved: 2026-08-07
---

# GH-030 — Portfolio issue review

## Plain-language focus

Assess the whole backlog for currency, overlap, order, blockers, and handoff—not one issue at a time.

## Problem and evidence

Single-issue triage cannot find cross-issue schema duplication, shared file collisions, stale
claims, or a coherent dispatch order. Manual audits already found all four.

## Proposed direction

Fan out the existing read-only issue card, verify claims against current code, synthesize
clusters and branch collisions, extract human questions, then emit a reviewed handoff.

## Relationships and collisions

Consumes #19's grouping and #26's readability. The workflow may belong in an issue-triage
companion skill rather than Spectacular core.

## Actionability

Actively prototyped by Source 011 ingestion; measure usefulness before specifying a command.

## Refactor relevance

This session is direct usage evidence for a future reusable refactor/portfolio-review skill.
